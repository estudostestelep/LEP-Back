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
	"time"

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

	// Setup routes with conditional middlewares
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
		r.Use(cors.New(cors.Config{
			AllowOrigins: []string{
				"http://localhost:5173",
				"http://localhost:5174",
				"https://lep-front.vercel.app/",
				"http://localhost:5173/",
				"http://localhost:5174/",
				"https://lep-front.vercel.app",
				"https://lep-front-git-main-leps-projects-a55eafc4.vercel.app/",
			},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Lpe-Organization-Id", "X-Lpe-Project-Id"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
		corsConfig.AllowAllOrigins = true
		log.Println("CORS: Allowing all origins (development mode)")
	} else {
		// Restrictive CORS for production
		corsConfig.AllowOrigins = []string{
			"http://localhost:5173",
			"http://localhost:5174",
			"https://lep-front.vercel.app/",
			"http://localhost:5173/",
			"http://localhost:5174/",
			"https://lep-front.vercel.app",
			"https://lep-front-git-main-leps-projects-a55eafc4.vercel.app/",
		}
		corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
		corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "X-Lpe-Organization-Id", "X-Lpe-Project-Id"}
		corsConfig.ExposeHeaders = []string{"Content-Length"}

		log.Println("CORS: Restricted origins (production mode)")
	}

	corsConfig.AllowCredentials = true
	corsConfig.AllowAllOrigins = true

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

	// Configure trusted proxies for Cloud Run
	if config.IsGCP() {
		// Cloud Run uses internal proxies - trust all for now
		r.SetTrustedProxies(nil)
	} else {
		// Local development - no trusted proxies
		r.SetTrustedProxies(nil)
	}

	log.Printf("🌟 LEP System starting on port %s", config.PORT)
	log.Printf("📍 Environment: %s", config.ENV)

	if config.IsLocalDev() {
		log.Printf("🔗 Health check: http://localhost:%s/health", config.PORT)
		log.Printf("🔗 API docs: http://localhost:%s/ping", config.PORT)
	}

	// Use standard HTTP server for better Cloud Run compatibility
	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	log.Printf("🚀 Server listening on %s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
