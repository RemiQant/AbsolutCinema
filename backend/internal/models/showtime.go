package models

import (
	"time"

	"gorm.io/gorm"
)

// Showtime represents a scheduled showing of a movie in a studio
type Showtime struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	MovieID   uint           `gorm:"not null;index" json:"movie_id"`
	StudioID  uint           `gorm:"not null;index:idx_studio_time" json:"studio_id"`
	StartTime time.Time      `gorm:"not null;index:idx_studio_time" json:"start_time"`
	EndTime   time.Time      `gorm:"not null;index" json:"end_time"` // Computed: StartTime + Movie.Duration + Buffer
	Price     float64        `gorm:"type:decimal(10,2);not null" json:"price"` // Fixed price for this specific slot
	
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// Relationships
	Movie   Movie    `gorm:"foreignKey:MovieID;constraint:OnDelete:CASCADE" json:"movie,omitempty"`
	Studio  Studio   `gorm:"foreignKey:StudioID;constraint:OnDelete:CASCADE" json:"studio,omitempty"`
	Tickets []Ticket `gorm:"foreignKey:ShowtimeID" json:"tickets,omitempty"`
}
