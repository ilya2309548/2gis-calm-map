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

func (s *UserParamsService) UpdateUserParamsByUserID(userID uint, req interface{}) (model.UserParams, error) {
	// Ожидаем map[string]interface{} или структуру model.UserParams / CreateUserParamsRequest.
	// Преобразуем в map для Updates.
	updates := map[string]interface{}{}
	switch v := req.(type) {
	case map[string]interface{}:
		updates = v
	case model.UserParams:
		updates = map[string]interface{}{
			"appearance": v.Appearance,
			"lighting": v.Lighting,
			"smell": v.Smell,
			"temperature": v.Temperature,
			"tactility": v.Tactility,
			"signage": v.Signage,
			"intuitiveness": v.Intuitiveness,
			"staff_attitude": v.StaffAttitude,
			"people_density": v.PeopleDensity,
			"self_service": v.SelfService,
			"calmness": v.Calmness,
		}
	default:
		// Попытаемся через reflection не заморачиваясь – пропускаем, оставляем пустым
	}
	return repository.UpdateUserParamsByUserID(userID, updates)
}
