package admin

import (
	"lep/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAdminRoutes configura todas as rotas administrativas
// Todas as rotas neste grupo requerem MasterAdmin
func SetupAdminRoutes(r gin.IRouter) {
	admin := r.Group("/admin")
	admin.Use(middleware.MasterAdminOnlyMiddleware())

	SetupMenuRoutes(admin)
	SetupCategoryRoutes(admin)
	SetupSubcategoryRoutes(admin)
	SetupRoleRoutes(admin)
	SetupModuleRoutes(admin)
	SetupPermissionRoutes(admin)
	SetupPackageRoutes(admin)
	SetupUserRoutes(admin)
	SetupPlanChangeRoutes(admin)
	SetupAuditLogRoutes(admin)
	SetupImagesRoutes(admin)
	SetupSidebarConfigRoutes(admin)
}
