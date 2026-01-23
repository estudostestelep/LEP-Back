package routes

import (
	"lep/middleware"
	"lep/resource"
	"lep/routes/admin"
	"lep/routes/client"
	"lep/routes/public"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas as rotas da aplicação
func SetupRoutes(r *gin.Engine) {
	// 1. Rotas públicas (sem autenticação)
	public.SetupPublicRoutes(r)

	// 2. Grupo protegido (auth + headers obrigatórios)
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.HeaderValidationMiddleware())

	// Rotas de auth protegidas
	protected.POST("/logout", resource.ServersControllers.SourceAuth.ServiceLogout)
	protected.POST("/checkToken", resource.ServersControllers.SourceAuth.ServiceValidateToken)

	// 3. Rotas admin (/admin/*) - MasterAdmin only
	// Middleware MasterAdminOnlyMiddleware aplicado no grupo admin
	admin.SetupAdminRoutes(protected)

	// 4. Rotas client - Role-based permissions
	// Middleware RolePermissionMiddleware aplicado em cada rota individualmente
	client.SetupClientRoutes(protected)
}
