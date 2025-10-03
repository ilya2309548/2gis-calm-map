package handler

import (
	"2gis-calm-map/api/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUsers godoc
// @Summary Get users
// @Tags users
// @Produce json
// @Success 200 {array} model.User
// @Router /users [get]
func GetUsers(c *gin.Context) {
	users, err := repository.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, users)
}
