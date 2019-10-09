// Package services contains the interfaces for all use cases in the business domain.
package services

import (
	"strings"

	"futuagro.com/pkg/domain/dtos"
	"futuagro.com/pkg/domain/models"
	"futuagro.com/pkg/store"
	"golang.org/x/crypto/bcrypt"
)

// AuthService implements use cases methods and domain business logic for authorizing users
type AuthService struct {
	userRepository *store.MongoUserRepository
}

// Login authenticates an user
func (s *AuthService) Login(dto *dtos.LoginDto) (*models.User, error) {
	user, err := s.userRepository.FindByEmail(strings.ToLower(dto.Email))
	if err != nil {
		return nil, err
	}
	if user != nil {
		err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(dto.Password))
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	return nil, nil
}

// NewAuthService creates an auth service with necessary dependencies.
func NewAuthService(userRepository *store.MongoUserRepository) *AuthService {
	return &AuthService{userRepository}
}
