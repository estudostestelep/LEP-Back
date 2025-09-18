package main

import (
	"lep/middleware"
	"lep/resource"
	"lep/routes"
	"net/http"

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

	// Authentication middleware
	r.Use(middleware.AuthMiddleware())

	// Setup all routes
	routes.SetupRoutes(r)

	r.Run(":8080")
}
