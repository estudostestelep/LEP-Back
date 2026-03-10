package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// UserWithRelations DTO para retornar usuário com suas organizações e projetos
type UserWithRelations struct {
	Id            uuid.UUID      `json:"id"`
	Name          string         `json:"name"`
	Email         string         `json:"email"`
	Permissions   pq.StringArray `json:"permissions"`
	Active        bool           `json:"active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     *time.Time     `json:"deleted_at,omitempty"`
	Organizations []UserOrgInfo  `json:"organizations"`
	Projects      []UserProjInfo `json:"projects"`
}

// UserOrgInfo informações do vínculo do usuário com organização
type UserOrgInfo struct {
	OrganizationId uuid.UUID `json:"organization_id"`
	Role           string    `json:"role"`
	Active         bool      `json:"active"`
}

// UserProjInfo informações do vínculo do usuário com projeto
type UserProjInfo struct {
	ProjectId uuid.UUID `json:"project_id"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
}

// UserWithRole DTO para retornar usuário com dados do cargo principal
type UserWithRole struct {
	User
	RoleName           string `json:"role_name"`
	RoleDisplayName    string `json:"role_display_name"`
	RoleHierarchyLevel int    `json:"role_hierarchy_level"`
}

// CreateUserRequest DTO para criar usuário (aceita tanto role quanto permissions)
type CreateUserRequest struct {
	Id          uuid.UUID `json:"id,omitempty"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Role        string    `json:"role,omitempty"`        // Campo usado pelo seed (manager, waiter, kitchen, admin)
	RoleId      string    `json:"role_id,omitempty"`     // ID do cargo a ser atribuído (sistema novo de permissões)
	Permissions []string  `json:"permissions,omitempty"` // Permissões diretas
	Active      bool      `json:"active"`
}

// ToUser converte CreateUserRequest para User, mapeando role para permissions se necessário
func (r *CreateUserRequest) ToUser() *User {
	user := &User{
		Id:       r.Id,
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
		Active:   r.Active,
	}

	// Se permissions foi fornecido, usar diretamente
	if len(r.Permissions) > 0 {
		user.Permissions = r.Permissions
	} else if r.Role != "" {
		// Converter role para permissions
		user.Permissions = RoleToPermissions(r.Role)
	}

	return user
}

// RoleToPermissions converte um role string em uma lista de permissions válidas do sistema
func RoleToPermissions(role string) []string {
	switch role {
	case "admin":
		// Admin tem acesso total ao restaurante
		return []string{
			"manage_projects", "manage_products", "manage_menus",
			"manage_orders", "manage_customers", "manage_tables",
			"manage_reservations", "manage_waitlists", "view_reports",
			"export_data", "manage_settings", "manage_notifications",
			"manage_tags", "manage_categories",
		}
	case "manager":
		// Manager pode gerenciar pedidos, mesas, reservas e clientes
		return []string{
			"manage_orders", "manage_customers", "manage_tables",
			"manage_reservations", "manage_waitlists", "view_reports",
		}
	case "waiter":
		// Garçom pode ver e criar pedidos, ver mesas e clientes
		return []string{
			"view_orders", "manage_orders", "view_tables", "view_customers",
		}
	case "kitchen":
		// Cozinha pode ver e atualizar status de pedidos
		return []string{
			"view_orders", "manage_orders",
		}
	case "sommelier":
		// Sommelier pode gerenciar pedidos e ver cardápio
		return []string{
			"view_orders", "manage_orders", "view_menus", "view_products",
		}
	case "owner":
		// Owner tem acesso total + master_admin
		return []string{
			"master_admin", "manage_users", "manage_organizations",
			"manage_projects", "manage_products", "manage_menus",
			"manage_orders", "manage_customers", "manage_tables",
			"manage_reservations", "manage_waitlists", "view_reports",
			"export_data", "manage_settings", "manage_notifications",
			"manage_tags", "manage_categories",
		}
	default:
		// Para roles desconhecidos, retornar como permissão direta
		// (pode ser uma permissão válida como "view_orders")
		return []string{role}
	}
}
