package handler

import (
	"time"

	"github.com/google/uuid"
)

// DTOs enriquecidos para resposta de login

type UserOrganizationWithName struct {
	Id             uuid.UUID  `json:"id"`
	UserId         uuid.UUID  `json:"user_id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	OrganizationName string   `json:"organization_name"` // NOVO
	Role           string     `json:"role"`
	Active         bool       `json:"active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type UserProjectWithName struct {
	Id               uuid.UUID  `json:"id"`
	UserId           uuid.UUID  `json:"user_id"`
	ProjectId        uuid.UUID  `json:"project_id"`
	ProjectName      string     `json:"project_name"`      // Nome do projeto
	OrganizationId   uuid.UUID  `json:"organization_id"`   // ID da organização pai
	OrganizationName string     `json:"organization_name"` // Nome da organização pai
	Role             string     `json:"role"`
	Active           bool       `json:"active"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

// UserAdminRoleInfo representa um cargo admin atribuído ao usuário
// Usado para informar o frontend sobre cargos de escopo "admin" do usuário
type UserAdminRoleInfo struct {
	Id              uuid.UUID `json:"id"`               // ID do user_role
	RoleId          uuid.UUID `json:"role_id"`          // ID do cargo
	RoleName        string    `json:"role_name"`        // Nome técnico (ex: "admin_support")
	RoleDisplayName string    `json:"role_display_name"` // Nome amigável (ex: "Suporte Técnico")
	Scope           string    `json:"scope"`            // Sempre "admin" para este DTO
	HierarchyLevel  int       `json:"hierarchy_level"`  // Nível de hierarquia (1-10)
	Active          bool      `json:"active"`           // Se está ativo
}

// ========== NOVOS DTOs para Admin e Client separados ==========

// AdminLoginRequest dados de login para admin
type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AdminLoginResponse resposta de login para admin
type AdminLoginResponse struct {
	Admin       interface{} `json:"admin"`       // Dados do admin (sem senha)
	Token       string      `json:"token"`
	UserType    string      `json:"user_type"`   // Sempre "admin"
	Permissions []string    `json:"permissions"`
}

// ClientLoginRequest dados de login para client
type ClientLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	OrgSlug  string `json:"org_slug" binding:"required"` // Slug da organização
}

// ClientLoginResponse resposta de login para client
type ClientLoginResponse struct {
	Client           interface{}         `json:"client"`            // Dados do cliente (sem senha)
	Token            string              `json:"token"`
	UserType         string              `json:"user_type"`         // Sempre "client"
	Organization     OrganizationInfo    `json:"organization"`      // Info da organização
	Projects         []ClientProjectInfo `json:"projects"`          // Projetos do cliente
	Permissions      []string            `json:"permissions"`
}

// OrganizationInfo informações básicas da organização
type OrganizationInfo struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}

// TenantResolveResponse resposta para resolver tenant
type TenantResolveResponse struct {
	Type             string `json:"type"`              // "admin" ou "client"
	OrganizationId   string `json:"organization_id,omitempty"`
	OrganizationName string `json:"organization_name,omitempty"`
	OrganizationSlug string `json:"organization_slug,omitempty"`
	LoginUrl         string `json:"login_url"`
}
