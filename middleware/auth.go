package middleware

import (
	"lep/resource"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isPublicRoute(c.Request.Method, c.Request.URL.Path) {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autorização necessário"})
			c.Abort()
			return
		}

		token := extractToken(authHeader)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			c.Abort()
			return
		}

		if !resource.ServersControllers.SourceAuth.ServiceValidateTokenIn(c) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func isPublicRoute(method, path string) bool {
	publicRoutes := map[string][]string{
		"POST": {"/login", "/user"},
		"GET":  {"/ping", "/health"},
	}

	// Add webhook routes
	if strings.HasPrefix(path, "/webhook/") {
		return true
	}

	if methods, exists := publicRoutes[method]; exists {
		for _, route := range methods {
			if path == route {
				return true
			}
		}
	}

	return false
}

func extractToken(authHeader string) string {
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return authHeader
}