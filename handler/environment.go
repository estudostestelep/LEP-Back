package handler

import (
	"errors"
	"fmt"
	"lep/repositories"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
)

type EnvironmentHandler struct {
	environmentRepo repositories.IEnvironmentRepository
}

type IEnvironmentHandler interface {
	GetEnvironmentById(id string) (*models.Environment, error)
	GetEnvironmentsByProject(orgId, projectId string) ([]models.Environment, error)
	CreateEnvironment(environment *models.Environment) error
	UpdateEnvironment(environment *models.Environment) error
	SoftDeleteEnvironment(id string) error
	GetActiveEnvironments(orgId, projectId string) ([]models.Environment, error)
}

func NewEnvironmentHandler(environmentRepo repositories.IEnvironmentRepository) IEnvironmentHandler {
	return &EnvironmentHandler{environmentRepo: environmentRepo}
}

// GetEnvironmentById busca ambiente por ID
func (h *EnvironmentHandler) GetEnvironmentById(id string) (*models.Environment, error) {
	envId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return h.environmentRepo.GetEnvironmentById(envId)
}

// GetEnvironmentsByProject busca ambientes por projeto
func (h *EnvironmentHandler) GetEnvironmentsByProject(orgId, projectId string) ([]models.Environment, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}

	projectUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}

	return h.environmentRepo.GetEnvironmentsByProject(orgUUID, projectUUID)
}

// CreateEnvironment cria novo ambiente
func (h *EnvironmentHandler) CreateEnvironment(environment *models.Environment) error {
	// Verificar se já existe ambiente com o mesmo nome no projeto
	exists, err := h.environmentRepo.CheckEnvironmentNameExists(environment.OrganizationId, environment.ProjectId, environment.Name, nil)
	if err != nil {
		return fmt.Errorf("erro ao verificar duplicata: %w", err)
	}
	if exists {
		return errors.New("already_exists: environment with this name already exists in this project")
	}

	environment.Id = uuid.New()
	environment.CreatedAt = time.Now()
	environment.UpdatedAt = time.Now()
	return h.environmentRepo.CreateEnvironment(environment)
}

// UpdateEnvironment atualiza ambiente existente
func (h *EnvironmentHandler) UpdateEnvironment(environment *models.Environment) error {
	// Verificar se já existe outro ambiente com o mesmo nome no projeto
	exists, err := h.environmentRepo.CheckEnvironmentNameExists(environment.OrganizationId, environment.ProjectId, environment.Name, &environment.Id)
	if err != nil {
		return fmt.Errorf("erro ao verificar duplicata: %w", err)
	}
	if exists {
		return errors.New("already_exists: environment with this name already exists in this project")
	}

	environment.UpdatedAt = time.Now()
	return h.environmentRepo.UpdateEnvironment(environment)
}

// SoftDeleteEnvironment remove ambiente logicamente
func (h *EnvironmentHandler) SoftDeleteEnvironment(id string) error {
	envId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return h.environmentRepo.SoftDeleteEnvironment(envId)
}

// GetActiveEnvironments busca ambientes ativos por projeto
func (h *EnvironmentHandler) GetActiveEnvironments(orgId, projectId string) ([]models.Environment, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}

	projectUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}

	return h.environmentRepo.GetActiveEnvironments(orgUUID, projectUUID)
}