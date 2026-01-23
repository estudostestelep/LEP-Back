package admin

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configura rotas de usuário para admin
// CRUD de usuários administrativos + Gerenciamento de acesso
func SetupUserRoutes(r gin.IRouter) {
	user := r.Group("/user")
	{
		// CRUD de usuários admin
		user.GET("", resource.ServersControllers.SourceUsers.ServiceListUsers)
		user.GET("/:id", resource.ServersControllers.SourceUsers.ServiceGetUser)
		user.POST("", resource.ServersControllers.SourceUsers.ServiceCreateUser)
		user.PUT("/:id", resource.ServersControllers.SourceUsers.ServiceUpdateUser)
		user.DELETE("/:id", resource.ServersControllers.SourceUsers.ServiceDeleteUser)

		// User Access Management (organizações e projetos)
		user.GET("/:id/organizations-projects", resource.ServersControllers.SourceUserAccess.ServiceGetUserOrganizationsAndProjects)
		user.POST("/:id/organizations-projects", resource.ServersControllers.SourceUserAccess.ServiceUpdateUserOrganizationsAndProjects)
	}
}
