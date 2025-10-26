package repositories

import (
	"lep/repositories/migrate"

	"gorm.io/gorm"
)

type DBconn struct {
	AuditLogs            IAuditLogsRepository
	BannedLists          IBannedListsRepository
	Customers            ICustomersRepository
	LoggedLists          ILoggedListsRepository
	Orders               IOrderRepository
	Organizations        IOrganizationRepository
	Products             IProductRepository
	Reservations         IReservationRepository
	Tables               ITableRepository
	User                 IUserRepository
	UserOrganizations    IUserOrganizationRepository
	UserProjects         IUserProjectRepository
	Waitlists            WaitlistRepositoryInterface
	KitchenQueue         IKitchenQueueRepository
	Projects             IProjectRepository
	Settings             ISettingsRepository
	Environments         IEnvironmentRepository
	Notifications        INotificationRepository
	Tags                 ITagRepository
	Menus                IMenuRepository
	Categories           ICategoryRepository
	Subcategories        ISubcategoryRepository
	EntityFileReference  IEntityFileReferenceRepository
	Migrate              migrate.IMigrate
}

func (r *DBconn) InjectProstgres(db *gorm.DB) {
	r.Migrate = migrate.NewConnMigrate(db)
	r.Organizations = NewConnOrganization(db)
	r.User = NewUserRepository(db)
	r.UserOrganizations = NewUserOrganizationRepository(db)
	r.UserProjects = NewUserProjectRepository(db)
	r.BannedLists = NewConnBannedLists(db)
	r.LoggedLists = NewConnLoggedLists(db)
	r.Products = NewConnProduct(db)
	r.Customers = NewConnCustomer(db)
	r.Orders = NewConnOrder(db)
	r.Tables = NewConnTable(db)
	r.AuditLogs = NewConnAuditLog(db)
	r.Reservations = NewConnReservation(db)
	r.Waitlists = NewWaitlistRepository(db)
	r.KitchenQueue = NewKitchenQueueRepository(db)
	r.Projects = NewProjectRepository(db)
	r.Settings = NewSettingsRepository(db)
	r.Environments = NewEnvironmentRepository(db)
	r.Notifications = NewNotificationRepository(db)
	r.Tags = NewConnTag(db)
	r.Menus = NewConnMenu(db)
	r.Categories = NewConnCategory(db)
	r.Subcategories = NewConnSubcategory(db)
	r.EntityFileReference = NewEntityFileReferenceRepository(db)

}
