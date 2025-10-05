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
	var counted int
	addIfPositive := func(val float64) {
		if val > 0 { // игнорируем нули как просили
			sum += val
			counted++
		}
	}
	for _, raw := range params {
		name := strings.ToLower(raw)
		switch name {
		case "appearance":
			addIfPositive(p.AppearanceAvg)
		case "lighting":
			addIfPositive(p.LightingAvg)
		case "smell":
			addIfPositive(p.SmellAvg)
		case "temperature":
			addIfPositive(p.TemperatureAvg)
		case "tactility":
			addIfPositive(p.TactilityAvg)
		case "signage":
			addIfPositive(p.SignageAvg)
		case "intuitiveness":
			addIfPositive(p.IntuitivenessAvg)
		case "staffattitude", "staff_attitude":
			addIfPositive(p.StaffAttitudeAvg)
		case "peopledensity", "people_density":
			addIfPositive(p.PeopleDensityAvg)
		case "selfservice", "self_service":
			addIfPositive(p.SelfServiceAvg)
		case "calmness":
			addIfPositive(p.CalmnessAvg)
		default:
			return 0, fmt.Errorf("unknown param: %s", raw)
		}
	}
	if counted == 0 {
		return 0, nil // все выбранные параметры имели среднее 0
	}
	return sum / float64(counted), nil
}
