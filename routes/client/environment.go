package client

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupEnvironmentRoutes configura rotas de ambiente para client
func SetupEnvironmentRoutes(r gin.IRouter) {
	environment := r.Group("/environment")
	{
		environment.GET("/:id", resource.ServersControllers.SourceEnvironment.GetEnvironmentById)
		environment.GET("", resource.ServersControllers.SourceEnvironment.GetEnvironmentsByProject)
		environment.GET("/active", resource.ServersControllers.SourceEnvironment.GetActiveEnvironments)
		environment.POST("", resource.ServersControllers.SourceEnvironment.CreateEnvironment)
		environment.PUT("/:id", resource.ServersControllers.SourceEnvironment.UpdateEnvironment)
		environment.DELETE("/:id", resource.ServersControllers.SourceEnvironment.SoftDeleteEnvironment)
	}
}
