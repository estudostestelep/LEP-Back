package client

import (
	"github.com/gin-gonic/gin"
)

// SetupClientRoutes configura todas as rotas de cliente
// Requer AuthMiddleware + HeaderValidationMiddleware (aplicados no routes.go principal)
func SetupClientRoutes(r gin.IRouter) {
	// Rotas de leitura (sem middleware específico de permissão)
	SetupMenuReadRoutes(r)
	SetupCategoryReadRoutes(r)
	SetupSubcategoryReadRoutes(r)
	SetupRoleClientRoutes(r)
	SetupModuleReadRoutes(r)
	SetupPermissionReadRoutes(r)
	SetupPackageClientRoutes(r)
	SetupOrganizationRoutes(r)
	SetupProjectRoutes(r)
	SetupUserRoutes(r)
	SetupSettingsRoutes(r)
	SetupReportsRoutes(r)
	SetupNotificationRoutes(r)
	SetupClientAuditLogRoutes(r)
	SetupPlanChangeClientRoutes(r)
	SetupOnboardingRoutes(r)
	SetupEnvironmentRoutes(r)

	// Rotas com RolePermissionMiddleware
	SetupProductRoutes(r)
	SetupTableRoutes(r)
	SetupReservationRoutes(r)
	SetupCustomerRoutes(r)
	SetupOrderRoutes(r)
	SetupWaitlistRoutes(r)
	SetupTagRoutes(r)
}
