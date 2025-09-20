package main

import (
	"fmt"
	"lep/config"
	"lep/middleware"
	"lep/resource"
	"lep/routes"
	"lep/utils"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize configuration and environment
	initializeEnvironment()

	// Initialize resources and database
	resource.Inject()

	// Setup Gin framework
	gin.SetMode(config.GIN_MODE)
	r := gin.Default()

	// Setup CORS based on environment
	setupCORS(r)

	// Setup basic health endpoints
	setupHealthEndpoints(r)

	// Initialize cron jobs if enabled
	initializeCronJobs()

	// Setup middlewares and routes
	setupMiddlewares(r)
	setupRoutes(r)

	// Start server
	startServer(r)
}

func initializeEnvironment() {
	log.Printf("Starting LEP System in environment: %s", config.ENV)
	log.Printf("Port: %s", config.PORT)
	log.Printf("Gin Mode: %s", config.GIN_MODE)
	log.Printf("Log Level: %s", config.LOG_LEVEL)
	log.Printf("Cron Jobs Enabled: %t", config.ENABLE_CRON_JOBS)

	// Environment-specific logging
	if config.IsLocalDev() {
		log.Println("🚀 Running in LOCAL DEVELOPMENT mode")
		log.Printf("Database: %s@%s:%s/%s", config.DB_USER, config.DB_HOST, config.DB_PORT, config.DB_NAME)
	} else if config.IsGCP() {
		log.Printf("☁️ Running on GOOGLE CLOUD PLATFORM - %s", config.ENV)
		if config.INSTANCE_UNIX_SOCKET != "" {
			log.Printf("Database: Cloud SQL via socket %s", config.INSTANCE_UNIX_SOCKET)
		}
	}
}

func setupCORS(r *gin.Engine) {
	corsConfig := cors.DefaultConfig()

	if config.IsDev() {
		// Permissive CORS for development
		corsConfig.AllowAllOrigins = true
		log.Println("CORS: Allowing all origins (development mode)")
	} else {
		// Restrictive CORS for production
		corsConfig.AllowOrigins = []string{
			"https://yourdomain.com",
			"https://www.yourdomain.com",
		}
		log.Println("CORS: Restricted origins (production mode)")
	}

	r.Use(cors.New(corsConfig))
}

func setupHealthEndpoints(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/health", func(c *gin.Context) {
		response := gin.H{
			"status":      "healthy",
			"environment": config.ENV,
			"version":     "1.0.0", // You can make this dynamic
		}

		// Add environment-specific health info
		if config.IsLocalDev() {
			response["mode"] = "local-development"
		} else if config.IsGCP() {
			response["mode"] = "gcp"
			response["platform"] = "cloud-run"
		}

		c.JSON(http.StatusOK, response)
	})
}

func initializeCronJobs() {
	if config.ENABLE_CRON_JOBS {
		log.Println("Initializing cron services...")

		// Use the global Repository that was initialized in resource.Inject()
		cronService := utils.NewCronService(&resource.Repository)
		cronService.StartCronJobs()
		log.Println("Cron services started successfully")
	} else {
		log.Printf("Cron jobs disabled for environment: %s", config.ENV)
	}
}

func setupMiddlewares(r *gin.Engine) {
	// Authentication middleware
	r.Use(middleware.AuthMiddleware())

	// Header validation middleware
	r.Use(middleware.HeaderValidationMiddleware())
}

func setupRoutes(r *gin.Engine) {
	// Setup all application routes
	routes.SetupRoutes(r)
}

func startServer(r *gin.Engine) {
	port := fmt.Sprintf(":%s", config.PORT)

	log.Printf("🌟 LEP System starting on port %s", config.PORT)
	log.Printf("📍 Environment: %s", config.ENV)

	if config.IsLocalDev() {
		log.Printf("🔗 Health check: http://localhost:%s/health", config.PORT)
		log.Printf("🔗 API docs: http://localhost:%s/ping", config.PORT)
	}

	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
