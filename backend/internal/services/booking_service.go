package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"absolutcinema-backend/internal/models"
)

const (
	// PostgreSQL error code for unique violation
	PgUniqueViolationCode = "23505"

	// Booking status constants
	BookingStatusPending   = "PENDING"
	BookingStatusPaid      = "PAID"
	BookingStatusCancelled = "CANCELLED"
	BookingStatusExpired   = "EXPIRED"
)

// BookingService handles booking operations
type BookingService struct {
	db             *gorm.DB
	paymentService *PaymentService
}

// CreateBookingRequest represents the request to create a booking
type CreateBookingRequest struct {
	ShowtimeID  FlexibleUint `json:"showtime_id" binding:"required"`
	SeatNumbers []string     `json:"seat_numbers" binding:"required,min=1"`
}

// FlexibleUint is a uint that can be unmarshaled from both string and number JSON values
type FlexibleUint uint

func (f *FlexibleUint) UnmarshalJSON(data []byte) error {
	// Remove quotes if present (string value)
	s := strings.Trim(string(data), "\"")
	
	// Parse as uint
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}
	
	*f = FlexibleUint(val)
	return nil
}

// BookingResult represents the result of a booking creation
type BookingResult struct {
	Booking    *models.Booking `json:"booking"`
	PaymentURL string          `json:"payment_url"`
	Message    string          `json:"message"`
}

// NewBookingService creates a new booking service
func NewBookingService(db *gorm.DB, paymentService *PaymentService) *BookingService {
	return &BookingService{
		db:             db,
		paymentService: paymentService,
	}
}

// CreateBooking creates a new booking with atomic transaction and race condition handling
func (bs *BookingService) CreateBooking(userID uuid.UUID, req *CreateBookingRequest) (*BookingResult, error) {
	// 1. Validate showtime exists and get price
	var showtime models.Showtime
	showtimeID := uint(req.ShowtimeID)
	if err := bs.db.Preload("Studio").First(&showtime, showtimeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("showtime not found")
		}
		return nil, fmt.Errorf("failed to fetch showtime: %w", err)
	}

	// Check if showtime is in the past
	if showtime.StartTime.Before(time.Now()) {
		return nil, errors.New("cannot book seats for a showtime that has already started")
	}

	// 2. Validate seat numbers against studio dimensions
	if err := bs.validateSeatNumbers(req.SeatNumbers, showtime.Studio.TotalRows, showtime.Studio.TotalCols); err != nil {
		return nil, err
	}

	// 3. Remove duplicates from seat numbers
	uniqueSeats := removeDuplicateSeats(req.SeatNumbers)

	// 4. Calculate total amount
	totalAmount := showtime.Price * float64(len(uniqueSeats))

	// 5. Generate invoice number
	invoiceNumber := generateInvoiceNumber()

	// 6. Get user email for payment invoice
	var user models.User
	if err := bs.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// 7. Create booking in a transaction
	var booking models.Booking
	err := bs.db.Transaction(func(tx *gorm.DB) error {
		// Create booking record
		booking = models.Booking{
			UserID:        userID,
			InvoiceNumber: invoiceNumber,
			TotalAmount:   totalAmount,
			Status:        BookingStatusPending,
		}

		if err := tx.Create(&booking).Error; err != nil {
			return fmt.Errorf("failed to create booking: %w", err)
		}

		// Create ticket records for each seat
		for _, seatNumber := range uniqueSeats {
			ticket := models.Ticket{
				BookingID:  booking.ID,
				ShowtimeID: showtimeID,
				SeatNumber: strings.ToUpper(seatNumber), // Normalize to uppercase
			}

			if err := tx.Create(&ticket).Error; err != nil {
				// Check for PostgreSQL unique violation (race condition guard)
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) && pgErr.Code == PgUniqueViolationCode {
					return &SeatConflictError{SeatNumber: seatNumber}
				}

				// Also check for GORM's duplicate key error (from BeforeCreate hook)
				if errors.Is(err, gorm.ErrDuplicatedKey) {
					return &SeatConflictError{SeatNumber: seatNumber}
				}

				return fmt.Errorf("failed to create ticket for seat %s: %w", seatNumber, err)
			}
		}

		return nil
	})

	if err != nil {
		// Check if it's a seat conflict error
		var conflictErr *SeatConflictError
		if errors.As(err, &conflictErr) {
			return nil, conflictErr
		}
		return nil, err
	}

	// 8. Load the complete booking with tickets
	if err := bs.db.Preload("Tickets").First(&booking, "id = ?", booking.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load booking: %w", err)
	}

	// 9. Create payment invoice via Xendit
	result := &BookingResult{
		Booking: &booking,
		Message: "Booking created successfully",
	}

	if bs.paymentService != nil {
		invoiceResult, err := bs.paymentService.CreateInvoice(&booking, user.Email)
		if err != nil {
			// Payment creation failed, but booking is created with PENDING status
			// User can retry payment later
			result.Message = "Booking created but payment link generation failed. Please try again later."
			return result, nil
		}

		// Update booking with payment URL and payment ID
		updateResult := bs.db.Model(&models.Booking{}).Where("id = ?", booking.ID).Updates(map[string]interface{}{
			"payment_url": invoiceResult.InvoiceURL,
			"payment_id":  invoiceResult.InvoiceID,
		})
		if updateResult.Error != nil {
			result.Message = fmt.Sprintf("Booking created but failed to save payment URL: %v", updateResult.Error)
			result.PaymentURL = invoiceResult.InvoiceURL
			return result, nil
		}

		result.PaymentURL = invoiceResult.InvoiceURL
		result.Booking.PaymentURL = invoiceResult.InvoiceURL
	}

	return result, nil
}

// GetUserBookings retrieves all bookings for a user
func (bs *BookingService) GetUserBookings(userID uuid.UUID) ([]models.Booking, error) {
	var bookings []models.Booking
	err := bs.db.
		Preload("Tickets").
		Preload("Tickets.Showtime").
		Preload("Tickets.Showtime.Movie").
		Preload("Tickets.Showtime.Studio").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&bookings).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch bookings: %w", err)
	}

	return bookings, nil
}

// GetBookingByID retrieves a booking by ID
func (bs *BookingService) GetBookingByID(bookingID uuid.UUID, userID uuid.UUID) (*models.Booking, error) {
	var booking models.Booking
	err := bs.db.
		Preload("Tickets").
		Preload("Tickets.Showtime").
		Preload("Tickets.Showtime.Movie").
		Preload("Tickets.Showtime.Studio").
		Where("id = ? AND user_id = ?", bookingID, userID).
		First(&booking).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("booking not found")
		}
		return nil, fmt.Errorf("failed to fetch booking: %w", err)
	}

	return &booking, nil
}

// GetOccupiedSeats retrieves all occupied seats for a showtime
func (bs *BookingService) GetOccupiedSeats(showtimeID uint) ([]string, error) {
	var tickets []models.Ticket
	err := bs.db.
		Joins("JOIN bookings ON bookings.id = tickets.booking_id").
		Where("tickets.showtime_id = ?", showtimeID).
		Where("bookings.status != ?", BookingStatusCancelled).
		Find(&tickets).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch occupied seats: %w", err)
	}

	seats := make([]string, len(tickets))
	for i, ticket := range tickets {
		seats[i] = ticket.SeatNumber
	}

	return seats, nil
}

// UpdateBookingStatus updates the status of a booking
func (bs *BookingService) UpdateBookingStatus(bookingID uuid.UUID, status string) error {
	result := bs.db.Model(&models.Booking{}).
		Where("id = ?", bookingID).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to update booking status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("booking not found")
	}

	return nil
}

// CancelBooking cancels a booking and releases the seats
func (bs *BookingService) CancelBooking(bookingID uuid.UUID, userID uuid.UUID) error {
	return bs.db.Transaction(func(tx *gorm.DB) error {
		// Find the booking
		var booking models.Booking
		err := tx.Where("id = ? AND user_id = ?", bookingID, userID).First(&booking).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("booking not found")
			}
			return err
		}

		// Check if already cancelled
		if booking.Status == BookingStatusCancelled {
			return errors.New("booking is already cancelled")
		}

		// Check if already paid
		if booking.Status == BookingStatusPaid {
			return errors.New("cannot cancel a paid booking")
		}

		// Delete associated tickets to release seats
		if err := tx.Where("booking_id = ?", bookingID).Delete(&models.Ticket{}).Error; err != nil {
			return fmt.Errorf("failed to delete tickets: %w", err)
		}

		// Update booking status
		if err := tx.Model(&booking).Update("status", BookingStatusCancelled).Error; err != nil {
			return fmt.Errorf("failed to update booking status: %w", err)
		}

		return nil
	})
}

// RetryPayment retries payment for a pending booking
func (bs *BookingService) RetryPayment(bookingID uuid.UUID, userID uuid.UUID) (*BookingResult, error) {
	// Get booking
	booking, err := bs.GetBookingByID(bookingID, userID)
	if err != nil {
		return nil, err
	}

	// Check if booking is pending
	if booking.Status != BookingStatusPending {
		return nil, errors.New("can only retry payment for pending bookings")
	}

	// Get user
	var user models.User
	if err := bs.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// Check if there's already a valid payment URL
	if booking.PaymentURL != "" {
		// Try to get invoice status
		if bs.paymentService != nil {
			invoiceResult, err := bs.paymentService.GetInvoiceByExternalID(bookingID)
			if err == nil && invoiceResult.Status == "PENDING" {
				// Invoice still valid
				return &BookingResult{
					Booking:    booking,
					PaymentURL: booking.PaymentURL,
					Message:    "Existing payment link is still valid",
				}, nil
			}
		}
	}

	// Create new invoice
	if bs.paymentService == nil {
		return nil, errors.New("payment service is not available")
	}

	invoiceResult, err := bs.paymentService.CreateInvoice(booking, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment invoice: %w", err)
	}

	// Update booking with new payment URL
	if err := bs.db.Model(booking).Update("payment_url", invoiceResult.InvoiceURL).Error; err != nil {
		return nil, fmt.Errorf("failed to save payment URL: %w", err)
	}

	booking.PaymentURL = invoiceResult.InvoiceURL

	return &BookingResult{
		Booking:    booking,
		PaymentURL: invoiceResult.InvoiceURL,
		Message:    "New payment link generated successfully",
	}, nil
}

// validateSeatNumbers validates that all seat numbers are within studio dimensions
func (bs *BookingService) validateSeatNumbers(seatNumbers []string, totalRows, totalCols int) error {
	for _, seat := range seatNumbers {
		if err := validateSeatNumber(seat, totalRows, totalCols); err != nil {
			return err
		}
	}
	return nil
}

// validateSeatNumber validates a single seat number against studio dimensions
func validateSeatNumber(seatNumber string, totalRows, totalCols int) error {
	seatNumber = strings.ToUpper(strings.TrimSpace(seatNumber))
	
	if len(seatNumber) < 2 {
		return fmt.Errorf("invalid seat number format: %s", seatNumber)
	}

	// Parse row (letter) and column (number)
	row := seatNumber[0]
	colStr := seatNumber[1:]

	// Validate row (A-Z maps to 0-25)
	if row < 'A' || row > 'Z' {
		return fmt.Errorf("invalid seat row: %s", seatNumber)
	}

	rowIndex := int(row - 'A')
	if rowIndex >= totalRows {
		return fmt.Errorf("seat row %c exceeds studio capacity (max row: %c)", row, rune('A'+totalRows-1))
	}

	// Validate column
	col, err := strconv.Atoi(colStr)
	if err != nil {
		return fmt.Errorf("invalid seat column: %s", seatNumber)
	}

	if col < 1 || col > totalCols {
		return fmt.Errorf("seat column %d exceeds studio capacity (1-%d)", col, totalCols)
	}

	return nil
}

// generateInvoiceNumber generates a unique invoice number
func generateInvoiceNumber() string {
	now := time.Now()
	return fmt.Sprintf("INV-%d%02d%02d-%s",
		now.Year(),
		now.Month(),
		now.Day(),
		strings.ToUpper(uuid.New().String()[:8]),
	)
}

// removeDuplicateSeats removes duplicate seat numbers
func removeDuplicateSeats(seats []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(seats))

	for _, seat := range seats {
		normalized := strings.ToUpper(strings.TrimSpace(seat))
		if !seen[normalized] {
			seen[normalized] = true
			result = append(result, normalized)
		}
	}

	return result
}

// SeatConflictError represents a seat already taken error
type SeatConflictError struct {
	SeatNumber string
}

func (e *SeatConflictError) Error() string {
	return fmt.Sprintf("seat %s is already taken", e.SeatNumber)
}

// IsConflictError checks if an error is a seat conflict error
func IsConflictError(err error) bool {
	var conflictErr *SeatConflictError
	return errors.As(err, &conflictErr)
}

// HandleInvoiceCallback processes Xendit invoice webhook callbacks
// This is idempotent - calling it multiple times with the same payload has no additional effect
func (bs *BookingService) HandleInvoiceCallback(payload *models.XenditInvoiceCallback, callbackToken string) error {
	// 1. Security: Validate callback token first
	if err := ValidateCallbackToken(callbackToken); err != nil {
		return err
	}

	// 2. Parse external_id (booking UUID)
	bookingID, err := uuid.Parse(payload.ExternalID)
	if err != nil {
		return NewWebhookError(ErrCodeBookingNotFound, "invalid external_id format")
	}

	// 3. Find the booking
	var booking models.Booking
	if err := bs.db.First(&booking, "id = ?", bookingID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewWebhookError(ErrCodeBookingNotFound, "booking not found")
		}
		return fmt.Errorf("failed to fetch booking: %w", err)
	}

	// 4. Process based on Xendit status
	switch payload.Status {
	case models.XenditStatusPaid, models.XenditStatusSettled:
		return bs.handlePaymentSuccess(&booking, payload)

	case models.XenditStatusExpired:
		return bs.handlePaymentExpired(&booking, payload)

	default:
		// Unknown status - log but don't error (Xendit may send other statuses we don't care about)
		return nil
	}
}

// handlePaymentSuccess handles successful payment (PAID or SETTLED status)
func (bs *BookingService) handlePaymentSuccess(booking *models.Booking, payload *models.XenditInvoiceCallback) error {
	// Idempotency check: if already paid, skip processing
	if booking.Status == BookingStatusPaid {
		// Already processed - return nil (success) for idempotency
		return nil
	}

	// Update booking status to PAID
	result := bs.db.Model(booking).Updates(map[string]interface{}{
		"status":     BookingStatusPaid,
		"payment_id": payload.ID, // Store Xendit invoice ID for reference
	})

	if result.Error != nil {
		return fmt.Errorf("failed to update booking status to PAID: %w", result.Error)
	}

	// TODO: Trigger ticket generation or send confirmation email here
	// Example:
	// go bs.sendBookingConfirmationEmail(booking.ID)
	// go bs.generateTicketPDF(booking.ID)

	return nil
}

// handlePaymentExpired handles expired payment (EXPIRED status)
func (bs *BookingService) handlePaymentExpired(booking *models.Booking, payload *models.XenditInvoiceCallback) error {
	// Idempotency check: if already cancelled or expired, skip processing
	if booking.Status == BookingStatusCancelled || booking.Status == BookingStatusExpired {
		return nil
	}

	// Don't expire if already paid (edge case where PAID and EXPIRED come out of order)
	if booking.Status == BookingStatusPaid {
		return nil
	}

	// Update booking status to CANCELLED and release seats
	return bs.db.Transaction(func(tx *gorm.DB) error {
		// Delete tickets to release seats
		if err := tx.Where("booking_id = ?", booking.ID).Delete(&models.Ticket{}).Error; err != nil {
			return fmt.Errorf("failed to delete tickets: %w", err)
		}

		// Update booking status
		if err := tx.Model(booking).Update("status", BookingStatusCancelled).Error; err != nil {
			return fmt.Errorf("failed to update booking status to CANCELLED: %w", err)
		}

		return nil
	})
}
