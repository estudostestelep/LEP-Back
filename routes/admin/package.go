package admin

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupPackageRoutes configura rotas de pacote/plano para admin
// CRUD de pacotes e gerenciamento de assinaturas
func SetupPackageRoutes(r gin.IRouter) {
	pkg := r.Group("/package")
	{
		// CRUD de pacotes
		pkg.POST("", resource.ServersControllers.SourceRole.CreatePackage)
		pkg.PUT("/:id", resource.ServersControllers.SourceRole.UpdatePackage)
		pkg.DELETE("/:id", resource.ServersControllers.SourceRole.DeletePackage)

		// Gerenciar módulos em pacotes
		pkg.POST("/:id/modules/:moduleId", resource.ServersControllers.SourceRole.AddModuleToPackage)
		pkg.DELETE("/:id/modules/:moduleId", resource.ServersControllers.SourceRole.RemoveModuleFromPackage)

		// Gerenciar limites de pacotes
		pkg.POST("/:id/limits", resource.ServersControllers.SourceRole.SetPackageLimit)

		// Lista de todas as assinaturas
		pkg.GET("/subscriptions", resource.ServersControllers.SourceRole.ListAllSubscriptions)

		// Gerenciar assinatura de organização específica
		pkg.POST("/subscription/:orgId", resource.ServersControllers.SourceRole.CreateOrganizationSubscription)
		pkg.PUT("/subscription/:orgId", resource.ServersControllers.SourceRole.UpdateOrganizationSubscription)
		pkg.DELETE("/subscription/:orgId", resource.ServersControllers.SourceRole.CancelOrganizationSubscription)
	}
}
