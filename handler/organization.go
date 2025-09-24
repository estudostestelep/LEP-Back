package handler

import (
	"errors"
	"fmt"
	"lep/repositories"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type resourceOrganization struct {
	repo *repositories.DBconn
}

type OrganizationBootstrapResponse struct {
	Organization struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"organization"`
	Project struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"project"`
	User struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	} `json:"user"`
	Message string `json:"message"`
}

type IHandlerOrganization interface {
	GetOrganizationById(id string) (*models.Organization, error)
	GetOrganizationByEmail(email string) (*models.Organization, error)
	ListOrganizations() ([]models.Organization, error)
	ListActiveOrganizations() ([]models.Organization, error)
	CreateOrganization(organization *models.Organization) error
	UpdateOrganization(organization *models.Organization) error
	SoftDeleteOrganization(id string) error
	HardDeleteOrganization(id string) error
	CreateOrganizationBootstrap(name, password string) (*OrganizationBootstrapResponse, error)
}

func (r *resourceOrganization) GetOrganizationById(id string) (*models.Organization, error) {
	// Validar UUID
	organizationId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	resp, err := r.repo.Organizations.GetOrganizationById(organizationId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceOrganization) GetOrganizationByEmail(email string) (*models.Organization, error) {
	resp, err := r.repo.Organizations.GetOrganizationByEmail(email)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceOrganization) ListOrganizations() ([]models.Organization, error) {
	resp, err := r.repo.Organizations.ListOrganizations()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceOrganization) ListActiveOrganizations() ([]models.Organization, error) {
	resp, err := r.repo.Organizations.ListActiveOrganizations()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceOrganization) CreateOrganization(organization *models.Organization) error {
	organization.Id = uuid.New()
	err := r.repo.Organizations.CreateOrganization(organization)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceOrganization) UpdateOrganization(organization *models.Organization) error {
	err := r.repo.Organizations.UpdateOrganization(organization)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceOrganization) SoftDeleteOrganization(id string) error {
	// Validar UUID
	organizationId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	err = r.repo.Organizations.SoftDeleteOrganization(organizationId)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceOrganization) HardDeleteOrganization(id string) error {
	// Validar UUID
	organizationId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	err = r.repo.Organizations.HardDeleteOrganization(organizationId)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceOrganization) CreateOrganizationBootstrap(name, password string) (*OrganizationBootstrapResponse, error) {
	// Validar senha hard-coded
	if password != "senha123" {
		return nil, errors.New("senha inválida")
	}

	// Validar nome
	if name == "" {
		return nil, errors.New("nome da organização é obrigatório")
	}

	// Verificar se já existe organização com esse nome
	existingOrg, _ := r.repo.Organizations.GetOrganizationByEmail(fmt.Sprintf("%s@lep.com", name))
	if existingOrg != nil {
		return nil, errors.New("já existe uma organização com esse nome")
	}

	// Criar organização
	org := &models.Organization{
		Id:          uuid.New(),
		Name:        name,
		Email:       fmt.Sprintf("%s@lep.com", name),
		Description: fmt.Sprintf("Organização %s", name),
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := r.repo.Organizations.CreateOrganization(org)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar organização: %v", err)
	}

	// Criar projeto padrão
	project := &models.Project{
		Id:             uuid.New(),
		OrganizationId: org.Id,
		Name:           fmt.Sprintf("Projeto %s", name),
		Description:    fmt.Sprintf("Projeto padrão da organização %s", name),
		Active:         true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = r.repo.Projects.CreateProject(project)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar projeto: %v", err)
	}

	// Criar usuário admin
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar hash da senha: %v", err)
	}

	user := &models.User{
		Id:             uuid.New(),
		OrganizationId: org.Id,
		ProjectId:      project.Id,
		Name:           name,
		Email:          fmt.Sprintf("%s@lep.com", name),
		Password:       string(hashedPassword),
		Role:           "admin",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = r.repo.User.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar usuário: %v", err)
	}

	// Montar resposta
	response := &OrganizationBootstrapResponse{
		Organization: struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		}{
			ID:    org.Id.String(),
			Name:  org.Name,
			Email: org.Email,
		},
		Project: struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}{
			ID:   project.Id.String(),
			Name: project.Name,
		},
		User: struct {
			ID       string `json:"id"`
			Email    string `json:"email"`
			Name     string `json:"name"`
			Password string `json:"password"`
		}{
			ID:       user.Id.String(),
			Email:    user.Email,
			Name:     user.Name,
			Password: "senha123", // Retorna senha em texto claro para login
		},
		Message: "Organização criada com sucesso! Você pode fazer login com as credenciais fornecidas.",
	}

	return response, nil
}

func NewSourceHandlerOrganization(repo *repositories.DBconn) IHandlerOrganization {
	return &resourceOrganization{repo: repo}
}
