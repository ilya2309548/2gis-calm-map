package handler

import (
	"2gis-calm-map/api/internal/model"
	"2gis-calm-map/api/internal/service"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v5"
)

type RegisterRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type RegisterResponse struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

var userServiceReg = service.NewUserService()

// Register godoc
// @Summary Register new user
// @Description Creates a new user and returns a JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param input body RegisterRequest true "User registration data"
// @Success 200 {object} RegisterResponse
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Server error"
// @Router /register [post]
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := userServiceReg.CreateUser(req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Лучше вынести секрет в .env!
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, RegisterResponse{
		Token: tokenStr,
		User:  user,
	})
}
