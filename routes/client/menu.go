package client

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupMenuReadRoutes configura rotas de leitura de menu para client
func SetupMenuReadRoutes(r gin.IRouter) {
	menu := r.Group("/menu")
	{
		// Rotas específicas primeiro
		menu.GET("/active-now", resource.ServersControllers.SourceMenu.ServiceGetActiveMenu)
		menu.GET("/active", resource.ServersControllers.SourceMenu.ServiceListActiveMenus)
		menu.GET("/options", resource.ServersControllers.SourceMenu.ServiceGetMenuOptions)

		// Rotas genéricas
		menu.GET("/:id", resource.ServersControllers.SourceMenu.ServiceGetMenu)
		menu.GET("", resource.ServersControllers.SourceMenu.ServiceListMenus)
	}
}
