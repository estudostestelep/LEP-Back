package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ClientAuditLog - Log de auditoria para ações de clientes
// Registra operações CRUD em entidades dentro de uma organização/projeto
// Este módulo é opcional e configurável por organização
type ClientAuditLog struct {
	Id             uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrganizationId uuid.UUID      `gorm:"type:uuid;not null;index" json:"organization_id"`
	ProjectId      uuid.UUID      `gorm:"type:uuid;not null;index" json:"project_id"`
	UserId         *uuid.UUID     `gorm:"type:uuid;index" json:"user_id,omitempty"`
	UserEmail      string         `gorm:"type:varchar(255)" json:"user_email,omitempty"`
	Action         string         `gorm:"type:varchar(50);not null;index" json:"action"` // CREATE, UPDATE, DELETE, STATUS_CHANGE
	EntityType     string         `gorm:"type:varchar(50);not null;index" json:"entity_type"` // reservation, order, customer, product, table
	EntityId       uuid.UUID      `gorm:"type:uuid;not null" json:"entity_id"`
	ModuleCode     string         `gorm:"type:varchar(50);not null;index" json:"module_code"` // client_reservations, client_orders, etc.
	OldValues      datatypes.JSON `gorm:"type:jsonb" json:"old_values,omitempty"`
	NewValues      datatypes.JSON `gorm:"type:jsonb" json:"new_values,omitempty"`
	ChangedFields  pq.StringArray `gorm:"type:text[]" json:"changed_fields"`
	Description    string         `gorm:"type:text" json:"description,omitempty"`
	IpAddress      string         `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	CreatedAt      time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
}

// TableName define o nome da tabela no banco de dados
func (ClientAuditLog) TableName() string {
	return "client_audit_logs"
}

// ClientAuditConfig - Configuração do módulo de log de auditoria para clientes
// Cada organização pode ter sua própria configuração
type ClientAuditConfig struct {
	Id             uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrganizationId uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"organization_id"`
	Enabled        bool           `gorm:"default:false" json:"enabled"` // Módulo ativado?
	MaxLogsStored  int            `gorm:"default:10000" json:"max_logs_stored"` // Limite de logs salvos
	RetentionDays  int            `gorm:"default:90" json:"retention_days"` // Dias para manter logs
	EnabledModules pq.StringArray `gorm:"type:text[]" json:"enabled_modules"` // Módulos ativos: ["reservations", "orders", "customers"]
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName define o nome da tabela no banco de dados
func (ClientAuditConfig) TableName() string {
	return "client_audit_configs"
}

// Constantes para os tipos de ação de cliente
const (
	ClientAuditActionCreate       = "CREATE"
	ClientAuditActionUpdate       = "UPDATE"
	ClientAuditActionDelete       = "DELETE"
	ClientAuditActionStatusChange = "STATUS_CHANGE"
)

// Constantes para tipos de entidade de cliente
const (
	ClientAuditEntityReservation = "reservation"
	ClientAuditEntityOrder       = "order"
	ClientAuditEntityCustomer    = "customer"
	ClientAuditEntityProduct     = "product"
	ClientAuditEntityTable       = "table"
	ClientAuditEntityMenu        = "menu"
	ClientAuditEntityCategory    = "category"
	ClientAuditEntitySubcategory = "subcategory"
	ClientAuditEntityTag         = "tag"
	ClientAuditEntityWaitlist    = "waitlist"
	ClientAuditEntityUser        = "user"
)

// Constantes para códigos de módulos
const (
	ClientAuditModuleReservations = "client_reservations"
	ClientAuditModuleOrders       = "client_orders"
	ClientAuditModuleCustomers    = "client_customers"
	ClientAuditModuleProducts     = "client_products"
	ClientAuditModuleTables       = "client_tables"
	ClientAuditModuleMenus        = "client_menus"
	ClientAuditModuleWaitlist     = "client_waitlist"
	ClientAuditModuleTags         = "client_tags"
	ClientAuditModuleUsers        = "client_users"
)

// ClientAuditLogFilters - Filtros para listagem de logs de cliente
type ClientAuditLogFilters struct {
	StartDate  *time.Time `json:"start_date,omitempty"`
	EndDate    *time.Time `json:"end_date,omitempty"`
	UserId     *uuid.UUID `json:"user_id,omitempty"`
	UserEmail  string     `json:"user_email,omitempty"`
	Action     string     `json:"action,omitempty"`
	EntityType string     `json:"entity_type,omitempty"`
	ModuleCode string     `json:"module_code,omitempty"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
}

// ClientAuditLogPaginatedResponse - Resposta paginada de logs de cliente
type ClientAuditLogPaginatedResponse struct {
	Data       []ClientAuditLog `json:"data"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

// AvailableClientAuditModules - Lista de módulos disponíveis para auditoria de cliente
var AvailableClientAuditModules = []struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}{
	{ClientAuditModuleReservations, "Reservas", "Registra criação, alteração e cancelamento de reservas"},
	{ClientAuditModuleOrders, "Pedidos", "Registra criação, alteração e status de pedidos"},
	{ClientAuditModuleCustomers, "Clientes", "Registra criação e alteração de dados de clientes"},
	{ClientAuditModuleProducts, "Produtos", "Registra criação, alteração e exclusão de produtos"},
	{ClientAuditModuleTables, "Mesas", "Registra criação, alteração e status de mesas"},
	{ClientAuditModuleMenus, "Cardápios", "Registra criação, alteração e exclusão de cardápios e categorias"},
	{ClientAuditModuleWaitlist, "Fila de Espera", "Registra entrada, alteração de status e remoção da fila de espera"},
	{ClientAuditModuleTags, "Tags", "Registra criação, alteração e exclusão de tags/etiquetas"},
	{ClientAuditModuleUsers, "Usuários", "Registra adição, alteração de permissões e remoção de usuários"},
}
