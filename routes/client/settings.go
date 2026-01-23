package client

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupSettingsRoutes configura rotas de configurações para client
func SetupSettingsRoutes(r gin.IRouter) {
	// Settings gerais
	settings := r.Group("/settings")
	{
		settings.GET("", resource.ServersControllers.SourceSettings.GetSettingsByProject)
		settings.PUT("", resource.ServersControllers.SourceSettings.UpdateSettings)
	}

	// Display Settings
	displaySettings := r.Group("/project/settings/display")
	{
		displaySettings.GET("", resource.ServersControllers.SourceDisplaySettings.GetDisplaySettings)
		displaySettings.PUT("", resource.ServersControllers.SourceDisplaySettings.UpdateDisplaySettings)
		displaySettings.POST("/reset", resource.ServersControllers.SourceDisplaySettings.ResetDisplaySettings)
	}

	// Theme Customization
	theme := r.Group("/project/settings/theme")
	{
		theme.GET("", resource.ServersControllers.SourceThemeCustomization.GetTheme)
		theme.POST("", resource.ServersControllers.SourceThemeCustomization.CreateOrUpdateTheme)
		theme.PUT("", resource.ServersControllers.SourceThemeCustomization.CreateOrUpdateTheme)
		theme.POST("/reset", resource.ServersControllers.SourceThemeCustomization.ResetTheme)
		theme.DELETE("", resource.ServersControllers.SourceThemeCustomization.DeleteTheme)
	}

	// Sidebar Config (leitura apenas - escrita é admin)
	sidebarConfig := r.Group("/sidebar-config")
	{
		sidebarConfig.GET("", resource.ServersControllers.SourceSidebarConfig.GetConfig)
	}
}
