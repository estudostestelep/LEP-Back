package public

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes configura rotas de autenticação públicas
func SetupAuthRoutes(r *gin.Engine) {
	// Login
	r.POST("/login", resource.ServersControllers.SourceAuth.ServiceLogin)

	// Cadastro público
	r.POST("/user", resource.ServersControllers.SourceUsers.ServiceCreateUser)

	// Criar organização (bootstrap)
	r.POST("/create-organization", resource.ServersControllers.SourceOrganization.ServiceCreateOrganizationBootstrap)

	// Admin routes (temporário - para reset de senha)
	r.POST("/admin/reset-passwords", resource.ServersControllers.SourceAdmin.ServiceResetPasswords)

	// Dev migration endpoint
	r.POST("/admin/run-migration", resource.ServersControllers.SourceAdmin.ServiceRunDevMigration)
}
