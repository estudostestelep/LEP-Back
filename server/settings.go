package server

import (
	"lep/handler"
	"lep/repositories/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type SettingsServer struct {
	handler handler.ISettingsHandler
}

type ISettingsServer interface {
	GetSettingsByProject(c *gin.Context)
	UpdateSettings(c *gin.Context)
}

func NewSettingsServer(handler handler.ISettingsHandler) ISettingsServer {
	return &SettingsServer{handler: handler}
}

// GetSettingsByProject busca configurações por projeto
func (s *SettingsServer) GetSettingsByProject(c *gin.Context) {
	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Organization-Id' cannot be empty",
		})
		return
	}

	projectId := c.GetHeader("X-Lpe-Project-Id")
	if strings.TrimSpace(projectId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Project-Id' cannot be empty",
		})
		return
	}

	settings, err := s.handler.GetOrCreateSettings(organizationId, projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateSettings atualiza configurações
func (s *SettingsServer) UpdateSettings(c *gin.Context) {
	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Organization-Id' cannot be empty",
		})
		return
	}

	projectId := c.GetHeader("X-Lpe-Project-Id")
	if strings.TrimSpace(projectId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Project-Id' cannot be empty",
		})
		return
	}

	// Buscar configurações existentes ou criar padrão
	existingSettings, err := s.handler.GetOrCreateSettings(organizationId, projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching settings"})
		return
	}

	var updateData models.Settings
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validações básicas
	if updateData.MinAdvanceHours < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "MinAdvanceHours must be >= 0"})
		return
	}
	if updateData.MaxAdvanceDays < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "MaxAdvanceDays must be >= 1"})
		return
	}
	if updateData.DefaultNotificationChannel != "sms" &&
	   updateData.DefaultNotificationChannel != "email" &&
	   updateData.DefaultNotificationChannel != "whatsapp" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "DefaultNotificationChannel must be sms, email, or whatsapp"})
		return
	}

	// Manter dados imutáveis
	updateData.Id = existingSettings.Id
	updateData.OrganizationId = existingSettings.OrganizationId
	updateData.ProjectId = existingSettings.ProjectId
	updateData.CreatedAt = existingSettings.CreatedAt

	err = s.handler.UpdateSettings(&updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating settings"})
		return
	}

	c.JSON(http.StatusOK, updateData)
}