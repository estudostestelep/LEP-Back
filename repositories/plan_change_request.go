package repositories

import (
	"errors"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlanChangeRequestRepository struct {
	db *gorm.DB
}

type IPlanChangeRequestRepository interface {
	// Create
	Create(request *models.PlanChangeRequest) error

	// Read
	GetById(id uuid.UUID) (*models.PlanChangeRequest, error)
	GetByOrganization(orgId uuid.UUID, status string) ([]*models.PlanChangeRequest, error)
	GetAllPending() ([]*models.PlanChangeRequest, error)
	GetAll(status string) ([]*models.PlanChangeRequest, error)
	GetByRequestedBy(userId uuid.UUID) ([]*models.PlanChangeRequest, error)

	// Update
	Update(request *models.PlanChangeRequest) error
	Approve(id uuid.UUID, reviewedBy uuid.UUID, reviewNotes string) error
	Reject(id uuid.UUID, reviewedBy uuid.UUID, reviewNotes string) error
	Cancel(id uuid.UUID) error

	// Delete
	Delete(id uuid.UUID) error
}

func NewPlanChangeRequestRepository(db *gorm.DB) IPlanChangeRequestRepository {
	return &PlanChangeRequestRepository{db: db}
}

// Create cria uma nova solicitação de mudança de plano
func (r *PlanChangeRequestRepository) Create(request *models.PlanChangeRequest) error {
	if request.Id == uuid.Nil {
		request.Id = uuid.New()
	}
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()
	request.Status = models.PlanChangeStatusPending

	return r.db.Create(request).Error
}

// GetById busca uma solicitação por ID
func (r *PlanChangeRequestRepository) GetById(id uuid.UUID) (*models.PlanChangeRequest, error) {
	var request models.PlanChangeRequest
	err := r.db.Where("id = ?", id).First(&request).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// GetByOrganization busca solicitações de uma organização
func (r *PlanChangeRequestRepository) GetByOrganization(orgId uuid.UUID, status string) ([]*models.PlanChangeRequest, error) {
	var requests []*models.PlanChangeRequest
	query := r.db.Where("organization_id = ?", orgId)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("created_at DESC").Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

// GetAllPending busca todas as solicitações pendentes (para admin)
func (r *PlanChangeRequestRepository) GetAllPending() ([]*models.PlanChangeRequest, error) {
	var requests []*models.PlanChangeRequest
	err := r.db.Where("status = ?", models.PlanChangeStatusPending).
		Order("created_at DESC").
		Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

// GetAll busca todas as solicitações com filtro opcional de status
func (r *PlanChangeRequestRepository) GetAll(status string) ([]*models.PlanChangeRequest, error) {
	var requests []*models.PlanChangeRequest
	query := r.db.Order("created_at DESC")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

// GetByRequestedBy busca solicitações criadas por um usuário específico
func (r *PlanChangeRequestRepository) GetByRequestedBy(userId uuid.UUID) ([]*models.PlanChangeRequest, error) {
	var requests []*models.PlanChangeRequest
	err := r.db.Where("requested_by = ?", userId).
		Order("created_at DESC").
		Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

// Update atualiza uma solicitação
func (r *PlanChangeRequestRepository) Update(request *models.PlanChangeRequest) error {
	request.UpdatedAt = time.Now()
	return r.db.Save(request).Error
}

// Approve aprova uma solicitação
func (r *PlanChangeRequestRepository) Approve(id uuid.UUID, reviewedBy uuid.UUID, reviewNotes string) error {
	request, err := r.GetById(id)
	if err != nil {
		return err
	}

	if !request.IsPending() {
		return errors.New("only pending requests can be approved")
	}

	now := time.Now()
	request.Status = models.PlanChangeStatusApproved
	request.ReviewedBy = &reviewedBy
	request.ReviewedAt = &now
	request.ReviewNotes = reviewNotes
	request.UpdatedAt = now

	return r.db.Save(request).Error
}

// Reject rejeita uma solicitação
func (r *PlanChangeRequestRepository) Reject(id uuid.UUID, reviewedBy uuid.UUID, reviewNotes string) error {
	request, err := r.GetById(id)
	if err != nil {
		return err
	}

	if !request.IsPending() {
		return errors.New("only pending requests can be rejected")
	}

	now := time.Now()
	request.Status = models.PlanChangeStatusRejected
	request.ReviewedBy = &reviewedBy
	request.ReviewedAt = &now
	request.ReviewNotes = reviewNotes
	request.UpdatedAt = now

	return r.db.Save(request).Error
}

// Cancel cancela uma solicitação (usuário pode cancelar sua própria solicitação)
func (r *PlanChangeRequestRepository) Cancel(id uuid.UUID) error {
	request, err := r.GetById(id)
	if err != nil {
		return err
	}

	if !request.IsPending() {
		return errors.New("only pending requests can be cancelled")
	}

	request.Status = models.PlanChangeStatusCancelled
	request.UpdatedAt = time.Now()

	return r.db.Save(request).Error
}

// Delete remove uma solicitação (soft delete)
func (r *PlanChangeRequestRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&models.PlanChangeRequest{}).Error
}
