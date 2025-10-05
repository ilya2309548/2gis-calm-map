package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Response item for comment list
type OrganizationCommentListItem struct {
	ID       uint     `json:"id"`
	UserID   uint     `json:"user_id"`
	UserName string   `json:"user_name"`
	Text     *string  `json:"text"`
	AvgValue *float64 `json:"avg_val"`
}

type OrganizationCommentListResponse struct {
	OrganizationID uint                          `json:"organization_id"`
	Items          []OrganizationCommentListItem `json:"items"`
}

// GetOrganizationComments godoc
// @Summary List comments for organization
// @Description Возвращает список комментариев организации с автором и средней оценкой.
// @Tags organization-comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organization_id path int true "Organization ID"
// @Success 200 {object} OrganizationCommentListResponse
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /organization/{organization_id}/comments [get]
func GetOrganizationComments(c *gin.Context) {
	// any authenticated user can view
	if _, exists := c.Get("user_id"); !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	orgIDParam := c.Param("organization_id")
	var orgID uint
	if _, err := fmt.Sscan(orgIDParam, &orgID); err != nil || orgID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization_id"})
		return
	}

	// ensure org exists
	if _, err := orgService.GetByID(orgID); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "organization not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	list, err := orgCommentService.ListByOrganization(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]OrganizationCommentListItem, 0, len(list))
	for _, cmt := range list {
		name := cmt.User.Name
		items = append(items, OrganizationCommentListItem{
			ID:       cmt.ID,
			UserID:   cmt.UserID,
			UserName: name,
			Text:     cmt.Text,
			AvgValue: cmt.AvgValue,
		})
	}

	c.JSON(http.StatusOK, OrganizationCommentListResponse{OrganizationID: orgID, Items: items})
}
