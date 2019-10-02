// Package services contains the interfaces for all use cases in the business domain.
package services

import (
	"futuagro.com/pkg/domain/dtos"
	"futuagro.com/pkg/domain/models"
	"futuagro.com/pkg/store"
)

// SupplierService implements use cases methods and domain business logic for suppliers
type SupplierService struct {
	repository *store.MongoSupplierRepository
}

// FindSupplierByID returns a supplier by its ID
func (s *SupplierService) FindSupplierByID(id string) (*models.Supplier, error) {
	return s.repository.FindByID(id)
}

// PopulateSupplierByID return a supplier with the crops property populated with the variant data
func (s *SupplierService) PopulateSupplierByID(id string) (*models.Supplier, error) {
	return s.repository.PopulateSupplierByID(id)
}

// FindAllSuppliers returns a list of suppliers
func (s *SupplierService) FindAllSuppliers() ([]*models.Supplier, error) {
	return s.repository.FindAll()
}

// CreateSupplier create a new supplier record
func (s *SupplierService) CreateSupplier(dto *dtos.SupplierDto) (*models.Supplier, error) {
	result, err := s.repository.Insert(dto)
	if err != nil {
		return nil, err
	}

	supplier, err := s.repository.FindByID(result)
	if err != nil {
		return nil, err
	}

	return supplier, nil
}

// UpdateSupplierByID update a supplier data by its id
func (s *SupplierService) UpdateSupplierByID(id string, dto *dtos.SupplierDto) (*models.Supplier, error) {
	result, err := s.repository.Update(id, dto)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	supplier, err := s.repository.PopulateSupplierByID(id)
	if err != nil {
		return nil, err
	}

	return supplier, nil
}

// DeleteSupplier delete a suplier by id
func (s *SupplierService) DeleteSupplier(id string) (bool, error) {
	return s.repository.Delete(id)
}

// AddCrop add a new Crop to a supplier
func (s *SupplierService) AddCrop(supplierID string, cropDto dtos.CropDto) (*models.Supplier, error) {

	result, err := s.repository.InsertCrop(supplierID, cropDto)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	supplier, err := s.repository.PopulateSupplierByID(supplierID)
	if err != nil {
		return nil, err
	}

	return supplier, nil
}

// NewSupplierService creates a supplier service with necessary dependencies.
func NewSupplierService(supplierRepository *store.MongoSupplierRepository) *SupplierService {
	return &SupplierService{supplierRepository}
}
