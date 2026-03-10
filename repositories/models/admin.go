package models

import (
	"time"

	"github.com/google/uuid"
)

// Admin representa um usuário administrador do sistema
// Administradores têm acesso ao painel administrativo e suas permissões são definidas via Roles
type Admin struct {
	Id           uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name         string     `json:"name" gorm:"not null"`
	Email        string     `gorm:"unique;not null" json:"email"`
	Password     string     `json:"-" gorm:"not null"` // hash bcrypt, nunca serializado
	Active       bool       `gorm:"default:true" json:"active"`
	LastAccessAt *time.Time `json:"last_access_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relacionamentos - permissões são via roles
	Roles []Role `json:"roles,omitempty" gorm:"many2many:admin_roles;joinForeignKey:AdminId;joinReferences:RoleId"`
}

// TableName define o nome da tabela no banco de dados
func (Admin) TableName() string {
	return "admins"
}

// IsActive verifica se o admin está ativo e não foi deletado
func (a *Admin) IsActive() bool {
	return a.Active && a.DeletedAt == nil
}

// IsMasterAdmin verifica se é um master admin baseado nos roles
// Requer que Roles esteja preloaded
func (a *Admin) IsMasterAdmin() bool {
	for _, role := range a.Roles {
		if role.IsMasterAdmin() {
			return true
		}
	}
	return false
}

// GetMaxHierarchyLevel retorna o maior nível de hierarquia entre os roles
// Requer que Roles esteja preloaded
func (a *Admin) GetMaxHierarchyLevel() int {
	maxLevel := 0
	for _, role := range a.Roles {
		if role.HierarchyLevel > maxLevel {
			maxLevel = role.HierarchyLevel
		}
	}
	return maxLevel
}

// HasRole verifica se o admin tem um role específico
// Requer que Roles esteja preloaded
func (a *Admin) HasRole(roleName string) bool {
	for _, role := range a.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}
