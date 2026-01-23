package admin

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupMenuRoutes configura rotas de menu para admin
// POST, PUT, DELETE - operações de escrita
func SetupMenuRoutes(r gin.IRouter) {
	menu := r.Group("/menu")
	{
		menu.POST("", resource.ServersControllers.SourceMenu.ServiceCreateMenu)
		menu.PUT("/:id", resource.ServersControllers.SourceMenu.ServiceUpdateMenu)
		menu.PUT("/:id/order", resource.ServersControllers.SourceMenu.ServiceUpdateMenuOrder)
		menu.PUT("/:id/status", resource.ServersControllers.SourceMenu.ServiceUpdateMenuStatus)
		menu.PUT("/:id/manual-override", resource.ServersControllers.SourceMenu.ServiceSetMenuAsManualOverride)
		menu.DELETE("/manual-override", resource.ServersControllers.SourceMenu.ServiceRemoveManualOverride)
		menu.DELETE("/:id", resource.ServersControllers.SourceMenu.ServiceDeleteMenu)
	}
}
