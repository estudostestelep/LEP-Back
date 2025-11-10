package middleware

import (
	"fmt"
	"lep/config"
	"lep/resource"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
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

		// For POST /menu and POST /category, validate token AND user access to org/project
		if (path == "/menu" || path == "/category") && method == "POST" {
			tokenString := c.GetHeader("Authorization")
			if tokenString == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Authorization token is required",
				})
				c.Abort()
				return
			}
			// Validate user has access to organization and project
			if err := validateUserAccess(c, organizationId, projectId); err != nil {
				c.JSON(http.StatusForbidden, gin.H{
					"error": fmt.Sprintf("Access denied: %v", err),
				})
				c.Abort()
				return
			}
			c.Next()
			return
		}

		// Validar se o usuário logado tem acesso à org e projeto (exceto rotas públicas)
		if err := validateUserAccess(c, organizationId, projectId); err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": fmt.Sprintf("Access denied: %v", err),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// validateUserAccess verifica se o usuário tem acesso à organização e projeto
func validateUserAccess(c *gin.Context, orgId, projId string) error {
	// Extrair user_id do token JWT
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		// Se não há token, retornar erro
		return fmt.Errorf("token de autenticação não encontrado")
	}

	// Remover "Bearer " se presente
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	// Parse do token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT_SECRET_PRIVATE_KEY), nil
	})

	if err != nil {
		return fmt.Errorf("token inválido: %v", err)
	}

	// Verificar se o token é válido
	if !token.Valid {
		return fmt.Errorf("token expirado ou inválido")
	}

	// Extrair claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("claims inválidas no token")
	}

	userIdInterface, exists := claims["user_id"]
	if !exists {
		return fmt.Errorf("user_id não encontrado no token")
	}

	userId, ok := userIdInterface.(string)
	if !ok {
		return fmt.Errorf("user_id inválido no token")
	}

	// Verificar expiração do token
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return fmt.Errorf("token expirado")
		}
	}

	// Armazenar user_id no contexto para uso posterior
	c.Set("user_id", userId)

	// Validar no banco se o usuário tem acesso à org e projeto
	// Usar o Repository global que foi inicializado em resource.Inject()

	// Validar que o usuário tem acesso à organização
	hasOrgAccess, err := resource.Repository.UserOrganizations.UserBelongsToOrganization(userId, orgId)
	if err != nil {
		return fmt.Errorf("erro ao validar acesso à organização: %v", err)
	}
	if !hasOrgAccess {
		return fmt.Errorf("usuário não tem acesso à organização: %s", orgId)
	}

	// Validar que o usuário tem acesso ao projeto (se projId foi fornecido)
	if projId != "" {
		hasProjectAccess, err := resource.Repository.UserProjects.UserBelongsToProject(userId, projId)
		if err != nil {
			return fmt.Errorf("erro ao validar acesso ao projeto: %v", err)
		}
		if !hasProjectAccess {
			return fmt.Errorf("usuário não tem acesso ao projeto: %s", projId)
		}
	}

	return nil
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

		// Public menu routes (no auth required)
		{"/public/menu/*", "*"},
		{"/public/categories/*", "*"},
		{"/public/menus/*", "*"},
		{"/public/project/*", "*"},
		{"/public/times/*", "*"},
		{"/public/reservation/*", "*"},

		// Public upload/static routes (no auth required)
		{"/uploads/*", "GET"},
		{"/static/*", "GET"},

		// User organizational access (has its own validation)
		{"/user/*/organizations-projects", "GET"},  // Get user access - does its own validation
		{"/user/*/organizations-projects", "POST"}, // Update user access - does its own validation
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
			// Check for patterns with wildcard in the middle (e.g., /user/*/organizations-projects)
			if strings.Contains(route.Path, "/*") && !strings.HasSuffix(route.Path, "/*") {
				// Pattern like /user/*/organizations-projects
				// Split by * and check if path matches the pattern
				parts := strings.Split(route.Path, "*")
				if len(parts) == 2 {
					prefix := parts[0]
					suffix := parts[1]
					if strings.HasPrefix(path, prefix) && strings.HasSuffix(path, suffix) {
						return true
					}
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
