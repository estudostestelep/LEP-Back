package admin

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupSidebarConfigRoutes configura rotas de configuração da sidebar para admin
func SetupSidebarConfigRoutes(r gin.IRouter) {
	sidebarConfig := r.Group("/sidebar-config")
	{
		sidebarConfig.PUT("", resource.ServersControllers.SourceSidebarConfig.UpdateConfig)
		sidebarConfig.POST("/reset", resource.ServersControllers.SourceSidebarConfig.ResetConfig)
	}
}
