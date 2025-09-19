package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ILeadRepository interface {
	CreateLead(lead *models.Lead) error
	GetLeadById(id uuid.UUID) (*models.Lead, error)
	GetLeadsByProject(orgId, projectId uuid.UUID) ([]models.Lead, error)
	UpdateLead(lead *models.Lead) error
	SoftDeleteLead(id uuid.UUID) error
	GetLeadsByStatus(orgId, projectId uuid.UUID, status string) ([]models.Lead, error)
	GetLeadsBySource(orgId, projectId uuid.UUID, source string) ([]models.Lead, error)
	SearchLeads(orgId, projectId uuid.UUID, query string) ([]models.Lead, error)
}

type LeadRepository struct {
	db *gorm.DB
}

func NewLeadRepository(db *gorm.DB) ILeadRepository {
	return &LeadRepository{db: db}
}

func (r *LeadRepository) CreateLead(lead *models.Lead) error {
	return r.db.Create(lead).Error
}

func (r *LeadRepository) GetLeadById(id uuid.UUID) (*models.Lead, error) {
	var lead models.Lead
	err := r.db.First(&lead, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &lead, nil
}

func (r *LeadRepository) GetLeadsByProject(orgId, projectId uuid.UUID) ([]models.Lead, error) {
	var leads []models.Lead
	err := r.db.Where("organization_id = ? AND project_id = ? AND deleted_at IS NULL", orgId, projectId).
		Order("created_at DESC").Find(&leads).Error
	return leads, err
}

func (r *LeadRepository) UpdateLead(lead *models.Lead) error {
	lead.UpdatedAt = time.Now()
	return r.db.Save(lead).Error
}

func (r *LeadRepository) SoftDeleteLead(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.Lead{}).Where("id = ?", id).Update("deleted_at", now).Error
}

func (r *LeadRepository) GetLeadsByStatus(orgId, projectId uuid.UUID, status string) ([]models.Lead, error) {
	var leads []models.Lead
	err := r.db.Where("organization_id = ? AND project_id = ? AND status = ? AND deleted_at IS NULL", orgId, projectId, status).
		Order("created_at DESC").Find(&leads).Error
	return leads, err
}

func (r *LeadRepository) GetLeadsBySource(orgId, projectId uuid.UUID, source string) ([]models.Lead, error) {
	var leads []models.Lead
	err := r.db.Where("organization_id = ? AND project_id = ? AND source = ? AND deleted_at IS NULL", orgId, projectId, source).
		Order("created_at DESC").Find(&leads).Error
	return leads, err
}

func (r *LeadRepository) SearchLeads(orgId, projectId uuid.UUID, query string) ([]models.Lead, error) {
	var leads []models.Lead
	searchPattern := "%" + query + "%"
	err := r.db.Where("organization_id = ? AND project_id = ? AND deleted_at IS NULL", orgId, projectId).
		Where("name ILIKE ? OR email ILIKE ? OR phone ILIKE ?", searchPattern, searchPattern, searchPattern).
		Order("created_at DESC").Find(&leads).Error
	return leads, err
}