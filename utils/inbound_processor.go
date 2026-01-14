package utils

import (
	"fmt"
	"lep/repositories"
	"lep/repositories/models"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

// InboundProcessorService processa mensagens recebidas de clientes
type InboundProcessorService struct {
	notificationRepo repositories.INotificationRepository
	reservationRepo  repositories.IReservationRepository
	customerRepo     repositories.ICustomersRepository
	tableRepo        repositories.ITableRepository
	settingsRepo     repositories.ISettingsRepository
	classifier       *ResponseClassifierService
}

// ProcessingResult resultado do processamento de uma mensagem
type ProcessingResult struct {
	Success       bool
	Action        string // "reservation_confirmed", "reservation_cancelled", "queued_for_review", "no_action", "customer_not_found", "reservation_not_found"
	ReservationId *uuid.UUID
	CustomerId    *uuid.UUID
	Message       string
}

// NewInboundProcessorService cria nova instância do processador
func NewInboundProcessorService(
	notificationRepo repositories.INotificationRepository,
	reservationRepo repositories.IReservationRepository,
	customerRepo repositories.ICustomersRepository,
	tableRepo repositories.ITableRepository,
	settingsRepo repositories.ISettingsRepository,
) *InboundProcessorService {
	return &InboundProcessorService{
		notificationRepo: notificationRepo,
		reservationRepo:  reservationRepo,
		customerRepo:     customerRepo,
		tableRepo:        tableRepo,
		settingsRepo:     settingsRepo,
		classifier:       NewResponseClassifierService(),
	}
}

// ProcessInboundMessage processa uma mensagem inbound e toma ação apropriada
func (s *InboundProcessorService) ProcessInboundMessage(inbound *models.NotificationInbound) ProcessingResult {
	// Passo 1: Normaliza o telefone e busca o cliente
	normalizedPhone := s.normalizePhone(inbound.From)
	customer, err := s.customerRepo.GetCustomerByPhone(inbound.OrganizationId, inbound.ProjectId, normalizedPhone)
	if err != nil || customer == nil {
		log.Printf("Cliente não encontrado para telefone %s: %v", inbound.From, err)
		return ProcessingResult{
			Success: false,
			Action:  "customer_not_found",
			Message: "Cliente não encontrado para o número de telefone",
		}
	}

	inbound.CustomerId = &customer.Id

	// Passo 2: Busca reserva pendente de confirmação
	reservation, err := s.reservationRepo.GetPendingConfirmationReservation(inbound.OrganizationId, inbound.ProjectId, customer.Id)
	if err != nil || reservation == nil {
		log.Printf("Nenhuma reserva pendente encontrada para cliente %s: %v", customer.Id, err)
		return ProcessingResult{
			Success:    false,
			Action:     "reservation_not_found",
			CustomerId: &customer.Id,
			Message:    "Nenhuma reserva pendente encontrada para confirmação",
		}
	}

	inbound.ReservationId = &reservation.Id

	// Passo 3: Classifica a resposta
	classification := s.classifier.ClassifyResponse(inbound.Body)
	inbound.ResponseType = classification.ResponseType
	inbound.ConfidenceScore = classification.ConfidenceScore
	inbound.ProcessingMethod = classification.Method

	// Passo 4: Verifica configuração de processamento
	settings, _ := s.settingsRepo.GetSettingsByProject(inbound.OrganizationId, inbound.ProjectId)
	processingMode := "automatic"
	if settings != nil && settings.ResponseProcessingMode != "" {
		processingMode = settings.ResponseProcessingMode
	}

	// Passo 5: Processa baseado no modo configurado
	result := ProcessingResult{
		Success:       true,
		ReservationId: &reservation.Id,
		CustomerId:    &customer.Id,
	}

	switch processingMode {
	case "automatic":
		result = s.processAutomatic(inbound, reservation, classification)

	case "ai_assisted":
		result = s.processAiAssisted(inbound, reservation, customer, classification)

	case "manual":
		result = s.processManual(inbound, reservation, customer)

	default:
		result = s.processAutomatic(inbound, reservation, classification)
	}

	result.ReservationId = &reservation.Id
	result.CustomerId = &customer.Id

	return result
}

// processAutomatic processa automaticamente com base na classificação
func (s *InboundProcessorService) processAutomatic(inbound *models.NotificationInbound, reservation *models.Reservation, classification ClassificationResult) ProcessingResult {
	result := ProcessingResult{Success: true}

	switch classification.ResponseType {
	case "confirmed":
		if err := s.confirmReservation(reservation); err != nil {
			result.Success = false
			result.Action = "error"
			result.Message = fmt.Sprintf("Erro ao confirmar reserva: %v", err)
		} else {
			result.Action = "reservation_confirmed"
			result.Message = "Reserva confirmada com sucesso"
			inbound.ActionTaken = "reservation_confirmed"
		}

	case "cancelled":
		if err := s.cancelReservation(reservation); err != nil {
			result.Success = false
			result.Action = "error"
			result.Message = fmt.Sprintf("Erro ao cancelar reserva: %v", err)
		} else {
			result.Action = "reservation_cancelled"
			result.Message = "Reserva cancelada conforme solicitação"
			inbound.ActionTaken = "reservation_cancelled"
		}

	default:
		result.Action = "no_action"
		result.Message = "Não foi possível entender a resposta. Por favor, responda SIM para confirmar ou NÃO para cancelar."
		inbound.ActionTaken = "no_action"
	}

	return result
}

// processAiAssisted cria item na fila com sugestão da IA
func (s *InboundProcessorService) processAiAssisted(inbound *models.NotificationInbound, reservation *models.Reservation, customer *models.Customer, classification ClassificationResult) ProcessingResult {
	result := ProcessingResult{Success: true}

	// Mapeia tipo de resposta para ação sugerida
	suggestedAction := "none"
	if classification.ResponseType == "confirmed" {
		suggestedAction = "confirm"
	} else if classification.ResponseType == "cancelled" {
		suggestedAction = "cancel"
	}

	// Cria item na fila de revisão
	queueItem := &models.ResponseReviewQueue{
		OrganizationId:  inbound.OrganizationId,
		ProjectId:       inbound.ProjectId,
		InboundId:       inbound.Id,
		ReservationId:   reservation.Id,
		CustomerId:      customer.Id,
		CustomerName:    customer.Name,
		CustomerPhone:   customer.Phone,
		MessageBody:     inbound.Body,
		SuggestedAction: suggestedAction,
		ConfidenceScore: classification.ConfidenceScore,
		Status:          "pending_review",
	}

	if err := s.notificationRepo.CreateReviewQueueItem(queueItem); err != nil {
		result.Success = false
		result.Action = "error"
		result.Message = fmt.Sprintf("Erro ao criar item na fila de revisão: %v", err)
	} else {
		result.Action = "queued_for_review"
		result.Message = "Mensagem enviada para fila de revisão com sugestão"
		inbound.ActionTaken = "queued_for_review"
	}

	return result
}

// processManual cria item na fila sem sugestão
func (s *InboundProcessorService) processManual(inbound *models.NotificationInbound, reservation *models.Reservation, customer *models.Customer) ProcessingResult {
	result := ProcessingResult{Success: true}

	// Cria item na fila de revisão sem sugestão
	queueItem := &models.ResponseReviewQueue{
		OrganizationId:  inbound.OrganizationId,
		ProjectId:       inbound.ProjectId,
		InboundId:       inbound.Id,
		ReservationId:   reservation.Id,
		CustomerId:      customer.Id,
		CustomerName:    customer.Name,
		CustomerPhone:   customer.Phone,
		MessageBody:     inbound.Body,
		SuggestedAction: "",
		ConfidenceScore: 0,
		Status:          "pending_review",
	}

	if err := s.notificationRepo.CreateReviewQueueItem(queueItem); err != nil {
		result.Success = false
		result.Action = "error"
		result.Message = fmt.Sprintf("Erro ao criar item na fila de revisão: %v", err)
	} else {
		result.Action = "queued_for_review"
		result.Message = "Mensagem enviada para fila de revisão manual"
		inbound.ActionTaken = "queued_for_review"
	}

	return result
}

// normalizePhone remove prefixos e formatação do telefone
func (s *InboundProcessorService) normalizePhone(phone string) string {
	// Remove prefixo whatsapp:
	phone = strings.TrimPrefix(phone, "whatsapp:")
	// Remove espaços, traços, parênteses
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")
	// Remove + inicial
	phone = strings.TrimPrefix(phone, "+")
	return phone
}

// confirmReservation atualiza status da reserva para confirmado
func (s *InboundProcessorService) confirmReservation(reservation *models.Reservation) error {
	reservation.Status = "confirmed"
	reservation.UpdatedAt = time.Now()
	return s.reservationRepo.UpdateReservation(reservation)
}

// cancelReservation cancela a reserva e libera a mesa
func (s *InboundProcessorService) cancelReservation(reservation *models.Reservation) error {
	// Atualiza status da reserva
	reservation.Status = "cancelled"
	reservation.UpdatedAt = time.Now()
	if err := s.reservationRepo.UpdateReservation(reservation); err != nil {
		return err
	}

	// Libera a mesa
	table, err := s.tableRepo.GetTableById(reservation.TableId)
	if err != nil {
		return err
	}

	table.Status = "livre"
	table.UpdatedAt = time.Now()
	return s.tableRepo.UpdateTable(table)
}

// ApproveReviewItem aprova item da fila e executa ação
func (s *InboundProcessorService) ApproveReviewItem(itemId uuid.UUID, reviewedBy uuid.UUID) error {
	item, err := s.notificationRepo.GetReviewQueueItemById(itemId)
	if err != nil {
		return err
	}

	reservation, err := s.reservationRepo.GetReservationById(item.ReservationId)
	if err != nil {
		return err
	}

	// Executa a ação sugerida
	switch item.SuggestedAction {
	case "confirm":
		if err := s.confirmReservation(reservation); err != nil {
			return err
		}
		item.ActionTaken = "reservation_confirmed"
	case "cancel":
		if err := s.cancelReservation(reservation); err != nil {
			return err
		}
		item.ActionTaken = "reservation_cancelled"
	default:
		return fmt.Errorf("nenhuma ação sugerida para aprovar")
	}

	// Atualiza item da fila
	item.Status = "approved"
	item.ReviewedBy = &reviewedBy
	now := time.Now()
	item.ReviewedAt = &now

	return s.notificationRepo.UpdateReviewQueueItem(item)
}

// RejectReviewItem rejeita item da fila sem ação
func (s *InboundProcessorService) RejectReviewItem(itemId uuid.UUID, reviewedBy uuid.UUID, notes string) error {
	item, err := s.notificationRepo.GetReviewQueueItemById(itemId)
	if err != nil {
		return err
	}

	item.Status = "rejected"
	item.ReviewedBy = &reviewedBy
	item.Notes = notes
	item.ActionTaken = "no_action"
	now := time.Now()
	item.ReviewedAt = &now

	return s.notificationRepo.UpdateReviewQueueItem(item)
}

// ExecuteCustomAction executa ação customizada no item da fila
func (s *InboundProcessorService) ExecuteCustomAction(itemId uuid.UUID, reviewedBy uuid.UUID, action string, notes string) error {
	item, err := s.notificationRepo.GetReviewQueueItemById(itemId)
	if err != nil {
		return err
	}

	reservation, err := s.reservationRepo.GetReservationById(item.ReservationId)
	if err != nil {
		return err
	}

	// Executa a ação escolhida
	switch action {
	case "confirm":
		if err := s.confirmReservation(reservation); err != nil {
			return err
		}
		item.ActionTaken = "reservation_confirmed"
	case "cancel":
		if err := s.cancelReservation(reservation); err != nil {
			return err
		}
		item.ActionTaken = "reservation_cancelled"
	default:
		item.ActionTaken = "no_action"
	}

	// Atualiza item da fila
	item.Status = "approved"
	item.ReviewedBy = &reviewedBy
	item.Notes = notes
	now := time.Now()
	item.ReviewedAt = &now

	return s.notificationRepo.UpdateReviewQueueItem(item)
}
