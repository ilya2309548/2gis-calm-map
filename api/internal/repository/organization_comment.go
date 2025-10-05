package repository

import (
	"2gis-calm-map/api/internal/db"
	"2gis-calm-map/api/internal/model"
)

func CreateOrganizationComment(c *model.OrganizationComment) error {
	return db.DB.Create(c).Error
}

func ListOrganizationComments(orgID uint) ([]model.OrganizationComment, error) {
	var list []model.OrganizationComment
	err := db.DB.Preload("User").Where("organization_id = ?", orgID).Order("id DESC").Find(&list).Error
	return list, err
}
