package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Admin representa um usuário administrador do sistema
// Administradores têm acesso a todas as organizações e projetos
type Admin struct {
	Id           uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name         string         `json:"name" gorm:"not null"`
	Email        string         `gorm:"unique;not null" json:"email"`
	Password     string         `json:"password,omitempty" gorm:"not null"` // hash bcrypt
	Permissions  pq.StringArray `gorm:"type:text[]" json:"permissions"`     // ex: ["master_admin"]
	Active       bool           `gorm:"default:true" json:"active"`
	LastAccessAt *time.Time     `json:"last_access_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    *time.Time     `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName define o nome da tabela no banco de dados
func (Admin) TableName() string {
	return "admins"
}

// IsActive verifica se o admin está ativo e não foi deletado
func (a *Admin) IsActive() bool {
	return a.Active && a.DeletedAt == nil
}

// HasPermission verifica se o admin tem uma permissão específica
func (a *Admin) HasPermission(permission string) bool {
	for _, p := range a.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// IsMasterAdmin verifica se é um master admin
func (a *Admin) IsMasterAdmin() bool {
	return a.HasPermission("master_admin")
}
