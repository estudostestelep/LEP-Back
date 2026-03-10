package handler

import (
	"errors"
	"fmt"
	"lep/repositories"
	"lep/repositories/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type resourceClientUser struct {
	repo *repositories.DBconn
}

type IHandlerClientUser interface {
	GetClientById(id string) (*models.Client, error)
	GetClientByEmail(email string) (*models.Client, error)
	GetClientByEmailAndOrg(email string, orgId string) (*models.Client, error)
	ListClients() ([]models.Client, error)
	ListClientsByOrgId(orgId string) ([]models.Client, error)
	ListClientsByOrganization(orgId string) ([]models.Client, error)
	ListClientsByProject(orgId string, projectId string) ([]models.Client, error)
	CreateClient(client *models.Client) error
	UpdateClient(client *models.Client) error
	DeleteClient(id string) error
	UpdateLastAccess(clientId string) error
	ValidateClientCredentials(email, password, orgSlug string) (*models.Client, *models.Organization, error)
	AddProjectToClient(clientId string, projectId string) error
	RemoveProjectFromClient(clientId string, projectId string) error
	GetClientWithOrganization(id string) (*models.ClientWithOrganization, error)
	GetClientProjects(client *models.Client) ([]ClientProjectInfo, error)
}

// ClientProjectInfo contém informações básicas de um projeto do cliente
type ClientProjectInfo struct {
	ProjectId   uuid.UUID `json:"project_id"`
	ProjectName string    `json:"project_name"`
}

func (r *resourceClientUser) GetClientById(id string) (*models.Client, error) {
	return r.repo.Clients.GetClientById(id)
}

func (r *resourceClientUser) GetClientByEmail(email string) (*models.Client, error) {
	return r.repo.Clients.GetClientByEmail(email)
}

func (r *resourceClientUser) GetClientByEmailAndOrg(email string, orgId string) (*models.Client, error) {
	return r.repo.Clients.GetClientByEmailAndOrg(email, orgId)
}

func (r *resourceClientUser) ListClients() ([]models.Client, error) {
	return r.repo.Clients.ListClients()
}

func (r *resourceClientUser) ListClientsByOrgId(orgId string) ([]models.Client, error) {
	return r.repo.Clients.ListClientsByOrganization(orgId)
}

func (r *resourceClientUser) ListClientsByOrganization(orgId string) ([]models.Client, error) {
	return r.repo.Clients.ListClientsByOrganization(orgId)
}

func (r *resourceClientUser) ListClientsByProject(orgId string, projectId string) ([]models.Client, error) {
	return r.repo.Clients.ListClientsByProject(orgId, projectId)
}

func (r *resourceClientUser) CreateClient(client *models.Client) error {
	// Verificar se email já existe na organização
	exists, err := r.repo.Clients.ClientEmailExistsInOrg(client.Email, client.OrgId.String())
	if err != nil {
		return fmt.Errorf("erro ao verificar email: %v", err)
	}
	if exists {
		return errors.New("email já cadastrado nesta organização")
	}

	// Verificar se a organização existe
	org, err := r.repo.Organizations.GetOrganizationById(client.OrgId)
	if err != nil || org == nil {
		return errors.New("organização não encontrada")
	}

	// Gerar UUID se não fornecido
	if client.Id == uuid.Nil {
		client.Id = uuid.New()
	}

	// Hash da senha
	if client.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(client.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("erro ao gerar hash da senha: %v", err)
		}
		client.Password = string(hashedPassword)
	}

	return r.repo.Clients.CreateClient(client)
}

func (r *resourceClientUser) UpdateClient(client *models.Client) error {
	// Se a senha foi fornecida, fazer hash
	if client.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(client.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("erro ao gerar hash da senha: %v", err)
		}
		client.Password = string(hashedPassword)
	}

	return r.repo.Clients.UpdateClient(client)
}

func (r *resourceClientUser) DeleteClient(id string) error {
	return r.repo.Clients.SoftDeleteClient(id)
}

func (r *resourceClientUser) UpdateLastAccess(clientId string) error {
	return r.repo.Clients.UpdateLastAccess(clientId)
}

func (r *resourceClientUser) ValidateClientCredentials(email, password, orgSlug string) (*models.Client, *models.Organization, error) {
	// Buscar organização pelo slug
	org, err := r.repo.Organizations.GetOrganizationBySlug(orgSlug)
	if err != nil || org == nil {
		return nil, nil, errors.New("organização não encontrada")
	}

	// Buscar cliente pelo email e organização
	client, err := r.repo.Clients.GetClientByEmailAndOrg(email, org.Id.String())
	if err != nil {
		return nil, nil, errors.New("credenciais inválidas")
	}

	if client == nil || !client.IsActive() {
		return nil, nil, errors.New("credenciais inválidas")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(client.Password), []byte(password)); err != nil {
		return nil, nil, errors.New("credenciais inválidas")
	}

	return client, org, nil
}

func (r *resourceClientUser) AddProjectToClient(clientId string, projectId string) error {
	// Verificar se o cliente existe
	client, err := r.repo.Clients.GetClientById(clientId)
	if err != nil || client == nil {
		return errors.New("cliente não encontrado")
	}

	// Verificar se o projeto existe e pertence à organização do cliente
	project, err := r.repo.Projects.GetProjectById(uuid.MustParse(projectId))
	if err != nil || project == nil {
		return errors.New("projeto não encontrado")
	}

	if project.OrganizationId != client.OrgId {
		return errors.New("projeto não pertence à organização do cliente")
	}

	return r.repo.Clients.AddProjectToClient(clientId, projectId)
}

func (r *resourceClientUser) RemoveProjectFromClient(clientId string, projectId string) error {
	return r.repo.Clients.RemoveProjectFromClient(clientId, projectId)
}

func (r *resourceClientUser) GetClientWithOrganization(id string) (*models.ClientWithOrganization, error) {
	return r.repo.Clients.GetClientWithOrganization(id)
}

func (r *resourceClientUser) GetClientProjects(client *models.Client) ([]ClientProjectInfo, error) {
	if len(client.ProjIds) == 0 {
		return []ClientProjectInfo{}, nil
	}

	result := make([]ClientProjectInfo, 0, len(client.ProjIds))
	for _, projIdStr := range client.ProjIds {
		projId, err := uuid.Parse(projIdStr)
		if err != nil {
			continue
		}

		project, err := r.repo.Projects.GetProjectById(projId)
		if err != nil || project == nil {
			continue
		}

		result = append(result, ClientProjectInfo{
			ProjectId:   project.Id,
			ProjectName: project.Name,
		})
	}

	return result, nil
}

func NewClientUserHandler(repo *repositories.DBconn) IHandlerClientUser {
	return &resourceClientUser{repo: repo}
}
