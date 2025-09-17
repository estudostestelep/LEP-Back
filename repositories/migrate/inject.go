package migrate

import "gorm.io/gorm"

type Migrate struct {
	MigratePostgres IMigrate
}

func (h *Migrate) Inject(db *gorm.DB) {
	h.MigratePostgres = NewConnMigrate(db)
}
