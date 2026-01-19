package server

import (
	"lep/repositories/migrate"
	"lep/repositories/models"

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
		&models.UserOrganization{}, // Relacionamento usuário-organização
		&models.UserProject{},      // Relacionamento usuário-projeto
		&models.Customer{},
		&models.Table{},
		&models.Product{},
		&models.Reservation{},
		&models.Waitlist{},
		&models.Order{},
		&models.AuditLog{},
		&models.AccessLog{}, // User access/login logs

		// SPRINT 1 models
		&models.Project{},
		&models.Settings{},
		&models.ProjectDisplaySettings{}, // Product display visibility settings
		&models.ThemeCustomization{},     // Theme customization colors
		&models.Environment{},

		// SPRINT 2 models (Notification System)
		&models.NotificationConfig{},
		&models.NotificationTemplate{},
		&models.NotificationLog{},
		&models.NotificationEvent{},
		&models.NotificationInbound{},
		&models.NotificationSchedule{},
		&models.ResponseReviewQueue{},

		// SPRINT 4 models (Advanced Validations)
		&models.BlockedPeriod{},

		// SPRINT 5 models (Advanced Features)
		&models.Lead{},
		&models.ReportMetric{},

		// Menu System models
		&models.Tag{},
		&models.ProductTag{},
		&models.Menu{},
		&models.Category{},
		&models.Subcategory{},
		&models.SubcategoryCategory{},

		// Image Management models (Deduplication & References)
		&models.FileReference{},
		&models.EntityFileReference{},

		// Role & Permission System models
		&models.Module{},
		&models.Permission{},
		&models.Role{},
		&models.UserRole{},
		&models.RolePermissionLevel{},
		&models.Package{},
		&models.PackageModule{},
		&models.PackageLimit{},
		&models.PackageBundle{},
		&models.OrganizationPackage{},

		// Plan Change Request System
		&models.PlanChangeRequest{},

		// Sidebar Config System
		&models.SidebarConfig{},

		// Admin Audit Log System (read-only logs for administrative actions)
		&models.AdminAuditLog{},

		// Client Audit Log System (optional module for client-side logging)
		&models.ClientAuditLog{},
		&models.ClientAuditConfig{},
	}

	// Usar migrate customizado para lidar com alterações no Product
	migrator := migrate.NewConnMigrate(db)
	migrator.MigrateRun(modelsToMigrate...)
}
