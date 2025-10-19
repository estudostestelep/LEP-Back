package models

import (
	"time"

	"github.com/google/uuid"
)

// --- Log/Audit (auditoria de ações) ---
type AuditLog struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	UserId         *uuid.UUID `json:"user_id,omitempty"`
	Action         string     `json:"action"` // ex: "create_reservation", "cancel_order"
	Entity         string     `json:"entity"` // ex: "reservation", "order"
	EntityId       uuid.UUID  `json:"entity_id"`
	Description    string     `json:"description,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}
