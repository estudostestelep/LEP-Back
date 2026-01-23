package client

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupModuleReadRoutes configura rotas de leitura de módulo para client
func SetupModuleReadRoutes(r gin.IRouter) {
	module := r.Group("/module")
	{
		module.GET("", resource.ServersControllers.SourceRole.ListModules)
		module.GET("/with-permissions", resource.ServersControllers.SourceRole.ListModulesWithPermissions)
		module.GET("/available", resource.ServersControllers.SourceRole.GetOrganizationModules)
		module.GET("/:id", resource.ServersControllers.SourceRole.GetModule)
	}
}
