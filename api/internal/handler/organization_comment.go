package handler

import (
	"net/http"

	"2gis-calm-map/api/internal/model"
	"2gis-calm-map/api/internal/repository"
	"2gis-calm-map/api/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var orgCommentService = service.NewOrganizationCommentService()

// OrganizationCommentCreateRequest swagger request model
// Allows providing any subset of values/comments; zeros or omitted numeric values are ignored (only non-nil & >0 values aggregate).
type OrganizationCommentCreateRequest struct {
	OrganizationID       uint    `json:"organization_id" binding:"required"`
	Text                 *string `json:"text"`
	AppearanceValue      *uint   `json:"appearance_value"`
	AppearanceComment    *string `json:"appearance_comment"`
	LightingValue        *uint   `json:"lighting_value"`
	LightingComment      *string `json:"lighting_comment"`
	SmellValue           *uint   `json:"smell_value"`
	SmellComment         *string `json:"smell_comment"`
	TemperatureValue     *uint   `json:"temperature_value"`
	TemperatureComment   *string `json:"temperature_comment"`
	TactilityValue       *uint   `json:"tactility_value"`
	TactilityComment     *string `json:"tactility_comment"`
	SignageValue         *uint   `json:"signage_value"`
	SignageComment       *string `json:"signage_comment"`
	IntuitivenessValue   *uint   `json:"intuitiveness_value"`
	IntuitivenessComment *string `json:"intuitiveness_comment"`
	StaffAttitudeValue   *uint   `json:"staff_attitude_value"`
	StaffAttitudeComment *string `json:"staff_attitude_comment"`
	PeopleDensityValue   *uint   `json:"people_density_value"`
	PeopleDensityComment *string `json:"people_density_comment"`
	SelfServiceValue     *uint   `json:"self_service_value"`
	SelfServiceComment   *string `json:"self_service_comment"`
	CalmnessValue        *uint   `json:"calmness_value"`
	CalmnessComment      *string `json:"calmness_comment"`
}

// OrganizationCommentCreateResponse response with created comment and updated aggregates snapshot
type OrganizationCommentCreateResponse struct {
	Comment *model.OrganizationComment `json:"comment"`
	Updated *model.OrganizationParams  `json:"updated_aggregates"`
}

// CreateOrganizationComment godoc
// @Summary Create comment for organization with optional parameter ratings
// @Description Создаёт комментарий. user_id берётся из токена автоматически. Каждый не-nil и >0 value обновляет агрегаты (sum,count,avg).
// @Tags organization-comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body OrganizationCommentCreateRequest true "Create comment"
// @Success 201 {object} OrganizationCommentCreateResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /organization/comment [post]
func CreateOrganizationComment(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	tokenUserID := userIDVal.(uint)

	var req OrganizationCommentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure organization exists
	if _, err := orgService.GetByID(req.OrganizationID); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "organization not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// fetch or create params aggregate
	paramsModel, err := orgParamsService.GetOrCreate(req.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	comment := &model.OrganizationComment{
		OrganizationID:       req.OrganizationID,
		UserID:               tokenUserID,
		Text:                 req.Text,
		AppearanceValue:      req.AppearanceValue,
		AppearanceComment:    req.AppearanceComment,
		LightingValue:        req.LightingValue,
		LightingComment:      req.LightingComment,
		SmellValue:           req.SmellValue,
		SmellComment:         req.SmellComment,
		TemperatureValue:     req.TemperatureValue,
		TemperatureComment:   req.TemperatureComment,
		TactilityValue:       req.TactilityValue,
		TactilityComment:     req.TactilityComment,
		SignageValue:         req.SignageValue,
		SignageComment:       req.SignageComment,
		IntuitivenessValue:   req.IntuitivenessValue,
		IntuitivenessComment: req.IntuitivenessComment,
		StaffAttitudeValue:   req.StaffAttitudeValue,
		StaffAttitudeComment: req.StaffAttitudeComment,
		PeopleDensityValue:   req.PeopleDensityValue,
		PeopleDensityComment: req.PeopleDensityComment,
		SelfServiceValue:     req.SelfServiceValue,
		SelfServiceComment:   req.SelfServiceComment,
		CalmnessValue:        req.CalmnessValue,
		CalmnessComment:      req.CalmnessComment,
	}

	// Compute average value across provided non-nil and >0 parameter values
	var sum uint
	var count uint
	acc := func(v *uint) {
		if v != nil && *v > 0 {
			sum += *v
			count++
		}
	}
	acc(req.AppearanceValue)
	acc(req.LightingValue)
	acc(req.SmellValue)
	acc(req.TemperatureValue)
	acc(req.TactilityValue)
	acc(req.SignageValue)
	acc(req.IntuitivenessValue)
	acc(req.StaffAttitudeValue)
	acc(req.PeopleDensityValue)
	acc(req.SelfServiceValue)
	acc(req.CalmnessValue)
	if count > 0 {
		avg := float64(sum) / float64(count)
		comment.AvgValue = &avg
	}

	if err := orgCommentService.CreateWithAggregation(comment, &paramsModel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// reload updated params
	updated, err := repository.GetOrganizationParams(req.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, OrganizationCommentCreateResponse{Comment: comment, Updated: &updated})
}
