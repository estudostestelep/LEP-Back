package client

import (
	"lep/handler"
	"lep/middleware"
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupProductRoutes configura rotas de produto para client
// Todas as rotas requerem RolePermissionMiddleware
func SetupProductRoutes(r gin.IRouter) {
	product := r.Group("/product")
	{
		// Rotas de leitura - requer permissão client_products_view
		product.GET("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1),
			resource.ServersControllers.SourceProducts.ServiceGetProduct)
		product.GET("/purchase/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1),
			resource.ServersControllers.SourceProducts.ServiceGetProductByPurchase)
		product.GET("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1),
			resource.ServersControllers.SourceProducts.ServiceListProducts)
		product.GET("/by-tag",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1),
			resource.ServersControllers.SourceProducts.ServiceGetProductsByTag)

		// Rotas de escrita - requerem permissões específicas + limite de pacote
		product.POST("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_create", 1),
			middleware.PackageLimitMiddleware(resource.Handlers.HandlerLimits, handler.LimitProducts),
			resource.ServersControllers.SourceProducts.ServiceCreateProduct)
		product.PUT("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_edit", 1),
			resource.ServersControllers.SourceProducts.ServiceUpdateProduct)
		product.PUT("/:id/image",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_edit", 1),
			resource.ServersControllers.SourceProducts.ServiceUpdateProductImage)
		product.DELETE("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_delete", 1),
			resource.ServersControllers.SourceProducts.ServiceDeleteProduct)

		// Tag management - requer permissão de edição
		product.GET("/:id/tags",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1),
			resource.ServersControllers.SourceProducts.ServiceGetProductTags)
		product.POST("/:id/tags",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_edit", 1),
			resource.ServersControllers.SourceProducts.ServiceAddTagToProduct)
		product.DELETE("/:id/tags/:tagId",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_edit", 1),
			resource.ServersControllers.SourceProducts.ServiceRemoveTagFromProduct)

		// Order and status management - requer permissão de edição
		product.PUT("/:id/order",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_edit", 1),
			resource.ServersControllers.SourceProducts.ServiceUpdateProductOrder)
		product.PUT("/:id/status",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_edit", 1),
			resource.ServersControllers.SourceProducts.ServiceUpdateProductStatus)

		// Filtering by menu structure - requer permissão de visualização
		product.GET("/type/:type",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1),
			resource.ServersControllers.SourceProducts.ServiceGetProductsByType)
		product.GET("/category/:categoryId",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1),
			resource.ServersControllers.SourceProducts.ServiceGetProductsByCategory)
		product.GET("/subcategory/:subcategoryId",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1),
			resource.ServersControllers.SourceProducts.ServiceGetProductsBySubcategory)
	}
}
