package service

import (
	"2gis-calm-map/api/internal/model"
	"2gis-calm-map/api/internal/repository"
)

type OrganizationCommentService struct{}

func NewOrganizationCommentService() *OrganizationCommentService {
	return &OrganizationCommentService{}
}

// CreateWithAggregation creates a comment then updates OrganizationParams aggregate.
func (s *OrganizationCommentService) CreateWithAggregation(c *model.OrganizationComment, p *model.OrganizationParams) error {
	if err := repository.CreateOrganizationComment(c); err != nil {
		return err
	}
	// update aggregates for each non-nil numeric field
	apply := func(val *uint, sum *uint, count *uint, avg *float64) {
		if val == nil || *val == 0 { // treat 0 as not provided per spec ("ненулевых")
			return
		}
		*sum += *val
		*count += 1
		*avg = float64(*sum) / float64(*count)
	}
	apply(c.AppearanceValue, &p.AppearanceSum, &p.AppearanceCount, &p.AppearanceAvg)
	apply(c.LightingValue, &p.LightingSum, &p.LightingCount, &p.LightingAvg)
	apply(c.SmellValue, &p.SmellSum, &p.SmellCount, &p.SmellAvg)
	apply(c.TemperatureValue, &p.TemperatureSum, &p.TemperatureCount, &p.TemperatureAvg)
	apply(c.TactilityValue, &p.TactilitySum, &p.TactilityCount, &p.TactilityAvg)
	apply(c.SignageValue, &p.SignageSum, &p.SignageCount, &p.SignageAvg)
	apply(c.IntuitivenessValue, &p.IntuitivenessSum, &p.IntuitivenessCount, &p.IntuitivenessAvg)
	apply(c.StaffAttitudeValue, &p.StaffAttitudeSum, &p.StaffAttitudeCount, &p.StaffAttitudeAvg)
	apply(c.PeopleDensityValue, &p.PeopleDensitySum, &p.PeopleDensityCount, &p.PeopleDensityAvg)
	apply(c.SelfServiceValue, &p.SelfServiceSum, &p.SelfServiceCount, &p.SelfServiceAvg)
	apply(c.CalmnessValue, &p.CalmnessSum, &p.CalmnessCount, &p.CalmnessAvg)

	// persist updated aggregates
	return repository.UpdateOrganizationParams(p)
}

func (s *OrganizationCommentService) ListByOrganization(orgID uint) ([]model.OrganizationComment, error) {
	return repository.ListOrganizationComments(orgID)
}
