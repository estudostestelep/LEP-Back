package main

import (
	"lep/middleware"
	"lep/resource"
	"lep/routes"
	"lep/utils"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	resource.Inject()
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	r.Use(cors.New(config))

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Initialize and start cron jobs if enabled
	enableCronJobs, _ := strconv.ParseBool(os.Getenv("ENABLE_CRON_JOBS"))
	if enableCronJobs {
		log.Println("Initializing cron services...")

		// Use the global Repository that was initialized in resource.Inject()
		cronService := utils.NewCronService(&resource.Repository)
		cronService.StartCronJobs()
		log.Println("Cron services started successfully")
	} else {
		log.Println("Cron jobs disabled - set ENABLE_CRON_JOBS=true to enable")
	}

	// Authentication middleware
	r.Use(middleware.AuthMiddleware())

	// Header validation middleware
	r.Use(middleware.HeaderValidationMiddleware())

	// Setup all routes
	routes.SetupRoutes(r)

	r.Run(":8080")
}
