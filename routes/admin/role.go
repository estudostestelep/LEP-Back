package admin

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupRoleRoutes configura rotas de cargo para admin
// POST, PUT, DELETE - operações de escrita
func SetupRoleRoutes(r gin.IRouter) {
	role := r.Group("/role")
	{
		role.POST("", resource.ServersControllers.SourceRole.CreateRole)
		role.PUT("/:id", resource.ServersControllers.SourceRole.UpdateRole)
		role.DELETE("/:id", resource.ServersControllers.SourceRole.DeleteRole)
	}
}
