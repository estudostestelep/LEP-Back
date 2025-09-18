package handler

import (
	"encoding/json"
	"fmt"
	"lep/repositories"
	"lep/repositories/models"
	"lep/utils"
	"time"

	"github.com/google/uuid"
)

type NotificationHandler struct {
	notificationRepo repositories.INotificationRepository
	projectRepo      repositories.IProjectRepository
}

func NewNotificationHandler(notificationRepo repositories.INotificationRepository, projectRepo repositories.IProjectRepository) *NotificationHandler {
	return &NotificationHandler{
		notificationRepo: notificationRepo,
		projectRepo:      projectRepo,
	}
}

type SendNotificationRequest struct {
	OrganizationId uuid.UUID         `json:"organization_id"`
	ProjectId      uuid.UUID         `json:"project_id"`
	EventType      string            `json:"event_type"`
	EntityType     string            `json:"entity_type"`
	EntityId       uuid.UUID         `json:"entity_id"`
	Recipient      string            `json:"recipient"`
	Channel        string            `json:"channel"`
	Variables      map[string]string `json:"variables,omitempty"`
}

type ProcessNotificationEventRequest struct {
	OrganizationId uuid.UUID              `json:"organization_id"`
	ProjectId      uuid.UUID              `json:"project_id"`
	EventType      string                 `json:"event_type"`
	EntityType     string                 `json:"entity_type"`
	EntityId       uuid.UUID              `json:"entity_id"`
	Data           map[string]interface{} `json:"data"`
}

func (h *NotificationHandler) SendNotification(req SendNotificationRequest) error {
	// Buscar projeto
	project, err := h.projectRepo.GetProjectById(req.ProjectId)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	// Buscar configuração de notificação para o evento
	config, err := h.notificationRepo.GetNotificationConfigByEvent(req.OrganizationId, req.ProjectId, req.EventType)
	if err != nil {
		return fmt.Errorf("notification config not found for event %s: %w", req.EventType, err)
	}

	if !config.Enabled {
		return fmt.Errorf("notifications disabled for event: %s", req.EventType)
	}

	// Verificar se o canal está habilitado
	channelEnabled := false
	for _, channel := range config.Channels {
		if channel == req.Channel {
			channelEnabled = true
			break
		}
	}

	if !channelEnabled {
		return fmt.Errorf("channel %s not enabled for event %s", req.Channel, req.EventType)
	}

	// Buscar template para o canal
	template, err := h.notificationRepo.GetNotificationTemplateByChannel(req.OrganizationId, req.ProjectId, req.Channel)
	if err != nil {
		return fmt.Errorf("template not found for channel %s: %w", req.Channel, err)
	}

	// Criar service de notificação
	notificationService := utils.NewNotificationService()

	// Preparar request de notificação
	notificationReq := utils.NotificationRequest{
		Channel:   req.Channel,
		Recipient: req.Recipient,
		Subject:   template.Subject,
		Message:   template.Body,
		Variables: req.Variables,
	}

	// Enviar notificação
	result, err := notificationService.SendNotification(notificationReq, project)

	// Criar log independente do resultado
	logEntry := &models.NotificationLog{
		Id:             uuid.New(),
		OrganizationId: req.OrganizationId,
		ProjectId:      req.ProjectId,
		EventType:      req.EventType,
		Channel:        req.Channel,
		Recipient:      req.Recipient,
		Subject:        template.Subject,
		Message:        template.Body,
		Status:         result.Status,
		ExternalId:     result.ExternalId,
		ErrorMessage:   result.ErrorMessage,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Salvar log
	if logErr := h.notificationRepo.CreateNotificationLog(logEntry); logErr != nil {
		// Log do erro mas não interromper o fluxo
		fmt.Printf("Error creating notification log: %v\n", logErr)
	}

	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	return nil
}

func (h *NotificationHandler) ProcessNotificationEvent(req ProcessNotificationEventRequest) error {
	// Criar evento de notificação
	dataJSON, err := json.Marshal(req.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	event := &models.NotificationEvent{
		Id:             uuid.New(),
		OrganizationId: req.OrganizationId,
		ProjectId:      req.ProjectId,
		EventType:      req.EventType,
		EntityType:     req.EntityType,
		EntityId:       req.EntityId,
		Data:           string(dataJSON),
		Processed:      false,
		CreatedAt:      time.Now(),
	}

	return h.notificationRepo.CreateNotificationEvent(event)
}

func (h *NotificationHandler) UpdateNotificationStatus(externalId, status string, deliveredAt *time.Time) error {
	return h.notificationRepo.UpdateNotificationLogStatus(externalId, status, deliveredAt)
}

func (h *NotificationHandler) ProcessInboundMessage(orgId, projectId uuid.UUID, channel, from, to, body, externalId string) error {
	inbound := &models.NotificationInbound{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projectId,
		Channel:        channel,
		From:           from,
		To:             to,
		Body:           body,
		ExternalId:     externalId,
		Processed:      false,
		CreatedAt:      time.Now(),
	}

	return h.notificationRepo.CreateNotificationInbound(inbound)
}

func (h *NotificationHandler) GetNotificationLogs(orgId, projectId uuid.UUID, limit int) ([]models.NotificationLog, error) {
	return h.notificationRepo.GetNotificationLogsByProject(orgId, projectId, limit)
}

func (h *NotificationHandler) GetNotificationTemplates(orgId, projectId uuid.UUID) ([]models.NotificationTemplate, error) {
	return h.notificationRepo.GetNotificationTemplatesByProject(orgId, projectId)
}

func (h *NotificationHandler) CreateNotificationTemplate(template *models.NotificationTemplate) error {
	return h.notificationRepo.CreateNotificationTemplate(template)
}

func (h *NotificationHandler) UpdateNotificationTemplate(template *models.NotificationTemplate) error {
	return h.notificationRepo.UpdateNotificationTemplate(template)
}

func (h *NotificationHandler) CreateOrUpdateNotificationConfig(config *models.NotificationConfig) error {
	return h.notificationRepo.CreateOrUpdateNotificationConfig(config)
}
