package client

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupNotificationRoutes configura rotas de notificação para client
func SetupNotificationRoutes(r gin.IRouter) {
	notification := r.Group("/notification")
	{
		// Enviar notificação manual
		notification.POST("/send", resource.ServersControllers.SourceNotification.SendNotification)
		// Processar evento de notificação
		notification.POST("/event", resource.ServersControllers.SourceNotification.ProcessEvent)
		// Logs de notificação
		notification.GET("/logs/:orgId/:projectId", resource.ServersControllers.SourceNotification.GetNotificationLogs)
		// Templates
		notification.GET("/templates/:orgId/:projectId", resource.ServersControllers.SourceNotification.GetNotificationTemplates)
		notification.POST("/template", resource.ServersControllers.SourceNotification.CreateNotificationTemplate)
		notification.PUT("/template", resource.ServersControllers.SourceNotification.UpdateNotificationTemplate)
		// Configurações
		notification.POST("/config", resource.ServersControllers.SourceNotification.CreateOrUpdateNotificationConfig)

		// Fila de revisão de respostas
		notification.GET("/review-queue/:orgId/:projectId", resource.ServersControllers.SourceNotification.GetReviewQueue)
		notification.POST("/review-queue/:id/approve", resource.ServersControllers.SourceNotification.ApproveReviewItem)
		notification.POST("/review-queue/:id/reject", resource.ServersControllers.SourceNotification.RejectReviewItem)
		notification.POST("/review-queue/:id/custom", resource.ServersControllers.SourceNotification.ExecuteCustomAction)

		// Lembretes customizados
		notification.GET("/reminders/:orgId/:projectId", resource.ServersControllers.SourceNotification.GetNotificationReminders)
		notification.POST("/reminder", resource.ServersControllers.SourceNotification.CreateNotificationReminder)
		notification.PUT("/reminder", resource.ServersControllers.SourceNotification.UpdateNotificationReminder)
		notification.DELETE("/reminder/:id", resource.ServersControllers.SourceNotification.DeleteNotificationReminder)

		// Debug/Admin - Executar job de notificações manualmente
		notification.POST("/trigger-scheduled", resource.ServersControllers.SourceNotification.TriggerScheduledNotifications)
	}
}
