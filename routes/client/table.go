package client

import (
	"lep/handler"
	"lep/middleware"
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupTableRoutes configura rotas de mesa para client
// Todas as rotas requerem RolePermissionMiddleware
func SetupTableRoutes(r gin.IRouter) {
	table := r.Group("/table")
	{
		// Rotas de leitura - requer permissão client_tables_view
		table.GET("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tables_view", 1),
			resource.ServersControllers.SourceTables.ServiceGetTable)
		table.GET("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tables_view", 1),
			resource.ServersControllers.SourceTables.ServiceListTables)

		// Rotas de escrita - requerem permissões específicas + limite de pacote
		table.POST("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tables_create", 1),
			middleware.PackageLimitMiddleware(resource.Handlers.HandlerLimits, handler.LimitTables),
			resource.ServersControllers.SourceTables.ServiceCreateTable)
		table.PUT("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tables_edit", 1),
			resource.ServersControllers.SourceTables.ServiceUpdateTable)
		table.DELETE("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tables_delete", 1),
			resource.ServersControllers.SourceTables.ServiceDeleteTable)
	}
}
