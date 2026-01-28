package admin

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupAdminUserRoutes configura rotas de CRUD para admins (tabela admins)
// Apenas master_admin pode gerenciar outros admins
func SetupAdminUserRoutes(r gin.IRouter) {
	adminUser := r.Group("/admin-user")
	{
		adminUser.GET("", resource.ServersControllers.SourceAdminUsers.ServiceListAdmins)
		adminUser.GET("/:id", resource.ServersControllers.SourceAdminUsers.ServiceGetAdmin)
		adminUser.POST("", resource.ServersControllers.SourceAdminUsers.ServiceCreateAdmin)
		adminUser.PUT("/:id", resource.ServersControllers.SourceAdminUsers.ServiceUpdateAdmin)
		adminUser.DELETE("/:id", resource.ServersControllers.SourceAdminUsers.ServiceDeleteAdmin)
	}
}
