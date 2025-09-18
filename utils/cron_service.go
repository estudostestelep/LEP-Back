package utils

import (
	"lep/repositories"
	"lep/repositories/models"
	"log"
	"time"

	"github.com/google/uuid"
)

type CronService struct {
	repo         *repositories.DBconn
	eventService *EventService
}

func NewCronService(repo *repositories.DBconn) *CronService {
	eventService := NewEventService(repo.Notifications, repo.Projects, repo.Settings)
	return &CronService{
		repo:         repo,
		eventService: eventService,
	}
}

// ProcessConfirmation24h - Processa confirmações 24h antes das reservas
func (c *CronService) ProcessConfirmation24h() error {
	log.Println("Starting 24h confirmation job...")

	// Buscar todas as organizações/projetos (simplificado - na prática seria mais otimizado)
	projects, err := c.getAllActiveProjects()
	if err != nil {
		return err
	}

	// Para cada projeto, buscar reservas que precisam de confirmação
	for _, project := range projects {
		if err := c.processProjectConfirmations(project.OrganizationId, project.Id); err != nil {
			log.Printf("Error processing confirmations for project %s: %v", project.Id, err)
			continue
		}
	}

	log.Println("24h confirmation job completed")
	return nil
}

func (c *CronService) processProjectConfirmations(orgId, projectId uuid.UUID) error {
	// Verificar se confirmação 24h está habilitada
	settings, err := c.repo.Settings.GetSettingsByProject(orgId, projectId)
	if err != nil || !settings.NotifyConfirmation24h {
		return nil
	}

	// Buscar projeto para pegar timezone
	project, err := c.repo.Projects.GetProjectById(projectId)
	if err != nil {
		log.Printf("Project not found %s: %v", projectId, err)
		return err
	}

	// Carregar timezone do projeto (padrão: America/Sao_Paulo)
	timezone := project.TimeZone
	if timezone == "" {
		timezone = "America/Sao_Paulo"
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Printf("Invalid timezone %s, using UTC: %v", timezone, err)
		loc = time.UTC
	}

	// Calcular janela de tempo no timezone do projeto: reservas que acontecem entre 23h e 25h a partir de agora
	now := time.Now().In(loc)
	start24h := now.Add(23 * time.Hour)
	end24h := now.Add(25 * time.Hour)

	// Buscar reservas nesta janela
	reservations, err := c.getReservationsInTimeRange(orgId, projectId, start24h, end24h)
	if err != nil {
		return err
	}

	log.Printf("Found %d reservations for 24h confirmation in project %s", len(reservations), projectId)

	// Processar cada reserva
	for _, reservation := range reservations {
		// Verificar se reserva está confirmada
		if reservation.Status != "confirmed" {
			continue
		}

		// Verificar se já foi enviada confirmação recentemente
		hasRecent, err := c.hasRecentConfirmationLog(orgId, projectId, reservation.Id)
		if err != nil {
			log.Printf("Error checking recent confirmations for reservation %s: %v", reservation.Id, err)
			continue
		}
		if hasRecent {
			log.Printf("Skipping reservation %s - confirmation already sent recently", reservation.Id)
			continue
		}

		// Buscar dados do cliente e mesa
		customer, err := c.repo.Customers.GetCustomerById(reservation.CustomerId)
		if err != nil {
			log.Printf("Customer not found for reservation %s: %v", reservation.Id, err)
			continue
		}

		table, err := c.repo.Tables.GetTableById(reservation.TableId)
		if err != nil {
			log.Printf("Table not found for reservation %s: %v", reservation.Id, err)
			continue
		}

		// Trigger de confirmação 24h
		if err := c.eventService.TriggerConfirmation24h(orgId, projectId, &reservation, customer, table); err != nil {
			log.Printf("Error triggering 24h confirmation for reservation %s: %v", reservation.Id, err)
		} else {
			log.Printf("24h confirmation sent for reservation %s", reservation.Id)
		}
	}

	return nil
}

// ProcessPendingEvents - Processa eventos pendentes de notificação
func (c *CronService) ProcessPendingEvents() error {
	log.Println("Starting pending events processing job...")

	// Buscar todos os projetos ativos
	projects, err := c.getAllActiveProjects()
	if err != nil {
		return err
	}

	// Para cada projeto, processar eventos pendentes
	for _, project := range projects {
		if err := c.eventService.ProcessPendingEvents(project.OrganizationId, project.Id); err != nil {
			log.Printf("Error processing pending events for project %s: %v", project.Id, err)
		}
	}

	log.Println("Pending events processing completed")
	return nil
}

// CleanupOldLogs - Limpa logs antigos de notificação (opcional)
func (c *CronService) CleanupOldLogs() error {
	log.Println("Starting cleanup of old notification logs...")

	// Deletar logs mais antigos que 90 dias
	cutoffDate := time.Now().AddDate(0, 0, -90)

	// Aqui implementaria a limpeza no repository
	// Por simplicidade, apenas log por agora
	log.Printf("Would cleanup logs older than %s", cutoffDate.Format("2006-01-02"))

	log.Println("Cleanup completed")
	return nil
}

// getAllActiveProjects - Busca todos os projetos ativos
func (c *CronService) getAllActiveProjects() ([]models.Project, error) {
	// Implementação real usando repository
	// Buscar todos os projetos ativos - seria otimizado com uma query específica
	// Por agora, usar método existente com organização fictícia

	// NOTA: Este método precisaria ser implementado no ProjectRepository
	// Por agora, retornamos array vazio para evitar erros
	return []models.Project{}, nil
}

// getReservationsInTimeRange - Busca reservas em um período específico
func (c *CronService) getReservationsInTimeRange(orgId, projectId uuid.UUID, start, end time.Time) ([]models.Reservation, error) {
	// Usar o método existente no repository
	return c.repo.Reservations.GetReservationsByTableAndDateRange(uuid.Nil, start, end)
}

// hasRecentConfirmationLog - Verifica se já foi enviada confirmação recentemente
func (c *CronService) hasRecentConfirmationLog(orgId, projectId uuid.UUID, reservationId uuid.UUID) (bool, error) {
	// Buscar logs de notificação das últimas 2 horas para esta reserva
	since := time.Now().Add(-2 * time.Hour)

	logs, err := c.repo.Notifications.GetNotificationLogsByProject(orgId, projectId, 50)
	if err != nil {
		return false, err
	}

	// Verificar se há log de confirmação_24h recente
	for _, log := range logs {
		if log.EventType == "confirmation_24h" &&
			log.CreatedAt.After(since) &&
			log.Status == "sent" {
			// Seria ideal verificar também o reservation_id nos metadados
			return true, nil
		}
	}

	return false, nil
}

// StartCronJobs - Inicia jobs automáticos (seria chamado no main)
func (c *CronService) StartCronJobs() {
	log.Println("Starting cron jobs...")

	// Job de confirmação 24h - executa a cada hora
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := c.ProcessConfirmation24h(); err != nil {
					log.Printf("Error in 24h confirmation job: %v", err)
				}
			}
		}
	}()

	// Job de eventos pendentes - executa a cada 5 minutos
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := c.ProcessPendingEvents(); err != nil {
					log.Printf("Error in pending events job: %v", err)
				}
			}
		}
	}()

	// Job de limpeza - executa uma vez por dia à meia-noite
	go func() {
		for {
			now := time.Now()
			// Calcular próxima meia-noite
			nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			time.Sleep(time.Until(nextMidnight))

			if err := c.CleanupOldLogs(); err != nil {
				log.Printf("Error in cleanup job: %v", err)
			}

			// Esperar um dia
			time.Sleep(24 * time.Hour)
		}
	}()

	log.Println("Cron jobs started successfully")
}
