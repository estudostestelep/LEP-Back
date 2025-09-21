package server

import (
	"fmt"
	"lep/repositories/models"

	"gorm.io/gorm"
)

func Start(db *gorm.DB) {
	var user models.User
	result := db.Where("id = ?", 1).First(&user)
	if result.Error != nil {
		fmt.Println("error:", result.Error)
	}

	modelsToMigrate := []interface{}{
		&models.BannedLists{},
		&models.LoggedLists{},

		// Base organization model (must be first due to FK relationships)
		&models.Organization{},

		// Core models
		&models.User{},
		&models.Customer{},
		&models.Table{},
		&models.Product{},
		&models.Reservation{},
		&models.Waitlist{},
		&models.Order{},
		&models.AuditLog{},

		// SPRINT 1 models
		&models.Project{},
		&models.Settings{},
		&models.Environment{},

		// SPRINT 2 models (Notification System)
		&models.NotificationConfig{},
		&models.NotificationTemplate{},
		&models.NotificationLog{},
		&models.NotificationEvent{},
		&models.NotificationInbound{},

		// SPRINT 4 models (Advanced Validations)
		&models.BlockedPeriod{},

		// SPRINT 5 models (Advanced Features)
		&models.Lead{},
		&models.ReportMetric{},
	}

	for _, model := range modelsToMigrate {
		migrator := db.Migrator()
		if !migrator.HasTable(model) {
			if err := db.AutoMigrate(model); err != nil {
				panic("error during migration")
			}
		}
	}

	newUser := &models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "$2a$10$6hOKDVLp9LWPa3MslIorkuzntcXH49TcAVo.3ZrLMn2r5gJYCrXiK", //12345
	}

	if err := db.Create(newUser).Error; err != nil {
		fmt.Println("error creating user or user already exists")
	}
}
