package handler

import (
	"errors"
	"lep/repositories"
	"lep/repositories/models"

	"github.com/google/uuid"
)

type PlanChangeRequestHandler struct {
	planChangeRequestRepo repositories.IPlanChangeRequestRepository
	roleHandler          *RoleHandler
}

type IPlanChangeRequestHandler interface {
	// Client endpoints
	CreateRequest(request *models.PlanChangeRequest, userId string) error
	GetMyRequests(userId string) ([]*models.PlanChangeRequest, error)
	GetRequestById(id string, userId string) (*models.PlanChangeRequest, error)
	CancelRequest(id string, userId string) error

	// Admin endpoints
	GetAllRequests(status string) ([]*models.PlanChangeRequest, error)
	GetPendingRequests() ([]*models.PlanChangeRequest, error)
	GetRequestsByOrganization(orgId string, status string) ([]*models.PlanChangeRequest, error)
	ApproveRequest(id string, reviewedBy string, reviewNotes string) (*models.PlanChangeRequest, error)
	RejectRequest(id string, reviewedBy string, reviewNotes string) (*models.PlanChangeRequest, error)
}

func NewPlanChangeRequestHandler(
	planChangeRequestRepo repositories.IPlanChangeRequestRepository,
	roleHandler *RoleHandler,
) IPlanChangeRequestHandler {
	return &PlanChangeRequestHandler{
		planChangeRequestRepo: planChangeRequestRepo,
		roleHandler:          roleHandler,
	}
}

// ==================== Client Methods ====================

// CreateRequest cria uma nova solicitação de mudança de plano
func (h *PlanChangeRequestHandler) CreateRequest(request *models.PlanChangeRequest, userId string) error {
	if request == nil {
		return errors.New("request cannot be nil")
	}

	// Validações básicas
	if request.RequestedPackageId == uuid.Nil {
		return errors.New("requested package ID is required")
	}

	if request.OrganizationId == uuid.Nil {
		return errors.New("organization ID is required")
	}

	// Parse userId
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return errors.New("invalid user ID")
	}

	request.RequestedBy = userUUID
	return h.planChangeRequestRepo.Create(request)
}

// GetMyRequests retorna as solicitações do usuário logado
func (h *PlanChangeRequestHandler) GetMyRequests(userId string) ([]*models.PlanChangeRequest, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	return h.planChangeRequestRepo.GetByRequestedBy(userUUID)
}

// GetRequestById busca uma solicitação por ID
func (h *PlanChangeRequestHandler) GetRequestById(id string, userId string) (*models.PlanChangeRequest, error) {
	requestUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid request ID")
	}

	request, err := h.planChangeRequestRepo.GetById(requestUUID)
	if err != nil {
		return nil, err
	}

	// Verificar se o usuário tem permissão para ver esta solicitação
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Cliente só pode ver suas próprias solicitações
	if request.RequestedBy != userUUID {
		// Aqui você pode adicionar verificação de admin
		return nil, errors.New("unauthorized to view this request")
	}

	return request, nil
}

// CancelRequest cancela uma solicitação (apenas o usuário que criou pode cancelar)
func (h *PlanChangeRequestHandler) CancelRequest(id string, userId string) error {
	requestUUID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid request ID")
	}

	request, err := h.planChangeRequestRepo.GetById(requestUUID)
	if err != nil {
		return err
	}

	// Verificar se o usuário é o dono da solicitação
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return errors.New("invalid user ID")
	}

	if request.RequestedBy != userUUID {
		return errors.New("unauthorized to cancel this request")
	}

	return h.planChangeRequestRepo.Cancel(requestUUID)
}

// ==================== Admin Methods ====================

// GetAllRequests retorna todas as solicitações (admin)
func (h *PlanChangeRequestHandler) GetAllRequests(status string) ([]*models.PlanChangeRequest, error) {
	return h.planChangeRequestRepo.GetAll(status)
}

// GetPendingRequests retorna todas as solicitações pendentes (admin)
func (h *PlanChangeRequestHandler) GetPendingRequests() ([]*models.PlanChangeRequest, error) {
	return h.planChangeRequestRepo.GetAllPending()
}

// GetRequestsByOrganization retorna solicitações de uma organização específica (admin)
func (h *PlanChangeRequestHandler) GetRequestsByOrganization(orgId string, status string) ([]*models.PlanChangeRequest, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, errors.New("invalid organization ID")
	}

	return h.planChangeRequestRepo.GetByOrganization(orgUUID, status)
}

// ApproveRequest aprova uma solicitação (admin)
func (h *PlanChangeRequestHandler) ApproveRequest(id string, reviewedBy string, reviewNotes string) (*models.PlanChangeRequest, error) {
	requestUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid request ID")
	}

	reviewedByUUID, err := uuid.Parse(reviewedBy)
	if err != nil {
		return nil, errors.New("invalid reviewer ID")
	}

	err = h.planChangeRequestRepo.Approve(requestUUID, reviewedByUUID, reviewNotes)
	if err != nil {
		return nil, err
	}

	// Retornar a solicitação atualizada
	return h.planChangeRequestRepo.GetById(requestUUID)
}

// RejectRequest rejeita uma solicitação (admin)
func (h *PlanChangeRequestHandler) RejectRequest(id string, reviewedBy string, reviewNotes string) (*models.PlanChangeRequest, error) {
	requestUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid request ID")
	}

	reviewedByUUID, err := uuid.Parse(reviewedBy)
	if err != nil {
		return nil, errors.New("invalid reviewer ID")
	}

	err = h.planChangeRequestRepo.Reject(requestUUID, reviewedByUUID, reviewNotes)
	if err != nil {
		return nil, err
	}

	// Retornar a solicitação atualizada
	return h.planChangeRequestRepo.GetById(requestUUID)
}
