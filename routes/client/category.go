package client

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupCategoryReadRoutes configura rotas de leitura de categoria para client
func SetupCategoryReadRoutes(r gin.IRouter) {
	category := r.Group("/category")
	{
		category.GET("/:id", resource.ServersControllers.SourceCategory.ServiceGetCategory)
		category.GET("", resource.ServersControllers.SourceCategory.ServiceListCategories)
		category.GET("/active", resource.ServersControllers.SourceCategory.ServiceListActiveCategories)
		category.GET("/menu/:menuId", resource.ServersControllers.SourceCategory.ServiceGetCategoriesByMenu)
	}
}
