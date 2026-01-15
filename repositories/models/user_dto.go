package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// UserWithRelations DTO para retornar usuário com suas organizações e projetos
type UserWithRelations struct {
	Id            uuid.UUID      `json:"id"`
	Name          string         `json:"name"`
	Email         string         `json:"email"`
	Permissions   pq.StringArray `json:"permissions"`
	Active        bool           `json:"active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     *time.Time     `json:"deleted_at,omitempty"`
	Organizations []UserOrgInfo  `json:"organizations"`
	Projects      []UserProjInfo `json:"projects"`
}

// UserOrgInfo informações do vínculo do usuário com organização
type UserOrgInfo struct {
	OrganizationId uuid.UUID `json:"organization_id"`
	Role           string    `json:"role"`
	Active         bool      `json:"active"`
}

// UserProjInfo informações do vínculo do usuário com projeto
type UserProjInfo struct {
	ProjectId uuid.UUID `json:"project_id"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
}

// CreateUserRequest DTO para criar usuário (aceita tanto role quanto permissions)
type CreateUserRequest struct {
	Id          uuid.UUID `json:"id,omitempty"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Role        string    `json:"role,omitempty"`        // Campo usado pelo seed (manager, waiter, kitchen, admin)
	RoleId      string    `json:"role_id,omitempty"`     // ID do cargo a ser atribuído (sistema novo de permissões)
	Permissions []string  `json:"permissions,omitempty"` // Permissões diretas
	Active      bool      `json:"active"`
}

// ToUser converte CreateUserRequest para User, mapeando role para permissions se necessário
func (r *CreateUserRequest) ToUser() *User {
	user := &User{
		Id:       r.Id,
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
		Active:   r.Active,
	}

	// Se permissions foi fornecido, usar diretamente
	if len(r.Permissions) > 0 {
		user.Permissions = r.Permissions
	} else if r.Role != "" {
		// Converter role para permissions
		user.Permissions = RoleToPermissions(r.Role)
	}

	return user
}

// RoleToPermissions converte um role string em uma lista de permissions
func RoleToPermissions(role string) []string {
	switch role {
	case "admin":
		return []string{"admin", "manager", "waiter", "kitchen"}
	case "manager":
		return []string{"manager", "waiter"}
	case "waiter":
		return []string{"waiter"}
	case "kitchen":
		return []string{"kitchen"}
	case "owner":
		return []string{"admin", "manager", "waiter", "kitchen", "owner"}
	default:
		return []string{role}
	}
}
