package client

import (
	"lep/middleware"
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupOrderRoutes configura rotas de pedido para client
// Todas as rotas requerem RolePermissionMiddleware
func SetupOrderRoutes(r gin.IRouter) {
	order := r.Group("/order")
	{
		// Rotas de leitura
		order.GET("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_orders_view", 1),
			resource.ServersControllers.SourceOrders.GetOrderById)
		order.GET("/:id/progress",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_orders_view", 1),
			resource.ServersControllers.SourceOrders.GetOrderProgress)
		order.GET("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_orders_view", 1),
			resource.ServersControllers.SourceOrders.ListOrders)

		// Rotas de escrita
		order.POST("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_orders_create", 1),
			resource.ServersControllers.SourceOrders.CreateOrder)
		order.PUT("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_orders_edit", 1),
			resource.ServersControllers.SourceOrders.UpdateOrder)
		order.PUT("/:id/status",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_orders_edit", 1),
			resource.ServersControllers.SourceOrders.UpdateOrderStatus)
		order.DELETE("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_orders_delete", 1),
			resource.ServersControllers.SourceOrders.SoftDeleteOrder)
	}

	// Kitchen specific routes
	kitchen := r.Group("/kitchen")
	{
		kitchen.GET("/queue", resource.ServersControllers.SourceOrders.GetKitchenQueue)
	}
}
