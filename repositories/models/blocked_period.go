package models

import (
	"time"

	"github.com/google/uuid"
)

// --- SPRINT 4: Validações Avançadas ---

// BlockedPeriod - Períodos bloqueados para reservas
type BlockedPeriod struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	Name           string     `json:"name"` // ex: "Manutenção", "Evento Privado"
	Description    string     `json:"description,omitempty"`
	StartDateTime  time.Time  `json:"start_datetime"`
	EndDateTime    time.Time  `json:"end_datetime"`
	RecurringType  string     `json:"recurring_type,omitempty"` // "none", "weekly", "monthly"
	Active         bool       `json:"active" gorm:"default:true"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}
