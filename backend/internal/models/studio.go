package models

import (
	"gorm.io/gorm"
)

// Studio represents a cinema studio/room with seating configuration
type Studio struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"` // e.g. Studio 1
	TotalRows int            `gorm:"not null" json:"total_rows"`             // e.g. 10 (A-J)
	TotalCols int            `gorm:"not null" json:"total_cols"`             // e.g. 8 (1-8)
	
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// Relationships
	Showtimes []Showtime `gorm:"foreignKey:StudioID" json:"showtimes,omitempty"`
}
