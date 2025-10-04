package handler

import (
	"net/http"

	"2gis-calm-map/api/internal/service"

	"github.com/gin-gonic/gin"
)

var orgParamsService = service.NewOrganizationParamsService()

type OrganizationParamsAverageRequest struct {
	OrganizationID uint     `json:"organization_id" binding:"required"`
	Params         []string `json:"params" binding:"required,min=1"`
}

type OrganizationParamsAverageResponse struct {
	OrganizationID uint     `json:"organization_id"`
	Params         []string `json:"params"`
	Average        float64  `json:"average"`
}

// GetOrganizationParamsAverage godoc
// @Summary Compute average across selected organization params
// @Description Returns (avg(param1)+...)/N for specified params. Requires JWT (admin/owner of org).
// @Tags organization-params
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body OrganizationParamsAverageRequest true "Organization params average request"
// @Success 200 {object} OrganizationParamsAverageResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /organization/params/average [post]
func GetOrganizationParamsAverage(c *gin.Context) {
	var req OrganizationParamsAverageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Простая проверка роли (из JWT middleware)
	roleValue, _ := c.Get("role")
	if roleValue != "admin" && roleValue != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	paramsModel, err := orgParamsService.GetOrCreate(req.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	avg, err := orgParamsService.ComputeAverageAcross(paramsModel, req.Params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, OrganizationParamsAverageResponse{
		OrganizationID: req.OrganizationID,
		Params:         req.Params,
		Average:        avg,
	})
}
