package admin

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupPermissionRoutes configura rotas de permissão para admin
// POST, PUT, DELETE - operações de escrita
func SetupPermissionRoutes(r gin.IRouter) {
	permission := r.Group("/permission")
	{
		permission.POST("", resource.ServersControllers.SourceRole.CreatePermission)
		permission.PUT("/:id", resource.ServersControllers.SourceRole.UpdatePermission)
		permission.DELETE("/:id", resource.ServersControllers.SourceRole.DeletePermission)
	}
}
