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

	// Extrair organization_id do header se não vier no body
	if template.OrganizationId == uuid.Nil {
		orgIdStr := c.GetHeader("X-Lpe-Organization-Id")
		if orgId, err := uuid.Parse(orgIdStr); err == nil {
			template.OrganizationId = orgId
		}
	}

	// Gerar ID se não existir
	if template.Id == uuid.Nil {
		template.Id = uuid.New()
	}

	// Definir timestamps
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()
	template.Active = true

	err := s.notificationHandler.CreateNotificationTemplate(&template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retornar resposta padronizada com o template criado (inclui ID)
	utils.SendCreatedSuccess(c, "Template created successfully", template)
}

// UpdateNotificationTemplate - Atualizar template de notificação
func (s *NotificationServer) UpdateNotificationTemplate(c *gin.Context) {
	var template models.NotificationTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extrair organization_id do header se não vier no body
	if template.OrganizationId == uuid.Nil {
		orgIdStr := c.GetHeader("X-Lpe-Organization-Id")
		if orgId, err := uuid.Parse(orgIdStr); err == nil {
			template.OrganizationId = orgId
		}
	}

	template.UpdatedAt = time.Now()

	err := s.notificationHandler.UpdateNotificationTemplate(&template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

// CreateOrUpdateNotificationConfig - Criar/atualizar configuração de notificação
func (s *NotificationServer) CreateOrUpdateNotificationConfig(c *gin.Context) {
	var config models.NotificationConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extrair organization_id do header se não vier no body
	if config.OrganizationId == uuid.Nil {
		orgIdStr := c.GetHeader("X-Lpe-Organization-Id")
		if orgId, err := uuid.Parse(orgIdStr); err == nil {
			config.OrganizationId = orgId
		}
	}

	// Extrair project_id do header se não vier no body
	if config.ProjectId == uuid.Nil {
		projIdStr := c.GetHeader("X-Lpe-Project-Id")
		if projId, err := uuid.Parse(projIdStr); err == nil {
			config.ProjectId = projId
		}
	}

	err := s.notificationHandler.CreateOrUpdateNotificationConfig(&config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// === REVIEW QUEUE ENDPOINTS ===

// GetReviewQueue - Lista itens pendentes na fila de revisão
func (s *NotificationServer) GetReviewQueue(c *gin.Context) {
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

	items, err := s.notificationHandler.GetPendingReviewItems(orgId, projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

// ApproveReviewItem - Aprova um item da fila e executa a ação sugerida
func (s *NotificationServer) ApproveReviewItem(c *gin.Context) {
	itemIdStr := c.Param("id")
	itemId, err := uuid.Parse(itemIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	// Pegar user ID do contexto de autenticação
	userIdStr := c.GetString("user_id")
	userId, _ := uuid.Parse(userIdStr)

	err = s.notificationHandler.ApproveReviewItem(itemId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item approved and action executed"})
}

// RejectReviewItem - Rejeita um item da fila sem ação
func (s *NotificationServer) RejectReviewItem(c *gin.Context) {
	itemIdStr := c.Param("id")
	itemId, err := uuid.Parse(itemIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}
	c.ShouldBindJSON(&req)

	// Pegar user ID do contexto de autenticação
	userIdStr := c.GetString("user_id")
	userId, _ := uuid.Parse(userIdStr)

	err = s.notificationHandler.RejectReviewItem(itemId, userId, req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item rejected"})
}

// ExecuteCustomAction - Executa uma ação customizada no item da fila
func (s *NotificationServer) ExecuteCustomAction(c *gin.Context) {
	itemIdStr := c.Param("id")
	itemId, err := uuid.Parse(itemIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req struct {
		Action string `json:"action"` // "confirm" ou "cancel"
		Notes  string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pegar user ID do contexto de autenticação
	userIdStr := c.GetString("user_id")
	userId, _ := uuid.Parse(userIdStr)

	err = s.notificationHandler.ExecuteCustomAction(itemId, userId, req.Action, req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Custom action executed"})
}