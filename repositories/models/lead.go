package models

import (
	"time"

	"github.com/google/uuid"
)

// --- SPRINT 5: Features Avançadas ---

// Lead - Sistema básico de CRM
type Lead struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	Name           string     `json:"name"`
	Email          string     `json:"email,omitempty"`
	Phone          string     `json:"phone,omitempty"`
	Source         string     `json:"source"` // "waitlist", "reservation", "walk_in"
	Status         string     `json:"status"` // "new", "contacted", "converted", "lost"
	Notes          string     `json:"notes,omitempty"`
	LastContact    *time.Time `json:"last_contact,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}
