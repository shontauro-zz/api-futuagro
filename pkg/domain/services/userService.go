// Package services contains the interfaces for all use cases in the business domain.
package services

import (
	"futuagro.com/pkg/domain/dtos"
	"futuagro.com/pkg/domain/models"
	"futuagro.com/pkg/store"
)

// UserService implements use cases methods and domain business logic for users
type UserService struct {
	repository *store.MongoUserRepository
}

// FindUserByID returns an user by its ID
func (s *UserService) FindUserByID(id string) (*models.User, error) {
	return s.repository.FindByID(id)
}

// PopulateUserByID return an user with the crops property populated with the variant data
func (s *UserService) PopulateUserByID(id string) (*models.User, error) {
	return s.repository.PopulateUserByID(id)
}

// FindAllUsers returns a list of users
func (s *UserService) FindAllUsers() ([]*models.User, error) {
	return s.repository.FindAll()
}

// Signup create a new user record
func (s *UserService) Signup(dto *dtos.UserDto) (*models.User, error) {
	result, err := s.repository.Insert(dto)
	if err != nil {
		return nil, err
	}

	user, err := s.repository.FindByID(result)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUserByID update an user data by its id
func (s *UserService) UpdateUserByID(id string, dto *dtos.UserDto) (*models.User, error) {
	result, err := s.repository.Update(id, dto)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	user, err := s.repository.PopulateUserByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser delete an user by id
func (s *UserService) DeleteUser(id string) (bool, error) {
	return s.repository.Delete(id)
}

// NewUserService creates an user service with necessary dependencies.
func NewUserService(repository *store.MongoUserRepository) *UserService {
	return &UserService{repository}
}
