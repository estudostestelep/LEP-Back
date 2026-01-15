package handler

import (
	"errors"
	"fmt"
	"lep/repositories"
	"lep/repositories/models"
	"strings"
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
	// Gerar ID apenas se não foi fornecido
	if organization.Id == uuid.Nil {
		organization.Id = uuid.New()
	}

	err := r.repo.Organizations.CreateOrganization(organization)
	if err != nil {
		return err
	}

	// 🎯 REGRA DE NEGÓCIO: Atribuir plano gratuito automaticamente
	// Toda organização começa com o plano "free"
	if err := r.assignFreePackage(organization.Id); err != nil {
		fmt.Printf("⚠️ Aviso: erro ao atribuir plano gratuito: %v\n", err)
		// Não falha a criação, apenas registra o aviso
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

	// Usar cascade delete para deletar todos os dados relacionados
	err = r.repo.CascadeDelete.SoftDeleteOrganizationCascade(organizationId)
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

	// Usar cascade delete para deletar permanentemente todos os dados relacionados
	err = r.repo.CascadeDelete.HardDeleteOrganizationCascade(organizationId)
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

	// Normalizar nome para gerar email válido (sem espaços, lowercase)
	normalizedName := strings.ToLower(strings.ReplaceAll(name, " ", "-"))

	// Verificar se já existe organização com esse nome
	existingOrg, _ := r.repo.Organizations.GetOrganizationByEmail(fmt.Sprintf("%s@lep.com", normalizedName))
	if existingOrg != nil {
		return nil, errors.New("já existe uma organização com esse nome")
	}

	// Criar organização
	org := &models.Organization{
		Id:          uuid.New(),
		Name:        name,
		Email:       fmt.Sprintf("%s@lep.com", normalizedName),
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
		Id:          uuid.New(),
		Name:        name,
		Email:       fmt.Sprintf("%s@lep.com", normalizedName),
		Password:    string(hashedPassword),
		Permissions: []string{"admin"},
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = r.repo.User.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar usuário: %v", err)
	}

	// Criar relacionamento usuário-organização
	userOrg := &models.UserOrganization{
		Id:             uuid.New(),
		UserId:         user.Id,
		OrganizationId: org.Id,
		Role:           "owner",
		Active:         true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = r.repo.UserOrganizations.Create(userOrg)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar relacionamento usuário-organização: %v", err)
	}

	// Criar relacionamento usuário-projeto
	userProj := &models.UserProject{
		Id:        uuid.New(),
		UserId:    user.Id,
		ProjectId: project.Id,
		Role:      "admin",
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = r.repo.UserProjects.Create(userProj)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar relacionamento usuário-projeto: %v", err)
	}

	// 🔑 REGRA DE NEGÓCIO: Adicionar master admins automaticamente à nova org
	// Master admins têm acesso a todas as organizações
	if err := r.addMasterAdminsToOrganization(org.Id, project.Id); err != nil {
		return nil, fmt.Errorf("erro ao adicionar master admins: %v", err)
	}

	// 🎯 REGRA DE NEGÓCIO: Atribuir plano gratuito automaticamente
	// Toda organização começa com o plano "free"
	if err := r.assignFreePackage(org.Id); err != nil {
		fmt.Printf("⚠️ Aviso: erro ao atribuir plano gratuito: %v\n", err)
		// Não falha a criação, apenas registra o aviso
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

// addMasterAdminsToOrganization adiciona automaticamente todos os master admins
// à nova organização criada (REGRA DE NEGÓCIO)
func (r *resourceOrganization) addMasterAdminsToOrganization(
	organizationId, projectId uuid.UUID,
) error {
	// Master admin IDs (definidos no seed_data.go)
	masterAdminEmails := []string{
		"pablo@lep.com",
		"luan@lep.com",
		"eduardo@lep.com",
	}

	now := time.Now()

	// Buscar cada master admin e criar relacionamentos
	for _, email := range masterAdminEmails {
		// Buscar usuário por email
		user, err := r.repo.User.GetUserByEmail(email)
		if err != nil {
			// Se não encontrar, pular (pode não estar cadastrado ainda)
			continue
		}

		// Criar relacionamento usuário-organização (se não existir)
		userOrg := &models.UserOrganization{
			Id:             uuid.New(),
			UserId:         user.Id,
			OrganizationId: organizationId,
			Role:           "admin", // Master admins são admins da org
			Active:         true,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		// Ignorar erro se já existe (idempotente)
		_ = r.repo.UserOrganizations.Create(userOrg)

		// Criar relacionamento usuário-projeto (se não existir)
		userProj := &models.UserProject{
			Id:        uuid.New(),
			UserId:    user.Id,
			ProjectId: projectId,
			Role:      "admin",
			Active:    true,
			CreatedAt: now,
			UpdatedAt: now,
		}

		// Ignorar erro se já existe (idempotente)
		_ = r.repo.UserProjects.Create(userProj)
	}

	return nil
}

// assignFreePackage atribui o plano gratuito (free) para uma nova organização
// Esta função é chamada automaticamente no bootstrap de organizações
func (r *resourceOrganization) assignFreePackage(orgId uuid.UUID) error {
	// Buscar o pacote gratuito pelo código
	freePkg, err := r.repo.Packages.GetByCodeName("free")
	if err != nil {
		return fmt.Errorf("plano gratuito não encontrado: %w", err)
	}

	// Criar a assinatura da organização
	now := time.Now()
	orgPackage := &models.OrganizationPackage{
		Id:             uuid.New(),
		OrganizationId: orgId,
		PackageId:      freePkg.Id,
		BillingCycle:   "monthly",
		Active:         true,
		StartedAt:      &now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	return r.repo.Packages.SubscribeOrganization(orgPackage)
}

func NewSourceHandlerOrganization(repo *repositories.DBconn) IHandlerOrganization {
	return &resourceOrganization{repo: repo}
}
