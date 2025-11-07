package server

import (
	"lep/handler"
	"lep/repositories/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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

	// Buscar configurações existentes primeiro
	existingSettings, err := s.handler.GetSettingsByProject(projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching existing settings"})
		return
	}

	// Manter dados imutáveis
	updateData.ID = existingSettings.ID
	updateData.ProjectID = existingSettings.ProjectID
	updateData.OrganizationID = existingSettings.OrganizationID
	updateData.CreatedAt = existingSettings.CreatedAt

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
