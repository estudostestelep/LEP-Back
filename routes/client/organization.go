package client

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupOrganizationRoutes configura rotas de organização para client
func SetupOrganizationRoutes(r gin.IRouter) {
	organization := r.Group("/organization")
	{
		organization.GET("/:id", resource.ServersControllers.SourceOrganization.GetOrganizationById)
		organization.GET("", resource.ServersControllers.SourceOrganization.ListOrganizations)
		organization.GET("/active", resource.ServersControllers.SourceOrganization.ListActiveOrganizations)
		organization.GET("/email", resource.ServersControllers.SourceOrganization.GetOrganizationByEmail)
		organization.PUT("/:id", resource.ServersControllers.SourceOrganization.UpdateOrganization)
		organization.DELETE("/:id", resource.ServersControllers.SourceOrganization.SoftDeleteOrganization)
		organization.DELETE("/:id/permanent", resource.ServersControllers.SourceOrganization.HardDeleteOrganization)
	}
}
