package migrate

import (
	"gorm.io/gorm"
)

type resourceMigrate struct {
	db *gorm.DB
}

type IMigrate interface {
	MigrateRun(modelsToMigrate ...interface{})
}

func (r *resourceMigrate) MigrateRun(modelsToMigrate ...interface{}) {
	for _, model := range modelsToMigrate {
		migrator := r.db.Migrator()
		if !migrator.HasTable(model) {
			if err := r.db.AutoMigrate(model); err != nil {
				panic("erro na migrate")
			}
		}
	}

	// user := &models.User{
	// 	Name:     "Test User",
	// 	Email:    "test@example.com",
	// 	Password: "$2a$10$6hOKDVLp9LWPa3MslIorkuzntcXH49TcAVo.3ZrLMn2r5gJYCrXiK", //12345
	// }

	// if err := r.db.Create(user).Error; err != nil {
	// 	fmt.Println("erro ao criar usuario ou usuario já existente")
	// }
}

func NewConnMigrate(db *gorm.DB) IMigrate {
	return &resourceMigrate{db: db}
}
