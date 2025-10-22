package models

import (
	"time"

	"github.com/google/uuid"
)

// --- Environment (ambientes do restaurante) ---
type Environment struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	Name           string     `json:"name"` // ex: "Salão Principal", "Varanda"
	Description    string     `json:"description,omitempty"`
	Capacity       int        `json:"capacity"` // capacidade máxima do ambiente
	Active         bool       `json:"active" gorm:"default:true"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}
