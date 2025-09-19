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