package models

import (
	"time"

	"github.com/google/uuid"
)

// --- Reservation (reserva de mesa) ---
type Reservation struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	CustomerId     uuid.UUID  `json:"customer_id"`
	TableId        *uuid.UUID `json:"table_id,omitempty"`
	Datetime       string     `json:"datetime"`
	PartySize      int        `json:"party_size"`
	Note           string     `json:"note,omitempty"`
	Status         string     `json:"status"` // "confirmed", "cancelled", "completed", "no_show", "pending", "not_approved"
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}
