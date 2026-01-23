package public

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupWebhookRoutes configura rotas de webhook públicas
func SetupWebhookRoutes(r *gin.Engine) {
	webhook := r.Group("/webhook")
	{
		// Twilio status callback
		webhook.POST("/twilio/status", resource.ServersControllers.SourceNotification.TwilioWebhookStatus)
		// Twilio inbound messages (com org/project na URL)
		webhook.POST("/twilio/inbound/:orgId/:projectId", resource.ServersControllers.SourceNotification.TwilioWebhookInbound)
	}
}
