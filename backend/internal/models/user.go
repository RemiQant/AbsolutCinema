package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system (admin or customer)
type User struct {
	ID       uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Username string         `gorm:"type:varchar(255);not null" json:"username"`
	Email    string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password string         `gorm:"type:varchar(255);not null" json:"-"` // "-" prevents password from being serialized in JSON
	Role     string         `gorm:"type:varchar(50);default:'customer'" json:"role"` // Values: admin, customer
	
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete support
	
	// Relationships
	Bookings []Booking `gorm:"foreignKey:UserID" json:"bookings,omitempty"`
}

// BeforeCreate hook to generate UUID if not set
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
