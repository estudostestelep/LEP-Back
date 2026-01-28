package admin

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupClientUserRoutes configura rotas de CRUD para clients (tabela clients)
// Admins podem gerenciar clients de qualquer organização
func SetupClientUserRoutes(r gin.IRouter) {
	clientUser := r.Group("/client-user")
	{
		clientUser.GET("", resource.ServersControllers.SourceClientUsers.ServiceListClients)
		clientUser.GET("/:id", resource.ServersControllers.SourceClientUsers.ServiceGetClient)
		clientUser.POST("", resource.ServersControllers.SourceClientUsers.ServiceCreateClient)
		clientUser.PUT("/:id", resource.ServersControllers.SourceClientUsers.ServiceUpdateClient)
		clientUser.DELETE("/:id", resource.ServersControllers.SourceClientUsers.ServiceDeleteClient)
	}
}
