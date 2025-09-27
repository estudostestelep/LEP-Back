package server

import (
	"fmt"
	"log"
	"lep/config"
	"lep/repositories/models"
	"lep/utils"

	"gorm.io/gorm"
)

func Start(db *gorm.DB) {
	//var user models.User
	//result := db.Where("id = ?", 1).First(&user)
	//if result.Error != nil {
	//	fmt.Println("error:", result.Error)
	//}

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
		if err := db.AutoMigrate(model); err != nil {
			panic("error during migration")
		}
	}

	// Check if this is the first run and auto-seed if enabled
	if config.IsAutoSeedEnabled() && isFirstRun(db) {
		log.Println("🌱 First run detected - running auto seed...")
		if err := runFirstTimeSeed(db); err != nil {
			log.Printf("Warning: Auto-seed failed: %v", err)
		} else {
			log.Println("✅ Auto-seed completed successfully")
		}
	}
}

// isFirstRun checks if this is the first time the application is running
// by checking if any users exist in the database
func isFirstRun(db *gorm.DB) bool {
	var count int64
	err := db.Model(&models.User{}).Count(&count).Error
	if err != nil {
		log.Printf("Warning: Could not check if first run: %v", err)
		return false
	}
	return count == 0
}

// runFirstTimeSeed runs the seeding process for first-time setup
func runFirstTimeSeed(db *gorm.DB) error {
	log.Println("  📊 Generating seed data...")
	seedData := utils.GenerateCompleteData()

	// Create organizations first
	log.Println("  🏢 Creating organizations...")
	for _, org := range seedData.Organizations {
		if err := db.Create(&org).Error; err != nil {
			return fmt.Errorf("failed to create organization: %w", err)
		}
	}

	// Create projects
	log.Println("  📁 Creating projects...")
	for _, project := range seedData.Projects {
		if err := db.Create(&project).Error; err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}
	}

	// Create users
	log.Println("  👥 Creating users...")
	for _, user := range seedData.Users {
		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Create environments
	log.Println("  🏗️ Creating environments...")
	for _, env := range seedData.Environments {
		if err := db.Create(&env).Error; err != nil {
			return fmt.Errorf("failed to create environment: %w", err)
		}
	}

	// Create tables
	log.Println("  🪑 Creating tables...")
	for _, table := range seedData.Tables {
		if err := db.Create(&table).Error; err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	// Create products
	log.Println("  🍽️ Creating products...")
	for _, product := range seedData.Products {
		if err := db.Create(&product).Error; err != nil {
			return fmt.Errorf("failed to create product: %w", err)
		}
	}

	// Create customers
	log.Println("  👤 Creating customers...")
	for _, customer := range seedData.Customers {
		if err := db.Create(&customer).Error; err != nil {
			return fmt.Errorf("failed to create customer: %w", err)
		}
	}

	// Create settings
	log.Println("  ⚙️ Creating settings...")
	for _, settings := range seedData.Settings {
		if err := db.Create(&settings).Error; err != nil {
			return fmt.Errorf("failed to create settings: %w", err)
		}
	}

	// Create notification templates
	log.Println("  📧 Creating notification templates...")
	for _, template := range seedData.Templates {
		if err := db.Create(&template).Error; err != nil {
			return fmt.Errorf("failed to create template: %w", err)
		}
	}

	log.Printf("  ✅ Successfully created %d users, %d products, %d tables, and more...",
		len(seedData.Users), len(seedData.Products), len(seedData.Tables))

	return nil
}
