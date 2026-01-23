package admin

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupModuleRoutes configura rotas de módulo para admin
// POST, PUT, DELETE - operações de escrita
func SetupModuleRoutes(r gin.IRouter) {
	module := r.Group("/module")
	{
		module.POST("", resource.ServersControllers.SourceRole.CreateModule)
		module.PUT("/:id", resource.ServersControllers.SourceRole.UpdateModule)
		module.DELETE("/:id", resource.ServersControllers.SourceRole.DeleteModule)
	}
}
