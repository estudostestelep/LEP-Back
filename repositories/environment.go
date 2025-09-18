package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EnvironmentRepository struct {
	db *gorm.DB
}

type IEnvironmentRepository interface {
	GetEnvironmentById(id uuid.UUID) (*models.Environment, error)
	GetEnvironmentsByProject(orgId, projectId uuid.UUID) ([]models.Environment, error)
	CreateEnvironment(environment *models.Environment) error
	UpdateEnvironment(environment *models.Environment) error
	SoftDeleteEnvironment(id uuid.UUID) error
	GetActiveEnvironments(orgId, projectId uuid.UUID) ([]models.Environment, error)
}

func NewEnvironmentRepository(db *gorm.DB) IEnvironmentRepository {
	return &EnvironmentRepository{db: db}
}

// GetEnvironmentById busca ambiente por ID
func (r *EnvironmentRepository) GetEnvironmentById(id uuid.UUID) (*models.Environment, error) {
	var environment models.Environment
	err := r.db.First(&environment, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &environment, nil
}

// GetEnvironmentsByProject busca ambientes por projeto
func (r *EnvironmentRepository) GetEnvironmentsByProject(orgId, projectId uuid.UUID) ([]models.Environment, error) {
	var environments []models.Environment
	err := r.db.Where("organization_id = ? AND project_id = ? AND deleted_at IS NULL", orgId, projectId).Find(&environments).Error
	return environments, err
}

// CreateEnvironment cria novo ambiente
func (r *EnvironmentRepository) CreateEnvironment(environment *models.Environment) error {
	return r.db.Create(environment).Error
}

// UpdateEnvironment atualiza ambiente existente
func (r *EnvironmentRepository) UpdateEnvironment(environment *models.Environment) error {
	environment.UpdatedAt = time.Now()
	return r.db.Save(environment).Error
}

// SoftDeleteEnvironment remove ambiente logicamente
func (r *EnvironmentRepository) SoftDeleteEnvironment(id uuid.UUID) error {
	return r.db.Model(&models.Environment{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

// GetActiveEnvironments busca ambientes ativos por projeto
func (r *EnvironmentRepository) GetActiveEnvironments(orgId, projectId uuid.UUID) ([]models.Environment, error) {
	var environments []models.Environment
	err := r.db.Where("organization_id = ? AND project_id = ? AND active = true AND deleted_at IS NULL", orgId, projectId).Find(&environments).Error
	return environments, err
}