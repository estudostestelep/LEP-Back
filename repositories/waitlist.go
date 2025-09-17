package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Interface para o repositório de Waitlist
type WaitlistRepositoryInterface interface {
	CreateWaitlist(waitlist *models.Waitlist) error
	GetWaitlistById(id uuid.UUID) (*models.Waitlist, error)
	ListWaitlists(OrganizationId, projectId uuid.UUID) ([]models.Waitlist, error)
	UpdateWaitlist(waitlist *models.Waitlist) error
	SoftDeleteWaitlist(id uuid.UUID) error
}

type WaitlistRepository struct {
	db *gorm.DB
}

func NewWaitlistRepository(db *gorm.DB) WaitlistRepositoryInterface {
	return &WaitlistRepository{db}
}

func (r *WaitlistRepository) CreateWaitlist(waitlist *models.Waitlist) error {
	return r.db.Create(waitlist).Error
}

func (r *WaitlistRepository) GetWaitlistById(id uuid.UUID) (*models.Waitlist, error) {
	var waitlist models.Waitlist
	err := r.db.First(&waitlist, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &waitlist, nil
}

func (r *WaitlistRepository) ListWaitlists(OrganizationId, projectId uuid.UUID) ([]models.Waitlist, error) {
	var waitlists []models.Waitlist
	err := r.db.Where("org_id = ? AND project_id = ? AND deleted_at IS NULL", OrganizationId, projectId).Find(&waitlists).Error
	return waitlists, err
}

func (r *WaitlistRepository) UpdateWaitlist(waitlist *models.Waitlist) error {
	return r.db.Save(waitlist).Error
}

func (r *WaitlistRepository) SoftDeleteWaitlist(id uuid.UUID) error {
	return r.db.Model(&models.Waitlist{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}
