package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"absolutcinema-backend/internal/services"
)

// BookingController handles booking-related HTTP requests
type BookingController struct {
	bookingService *services.BookingService
}

// NewBookingController creates a new booking controller
func NewBookingController(bookingService *services.BookingService) *BookingController {
	return &BookingController{
		bookingService: bookingService,
	}
}

// CreateBooking handles POST /api/bookings
// Creates a new booking with seat reservation and payment link generation
func (bc *BookingController) CreateBooking(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		// Try parsing as string
		userIDStr, ok := userIDValue.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user ID format",
			})
			return
		}
		var err error
		userID, err = uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user ID",
			})
			return
		}
	}

	// Parse request body
	var req services.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate minimum seat count
	if len(req.SeatNumbers) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "At least one seat must be selected",
		})
		return
	}

	// Create booking
	result, err := bc.bookingService.CreateBooking(userID, &req)
	if err != nil {
		// Check for seat conflict (race condition)
		if services.IsConflictError(err) {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Seat conflict",
				"details": err.Error(),
				"code":    "SEAT_ALREADY_TAKEN",
			})
			return
		}

		// Check for validation errors
		if err.Error() == "showtime not found" ||
			err.Error() == "cannot book seats for a showtime that has already started" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Check for invalid seat errors
		if len(err.Error()) > 0 && (err.Error()[:4] == "seat" || err.Error()[:7] == "invalid") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create booking",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     result.Message,
		"data":        result.Booking,
		"payment_url": result.PaymentURL,
	})
}

// GetBookings handles GET /api/bookings
// Returns all bookings for the authenticated user
func (bc *BookingController) GetBookings(c *gin.Context) {
	// Get user ID from context
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		userIDStr, ok := userIDValue.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user ID format",
			})
			return
		}
		var err error
		userID, err = uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user ID",
			})
			return
		}
	}

	// Get bookings
	bookings, err := bc.bookingService.GetUserBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve bookings",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bookings retrieved successfully",
		"data":    bookings,
		"count":   len(bookings),
	})
}

// GetBookingByID handles GET /api/bookings/:id
// Returns a specific booking by ID
func (bc *BookingController) GetBookingByID(c *gin.Context) {
	// Get user ID from context
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		userIDStr, ok := userIDValue.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user ID format",
			})
			return
		}
		var err error
		userID, err = uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user ID",
			})
			return
		}
	}

	// Parse booking ID
	bookingIDStr := c.Param("id")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid booking ID",
		})
		return
	}

	// Get booking
	booking, err := bc.bookingService.GetBookingByID(bookingID, userID)
	if err != nil {
		if err.Error() == "booking not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve booking",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking retrieved successfully",
		"data":    booking,
	})
}

// CancelBooking handles DELETE /api/bookings/:id
// Cancels a pending booking
func (bc *BookingController) CancelBooking(c *gin.Context) {
	// Get user ID from context
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		userIDStr, ok := userIDValue.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user ID format",
			})
			return
		}
		var err error
		userID, err = uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user ID",
			})
			return
		}
	}

	// Parse booking ID
	bookingIDStr := c.Param("id")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid booking ID",
		})
		return
	}

	// Cancel booking
	if err := bc.bookingService.CancelBooking(bookingID, userID); err != nil {
		if err.Error() == "booking not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err.Error() == "booking is already cancelled" ||
			err.Error() == "cannot cancel a paid booking" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to cancel booking",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking cancelled successfully",
	})
}

// RetryPayment handles POST /api/bookings/:id/retry-payment
// Generates a new payment link for a pending booking
func (bc *BookingController) RetryPayment(c *gin.Context) {
	// Get user ID from context
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		userIDStr, ok := userIDValue.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user ID format",
			})
			return
		}
		var err error
		userID, err = uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user ID",
			})
			return
		}
	}

	// Parse booking ID
	bookingIDStr := c.Param("id")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid booking ID",
		})
		return
	}

	// Retry payment
	result, err := bc.bookingService.RetryPayment(bookingID, userID)
	if err != nil {
		if err.Error() == "booking not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err.Error() == "can only retry payment for pending bookings" ||
			err.Error() == "payment service is not available" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retry payment",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     result.Message,
		"data":        result.Booking,
		"payment_url": result.PaymentURL,
	})
}
