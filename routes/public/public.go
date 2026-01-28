package public

import (
	"github.com/gin-gonic/gin"
)

// SetupPublicRoutes configura todas as rotas públicas
// Nenhuma autenticação necessária
func SetupPublicRoutes(r *gin.Engine) {
	SetupAuthRoutes(r)
	SetupTenantRoutes(r)
	SetupMenuRoutes(r)
	SetupWebhookRoutes(r)
	SetupUploadRoutes(r)
	SetupSeedingRoutes(r)
}
