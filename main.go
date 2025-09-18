package main

import (
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

	// Authentication middleware
	r.Use(func(c *gin.Context) {
		if c.Request.Method == "POST" && c.Request.URL.Path == "/login" {
			c.Next()
			return
		}
		if c.Request.Method == "POST" && c.Request.URL.Path == "/user" {
			c.Next()
			return
		}

		token := resource.ServersControllers.SourceAuth.ServiceValidateTokenIn(c)

		if token {
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
		}
	})

	// Setup all routes
	routes.SetupRoutes(r)

	r.Run(":8080")
}
