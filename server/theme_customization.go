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
		// ==================== CORES PRINCIPAIS - LIGHT MODE (7 campos) ====================
		PrimaryColorLight:        themeRequest.PrimaryColorLight,
		SecondaryColorLight:      themeRequest.SecondaryColorLight,
		BackgroundColorLight:     themeRequest.BackgroundColorLight,
		CardBackgroundColorLight: themeRequest.CardBackgroundColorLight,
		TextColorLight:           themeRequest.TextColorLight,
		TextSecondaryColorLight:  themeRequest.TextSecondaryColorLight,
		AccentColorLight:         themeRequest.AccentColorLight,

		// ==================== CORES PRINCIPAIS - DARK MODE (7 campos) ====================
		PrimaryColorDark:        themeRequest.PrimaryColorDark,
		SecondaryColorDark:      themeRequest.SecondaryColorDark,
		BackgroundColorDark:     themeRequest.BackgroundColorDark,
		CardBackgroundColorDark: themeRequest.CardBackgroundColorDark,
		TextColorDark:           themeRequest.TextColorDark,
		TextSecondaryColorDark:  themeRequest.TextSecondaryColorDark,
		AccentColorDark:         themeRequest.AccentColorDark,

		// ==================== CORES SEMÂNTICAS - LIGHT MODE (5 campos) ====================
		DestructiveColorLight: themeRequest.DestructiveColorLight,
		SuccessColorLight:     themeRequest.SuccessColorLight,
		WarningColorLight:     themeRequest.WarningColorLight,
		BorderColorLight:      themeRequest.BorderColorLight,
		PriceColorLight:       themeRequest.PriceColorLight,

		// ==================== CORES SEMÂNTICAS - DARK MODE (5 campos) ====================
		DestructiveColorDark: themeRequest.DestructiveColorDark,
		SuccessColorDark:     themeRequest.SuccessColorDark,
		WarningColorDark:     themeRequest.WarningColorDark,
		BorderColorDark:      themeRequest.BorderColorDark,
		PriceColorDark:       themeRequest.PriceColorDark,

		// ==================== SISTEMA - LIGHT MODE (2 campos) ====================
		FocusRingColorLight:      themeRequest.FocusRingColorLight,
		InputBackgroundColorLight: themeRequest.InputBackgroundColorLight,

		// ==================== SISTEMA - DARK MODE (2 campos) ====================
		FocusRingColorDark:       themeRequest.FocusRingColorDark,
		InputBackgroundColorDark: themeRequest.InputBackgroundColorDark,

		// ==================== CONFIGURAÇÕES NUMÉRICAS (2 campos) ====================
		DisabledOpacity: themeRequest.DisabledOpacity,
		ShadowIntensity: themeRequest.ShadowIntensity,

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
