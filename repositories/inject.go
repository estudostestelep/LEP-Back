package repositories

import (
	"gorm.io/gorm"
)

type DBconn struct {
	AuditLogs           IAuditLogsRepository
	AccessLogs          IAccessLogRepository
	BannedLists         IBannedListsRepository
	Customers           ICustomersRepository
	LoggedLists         ILoggedListsRepository
	Orders              IOrderRepository
	Organizations       IOrganizationRepository
	Products            IProductRepository
	Reservations        IReservationRepository
	Tables              ITableRepository
	User                IUserRepository
	// Novos repositórios para Admin e Client
	Admins              IAdminRepository
	Clients             IClientRepository
	Waitlists           WaitlistRepositoryInterface
	KitchenQueue        IKitchenQueueRepository
	Projects            IProjectRepository
	Settings            ISettingsRepository
	DisplaySettings     IDisplaySettingsRepository
	ThemeCustomization  IThemeCustomizationRepository
	Environments        IEnvironmentRepository
	Notifications       INotificationRepository
	Tags                ITagRepository
	Menus               IMenuRepository
	Categories          ICategoryRepository
	Subcategories       ISubcategoryRepository
	EntityFileReference IEntityFileReferenceRepository
	FileReference       IFileReferenceRepository
	// Role & Permission system
	Roles       IRoleRepository
	Permissions IPermissionRepository
	Modules     IModuleRepository
	Packages    IPackageRepository
	// Plan Change Request
	PlanChangeRequests IPlanChangeRequestRepository
	// Sidebar Config
	SidebarConfig ISidebarConfigRepository
	// Cascade Delete
	CascadeDelete ICascadeDeleteRepository
	// Admin Audit Logs (read-only)
	AdminAuditLogs IAdminAuditLogRepository
	// Client Audit Logs (optional module)
	ClientAuditLogs IClientAuditLogRepository
}

func (r *DBconn) InjectPostgres(db *gorm.DB) {
	r.Organizations = NewConnOrganization(db)
	r.User = NewUserRepository(db)
	// Novos repositórios para Admin e Client
	r.Admins = NewAdminRepository(db)
	r.Clients = NewClientRepository(db)
	r.BannedLists = NewConnBannedLists(db)
	r.LoggedLists = NewConnLoggedLists(db)
	r.Products = NewConnProduct(db)
	r.Customers = NewConnCustomer(db)
	r.Orders = NewConnOrder(db)
	r.Tables = NewConnTable(db)
	r.AuditLogs = NewConnAuditLog(db)
	r.AccessLogs = NewAccessLogRepository(db)
	r.Reservations = NewConnReservation(db)
	r.Waitlists = NewWaitlistRepository(db)
	r.KitchenQueue = NewKitchenQueueRepository(db)
	r.Projects = NewProjectRepository(db)
	r.Settings = NewSettingsRepository(db)
	r.DisplaySettings = NewDisplaySettingsRepository(db)
	r.ThemeCustomization = NewThemeCustomizationRepository(db)
	r.Environments = NewEnvironmentRepository(db)
	r.Notifications = NewNotificationRepository(db)
	r.Tags = NewConnTag(db)
	r.Menus = NewConnMenu(db)
	r.Categories = NewConnCategory(db)
	r.Subcategories = NewConnSubcategory(db)
	r.EntityFileReference = NewEntityFileReferenceRepository(db)
	r.FileReference = NewFileReferenceRepository(db)
	// Role & Permission system
	r.Roles = NewRoleRepository(db)
	r.Permissions = NewPermissionRepository(db)
	r.Modules = NewModuleRepository(db)
	r.Packages = NewPackageRepository(db)
	// Plan Change Request
	r.PlanChangeRequests = NewPlanChangeRequestRepository(db)
	// Sidebar Config
	r.SidebarConfig = NewSidebarConfigRepository(db)
	// Cascade Delete
	r.CascadeDelete = NewCascadeDeleteRepository(db)
	// Admin Audit Logs (read-only)
	r.AdminAuditLogs = NewAdminAuditLogRepository(db)
	// Client Audit Logs (optional module)
	r.ClientAuditLogs = NewClientAuditLogRepository(db)
}
