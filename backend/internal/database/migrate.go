package database

import (
	"absolutcinema-backend/internal/models"
	"log"
)

// Migrate runs database migrations for all models
func (s *service) Migrate() error {
	log.Println("Running database migrations...")
	
	// AutoMigrate will create tables, missing columns, and missing indexes
	// It will NOT change existing column types or delete unused columns
	err := s.gormDB.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.Movie{},
		&models.Studio{},
		&models.Showtime{},
		&models.Booking{},
		&models.Ticket{},
	)
	
	if err != nil {
		log.Printf("Migration failed: %v", err)
		return err
	}
	
	// Create unique composite index for tickets (showtime_id, seat_number)
	// This enforces that a seat can only be booked once per showtime
	err = s.gormDB.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_tickets_showtime_seat 
		ON tickets(showtime_id, seat_number)
	`).Error
	
	if err != nil {
		log.Printf("Failed to create unique index on tickets: %v", err)
		return err
	}
	
	log.Println("Database migrations completed successfully!")
	return nil
}
