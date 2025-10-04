package repository

import (
	"2gis-calm-map/api/internal/db"
	"2gis-calm-map/api/internal/model"
)

func CreateOrganization(org *model.Organization) error {
	return db.DB.Create(org).Error
}

func GetOrganizationByOwner(ownerID uint) (model.Organization, error) {
	var org model.Organization
	err := db.DB.Where("owner_id = ?", ownerID).First(&org).Error
	return org, err
}

func UpdateOrganizationByOwner(ownerID uint, updates map[string]interface{}) (model.Organization, error) {
	var org model.Organization
	if err := db.DB.Where("owner_id = ?", ownerID).First(&org).Error; err != nil {
		return org, err
	}
	if err := db.DB.Model(&org).Updates(updates).Error; err != nil {
		return org, err
	}
	return org, nil
}
