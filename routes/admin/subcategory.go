package admin

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupSubcategoryRoutes configura rotas de subcategoria para admin
// POST, PUT, DELETE - operações de escrita
func SetupSubcategoryRoutes(r gin.IRouter) {
	subcategory := r.Group("/subcategory")
	{
		subcategory.POST("", resource.ServersControllers.SourceSubcategory.ServiceCreateSubcategory)
		subcategory.PUT("/:id", resource.ServersControllers.SourceSubcategory.ServiceUpdateSubcategory)
		subcategory.PUT("/:id/order", resource.ServersControllers.SourceSubcategory.ServiceUpdateSubcategoryOrder)
		subcategory.PUT("/:id/status", resource.ServersControllers.SourceSubcategory.ServiceUpdateSubcategoryStatus)
		subcategory.DELETE("/:id", resource.ServersControllers.SourceSubcategory.ServiceDeleteSubcategory)

		// Category relationship management
		subcategory.POST("/:id/category/:categoryId", resource.ServersControllers.SourceSubcategory.ServiceAddCategoryToSubcategory)
		subcategory.DELETE("/:id/category/:categoryId", resource.ServersControllers.SourceSubcategory.ServiceRemoveCategoryFromSubcategory)
	}
}
