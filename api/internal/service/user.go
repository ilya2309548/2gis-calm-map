package service

import (
	"2gis-calm-map/api/internal/model"
	"2gis-calm-map/api/internal/repository"
)

// UserService handles user business logic
type UserService struct{}

// NewUserService creates a new UserService
func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	return repository.GetAllUsers()
}
