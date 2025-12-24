package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"absolutcinema-backend/internal/models"
	"absolutcinema-backend/internal/services"
)

// WebhookController handles external webhook callbacks
type WebhookController struct {
	bookingService *services.BookingService
}

// NewWebhookController creates a new webhook controller
func NewWebhookController(bookingService *services.BookingService) *WebhookController {
	return &WebhookController{
		bookingService: bookingService,
	}
}

// HandleXenditCallback handles Xendit invoice webhook callbacks
// @Router /api/webhooks/xendit [post]
func (wc *WebhookController) HandleXenditCallback(c *gin.Context) {
	// 1. Extract callback token from header
	// Xendit sends this token in the x-callback-token header
	callbackToken := c.GetHeader("x-callback-token")

	// 2. Bind JSON payload
	var payload models.XenditInvoiceCallback
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Printf("[Webhook] Failed to parse payload: %v", err)
		// Still return 200 to prevent Xendit from retrying malformed requests
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "invalid payload format",
		})
		return
	}

	// Log incoming webhook for debugging (exclude sensitive data in production)
	log.Printf("[Webhook] Received Xendit callback - ID: %s, ExternalID: %s, Status: %s",
		payload.ID, payload.ExternalID, payload.Status)

	// 3. Process the webhook via service
	err := wc.bookingService.HandleInvoiceCallback(&payload, callbackToken)
	if err != nil {
		// Check error type for appropriate response
		var webhookErr *services.WebhookError
		if errors.As(err, &webhookErr) {
			switch webhookErr.Code {
			case services.ErrCodeUnauthorized:
				// IMPORTANT: For security, we still return 200 to not reveal that the token is invalid
				// But we log it for monitoring
				log.Printf("[Webhook] SECURITY WARNING - Invalid callback token attempt for invoice %s", payload.ID)
				c.JSON(http.StatusOK, gin.H{
					"status":  "error",
					"message": "unauthorized",
				})
				return

			case services.ErrCodeBookingNotFound:
				log.Printf("[Webhook] Booking not found for external_id: %s", payload.ExternalID)
				// Return 200 to prevent retries for non-existent bookings
				c.JSON(http.StatusOK, gin.H{
					"status":  "error",
					"message": "booking not found",
				})
				return

			default:
				log.Printf("[Webhook] Error processing webhook: %v", err)
			}
		} else {
			// Unexpected error - log it
			log.Printf("[Webhook] Unexpected error: %v", err)
		}

		// For unexpected errors, we still return 200 but log for investigation
		// This prevents infinite retry loops from Xendit
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "processing error",
		})
		return
	}

	// 4. Success response
	log.Printf("[Webhook] Successfully processed callback for booking %s, status: %s",
		payload.ExternalID, payload.Status)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "webhook processed",
	})
}
