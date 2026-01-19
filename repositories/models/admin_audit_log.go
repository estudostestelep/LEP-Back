package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// AdminAuditLog - Log de auditoria para ações administrativas
// Registra todas as operações de CRUD realizadas por Master Admins
// IMPORTANTE: Este modelo é READ-ONLY após criação - não permite Update ou Delete
type AdminAuditLog struct {
	Id               uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ActorId          uuid.UUID      `gorm:"type:uuid;not null;index" json:"actor_id"`                      // ID do admin que executou a ação
	ActorEmail       string         `gorm:"type:varchar(255);not null" json:"actor_email"`                 // Email do admin para referência rápida
	TargetId         uuid.UUID      `gorm:"type:uuid;not null;index" json:"target_id"`                     // ID do usuário/entidade afetada
	TargetEmail      string         `gorm:"type:varchar(255)" json:"target_email"`                         // Email do alvo (quando aplicável)
	Action           string         `gorm:"type:varchar(50);not null;index" json:"action"`                 // CREATE, UPDATE, DELETE, RESET_PASSWORD
	EntityType       string         `gorm:"type:varchar(50);not null;index" json:"entity_type"`            // "user", "organization", "project"
	OrganizationId   *uuid.UUID     `gorm:"type:uuid;index" json:"organization_id,omitempty"`              // Null = zona administrativa global
	OrganizationName string         `gorm:"type:varchar(255)" json:"organization_name,omitempty"`          // Nome da org para exibição
	ProjectId        *uuid.UUID     `gorm:"type:uuid;index" json:"project_id,omitempty"`                   // Null = zona administrativa global
	ProjectName      string         `gorm:"type:varchar(255)" json:"project_name,omitempty"`               // Nome do projeto para exibição
	IsAdminZone      bool           `gorm:"default:false" json:"is_admin_zone"`                            // true = "Executado em Zona Administrativa"
	OldValues        datatypes.JSON `gorm:"type:jsonb" json:"old_values,omitempty"`                        // Estado anterior (JSON) - "De"
	NewValues        datatypes.JSON `gorm:"type:jsonb" json:"new_values,omitempty"`                        // Estado novo (JSON) - "Para"
	ChangedFields    pq.StringArray `gorm:"type:text[]" json:"changed_fields"`                             // Lista de campos alterados
	IpAddress        string         `gorm:"type:varchar(45)" json:"ip_address,omitempty"`                  // IP do admin (IPv4/IPv6)
	UserAgent        string         `gorm:"type:varchar(500)" json:"user_agent,omitempty"`                 // User Agent do navegador
	CreatedAt        time.Time      `gorm:"autoCreateTime;index" json:"created_at"`                        // Timestamp da ação
}

// TableName define o nome da tabela no banco de dados
func (AdminAuditLog) TableName() string {
	return "admin_audit_logs"
}

// Constantes para os tipos de ação
const (
	AdminAuditActionCreate        = "CREATE"
	AdminAuditActionUpdate        = "UPDATE"
	AdminAuditActionDelete        = "DELETE"
	AdminAuditActionResetPassword = "RESET_PASSWORD"
	AdminAuditActionAssign        = "ASSIGN"   // Atribuição (ex: cargo a usuário)
	AdminAuditActionRemove        = "REMOVE"   // Remoção (ex: cargo de usuário)
	AdminAuditActionApprove       = "APPROVE"  // Aprovação (ex: solicitação de plano)
	AdminAuditActionReject        = "REJECT"   // Rejeição (ex: solicitação de plano)
	AdminAuditActionCancel        = "CANCEL"   // Cancelamento
	AdminAuditActionCleanup       = "CLEANUP"  // Limpeza de dados
	AdminAuditActionReset         = "RESET"    // Reset de configuração
)

// Constantes para tipos de entidade
const (
	AdminAuditEntityUser             = "user"
	AdminAuditEntityOrganization     = "organization"
	AdminAuditEntityProject          = "project"
	AdminAuditEntityRole             = "role"
	AdminAuditEntityUserRole         = "user_role"         // Atribuição de cargo a usuário
	AdminAuditEntityModule           = "module"            // Módulos do sistema
	AdminAuditEntityPermission       = "permission"        // Permissões do sistema
	AdminAuditEntityPackage          = "package"           // Pacotes/Planos
	AdminAuditEntitySubscription     = "subscription"      // Assinaturas
	AdminAuditEntityPackageModule    = "package_module"    // Módulos em pacotes
	AdminAuditEntityPackageLimit     = "package_limit"     // Limites de pacotes
	AdminAuditEntityPlanRequest      = "plan_request"      // Solicitações de mudança de plano
	AdminAuditEntityMenu             = "menu"              // Menus
	AdminAuditEntityCategory         = "category"          // Categorias
	AdminAuditEntitySidebarConfig    = "sidebar_config"    // Configuração de sidebar
	AdminAuditEntityClientAuditConfig = "client_audit_config" // Configuração de auditoria de cliente
	AdminAuditEntityImage            = "image"             // Imagens
	AdminAuditEntityUserAccess       = "user_access"       // Acesso do usuário a orgs/projetos
	AdminAuditEntityRolePermission   = "role_permission"   // Permissões de cargo
)

// AdminAuditLogFilters - Filtros para listagem de logs
type AdminAuditLogFilters struct {
	StartDate  *time.Time `json:"start_date,omitempty"`
	EndDate    *time.Time `json:"end_date,omitempty"`
	ActorId    *uuid.UUID `json:"actor_id,omitempty"`
	ActorEmail string     `json:"actor_email,omitempty"`
	Action     string     `json:"action,omitempty"`
	EntityType string     `json:"entity_type,omitempty"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
}

// AdminAuditLogPaginatedResponse - Resposta paginada de logs
type AdminAuditLogPaginatedResponse struct {
	Data       []AdminAuditLog `json:"data"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}
