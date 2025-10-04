package repository

import (
	"2gis-calm-map/api/internal/db"
	"2gis-calm-map/api/internal/model"
)

func GetOrganizationParams(orgID uint) (model.OrganizationParams, error) {
	var p model.OrganizationParams
	err := db.DB.Where("organization_id = ?", orgID).First(&p).Error
	return p, err
}

func CreateEmptyOrganizationParams(orgID uint) (model.OrganizationParams, error) {
	p := model.OrganizationParams{OrganizationID: orgID}
	err := db.DB.Create(&p).Error
	return p, err
}
