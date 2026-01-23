package client

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupPermissionReadRoutes configura rotas de leitura de permissão para client
func SetupPermissionReadRoutes(r gin.IRouter) {
	permission := r.Group("/permission")
	{
		permission.GET("", resource.ServersControllers.SourceRole.ListPermissions)
		permission.GET("/:id", resource.ServersControllers.SourceRole.GetPermission)
	}
}
