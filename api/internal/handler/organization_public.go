package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrganizationByAddressRequest request body for public lookup
type OrganizationByAddressRequest struct {
	Address string `json:"address" binding:"required"`
}

type OrganizationByAddressResponse struct {
	ID               uint     `json:"id"`
	Address          string   `json:"address"`
	OrganizationType string   `json:"organization_type"`
	Longitude        *float64 `json:"longitude"`
	Latitude         *float64 `json:"latitude"`
	MapPath          *string  `json:"map_path"`
	PicturePath      *string  `json:"picture_path"`
}

// GetOrganizationByAddressPublic godoc
// @Summary Public organization lookup by address
// @Description Возвращает организацию (id и основные поля) по точному адресу. Публично.
// @Tags organization
// @Accept json
// @Produce json
// @Param input body OrganizationByAddressRequest true "Address lookup"
// @Success 200 {object} OrganizationByAddressResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /organization/public/by-address [post]
func GetOrganizationByAddressPublic(c *gin.Context) {
	var req OrganizationByAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	org, err := organizationService.GetByAddress(req.Address)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := OrganizationByAddressResponse{
		ID:               org.ID,
		Address:          org.Address,
		OrganizationType: org.OrganizationType,
		Longitude:        org.Longitude,
		Latitude:         org.Latitude,
		MapPath:          org.MapPath,
		PicturePath:      org.PicturePath,
	}
	c.JSON(http.StatusOK, resp)
}
