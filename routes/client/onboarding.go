package client

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupOnboardingRoutes configura rotas para status de onboarding
func SetupOnboardingRoutes(r gin.IRouter) {
	onboarding := r.Group("/onboarding")
	{
		onboarding.GET("/status", resource.ServersControllers.SourceOnboarding.GetOnboardingStatus)
	}
}
