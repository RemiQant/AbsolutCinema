package models

import (
	"gorm.io/gorm"
)

type Studio struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	TotalRows int            `gorm:"not null" json:"total_rows"`       
	TotalCols int            `gorm:"not null" json:"total_cols"`
	
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	Showtimes []Showtime `gorm:"foreignKey:StudioID" json:"-"`
}
