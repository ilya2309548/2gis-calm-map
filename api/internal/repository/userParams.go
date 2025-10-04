package repository

import (
	"2gis-calm-map/api/internal/db"
	"2gis-calm-map/api/internal/model"
)

func CreateUserParams(params model.UserParams) (model.UserParams, error) {
	err := db.DB.Create(&params).Error
	return params, err
}

func GetUserParamsByUserID(userID uint) (model.UserParams, error) {
	var params model.UserParams
	err := db.DB.Where("user_id = ?", userID).First(&params).Error
	return params, err
}

func UpdateUserParamsByUserID(userID uint, updates map[string]interface{}) (model.UserParams, error) {
	var params model.UserParams
	if err := db.DB.Where("user_id = ?", userID).First(&params).Error; err != nil {
		return params, err
	}
	if err := db.DB.Model(&params).Updates(updates).Error; err != nil {
		return params, err
	}
	return params, nil
}
