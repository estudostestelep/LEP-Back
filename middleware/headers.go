package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// HeaderValidationMiddleware valida headers obrigatórios X-Lpe-Organization-Id e X-Lpe-Project-Id
func HeaderValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip validation for specific routes
		if isExemptRoute(c.Request.URL.Path, c.Request.Method) {
			c.Next()
			return
		}

		// Validate Organization ID
		organizationId := c.GetHeader("X-Lpe-Organization-Id")
		if strings.TrimSpace(organizationId) == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "X-Lpe-Organization-Id header is required",
			})
			c.Abort()
			return
		}

		// Validate Project ID
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

// isExemptRoute verifica se a rota deve ser isenta da validação de headers
func isExemptRoute(path, method string) bool {
	// Routes that don't require organization/project headers
	exemptRoutes := []RoutePattern{
		{"/login", "POST"},
		{"/user", "POST"},           // Public user creation
		{"/ping", "GET"},
		{"/health", "GET"},
		{"/webhook/*", "*"},         // All webhook routes
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