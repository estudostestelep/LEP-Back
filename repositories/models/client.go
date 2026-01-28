package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Client representa um usuário cliente do sistema
// Clientes pertencem a UMA organização e podem ter acesso a múltiplos projetos
type Client struct {
	Id           uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name         string         `json:"name" gorm:"not null"`
	Email        string         `gorm:"unique;not null" json:"email"`
	Password     string         `json:"password,omitempty" gorm:"not null"` // hash bcrypt
	OrgId        uuid.UUID      `json:"org_id" gorm:"type:uuid;not null;index"`
	ProjIds      pq.StringArray `json:"proj_ids" gorm:"type:text[]"` // Array de UUIDs dos projetos
	Permissions  pq.StringArray `gorm:"type:text[]" json:"permissions"`
	Active       bool           `gorm:"default:true" json:"active"`
	LastAccessAt *time.Time     `json:"last_access_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    *time.Time     `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName define o nome da tabela no banco de dados
func (Client) TableName() string {
	return "clients"
}

// IsActive verifica se o cliente está ativo e não foi deletado
func (c *Client) IsActive() bool {
	return c.Active && c.DeletedAt == nil
}

// HasProjectAccess verifica se o cliente tem acesso a um projeto específico
func (c *Client) HasProjectAccess(projectId string) bool {
	for _, pid := range c.ProjIds {
		if pid == projectId {
			return true
		}
	}
	return false
}

// HasPermission verifica se o cliente tem uma permissão específica
func (c *Client) HasPermission(permission string) bool {
	for _, p := range c.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// AddProject adiciona um projeto ao cliente se ainda não existir
func (c *Client) AddProject(projectId string) {
	if !c.HasProjectAccess(projectId) {
		c.ProjIds = append(c.ProjIds, projectId)
	}
}

// RemoveProject remove um projeto do cliente
func (c *Client) RemoveProject(projectId string) {
	newProjIds := make([]string, 0, len(c.ProjIds))
	for _, pid := range c.ProjIds {
		if pid != projectId {
			newProjIds = append(newProjIds, pid)
		}
	}
	c.ProjIds = newProjIds
}

// ClientWithOrganization é um DTO que inclui informações da organização
type ClientWithOrganization struct {
	Client
	OrganizationName string `json:"organization_name"`
	OrganizationSlug string `json:"organization_slug"`
}

// ClientWithProjects é um DTO que inclui informações dos projetos
type ClientWithProjects struct {
	Client
	Projects []ClientProjectInfo `json:"projects"`
}

// ClientProjectInfo contém informações básicas de um projeto do cliente
type ClientProjectInfo struct {
	ProjectId   uuid.UUID `json:"project_id"`
	ProjectName string    `json:"project_name"`
}
