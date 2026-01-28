package public

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupTenantRoutes configura rotas de tenant públicas
func SetupTenantRoutes(r *gin.Engine) {
	// Resolver tenant por slug
	r.GET("/tenant/resolve", resource.ServersControllers.SourceTenant.ServiceResolveTenant)
}
