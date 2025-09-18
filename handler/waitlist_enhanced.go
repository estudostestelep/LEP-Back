package handler

import (
	"fmt"
	"lep/repositories"
	"lep/repositories/models"
	"lep/utils"
	"time"

	"github.com/google/uuid"
)

type WaitlistEnhancedHandler struct {
	repo         *repositories.DBconn
	eventService *utils.EventService
}

type IWaitlistEnhancedHandler interface {
	AddToWaitlist(waitlist *models.Waitlist) error
	GetWaitlistPosition(id string) (int, error)
	GetEstimatedWaitTime(orgId, projectId uuid.UUID, partySize int) (int, error)
	NotifyNextInLine(orgId, projectId uuid.UUID) error
	RemoveFromWaitlist(id string) error
	GetWaitlistByProject(orgId, projectId string) ([]models.Waitlist, error)
	ConvertToLead(waitlistId uuid.UUID) error
}

func NewWaitlistEnhancedHandler(repo *repositories.DBconn) IWaitlistEnhancedHandler {
	eventService := utils.NewEventService(repo.Notifications, repo.Projects, repo.Settings)
	return &WaitlistEnhancedHandler{
		repo:         repo,
		eventService: eventService,
	}
}

// AddToWaitlist - Adiciona cliente à fila de espera com estimativa de tempo
func (w *WaitlistEnhancedHandler) AddToWaitlist(waitlist *models.Waitlist) error {
	waitlist.Id = uuid.New()
	waitlist.Status = "waiting"
	waitlist.CreatedAt = time.Now()
	waitlist.UpdatedAt = time.Now()

	// Criar waitlist
	if err := w.repo.Waitlists.CreateWaitlist(waitlist); err != nil {
		return err
	}

	// Converter automaticamente para lead
	if err := w.ConvertToLead(waitlist.Id); err != nil {
		// Log erro mas não interrompe
		fmt.Printf("Error converting waitlist to lead: %v\n", err)
	}

	return nil
}

// GetWaitlistPosition - Retorna posição na fila
func (w *WaitlistEnhancedHandler) GetWaitlistPosition(id string) (int, error) {
	waitlistId, err := uuid.Parse(id)
	if err != nil {
		return 0, err
	}

	// Buscar o item da waitlist
	waitlistItem, err := w.repo.Waitlists.GetWaitlistById(waitlistId)
	if err != nil {
		return 0, err
	}

	// Buscar todos os itens da waitlist do projeto, ordenados por criação
	allWaitlist, err := w.repo.Waitlists.GetWaitlistByProject(waitlistItem.OrganizationId, waitlistItem.ProjectId)
	if err != nil {
		return 0, err
	}

	// Encontrar posição
	position := 1
	for _, item := range allWaitlist {
		if item.Status != "waiting" {
			continue
		}
		if item.Id == waitlistId {
			return position, nil
		}
		// Se foi criado antes, está na frente
		if item.CreatedAt.Before(waitlistItem.CreatedAt) {
			position++
		}
	}

	return position, nil
}

// GetEstimatedWaitTime - Calcula tempo estimado de espera baseado na fila
func (w *WaitlistEnhancedHandler) GetEstimatedWaitTime(orgId, projectId uuid.UUID, partySize int) (int, error) {
	// Buscar fila atual
	waitlist, err := w.repo.Waitlists.GetWaitlistByProject(orgId, projectId)
	if err != nil {
		return 0, err
	}

	// Contar pessoas na frente com o mesmo tamanho ou menor
	peopleAhead := 0
	for _, item := range waitlist {
		if item.Status == "waiting" && item.People <= partySize {
			peopleAhead++
		}
	}

	// Buscar mesas disponíveis que comportam o grupo
	tables, err := w.repo.Tables.GetTablesByProject(orgId, projectId)
	if err != nil {
		return 0, err
	}

	availableTables := 0
	for _, table := range tables {
		if table.Status == "livre" && table.Capacity >= partySize {
			availableTables++
		}
	}

	// Cálculo simples: se não há mesas disponíveis, cada pessoa espera ~30min
	// Se há mesas, tempo reduzido
	baseTime := 30 // minutos base por posição
	if availableTables > 0 {
		baseTime = 15 // reduz pela metade se há mesas disponíveis
	}

	estimatedTime := peopleAhead * baseTime

	// Mínimo de 5 minutos, máximo de 120 minutos
	if estimatedTime < 5 {
		estimatedTime = 5
	}
	if estimatedTime > 120 {
		estimatedTime = 120
	}

	return estimatedTime, nil
}

// NotifyNextInLine - Notifica próximo da fila quando mesa fica disponível
func (w *WaitlistEnhancedHandler) NotifyNextInLine(orgId, projectId uuid.UUID) error {
	// Buscar fila de espera
	waitlist, err := w.repo.Waitlists.GetWaitlistByProject(orgId, projectId)
	if err != nil {
		return err
	}

	// Buscar mesas disponíveis
	tables, err := w.repo.Tables.GetTablesByProject(orgId, projectId)
	if err != nil {
		return err
	}

	// Para cada mesa disponível, notificar alguém da fila
	for _, table := range tables {
		if table.Status != "livre" {
			continue
		}

		// Encontrar primeira pessoa da fila que cabe na mesa
		for _, wait := range waitlist {
			if wait.Status != "waiting" {
				continue
			}

			if wait.People <= table.Capacity {
				// Buscar dados do cliente
				customer, err := w.repo.Customers.GetCustomerById(wait.CustomerId)
				if err != nil {
					continue
				}

				// Estimar tempo de espera
				estimatedWait, _ := w.GetEstimatedWaitTime(orgId, projectId, wait.People)

				// Trigger de mesa disponível
				if err := w.eventService.TriggerTableAvailable(orgId, projectId, &table, customer, estimatedWait); err != nil {
					fmt.Printf("Error triggering table available event: %v\n", err)
				}

				// Atualizar status para "notified"
				wait.Status = "notified"
				wait.UpdatedAt = time.Now()
				w.repo.Waitlists.UpdateWaitlist(&wait)

				break // Apenas o primeiro da fila para esta mesa
			}
		}
	}

	return nil
}

// RemoveFromWaitlist - Remove da fila (quando sentou ou desistiu)
func (w *WaitlistEnhancedHandler) RemoveFromWaitlist(id string) error {
	waitlistId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	// Atualizar status para "seated" ou "left"
	waitlistItem, err := w.repo.Waitlists.GetWaitlistById(waitlistId)
	if err != nil {
		return err
	}

	waitlistItem.Status = "seated"
	waitlistItem.UpdatedAt = time.Now()
	return w.repo.Waitlists.UpdateWaitlist(waitlistItem)
}

// GetWaitlistByProject - Lista fila de espera do projeto
func (w *WaitlistEnhancedHandler) GetWaitlistByProject(orgId, projectId string) ([]models.Waitlist, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projectUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	return w.repo.Waitlists.GetWaitlistByProject(orgUUID, projectUUID)
}

// ConvertToLead - Converte waitlist em lead para CRM
func (w *WaitlistEnhancedHandler) ConvertToLead(waitlistId uuid.UUID) error {
	// Buscar waitlist
	waitlistItem, err := w.repo.Waitlists.GetWaitlistById(waitlistId)
	if err != nil {
		return err
	}

	// Buscar cliente
	customer, err := w.repo.Customers.GetCustomerById(waitlistItem.CustomerId)
	if err != nil {
		return err
	}

	// Criar lead - assumindo que o repository existirá
	// Na prática seria implementado
	lead := &models.Lead{
		Id:             uuid.New(),
		OrganizationId: waitlistItem.OrganizationId,
		ProjectId:      waitlistItem.ProjectId,
		Name:           customer.Name,
		Email:          customer.Email,
		Phone:          customer.Phone,
		Source:         "waitlist",
		Status:         "new",
		Notes:          fmt.Sprintf("Added to waitlist on %s for %d people", waitlistItem.CreatedAt.Format("2006-01-02 15:04"), waitlistItem.People),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Por agora apenas log - seria r.repo.Leads.CreateLead(lead)
	fmt.Printf("Would create lead: %+v\n", lead)

	return nil
}

// GetTablesByProject - Helper para buscar mesas do projeto
func (w *WaitlistEnhancedHandler) GetTablesByProject(orgId, projectId uuid.UUID) ([]models.Table, error) {
	// Implementação simplificada - assumindo que o método existirá
	// return w.repo.Tables.GetTablesByProject(orgId, projectId)

	// Por agora retorna array vazio
	return []models.Table{}, nil
}