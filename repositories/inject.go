package repositories

import (
	"lep/repositories/migrate"

	"gorm.io/gorm"
)

type DBconn struct {
	AuditLogs    IAuditLogsRepository
	BannedLists  IBannedListsRepository
	Customers    ICustomersRepository
	LoggedLists  ILoggedListsRepository
	Orders       IOrderRepository
	Products     IProductRepository
	Reservations IReservationRepository
	Tables       ITableRepository
	User         IUserRepository
	Migrate      migrate.IMigrate
}

func (r *DBconn) InjectProstgres(db *gorm.DB) {
	r.Migrate = migrate.NewConnMigrate(db)
	r.User = NewUserRepository(db)
	r.BannedLists = NewConnBannedLists(db)
	r.LoggedLists = NewConnLoggedLists(db)
	r.Products = NewConnProduct(db)
	r.Customers = NewConnCustomer(db)
	r.Orders = NewConnOrder(db)
	r.Tables = NewConnTable(db)
	r.AuditLogs = NewConnAuditLog(db)
	r.Reservations = NewConnReservation(db)

}
