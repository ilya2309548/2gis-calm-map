package handler

import (
	"2gis-calm-map/api/internal/model"
	"2gis-calm-map/api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateUserParamsRequest struct {
	Appearance    bool `json:"appearance"`
	Lighting      bool `json:"lighting"`
	Smell         bool `json:"smell"`
	Temperature   bool `json:"temperature"`
	Tactility     bool `json:"tactility"`
	Signage       bool `json:"signage"`
	Intuitiveness bool `json:"intuitiveness"`
	StaffAttitude bool `json:"staff_attitude"`
	PeopleDensity bool `json:"people_density"`
	SelfService   bool `json:"self_service"`
	Calmness      bool `json:"calmness"`
}

var userParamsService = service.NewUserParamsService()

// CreateUserParams godoc
// @Summary Create user parameters
// @Description Saves user evaluation parameters. Requires JWT token.
// @Tags user-params
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body CreateUserParamsRequest true "User parameters"
// @Success 200 {object} model.UserParams
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Server error"
// @Router /user-params [post]
func CreateUserParams(c *gin.Context) {
	var req CreateUserParamsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user_id in context"})
		return
	}

	params, err := userParamsService.CreateUserParams(model.UserParams{
		UserID:        userID.(uint),
		Appearance:    req.Appearance,
		Lighting:      req.Lighting,
		Smell:         req.Smell,
		Temperature:   req.Temperature,
		Tactility:     req.Tactility,
		Signage:       req.Signage,
		Intuitiveness: req.Intuitiveness,
		StaffAttitude: req.StaffAttitude,
		PeopleDensity: req.PeopleDensity,
		SelfService:   req.SelfService,
		Calmness:      req.Calmness,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, params)
}
