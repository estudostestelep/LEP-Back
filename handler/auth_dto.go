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
