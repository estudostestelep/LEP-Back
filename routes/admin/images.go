package admin

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupImagesRoutes configura rotas de gerenciamento de imagens para admin
func SetupImagesRoutes(r gin.IRouter) {
	images := r.Group("/images")
	{
		images.POST("/cleanup", resource.ServersControllers.SourceImageManagement.ServiceCleanupOrphanedFiles)
		images.GET("/stats", resource.ServersControllers.SourceImageManagement.ServiceGetImageStats)
	}
}
