package repository

import (
	"2gis-calm-map/api/internal/db"
	"2gis-calm-map/api/internal/model"
)

func GetAllUsers() ([]model.User, error) {
	var users []model.User
	err := db.DB.Find(&users).Error
	return users, err
}
