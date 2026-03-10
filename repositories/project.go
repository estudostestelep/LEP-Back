package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

type IProjectRepository interface {
	GetProjectById(id uuid.UUID) (*models.Project, error)
	GetProjectByOrganization(orgId uuid.UUID) ([]models.Project, error)
	CreateProject(project *models.Project) error
	UpdateProject(project *models.Project) error
	SoftDeleteProject(id uuid.UUID) error
	GetActiveProjects(orgId uuid.UUID) ([]models.Project, error)
	GetProjectBySlug(orgId uuid.UUID, slug string) (*models.Project, error)
	GetDefaultProject(orgId uuid.UUID) (*models.Project, error)
	ProjectSlugExists(orgId uuid.UUID, slug string) (bool, error)
	SetDefaultProject(orgId uuid.UUID, projectId uuid.UUID) error
	ClearDefaultProject(orgId uuid.UUID) error
}

func NewProjectRepository(db *gorm.DB) IProjectRepository {
	return &ProjectRepository{db: db}
}

// GetProjectById busca projeto por ID
func (r *ProjectRepository) GetProjectById(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := r.db.First(&project, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// GetProjectByOrganization busca projetos por organização
func (r *ProjectRepository) GetProjectByOrganization(orgId uuid.UUID) ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Where("organization_id = ? AND deleted_at IS NULL", orgId).Find(&projects).Error
	return projects, err
}

// CreateProject cria novo projeto
func (r *ProjectRepository) CreateProject(project *models.Project) error {
	return r.db.Create(project).Error
}

// UpdateProject atualiza projeto existente
func (r *ProjectRepository) UpdateProject(project *models.Project) error {
	project.UpdatedAt = time.Now()
	return r.db.Save(project).Error
}

// SoftDeleteProject remove projeto logicamente
func (r *ProjectRepository) SoftDeleteProject(id uuid.UUID) error {
	return r.db.Model(&models.Project{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

// GetActiveProjects busca projetos ativos por organização
func (r *ProjectRepository) GetActiveProjects(orgId uuid.UUID) ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Where("organization_id = ? AND active = true AND deleted_at IS NULL", orgId).Find(&projects).Error
	return projects, err
}

// GetProjectBySlug busca projeto por slug dentro de uma organização
func (r *ProjectRepository) GetProjectBySlug(orgId uuid.UUID, slug string) (*models.Project, error) {
	var project models.Project
	err := r.db.First(&project, "organization_id = ? AND slug = ? AND deleted_at IS NULL", orgId, slug).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// GetDefaultProject busca o projeto padrão de uma organização
func (r *ProjectRepository) GetDefaultProject(orgId uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := r.db.First(&project, "organization_id = ? AND is_default = true AND deleted_at IS NULL", orgId).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// ProjectSlugExists verifica se um slug já existe na organização
func (r *ProjectRepository) ProjectSlugExists(orgId uuid.UUID, slug string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Project{}).
		Where("organization_id = ? AND slug = ? AND deleted_at IS NULL", orgId, slug).
		Count(&count).Error
	return count > 0, err
}

// SetDefaultProject define um projeto como padrão (e remove o padrão anterior)
func (r *ProjectRepository) SetDefaultProject(orgId uuid.UUID, projectId uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Remove is_default de todos os projetos da organização
		if err := tx.Model(&models.Project{}).
			Where("organization_id = ? AND deleted_at IS NULL", orgId).
			Update("is_default", false).Error; err != nil {
			return err
		}
		// Define o novo projeto como default
		if err := tx.Model(&models.Project{}).
			Where("id = ? AND organization_id = ?", projectId, orgId).
			Update("is_default", true).Error; err != nil {
			return err
		}
		return nil
	})
}

// ClearDefaultProject remove o projeto padrão de uma organização
func (r *ProjectRepository) ClearDefaultProject(orgId uuid.UUID) error {
	return r.db.Model(&models.Project{}).
		Where("organization_id = ? AND deleted_at IS NULL", orgId).
		Update("is_default", false).Error
}