package client

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupPlanChangeClientRoutes configura rotas de solicitações de mudança de plano para client
func SetupPlanChangeClientRoutes(r gin.IRouter) {
	planChange := r.Group("/plan-change-request")
	{
		planChange.POST("", resource.ServersControllers.SourcePlanChangeRequest.CreateRequest)
		planChange.GET("/my-requests", resource.ServersControllers.SourcePlanChangeRequest.GetMyRequests)
		planChange.GET("/:id", resource.ServersControllers.SourcePlanChangeRequest.GetRequestById)
		planChange.POST("/:id/cancel", resource.ServersControllers.SourcePlanChangeRequest.CancelRequest)
	}
}
