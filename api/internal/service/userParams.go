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

func (s *UserParamsService) GetUserParamsByUserID(userID uint) (model.UserParams, error) {
	return repository.GetUserParamsByUserID(userID)
}
