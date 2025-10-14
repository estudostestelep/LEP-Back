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
