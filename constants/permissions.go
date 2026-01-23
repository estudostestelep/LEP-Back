package constants

import "fmt"

// Permission constants for user permissions
type Permission string

// User-Level Permissions (stored in User.Permissions)
const (
	// Master Admin - Full system access
	PermissionMasterAdmin Permission = "master_admin"

	// User Management
	PermissionManageUsers Permission = "manage_users"
	PermissionViewUsers   Permission = "view_users"

	// Organization Management
	PermissionManageOrganizations Permission = "manage_organizations"
	PermissionViewOrganizations   Permission = "view_organizations"

	// Project Management
	PermissionManageProjects Permission = "manage_projects"
	PermissionViewProjects   Permission = "view_projects"

	// Product Management
	PermissionManageProducts Permission = "manage_products"
	PermissionViewProducts   Permission = "view_products"

	// Menu Management (Cardápio)
	PermissionManageMenus Permission = "manage_menus"
	PermissionViewMenus   Permission = "view_menus"

	// Order Management
	PermissionManageOrders Permission = "manage_orders"
	PermissionViewOrders   Permission = "view_orders"

	// Customer Management
	PermissionManageCustomers Permission = "manage_customers"
	PermissionViewCustomers   Permission = "view_customers"

	// Table Management
	PermissionManageTables Permission = "manage_tables"
	PermissionViewTables   Permission = "view_tables"

	// Reservation Management
	PermissionManageReservations Permission = "manage_reservations"
	PermissionViewReservations   Permission = "view_reservations"

	// Waitlist Management
	PermissionManageWaitlists Permission = "manage_waitlists"
	PermissionViewWaitlists   Permission = "view_waitlists"

	// Reports and Analytics
	PermissionViewReports Permission = "view_reports"
	PermissionExportData  Permission = "export_data"

	// Settings and Configuration
	PermissionManageSettings Permission = "manage_settings"
	PermissionViewSettings   Permission = "view_settings"

	// Notifications
	PermissionManageNotifications Permission = "manage_notifications"
	PermissionViewNotifications   Permission = "view_notifications"

	// Tags and Categories
	PermissionManageTags       Permission = "manage_tags"
	PermissionManageCategories Permission = "manage_categories"
)

// Organization-Level Roles (stored in UserOrganization.Role)
type OrganizationRole string

const (
	// Owner - Full organization access
	OrgRoleOwner OrganizationRole = "owner"

	// Admin - Can manage all organization resources
	OrgRoleAdmin OrganizationRole = "admin"

	// Manager - Can manage some resources
	OrgRoleManager OrganizationRole = "manager"

	// Member - Limited access
	OrgRoleMember OrganizationRole = "member"
)

// Project-Level Roles (stored in UserProject.Role)
type ProjectRole string

const (
	// Admin - Full project access
	ProjectRoleAdmin ProjectRole = "admin"

	// Manager - Can manage most project resources
	ProjectRoleManager ProjectRole = "manager"

	// Supervisor - Can view and update resources
	ProjectRoleSupervisor ProjectRole = "supervisor"

	// Waiter - Can view orders and tables
	ProjectRoleWaiter ProjectRole = "waiter"

	// Kitchen - Can view and update order status
	ProjectRoleKitchen ProjectRole = "kitchen"

	// Viewer - Read-only access
	ProjectRoleViewer ProjectRole = "viewer"
)

// Permission Groups - Common permission combinations
var (
	// Master Admin has all permissions
	MasterAdminPermissions = []Permission{
		PermissionMasterAdmin,
		PermissionManageUsers,
		PermissionManageOrganizations,
		PermissionManageProjects,
		PermissionManageProducts,
		PermissionManageMenus,
		PermissionManageOrders,
		PermissionManageCustomers,
		PermissionManageTables,
		PermissionManageReservations,
		PermissionManageWaitlists,
		PermissionViewReports,
		PermissionExportData,
		PermissionManageSettings,
		PermissionManageNotifications,
		PermissionManageTags,
		PermissionManageCategories,
	}

	// Restaurant Admin - Can manage everything except users and organizations
	RestaurantAdminPermissions = []Permission{
		PermissionManageProjects,
		PermissionManageProducts,
		PermissionManageMenus,
		PermissionManageOrders,
		PermissionManageCustomers,
		PermissionManageTables,
		PermissionManageReservations,
		PermissionManageWaitlists,
		PermissionViewReports,
		PermissionExportData,
		PermissionManageSettings,
		PermissionManageNotifications,
		PermissionManageTags,
		PermissionManageCategories,
	}

	// Manager - Can manage orders, tables, reservations, customers
	ManagerPermissions = []Permission{
		PermissionManageOrders,
		PermissionManageCustomers,
		PermissionManageTables,
		PermissionManageReservations,
		PermissionManageWaitlists,
		PermissionViewReports,
	}

	// Waiter - Can view orders and tables only
	WaiterPermissions = []Permission{
		PermissionViewOrders,
		PermissionViewTables,
		PermissionManageOrders,
		PermissionViewCustomers,
	}

	// Kitchen - Can view and update orders only
	KitchenPermissions = []Permission{
		PermissionViewOrders,
		PermissionManageOrders,
	}

	// Viewer - Read-only access
	ViewerPermissions = []Permission{
		PermissionViewProducts,
		PermissionViewMenus,
		PermissionViewOrders,
		PermissionViewCustomers,
		PermissionViewTables,
		PermissionViewReservations,
		PermissionViewWaitlists,
		PermissionViewReports,
		PermissionViewSettings,
		PermissionViewNotifications,
	}
)

// Helper functions

// HasPermission checks if a permission is in the list
func HasPermission(permissions []string, required Permission) bool {
	for _, p := range permissions {
		if p == string(required) {
			return true
		}
	}
	return false
}

// HasAnyPermission checks if any of the required permissions are present
func HasAnyPermission(permissions []string, required []Permission) bool {
	for _, req := range required {
		if HasPermission(permissions, req) {
			return true
		}
	}
	return false
}

// HasAllPermissions checks if all required permissions are present
func HasAllPermissions(permissions []string, required []Permission) bool {
	for _, req := range required {
		if !HasPermission(permissions, req) {
			return false
		}
	}
	return true
}

// IsMasterAdmin checks if user is a master admin
func IsMasterAdmin(permissions []string) bool {
	return HasPermission(permissions, PermissionMasterAdmin)
}

// ConvertPermissions converts []Permission to []string
func ConvertPermissions(permissions []Permission) []string {
	result := make([]string, len(permissions))
	for i, p := range permissions {
		result[i] = string(p)
	}
	return result
}

// AllValidPermissions lista todas as permissões válidas do sistema
var AllValidPermissions = []Permission{
	PermissionMasterAdmin,
	PermissionManageUsers,
	PermissionViewUsers,
	PermissionManageOrganizations,
	PermissionViewOrganizations,
	PermissionManageProjects,
	PermissionViewProjects,
	PermissionManageProducts,
	PermissionViewProducts,
	PermissionManageMenus,
	PermissionViewMenus,
	PermissionManageOrders,
	PermissionViewOrders,
	PermissionManageCustomers,
	PermissionViewCustomers,
	PermissionManageTables,
	PermissionViewTables,
	PermissionManageReservations,
	PermissionViewReservations,
	PermissionManageWaitlists,
	PermissionViewWaitlists,
	PermissionViewReports,
	PermissionExportData,
	PermissionManageSettings,
	PermissionViewSettings,
	PermissionManageNotifications,
	PermissionViewNotifications,
	PermissionManageTags,
	PermissionManageCategories,
}

// IsValidPermission verifica se uma string é uma permissão válida
func IsValidPermission(perm string) bool {
	for _, valid := range AllValidPermissions {
		if string(valid) == perm {
			return true
		}
	}
	return false
}

// ValidatePermissions valida uma lista de permissões
// Retorna erro se alguma permissão for inválida ou vazia
func ValidatePermissions(permissions []string) error {
	for _, perm := range permissions {
		if perm == "" {
			return fmt.Errorf("permissão não pode ser vazia")
		}
		if !IsValidPermission(perm) {
			return fmt.Errorf("permissão inválida: %s", perm)
		}
	}
	return nil
}
