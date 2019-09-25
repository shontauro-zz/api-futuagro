// Package services contains the interfaces for all use cases in the business domain.
package services

import (
	"futuagro.com/pkg/domain/dtos"
	"futuagro.com/pkg/domain/models"
	"futuagro.com/pkg/store"
)

// CityService represent the country's domain service contract
type CityService struct {
	repository *store.MongoCityRepository
}

// FindCityByID returns a city by its ID
func (s *CityService) FindCityByID(id string) (*models.City, error) {
	return s.repository.FindByID(id)
}

// FindAllCities returns a list of cities
func (s *CityService) FindAllCities() ([]*models.City, error) {
	return s.repository.FindAll()
}

//FindAllCitiesByCountryState return a list of cities that belongs to a countryState
func (s *CityService) FindAllCitiesByCountryState(stateID string) ([]*models.City, error) {
	return s.repository.FindCitiesByCountryState(stateID)
}

// CreateCity create a new city record
func (s *CityService) CreateCity(stateID string, city *dtos.CityDto) (string, error) {
	return s.repository.Insert(stateID, city)
}

// UpdateCityByID update a city's data by its id
func (s *CityService) UpdateCityByID(stateID string, cityID string, cityDto *dtos.CityDto) (*models.City, error) {
	return s.repository.Update(stateID, cityID, cityDto)
}

// DeleteCityByID delete a city by id
func (s *CityService) DeleteCityByID(stateID string, cityID string) (bool, error) {
	return s.repository.Delete(stateID, cityID)
}

// NewCityService creates a country service with necessary dependencies.
func NewCityService(cityRepository *store.MongoCityRepository) *CityService {
	return &CityService{cityRepository}
}
