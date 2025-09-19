package handler

import (
	"lep/repositories"
	"lep/repositories/models"
	"lep/utils"
	"time"

	"github.com/google/uuid"
)

type ProjectHandler struct {
	projectRepo      repositories.IProjectRepository
	settingsRepo     repositories.ISettingsRepository
	notificationRepo repositories.INotificationRepository
}

type IProjectHandler interface {
	GetProjectById(id string) (*models.Project, error)
	GetProjectsByOrganization(orgId string) ([]models.Project, error)
	CreateProject(project *models.Project) error
	UpdateProject(project *models.Project) error
	SoftDeleteProject(id string) error
	GetActiveProjects(orgId string) ([]models.Project, error)
}

func NewProjectHandler(projectRepo repositories.IProjectRepository, settingsRepo repositories.ISettingsRepository, notificationRepo repositories.INotificationRepository) IProjectHandler {
	return &ProjectHandler{
		projectRepo:      projectRepo,
		settingsRepo:     settingsRepo,
		notificationRepo: notificationRepo,
	}
}

// GetProjectById busca projeto por ID
func (h *ProjectHandler) GetProjectById(id string) (*models.Project, error) {
	projectId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return h.projectRepo.GetProjectById(projectId)
}

// GetProjectsByOrganization busca projetos por organização
func (h *ProjectHandler) GetProjectsByOrganization(orgId string) ([]models.Project, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	return h.projectRepo.GetProjectByOrganization(orgUUID)
}

// CreateProject cria novo projeto e configurações padrão
func (h *ProjectHandler) CreateProject(project *models.Project) error {
	project.Id = uuid.New()
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()

	// Cria o projeto
	err := h.projectRepo.CreateProject(project)
	if err != nil {
		return err
	}

	// Cria configurações padrão para o projeto
	_, err = h.settingsRepo.GetOrCreateSettings(project.OrganizationId, project.Id)
	if err != nil {
		return err
	}

	// Cria templates padrão de notificação
	defaultTemplates := utils.CreateDefaultNotificationTemplates(project.OrganizationId, project.Id)
	for _, template := range defaultTemplates {
		if err := h.notificationRepo.CreateNotificationTemplate(&template); err != nil {
			// Log do erro mas não interrompe o processo
			continue
		}
	}

	// Cria configurações padrão de notificação
	defaultConfigs := utils.CreateDefaultNotificationConfigs(project.OrganizationId, project.Id)
	for _, config := range defaultConfigs {
		if err := h.notificationRepo.CreateOrUpdateNotificationConfig(&config); err != nil {
			// Log do erro mas não interrompe o processo
			continue
		}
	}

	return nil
}

// UpdateProject atualiza projeto existente
func (h *ProjectHandler) UpdateProject(project *models.Project) error {
	project.UpdatedAt = time.Now()
	return h.projectRepo.UpdateProject(project)
}

// SoftDeleteProject remove projeto logicamente
func (h *ProjectHandler) SoftDeleteProject(id string) error {
	projectId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return h.projectRepo.SoftDeleteProject(projectId)
}

// GetActiveProjects busca projetos ativos por organização
func (h *ProjectHandler) GetActiveProjects(orgId string) ([]models.Project, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	return h.projectRepo.GetActiveProjects(orgUUID)
}