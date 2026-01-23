package admin

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupCategoryRoutes configura rotas de categoria para admin
// POST, PUT, DELETE - operações de escrita
func SetupCategoryRoutes(r gin.IRouter) {
	category := r.Group("/category")
	{
		category.POST("", resource.ServersControllers.SourceCategory.ServiceCreateCategory)
		category.PUT("/:id", resource.ServersControllers.SourceCategory.ServiceUpdateCategory)
		category.PUT("/:id/order", resource.ServersControllers.SourceCategory.ServiceUpdateCategoryOrder)
		category.PUT("/:id/status", resource.ServersControllers.SourceCategory.ServiceUpdateCategoryStatus)
		category.DELETE("/:id", resource.ServersControllers.SourceCategory.ServiceDeleteCategory)
	}
}
