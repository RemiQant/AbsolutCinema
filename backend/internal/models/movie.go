package models

import (
	"gorm.io/gorm"
)

type Movie struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Title           string         `gorm:"type:varchar(255);not null" json:"title"`
	Description     string         `gorm:"type:text" json:"description"`
	DurationMinutes int            `gorm:"not null" json:"duration_minutes"`
	PosterURL       string         `gorm:"type:varchar(500);column:poster_url" json:"poster_url"`
	Rating          string         `gorm:"type:varchar(50)" json:"rating"`
	
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	Showtimes []Showtime `gorm:"foreignKey:MovieID" json:"-"`
}
