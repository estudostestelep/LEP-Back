package server

import (
	"lep/handler"
	"lep/repositories/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ThemeCustomizationServer struct {
	handler handler.IThemeCustomizationHandler
}

type IThemeCustomizationServer interface {
	GetTheme(c *gin.Context)
	CreateOrUpdateTheme(c *gin.Context)
	ResetTheme(c *gin.Context)
	DeleteTheme(c *gin.Context)
}

func NewThemeCustomizationServer(handler handler.IThemeCustomizationHandler) IThemeCustomizationServer {
	return &ThemeCustomizationServer{handler: handler}
}

// GetTheme busca customização de tema do projeto
func (s *ThemeCustomizationServer) GetTheme(c *gin.Context) {
	projectId := c.GetHeader("X-Lpe-Project-Id")
	if strings.TrimSpace(projectId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Project-Id' cannot be empty",
		})
		return
	}

	theme, err := s.handler.GetThemeByProject(projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching theme customization"})
		return
	}

	c.JSON(http.StatusOK, theme)
}

// CreateOrUpdateTheme cria ou atualiza customização de tema
func (s *ThemeCustomizationServer) CreateOrUpdateTheme(c *gin.Context) {
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

	var themeRequest models.ThemeCustomizationRequest
	if err := c.ShouldBindJSON(&themeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Converter DTO para model
	themeData := &models.ThemeCustomization{
		// ==================== CORES PRINCIPAIS (7 campos) ====================
		PrimaryColor:        themeRequest.PrimaryColor,
		SecondaryColor:      themeRequest.SecondaryColor,
		BackgroundColor:     themeRequest.BackgroundColor,
		CardBackgroundColor: themeRequest.CardBackgroundColor,
		TextColor:           themeRequest.TextColor,
		TextSecondaryColor:  themeRequest.TextSecondaryColor,
		AccentColor:         themeRequest.AccentColor,

		// ==================== CORES SEMÂNTICAS (4 novos campos) ====================
		DestructiveColor: themeRequest.DestructiveColor,
		SuccessColor:     themeRequest.SuccessColor,
		WarningColor:     themeRequest.WarningColor,
		BorderColor:      themeRequest.BorderColor,

		// ==================== CONFIGURAÇÕES DO SISTEMA (4 novos campos) ====================
		DisabledOpacity:      themeRequest.DisabledOpacity,
		FocusRingColor:       themeRequest.FocusRingColor,
		InputBackgroundColor: themeRequest.InputBackgroundColor,
		ShadowIntensity:      themeRequest.ShadowIntensity,

		IsActive: themeRequest.IsActive,
	}

	theme, err := s.handler.CreateOrUpdateTheme(projectId, organizationId, themeData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, theme)
}

// ResetTheme reseta tema para valores padrão
func (s *ThemeCustomizationServer) ResetTheme(c *gin.Context) {
	projectId := c.GetHeader("X-Lpe-Project-Id")
	if strings.TrimSpace(projectId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Project-Id' cannot be empty",
		})
		return
	}

	theme, err := s.handler.ResetToDefaults(projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error resetting theme"})
		return
	}

	c.JSON(http.StatusOK, theme)
}

// DeleteTheme deleta customização de tema
func (s *ThemeCustomizationServer) DeleteTheme(c *gin.Context) {
	projectId := c.GetHeader("X-Lpe-Project-Id")
	if strings.TrimSpace(projectId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Project-Id' cannot be empty",
		})
		return
	}

	err := s.handler.DeleteTheme(projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting theme"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Theme deleted successfully"})
}
