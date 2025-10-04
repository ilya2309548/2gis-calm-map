package handler

import (
	"errors"
	"net/http"
	"strings"

	"2gis-calm-map/api/internal/model"
	"2gis-calm-map/api/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var organizationService = service.NewOrganizationService()

type OrganizationCreateRequest struct {
	Address          string   `json:"address" binding:"required"`
	Longitude        *float64 `json:"longitude"`
	Latitude         *float64 `json:"latitude"`
	OrganizationType string   `json:"organization_type" binding:"required"`
}

type OrganizationUpdateRequest struct {
	Address          *string  `json:"address"`
	Longitude        *float64 `json:"longitude"`
	Latitude         *float64 `json:"latitude"`
	OrganizationType *string  `json:"organization_type"`
}

func roleAllowed(c *gin.Context) (uint, bool) {
	uidRaw, ok := c.Get("user_id")
	if !ok {
		return 0, false
	}
	role, _ := c.Get("role")
	if role == "admin" || role == "owner" { // допустимые роли
		return uidRaw.(uint), true
	}
	return 0, false
}

// CreateOrganization godoc
// @Summary Create organization
// @Description Create organization for current user (1:1). Roles: owner, admin
// @Tags organization
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body OrganizationCreateRequest true "Organization data"
// @Success 200 {object} model.Organization
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /organization [post]
func CreateOrganization(c *gin.Context) {
	ownerID, allowed := roleAllowed(c)
	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req OrganizationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org := model.Organization{
		OwnerID:          ownerID,
		Address:          req.Address,
		Longitude:        req.Longitude, // Optional, may be nil
		Latitude:         req.Latitude,  // Optional, may be nil
		OrganizationType: req.OrganizationType,
	}
	if err := organizationService.Create(&org); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") || strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			c.JSON(http.StatusConflict, gin.H{"error": "organization already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, org)
}

// GetOrganization godoc
// @Summary Get organization (owner self)
// @Description Get organization for current user (admin can also fetch for self only here)
// @Tags organization
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.Organization
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /organization [get]
func GetOrganization(c *gin.Context) {
	ownerID, allowed := roleAllowed(c)
	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	org, err := organizationService.GetByOwner(ownerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, org)
}

// PatchOrganization godoc
// @Summary Update organization
// @Description Partially update organization (owner/admin)
// @Tags organization
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body OrganizationUpdateRequest true "Fields to update"
// @Success 200 {object} model.Organization
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /organization [patch]
func PatchOrganization(c *gin.Context) {
	ownerID, allowed := roleAllowed(c)
	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req OrganizationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]interface{}{}
	if req.Address != nil {
		updates["address"] = *req.Address
	}
	if req.Longitude != nil {
		updates["longitude"] = *req.Longitude
	}
	if req.Latitude != nil {
		updates["latitude"] = *req.Latitude
	}
	if req.OrganizationType != nil {
		updates["organization_type"] = *req.OrganizationType
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields"})
		return
	}

	org, err := organizationService.UpdateByOwner(ownerID, updates)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, org)
}
