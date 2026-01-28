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
