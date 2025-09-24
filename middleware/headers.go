package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// HeaderValidationMiddleware valida headers obrigatórios X-Lpe-Organization-Id e X-Lpe-Project-Id
func HeaderValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method

		// Skip validation for specific routes
		if isFullyExemptRoute(path, method) {
			c.Next()
			return
		}

		// Project creation requires only organization header
		if path == "/project" && method == "POST" {
			organizationId := c.GetHeader("X-Lpe-Organization-Id")
			if strings.TrimSpace(organizationId) == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "X-Lpe-Organization-Id header is required for project creation",
				})
				c.Abort()
				return
			}
			c.Set("organization_id", organizationId)
			c.Next()
			return
		}

		// All other routes require both headers
		organizationId := c.GetHeader("X-Lpe-Organization-Id")
		if strings.TrimSpace(organizationId) == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "X-Lpe-Organization-Id header is required",
			})
			c.Abort()
			return
		}

		projectId := c.GetHeader("X-Lpe-Project-Id")
		if strings.TrimSpace(projectId) == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "X-Lpe-Project-Id header is required",
			})
			c.Abort()
			return
		}

		// Store in context for easy access in handlers
		c.Set("organization_id", organizationId)
		c.Set("project_id", projectId)

		c.Next()
	}
}

// isFullyExemptRoute verifica se a rota deve ser totalmente isenta da validação de headers
func isFullyExemptRoute(path, method string) bool {
	// Routes that don't require any organization/project headers
	exemptRoutes := []RoutePattern{
		{"/login", "POST"},
		{"/user", "POST"},         // Public user creation
		{"/organization", "POST"}, // Organization creation (bootstrap)
		{"/ping", "GET"},
		{"/health", "GET"},
		{"/webhook/*", "*"}, // All webhook routes
	}

	for _, route := range exemptRoutes {
		if route.Method == "*" || route.Method == method {
			if route.Path == path {
				return true
			}
			// Check for wildcard paths
			if strings.HasSuffix(route.Path, "/*") {
				prefix := strings.TrimSuffix(route.Path, "/*")
				if strings.HasPrefix(path, prefix) {
					return true
				}
			}
		}
	}

	return false
}

type RoutePattern struct {
	Path   string
	Method string
}
