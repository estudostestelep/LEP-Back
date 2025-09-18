package handler

import (
	"fmt"
	"lep/repositories"
	"lep/repositories/models"
	"lep/utils"
	"time"

	"github.com/google/uuid"
)

type ReservationEnhancedHandler struct {
	repo         *repositories.DBconn
	eventService *utils.EventService
}

type IReservationEnhancedHandler interface {
	GetReservation(id string) (*models.Reservation, error)
	CreateReservationWithTriggers(reservation *models.Reservation) error
	UpdateReservationWithTriggers(updatedReservation *models.Reservation) error
	CancelReservationWithTriggers(id, reason string) error
	DeleteReservation(id string) error
	ListReservations(orgId, projectId string) ([]models.Reservation, error)
	ValidateReservation(reservation *models.Reservation) error
}

func NewReservationEnhancedHandler(repo *repositories.DBconn) IReservationEnhancedHandler {
	eventService := utils.NewEventService(repo.Notifications, repo.Projects, repo.Settings)
	return &ReservationEnhancedHandler{
		repo:         repo,
		eventService: eventService,
	}
}

func (r *ReservationEnhancedHandler) GetReservation(id string) (*models.Reservation, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return r.repo.Reservations.GetReservationById(uuid)
}

// CreateReservationWithTriggers - Cria reserva com validações e triggers automáticos
func (r *ReservationEnhancedHandler) CreateReservationWithTriggers(reservation *models.Reservation) error {
	// Validar reserva
	if err := r.ValidateReservation(reservation); err != nil {
		return err
	}

	// Configurar reserva
	reservation.Id = uuid.New()
	reservation.Status = "confirmed"
	reservation.CreatedAt = time.Now()
	reservation.UpdatedAt = time.Now()

	// Criar reserva
	if err := r.repo.Reservations.CreateReservation(reservation); err != nil {
		return err
	}

	// Atualizar status da mesa para "reservada"
	if err := r.updateTableStatus(reservation.TableId, "reservada"); err != nil {
		// Log do erro mas não interrompe o processo
		fmt.Printf("Error updating table status: %v\n", err)
	}

	// Buscar dados adicionais para o trigger
	customer, err := r.repo.Customers.GetCustomerById(reservation.CustomerId)
	if err != nil {
		return fmt.Errorf("customer not found: %w", err)
	}

	table, err := r.repo.Tables.GetTableById(reservation.TableId)
	if err != nil {
		return fmt.Errorf("table not found: %w", err)
	}

	// Trigger de notificação
	if err := r.eventService.TriggerReservationCreated(reservation.OrganizationId, reservation.ProjectId, reservation, customer, table); err != nil {
		// Log do erro mas não interrompe o processo
		fmt.Printf("Error triggering reservation created event: %v\n", err)
	}

	return nil
}

// UpdateReservationWithTriggers - Atualiza reserva com triggers automáticos
func (r *ReservationEnhancedHandler) UpdateReservationWithTriggers(updatedReservation *models.Reservation) error {
	// Validar reserva
	if err := r.ValidateReservation(updatedReservation); err != nil {
		return err
	}

	// Buscar reserva atual para comparação
	currentReservation, err := r.repo.Reservations.GetReservationById(updatedReservation.Id)
	if err != nil {
		return fmt.Errorf("current reservation not found: %w", err)
	}

	// Se mudou a mesa, atualizar status das mesas
	if currentReservation.TableId != updatedReservation.TableId {
		// Liberar mesa anterior
		if err := r.updateTableStatus(currentReservation.TableId, "livre"); err != nil {
			fmt.Printf("Error freeing previous table: %v\n", err)
		}
		// Reservar nova mesa
		if err := r.updateTableStatus(updatedReservation.TableId, "reservada"); err != nil {
			fmt.Printf("Error reserving new table: %v\n", err)
		}
	}

	updatedReservation.UpdatedAt = time.Now()
	if err := r.repo.Reservations.UpdateReservation(updatedReservation); err != nil {
		return err
	}

	// Buscar dados adicionais para o trigger
	customer, err := r.repo.Customers.GetCustomerById(updatedReservation.CustomerId)
	if err != nil {
		return fmt.Errorf("customer not found: %w", err)
	}

	table, err := r.repo.Tables.GetTableById(updatedReservation.TableId)
	if err != nil {
		return fmt.Errorf("table not found: %w", err)
	}

	// Trigger de notificação
	if err := r.eventService.TriggerReservationUpdated(updatedReservation.OrganizationId, updatedReservation.ProjectId, updatedReservation, customer, table); err != nil {
		fmt.Printf("Error triggering reservation updated event: %v\n", err)
	}

	return nil
}

// CancelReservationWithTriggers - Cancela reserva com triggers automáticos
func (r *ReservationEnhancedHandler) CancelReservationWithTriggers(id, reason string) error {
	reservationId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	// Buscar reserva atual
	reservation, err := r.repo.Reservations.GetReservationById(reservationId)
	if err != nil {
		return fmt.Errorf("reservation not found: %w", err)
	}

	// Atualizar status para cancelada
	reservation.Status = "cancelled"
	reservation.UpdatedAt = time.Now()
	if err := r.repo.Reservations.UpdateReservation(reservation); err != nil {
		return err
	}

	// Liberar mesa
	if err := r.updateTableStatus(reservation.TableId, "livre"); err != nil {
		fmt.Printf("Error freeing table after cancellation: %v\n", err)
	}

	// Buscar dados adicionais para o trigger
	customer, err := r.repo.Customers.GetCustomerById(reservation.CustomerId)
	if err != nil {
		return fmt.Errorf("customer not found: %w", err)
	}

	table, err := r.repo.Tables.GetTableById(reservation.TableId)
	if err != nil {
		return fmt.Errorf("table not found: %w", err)
	}

	// Trigger de notificação
	if err := r.eventService.TriggerReservationCancelled(reservation.OrganizationId, reservation.ProjectId, reservation, customer, table); err != nil {
		fmt.Printf("Error triggering reservation cancelled event: %v\n", err)
	}

	// Verificar fila de espera para a mesa que foi liberada
	r.checkWaitlistForTable(reservation.OrganizationId, reservation.ProjectId, reservation.TableId)

	return nil
}

func (r *ReservationEnhancedHandler) DeleteReservation(id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.repo.Reservations.DeleteReservation(uuid)
}

func (r *ReservationEnhancedHandler) ListReservations(orgId, projectId string) ([]models.Reservation, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projectUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	return r.repo.Reservations.GetReservationsByProject(orgUUID, projectUUID)
}

// ValidateReservation - Valida regras de negócio para reservas
func (r *ReservationEnhancedHandler) ValidateReservation(reservation *models.Reservation) error {
	// Buscar configurações do projeto
	settings, err := r.repo.Settings.GetSettingsByProject(reservation.OrganizationId, reservation.ProjectId)
	if err != nil {
		return fmt.Errorf("settings not found: %w", err)
	}

	// Validar antecedência mínima
	now := time.Now()
	minAdvanceTime := now.Add(time.Duration(settings.MinAdvanceHours) * time.Hour)
	if reservation.Datetime.Before(minAdvanceTime) {
		return fmt.Errorf("reservation must be at least %d hours in advance", settings.MinAdvanceHours)
	}

	// Validar antecedência máxima
	maxAdvanceTime := now.Add(time.Duration(settings.MaxAdvanceDays) * 24 * time.Hour)
	if reservation.Datetime.After(maxAdvanceTime) {
		return fmt.Errorf("reservation cannot be more than %d days in advance", settings.MaxAdvanceDays)
	}

	// Validar períodos bloqueados
	if err := r.checkBlockedPeriods(reservation.OrganizationId, reservation.ProjectId, reservation.Datetime); err != nil {
		return err
	}

	// Validar se mesa existe e está disponível
	table, err := r.repo.Tables.GetTableById(reservation.TableId)
	if err != nil {
		return fmt.Errorf("table not found: %w", err)
	}

	// Validar capacidade da mesa
	if reservation.PartySize > table.Capacity {
		return fmt.Errorf("party size (%d) exceeds table capacity (%d)", reservation.PartySize, table.Capacity)
	}

	// Validar conflitos de horário (só para mesa não "livre")
	if table.Status != "livre" {
		// Verificar se há conflito de horário
		if err := r.checkTimeConflicts(reservation); err != nil {
			return err
		}
	}

	return nil
}

// checkBlockedPeriods - Verifica se a data/hora está em período bloqueado
func (r *ReservationEnhancedHandler) checkBlockedPeriods(orgId, projectId uuid.UUID, datetime time.Time) error {
	// Implementação simplificada - assumindo que o repository existirá
	// Na prática seria: r.repo.BlockedPeriods.CheckPeriodBlocked(orgId, projectId, datetime)

	// Por agora, uma validação básica - não permitir reservas em domingos
	if datetime.Weekday() == time.Sunday {
		return fmt.Errorf("reservations are not allowed on Sundays")
	}

	// Não permitir reservas muito tarde (após 23h)
	if datetime.Hour() >= 23 {
		return fmt.Errorf("reservations are not allowed after 11 PM")
	}

	// Não permitir reservas muito cedo (antes das 10h)
	if datetime.Hour() < 10 {
		return fmt.Errorf("reservations are not allowed before 10 AM")
	}

	return nil
}

// updateTableStatus - Atualiza status da mesa
func (r *ReservationEnhancedHandler) updateTableStatus(tableId uuid.UUID, status string) error {
	table, err := r.repo.Tables.GetTableById(tableId)
	if err != nil {
		return err
	}

	table.Status = status
	table.UpdatedAt = time.Now()
	return r.repo.Tables.UpdateTable(table)
}

// checkTimeConflicts - Verifica conflitos de horário para a mesa
func (r *ReservationEnhancedHandler) checkTimeConflicts(reservation *models.Reservation) error {
	// Buscar reservas existentes para a mesa no mesmo dia
	dayStart := time.Date(reservation.Datetime.Year(), reservation.Datetime.Month(), reservation.Datetime.Day(), 0, 0, 0, 0, reservation.Datetime.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	existingReservations, err := r.repo.Reservations.GetReservationsByTableAndDateRange(reservation.TableId, dayStart, dayEnd)
	if err != nil {
		return err
	}

	// Verificar conflitos (assumindo 2 horas por reserva)
	reservationStart := reservation.Datetime
	reservationEnd := reservationStart.Add(2 * time.Hour)

	for _, existing := range existingReservations {
		// Pular a própria reserva se for uma atualização
		if existing.Id == reservation.Id {
			continue
		}

		// Pular reservas canceladas
		if existing.Status == "cancelled" {
			continue
		}

		existingStart := existing.Datetime
		existingEnd := existingStart.Add(2 * time.Hour)

		// Verificar sobreposição
		if reservationStart.Before(existingEnd) && reservationEnd.After(existingStart) {
			return fmt.Errorf("time conflict with existing reservation at %s", existing.Datetime.Format("15:04"))
		}
	}

	return nil
}

// checkWaitlistForTable - Verifica fila de espera quando mesa fica disponível
func (r *ReservationEnhancedHandler) checkWaitlistForTable(orgId, projectId, tableId uuid.UUID) {
	// Buscar pessoas na fila de espera
	waitlist, err := r.repo.Waitlists.GetWaitlistByProject(orgId, projectId)
	if err != nil {
		return
	}

	// Buscar mesa
	table, err := r.repo.Tables.GetTableById(tableId)
	if err != nil {
		return
	}

	// Encontrar primeira pessoa da fila que cabe na mesa
	for _, wait := range waitlist {
		if wait.Status != "waiting" {
			continue
		}

		if wait.People <= table.Capacity {
			// Buscar dados do cliente
			customer, err := r.repo.Customers.GetCustomerById(wait.CustomerId)
			if err != nil {
				continue
			}

			// Estimar tempo de espera (simples: 5 minutos)
			estimatedWait := 5

			// Trigger de mesa disponível
			if err := r.eventService.TriggerTableAvailable(orgId, projectId, table, customer, estimatedWait); err != nil {
				fmt.Printf("Error triggering table available event: %v\n", err)
			}

			// Atualizar status na fila de espera para "notified"
			wait.Status = "notified"
			wait.UpdatedAt = time.Now()
			r.repo.Waitlists.UpdateWaitlist(&wait)

			break // Apenas o primeiro da fila
		}
	}
}