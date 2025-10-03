package handler

import (
	"2gis-calm-map/api/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get users
// @Produce json
// @Success 200 {array} model.User
// @Router /users [get]

func GetUsers(c *gin.Context) {
	users := []model.User{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Doe", Email: "jane@example.com"},
	}
	c.JSON(http.StatusOK, users)
}
