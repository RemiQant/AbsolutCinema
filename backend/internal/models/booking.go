package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Booking represents a booking/transaction header (receipt)
type Booking struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	
	InvoiceNumber string    `gorm:"type:varchar(100);uniqueIndex" json:"invoice_number"` // e.g. INV-2025-XXXX
	TotalAmount   float64   `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status        string    `gorm:"type:varchar(50);default:'PENDING'" json:"status"` // PENDING, PAID, CANCELLED
	
	// Phase 2 Prep: Store the Midtrans/Gateway URL here
	PaymentURL string `gorm:"type:varchar(500);column:payment_url" json:"payment_url,omitempty"`
	
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	
	// Relationships
	User    User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Tickets []Ticket `gorm:"foreignKey:BookingID" json:"tickets,omitempty"`
}

// BeforeCreate hook to generate UUID if not set
func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
