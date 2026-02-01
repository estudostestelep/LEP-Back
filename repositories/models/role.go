package models

import (
	"time"

	"github.com/google/uuid"
)

// ==================== Role ====================

// Role representa um cargo no sistema
// Cargos são atribuídos a usuários e definem suas permissões via tabela pivot role_permissions
type Role struct {
	Id             uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	Name           string     `json:"name" gorm:"size:100;not null;uniqueIndex"` // Nome técnico (ex: "owner", "manager")
	DisplayName    string     `json:"display_name" gorm:"size:100"`              // Nome amigável (ex: "Proprietário", "Gerente")
	Description    string     `json:"description" gorm:"size:500"`               // Descrição do cargo
	HierarchyLevel int        `json:"hierarchy_level" gorm:"default:1"`          // Nível de hierarquia (1-10, 10=master_admin)
	Scope          string     `json:"scope" gorm:"size:20;default:'client'"`     // "admin" ou "client"
	OrganizationId *uuid.UUID `json:"organization_id,omitempty" gorm:"type:uuid;index"`
	IsSystem       bool       `json:"is_system" gorm:"default:false"` // Cargo do sistema (não pode ser deletado)
	Active         bool       `json:"active" gorm:"default:true"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relacionamentos
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
}

func (Role) TableName() string {
	return "roles"
}

// IsMasterAdmin verifica se o cargo é de master admin baseado na hierarquia
func (r *Role) IsMasterAdmin() bool {
	return r.HierarchyLevel >= 10
}

// ==================== RolePermission (Pivot) ====================

// RolePermission representa a associação entre cargo e permissão
type RolePermission struct {
	RoleId       uuid.UUID `json:"role_id" gorm:"type:uuid;primaryKey"`
	PermissionId uuid.UUID `json:"permission_id" gorm:"type:uuid;primaryKey"`
	CreatedAt    time.Time `json:"created_at"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

// ==================== AdminRole (Pivot Admin-Role) ====================

// AdminRole representa a associação entre admin e cargo
type AdminRole struct {
	Id             uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	AdminId        uuid.UUID  `json:"admin_id" gorm:"type:uuid;not null;index"`
	RoleId         uuid.UUID  `json:"role_id" gorm:"type:uuid;not null;index"`
	OrganizationId *uuid.UUID `json:"organization_id,omitempty" gorm:"type:uuid;index"` // NULL = cargo admin global
	Active         bool       `json:"active" gorm:"default:true"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relacionamentos para preload
	Admin *Admin `json:"admin,omitempty" gorm:"foreignKey:AdminId"`
	Role  *Role  `json:"role,omitempty" gorm:"foreignKey:RoleId"`
}

func (AdminRole) TableName() string {
	return "admin_roles"
}

// ==================== ClientRole (Pivot Client-Role) ====================

// ClientRole representa a associação entre client e cargo
type ClientRole struct {
	Id             uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	ClientId       uuid.UUID  `json:"client_id" gorm:"type:uuid;not null;index"`
	RoleId         uuid.UUID  `json:"role_id" gorm:"type:uuid;not null;index"`
	OrganizationId uuid.UUID  `json:"organization_id" gorm:"type:uuid;not null;index"`
	ProjectId      *uuid.UUID `json:"project_id,omitempty" gorm:"type:uuid;index"` // Contexto do projeto (opcional)
	Active         bool       `json:"active" gorm:"default:true"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relacionamentos para preload
	Client *Client `json:"client,omitempty" gorm:"foreignKey:ClientId"`
	Role   *Role   `json:"role,omitempty" gorm:"foreignKey:RoleId"`
}

func (ClientRole) TableName() string {
	return "client_roles"
}

// ==================== Permission ====================

// Permission representa uma permissão no sistema no formato module:action
type Permission struct {
	Id          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	Code        string     `json:"code" gorm:"size:100;not null;uniqueIndex"`  // Código completo (ex: "orders:read")
	Module      string     `json:"module" gorm:"size:50;not null;index"`       // Módulo (ex: "orders")
	Action      string     `json:"action" gorm:"size:20;not null"`             // Ação (ex: "read")
	DisplayName string     `json:"display_name" gorm:"size:100"`               // Nome amigável
	Description string     `json:"description" gorm:"size:500"`                // Descrição
	Active      bool       `json:"active" gorm:"default:true"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

func (Permission) TableName() string {
	return "permissions"
}

// ==================== Module ====================

// Module representa um módulo do sistema (agrupamento de permissões)
type Module struct {
	Id           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	Code         string     `json:"code" gorm:"size:100;not null;uniqueIndex"` // Código (ex: "orders")
	Name         string     `json:"name" gorm:"size:100"`                      // Nome amigável (ex: "Pedidos")
	Description  string     `json:"description" gorm:"size:500"`
	Icon         string     `json:"icon" gorm:"size:50"`                   // Ícone do módulo (lucide icon name)
	Scope        string     `json:"scope" gorm:"size:20;default:'client'"` // "admin" ou "client"
	DisplayOrder int        `json:"display_order" gorm:"default:0"`        // Ordem de exibição
	IsFree       bool       `json:"is_free" gorm:"default:false"`          // Módulo gratuito (disponível para todos)
	Active       bool       `json:"active" gorm:"default:true"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

func (Module) TableName() string {
	return "modules"
}

// ==================== Plan (anteriormente Package) ====================

// Plan representa um plano de assinatura
type Plan struct {
	Id           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	Code         string     `json:"code" gorm:"size:100;not null;uniqueIndex"` // Código do plano
	Name         string     `json:"name" gorm:"size:100"`                      // Nome amigável (ex: "Plano Básico")
	Description  string     `json:"description" gorm:"size:500"`
	PriceMonthly float64    `json:"price_monthly" gorm:"type:decimal(10,2);default:0"`
	PriceYearly  float64    `json:"price_yearly" gorm:"type:decimal(10,2);default:0"`
	IsPublic     bool       `json:"is_public" gorm:"default:true"`  // Visível publicamente
	DisplayOrder int        `json:"display_order" gorm:"default:0"` // Ordem de exibição
	Active       bool       `json:"active" gorm:"default:true"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relacionamentos
	Modules []Module    `json:"modules,omitempty" gorm:"many2many:plan_modules;"`
	Limits  []PlanLimit `json:"limits,omitempty" gorm:"foreignKey:PlanId"`
}

func (Plan) TableName() string {
	return "plans"
}

// ==================== PlanModule (Pivot Plan-Module) ====================

// PlanModule associa módulos a planos
type PlanModule struct {
	PlanId    uuid.UUID `json:"plan_id" gorm:"type:uuid;primaryKey"`
	ModuleId  uuid.UUID `json:"module_id" gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `json:"created_at"`
}

func (PlanModule) TableName() string {
	return "plan_modules"
}

// ==================== PlanLimit ====================

// PlanLimit define limites de recursos por plano
type PlanLimit struct {
	Id         uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	PlanId     uuid.UUID  `json:"plan_id" gorm:"type:uuid;not null;index"`
	LimitType  string     `json:"limit_type" gorm:"size:50;not null"` // "users", "tables", "products", "reservations_per_day"
	LimitValue int        `json:"limit_value" gorm:"default:-1"`      // -1 = ilimitado
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relacionamentos
	Plan *Plan `json:"plan,omitempty" gorm:"foreignKey:PlanId"`
}

func (PlanLimit) TableName() string {
	return "plan_limits"
}

// ==================== OrganizationPlan ====================

// OrganizationPlan associa planos a organizações (assinatura)
type OrganizationPlan struct {
	Id             uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	OrganizationId uuid.UUID  `json:"organization_id" gorm:"type:uuid;not null;uniqueIndex"`
	PlanId         uuid.UUID  `json:"plan_id" gorm:"type:uuid;not null;index"`
	BillingCycle   string     `json:"billing_cycle" gorm:"size:20;default:'monthly'"`   // "monthly", "yearly"
	CustomPrice    *float64   `json:"custom_price,omitempty" gorm:"type:decimal(10,2)"` // Preço customizado (override)
	StartsAt       *time.Time `json:"starts_at,omitempty"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
	Active         bool       `json:"active" gorm:"default:true"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relacionamentos
	Organization *Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationId"`
	Plan         *Plan         `json:"plan,omitempty" gorm:"foreignKey:PlanId"`
}

func (OrganizationPlan) TableName() string {
	return "organization_plans"
}

// ==================== DTOs e Views ====================

// RoleWithPermissions representa um cargo com suas permissões carregadas
type RoleWithPermissions struct {
	Role
	PermissionCodes []string `json:"permission_codes,omitempty" gorm:"-"`
}

// PlanWithModules representa um plano com seus módulos carregados
type PlanWithModules struct {
	Plan
	ModuleCodes []string `json:"module_codes,omitempty" gorm:"-"`
}

