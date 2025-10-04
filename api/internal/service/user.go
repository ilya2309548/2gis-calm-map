package service

import (
	"2gis-calm-map/api/internal/model"
	"2gis-calm-map/api/internal/repository"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	return repository.GetAllUsers()
}

func (s *UserService) CreateUser(name, email, password, role string) (model.User, error) {
	return repository.CreateUser(name, email, password, role)
}
