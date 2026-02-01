package public

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes configura rotas de autenticação públicas
func SetupAuthRoutes(r *gin.Engine) {
	// Login legado (mantido para compatibilidade)
	r.POST("/login", resource.ServersControllers.SourceAuth.ServiceLogin)

	// ========== NOVOS ENDPOINTS DE LOGIN SEPARADOS ==========
	// Login de administradores
	r.POST("/admin/login", resource.ServersControllers.SourceAuthAdmin.ServiceAdminLogin)

	// Login de clientes (requer org_slug no body)
	r.POST("/client/login", resource.ServersControllers.SourceAuthClient.ServiceClientLogin)

	// Criar organização (bootstrap)
	r.POST("/create-organization", resource.ServersControllers.SourceOrganization.ServiceCreateOrganizationBootstrap)

	// Admin routes (temporário - para reset de senha)
	r.POST("/admin/reset-passwords", resource.ServersControllers.SourceAdmin.ServiceResetPasswords)

	// Dev migration endpoint
	r.POST("/admin/run-migration", resource.ServersControllers.SourceAdmin.ServiceRunDevMigration)
}
