package client

import (
	"lep/middleware"
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupCustomerRoutes configura rotas de cliente para client
// Todas as rotas requerem RolePermissionMiddleware
func SetupCustomerRoutes(r gin.IRouter) {
	customer := r.Group("/customer")
	{
		// Rotas de leitura
		customer.GET("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_customers_view", 1),
			resource.ServersControllers.SourceCustomer.ServiceGetCustomer)
		customer.GET("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_customers_view", 1),
			resource.ServersControllers.SourceCustomer.ServiceListCustomers)

		// Rotas de escrita
		customer.POST("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_customers_create", 1),
			resource.ServersControllers.SourceCustomer.ServiceCreateCustomer)
		customer.PUT("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_customers_edit", 1),
			resource.ServersControllers.SourceCustomer.ServiceUpdateCustomer)
		customer.DELETE("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_customers_delete", 1),
			resource.ServersControllers.SourceCustomer.ServiceDeleteCustomer)
	}
}
