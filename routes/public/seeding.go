package public

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupSeedingRoutes configura rotas para seeding (desenvolvimento)
func SetupSeedingRoutes(r *gin.Engine) {
	// Rotas para seeding
	r.POST("/organization", resource.ServersControllers.SourceOrganization.CreateOrganization)
	r.POST("/project", resource.ServersControllers.SourceProject.CreateProject)
}
