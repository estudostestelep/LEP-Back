package tests

import (
	"flag"
	"fmt"
	"lep/config"
	"lep/repositories/models"
	"lep/resource"
	"lep/routes"
	"lep/utils"
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TestConfig holds configuration for testing environment
type TestConfig struct {
	DB           *gorm.DB
	Router       *gin.Engine
	OrgID        string
	ProjectID    string
	AuthToken    string
	TestData     *TestData
	SeedData     *utils.SeedData
}

// NewTestConfig creates a new test configuration with real database and handlers
func NewTestConfig() *TestConfig {
	gin.SetMode(gin.TestMode)

	// Initialize configuration
	flag.Parse()
	config.LoadEnv()

	// Set test environment variables if not set
	if os.Getenv("DB_NAME") == "" {
		os.Setenv("DB_NAME", "lep_test")
	}

	// Connect to test database
	db, err := resource.OpenConnDBPostgres2()
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate test database
	err = runTestMigrations(db)
	if err != nil {
		log.Fatalf("Failed to run test migrations: %v", err)
	}

	// Setup real router with all routes
	router := setupTestRouter(db)

	// Generate test data
	testData := NewTestData()
	seedData := utils.GenerateCompleteData()

	return &TestConfig{
		DB:        db,
		Router:    router,
		OrgID:     testData.OrganizationID,
		ProjectID: testData.ProjectID,
		AuthToken: "test-jwt-token",
		TestData:  testData,
		SeedData:  seedData,
	}
}

// SetupTestDatabase prepares the test database with clean state
func (tc *TestConfig) SetupTestDatabase() error {
	// Clear all tables
	err := tc.clearTestData()
	if err != nil {
		return fmt.Errorf("failed to clear test data: %v", err)
	}

	// Seed with test data
	err = tc.seedTestData()
	if err != nil {
		return fmt.Errorf("failed to seed test data: %v", err)
	}

	return nil
}

// CleanupTestDatabase cleans up test data after tests
func (tc *TestConfig) CleanupTestDatabase() error {
	return tc.clearTestData()
}

// setupTestRouter creates a router with real handlers and middleware
func setupTestRouter(db *gorm.DB) *gin.Engine {
	// Initialize resource container
	resource.SetDB(db)

	// Create router with real routes
	router := gin.New()

	// Add CORS middleware for testing
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Lpe-Organization-Id, X-Lpe-Project-Id")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Setup all real routes
	routes.SetupRoutes(router)

	return router
}

// runTestMigrations runs database migrations for testing
func runTestMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Organization{},
		&models.Project{},
		&models.User{},
		&models.Customer{},
		&models.Product{},
		&models.Table{},
		&models.Order{},
		&models.Reservation{},
		&models.Waitlist{},
		&models.Environment{},
		&models.Settings{},
		&models.NotificationTemplate{},
		&models.NotificationConfig{},
		&models.NotificationLog{},
		&models.NotificationEvent{},
		&models.NotificationInbound{},
		&models.BannedLists{},
		&models.LoggedLists{},
		&models.AuditLog{},
	)
}

// clearTestData removes all test data from database
func (tc *TestConfig) clearTestData() error {
	// Order matters due to foreign key constraints
	tables := []string{
		"audit_logs",
		"logged_lists",
		"banned_lists",
		"notification_inbound",
		"notification_events",
		"notification_logs",
		"notification_configs",
		"notification_templates",
		"settings",
		"environments",
		"waitlists",
		"reservations",
		"orders",
		"tables",
		"products",
		"customers",
		"users",
		"projects",
		"organizations",
	}

	for _, table := range tables {
		result := tc.DB.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if result.Error != nil {
			return fmt.Errorf("failed to clear table %s: %v", table, result.Error)
		}
	}

	return nil
}

// seedTestData populates database with test data
func (tc *TestConfig) seedTestData() error {
	// Seed organizations
	for _, org := range tc.SeedData.Organizations {
		if err := tc.DB.Create(&org).Error; err != nil {
			return fmt.Errorf("failed to seed organization: %v", err)
		}
	}

	// Seed projects
	for _, project := range tc.SeedData.Projects {
		if err := tc.DB.Create(&project).Error; err != nil {
			return fmt.Errorf("failed to seed project: %v", err)
		}
	}

	// Seed users
	for _, user := range tc.SeedData.Users {
		if err := tc.DB.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to seed user: %v", err)
		}
	}

	// Seed environments
	for _, env := range tc.SeedData.Environments {
		if err := tc.DB.Create(&env).Error; err != nil {
			return fmt.Errorf("failed to seed environment: %v", err)
		}
	}

	// Seed customers
	for _, customer := range tc.SeedData.Customers {
		if err := tc.DB.Create(&customer).Error; err != nil {
			return fmt.Errorf("failed to seed customer: %v", err)
		}
	}

	// Seed products
	for _, product := range tc.SeedData.Products {
		if err := tc.DB.Create(&product).Error; err != nil {
			return fmt.Errorf("failed to seed product: %v", err)
		}
	}

	// Seed tables
	for _, table := range tc.SeedData.Tables {
		if err := tc.DB.Create(&table).Error; err != nil {
			return fmt.Errorf("failed to seed table: %v", err)
		}
	}

	// Seed settings
	for _, setting := range tc.SeedData.Settings {
		if err := tc.DB.Create(&setting).Error; err != nil {
			return fmt.Errorf("failed to seed settings: %v", err)
		}
	}

	// Seed notification templates
	for _, template := range tc.SeedData.Templates {
		if err := tc.DB.Create(&template).Error; err != nil {
			return fmt.Errorf("failed to seed notification template: %v", err)
		}
	}

	return nil
}

// TestMain sets up and tears down the test environment
func TestMain(m *testing.M) {
	// Setup
	fmt.Println("🧪 Setting up test environment...")

	// Run tests
	code := m.Run()

	// Cleanup
	fmt.Println("🧹 Cleaning up test environment...")

	os.Exit(code)
}