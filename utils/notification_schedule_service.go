package utils

import (
	"encoding/json"
	"lep/repositories"
	"lep/repositories/models"
	"log"
	"time"

	"github.com/google/uuid"
)

// NotificationScheduleService gerencia agendamentos de notificações
type NotificationScheduleService struct {
	notificationRepo repositories.INotificationRepository
	reservationRepo  repositories.IReservationRepository
	customerRepo     repositories.ICustomersRepository
	tableRepo        repositories.ITableRepository
	settingsRepo     repositories.ISettingsRepository
	projectRepo      repositories.IProjectRepository
	eventService     *EventService
}

// ScheduleMetadata metadados salvos no agendamento
type ScheduleMetadata struct {
	CustomerName  string    `json:"customer_name"`
	CustomerPhone string    `json:"customer_phone"`
	CustomerEmail string    `json:"customer_email"`
	TableNumber   int       `json:"table_number"`
	DateTime      time.Time `json:"datetime"`
	PartySize     int       `json:"party_size"`
}

// NewNotificationScheduleService cria nova instância do serviço
func NewNotificationScheduleService(
	notificationRepo repositories.INotificationRepository,
	reservationRepo repositories.IReservationRepository,
	customerRepo repositories.ICustomersRepository,
	tableRepo repositories.ITableRepository,
	settingsRepo repositories.ISettingsRepository,
	projectRepo repositories.IProjectRepository,
) *NotificationScheduleService {
	eventService := NewEventService(notificationRepo, projectRepo, settingsRepo)
	return &NotificationScheduleService{
		notificationRepo: notificationRepo,
		reservationRepo:  reservationRepo,
		customerRepo:     customerRepo,
		tableRepo:        tableRepo,
		settingsRepo:     settingsRepo,
		projectRepo:      projectRepo,
		eventService:     eventService,
	}
}

// ScheduleReservationNotifications cria agendamentos para uma nova reserva
func (s *NotificationScheduleService) ScheduleReservationNotifications(reservation *models.Reservation, customer *models.Customer, table *models.Table) error {
	settings, err := s.settingsRepo.GetSettingsByProject(reservation.OrganizationId, reservation.ProjectId)
	if err != nil {
		log.Printf("Erro ao buscar settings para agendamento: %v", err)
		return err
	}

	// Parse datetime da reserva
	reservationTime, err := time.Parse(time.RFC3339, reservation.Datetime)
	if err != nil {
		log.Printf("Erro ao parsear datetime da reserva: %v", err)
		return err
	}

	// Cria metadados
	metadata := ScheduleMetadata{
		CustomerName:  customer.Name,
		CustomerPhone: customer.Phone,
		CustomerEmail: customer.Email,
		TableNumber:   table.Number,
		DateTime:      reservationTime,
		PartySize:     reservation.PartySize,
	}
	metadataJSON, _ := json.Marshal(metadata)

	// Agenda confirmação (se habilitado e horas > 0)
	if settings.NotifyConfirmation24h && settings.ConfirmationHoursBefore > 0 {
		confirmTime := reservationTime.Add(-time.Duration(settings.ConfirmationHoursBefore) * time.Hour)

		// Só agenda se for no futuro
		if confirmTime.After(time.Now()) {
			schedule := &models.NotificationSchedule{
				OrganizationId: reservation.OrganizationId,
				ProjectId:      reservation.ProjectId,
				EventType:      "confirmation_request",
				EntityType:     "reservation",
				EntityId:       reservation.Id,
				ScheduledFor:   confirmTime,
				Status:         "pending",
				Metadata:       string(metadataJSON),
			}

			if err := s.notificationRepo.CreateNotificationSchedule(schedule); err != nil {
				log.Printf("Erro ao agendar confirmação: %v", err)
			} else {
				log.Printf("Agendamento de confirmação criado para %s", confirmTime.Format("02/01/2006 15:04"))
			}
		}
	}

	// Agenda lembrete (se configurado e horas > 0)
	if settings.ReminderHoursBefore > 0 {
		reminderTime := reservationTime.Add(-time.Duration(settings.ReminderHoursBefore) * time.Hour)

		if reminderTime.After(time.Now()) {
			schedule := &models.NotificationSchedule{
				OrganizationId: reservation.OrganizationId,
				ProjectId:      reservation.ProjectId,
				EventType:      "reminder",
				EntityType:     "reservation",
				EntityId:       reservation.Id,
				ScheduledFor:   reminderTime,
				Status:         "pending",
				Metadata:       string(metadataJSON),
			}

			if err := s.notificationRepo.CreateNotificationSchedule(schedule); err != nil {
				log.Printf("Erro ao agendar lembrete: %v", err)
			} else {
				log.Printf("Agendamento de lembrete criado para %s", reminderTime.Format("02/01/2006 15:04"))
			}
		}
	}

	return nil
}

// CancelReservationSchedules cancela agendamentos pendentes de uma reserva
func (s *NotificationScheduleService) CancelReservationSchedules(reservationId uuid.UUID) error {
	return s.notificationRepo.CancelSchedulesByEntity("reservation", reservationId)
}

// ProcessDueSchedules processa agendamentos que chegaram na hora
func (s *NotificationScheduleService) ProcessDueSchedules() error {
	schedules, err := s.notificationRepo.GetDueSchedules(time.Now())
	if err != nil {
		return err
	}

	log.Printf("Processando %d agendamentos pendentes", len(schedules))

	for _, schedule := range schedules {
		if err := s.processSchedule(&schedule); err != nil {
			log.Printf("Erro ao processar agendamento %s: %v", schedule.Id, err)
			continue
		}
	}

	return nil
}

// processSchedule processa um agendamento individual
func (s *NotificationScheduleService) processSchedule(schedule *models.NotificationSchedule) error {
	// Busca a reserva
	reservation, err := s.reservationRepo.GetReservationById(schedule.EntityId)
	if err != nil {
		log.Printf("Reserva não encontrada para agendamento %s: %v", schedule.Id, err)
		return s.markScheduleStatus(schedule.Id, "skipped")
	}

	// Pula se reserva já foi cancelada ou completada
	if reservation.Status == "cancelled" || reservation.Status == "completed" || reservation.Status == "no_show" {
		log.Printf("Pulando agendamento %s - reserva com status %s", schedule.Id, reservation.Status)
		return s.markScheduleStatus(schedule.Id, "skipped")
	}

	// Busca cliente e mesa
	customer, err := s.customerRepo.GetCustomerById(reservation.CustomerId)
	if err != nil {
		log.Printf("Cliente não encontrado para agendamento %s: %v", schedule.Id, err)
		return s.markScheduleStatus(schedule.Id, "skipped")
	}

	table, err := s.tableRepo.GetTableById(reservation.TableId)
	if err != nil {
		log.Printf("Mesa não encontrada para agendamento %s: %v", schedule.Id, err)
		return s.markScheduleStatus(schedule.Id, "skipped")
	}

	// Envia notificação baseado no tipo de evento
	var triggerErr error
	switch schedule.EventType {
	case "confirmation_request":
		triggerErr = s.eventService.TriggerConfirmation24h(
			schedule.OrganizationId,
			schedule.ProjectId,
			reservation,
			customer,
			table,
		)
	case "reminder":
		// Usa o mesmo trigger de confirmação para lembrete
		// Pode ser expandido para template diferente
		triggerErr = s.eventService.TriggerConfirmation24h(
			schedule.OrganizationId,
			schedule.ProjectId,
			reservation,
			customer,
			table,
		)
	default:
		log.Printf("Tipo de evento desconhecido: %s", schedule.EventType)
		return s.markScheduleStatus(schedule.Id, "skipped")
	}

	if triggerErr != nil {
		log.Printf("Erro ao disparar notificação do agendamento %s: %v", schedule.Id, triggerErr)
		return s.markScheduleStatus(schedule.Id, "failed")
	}

	log.Printf("Agendamento %s processado com sucesso", schedule.Id)
	return s.markScheduleStatus(schedule.Id, "sent")
}

// markScheduleStatus atualiza o status de um agendamento
func (s *NotificationScheduleService) markScheduleStatus(id uuid.UUID, status string) error {
	return s.notificationRepo.UpdateScheduleStatus(id, status)
}

// GetSchedulesByReservation retorna agendamentos de uma reserva
func (s *NotificationScheduleService) GetSchedulesByReservation(reservationId uuid.UUID) ([]models.NotificationSchedule, error) {
	return s.notificationRepo.GetSchedulesByReservation(reservationId)
}
