package client

import (
	"lep/middleware"
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupTagRoutes configura rotas de tag para client
// Todas as rotas requerem RolePermissionMiddleware
func SetupTagRoutes(r gin.IRouter) {
	tag := r.Group("/tag")
	{
		// Rotas de leitura - requer permissão client_tags_view (nível 1)
		tag.GET("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tags_view", 1),
			resource.ServersControllers.SourceTag.ServiceGetTag)
		tag.GET("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tags_view", 1),
			resource.ServersControllers.SourceTag.ServiceListTags)
		tag.GET("/active",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tags_view", 1),
			resource.ServersControllers.SourceTag.ServiceListActiveTags)
		tag.GET("/entity/:entityType",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tags_view", 1),
			resource.ServersControllers.SourceTag.ServiceGetTagsByEntityType)

		// Rotas de escrita - requerem permissões específicas (nível 1 = habilitado)
		tag.POST("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tags_create", 1),
			resource.ServersControllers.SourceTag.ServiceCreateTag)
		tag.PUT("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tags_edit", 1),
			resource.ServersControllers.SourceTag.ServiceUpdateTag)
		tag.DELETE("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tags_delete", 1),
			resource.ServersControllers.SourceTag.ServiceDeleteTag)
	}
}
