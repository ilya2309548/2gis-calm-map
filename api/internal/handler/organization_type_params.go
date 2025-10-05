package handler

import (
	"net/http"

	"2gis-calm-map/api/internal/service"

	"github.com/gin-gonic/gin"
)

type OrganizationsParamsAverageByTypeRequest struct {
	OrganizationType string   `json:"organization_type" binding:"required"`
	Params           []string `json:"params" binding:"required,min=1"`
	// Optional threshold; if omitted, defaults to 3.0
	Threshold *float64 `json:"threshold"`
}

type OrganizationWithSelectedAverage struct {
	Organization interface{} `json:"organization"`
	Average      float64     `json:"average"`
	Params       []string    `json:"params"`
}

type OrganizationsParamsAverageByTypeResponse struct {
	OrganizationType string                            `json:"organization_type"`
	Items            []OrganizationWithSelectedAverage `json:"items"`
}

// Reuse existing services
var orgService = service.NewOrganizationService()
var orgParamsAggService = service.NewOrganizationParamsService()

// GetOrganizationsParamsAverageByType godoc
// @Summary Compute averages for each organization of given type
// @Description For every organization of a specified type computes (avg(p1)+...)/N and returns only those with average > threshold (default 3.0).
// @Tags organization-params
// @Accept json
// @Produce json
// @Param input body OrganizationsParamsAverageByTypeRequest true "Request"
// @Success 200 {object} OrganizationsParamsAverageByTypeResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /organization/params/average/by-type [post]
func GetOrganizationsParamsAverageByType(c *gin.Context) {
	// публичный доступ: без проверки роли

	var req OrganizationsParamsAverageByTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgs, err := orgService.GetByType(req.OrganizationType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]OrganizationWithSelectedAverage, 0, len(orgs))
	defaultThreshold := 3.0
	threshold := defaultThreshold
	if req.Threshold != nil {
		threshold = *req.Threshold
	}
	for _, org := range orgs {
		// Ensure params record exists
		paramsModel, err := orgParamsAggService.GetOrCreate(org.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		avg, err := orgParamsAggService.ComputeAverageAcross(paramsModel, req.Params)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if avg > threshold {
			items = append(items, OrganizationWithSelectedAverage{
				Organization: org,
				Average:      avg,
				Params:       req.Params,
			})
		}
	}

	c.JSON(http.StatusOK, OrganizationsParamsAverageByTypeResponse{
		OrganizationType: req.OrganizationType,
		Items:            items,
	})
}
