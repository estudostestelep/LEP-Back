package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type resourceOrganization struct {
	db *gorm.DB
}

type IOrganizationRepository interface {
	GetOrganizationById(id uuid.UUID) (*models.Organization, error)
	GetOrganizationByEmail(email string) (*models.Organization, error)
	ListOrganizations() ([]models.Organization, error)
	ListActiveOrganizations() ([]models.Organization, error)
	CreateOrganization(organization *models.Organization) error
	UpdateOrganization(organization *models.Organization) error
	SoftDeleteOrganization(id uuid.UUID) error
	HardDeleteOrganization(id uuid.UUID) error
}

func NewConnOrganization(db *gorm.DB) IOrganizationRepository {
	return &resourceOrganization{db: db}
}

func (r *resourceOrganization) GetOrganizationById(id uuid.UUID) (*models.Organization, error) {
	var organization models.Organization
	err := r.db.Preload("Projects").First(&organization, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

func (r *resourceOrganization) GetOrganizationByEmail(email string) (*models.Organization, error) {
	var organization models.Organization
	err := r.db.First(&organization, "email = ? AND deleted_at IS NULL", email).Error
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

func (r *resourceOrganization) ListOrganizations() ([]models.Organization, error) {
	var organizations []models.Organization
	err := r.db.Where("deleted_at IS NULL").Order("created_at DESC").Find(&organizations).Error
	return organizations, err
}

func (r *resourceOrganization) ListActiveOrganizations() ([]models.Organization, error) {
	var organizations []models.Organization
	err := r.db.Where("active = ? AND deleted_at IS NULL", true).Order("created_at DESC").Find(&organizations).Error
	return organizations, err
}

func (r *resourceOrganization) CreateOrganization(organization *models.Organization) error {
	// Generate UUID if not provided
	if organization.Id == uuid.Nil {
		organization.Id = uuid.New()
	}

	// Set default values
	organization.Active = true
	organization.CreatedAt = time.Now()
	organization.UpdatedAt = time.Now()

	return r.db.Create(organization).Error
}

func (r *resourceOrganization) UpdateOrganization(organization *models.Organization) error {
	organization.UpdatedAt = time.Now()
	return r.db.Save(organization).Error
}

func (r *resourceOrganization) SoftDeleteOrganization(id uuid.UUID) error {
	return r.db.Model(&models.Organization{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

func (r *resourceOrganization) HardDeleteOrganization(id uuid.UUID) error {
	return r.db.Unscoped().Delete(&models.Organization{}, "id = ?", id).Error
}