package models

import (
	"time"

	"github.com/google/uuid"
)

// --- Table (mesa) ---
type Table struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	EnvironmentId  *uuid.UUID `json:"environment_id,omitempty"` // vinculação com ambiente
	Number         int        `json:"number"`
	Capacity       int        `json:"capacity"`
	Location       string     `json:"location,omitempty"`            // descrição adicional da localização
	Status         string     `json:"status" gorm:"default:'livre'"` // "livre", "ocupada", "reservada"
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}
