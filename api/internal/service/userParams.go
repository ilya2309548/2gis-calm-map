package service

import (
	"2gis-calm-map/api/internal/model"
	"2gis-calm-map/api/internal/repository"
)

type UserParamsService struct{}

func NewUserParamsService() *UserParamsService {
	return &UserParamsService{}
}

// CreateUserParams сохраняет параметры пользователя.
// Принимает уже собранную модель без небезопасных преобразований типов.
func (s *UserParamsService) CreateUserParams(params model.UserParams) (model.UserParams, error) {
	return repository.CreateUserParams(params)
}
