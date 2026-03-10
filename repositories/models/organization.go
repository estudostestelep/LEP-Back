package models

import (
	"time"

	"github.com/google/uuid"
)

// --- Organization (organização mãe) ---
type Organization struct {
	Id          uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string     `json:"name" gorm:"not null"`
	Slug        string     `json:"slug" gorm:"unique;size:100"` // Identificador único para subdomínio
	Email       string     `gorm:"unique" json:"email"`
	Phone       string     `json:"phone,omitempty"`
	Address     string     `json:"address,omitempty"`
	Website     string     `json:"website,omitempty"`
	Description string     `json:"description,omitempty"`
	Active      bool       `gorm:"default:true" json:"active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
