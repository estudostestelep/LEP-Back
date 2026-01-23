package client

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupClientAuditLogRoutes configura rotas de logs de auditoria de cliente
func SetupClientAuditLogRoutes(r gin.IRouter) {
	auditLogs := r.Group("/client-audit-logs")
	{
		auditLogs.GET("", resource.ServersControllers.SourceClientAuditLog.ServiceListClientAuditLogs)
		auditLogs.GET("/:id", resource.ServersControllers.SourceClientAuditLog.ServiceGetClientAuditLog)
	}
}
