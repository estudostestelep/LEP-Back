package admin

import (
	"lep/handler"
	"lep/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAdminRoutes configura todas as rotas administrativas
// Todas as rotas neste grupo requerem escopo admin (master_admin ou cargos admin como admin_support)
func SetupAdminRoutes(r gin.IRouter, authHandler handler.IHandlerAuth) {
	admin := r.Group("/admin")
	admin.Use(middleware.AdminScopeMiddleware(authHandler))

	SetupMenuRoutes(admin)
	SetupCategoryRoutes(admin)
	SetupSubcategoryRoutes(admin)
	SetupRoleRoutes(admin)
	SetupModuleRoutes(admin)
	SetupPermissionRoutes(admin)
	SetupPackageRoutes(admin) // Rotas de Plans (mantendo nome legacy "package")
	SetupPlanChangeRoutes(admin)
	SetupAuditLogRoutes(admin)
	SetupImagesRoutes(admin)
	SetupSidebarConfigRoutes(admin)

	// Sistema de usuários (admins e clients separados)
	SetupAdminUserRoutes(admin)
	SetupClientUserRoutes(admin)
}
