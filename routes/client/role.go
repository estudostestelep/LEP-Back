package client

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupRoleClientRoutes configura rotas de cargo para client
// Leitura + atribuição/remoção de cargos
func SetupRoleClientRoutes(r gin.IRouter) {
	role := r.Group("/role")
	{
		// Rotas específicas primeiro
		role.GET("/system", resource.ServersControllers.SourceRole.ListSystemRoles)
		role.GET("/check", resource.ServersControllers.SourceRole.CheckPermission)
		role.GET("/my-permissions", resource.ServersControllers.SourceRole.GetMyPermissions)

		// Atribuição de cargos a usuários
		role.POST("/assign", resource.ServersControllers.SourceRole.AssignRoleToUser)
		role.POST("/remove", resource.ServersControllers.SourceRole.RemoveRoleFromUser)
		role.GET("/user/:userId", resource.ServersControllers.SourceRole.GetUserRoles)
		role.GET("/user/:userId/details", resource.ServersControllers.SourceRole.GetUserRolesWithDetails)

		// Níveis de permissão por cargo
		role.POST("/permission-level", resource.ServersControllers.SourceRole.SetPermissionLevel)

		// Leitura de cargos
		role.GET("", resource.ServersControllers.SourceRole.ListRoles)
		role.GET("/:id", resource.ServersControllers.SourceRole.GetRole)
		role.GET("/:id/permissions", resource.ServersControllers.SourceRole.GetRolePermissions)
	}
}
