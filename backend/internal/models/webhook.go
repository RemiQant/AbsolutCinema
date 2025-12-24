package models

import "time"

// XenditInvoiceCallback represents the webhook payload from Xendit for invoice callbacks
// Documentation: https://developers.xendit.co/api-reference/#invoice-callback
type XenditInvoiceCallback struct {
	// ID is the unique identifier for the invoice from Xendit
	ID string `json:"id"`

	// ExternalID is our booking ID (UUID) that we passed when creating the invoice
	ExternalID string `json:"external_id"`

	// UserID is the Xendit user ID (merchant)
	UserID string `json:"user_id"`

	// Status of the invoice: PAID, SETTLED, EXPIRED
	Status string `json:"status"`

	// MerchantName is the name of the merchant
	MerchantName string `json:"merchant_name"`

	// Amount is the invoice amount
	Amount float64 `json:"amount"`

	// PaidAmount is the amount that was actually paid
	PaidAmount float64 `json:"paid_amount,omitempty"`

	// BankCode is the bank used for payment (if applicable)
	BankCode string `json:"bank_code,omitempty"`

	// PaidAt is the timestamp when the invoice was paid
	PaidAt *time.Time `json:"paid_at,omitempty"`

	// PayerEmail is the email of the payer
	PayerEmail string `json:"payer_email,omitempty"`

	// Description of the invoice
	Description string `json:"description,omitempty"`

	// AdjustedReceivedAmount after fees
	AdjustedReceivedAmount float64 `json:"adjusted_received_amount,omitempty"`

	// FeesPaidAmount is the fees deducted
	FeesPaidAmount float64 `json:"fees_paid_amount,omitempty"`

	// Updated timestamp
	Updated time.Time `json:"updated,omitempty"`

	// Created timestamp
	Created time.Time `json:"created,omitempty"`

	// Currency of the invoice
	Currency string `json:"currency,omitempty"`

	// PaymentMethod used (e.g., "BANK_TRANSFER", "CREDIT_CARD", etc.)
	PaymentMethod string `json:"payment_method,omitempty"`

	// PaymentChannel specific channel used
	PaymentChannel string `json:"payment_channel,omitempty"`

	// PaymentDestination (e.g., virtual account number)
	PaymentDestination string `json:"payment_destination,omitempty"`
}

// Xendit Invoice Status Constants
const (
	XenditStatusPaid    = "PAID"
	XenditStatusSettled = "SETTLED"
	XenditStatusExpired = "EXPIRED"
)
