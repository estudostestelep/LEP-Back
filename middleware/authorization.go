package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleBasedAuthMiddleware verifica permissões baseadas em roles
func RoleBasedAuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (set by auth middleware)
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User role not found in context",
			})
			c.Abort()
			return
		}

		if !hasPermission(userRole.(string), requiredRole) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// hasPermission verifica se o role atual tem permissão para o role requerido
func hasPermission(userRole, requiredRole string) bool {
	// Hierarquia de roles: admin > manager > waiter > viewer
	roleHierarchy := map[string]int{
		"admin":   4,
		"manager": 3,
		"waiter":  2,
		"viewer":  1,
	}

	userLevel := roleHierarchy[userRole]
	requiredLevel := roleHierarchy[requiredRole]

	return userLevel >= requiredLevel
}

// PermissionMiddleware verifica permissões específicas
func PermissionMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user permissions from context
		userPermissions, exists := c.Get("user_permissions")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User permissions not found in context",
			})
			c.Abort()
			return
		}

		permissions := userPermissions.([]string)
		if !contains(permissions, requiredPermission) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Permission denied: " + requiredPermission,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// contains verifica se uma permissão está na lista
func contains(permissions []string, permission string) bool {
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}