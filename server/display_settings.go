package server

import (
	"lep/handler"
	"lep/repositories/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DisplaySettingsServer struct {
	handler handler.IDisplaySettingsHandler
}

type IDisplaySettingsServer interface {
	GetDisplaySettings(c *gin.Context)
	UpdateDisplaySettings(c *gin.Context)
	ResetDisplaySettings(c *gin.Context)
}

func NewDisplaySettingsServer(handler handler.IDisplaySettingsHandler) IDisplaySettingsServer {
	return &DisplaySettingsServer{handler: handler}
}

// GetDisplaySettings busca configurações de exibição de produtos
func (s *DisplaySettingsServer) GetDisplaySettings(c *gin.Context) {
	projectId := c.GetHeader("X-Lpe-Project-Id")
	if strings.TrimSpace(projectId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Project-Id' cannot be empty",
		})
		return
	}

	settings, err := s.handler.GetSettingsByProject(projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching display settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateDisplaySettings atualiza configurações de exibição
func (s *DisplaySettingsServer) UpdateDisplaySettings(c *gin.Context) {
	projectId := c.GetHeader("X-Lpe-Project-Id")
	if strings.TrimSpace(projectId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Project-Id' cannot be empty",
		})
		return
	}

	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Organization-Id' cannot be empty",
		})
		return
	}

	var updateData models.ProjectDisplaySettings
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Parse IDs from headers
	projectUUID, err := uuid.Parse(projectId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	orgUUID, err := uuid.Parse(organizationId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID format"})
		return
	}

	// Buscar configurações existentes primeiro
	existingSettings, err := s.handler.GetSettingsByProject(projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching existing settings"})
		return
	}

	// Manter dados imutáveis e garantir IDs corretos do header
	// Se ID for uuid.Nil (novo registro), gera um novo
	if existingSettings.ID == uuid.Nil {
		updateData.ID = uuid.New()
		updateData.CreatedAt = time.Now()
	} else {
		updateData.ID = existingSettings.ID
		updateData.CreatedAt = existingSettings.CreatedAt
	}
	updateData.ProjectID = projectUUID
	updateData.OrganizationID = orgUUID

	err = s.handler.UpdateSettings(&updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating display settings"})
		return
	}

	c.JSON(http.StatusOK, updateData)
}

// ResetDisplaySettings reseta configurações para valores padrão
func (s *DisplaySettingsServer) ResetDisplaySettings(c *gin.Context) {
	projectId := c.GetHeader("X-Lpe-Project-Id")
	if strings.TrimSpace(projectId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Project-Id' cannot be empty",
		})
		return
	}

	settings, err := s.handler.ResetToDefaults(projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error resetting display settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}
