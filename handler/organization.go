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
	"gorm.io/gorm"
)

type resourceOrganization struct {
	repo *repositories.DBconn
	db   *gorm.DB
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
	GetOrganizationBySlug(slug string) (*models.Organization, error)
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

func (r *resourceOrganization) GetOrganizationBySlug(slug string) (*models.Organization, error) {
	resp, err := r.repo.Organizations.GetOrganizationBySlug(slug)
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

	// Gerar slug automaticamente se não fornecido
	if organization.Slug == "" {
		organization.Slug = strings.ToLower(strings.ReplaceAll(organization.Name, " ", "-"))
	}

	err := r.repo.Organizations.CreateOrganization(organization)
	if err != nil {
		return err
	}

	// 🎯 REGRA DE NEGÓCIO: Atribuir plano automaticamente
	// Organização "demo" recebe plano demo, demais recebem plano free
	if organization.Slug == "demo" {
		if err := r.assignPlanByCode(organization.Id, "demo"); err != nil {
			fmt.Printf("⚠️ Aviso: erro ao atribuir plano demo: %v\n", err)
		}
	} else {
		if err := r.assignFreePlan(organization.Id); err != nil {
			fmt.Printf("⚠️ Aviso: erro ao atribuir plano gratuito: %v\n", err)
		}
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
	slug := normalizedName
	email := fmt.Sprintf("%s@lep.com", normalizedName)

	// Verificar se já existe organização com esse nome (ativa)
	existingOrg, _ := r.repo.Organizations.GetOrganizationByEmail(email)
	if existingOrg != nil {
		return nil, errors.New("já existe uma organização com esse nome")
	}

	// Remover organizações soft-deleted que bloqueiam slug/email (para permitir recriação)
	r.removeSoftDeletedOrgBySlugOrEmail(slug, email)

	// Buscar o role de "org_owner" ANTES da transação (opcional - não bloqueia criação)
	var ownerRole *models.Role
	ownerRole, _ = r.repo.Roles.GetByName("org_owner")
	if ownerRole == nil {
		// Se não encontrar org_owner, tentar org_admin
		ownerRole, _ = r.repo.Roles.GetByName("org_admin")
	}

	// Gerar hash da senha ANTES da transação
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar hash da senha: %v", err)
	}

	// Preparar entidades
	now := time.Now()
	org := &models.Organization{
		Id:          uuid.New(),
		Name:        name,
		Slug:        slug,
		Email:       email,
		Description: fmt.Sprintf("Organização %s", name),
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	project := &models.Project{
		Id:             uuid.New(),
		OrganizationId: org.Id,
		Name:           fmt.Sprintf("Projeto %s", name),
		Description:    fmt.Sprintf("Projeto padrão da organização %s", name),
		Active:         true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	client := &models.Client{
		Id:       uuid.New(),
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		OrgId:    org.Id,
		Active:   true,
	}

	// Executar tudo em uma transação atômica
	err = r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Criar organização
		if err := tx.Create(org).Error; err != nil {
			return fmt.Errorf("erro ao criar organização: %v", err)
		}

		// 2. Criar projeto
		if err := tx.Create(project).Error; err != nil {
			return fmt.Errorf("erro ao criar projeto: %v", err)
		}

		// 3. Criar cliente
		if err := tx.Create(client).Error; err != nil {
			return fmt.Errorf("erro ao criar cliente: %v", err)
		}

		// 4. Atribuir role ao cliente (se existir role no sistema)
		if ownerRole != nil {
			clientRole := &models.ClientRole{
				Id:             uuid.New(),
				ClientId:       client.Id,
				RoleId:         ownerRole.Id,
				OrganizationId: org.Id,
				ProjectId:      &project.Id,
				Active:         true,
			}
			if err := tx.Create(clientRole).Error; err != nil {
				fmt.Printf("⚠️ Aviso: erro ao atribuir cargo ao cliente: %v\n", err)
				// Não falha a transação, apenas registra o aviso
			}
		} else {
			fmt.Printf("⚠️ Aviso: nenhum cargo org_owner ou org_admin encontrado no sistema\n")
		}

		// 5. Atribuir plano gratuito
		freePlan := &models.Plan{}
		if err := tx.Where("code = ?", "free").First(freePlan).Error; err != nil {
			fmt.Printf("⚠️ Aviso: plano gratuito não encontrado: %v\n", err)
			// Não falha a transação, apenas registra o aviso
		} else {
			orgPlan := &models.OrganizationPlan{
				Id:             uuid.New(),
				OrganizationId: org.Id,
				PlanId:         freePlan.Id,
				BillingCycle:   "monthly",
				Active:         true,
				StartsAt:       &now,
			}
			if err := tx.Create(orgPlan).Error; err != nil {
				fmt.Printf("⚠️ Aviso: erro ao atribuir plano gratuito: %v\n", err)
				// Não falha a transação, apenas registra o aviso
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Montar resposta (após transação bem-sucedida)
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
			ID:       client.Id.String(),
			Email:    client.Email,
			Name:     client.Name,
			Password: "senha123", // Retorna senha em texto claro para login
		},
		Message: "Organização criada com sucesso! Você pode fazer login com as credenciais fornecidas.",
	}

	return response, nil
}

// removeSoftDeletedOrgBySlugOrEmail remove organizações soft-deleted que bloqueiam slug ou email
func (r *resourceOrganization) removeSoftDeletedOrgBySlugOrEmail(slug, email string) {
	softDeletedOrgs, err := r.repo.Organizations.FindSoftDeletedBySlugOrEmail(slug, email)
	if err != nil || len(softDeletedOrgs) == 0 {
		return
	}
	for _, org := range softDeletedOrgs {
		fmt.Printf("🗑️ Removendo organização soft-deleted bloqueante: %s (slug=%s, email=%s)\n", org.Name, org.Slug, org.Email)
		_ = r.repo.CascadeDelete.HardDeleteOrganizationCascade(org.Id)
	}
}

// addMasterAdminsToOrganization adiciona automaticamente todos os master admins
// à nova organização criada (REGRA DE NEGÓCIO)
// NOTA: Master admins têm acesso global via permissão master_admin, não precisam de UserRole por org
func (r *resourceOrganization) addMasterAdminsToOrganization(
	organizationId, projectId uuid.UUID,
) error {
	// Master admins têm acesso global via a permissão "master_admin" no array de permissions do User
	// O middleware de autorização verifica essa permissão e concede acesso total
	// Não é mais necessário criar UserRole para cada organização
	fmt.Printf("ℹ️ Master admins têm acesso global via permissão master_admin\n")
	return nil
}

// assignFreePlan atribui o plano gratuito (free) para uma nova organização
// Esta função é chamada automaticamente no bootstrap de organizações
func (r *resourceOrganization) assignFreePlan(orgId uuid.UUID) error {
	// Buscar o plano gratuito pelo código
	freePlan, err := r.repo.Plans.GetByCode("free")
	if err != nil {
		return fmt.Errorf("plano gratuito não encontrado: %w", err)
	}

	// Criar a assinatura da organização
	now := time.Now()
	orgPlan := &models.OrganizationPlan{
		Id:             uuid.New(),
		OrganizationId: orgId,
		PlanId:         freePlan.Id,
		BillingCycle:   "monthly",
		Active:         true,
		StartsAt:       &now,
	}

	return r.repo.Plans.SubscribeOrganization(orgPlan)
}

// assignPlanByCode atribui um plano específico para uma organização pelo código do plano
func (r *resourceOrganization) assignPlanByCode(orgId uuid.UUID, planCode string) error {
	plan, err := r.repo.Plans.GetByCode(planCode)
	if err != nil {
		return fmt.Errorf("plano '%s' não encontrado: %w", planCode, err)
	}

	now := time.Now()
	orgPlan := &models.OrganizationPlan{
		Id:             uuid.New(),
		OrganizationId: orgId,
		PlanId:         plan.Id,
		BillingCycle:   "monthly",
		Active:         true,
		StartsAt:       &now,
	}

	return r.repo.Plans.SubscribeOrganization(orgPlan)
}

func NewSourceHandlerOrganization(repo *repositories.DBconn, db *gorm.DB) IHandlerOrganization {
	return &resourceOrganization{repo: repo, db: db}
}
