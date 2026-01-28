package middleware

import (
	"fmt"
	"lep/handler"
	"net/http"
	"reflect"

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
	// Hierarquia de roles: master_admin > admin > manager > waiter > viewer
	roleHierarchy := map[string]int{
		"master_admin": 5, // Nível mais alto - acesso total ao sistema
		"admin":        4,
		"manager":      3,
		"waiter":       2,
		"viewer":       1,
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

		// Handle both []string and pq.StringArray types
		var permissions []string
		switch v := userPermissions.(type) {
		case []string:
			permissions = v
		default:
			// Try to convert to []string for other array types
			if arr, ok := v.([]interface{}); ok {
				for _, item := range arr {
					if str, ok := item.(string); ok {
						permissions = append(permissions, str)
					}
				}
			}
		}

		// Master Admins têm acesso a tudo - bypass automático
		if contains(permissions, "master_admin") {
			c.Next()
			return
		}

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

// IsMasterAdmin verifica se um usuário tem permissão de Master Admin
// Esta função pode ser usada em handlers para lógicas específicas
func IsMasterAdmin(c *gin.Context) bool {
	userPermissions, exists := c.Get("user_permissions")
	if !exists {
		return false
	}

	// Handle both []string and pq.StringArray types
	var permissions []string

	// Try direct cast first
	if strArr, ok := userPermissions.([]string); ok {
		permissions = strArr
	} else {
		// Use reflection to iterate over pq.StringArray
		// pq.StringArray is an array type, so we can use reflection to access elements
		val := reflect.ValueOf(userPermissions)
		if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
			for i := 0; i < val.Len(); i++ {
				elem := val.Index(i)
				if elem.Kind() == reflect.String {
					permissions = append(permissions, elem.String())
				}
			}
		}
	}

	return contains(permissions, "master_admin")
}

// MasterAdminOnlyMiddleware middleware que permite apenas Master Admins
func MasterAdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log para diagnóstico
		userId := c.GetString("user_id")
		userPermissions, _ := c.Get("user_permissions")
		fmt.Printf("[DEBUG MasterAdminOnlyMiddleware] Path=%s, Method=%s, UserId=%s, Permissions=%v\n",
			c.Request.URL.Path, c.Request.Method, userId, userPermissions)

		if !IsMasterAdmin(c) {
			fmt.Printf("[DEBUG MasterAdminOnlyMiddleware] BLOCKED - User %s is not Master Admin\n", userId)
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied: Master Admin only",
			})
			c.Abort()
			return
		}
		fmt.Printf("[DEBUG MasterAdminOnlyMiddleware] ALLOWED - User %s is Master Admin\n", userId)
		c.Next()
	}
}

// AdminScopeMiddleware middleware que permite acesso a usuários com cargos de escopo "admin"
// Isso inclui master_admin (via permissions) e cargos como admin_support, admin_sales
func AdminScopeMiddleware(authHandler handler.IHandlerAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString("user_id")
		userPermissions, _ := c.Get("user_permissions")
		fmt.Printf("[DEBUG AdminScopeMiddleware] Path=%s, Method=%s, UserId=%s, Permissions=%v\n",
			c.Request.URL.Path, c.Request.Method, userId, userPermissions)

		// 1. Master admin tem bypass total (verificação pelo array de permissions legado)
		if IsMasterAdmin(c) {
			fmt.Printf("[DEBUG AdminScopeMiddleware] ALLOWED - User %s is Master Admin (bypass)\n", userId)
			c.Next()
			return
		}

		// 2. Verificar se tem cargo de escopo admin (via tabela user_roles)
		hasAdminScope, err := authHandler.UserHasAdminScopeRole(userId)
		if err != nil {
			fmt.Printf("[DEBUG AdminScopeMiddleware] ERROR checking admin scope for user %s: %v\n", userId, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error checking admin permissions",
			})
			c.Abort()
			return
		}

		if !hasAdminScope {
			fmt.Printf("[DEBUG AdminScopeMiddleware] BLOCKED - User %s has no admin scope role\n", userId)
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied: Admin scope required",
			})
			c.Abort()
			return
		}

		fmt.Printf("[DEBUG AdminScopeMiddleware] ALLOWED - User %s has admin scope role\n", userId)
		c.Next()
	}
}
