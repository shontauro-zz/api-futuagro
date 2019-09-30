// Package services contains the interfaces for all use cases in the business domain.
package services

import (
	"futuagro.com/pkg/domain/dtos"
	"futuagro.com/pkg/domain/models"
	"futuagro.com/pkg/store"
)

// ItemService implements use cases methods and domain business logic for items
type ItemService struct {
	repository *store.MongoItemRepository
}

// FindItemByID returns an Item by its ID
func (s *ItemService) FindItemByID(id string) (*models.Item, error) {
	return s.repository.FindByID(id)
}

// FindAllItems returns a list of items
func (s *ItemService) FindAllItems() ([]*models.Item, error) {
	return s.repository.FindAll()
}

// CreateItem create a new Item record
func (s *ItemService) CreateItem(dto *dtos.ItemDto) (string, error) {
	return s.repository.Insert(dto)
}

// UpdateItemByID update an item data by its id
func (s *ItemService) UpdateItemByID(id string, itemDto *dtos.ItemDto) (*models.Item, error) {
	return s.repository.Update(id, itemDto)
}

// DeleteItemByID delete an item by id
func (s *ItemService) DeleteItemByID(id string) (bool, error) {
	return s.repository.Delete(id)
}

// NewItemService creates an Item service with necessary dependencies.
func NewItemService(repository *store.MongoItemRepository) *ItemService {
	return &ItemService{repository}
}
