package main

import (
	"fmt"
	"lep/repositories/models"
	"lep/resource"
	"log"
	"os"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var (
	verbose bool
	dryRun  bool
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "migrate",
		Short: "LEP Database Migration Tool",
		Long: `Executes database migrations using GORM AutoMigrate.

This tool creates or updates all tables in the database according to the models defined in repositories/models/.

It is safe to run multiple times as GORM AutoMigrate:
- Creates tables if they don't exist
- Adds missing columns
- Adds missing indexes
- Does NOT remove columns or alter existing column types`,
		Run: runMigration,
	}

	rootCmd.Flags().BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be migrated without making changes")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runMigration(cmd *cobra.Command, args []string) {
	fmt.Println("\n🔄 LEP Database Migration Tool")
	fmt.Println("================================")

	// Get environment
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "dev"
	}

	fmt.Printf("Environment: %s\n", env)
	fmt.Printf("Verbose: %t\n", verbose)
	fmt.Printf("Dry Run: %t\n\n", dryRun)

	if dryRun {
		fmt.Println("⚠️  DRY RUN MODE - No changes will be made")
		fmt.Println("")
	}

	// Connect to database
	fmt.Println("📡 Connecting to database...")
	db, err := resource.OpenConnDBPostgres()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Get database name
	var dbName string
	db.Raw("SELECT current_database()").Scan(&dbName)
	fmt.Printf("✅ Connected to database: %s\n\n", dbName)

	if dryRun {
		fmt.Println("📋 Would migrate the following tables:")
		printModelList()
		fmt.Println("\n⚠️  Use without --dry-run to apply changes")
		return
	}

	// Run auto-migration
	fmt.Println("🔧 Running auto-migration...")
	err = runMigrations(db)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("✅ Migration completed successfully!")

	// Print summary
	printMigrationSummary(db)

	fmt.Println("\n🎉 All tables are up to date!")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Run seed to populate data: go run cmd/seed/main.go")
	fmt.Println("  2. Or use remote seed: ./lep-seed-remote.exe --url <API_URL>")
}

func runMigrations(db *gorm.DB) error {
	models := []interface{}{
		// Core entities
		&models.Organization{},
		&models.Project{},
		&models.User{},
		&models.UserOrganization{},
		&models.UserProject{},
		&models.Customer{},

		// Menu system
		&models.Menu{},
		&models.Category{},
		&models.Subcategory{},
		&models.SubcategoryCategory{},

		// Tags
		&models.Tag{},
		&models.Product{},
		&models.ProductTag{},

		// Operations
		&models.Table{},
		&models.Order{},
		&models.Reservation{},
		&models.Waitlist{},
		&models.Environment{},
		&models.Settings{},

		// Notifications
		&models.NotificationTemplate{},
		&models.NotificationConfig{},
		&models.NotificationLog{},
		&models.NotificationEvent{},
		&models.NotificationInbound{},

		// Advanced features
		&models.BlockedPeriod{},
		&models.Lead{},
		&models.ReportMetric{},

		// Auth
		&models.BannedLists{},
		&models.LoggedLists{},

		// Audit
		&models.AuditLog{},
	}

	for i, model := range models {
		modelName := getModelName(model)

		if verbose {
			fmt.Printf("  [%d/%d] Migrating %s...\n", i+1, len(models), modelName)
		}

		err := db.AutoMigrate(model)
		if err != nil {
			return fmt.Errorf("failed to migrate %s: %v", modelName, err)
		}

		if verbose {
			fmt.Printf("  ✓ %s migrated successfully\n", modelName)
		}
	}

	if !verbose {
		fmt.Printf("  ✓ Migrated %d models successfully\n", len(models))
	}

	return nil
}

func getModelName(model interface{}) string {
	switch model.(type) {
	case *models.Organization:
		return "Organization"
	case *models.Project:
		return "Project"
	case *models.User:
		return "User"
	case *models.UserOrganization:
		return "UserOrganization"
	case *models.UserProject:
		return "UserProject"
	case *models.Customer:
		return "Customer"
	case *models.Menu:
		return "Menu"
	case *models.Category:
		return "Category"
	case *models.Subcategory:
		return "Subcategory"
	case *models.SubcategoryCategory:
		return "SubcategoryCategory"
	case *models.Tag:
		return "Tag"
	case *models.Product:
		return "Product"
	case *models.ProductTag:
		return "ProductTag"
	case *models.Table:
		return "Table"
	case *models.Order:
		return "Order"
	case *models.Reservation:
		return "Reservation"
	case *models.Waitlist:
		return "Waitlist"
	case *models.Environment:
		return "Environment"
	case *models.Settings:
		return "Settings"
	case *models.NotificationTemplate:
		return "NotificationTemplate"
	case *models.NotificationConfig:
		return "NotificationConfig"
	case *models.NotificationLog:
		return "NotificationLog"
	case *models.NotificationEvent:
		return "NotificationEvent"
	case *models.NotificationInbound:
		return "NotificationInbound"
	case *models.BlockedPeriod:
		return "BlockedPeriod"
	case *models.Lead:
		return "Lead"
	case *models.ReportMetric:
		return "ReportMetric"
	case *models.BannedLists:
		return "BannedLists"
	case *models.LoggedLists:
		return "LoggedLists"
	case *models.AuditLog:
		return "AuditLog"
	default:
		return "Unknown"
	}
}

func printModelList() {
	models := []string{
		"organizations", "projects", "users", "user_organizations", "user_projects",
		"customers", "menus", "categories", "subcategories", "subcategory_categories",
		"tags", "products", "product_tags", "tables", "orders",
		"reservations", "waitlists", "environments", "settings",
		"notification_templates", "notification_configs", "notification_logs",
		"notification_events", "notification_inbounds", "blocked_periods",
		"leads", "report_metrics", "banned_lists", "logged_lists", "audit_logs",
	}

	for i, model := range models {
		fmt.Printf("  %2d. %s\n", i+1, model)
	}
}

func printMigrationSummary(db *gorm.DB) {
	fmt.Println("\n📊 Migration Summary:")
	fmt.Println("====================")

	tables := map[string]string{
		"organizations":          "🏢 Organizations",
		"projects":               "📁 Projects",
		"users":                  "👥 Users",
		"user_organizations":     "🔗 User-Organization Links",
		"user_projects":          "🔗 User-Project Links",
		"customers":              "👤 Customers",
		"menus":                  "📖 Menus",
		"categories":             "📂 Categories",
		"subcategories":          "📂 Subcategories",
		"subcategory_categories": "🔗 Subcategory-Category Links",
		"tags":                   "🏷️  Tags",
		"products":               "🍽️  Products",
		"product_tags":           "🔗 Product-Tag Links",
		"tables":                 "🪑 Tables",
		"orders":                 "📝 Orders",
		"reservations":           "🎫 Reservations",
		"waitlists":              "⏰ Waitlists",
		"environments":           "🏢 Environments",
		"settings":               "⚙️  Settings",
		"notification_templates": "📧 Notification Templates",
		"notification_configs":   "🔔 Notification Configs",
		"notification_logs":      "📜 Notification Logs",
		"notification_events":    "📨 Notification Events",
		"notification_inbounds":  "📥 Inbound Messages",
		"blocked_periods":        "🚫 Blocked Periods",
		"leads":                  "🎯 Leads",
		"report_metrics":         "📊 Report Metrics",
		"banned_lists":           "🚫 Banned Tokens",
		"logged_lists":           "✅ Active Tokens",
		"audit_logs":             "📋 Audit Logs",
	}

	existingCount := 0
	for table, label := range tables {
		var count int64
		result := db.Table(table).Count(&count)
		if result.Error == nil {
			existingCount++
			if verbose {
				fmt.Printf("%-25s %d records\n", label, count)
			}
		}
	}

	fmt.Printf("\n✅ %d/%d tables exist and are accessible\n", existingCount, len(tables))
}
