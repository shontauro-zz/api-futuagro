// Package services contains the interfaces for all use cases in the business domain.
package services

import (
	"futuagro.com/pkg/domain/dtos"
	"futuagro.com/pkg/domain/models"
	"futuagro.com/pkg/store"
)

// CountryService represent the country's domain service contract
type CountryService struct {
	repository *store.MongoCountryRepository
}

// FindCountryByID returns a country by its ID
func (s *CountryService) FindCountryByID(id string) (*models.Country, error) {
	return s.repository.FindByID(id)
}

// FindAllCountries returns a list of countries
func (s *CountryService) FindAllCountries() ([]*models.Country, error) {
	return s.repository.FindAll()
}

// CreateCountry create a new country record
func (s *CountryService) CreateCountry(country *models.Country) (string, error) {
	return s.repository.Insert(country)
}

// UpdateCountryByID update a country's data by its id
func (s *CountryService) UpdateCountryByID(id string, country *models.Country) (*models.Country, error) {
	return s.repository.Update(id, country)
}

// DeleteCountryByID delete a country by id
func (s *CountryService) DeleteCountryByID(id string) (bool, error) {
	return s.repository.Delete(id)
}

// AddState add a new state to a country
func (s *CountryService) AddState(countryID string, stateDto dtos.CountryStateDto) (*models.Country, error) {
	return s.repository.InsertCountryState(countryID, stateDto)
}

// UpdateState update a country state data
func (s *CountryService) UpdateState(countryID string, stateID string, stateDto dtos.CountryStateDto) (*models.Country, error) {
	return s.repository.UpdateCountryState(countryID, stateID, stateDto)
}

// DeleteState remove a state from a country
func (s *CountryService) DeleteState(countryID string, stateID string) (*models.Country, error) {
	return s.repository.DeleteCountryState(countryID, stateID)
}

// NewCountryService creates a country service with necessary dependencies.
func NewCountryService(countryRepository *store.MongoCountryRepository) *CountryService {
	return &CountryService{countryRepository}
}
