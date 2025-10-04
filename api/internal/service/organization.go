package service

import (
	"2gis-calm-map/api/internal/model"
	"2gis-calm-map/api/internal/repository"
)

type OrganizationService struct{}

func NewOrganizationService() *OrganizationService { return &OrganizationService{} }

func (s *OrganizationService) Create(org *model.Organization) error {
	return repository.CreateOrganization(org)
}

func (s *OrganizationService) GetByOwner(ownerID uint) (model.Organization, error) {
	return repository.GetOrganizationByOwner(ownerID)
}

func (s *OrganizationService) UpdateByOwner(ownerID uint, updates map[string]interface{}) (model.Organization, error) {
	return repository.UpdateOrganizationByOwner(ownerID, updates)
}

func (s *OrganizationService) GetByType(orgType string) ([]model.Organization, error) {
	return repository.GetOrganizationsByType(orgType)
}
