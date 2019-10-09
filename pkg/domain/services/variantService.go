// Package services contains the interfaces for all use cases in the business domain.
package services

import (
	"futuagro.com/pkg/domain/dtos"
	"futuagro.com/pkg/domain/models"
	"futuagro.com/pkg/store"
)

// VariantService implements use cases methods and domain business logic for variants
type VariantService struct {
	repository *store.MongoVariantRepository
}

//FindVariantByID return a variant by its ID
func (s *VariantService) FindVariantByID(ID string) (*models.Variant, error) {
	return s.repository.FindVariantByID(ID)
}

// FindOneVariantByItemID returns a variant by its ID and item ID
func (s *VariantService) FindOneVariantByItemID(itemID string, variantID string) (*models.Variant, error) {
	return s.repository.FindOneVariantByItemID(itemID, variantID)
}

// FindVariantsByItemID returns a list of variants that belongs to an item
func (s *VariantService) FindVariantsByItemID(itemID string) ([]*models.Variant, error) {
	return s.repository.FindVariantsByItemID(itemID)
}

// CreateVariant create a new Variant record
func (s *VariantService) CreateVariant(itemID string, dto *dtos.VariantDto) (string, error) {
	return s.repository.Insert(itemID, dto)
}

// UpdateVariant update a variant data
func (s *VariantService) UpdateVariant(itemID string, variantID string, itemDto *dtos.VariantDto) (*models.Variant, error) {
	return s.repository.Update(itemID, variantID, itemDto)
}

// DeleteVariant delete a variant by id
func (s *VariantService) DeleteVariant(itemID string, variantID string) (bool, error) {
	return s.repository.Delete(itemID, variantID)
}

// NewVariantService creates a variant service with necessary dependencies.
func NewVariantService(repository *store.MongoVariantRepository) *VariantService {
	return &VariantService{repository}
}
