package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"2gis-calm-map/api/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Directories (relative to project root where binary runs). Adjust as needed.
const (
	mediaBaseDir  = "media"
	orgMapDir     = "media/org_maps"
	orgPictureDir = "media/org_pictures"
)

var allowedImageExt = map[string]struct{}{".png": {}, ".jpg": {}, ".jpeg": {}, ".webp": {}, ".gif": {}}

// ensureDirs makes sure storage directories exist.
func ensureDirs() error {
	for _, d := range []string{mediaBaseDir, orgMapDir, orgPictureDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return err
		}
	}
	return nil
}

// uploadOrganizationImage handles generic upload and model update.
func uploadOrganizationImage(c *gin.Context, kind string) {
	if err := ensureDirs(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// auth check (any authenticated user can view, but we restrict upload to owner/admin)
	roleVal, _ := c.Get("role")
	if roleVal != "admin" && roleVal != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	orgIDParam := c.Param("organization_id")
	var orgID uint
	if _, err := fmt.Sscan(orgIDParam, &orgID); err != nil || orgID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization_id"})
		return
	}

	org, err := orgService.GetByID(orgID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "organization not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if _, ok := allowedImageExt[ext]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported extension"})
		return
	}

	// unique filename
	fname := fmt.Sprintf("org_%d_%s_%d%s", org.ID, kind, time.Now().UnixNano(), ext)
	var targetDir string
	if kind == "map" {
		targetDir = orgMapDir
	} else {
		targetDir = orgPictureDir
	}
	fullPath := filepath.Join(targetDir, fname)

	out, err := os.Create(fullPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer out.Close()
	if _, err := io.Copy(out, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	relPath := fullPath // stored as relative path (already relative to run dir)
	updates := map[string]interface{}{}
	if kind == "map" {
		updates["map_path"] = relPath
	} else {
		updates["picture_path"] = relPath
	}

	// direct update via repository
	if err := repository.UpdateOrganizationFields(org.ID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"path": relPath, "type": kind})
}

// GetOrganizationImage serves a stored image (map or picture) by kind.
func GetOrganizationImage(c *gin.Context) {
	orgIDParam := c.Param("organization_id")
	var orgID uint
	if _, err := fmt.Sscan(orgIDParam, &orgID); err != nil || orgID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization_id"})
		return
	}
	kind := c.Param("kind") // map | picture

	org, err := orgService.GetByID(orgID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "organization not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var pathPtr *string
	if kind == "map" {
		pathPtr = org.MapPath
	} else if kind == "picture" {
		pathPtr = org.PicturePath
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kind"})
		return
	}
	if pathPtr == nil || *pathPtr == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not set"})
		return
	}

	// Serve file
	c.File(*pathPtr)
}

// Upload endpoints wrappers
// @Summary Upload organization map image
// @Description Загрузка файла карты организации (png/jpg/jpeg/webp/gif). Требует роль owner/admin.
// @Tags organization-media
// @Accept mpfd
// @Produce json
// @Security BearerAuth
// @Param organization_id path int true "Organization ID"
// @Param file formData file true "Image file"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /organization/{organization_id}/map/upload [post]
func UploadOrganizationMap(c *gin.Context) { uploadOrganizationImage(c, "map") }

// @Summary Upload organization picture image
// @Description Загрузка основной картинки организации. Требует роль owner/admin.
// @Tags organization-media
// @Accept mpfd
// @Produce json
// @Security BearerAuth
// @Param organization_id path int true "Organization ID"
// @Param file formData file true "Image file"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /organization/{organization_id}/picture/upload [post]
func UploadOrganizationPicture(c *gin.Context) { uploadOrganizationImage(c, "picture") }

// @Summary Get organization image (map or picture)
// @Description Возвращает файл изображения по типу (map | picture).
// @Tags organization-media
// @Produce octet-stream
// @Security BearerAuth
// @Param organization_id path int true "Organization ID"
// @Param kind path string true "Kind (map|picture)"
// @Success 200 {file} byte
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /organization/{organization_id}/image/{kind} [get]
func GetOrganizationImageHandler(c *gin.Context) { GetOrganizationImage(c) }
