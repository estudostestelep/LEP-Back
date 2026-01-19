package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// --- User (funcionário/admin) ---
type User struct {
	Id           uuid.UUID      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string         `json:"name"`
	Email        string         `gorm:"unique" json:"email"`
	Password     string         `json:"password"` // armazenar hash!
	Permissions  pq.StringArray `gorm:"type:text[]" json:"permissions"`
	Active       bool           `gorm:"default:true" json:"active"`
	LastAccessAt *time.Time     `json:"last_access_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    *time.Time     `json:"deleted_at,omitempty"`
}

// --- UserOrganization (relacionamento usuário-organização) ---
type UserOrganization struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId         uuid.UUID  `json:"user_id" gorm:"not null"`
	OrganizationId uuid.UUID  `json:"organization_id" gorm:"not null"`
	Role           string     `json:"role"`                       // ex: "owner", "admin", "member"
	Active         bool       `json:"active" gorm:"default:true"` // permite desativar sem deletar
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// --- UserProject (relacionamento usuário-projeto) ---
type UserProject struct {
	Id        uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    uuid.UUID  `json:"user_id" gorm:"not null"`
	ProjectId uuid.UUID  `json:"project_id" gorm:"not null"`
	Role      string     `json:"role"`                       // ex: "manager", "waiter", "admin"
	Active    bool       `json:"active" gorm:"default:true"` // permite desativar sem deletar
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
