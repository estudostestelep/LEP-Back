package client

import (
	"lep/middleware"
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupWaitlistRoutes configura rotas de fila de espera para client
// Requer módulo client_waitlist + RolePermissionMiddleware
func SetupWaitlistRoutes(r gin.IRouter) {
	waitlist := r.Group("/waitlist")
	{
		// Verificar se módulo de fila de espera está disponível
		waitlist.Use(middleware.ModuleRequiredMiddleware(resource.Handlers.HandlerLimits, "client_waitlist"))

		// Rotas de leitura
		waitlist.GET("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_waitlist_view", 1),
			resource.ServersControllers.SourceWaitlist.ServiceGetWaitlist)
		waitlist.GET("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_waitlist_view", 1),
			resource.ServersControllers.SourceWaitlist.ServiceListWaitlists)

		// Rotas de escrita
		waitlist.POST("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_waitlist_create", 1),
			resource.ServersControllers.SourceWaitlist.ServiceCreateWaitlist)
		waitlist.PUT("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_waitlist_edit", 1),
			resource.ServersControllers.SourceWaitlist.ServiceUpdateWaitlist)
		waitlist.DELETE("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_waitlist_delete", 1),
			resource.ServersControllers.SourceWaitlist.ServiceDeleteWaitlist)
	}
}
