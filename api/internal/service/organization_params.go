package service

import (
	"errors"
	"fmt"
	"strings"

	"2gis-calm-map/api/internal/model"
	"2gis-calm-map/api/internal/repository"

	"gorm.io/gorm"
)

type OrganizationParamsService struct{}

func NewOrganizationParamsService() *OrganizationParamsService { return &OrganizationParamsService{} }

// GetOrCreate returns existing params or creates empty one.
func (s *OrganizationParamsService) GetOrCreate(orgID uint) (model.OrganizationParams, error) {
	p, err := repository.GetOrganizationParams(orgID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repository.CreateEmptyOrganizationParams(orgID)
		}
		return model.OrganizationParams{}, err
	}
	return p, nil
}

// ComputeAverageAcross returns (avg(param1)+avg(param2)+...)/n for provided param names.
func (s *OrganizationParamsService) ComputeAverageAcross(p model.OrganizationParams, params []string) (float64, error) {
	if len(params) == 0 {
		return 0, fmt.Errorf("no params provided")
	}
	var sum float64
	for _, raw := range params {
		name := strings.ToLower(raw)
		switch name {
		case "appearance":
			sum += p.AppearanceAvg
		case "lighting":
			sum += p.LightingAvg
		case "smell":
			sum += p.SmellAvg
		case "temperature":
			sum += p.TemperatureAvg
		case "tactility":
			sum += p.TactilityAvg
		case "signage":
			sum += p.SignageAvg
		case "intuitiveness":
			sum += p.IntuitivenessAvg
		case "staffattitude", "staff_attitude":
			sum += p.StaffAttitudeAvg
		case "peopledensity", "people_density":
			sum += p.PeopleDensityAvg
		case "selfservice", "self_service":
			sum += p.SelfServiceAvg
		case "calmness":
			sum += p.CalmnessAvg
		default:
			return 0, fmt.Errorf("unknown param: %s", raw)
		}
	}
	return sum / float64(len(params)), nil
}
