package utils

import (
	"encoding/json"
	"fmt"
	"lep/repositories"
	"lep/repositories/models"
	"log"
	"time"

	"github.com/google/uuid"
)

type EventService struct {
	notificationRepo repositories.INotificationRepository
	projectRepo      repositories.IProjectRepository
	settingsRepo     repositories.ISettingsRepository
}

type EventData struct {
	ReservationId *uuid.UUID `json:"reservation_id,omitempty"`
	CustomerId    *uuid.UUID `json:"customer_id,omitempty"`
	TableId       *uuid.UUID `json:"table_id,omitempty"`
	CustomerName  string     `json:"customer_name,omitempty"`
	CustomerPhone string     `json:"customer_phone,omitempty"`
	CustomerEmail string     `json:"customer_email,omitempty"`
	TableNumber   int        `json:"table_number,omitempty"`
	DateTime      *time.Time `json:"datetime,omitempty"`
	PartySize     int        `json:"party_size,omitempty"`
	Status        string     `json:"status,omitempty"`
	EstimatedWait int        `json:"estimated_wait,omitempty"` // em minutos
	Environment   string     `json:"environment,omitempty"`
}

func NewEventService(notificationRepo repositories.INotificationRepository, projectRepo repositories.IProjectRepository, settingsRepo repositories.ISettingsRepository) *EventService {
	return &EventService{
		notificationRepo: notificationRepo,
		projectRepo:      projectRepo,
		settingsRepo:     settingsRepo,
	}
}

// TriggerReservationCreated - Evento quando reserva é criada
func (e *EventService) TriggerReservationCreated(orgId, projectId uuid.UUID, reservation *models.Reservation, customer *models.Customer, table *models.Table) error {
	settings, err := e.settingsRepo.GetSettingsByProject(orgId, projectId)
	if err != nil || !settings.NotifyReservationCreate {
		return nil // Notificação desabilitada
	}

	eventData := EventData{
		ReservationId: &reservation.Id,
		CustomerId:    &reservation.CustomerId,
		TableId:       &reservation.TableId,
		CustomerName:  customer.Name,
		CustomerPhone: customer.Phone,
		CustomerEmail: customer.Email,
		TableNumber:   table.Number,
		DateTime:      parseTime(reservation.Datetime),
		PartySize:     reservation.PartySize,
		Status:        reservation.Status,
	}

	return e.createAndProcessEvent(orgId, projectId, "reservation_create", "reservation", reservation.Id, eventData)
}

// TriggerReservationUpdated - Evento quando reserva é atualizada
func (e *EventService) TriggerReservationUpdated(orgId, projectId uuid.UUID, reservation *models.Reservation, customer *models.Customer, table *models.Table) error {
	settings, err := e.settingsRepo.GetSettingsByProject(orgId, projectId)
	if err != nil || !settings.NotifyReservationUpdate {
		return nil
	}

	eventData := EventData{
		ReservationId: &reservation.Id,
		CustomerId:    &reservation.CustomerId,
		TableId:       &reservation.TableId,
		CustomerName:  customer.Name,
		CustomerPhone: customer.Phone,
		CustomerEmail: customer.Email,
		TableNumber:   table.Number,
		DateTime:      parseTime(reservation.Datetime),
		PartySize:     reservation.PartySize,
		Status:        reservation.Status,
	}

	return e.createAndProcessEvent(orgId, projectId, "reservation_update", "reservation", reservation.Id, eventData)
}

// TriggerReservationCancelled - Evento quando reserva é cancelada
func (e *EventService) TriggerReservationCancelled(orgId, projectId uuid.UUID, reservation *models.Reservation, customer *models.Customer, table *models.Table) error {
	settings, err := e.settingsRepo.GetSettingsByProject(orgId, projectId)
	if err != nil || !settings.NotifyReservationCancel {
		return nil
	}

	eventData := EventData{
		ReservationId: &reservation.Id,
		CustomerId:    &reservation.CustomerId,
		TableId:       &reservation.TableId,
		CustomerName:  customer.Name,
		CustomerPhone: customer.Phone,
		CustomerEmail: customer.Email,
		TableNumber:   table.Number,
		DateTime:      parseTime(reservation.Datetime),
		PartySize:     reservation.PartySize,
		Status:        reservation.Status,
	}

	return e.createAndProcessEvent(orgId, projectId, "reservation_cancel", "reservation", reservation.Id, eventData)
}

// TriggerTableAvailable - Evento quando mesa fica disponível (para fila de espera)
func (e *EventService) TriggerTableAvailable(orgId, projectId uuid.UUID, table *models.Table, customer *models.Customer, estimatedWait int) error {
	settings, err := e.settingsRepo.GetSettingsByProject(orgId, projectId)
	if err != nil || !settings.NotifyTableAvailable {
		return nil
	}

	eventData := EventData{
		CustomerId:    &customer.Id,
		TableId:       &table.Id,
		CustomerName:  customer.Name,
		CustomerPhone: customer.Phone,
		CustomerEmail: customer.Email,
		TableNumber:   table.Number,
		EstimatedWait: estimatedWait,
	}

	return e.createAndProcessEvent(orgId, projectId, "table_available", "table", table.Id, eventData)
}

// TriggerConfirmation24h - Evento de confirmação 24h antes
func (e *EventService) TriggerConfirmation24h(orgId, projectId uuid.UUID, reservation *models.Reservation, customer *models.Customer, table *models.Table) error {
	settings, err := e.settingsRepo.GetSettingsByProject(orgId, projectId)
	if err != nil || !settings.NotifyConfirmation24h {
		return nil
	}

	eventData := EventData{
		ReservationId: &reservation.Id,
		CustomerId:    &reservation.CustomerId,
		TableId:       &reservation.TableId,
		CustomerName:  customer.Name,
		CustomerPhone: customer.Phone,
		CustomerEmail: customer.Email,
		TableNumber:   table.Number,
		DateTime:      parseTime(reservation.Datetime),
		PartySize:     reservation.PartySize,
		Status:        reservation.Status,
	}

	return e.createAndProcessEvent(orgId, projectId, "confirmation_24h", "reservation", reservation.Id, eventData)
}

func parseTime(datetimeStr string) *time.Time {
	if datetimeStr == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, datetimeStr)
	if err != nil {
		return nil
	}
	return &t
}

// createAndProcessEvent - Cria evento e processa notificações automaticamente
func (e *EventService) createAndProcessEvent(orgId, projectId uuid.UUID, eventType, entityType string, entityId uuid.UUID, eventData EventData) error {
	// Criar evento
	dataJSON, err := json.Marshal(eventData)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	event := &models.NotificationEvent{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projectId,
		EventType:      eventType,
		EntityType:     entityType,
		EntityId:       entityId,
		Data:           string(dataJSON),
		Processed:      false,
		CreatedAt:      time.Now(),
	}

	if err := e.notificationRepo.CreateNotificationEvent(event); err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	// Processar evento imediatamente
	return e.ProcessEvent(event)
}

// ProcessEvent - Processa um evento e envia notificações
func (e *EventService) ProcessEvent(event *models.NotificationEvent) error {
	// Buscar configuração do evento
	config, err := e.notificationRepo.GetNotificationConfigByEvent(event.OrganizationId, event.ProjectId, event.EventType)
	if err != nil {
		log.Printf("No notification config found for event %s: %v", event.EventType, err)
		return e.markEventAsProcessed(event.Id)
	}

	if !config.Enabled {
		log.Printf("Notifications disabled for event %s", event.EventType)
		return e.markEventAsProcessed(event.Id)
	}

	// Parse event data
	var eventData EventData
	if err := json.Unmarshal([]byte(event.Data), &eventData); err != nil {
		return fmt.Errorf("failed to parse event data: %w", err)
	}

	// Buscar projeto
	project, err := e.projectRepo.GetProjectById(event.ProjectId)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	// Criar service de notificação
	notificationService := NewNotificationService()

	// Enviar notificação para cada canal habilitado
	for _, channel := range config.Channels {
		// Determinar destinatário baseado no canal
		var recipient string
		switch channel {
		case "sms", "whatsapp":
			recipient = eventData.CustomerPhone
		case "email":
			recipient = eventData.CustomerEmail
		}

		if recipient == "" {
			log.Printf("No recipient found for channel %s", channel)
			continue
		}

		// Buscar template para o canal
		template, err := e.notificationRepo.GetNotificationTemplateByChannel(event.OrganizationId, event.ProjectId, channel)
		if err != nil {
			log.Printf("No template found for channel %s: %v", channel, err)
			continue
		}

		// Preparar variáveis do template
		variables := e.buildTemplateVariables(eventData)

		// Preparar request de notificação
		notificationReq := NotificationRequest{
			Channel:   channel,
			Recipient: recipient,
			Subject:   template.Subject,
			Message:   template.Body,
			Variables: variables,
		}

		// Enviar notificação
		result, err := notificationService.SendNotification(notificationReq, project)

		// Criar log
		logEntry := &models.NotificationLog{
			Id:             uuid.New(),
			OrganizationId: event.OrganizationId,
			ProjectId:      event.ProjectId,
			EventType:      event.EventType,
			Channel:        channel,
			Recipient:      recipient,
			Subject:        template.Subject,
			Message:        template.Body,
			Status:         result.Status,
			ExternalId:     result.ExternalId,
			ErrorMessage:   result.ErrorMessage,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if logErr := e.notificationRepo.CreateNotificationLog(logEntry); logErr != nil {
			log.Printf("Error creating notification log: %v", logErr)
		}

		if err != nil {
			log.Printf("Failed to send notification via %s: %v", channel, err)
		}
	}

	// Marcar evento como processado
	return e.markEventAsProcessed(event.Id)
}

// buildTemplateVariables - Constrói variáveis para o template
func (e *EventService) buildTemplateVariables(eventData EventData) map[string]string {
	variables := make(map[string]string)

	if eventData.CustomerName != "" {
		variables["nome"] = eventData.CustomerName
		variables["cliente"] = eventData.CustomerName
	}

	if eventData.TableNumber > 0 {
		variables["mesa"] = fmt.Sprintf("%d", eventData.TableNumber)
		variables["numero_mesa"] = fmt.Sprintf("%d", eventData.TableNumber)
	}

	if eventData.DateTime != nil {
		variables["data"] = eventData.DateTime.Format("02/01/2006")
		variables["hora"] = eventData.DateTime.Format("15:04")
		variables["data_hora"] = eventData.DateTime.Format("02/01/2006 às 15:04")
	}

	if eventData.PartySize > 0 {
		variables["pessoas"] = fmt.Sprintf("%d", eventData.PartySize)
		variables["quantidade_pessoas"] = fmt.Sprintf("%d", eventData.PartySize)
	}

	if eventData.EstimatedWait > 0 {
		variables["tempo_espera"] = fmt.Sprintf("%d minutos", eventData.EstimatedWait)
	}

	if eventData.Status != "" {
		variables["status"] = eventData.Status
	}

	return variables
}

// markEventAsProcessed - Marca evento como processado
func (e *EventService) markEventAsProcessed(eventId uuid.UUID) error {
	return e.notificationRepo.MarkEventAsProcessed(eventId)
}

// ProcessPendingEvents - Processa eventos pendentes (para jobs/cron)
func (e *EventService) ProcessPendingEvents(orgId, projectId uuid.UUID) error {
	events, err := e.notificationRepo.GetUnprocessedEvents(orgId, projectId)
	if err != nil {
		return err
	}

	for _, event := range events {
		if err := e.ProcessEvent(&event); err != nil {
			log.Printf("Error processing event %s: %v", event.Id, err)
		}
	}

	return nil
}
