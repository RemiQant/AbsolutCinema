package services

import (
	"errors"

	"gorm.io/gorm"

	"absolutcinema-backend/internal/models"
)

type MovieService struct {
	db *gorm.DB
}

func NewMovieService(db *gorm.DB) *MovieService {
	return &MovieService{db: db}
}

// CreateMovie creates a new movie with validation
func (s *MovieService) CreateMovie(movie *models.Movie) error {
	// Validation
	if err := s.validateMovie(movie); err != nil {
		return err
	}
	
	return s.db.Create(movie).Error
}

// GetAllMovies retrieves all movies
func (s *MovieService) GetAllMovies() ([]models.Movie, error) {
	var movies []models.Movie
	if err := s.db.Find(&movies).Error; err != nil {
		return nil, err
	}
	return movies, nil
}

// GetMovieByID retrieves a movie by ID
func (s *MovieService) GetMovieByID(id uint) (*models.Movie, error) {
	var movie models.Movie
	if err := s.db.First(&movie, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("movie not found")
		}
		return nil, err
	}
	return &movie, nil
}

// UpdateMovie updates an existing movie
func (s *MovieService) UpdateMovie(id uint, updates *models.Movie) error {
	// Validation
	if err := s.validateMovie(updates); err != nil {
		return err
	}
	
	// Check if movie exists
	var movie models.Movie
	if err := s.db.First(&movie, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("movie not found")
		}
		return err
	}
	
	// Update fields
	movie.Title = updates.Title
	movie.Description = updates.Description
	movie.DurationMinutes = updates.DurationMinutes
	movie.PosterURL = updates.PosterURL
	movie.Rating = updates.Rating
	
	return s.db.Save(&movie).Error
}

// DeleteMovie soft deletes a movie
func (s *MovieService) DeleteMovie(id uint) error {
	result := s.db.Delete(&models.Movie{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("movie not found")
	}
	return nil
}

// validateMovie validates movie data
func (s *MovieService) validateMovie(movie *models.Movie) error {
	if movie.Title == "" {
		return errors.New("movie title is required")
	}
	
	if movie.DurationMinutes <= 0 {
		return errors.New("duration must be greater than 0")
	}
	
	return nil
}
