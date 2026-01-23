package public

import (
	"github.com/gin-gonic/gin"
)

// SetupPublicRoutes configura todas as rotas públicas
// Nenhuma autenticação necessária
func SetupPublicRoutes(r *gin.Engine) {
	SetupAuthRoutes(r)
	SetupMenuRoutes(r)
	SetupWebhookRoutes(r)
	SetupUploadRoutes(r)
	SetupSeedingRoutes(r)
}
