package services

import (
	"errors"

	"gorm.io/gorm"

	"absolutcinema-backend/internal/models"
)

type StudioService struct {
	db *gorm.DB
}

func NewStudioService(db *gorm.DB) *StudioService {
	return &StudioService{db: db}
}

// CreateStudio creates a new studio with validation
func (s *StudioService) CreateStudio(studio *models.Studio) error {
	// Validation
	if err := s.validateStudio(studio); err != nil {
		return err
	}
	
	return s.db.Create(studio).Error
}

// GetAllStudios retrieves all studios (including soft deleted if needed)
func (s *StudioService) GetAllStudios() ([]models.Studio, error) {
	var studios []models.Studio
	if err := s.db.Find(&studios).Error; err != nil {
		return nil, err
	}
	return studios, nil
}

// GetStudioByID retrieves a studio by ID
func (s *StudioService) GetStudioByID(id uint) (*models.Studio, error) {
	var studio models.Studio
	if err := s.db.First(&studio, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("studio not found")
		}
		return nil, err
	}
	return &studio, nil
}

// UpdateStudio updates an existing studio
func (s *StudioService) UpdateStudio(id uint, updates *models.Studio) error {
	// Validation
	if err := s.validateStudio(updates); err != nil {
		return err
	}
	
	// Check if studio exists
	var studio models.Studio
	if err := s.db.First(&studio, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("studio not found")
		}
		return err
	}
	
	// Update fields
	studio.Name = updates.Name
	studio.TotalRows = updates.TotalRows
	studio.TotalCols = updates.TotalCols
	
	return s.db.Save(&studio).Error
}

// DeleteStudio soft deletes a studio
func (s *StudioService) DeleteStudio(id uint) error {
	result := s.db.Delete(&models.Studio{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("studio not found")
	}
	return nil
}

// validateStudio validates studio data
func (s *StudioService) validateStudio(studio *models.Studio) error {
	if studio.Name == "" {
		return errors.New("studio name is required")
	}
	
	if studio.TotalRows <= 0 || studio.TotalRows > 20 {
		return errors.New("total rows must be between 1 and 20")
	}
	
	if studio.TotalCols <= 0 || studio.TotalCols > 20 {
		return errors.New("total columns must be between 1 and 20")
	}
	
	return nil
}
