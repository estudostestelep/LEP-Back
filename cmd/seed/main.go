package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"lep/repositories/models"
	"lep/resource"
	"lep/routes"
	"lep/utils"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var (
	clearFirst  bool
	environment string
	verbose     bool
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "seed",
		Short: "LEP Database Seeder",
		Long:  `Populate the LEP database with realistic sample data for development and testing.`,
		Run:   runSeed,
	}

	rootCmd.Flags().BoolVar(&clearFirst, "clear-first", false, "Clear existing data before seeding")
	rootCmd.Flags().StringVar(&environment, "environment", "dev", "Environment to seed (dev, test, staging)")
	rootCmd.Flags().BoolVar(&verbose, "verbose", false, "Enable verbose logging")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runSeed(cmd *cobra.Command, args []string) {
	fmt.Println("\n🌱 LEP Database Seeder")
	fmt.Println("======================")
	fmt.Printf("Environment: %s\n", environment)
	fmt.Printf("Clear first: %t\n", clearFirst)
	fmt.Printf("Verbose: %t\n\n", verbose)

	// Initialize configuration
	flag.Parse()
	//config.LoadEnv()

	// Connect to database
	fmt.Println("📡 Connecting to database...")
	db, err := resource.OpenConnDBPostgres2()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate to ensure all tables exist
	fmt.Println("🔄 Running auto-migration...")
	err = runMigrations(db)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Clear data if requested
	if clearFirst {
		fmt.Println("🧹 Clearing existing data...")
		err = clearExistingData(db)
		if err != nil {
			log.Fatalf("Failed to clear data: %v", err)
		}
	}

	// Generate seed data
	fmt.Println("📊 Generating seed data...")
	seedData := utils.GenerateCompleteData()

	// Initialize server handlers for seeding
	fmt.Println("🔧 Initializing server handlers...")
	router := setupTestRouter()

	// Seed the database using real server routes
	fmt.Println("🌱 Seeding database via server routes...")
	err = seedDatabaseViaServer(router, seedData)
	if err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// Print summary
	printSeedingSummary(db)

	fmt.Println("\n✅ Database seeding completed successfully!")
	fmt.Println("🚀 You can now start the LEP backend and see the data in action.")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Start the backend: go run main.go")
	fmt.Println("  2. Check health: curl http://localhost:8080/health")
	fmt.Println("  3. Login with: admin@lep-demo.com / password")
}

func runMigrations(db *gorm.DB) error {
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
		&models.BannedLists{},
		&models.LoggedLists{},
		&models.AuditLog{},
	)
}

func clearExistingData(db *gorm.DB) error {
	// Order matters due to foreign key constraints
	tables := []string{
		"audit_logs",
		"logged_lists",
		"banned_lists",
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
		if verbose {
			fmt.Printf("  Clearing table: %s\n", table)
		}

		// Use TRUNCATE for better performance, fall back to DELETE if not supported
		result := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table))
		if result.Error != nil {
			// Fallback to DELETE
			result = db.Exec(fmt.Sprintf("DELETE FROM %s", table))
			if result.Error != nil {
				return fmt.Errorf("failed to clear table %s: %v", table, result.Error)
			}
		}
	}

	return nil
}

func seedDatabase(db *gorm.DB, data *utils.SeedData) error {
	// Seed in dependency order

	fmt.Println("  📋 Seeding organizations...")
	for _, org := range data.Organizations {
		if err := createOrUpdate(db, &models.Organization{}, org.Id, org); err != nil {
			return fmt.Errorf("failed to seed organization: %v", err)
		}
		if verbose {
			fmt.Printf("    ✓ %s\n", org.Name)
		}
	}

	fmt.Println("  📁 Seeding projects...")
	for _, project := range data.Projects {
		if err := createOrUpdate(db, &models.Project{}, project.Id, project); err != nil {
			return fmt.Errorf("failed to seed project: %v", err)
		}
		if verbose {
			fmt.Printf("    ✓ %s\n", project.Name)
		}
	}

	fmt.Println("  👥 Seeding users...")
	for _, user := range data.Users {
		if err := createOrUpdate(db, &models.User{}, user.Id, user); err != nil {
			return fmt.Errorf("failed to seed user: %v", err)
		}
		if verbose {
			fmt.Printf("    ✓ %s (%s)\n", user.Name, user.Role)
		}
	}

	fmt.Println("  🏢 Seeding environments...")
	for _, env := range data.Environments {
		if err := createOrUpdate(db, &models.Environment{}, env.Id, env); err != nil {
			return fmt.Errorf("failed to seed environment: %v", err)
		}
		if verbose {
			fmt.Printf("    ✓ %s\n", env.Name)
		}
	}

	fmt.Println("  👤 Seeding customers...")
	for _, customer := range data.Customers {
		if err := createOrUpdate(db, &models.Customer{}, customer.Id, customer); err != nil {
			return fmt.Errorf("failed to seed customer: %v", err)
		}
		if verbose {
			fmt.Printf("    ✓ %s\n", customer.Name)
		}
	}

	fmt.Println("  🍽️  Seeding products...")
	for _, product := range data.Products {
		if err := createOrUpdate(db, &models.Product{}, product.Id, product); err != nil {
			return fmt.Errorf("failed to seed product: %v", err)
		}
		if verbose {
			fmt.Printf("    ✓ %s - R$ %.2f\n", product.Name, product.Price)
		}
	}

	fmt.Println("  🪑 Seeding tables...")
	for _, table := range data.Tables {
		if err := createOrUpdate(db, &models.Table{}, table.Id, table); err != nil {
			return fmt.Errorf("failed to seed table: %v", err)
		}
		if verbose {
			fmt.Printf("    ✓ Mesa %d (%s) - %s\n", table.Number, table.Status, table.Location)
		}
	}

	fmt.Println("  📝 Seeding orders...")
	for _, order := range data.Orders {
		if err := createOrUpdate(db, &models.Order{}, order.Id, order); err != nil {
			return fmt.Errorf("failed to seed order: %v", err)
		}
		if verbose {
			fmt.Printf("    ✓ Order %s - %s (R$ %.2f)\n", order.Id.String()[:8], order.Status, order.TotalAmount)
		}
	}

	fmt.Println("  🎫 Seeding reservations...")
	for _, reservation := range data.Reservations {
		if err := createOrUpdate(db, &models.Reservation{}, reservation.Id, reservation); err != nil {
			return fmt.Errorf("failed to seed reservation: %v", err)
		}
		if verbose {
			fmt.Printf("    ✓ %s - %s (%d pessoas)\n", reservation.Datetime, reservation.Status, reservation.PartySize)
		}
	}

	fmt.Println("  ⏰ Seeding waitlists...")
	for _, waitlist := range data.Waitlists {
		if err := createOrUpdate(db, &models.Waitlist{}, waitlist.Id, waitlist); err != nil {
			return fmt.Errorf("failed to seed waitlist: %v", err)
		}
		if verbose {
			fmt.Printf("    ✓ %s - %d pessoas (%d min)\n", waitlist.Status)
		}
	}

	fmt.Println("  ⚙️  Seeding settings...")
	for _, setting := range data.Settings {
		if err := createOrUpdate(db, &models.Settings{}, setting.Id, setting); err != nil {
			return fmt.Errorf("failed to seed settings: %v", err)
		}
		if verbose {
			fmt.Printf("    ✓ Settings configured\n")
		}
	}

	fmt.Println("  📧 Seeding notification templates...")
	for _, template := range data.Templates {
		if err := createOrUpdate(db, &models.NotificationTemplate{}, template.Id, template); err != nil {
			return fmt.Errorf("failed to seed notification template: %v", err)
		}
		if verbose {
			fmt.Printf("    ✓ %s (%s)\n", template.Channel)
		}
	}

	return nil
}

func setupTestRouter() *gin.Engine {

	// Initialize all resources and handlers
	resource.Inject()

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create router and setup routes
	router := gin.New()
	router.Use(gin.Recovery())

	// Setup all routes
	routes.SetupRoutes(router)

	return router
}

func seedDatabaseViaServer(router *gin.Engine, data *utils.SeedData) error {
	// 1. Criar organização primeiro (sem auth)
	fmt.Println("  📋 Criando organização...")
	orgId, err := createOrganization(router, data.Organizations[0])
	if err != nil {
		return fmt.Errorf("failed to create organization: %v", err)
	}

	// 2. Criar projeto usando o ID da organização
	fmt.Println("  📁 Criando projeto...")
	projectId, err := createProject(router, data.Projects[0], orgId)
	if err != nil {
		return fmt.Errorf("failed to create project: %v", err)
	}

	// 3. Criar usuário admin usando org e project IDs
	fmt.Println("  👥 Criando usuário admin...")
	adminUser := data.Users[0]
	adminUser.OrganizationId = orgId
	adminUser.ProjectId = projectId
	adminToken, err := createAdminUser(router, adminUser)
	if err != nil {
		return fmt.Errorf("failed to create admin user: %v", err)
	}

	// 4. Configurar headers para requests autenticados
	headers := map[string]string{
		"Content-Type":          "application/json",
		"Authorization":         "Bearer " + adminToken,
		"X-Lpe-Organization-Id": orgId.String(),
		"X-Lpe-Project-Id":      projectId.String(),
	}

	// 5. Criar demais usuários (com auth)
	if len(data.Users) > 1 {
		fmt.Println("  👥 Criando demais usuários...")
		for i := 1; i < len(data.Users); i++ {
			user := data.Users[i]
			user.OrganizationId = orgId
			user.ProjectId = projectId
			if err := createUser(router, user, headers); err != nil {
				return fmt.Errorf("failed to create user %s: %v", user.Name, err)
			}
		}
	}

	// 6. Criar environments
	if len(data.Environments) > 0 {
		fmt.Println("  🏢 Criando environments...")
		for _, env := range data.Environments {
			env.OrganizationId = orgId
			env.ProjectId = projectId
			if err := createEnvironment(router, env, headers); err != nil {
				return fmt.Errorf("failed to create environment %s: %v", env.Name, err)
			}
		}
	}

	// 7. Criar customers
	if len(data.Customers) > 0 {
		fmt.Println("  👤 Criando customers...")
		for _, customer := range data.Customers {
			customer.OrganizationId = orgId
			customer.ProjectId = projectId
			if err := createCustomer(router, customer, headers); err != nil {
				return fmt.Errorf("failed to create customer %s: %v", customer.Name, err)
			}
		}
	}

	// 8. Criar products
	if len(data.Products) > 0 {
		fmt.Println("  🍽️  Criando products...")
		for _, product := range data.Products {
			product.OrganizationId = orgId
			product.ProjectId = projectId
			if err := createProduct(router, product, headers); err != nil {
				return fmt.Errorf("failed to create product %s: %v", product.Name, err)
			}
		}
	}

	// 9. Criar tables
	if len(data.Tables) > 0 {
		fmt.Println("  🪑 Criando tables...")
		for _, table := range data.Tables {
			table.OrganizationId = orgId
			table.ProjectId = projectId
			if err := createTable(router, table, headers); err != nil {
				return fmt.Errorf("failed to create table %d: %v", table.Number, err)
			}
		}
	}

	// 10. Criar orders
	if len(data.Orders) > 0 {
		fmt.Println("  📝 Criando orders...")
		for _, order := range data.Orders {
			order.OrganizationId = orgId
			order.ProjectId = projectId
			if err := createOrder(router, order, headers); err != nil {
				return fmt.Errorf("failed to create order: %v", err)
			}
		}
	}

	// 11. Criar reservations
	if len(data.Reservations) > 0 {
		fmt.Println("  🎫 Criando reservations...")
		for _, reservation := range data.Reservations {
			reservation.OrganizationId = orgId
			reservation.ProjectId = projectId
			if err := createReservation(router, reservation, headers); err != nil {
				return fmt.Errorf("failed to create reservation: %v", err)
			}
		}
	}

	// 12. Criar waitlists
	if len(data.Waitlists) > 0 {
		fmt.Println("  ⏰ Criando waitlists...")
		for _, waitlist := range data.Waitlists {
			waitlist.OrganizationId = orgId
			waitlist.ProjectId = projectId
			if err := createWaitlist(router, waitlist, headers); err != nil {
				return fmt.Errorf("failed to create waitlist: %v", err)
			}
		}
	}

	//// 13. Criar settings
	//if len(data.Settings) > 0 {
	//	fmt.Println("  ⚙️  Criando settings...")
	//	setting := data.Settings[0]
	//	setting.OrganizationId = orgId
	//	setting.ProjectId = projectId
	//	if err := createSettings(router, setting, headers); err != nil {
	//		return fmt.Errorf("failed to create settings: %v", err)
	//	}
	//}
	//
	//// 14. Criar notification templates
	//if len(data.Templates) > 0 {
	//	fmt.Println("  📧 Criando notification templates...")
	//	for _, template := range data.Templates {
	//		template.OrganizationId = orgId
	//		template.ProjectId = projectId
	//		if err := createNotificationTemplate(router, template, headers); err != nil {
	//			return fmt.Errorf("failed to create notification template: %v", err)
	//		}
	//	}
	//}

	fmt.Printf("\n✅ Seeding concluído com sucesso!")
	fmt.Printf("\n📋 Organization ID: %s", orgId)
	fmt.Printf("\n📁 Project ID: %s", projectId)

	return nil
}

// Criar organização e retornar o ID (sem headers - bootstrap)
func createOrganization(router *gin.Engine, org models.Organization) (uuid.UUID, error) {
	body, _ := json.Marshal(org)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/organization", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// Sem headers - rota isenta de validação
	router.ServeHTTP(w, req)

	if w.Code != 201 && w.Code != 409 { // 409 if already exists
		return uuid.Nil, fmt.Errorf("failed to create organization: status %d - %s", w.Code, w.Body.String())
	}

	if verbose {
		fmt.Printf("    ✓ %s\n", org.Name)
	}

	return org.Id, nil
}

// Criar projeto usando organization ID e retornar project ID (apenas org header)
func createProject(router *gin.Engine, project models.Project, orgId uuid.UUID) (uuid.UUID, error) {
	project.OrganizationId = orgId
	body, _ := json.Marshal(project)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/project", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Lpe-Organization-Id", orgId.String())
	// Apenas organization header - project header não é necessário
	router.ServeHTTP(w, req)

	if w.Code != 201 && w.Code != 409 {
		return uuid.Nil, fmt.Errorf("failed to create project: status %d - %s", w.Code, w.Body.String())
	}

	if verbose {
		fmt.Printf("    ✓ %s\n", project.Name)
	}

	return project.Id, nil
}

// Criar usuário admin e retornar token de login (sem headers - rota pública)
func createAdminUser(router *gin.Engine, user models.User) (string, error) {
	// Criar usuário
	body, _ := json.Marshal(user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// Sem headers - rota pública de criação de usuário
	router.ServeHTTP(w, req)

	if w.Code != 201 && w.Code != 409 {
		return "", fmt.Errorf("failed to create admin user: status %d - %s", w.Code, w.Body.String())
	}

	// Fazer login para obter token
	loginData := map[string]string{
		"email":    user.Email,
		"password": "password", // Senha padrão do seed data
	}
	loginBody, _ := json.Marshal(loginData)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		return "", fmt.Errorf("failed to login admin user: status %d - %s", w.Code, w.Body.String())
	}

	var loginResponse map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &loginResponse); err != nil {
		return "", fmt.Errorf("failed to parse login response: %v", err)
	}

	token, ok := loginResponse["token"].(string)
	if !ok {
		return "", fmt.Errorf("no token in login response")
	}

	if verbose {
		fmt.Printf("    ✓ %s (%s)\n", user.Name, user.Role)
	}

	return token, nil
}

// Criar usuário comum (com autenticação)
func createUser(router *gin.Engine, user models.User, headers map[string]string) error {
	body, _ := json.Marshal(user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != 201 && w.Code != 409 {
		return fmt.Errorf("status %d - %s", w.Code, w.Body.String())
	}

	if verbose {
		fmt.Printf("    ✓ %s (%s)\n", user.Name, user.Role)
	}

	return nil
}

// Criar environment
func createEnvironment(router *gin.Engine, env models.Environment, headers map[string]string) error {
	body, _ := json.Marshal(env)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/environment", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	router.ServeHTTP(w, req)

	if w.Code != 201 && w.Code != 409 {
		return fmt.Errorf("status %d - %s", w.Code, w.Body.String())
	}

	if verbose {
		fmt.Printf("    ✓ %s\n", env.Name)
	}

	return nil
}

// Criar customer
func createCustomer(router *gin.Engine, customer models.Customer, headers map[string]string) error {
	body, _ := json.Marshal(customer)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/customer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	router.ServeHTTP(w, req)

	if w.Code != 201 && w.Code != 409 {
		return fmt.Errorf("status %d - %s", w.Code, w.Body.String())
	}

	if verbose {
		fmt.Printf("    ✓ %s\n", customer.Name)
	}

	return nil
}

// Criar product
func createProduct(router *gin.Engine, product models.Product, headers map[string]string) error {
	body, _ := json.Marshal(product)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	router.ServeHTTP(w, req)

	if w.Code != 201 && w.Code != 409 {
		return fmt.Errorf("status %d - %s", w.Code, w.Body.String())
	}

	if verbose {
		fmt.Printf("    ✓ %s - R$ %.2f\n", product.Name, product.Price)
	}

	return nil
}

// Criar table
func createTable(router *gin.Engine, table models.Table, headers map[string]string) error {
	body, _ := json.Marshal(table)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/table", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	router.ServeHTTP(w, req)

	if w.Code != 201 && w.Code != 409 {
		return fmt.Errorf("status %d - %s", w.Code, w.Body.String())
	}

	if verbose {
		fmt.Printf("    ✓ Mesa %d (%s) - %s\n", table.Number, table.Status, table.Location)
	}

	return nil
}

// Criar order
func createOrder(router *gin.Engine, order models.Order, headers map[string]string) error {
	body, _ := json.Marshal(order)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/order", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	router.ServeHTTP(w, req)

	if w.Code != 201 && w.Code != 409 {
		return fmt.Errorf("status %d - %s", w.Code, w.Body.String())
	}

	if verbose {
		fmt.Printf("    ✓ Order %s - %s (R$ %.2f)\n", order.Id.String()[:8], order.Status, order.TotalAmount)
	}

	return nil
}

// Criar reservation
func createReservation(router *gin.Engine, reservation models.Reservation, headers map[string]string) error {
	body, _ := json.Marshal(reservation)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/reservation", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	router.ServeHTTP(w, req)

	if w.Code != 201 && w.Code != 409 {
		return fmt.Errorf("status %d - %s", w.Code, w.Body.String())
	}

	if verbose {
		fmt.Printf("    ✓ %s - %s (%d pessoas)\n", reservation.Datetime, reservation.Status, reservation.PartySize)
	}

	return nil
}

// Criar waitlist
func createWaitlist(router *gin.Engine, waitlist models.Waitlist, headers map[string]string) error {
	body, _ := json.Marshal(waitlist)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/waitlist", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	router.ServeHTTP(w, req)

	if w.Code != 201 && w.Code != 409 {
		return fmt.Errorf("status %d - %s", w.Code, w.Body.String())
	}

	if verbose {
		fmt.Printf("    ✓ %s\n", waitlist.Status)
	}

	return nil
}

// Criar settings
func createSettings(router *gin.Engine, settings models.Settings, headers map[string]string) error {
	body, _ := json.Marshal(settings)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/settings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	router.ServeHTTP(w, req)

	if w.Code != 200 && w.Code != 201 {
		return fmt.Errorf("status %d - %s", w.Code, w.Body.String())
	}

	if verbose {
		fmt.Printf("    ✓ Settings configurado\n")
	}

	return nil
}

// Criar notification template
func createNotificationTemplate(router *gin.Engine, template models.NotificationTemplate, headers map[string]string) error {
	body, _ := json.Marshal(template)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/notification/template", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	router.ServeHTTP(w, req)

	if w.Code != 201 && w.Code != 409 {
		return fmt.Errorf("status %d - %s", w.Code, w.Body.String())
	}

	if verbose {
		fmt.Printf("    ✓ %s\n", template.Channel)
	}

	return nil
}

func createOrUpdate(db *gorm.DB, model interface{}, id interface{}, data interface{}) error {
	// Check if record exists
	result := db.Where("id = ?", id).First(model)

	if result.Error == gorm.ErrRecordNotFound {
		// Create new record
		return db.Create(data).Error
	} else if result.Error != nil {
		// Other error
		return result.Error
	} else {
		// Record exists, update it
		return db.Where("id = ?", id).Updates(data).Error
	}
}

func printSeedingSummary(db *gorm.DB) {
	fmt.Println("\n📊 Seeding Summary:")
	fmt.Println("==================")

	tables := map[string]string{
		"organizations":          "🏢 Organizations",
		"projects":               "📁 Projects",
		"users":                  "👥 Users",
		"customers":              "👤 Customers",
		"products":               "🍽️  Products",
		"tables":                 "🪑 Tables",
		"orders":                 "📝 Orders",
		"reservations":           "🎫 Reservations",
		"waitlists":              "⏰ Waitlists",
		"environments":           "🏢 Environments",
		"settings":               "⚙️  Settings",
		"notification_templates": "📧 Templates",
	}

	for table, label := range tables {
		var count int64
		db.Table(table).Count(&count)
		fmt.Printf("%-20s %d records\n", label, count)
	}

	fmt.Println("\n🎯 Sample Data Available:")
	fmt.Println("========================")
	fmt.Println("📧 Login Credentials:")
	fmt.Println("  • admin@lep-demo.com / password (Admin)")
	fmt.Println("  • garcom@lep-demo.com / password (Waiter)")
	fmt.Println("  • gerente@lep-demo.com / password (Manager)")
	fmt.Println("")
	fmt.Println("📊 Data Highlights:")
	fmt.Println("  • 12 products across 3 categories")
	fmt.Println("  • 8 tables with different statuses")
	fmt.Println("  • 4 active orders in various stages")
	fmt.Println("  • 6 reservations (past, present, future)")
	fmt.Println("  • 3 waitlist entries")
	fmt.Println("  • 5 customers with preferences")
	fmt.Println("  • 5 notification templates")
}
