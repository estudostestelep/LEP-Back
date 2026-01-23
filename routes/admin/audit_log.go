package admin

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupAuditLogRoutes configura rotas de logs de auditoria para admin
func SetupAuditLogRoutes(r gin.IRouter) {
	// Admin Audit Logs
	auditLogs := r.Group("/audit-logs")
	{
		auditLogs.GET("", resource.ServersControllers.SourceAdminAuditLog.ServiceListAdminAuditLogs)
		auditLogs.GET("/:id", resource.ServersControllers.SourceAdminAuditLog.ServiceGetAdminAuditLog)
		auditLogs.DELETE("/cleanup", resource.ServersControllers.SourceAdminAuditLog.ServiceDeleteOldLogs)
	}

	// Client Audit Config
	clientAuditConfig := r.Group("/client-audit-config")
	{
		clientAuditConfig.GET("", resource.ServersControllers.SourceClientAuditLog.ServiceGetClientAuditConfig)
		clientAuditConfig.PUT("", resource.ServersControllers.SourceClientAuditLog.ServiceUpdateClientAuditConfig)
		clientAuditConfig.GET("/modules", resource.ServersControllers.SourceClientAuditLog.ServiceGetAvailableModules)
	}
}
