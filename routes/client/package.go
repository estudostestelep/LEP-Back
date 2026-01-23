package client

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupPackageClientRoutes configura rotas de pacote para client
// Leitura + assinatura
func SetupPackageClientRoutes(r gin.IRouter) {
	pkg := r.Group("/package")
	{
		// Rotas específicas primeiro
		pkg.GET("/subscription", resource.ServersControllers.SourceRole.GetOrganizationSubscription)
		pkg.GET("/usage-limits", resource.ServersControllers.SourceRole.GetUsageAndLimits)
		pkg.POST("/subscribe", resource.ServersControllers.SourceRole.SubscribeOrganization)

		// Leitura de pacotes
		pkg.GET("", resource.ServersControllers.SourceRole.ListPackages)
		pkg.GET("/:id", resource.ServersControllers.SourceRole.GetPackageWithModules)
		pkg.GET("/:id/limits", resource.ServersControllers.SourceRole.GetPackageLimits)
	}
}
