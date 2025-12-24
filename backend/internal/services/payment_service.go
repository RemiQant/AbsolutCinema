package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	xendit "github.com/xendit/xendit-go/v6"
	invoice "github.com/xendit/xendit-go/v6/invoice"

	"absolutcinema-backend/internal/models"
)

// PaymentService handles payment operations via Xendit
type PaymentService struct {
	client *xendit.APIClient
}

// InvoiceResult contains the result of creating an invoice
type InvoiceResult struct {
	InvoiceID   string  `json:"invoice_id"`
	InvoiceURL  string  `json:"invoice_url"`
	ExternalID  string  `json:"external_id"`
	Amount      float64 `json:"amount"`
	Status      string  `json:"status"`
	ExpiryDate  string  `json:"expiry_date"`
}

// NewPaymentService creates a new payment service with Xendit client
func NewPaymentService() (*PaymentService, error) {
	secretKey := os.Getenv("XENDIT_SECRET_KEY")
	if secretKey == "" {
		return nil, errors.New("XENDIT_SECRET_KEY environment variable is not set")
	}

	client := xendit.NewClient(secretKey)

	return &PaymentService{
		client: client,
	}, nil
}

// CreateInvoice creates a Xendit invoice for a booking
func (ps *PaymentService) CreateInvoice(booking *models.Booking, userEmail string) (*InvoiceResult, error) {
	if ps.client == nil {
		return nil, errors.New("payment client is not initialized")
	}

	ctx := context.Background()

	// External ID should be the booking ID (UUID)
	externalID := booking.ID.String()

	// Invoice description
	description := fmt.Sprintf("AbsolutCinema Booking - Invoice %s", booking.InvoiceNumber)

	// Set invoice expiry (24 hours from now)
	expiryDate := time.Now().Add(24 * time.Hour)

	// Build invoice items from tickets
	items := buildInvoiceItems(booking)

	// Create invoice request
	invoiceRequest := *invoice.NewCreateInvoiceRequest(externalID, float64(booking.TotalAmount))
	invoiceRequest.SetDescription(description)
	invoiceRequest.SetPayerEmail(userEmail)
	invoiceRequest.SetCurrency("IDR")
	invoiceRequest.SetSuccessRedirectUrl(getSuccessRedirectURL())
	invoiceRequest.SetFailureRedirectUrl(getFailureRedirectURL())
	
	if len(items) > 0 {
		invoiceRequest.SetItems(items)
	}

	// Create the invoice via Xendit API
	resp, _, err := ps.client.InvoiceApi.CreateInvoice(ctx).
		CreateInvoiceRequest(invoiceRequest).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("failed to create Xendit invoice: %w", err)
	}

	result := &InvoiceResult{
		InvoiceID:   *resp.Id,
		InvoiceURL:  resp.InvoiceUrl,
		ExternalID:  resp.ExternalId,
		Amount:      resp.Amount,
		Status:      string(resp.Status),
		ExpiryDate:  expiryDate.Format(time.RFC3339),
	}

	return result, nil
}

// GetInvoice retrieves an invoice by ID
func (ps *PaymentService) GetInvoice(invoiceID string) (*InvoiceResult, error) {
	if ps.client == nil {
		return nil, errors.New("payment client is not initialized")
	}

	ctx := context.Background()

	resp, _, err := ps.client.InvoiceApi.GetInvoiceById(ctx, invoiceID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get Xendit invoice: %w", err)
	}

	result := &InvoiceResult{
		InvoiceID:   *resp.Id,
		InvoiceURL:  resp.InvoiceUrl,
		ExternalID:  resp.ExternalId,
		Amount:      resp.Amount,
		Status:      string(resp.Status),
	}

	return result, nil
}

// GetInvoiceByExternalID retrieves an invoice by external ID (booking ID)
func (ps *PaymentService) GetInvoiceByExternalID(bookingID uuid.UUID) (*InvoiceResult, error) {
	if ps.client == nil {
		return nil, errors.New("payment client is not initialized")
	}

	ctx := context.Background()

	resp, _, err := ps.client.InvoiceApi.GetInvoices(ctx).
		ExternalId(bookingID.String()).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get Xendit invoice by external ID: %w", err)
	}

	if len(resp) == 0 {
		return nil, errors.New("invoice not found")
	}

	// Get the most recent invoice
	inv := resp[0]

	result := &InvoiceResult{
		InvoiceID:   *inv.Id,
		InvoiceURL:  inv.InvoiceUrl,
		ExternalID:  inv.ExternalId,
		Amount:      inv.Amount,
		Status:      string(inv.Status),
	}

	return result, nil
}

// ExpireInvoice expires/cancels an invoice
func (ps *PaymentService) ExpireInvoice(invoiceID string) error {
	if ps.client == nil {
		return errors.New("payment client is not initialized")
	}

	ctx := context.Background()

	_, _, err := ps.client.InvoiceApi.ExpireInvoice(ctx, invoiceID).Execute()
	if err != nil {
		return fmt.Errorf("failed to expire Xendit invoice: %w", err)
	}

	return nil
}

// buildInvoiceItems creates invoice items from booking tickets
func buildInvoiceItems(booking *models.Booking) []invoice.InvoiceItem {
	if len(booking.Tickets) == 0 {
		return nil
	}

	items := make([]invoice.InvoiceItem, 0, len(booking.Tickets))

	// Group tickets by showtime for better invoice display
	// For simplicity, we create one item per ticket
	pricePerSeat := float32(booking.TotalAmount / float64(len(booking.Tickets)))
	
	for _, ticket := range booking.Tickets {
		item := *invoice.NewInvoiceItem(
			fmt.Sprintf("Seat %s", ticket.SeatNumber),
			1, // Quantity
			pricePerSeat,
		)
		items = append(items, item)
	}

	return items
}

// getSuccessRedirectURL returns the URL to redirect after successful payment
func getSuccessRedirectURL() string {
	baseURL := os.Getenv("FRONTEND_URL")
	if baseURL == "" {
		baseURL = "http://localhost:5173"
	}
	return baseURL + "/booking/success"
}

// getFailureRedirectURL returns the URL to redirect after failed payment
func getFailureRedirectURL() string {
	baseURL := os.Getenv("FRONTEND_URL")
	if baseURL == "" {
		baseURL = "http://localhost:5173"
	}
	return baseURL + "/booking/failed"
}

// WebhookError represents errors that can occur during webhook processing
type WebhookError struct {
	Code    string
	Message string
}

func (e *WebhookError) Error() string {
	return e.Message
}

// Webhook error codes
const (
	ErrCodeUnauthorized     = "UNAUTHORIZED"
	ErrCodeBookingNotFound  = "BOOKING_NOT_FOUND"
	ErrCodeInvalidStatus    = "INVALID_STATUS"
	ErrCodeAlreadyProcessed = "ALREADY_PROCESSED"
)

// NewWebhookError creates a new webhook error
func NewWebhookError(code, message string) *WebhookError {
	return &WebhookError{Code: code, Message: message}
}

// ValidateCallbackToken validates the Xendit callback token from the request header
// This is critical for security - ensures the webhook is actually from Xendit
func ValidateCallbackToken(token string) error {
	expectedToken := os.Getenv("XENDIT_WEBHOOK_TOKEN")
	if expectedToken == "" {
		return NewWebhookError(ErrCodeUnauthorized, "webhook token not configured on server")
	}

	if token == "" {
		return NewWebhookError(ErrCodeUnauthorized, "missing callback token in request")
	}

	if token != expectedToken {
		return NewWebhookError(ErrCodeUnauthorized, "invalid callback token")
	}

	return nil
}
