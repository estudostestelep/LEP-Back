package client

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configura rotas de usuário para client
func SetupUserRoutes(r gin.IRouter) {
	user := r.Group("/user")
	{
		user.GET("/:id", resource.ServersControllers.SourceUsers.ServiceGetUser)
		user.GET("/group/:id", resource.ServersControllers.SourceUsers.ServiceGetUserByGroup)
		user.GET("", resource.ServersControllers.SourceUsers.ServiceListUsers)
		user.PUT("/:id", resource.ServersControllers.SourceUsers.ServiceUpdateUser)
		user.DELETE("/:id", resource.ServersControllers.SourceUsers.ServiceDeleteUser)

		// Access Logs - Histórico de acessos do usuário
		user.GET("/:id/access-logs", resource.ServersControllers.SourceAccessLog.ServiceGetUserAccessLogs)
	}
}
