package middleware

import (
	"lep/config"
	"lep/resource"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
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

		tokenString := extractToken(authHeader)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			c.Abort()
			return
		}

		// Extrair claims do token para identificar tipo de usuário
		claims, err := parseTokenClaims(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
			c.Abort()
			return
		}

		// Identificar tipo de usuário do token
		userType, _ := claims["user_type"].(string)
		userId, _ := claims["user_id"].(string)
		userEmail, _ := claims["email"].(string)

		// Se for token novo (com user_type), usar validação por tipo
		if userType == "admin" {
			// Validar admin
			admin, err := resource.Handlers.HandlerAdminUser.GetAdminById(userId)
			if err != nil || admin == nil || !admin.IsActive() {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
				c.Abort()
				return
			}
			c.Set("user_id", admin.Id.String())
			c.Set("user_email", admin.Email)
			// Buscar permissões via roles
			var adminPermissions []string
			adminRoles, _ := resource.Handlers.HandlerAdminUser.GetAdminRoles(admin.Id.String())
			if len(adminRoles) > 0 {
				for _, ar := range adminRoles {
					if ar.Active {
						perms, _ := resource.Handlers.HandlerAdminUser.GetPermissionsFromRole(ar.RoleId.String())
						adminPermissions = append(adminPermissions, perms...)
						break
					}
				}
			}
			c.Set("user_permissions", adminPermissions)
			c.Set("user_type", "admin")
			c.Set("admin", admin)
		} else if userType == "client" {
			// Validar client
			client, err := resource.Handlers.HandlerClientUser.GetClientById(userId)
			if err != nil || client == nil || !client.IsActive() {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
				c.Abort()
				return
			}
			c.Set("user_id", client.Id.String())
			c.Set("user_email", client.Email)
			// Buscar permissões via roles
			var clientPermissions []string
			clientRoles, _ := resource.Handlers.HandlerRole.GetClientRoles(client.Id.String(), client.OrgId.String())
			if len(clientRoles) > 0 {
				for _, cr := range clientRoles {
					if cr.Active {
						perms, _ := resource.Handlers.HandlerAdminUser.GetPermissionsFromRole(cr.RoleId.String())
						clientPermissions = append(clientPermissions, perms...)
						break
					}
				}
			}
			c.Set("user_permissions", clientPermissions)
			c.Set("user_type", "client")
			c.Set("org_id", client.OrgId.String())
			c.Set("proj_ids", client.ProjIds)
			c.Set("client", client)
		} else {
			// Fallback: Token legado (sem user_type) - usar validação antiga
			user, err := resource.Handlers.HandlerAuth.VerificationToken(tokenString)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
				c.Abort()
				return
			}

			if user != nil {
				c.Set("user_id", user.Id.String())
				c.Set("user_email", user.Email)
				c.Set("user_type", user.UserType)
				c.Set("user", user)
			}
		}

		// Definir user_id e user_email de qualquer forma para compatibilidade
		if userId != "" {
			c.Set("user_id", userId)
		}
		if userEmail != "" {
			c.Set("user_email", userEmail)
		}

		c.Next()
	}
}

// parseTokenClaims extrai os claims do token JWT
func parseTokenClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT_SECRET_PRIVATE_KEY), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

func isPublicRoute(method, path string) bool {
	publicRoutes := map[string][]string{
		"POST": {"/login", "/user", "/admin/login", "/client/login"},
		"GET":  {"/ping", "/health", "/tenant/resolve"},
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