package constants

import "strings"

// ==================== Modules ====================

// Modules disponíveis no sistema
const (
	ModuleOrders        = "orders"
	ModuleMenu          = "menu"
	ModuleProducts      = "products"
	ModuleTables        = "tables"
	ModuleReservations  = "reservations"
	ModuleWaitlist      = "waitlist"
	ModuleCustomers     = "customers"
	ModuleUsers         = "users"
	ModuleReports       = "reports"
	ModuleSettings      = "settings"
	ModuleNotifications = "notifications"
	// Admin modules
	ModuleOrganizations = "organizations"
	ModulePlans         = "plans"
)

// AllModules lista todos os módulos disponíveis
var AllModules = []string{
	ModuleOrders,
	ModuleMenu,
	ModuleProducts,
	ModuleTables,
	ModuleReservations,
	ModuleWaitlist,
	ModuleCustomers,
	ModuleUsers,
	ModuleReports,
	ModuleSettings,
	ModuleNotifications,
	ModuleOrganizations,
	ModulePlans,
}

// ClientModules lista módulos disponíveis para clientes
var ClientModules = []string{
	ModuleOrders,
	ModuleMenu,
	ModuleProducts,
	ModuleTables,
	ModuleReservations,
	ModuleWaitlist,
	ModuleCustomers,
	ModuleUsers,
	ModuleReports,
	ModuleSettings,
	ModuleNotifications,
}

// AdminModules lista módulos exclusivos de admin
var AdminModules = []string{
	ModuleOrganizations,
	ModulePlans,
}

// ==================== Actions ====================

// Actions padrão CRUD
const (
	ActionRead   = "read"
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionExport = "export"
	ActionSend   = "send"
)

// AllActions lista todas as ações disponíveis
var AllActions = []string{
	ActionRead,
	ActionCreate,
	ActionUpdate,
	ActionDelete,
	ActionExport,
	ActionSend,
}

// ==================== Permissions ====================

// Permission no formato module:action
type Permission string

// Orders permissions
const (
	PermOrdersRead   Permission = "orders:read"
	PermOrdersCreate Permission = "orders:create"
	PermOrdersUpdate Permission = "orders:update"
	PermOrdersDelete Permission = "orders:delete"
)

// Menu permissions
const (
	PermMenuRead   Permission = "menu:read"
	PermMenuCreate Permission = "menu:create"
	PermMenuUpdate Permission = "menu:update"
	PermMenuDelete Permission = "menu:delete"
)

// Products permissions
const (
	PermProductsRead   Permission = "products:read"
	PermProductsCreate Permission = "products:create"
	PermProductsUpdate Permission = "products:update"
	PermProductsDelete Permission = "products:delete"
)

// Tables permissions
const (
	PermTablesRead   Permission = "tables:read"
	PermTablesCreate Permission = "tables:create"
	PermTablesUpdate Permission = "tables:update"
	PermTablesDelete Permission = "tables:delete"
)

// Reservations permissions
const (
	PermReservationsRead   Permission = "reservations:read"
	PermReservationsCreate Permission = "reservations:create"
	PermReservationsUpdate Permission = "reservations:update"
	PermReservationsDelete Permission = "reservations:delete"
)

// Waitlist permissions
const (
	PermWaitlistRead   Permission = "waitlist:read"
	PermWaitlistCreate Permission = "waitlist:create"
	PermWaitlistUpdate Permission = "waitlist:update"
	PermWaitlistDelete Permission = "waitlist:delete"
)

// Customers permissions
const (
	PermCustomersRead   Permission = "customers:read"
	PermCustomersCreate Permission = "customers:create"
	PermCustomersUpdate Permission = "customers:update"
	PermCustomersDelete Permission = "customers:delete"
)

// Users permissions
const (
	PermUsersRead   Permission = "users:read"
	PermUsersCreate Permission = "users:create"
	PermUsersUpdate Permission = "users:update"
	PermUsersDelete Permission = "users:delete"
)

// Reports permissions
const (
	PermReportsRead   Permission = "reports:read"
	PermReportsExport Permission = "reports:export"
)

// Settings permissions
const (
	PermSettingsRead   Permission = "settings:read"
	PermSettingsUpdate Permission = "settings:update"
)

// Notifications permissions
const (
	PermNotificationsRead   Permission = "notifications:read"
	PermNotificationsCreate Permission = "notifications:create"
	PermNotificationsUpdate Permission = "notifications:update"
	PermNotificationsDelete Permission = "notifications:delete"
	PermNotificationsSend   Permission = "notifications:send"
)

// Admin - Organizations permissions
const (
	PermOrganizationsRead   Permission = "organizations:read"
	PermOrganizationsCreate Permission = "organizations:create"
	PermOrganizationsUpdate Permission = "organizations:update"
	PermOrganizationsDelete Permission = "organizations:delete"
)

// Admin - Plans permissions
const (
	PermPlansRead   Permission = "plans:read"
	PermPlansCreate Permission = "plans:create"
	PermPlansUpdate Permission = "plans:update"
	PermPlansDelete Permission = "plans:delete"
)

// ==================== Permission Groups ====================

// AllPermissions lista todas as permissões do sistema
var AllPermissions = []Permission{
	// Orders
	PermOrdersRead, PermOrdersCreate, PermOrdersUpdate, PermOrdersDelete,
	// Menu
	PermMenuRead, PermMenuCreate, PermMenuUpdate, PermMenuDelete,
	// Products
	PermProductsRead, PermProductsCreate, PermProductsUpdate, PermProductsDelete,
	// Tables
	PermTablesRead, PermTablesCreate, PermTablesUpdate, PermTablesDelete,
	// Reservations
	PermReservationsRead, PermReservationsCreate, PermReservationsUpdate, PermReservationsDelete,
	// Waitlist
	PermWaitlistRead, PermWaitlistCreate, PermWaitlistUpdate, PermWaitlistDelete,
	// Customers
	PermCustomersRead, PermCustomersCreate, PermCustomersUpdate, PermCustomersDelete,
	// Users
	PermUsersRead, PermUsersCreate, PermUsersUpdate, PermUsersDelete,
	// Reports
	PermReportsRead, PermReportsExport,
	// Settings
	PermSettingsRead, PermSettingsUpdate,
	// Notifications
	PermNotificationsRead, PermNotificationsCreate, PermNotificationsUpdate, PermNotificationsDelete, PermNotificationsSend,
	// Admin - Organizations
	PermOrganizationsRead, PermOrganizationsCreate, PermOrganizationsUpdate, PermOrganizationsDelete,
	// Admin - Plans
	PermPlansRead, PermPlansCreate, PermPlansUpdate, PermPlansDelete,
}

// PermissionsByModule agrupa permissões por módulo
var PermissionsByModule = map[string][]Permission{
	ModuleOrders:        {PermOrdersRead, PermOrdersCreate, PermOrdersUpdate, PermOrdersDelete},
	ModuleMenu:          {PermMenuRead, PermMenuCreate, PermMenuUpdate, PermMenuDelete},
	ModuleProducts:      {PermProductsRead, PermProductsCreate, PermProductsUpdate, PermProductsDelete},
	ModuleTables:        {PermTablesRead, PermTablesCreate, PermTablesUpdate, PermTablesDelete},
	ModuleReservations:  {PermReservationsRead, PermReservationsCreate, PermReservationsUpdate, PermReservationsDelete},
	ModuleWaitlist:      {PermWaitlistRead, PermWaitlistCreate, PermWaitlistUpdate, PermWaitlistDelete},
	ModuleCustomers:     {PermCustomersRead, PermCustomersCreate, PermCustomersUpdate, PermCustomersDelete},
	ModuleUsers:         {PermUsersRead, PermUsersCreate, PermUsersUpdate, PermUsersDelete},
	ModuleReports:       {PermReportsRead, PermReportsExport},
	ModuleSettings:      {PermSettingsRead, PermSettingsUpdate},
	ModuleNotifications: {PermNotificationsRead, PermNotificationsCreate, PermNotificationsUpdate, PermNotificationsDelete, PermNotificationsSend},
	ModuleOrganizations: {PermOrganizationsRead, PermOrganizationsCreate, PermOrganizationsUpdate, PermOrganizationsDelete},
	ModulePlans:         {PermPlansRead, PermPlansCreate, PermPlansUpdate, PermPlansDelete},
}

// ==================== Role Presets ====================

// OwnerPermissions - permissões para proprietário da organização
var OwnerPermissions = []Permission{
	// Full CRUD em todos os módulos de cliente
	PermOrdersRead, PermOrdersCreate, PermOrdersUpdate, PermOrdersDelete,
	PermMenuRead, PermMenuCreate, PermMenuUpdate, PermMenuDelete,
	PermProductsRead, PermProductsCreate, PermProductsUpdate, PermProductsDelete,
	PermTablesRead, PermTablesCreate, PermTablesUpdate, PermTablesDelete,
	PermReservationsRead, PermReservationsCreate, PermReservationsUpdate, PermReservationsDelete,
	PermWaitlistRead, PermWaitlistCreate, PermWaitlistUpdate, PermWaitlistDelete,
	PermCustomersRead, PermCustomersCreate, PermCustomersUpdate, PermCustomersDelete,
	PermUsersRead, PermUsersCreate, PermUsersUpdate, PermUsersDelete,
	PermReportsRead, PermReportsExport,
	PermSettingsRead, PermSettingsUpdate,
	PermNotificationsRead, PermNotificationsCreate, PermNotificationsUpdate, PermNotificationsDelete, PermNotificationsSend,
}

// ManagerPermissions - permissões para gerente
var ManagerPermissions = []Permission{
	PermOrdersRead, PermOrdersCreate, PermOrdersUpdate,
	PermMenuRead,
	PermProductsRead,
	PermTablesRead, PermTablesCreate, PermTablesUpdate,
	PermReservationsRead, PermReservationsCreate, PermReservationsUpdate,
	PermWaitlistRead, PermWaitlistCreate, PermWaitlistUpdate,
	PermCustomersRead, PermCustomersCreate, PermCustomersUpdate,
	PermUsersRead,
	PermReportsRead,
	PermSettingsRead,
	PermNotificationsRead, PermNotificationsSend,
}

// WaiterPermissions - permissões para garçom
var WaiterPermissions = []Permission{
	PermOrdersRead, PermOrdersCreate, PermOrdersUpdate,
	PermMenuRead,
	PermProductsRead,
	PermTablesRead, PermTablesUpdate,
	PermCustomersRead,
}

// KitchenPermissions - permissões para cozinha
var KitchenPermissions = []Permission{
	PermOrdersRead, PermOrdersUpdate,
	PermMenuRead,
	PermProductsRead,
}

// ViewerPermissions - permissões somente leitura
var ViewerPermissions = []Permission{
	PermOrdersRead,
	PermMenuRead,
	PermProductsRead,
	PermTablesRead,
	PermReservationsRead,
	PermWaitlistRead,
	PermCustomersRead,
	PermReportsRead,
	PermSettingsRead,
}

// AdminFullPermissions - todas as permissões para admin do sistema
var AdminFullPermissions = []Permission{
	// Todas as permissões de cliente
	PermOrdersRead, PermOrdersCreate, PermOrdersUpdate, PermOrdersDelete,
	PermMenuRead, PermMenuCreate, PermMenuUpdate, PermMenuDelete,
	PermProductsRead, PermProductsCreate, PermProductsUpdate, PermProductsDelete,
	PermTablesRead, PermTablesCreate, PermTablesUpdate, PermTablesDelete,
	PermReservationsRead, PermReservationsCreate, PermReservationsUpdate, PermReservationsDelete,
	PermWaitlistRead, PermWaitlistCreate, PermWaitlistUpdate, PermWaitlistDelete,
	PermCustomersRead, PermCustomersCreate, PermCustomersUpdate, PermCustomersDelete,
	PermUsersRead, PermUsersCreate, PermUsersUpdate, PermUsersDelete,
	PermReportsRead, PermReportsExport,
	PermSettingsRead, PermSettingsUpdate,
	PermNotificationsRead, PermNotificationsCreate, PermNotificationsUpdate, PermNotificationsDelete, PermNotificationsSend,
	// Admin específico
	PermOrganizationsRead, PermOrganizationsCreate, PermOrganizationsUpdate, PermOrganizationsDelete,
	PermPlansRead, PermPlansCreate, PermPlansUpdate, PermPlansDelete,
}

// ==================== Hierarchy ====================

// HierarchyLevel constants
const (
	HierarchyMasterAdmin = 10
	HierarchyAdmin       = 8
	HierarchyOwner       = 7
	HierarchyManager     = 5
	HierarchySupervisor  = 4
	HierarchyAttendant   = 3
	HierarchyWaiter      = 2
	HierarchyViewer      = 1
)

// IsMasterAdminLevel verifica se o nível de hierarquia é master admin
func IsMasterAdminLevel(level int) bool {
	return level >= HierarchyMasterAdmin
}

// ==================== Helper Functions ====================

// ParsePermission extrai module e action de uma permissão
func ParsePermission(perm Permission) (module, action string) {
	parts := strings.Split(string(perm), ":")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", ""
}

// BuildPermission constrói uma permissão a partir de module e action
func BuildPermission(module, action string) Permission {
	return Permission(module + ":" + action)
}

// HasPermission verifica se uma permissão está na lista
func HasPermission(permissions []string, required Permission) bool {
	for _, p := range permissions {
		if p == string(required) {
			return true
		}
	}
	return false
}

// HasAnyPermission verifica se alguma das permissões requeridas está presente
func HasAnyPermission(permissions []string, required []Permission) bool {
	for _, req := range required {
		if HasPermission(permissions, req) {
			return true
		}
	}
	return false
}

// HasAllPermissions verifica se todas as permissões requeridas estão presentes
func HasAllPermissions(permissions []string, required []Permission) bool {
	for _, req := range required {
		if !HasPermission(permissions, req) {
			return false
		}
	}
	return true
}

// IsValidPermission verifica se uma string é uma permissão válida
func IsValidPermission(perm string) bool {
	for _, valid := range AllPermissions {
		if string(valid) == perm {
			return true
		}
	}
	return false
}

// IsValidModule verifica se um módulo é válido
func IsValidModule(module string) bool {
	for _, m := range AllModules {
		if m == module {
			return true
		}
	}
	return false
}

// GetModulePermissions retorna todas as permissões de um módulo
func GetModulePermissions(module string) []Permission {
	return PermissionsByModule[module]
}

// ConvertPermissionsToStrings converte slice de Permission para slice de string
func ConvertPermissionsToStrings(permissions []Permission) []string {
	result := make([]string, len(permissions))
	for i, p := range permissions {
		result[i] = string(p)
	}
	return result
}

// ConvertStringsToPermissions converte slice de string para slice de Permission
func ConvertStringsToPermissions(permissions []string) []Permission {
	result := make([]Permission, len(permissions))
	for i, p := range permissions {
		result[i] = Permission(p)
	}
	return result
}

// ==================== Display Names ====================

// ModuleDisplayNames mapeia códigos de módulo para nomes amigáveis
var ModuleDisplayNames = map[string]string{
	ModuleOrders:        "Pedidos",
	ModuleMenu:          "Cardápio",
	ModuleProducts:      "Produtos",
	ModuleTables:        "Mesas",
	ModuleReservations:  "Reservas",
	ModuleWaitlist:      "Fila de Espera",
	ModuleCustomers:     "Clientes",
	ModuleUsers:         "Usuários",
	ModuleReports:       "Relatórios",
	ModuleSettings:      "Configurações",
	ModuleNotifications: "Notificações",
	ModuleOrganizations: "Organizações",
	ModulePlans:         "Planos",
}

// ActionDisplayNames mapeia códigos de ação para nomes amigáveis
var ActionDisplayNames = map[string]string{
	ActionRead:   "Visualizar",
	ActionCreate: "Criar",
	ActionUpdate: "Editar",
	ActionDelete: "Excluir",
	ActionExport: "Exportar",
	ActionSend:   "Enviar",
}

// GetPermissionDisplayName retorna o nome amigável de uma permissão
func GetPermissionDisplayName(perm Permission) string {
	module, action := ParsePermission(perm)
	moduleName := ModuleDisplayNames[module]
	actionName := ActionDisplayNames[action]
	if moduleName == "" {
		moduleName = module
	}
	if actionName == "" {
		actionName = action
	}
	return actionName + " " + moduleName
}

