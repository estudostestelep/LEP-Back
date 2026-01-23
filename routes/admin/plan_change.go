package admin

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupPlanChangeRoutes configura rotas de solicitações de mudança de plano para admin
func SetupPlanChangeRoutes(r gin.IRouter) {
	planChange := r.Group("/plan-change-request")
	{
		planChange.GET("", resource.ServersControllers.SourcePlanChangeRequest.GetAllRequests)
		planChange.GET("/pending", resource.ServersControllers.SourcePlanChangeRequest.GetPendingRequests)
		planChange.GET("/organization/:orgId", resource.ServersControllers.SourcePlanChangeRequest.GetRequestsByOrganization)
		planChange.POST("/:id/approve", resource.ServersControllers.SourcePlanChangeRequest.ApproveRequest)
		planChange.POST("/:id/reject", resource.ServersControllers.SourcePlanChangeRequest.RejectRequest)
	}
}
