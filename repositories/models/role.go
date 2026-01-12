package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Role representa um cargo no sistema
// Cargos são atribuídos a usuários e definem suas permissões
type Role struct {
	Id             uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	Name           string         `json:"name" gorm:"size:100;not null"`                    // Nome técnico (ex: "admin", "manager")
	DisplayName    string         `json:"display_name" gorm:"size:100"`                     // Nome amigável (ex: "Administrador", "Gerente")
	Description    string         `json:"description" gorm:"size:500"`                      // Descrição do cargo
	HierarchyLevel int            `json:"hierarchy_level" gorm:"default:1"`                 // Nível de hierarquia (1-10)
	Scope          string         `json:"scope" gorm:"size:20;default:'client'"`            // "admin" ou "client"
	Permissions    pq.StringArray `json:"permissions" gorm:"type:text[]"`                   // Array de permissões (legacy)
	OrganizationId *uuid.UUID     `json:"organization_id,omitempty" gorm:"type:uuid;index"` // NULL = cargo global do sistema
	ProjectId      *uuid.UUID     `json:"project_id,omitempty" gorm:"type:uuid;index"`      // NULL = cargo da organização
	IsSystem       bool           `json:"is_system" gorm:"default:false"`                   // Cargo do sistema (não pode ser deletado)
	Active         bool           `json:"active" gorm:"default:true"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      *time.Time     `json:"deleted_at,omitempty" gorm:"index"`
}

func (Role) TableName() string {
	return "roles"
}

// UserRole representa a associação entre usuário e cargo (muitos para muitos)
type UserRole struct {
	Id             uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	UserId         uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	RoleId         uuid.UUID  `json:"role_id" gorm:"type:uuid;not null;index"`
	OrganizationId uuid.UUID  `json:"organization_id" gorm:"type:uuid;not null;index"` // Contexto da organização
	ProjectId      *uuid.UUID `json:"project_id,omitempty" gorm:"type:uuid;index"`     // Contexto do projeto (opcional)
	Active         bool       `json:"active" gorm:"default:true"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relacionamentos para preload
	User *User `json:"user,omitempty" gorm:"foreignKey:UserId"`
	Role *Role `json:"role,omitempty" gorm:"foreignKey:RoleId"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

// RoleWithPermissionLevels representa um cargo com níveis de permissão detalhados
type RoleWithPermissionLevels struct {
	Role
	PermissionLevels []RolePermissionLevel `json:"permission_levels,omitempty"`
}

// RolePermissionLevel define o nível de acesso (0, 1, 2) para cada permissão
type RolePermissionLevel struct {
	Id           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	RoleId       uuid.UUID  `json:"role_id" gorm:"type:uuid;not null;index"`
	PermissionId uuid.UUID  `json:"permission_id" gorm:"type:uuid;not null;index"`
	Level        int        `json:"level" gorm:"default:0"` // 0=sem acesso, 1=leitura, 2=escrita
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relacionamentos
	Role       *Role       `json:"role,omitempty" gorm:"foreignKey:RoleId"`
	Permission *Permission `json:"permission,omitempty" gorm:"foreignKey:PermissionId"`
}

func (RolePermissionLevel) TableName() string {
	return "role_permission_levels"
}

// Permission representa uma permissão no sistema (tabela separada para flexibilidade)
type Permission struct {
	Id          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	CodeName    string     `json:"code_name" gorm:"size:100;not null;uniqueIndex"` // Nome técnico (ex: "client_menu_view")
	DisplayName string     `json:"display_name" gorm:"size:100"`                   // Nome amigável (ex: "Visualizar Cardápio")
	Description string     `json:"description" gorm:"size:500"`
	ModuleId    uuid.UUID  `json:"module_id" gorm:"type:uuid;not null;index"` // Módulo ao qual pertence
	Active      bool       `json:"active" gorm:"default:true"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relacionamentos
	Module *Module `json:"module,omitempty" gorm:"foreignKey:ModuleId"`
}

func (Permission) TableName() string {
	return "permissions"
}

// Module representa um módulo do sistema (agrupamento de permissões)
type Module struct {
	Id           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	CodeName     string     `json:"code_name" gorm:"size:100;not null;uniqueIndex"` // Nome técnico (ex: "menu", "orders")
	DisplayName  string     `json:"display_name" gorm:"size:100"`                   // Nome amigável (ex: "Cardápio", "Pedidos")
	Description  string     `json:"description" gorm:"size:500"`
	Icon         string     `json:"icon" gorm:"size:50"`                       // Ícone do módulo (lucide icon name)
	Scope        string     `json:"scope" gorm:"size:20;default:'client'"`     // "admin" ou "client"
	DisplayOrder int        `json:"display_order" gorm:"default:0"`            // Ordem de exibição
	IsFree       bool       `json:"is_free" gorm:"default:false"`              // Módulo gratuito (disponível para todos)
	Active       bool       `json:"active" gorm:"default:true"`                // Se o módulo está disponível
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relacionamentos
	Permissions []Permission `json:"permissions,omitempty" gorm:"foreignKey:ModuleId"`
}

func (Module) TableName() string {
	return "modules"
}

// Package representa um pacote/plano de assinatura
type Package struct {
	Id           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	CodeName     string     `json:"code_name" gorm:"size:100;not null;uniqueIndex"`
	DisplayName  string     `json:"display_name" gorm:"size:100"`   // Nome amigável (ex: "Plano Básico")
	Description  string     `json:"description" gorm:"size:500"`
	PriceMonthly float64    `json:"price_monthly" gorm:"type:decimal(10,2);default:0"`
	PriceYearly  float64    `json:"price_yearly" gorm:"type:decimal(10,2);default:0"`
	IsPublic     bool       `json:"is_public" gorm:"default:true"`  // Visível publicamente
	DisplayOrder int        `json:"display_order" gorm:"default:0"` // Ordem de exibição
	Active       bool       `json:"active" gorm:"default:true"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relacionamentos (não mapeados diretamente pelo GORM para evitar ciclos)
	Modules []Module `json:"modules,omitempty" gorm:"-"` // Preenchido manualmente
}

func (Package) TableName() string {
	return "packages"
}

// PackageModule associa módulos a pacotes
type PackageModule struct {
	Id        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	PackageId uuid.UUID  `json:"package_id" gorm:"type:uuid;not null;index"`
	ModuleId  uuid.UUID  `json:"module_id" gorm:"type:uuid;not null;index"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relacionamentos
	Package *Package `json:"package,omitempty" gorm:"foreignKey:PackageId"`
	Module  *Module  `json:"module,omitempty" gorm:"foreignKey:ModuleId"`
}

func (PackageModule) TableName() string {
	return "package_modules"
}

// PackageBundle associa pacotes a bundles
type PackageBundle struct {
	Id              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	BundleId        uuid.UUID  `json:"bundle_id" gorm:"type:uuid;not null;index"`  // ID do bundle (Package com IsBundle=true)
	PackageId       uuid.UUID  `json:"package_id" gorm:"type:uuid;not null;index"` // ID do pacote incluído
	DiscountPercent float64    `json:"discount_percent" gorm:"type:decimal(5,2);default:0"`
	CreatedAt       time.Time  `json:"created_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relacionamentos
	Bundle  *Package `json:"bundle,omitempty" gorm:"foreignKey:BundleId"`
	Package *Package `json:"package,omitempty" gorm:"foreignKey:PackageId"`
}

func (PackageBundle) TableName() string {
	return "package_bundles"
}

// PackageLimit define limites de recursos por pacote
type PackageLimit struct {
	Id         uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	PackageId  uuid.UUID  `json:"package_id" gorm:"type:uuid;not null;index"`
	LimitType  string     `json:"limit_type" gorm:"size:50;not null"` // "users", "tables", "products", "reservations_per_day"
	LimitValue int        `json:"limit_value" gorm:"default:-1"`      // -1 = ilimitado
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relacionamentos
	Package *Package `json:"package,omitempty" gorm:"foreignKey:PackageId"`
}

func (PackageLimit) TableName() string {
	return "package_limits"
}

// OrganizationPackage associa pacotes a organizações (assinatura)
type OrganizationPackage struct {
	Id             uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	OrganizationId uuid.UUID  `json:"organization_id" gorm:"type:uuid;not null;index"`
	PackageId      uuid.UUID  `json:"package_id" gorm:"type:uuid;not null;index"`
	BillingCycle   string     `json:"billing_cycle" gorm:"size:20;default:'monthly'"` // "monthly", "yearly"
	CustomPrice    *float64   `json:"custom_price,omitempty" gorm:"type:decimal(10,2)"` // Preço customizado (override)
	StartedAt      *time.Time `json:"started_at,omitempty"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
	Active         bool       `json:"active" gorm:"default:true"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relacionamentos
	Organization *Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationId"`
	Package      *Package      `json:"package,omitempty" gorm:"foreignKey:PackageId"`
}

func (OrganizationPackage) TableName() string {
	return "organization_packages"
}
