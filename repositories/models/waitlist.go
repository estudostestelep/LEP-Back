package models

import (
	"time"

	"github.com/google/uuid"
)

// --- Waitlist (fila de espera) ---
type Waitlist struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	CustomerId     *uuid.UUID `json:"customer_id,omitempty"`
	CustomerName   string     `json:"customer_name"`
	CustomerPhone  string     `json:"customer_phone,omitempty"`
	CustomerEmail  string     `json:"customer_email,omitempty"`
	Notes          string     `json:"notes,omitempty"`
	People         int        `json:"party_size"`
	Status         string     `json:"status"` // ex: "waiting", "notified", "seated"
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}
