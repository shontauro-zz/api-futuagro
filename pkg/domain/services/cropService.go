// Package services contains the interfaces for all use cases in the business domain.
package services

import (
	"futuagro.com/pkg/domain/dtos"
	"futuagro.com/pkg/domain/models"
	"futuagro.com/pkg/store"
)

// CropService implements use cases methods and domain business logic for crops
type CropService struct {
	repository *store.MongoCropRepository
}

// FindCropByID returns a crop by its ID
func (s *CropService) FindCropByID(id string) (*models.Crop, error) {
	return s.repository.FindByID(id)
}

// FindAllCrops returns a list of crops
func (s *CropService) FindAllCrops() ([]*models.Crop, error) {
	return s.repository.FindAll()
}

// CreateCrop create a new crop record
func (s *CropService) CreateCrop(dto *dtos.CropDto) (*models.Crop, error) {
	result, err := s.repository.Insert(dto)
	if err != nil {
		return nil, err
	}

	crop, err := s.repository.FindByID(result)
	if err != nil {
		return nil, err
	}

	return crop, nil
}

// UpdateCropByID update a crop data by its id
func (s *CropService) UpdateCropByID(id string, dto *dtos.CropDto) (*models.Crop, error) {
	result, err := s.repository.Update(id, dto)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	crop, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	return crop, nil
}

// DeleteCropByID delete a crop by id
func (s *CropService) DeleteCropByID(id string) (bool, error) {
	return s.repository.Delete(id)
}

// NewCropService creates a crop service with necessary dependencies.
func NewCropService(repository *store.MongoCropRepository) *CropService {
	return &CropService{repository}
}
