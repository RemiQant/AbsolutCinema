package services

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"absolutcinema-backend/internal/models"
)

const (
	// CleanupBufferMinutes is the buffer time between showtimes for cleaning/preparation
	CleanupBufferMinutes = 15
)

type ShowtimeService struct {
	db *gorm.DB
}

func NewShowtimeService(db *gorm.DB) *ShowtimeService {
	return &ShowtimeService{db: db}
}

// CreateShowtime creates a new showtime with overlap validation
func (s *ShowtimeService) CreateShowtime(showtime *models.Showtime) error {
	// Validation
	if err := s.validateShowtime(showtime); err != nil {
		return err
	}

	// Fetch movie to get duration
	var movie models.Movie
	if err := s.db.First(&movie, showtime.MovieID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("movie not found")
		}
		return err
	}

	// Verify studio exists
	var studio models.Studio
	if err := s.db.First(&studio, showtime.StudioID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("studio not found")
		}
		return err
	}

	// Calculate EndTime: StartTime + Movie Duration + Cleanup Buffer
	totalDuration := time.Duration(movie.DurationMinutes+CleanupBufferMinutes) * time.Minute
	showtime.EndTime = showtime.StartTime.Add(totalDuration)

	// CRITICAL: Check for overlapping showtimes in the same studio
	if err := s.checkOverlap(0, showtime.StudioID, showtime.StartTime, showtime.EndTime); err != nil {
		return err
	}

	return s.db.Create(showtime).Error
}

// GetAllShowtimes retrieves all showtimes with optional filters
func (s *ShowtimeService) GetAllShowtimes(movieID *uint, date *time.Time) ([]models.Showtime, error) {
	var showtimes []models.Showtime
	query := s.db.Preload("Movie").Preload("Studio")

	// Filter by movie_id if provided
	if movieID != nil && *movieID > 0 {
		query = query.Where("movie_id = ?", *movieID)
	}

	// Filter by date if provided (showtimes starting on that date)
	if date != nil {
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)
		query = query.Where("start_time >= ? AND start_time < ?", startOfDay, endOfDay)
	}

	if err := query.Order("start_time ASC").Find(&showtimes).Error; err != nil {
		return nil, err
	}

	return showtimes, nil
}

// GetShowtimeByID retrieves a showtime by ID
func (s *ShowtimeService) GetShowtimeByID(id uint) (*models.Showtime, error) {
	var showtime models.Showtime
	if err := s.db.Preload("Movie").Preload("Studio").First(&showtime, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("showtime not found")
		}
		return nil, err
	}
	return &showtime, nil
}

// UpdateShowtime updates an existing showtime with overlap validation
func (s *ShowtimeService) UpdateShowtime(id uint, updates *models.Showtime) error {
	// Validation
	if err := s.validateShowtime(updates); err != nil {
		return err
	}

	// Check if showtime exists
	var showtime models.Showtime
	if err := s.db.First(&showtime, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("showtime not found")
		}
		return err
	}

	// Fetch movie to get duration
	var movie models.Movie
	if err := s.db.First(&movie, updates.MovieID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("movie not found")
		}
		return err
	}

	// Verify studio exists
	var studio models.Studio
	if err := s.db.First(&studio, updates.StudioID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("studio not found")
		}
		return err
	}

	// Calculate new EndTime
	totalDuration := time.Duration(movie.DurationMinutes+CleanupBufferMinutes) * time.Minute
	newEndTime := updates.StartTime.Add(totalDuration)

	// CRITICAL: Check for overlapping showtimes (excluding current showtime)
	if err := s.checkOverlap(id, updates.StudioID, updates.StartTime, newEndTime); err != nil {
		return err
	}

	// Update fields
	showtime.MovieID = updates.MovieID
	showtime.StudioID = updates.StudioID
	showtime.StartTime = updates.StartTime
	showtime.EndTime = newEndTime
	showtime.Price = updates.Price

	return s.db.Save(&showtime).Error
}

// DeleteShowtime soft deletes a showtime
func (s *ShowtimeService) DeleteShowtime(id uint) error {
	result := s.db.Delete(&models.Showtime{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("showtime not found")
	}
	return nil
}

// validateShowtime validates showtime data
func (s *ShowtimeService) validateShowtime(showtime *models.Showtime) error {
	if showtime.MovieID == 0 {
		return errors.New("movie_id is required")
	}

	if showtime.StudioID == 0 {
		return errors.New("studio_id is required")
	}

	if showtime.StartTime.IsZero() {
		return errors.New("start_time is required")
	}

	// Ensure start time is in the future (optional business rule)
	if showtime.StartTime.Before(time.Now()) {
		return errors.New("start_time must be in the future")
	}

	if showtime.Price <= 0 {
		return errors.New("price must be greater than 0")
	}

	return nil
}

// checkOverlap checks if a showtime overlaps with existing showtimes in the same studio
// CRITICAL FUNCTION: Prevents double-booking of studios
//
// Overlap Logic:
// Two time ranges [A_start, A_end] and [B_start, B_end] overlap if:
// (A_start < B_end) AND (A_end > B_start)
//
// Parameters:
//   - excludeID: ID of showtime to exclude from check (used during updates, 0 for creates)
//   - studioID: ID of the studio to check
//   - newStart: Proposed start time
//   - newEnd: Proposed end time (already includes movie duration + buffer)
func (s *ShowtimeService) checkOverlap(excludeID uint, studioID uint, newStart time.Time, newEnd time.Time) error {
	var count int64

	query := s.db.Model(&models.Showtime{}).
		Where("studio_id = ?", studioID).
		Where("start_time < ?", newEnd).   // Existing start is before new end
		Where("end_time > ?", newStart)    // Existing end is after new start

	// Exclude current showtime when updating
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("schedule conflict: studio is already occupied during this time slot")
	}

	return nil
}
