package handler

import (
	"fmt"
	"lep/repositories"
	"lep/repositories/models"
	"lep/utils"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type ProjectHandler struct {
	projectRepo       repositories.IProjectRepository
	settingsRepo      repositories.ISettingsRepository
	notificationRepo  repositories.INotificationRepository
	cascadeDeleteRepo repositories.ICascadeDeleteRepository
}

type IProjectHandler interface {
	GetProjectById(id string) (*models.Project, error)
	GetProjectsByOrganization(orgId string) ([]models.Project, error)
	CreateProject(project *models.Project) error
	UpdateProject(project *models.Project) error
	SoftDeleteProject(id string) error
	HardDeleteProject(id string) error
	GetActiveProjects(orgId string) ([]models.Project, error)
	GetProjectBySlug(orgId, slug string) (*models.Project, error)
	GetDefaultProject(orgId string) (*models.Project, error)
	ResolveProject(orgId, projectSlug string) (*models.Project, error)
	SetDefaultProject(orgId, projectId string) error
	GenerateSlug(name string) string
	ProjectSlugExists(orgId, slug string) (bool, error)
}

func NewProjectHandler(
	projectRepo repositories.IProjectRepository,
	settingsRepo repositories.ISettingsRepository,
	notificationRepo repositories.INotificationRepository,
	cascadeDeleteRepo repositories.ICascadeDeleteRepository,
) IProjectHandler {
	return &ProjectHandler{
		projectRepo:       projectRepo,
		settingsRepo:      settingsRepo,
		notificationRepo:  notificationRepo,
		cascadeDeleteRepo: cascadeDeleteRepo,
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
	// Gerar ID apenas se não foi fornecido
	if project.Id == uuid.Nil {
		project.Id = uuid.New()
	}

	// Auto-gerar slug se não fornecido
	if project.Slug == "" {
		project.Slug = h.GenerateSlug(project.Name)
	}

	// Garantir unicidade do slug adicionando sufixo se necessário
	baseSlug := project.Slug
	counter := 1
	for {
		exists, err := h.projectRepo.ProjectSlugExists(project.OrganizationId, project.Slug)
		if err != nil {
			break
		}
		if !exists {
			break
		}
		project.Slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
		if counter > 99 {
			project.Slug = fmt.Sprintf("%s-%s", baseSlug, uuid.New().String()[:8])
			break
		}
	}

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

// SoftDeleteProject remove projeto logicamente com cascade delete
func (h *ProjectHandler) SoftDeleteProject(id string) error {
	projectId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	// Usar cascade delete para deletar todos os dados relacionados
	return h.cascadeDeleteRepo.SoftDeleteProjectCascade(projectId)
}

// HardDeleteProject remove projeto permanentemente com cascade delete
func (h *ProjectHandler) HardDeleteProject(id string) error {
	projectId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	// Usar cascade delete para deletar permanentemente todos os dados relacionados
	return h.cascadeDeleteRepo.HardDeleteProjectCascade(projectId)
}

// GetActiveProjects busca projetos ativos por organização
func (h *ProjectHandler) GetActiveProjects(orgId string) ([]models.Project, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	return h.projectRepo.GetActiveProjects(orgUUID)
}

// GetProjectBySlug busca projeto por slug dentro de uma organização
func (h *ProjectHandler) GetProjectBySlug(orgId, slug string) (*models.Project, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	return h.projectRepo.GetProjectBySlug(orgUUID, slug)
}

// GetDefaultProject busca o projeto padrão de uma organização
func (h *ProjectHandler) GetDefaultProject(orgId string) (*models.Project, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	return h.projectRepo.GetDefaultProject(orgUUID)
}

// ResolveProject resolve projeto por slug ou retorna o default/primeiro ativo
func (h *ProjectHandler) ResolveProject(orgId, projectSlug string) (*models.Project, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}

	// Se slug fornecido, busca por slug
	if projectSlug != "" {
		project, err := h.projectRepo.GetProjectBySlug(orgUUID, projectSlug)
		if err == nil {
			return project, nil
		}
	}

	// Tenta buscar projeto default
	project, err := h.projectRepo.GetDefaultProject(orgUUID)
	if err == nil {
		return project, nil
	}

	// Fallback: retorna primeiro projeto ativo
	projects, err := h.projectRepo.GetActiveProjects(orgUUID)
	if err != nil {
		return nil, err
	}
	if len(projects) > 0 {
		return &projects[0], nil
	}

	return nil, err
}

// SetDefaultProject define um projeto como padrão
func (h *ProjectHandler) SetDefaultProject(orgId, projectId string) error {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return err
	}
	projectUUID, err := uuid.Parse(projectId)
	if err != nil {
		return err
	}
	return h.projectRepo.SetDefaultProject(orgUUID, projectUUID)
}

// GenerateSlug gera um slug a partir do nome
func (h *ProjectHandler) GenerateSlug(name string) string {
	// Remove acentos
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, name)

	// Converte para minúsculas
	result = strings.ToLower(result)

	// Substitui espaços e caracteres especiais por hífens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	result = reg.ReplaceAllString(result, "-")

	// Remove hífens do início e fim
	result = strings.Trim(result, "-")

	return result
}

// ProjectSlugExists verifica se um slug já existe na organização
func (h *ProjectHandler) ProjectSlugExists(orgId, slug string) (bool, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return false, err
	}
	return h.projectRepo.ProjectSlugExists(orgUUID, slug)
}