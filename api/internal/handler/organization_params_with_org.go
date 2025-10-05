package handler

import (
	"net/http"

	"2gis-calm-map/api/internal/repository"
	"2gis-calm-map/api/internal/service"

	"github.com/gin-gonic/gin"
)

var orgParamsWithOrgService = service.NewOrganizationParamsService()

// OrganizationParamsWithOrgRequest extends average request.
type OrganizationParamsWithOrgRequest struct {
	OrganizationID uint     `json:"organization_id" binding:"required"`
	Params         []string `json:"params" binding:"required,min=1"`
}

type OrganizationParamsWithOrgResponse struct {
	Organization struct {
		ID               uint     `json:"id"`
		Address          string   `json:"address"`
		OrganizationType string   `json:"organization_type"`
		Longitude        *float64 `json:"longitude"`
		Latitude         *float64 `json:"latitude"`
		MapPath          *string  `json:"map_path"`
		PicturePath      *string  `json:"picture_path"`
	} `json:"organization"`
	Params  []string `json:"params"`
	Average float64  `json:"average"`
}

// GetOrganizationParamsAverageWithOrganizationInfo godoc
// @Summary Compute average and return organization info
// @Description Как /organization/params/average, но дополняет адресом/типом/координатами и путями изображений. Публично.
// @Tags organization-params
// @Accept json
// @Produce json
// @Param input body OrganizationParamsWithOrgRequest true "Request"
// @Success 200 {object} OrganizationParamsWithOrgResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /organization/params/average/with-info [post]
func GetOrganizationParamsAverageWithOrganizationInfo(c *gin.Context) {
	var req OrganizationParamsWithOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := repository.GetOrganizationByID(req.OrganizationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "organization not found"})
		return
	}

	paramsModel, err := orgParamsWithOrgService.GetOrCreate(req.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	avg, err := orgParamsWithOrgService.ComputeAverageAcross(paramsModel, req.Params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := OrganizationParamsWithOrgResponse{Params: req.Params, Average: avg}
	resp.Organization.ID = org.ID
	resp.Organization.Address = org.Address
	resp.Organization.OrganizationType = org.OrganizationType
	resp.Organization.Longitude = org.Longitude
	resp.Organization.Latitude = org.Latitude
	resp.Organization.MapPath = org.MapPath
	resp.Organization.PicturePath = org.PicturePath

	c.JSON(http.StatusOK, resp)
}
