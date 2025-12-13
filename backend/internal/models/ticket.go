package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Ticket represents a transaction detail - a specific seat reservation
type Ticket struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingID  uuid.UUID `gorm:"type:uuid;not null;index" json:"booking_id"`
	ShowtimeID uint      `gorm:"not null;index" json:"showtime_id"`
	SeatNumber string    `gorm:"type:varchar(10);not null" json:"seat_number"` // e.g. A5
	
	// Relationships
	Booking  Booking  `gorm:"foreignKey:BookingID;constraint:OnDelete:CASCADE" json:"booking,omitempty"`
	Showtime Showtime `gorm:"foreignKey:ShowtimeID;constraint:OnDelete:CASCADE" json:"showtime,omitempty"`
}

// TableName specifies the table name for Ticket
func (Ticket) TableName() string {
	return "tickets"
}

// BeforeCreate hook - CRITICAL LOGIC: Ensure (showtime_id + seat_number) is unique
// This is enforced at the database level via a unique composite index
func (t *Ticket) BeforeCreate(tx *gorm.DB) error {
	// Check if this seat is already taken for this showtime
	var count int64
	err := tx.Model(&Ticket{}).
		Where("showtime_id = ? AND seat_number = ?", t.ShowtimeID, t.SeatNumber).
		Count(&count).Error
	
	if err != nil {
		return err
	}
	
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}
	
	return nil
}
