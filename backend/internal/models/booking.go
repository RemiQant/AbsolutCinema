package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Booking struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	
	InvoiceNumber string    `gorm:"type:varchar(100);uniqueIndex" json:"invoice_number"`
	TotalAmount   float64   `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status        string    `gorm:"type:varchar(50);default:'PENDING'" json:"status"`
	
	// Payment gateway fields (Xendit)
	PaymentURL string `gorm:"type:varchar(500);column:payment_url" json:"payment_url,omitempty"`
	PaymentID  string `gorm:"type:varchar(100);column:payment_id" json:"payment_id,omitempty"`
	
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	
	// Relationships (User hidden from JSON to reduce payload size)
	User    User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Tickets []Ticket `gorm:"foreignKey:BookingID" json:"tickets,omitempty"`
}

// BeforeCreate hook to generate UUID if not set
func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
