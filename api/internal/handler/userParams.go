package handler

import (
	"2gis-calm-map/api/internal/model"
	"2gis-calm-map/api/internal/service"
	"net/http"
	"strconv"

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

// GetUserParams godoc
// @Summary Get user parameters by user id
// @Description Returns user evaluation parameters. Requires JWT token.
// @Tags user-params
// @Produce json
// @Security BearerAuth
// @Param user_id path int true "User ID"
// @Success 200 {object} model.UserParams
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /user-params/{user_id} [get]
func GetUserParams(c *gin.Context) {
	// Можно разрешить получать только свои параметры: сравнить user_id из токена и path param
	// Для простоты сейчас просто возвращаем по path.
	pathID := c.Param("user_id")
	if pathID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}
	parsed, err := strconv.ParseUint(pathID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	// (Опционально) сверка с токеном:
	if tokenUID, ok := c.Get("user_id"); ok {
		if uint(parsed) != tokenUID.(uint) {
			// Если нужно ограничение – раскомментируйте:
			// c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			// return
		}
	}
	params, err := userParamsService.GetUserParamsByUserID(uint(parsed))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, params)
}
