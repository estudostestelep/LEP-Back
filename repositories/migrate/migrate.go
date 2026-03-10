package migrate

import (
	"fmt"

	"gorm.io/gorm"
)

type resourceMigrate struct {
	db *gorm.DB
}

type IMigrate interface {
	MigrateRun(modelsToMigrate ...interface{})
}

func (r *resourceMigrate) MigrateRun(modelsToMigrate ...interface{}) {
	// Migração automática para novas tabelas
	for _, model := range modelsToMigrate {
		migrator := r.db.Migrator()
		if !migrator.HasTable(model) {
			if err := r.db.AutoMigrate(model); err != nil {
				panic(fmt.Sprintf("erro ao migrar tabela %T: %v", model, err))
			}
		}
	}

	// AutoMigrate para adicionar novos campos em tabelas existentes
	if err := r.db.AutoMigrate(modelsToMigrate...); err != nil {
		panic(fmt.Sprintf("erro na migrate AutoMigrate: %v", err))
	}
}

func NewConnMigrate(db *gorm.DB) IMigrate {
	return &resourceMigrate{db: db}
}
