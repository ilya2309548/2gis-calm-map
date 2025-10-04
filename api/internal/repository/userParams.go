package repository

import (
	"2gis-calm-map/api/internal/db"
	"2gis-calm-map/api/internal/model"
)

func CreateUserParams(params model.UserParams) (model.UserParams, error) {
	err := db.DB.Create(&params).Error
	return params, err
}
