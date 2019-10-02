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
func (s *SupplierService) CreateSupplier(supplier *dtos.SupplierDto) (string, error) {
	return s.repository.Insert(supplier)
}

// UpdateSupplierByID update a supplier data by its id
func (s *SupplierService) UpdateSupplierByID(id string, supplier *dtos.SupplierDto) (*models.Supplier, error) {
	return s.repository.Update(id, supplier)
}

// DeleteSupplier delete a suplier by id
func (s *SupplierService) DeleteSupplier(id string) (bool, error) {
	return s.repository.Delete(id)
}

// AddCrop add a new Crop to a supplier
func (s *SupplierService) AddCrop(supplierID string, cropDto dtos.CropDto) (*models.Supplier, error) {
	return s.repository.InsertCrop(supplierID, cropDto)
}

// NewSupplierService creates a supplier service with necessary dependencies.
func NewSupplierService(supplierRepository *store.MongoSupplierRepository) *SupplierService {
	return &SupplierService{supplierRepository}
}
