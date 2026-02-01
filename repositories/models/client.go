package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Client representa um usuário cliente do sistema
// Clientes pertencem a UMA organização e podem ter acesso a múltiplos projetos
// Suas permissões são definidas via Roles na tabela client_roles
type Client struct {
	Id           uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name         string         `json:"name" gorm:"not null"`
	Email        string         `gorm:"unique;not null" json:"email"`
	Password     string         `json:"-" gorm:"not null"` // hash bcrypt, nunca serializado
	OrgId        uuid.UUID      `json:"org_id" gorm:"type:uuid;not null;index"`
	ProjIds      pq.StringArray `json:"proj_ids" gorm:"type:text[]"` // Array de UUIDs dos projetos
	Active       bool           `gorm:"default:true" json:"active"`
	LastAccessAt *time.Time     `json:"last_access_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    *time.Time     `json:"deleted_at,omitempty" gorm:"index"`

	// Relacionamentos - permissões são via roles
	Roles []Role `json:"roles,omitempty" gorm:"many2many:client_roles;joinForeignKey:ClientId;joinReferences:RoleId"`
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

// GetMaxHierarchyLevel retorna o maior nível de hierarquia entre os roles
// Requer que Roles esteja preloaded
func (c *Client) GetMaxHierarchyLevel() int {
	maxLevel := 0
	for _, role := range c.Roles {
		if role.HierarchyLevel > maxLevel {
			maxLevel = role.HierarchyLevel
		}
	}
	return maxLevel
}

// HasRole verifica se o cliente tem um role específico
// Requer que Roles esteja preloaded
func (c *Client) HasRole(roleName string) bool {
	for _, role := range c.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}

// ==================== DTOs ====================

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
