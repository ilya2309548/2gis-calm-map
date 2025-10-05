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

func GetOrganizationsByType(orgType string) ([]model.Organization, error) {
	var orgs []model.Organization
	err := db.DB.Preload("Params").Where("organization_type = ?", orgType).Find(&orgs).Error
	return orgs, err
}

func GetOrganizationByID(id uint) (model.Organization, error) {
	var org model.Organization
	err := db.DB.Preload("Params").First(&org, id).Error
	return org, err
}

func UpdateOrganizationFields(id uint, updates map[string]interface{}) error {
	return db.DB.Model(&model.Organization{}).Where("id = ?", id).Updates(updates).Error
}
