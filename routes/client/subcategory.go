package client

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupSubcategoryReadRoutes configura rotas de leitura de subcategoria para client
func SetupSubcategoryReadRoutes(r gin.IRouter) {
	subcategory := r.Group("/subcategory")
	{
		subcategory.GET("/:id", resource.ServersControllers.SourceSubcategory.ServiceGetSubcategory)
		subcategory.GET("", resource.ServersControllers.SourceSubcategory.ServiceListSubcategories)
		subcategory.GET("/active", resource.ServersControllers.SourceSubcategory.ServiceListActiveSubcategories)
		subcategory.GET("/category/:categoryId", resource.ServersControllers.SourceSubcategory.ServiceGetSubcategoriesByCategory)
		subcategory.GET("/:id/categories", resource.ServersControllers.SourceSubcategory.ServiceGetSubcategoryCategories)
	}
}
