package models

import (
	"time"

	"gorm.io/gorm"
)

type Showtime struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	MovieID   uint           `gorm:"not null;index" json:"movie_id"`
	StudioID  uint           `gorm:"not null;index:idx_studio_time" json:"studio_id"`
	StartTime time.Time      `gorm:"not null;index:idx_studio_time" json:"start_time"`
	EndTime   time.Time      `gorm:"not null;index" json:"end_time"`
	Price     float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	Movie   Movie    `gorm:"foreignKey:MovieID;constraint:OnDelete:CASCADE" json:"-"`
	Studio  Studio   `gorm:"foreignKey:StudioID;constraint:OnDelete:CASCADE" json:"-"`
	Tickets []Ticket `gorm:"foreignKey:ShowtimeID" json:"-"`
}
