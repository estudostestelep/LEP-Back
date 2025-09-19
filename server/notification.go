package server

import (
	"lep/handler"
	"lep/repositories/models"
	"lep/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationServer struct {
	notificationHandler *handler.NotificationHandler
}

func NewNotificationServer(notificationHandler *handler.NotificationHandler) *NotificationServer {
	return &NotificationServer{
		notificationHandler: notificationHandler,
	}
}

// === WEBHOOK ENDPOINTS ===

// TwilioWebhookStatus - Webhook para status de entrega do Twilio
func (s *NotificationServer) TwilioWebhookStatus(c *gin.Context) {
	var webhook utils.TwilioWebhookStatus
	if err := c.ShouldBindJSON(&webhook); err != nil {
		// Tentar bind de form data (Twilio envia como form)
		webhook.MessageSid = c.PostForm("MessageSid")
		webhook.MessageStatus = c.PostForm("MessageStatus")
		webhook.ErrorCode = c.PostForm("ErrorCode")
	}

	if webhook.MessageSid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "MessageSid is required"})
		return
	}

	// Determinar se foi entregue
	var deliveredAt *time.Time
	if webhook.MessageStatus == "delivered" || webhook.MessageStatus == "read" {
		now := time.Now()
		deliveredAt = &now
	}

	// Atualizar status no banco
	err := s.notificationHandler.UpdateNotificationStatus(webhook.MessageSid, webhook.MessageStatus, deliveredAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notification status"})
		return
	}

	// Resposta para Twilio
	c.XML(http.StatusOK, gin.H{"response": "ok"})
}

// TwilioWebhookInbound - Webhook para mensagens recebidas do Twilio
func (s *NotificationServer) TwilioWebhookInbound(c *gin.Context) {
	// Extrair parâmetros da URL para identificar organização/projeto
	orgIdStr := c.Param("orgId")
	projectIdStr := c.Param("projectId")

	orgId, err := uuid.Parse(orgIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	projectId, err := uuid.Parse(projectIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var webhook utils.TwilioWebhookInbound
	if err := c.ShouldBindJSON(&webhook); err != nil {
		// Tentar bind de form data (Twilio envia como form)
		webhook.From = c.PostForm("From")
		webhook.To = c.PostForm("To")
		webhook.Body = c.PostForm("Body")
		webhook.MessageSid = c.PostForm("MessageSid")
	}

	if webhook.From == "" || webhook.Body == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "From and Body are required"})
		return
	}

	// Determinar canal baseado no número "To"
	channel := "sms"
	if len(webhook.To) > 0 && webhook.To[:9] == "whatsapp:" {
		channel = "whatsapp"
	}

	// Processar mensagem recebida
	err = s.notificationHandler.ProcessInboundMessage(
		orgId, projectId, channel, webhook.From, webhook.To, webhook.Body, webhook.MessageSid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process inbound message"})
		return
	}

	// Resposta para Twilio
	c.XML(http.StatusOK, gin.H{"response": "ok"})
}

// === API ENDPOINTS ===

// SendNotification - Enviar notificação manual
func (s *NotificationServer) SendNotification(c *gin.Context) {
	var req handler.SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.notificationHandler.SendNotification(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification sent successfully"})
}

// ProcessEvent - Processar evento de notificação
func (s *NotificationServer) ProcessEvent(c *gin.Context) {
	var req handler.ProcessNotificationEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.notificationHandler.ProcessNotificationEvent(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event processed successfully"})
}

// GetNotificationLogs - Buscar logs de notificação
func (s *NotificationServer) GetNotificationLogs(c *gin.Context) {
	orgIdStr := c.Param("orgId")
	projectIdStr := c.Param("projectId")

	orgId, err := uuid.Parse(orgIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	projectId, err := uuid.Parse(projectIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	logs, err := s.notificationHandler.GetNotificationLogs(orgId, projectId, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

// GetNotificationTemplates - Buscar templates de notificação
func (s *NotificationServer) GetNotificationTemplates(c *gin.Context) {
	orgIdStr := c.Param("orgId")
	projectIdStr := c.Param("projectId")

	orgId, err := uuid.Parse(orgIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	projectId, err := uuid.Parse(projectIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	templates, err := s.notificationHandler.GetNotificationTemplates(orgId, projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

// CreateNotificationTemplate - Criar template de notificação
func (s *NotificationServer) CreateNotificationTemplate(c *gin.Context) {
	var template models.NotificationTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.notificationHandler.CreateNotificationTemplate(&template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"template": template})
}

// UpdateNotificationTemplate - Atualizar template de notificação
func (s *NotificationServer) UpdateNotificationTemplate(c *gin.Context) {
	var template models.NotificationTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.notificationHandler.UpdateNotificationTemplate(&template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"template": template})
}

// CreateOrUpdateNotificationConfig - Criar/atualizar configuração de notificação
func (s *NotificationServer) CreateOrUpdateNotificationConfig(c *gin.Context) {
	var config models.NotificationConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.notificationHandler.CreateOrUpdateNotificationConfig(&config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"config": config})
}