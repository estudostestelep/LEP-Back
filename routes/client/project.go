package client

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupProjectRoutes configura rotas de projeto para client
func SetupProjectRoutes(r gin.IRouter) {
	project := r.Group("/project")
	{
		project.GET("/:id", resource.ServersControllers.SourceProject.GetProjectById)
		project.GET("", resource.ServersControllers.SourceProject.GetProjectsByOrganization)
		project.GET("/organization/:orgId", resource.ServersControllers.SourceProject.GetProjectsByOrganizationId)
		project.GET("/active", resource.ServersControllers.SourceProject.GetActiveProjects)
		project.PUT("/:id", resource.ServersControllers.SourceProject.UpdateProject)
		project.DELETE("/:id", resource.ServersControllers.SourceProject.SoftDeleteProject)
		project.DELETE("/:id/permanent", resource.ServersControllers.SourceProject.HardDeleteProject)
	}

	// User-Project routes
	userProj := r.Group("/user-project")
	{
		userProj.DELETE("/user/:userId/proj/:projectId", resource.ServersControllers.SourceUserProject.ServiceRemoveUserFromProject)
		userProj.PUT("/:id", resource.ServersControllers.SourceUserProject.ServiceUpdateUserProject)
		userProj.GET("/user/:userId", resource.ServersControllers.SourceUserProject.ServiceGetUserProjects)
		userProj.GET("/user/:userId/org/:orgId", resource.ServersControllers.SourceUserProject.ServiceGetUserProjectsByOrganization)
		userProj.GET("/proj/:projectId", resource.ServersControllers.SourceUserProject.ServiceGetProjectUsers)
	}
}
