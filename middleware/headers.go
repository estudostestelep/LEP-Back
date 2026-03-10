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

		fmt.Printf("🔍 HeaderValidationMiddleware: path=%s, method=%s\n", path, method)

		// Skip validation for specific routes
		if isFullyExemptRoute(path, method) {
			fmt.Printf("⏭️ Rota exempta: %s %s\n", method, path)
			c.Next()
			return
		}

		// Admin users don't need org/project headers - check JWT early
		tokenString := c.GetHeader("Authorization")
		if tokenString != "" {
			cleanToken := strings.TrimPrefix(tokenString, "Bearer ")
			parsedToken, parseErr := jwt.Parse(cleanToken, func(t *jwt.Token) (interface{}, error) {
				return []byte(config.JWT_SECRET_PRIVATE_KEY), nil
			})
			if parseErr == nil && parsedToken.Valid {
				if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
					isAdmin := false

					// Verificar user_type == "admin" (admin login)
					userType, _ := claims["user_type"].(string)
					if userType == "admin" {
						isAdmin = true
					}

					// Verificar se é master_admin via hierarquia de roles
					if !isAdmin {
						if userId, ok := claims["user_id"].(string); ok && userId != "" {
							if userType == "" {
								userType = "client" // default para tokens legados
							}
							isMaster, _ := resource.Repository.Roles.IsMasterAdmin(userId, userType)
							if isMaster {
								isAdmin = true
							}
						}
					}

					if isAdmin {
						userId, _ := claims["user_id"].(string)
						c.Set("user_id", userId)
						c.Set("organization_id", c.GetHeader("X-Lpe-Organization-Id"))
						c.Set("project_id", c.GetHeader("X-Lpe-Project-Id"))
						fmt.Printf("⏭️ Admin user bypass: userId=%s\n", userId)
						c.Next()
						return
					}
				}
			}
		}

		// POST /user: Capturar headers mas não validar acesso (usuário está sendo criado)
		if path == "/user" && method == "POST" {
			organizationId := c.GetHeader("X-Lpe-Organization-Id")
			projectId := c.GetHeader("X-Lpe-Project-Id")
			fmt.Printf("🔐 POST /user - Headers capturados: orgId=%s, projectId=%s\n", organizationId, projectId)
			// Setar no contexto para o handler usar (podem estar vazios)
			c.Set("organization_id", organizationId)
			c.Set("project_id", projectId)
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

	// ========== NOVO: Verificar tipo de usuário ==========
	userType, _ := claims["user_type"].(string)

	// ADMIN: tem acesso a todas as organizações e projetos
	if userType == "admin" {
		return nil
	}

	// CLIENT: verificar acesso direto via campos do cliente
	if userType == "client" {
		// Buscar cliente pelo ID
		client, err := resource.Repository.Clients.GetClientById(userId)
		if err != nil || client == nil {
			return fmt.Errorf("cliente não encontrado")
		}

		// Verificar se a org do header é a org do cliente
		if client.OrgId.String() != orgId {
			return fmt.Errorf("usuário não tem acesso à organização: %s", orgId)
		}

		// Verificar se o projeto está na lista de projetos do cliente
		if projId != "" && !client.HasProjectAccess(projId) {
			return fmt.Errorf("usuário não tem acesso ao projeto: %s", projId)
		}

		return nil
	}

	// ========== Validar via client_roles (fallback para tokens legados) ==========
	// Validar que o usuário tem acesso à organização via client_roles
	clientRoles, err := resource.Repository.Roles.GetClientRoles(userId, orgId)
	if err != nil {
		return fmt.Errorf("erro ao validar acesso à organização: %v", err)
	}
	if len(clientRoles) == 0 {
		return fmt.Errorf("usuário não tem acesso à organização: %s", orgId)
	}

	// Validar que o usuário tem acesso ao projeto (se projId foi fornecido)
	if projId != "" {
		hasProjectAccess := false
		for _, role := range clientRoles {
			// Se a role não tem project_id, o usuário tem acesso a todos os projetos da org
			if role.ProjectId == nil {
				hasProjectAccess = true
				break
			}
			// Se a role tem project_id e é igual ao projeto solicitado
			if role.ProjectId != nil && role.ProjectId.String() == projId {
				hasProjectAccess = true
				break
			}
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
		{"/admin/login", "POST"},   // Login de admin
		{"/client/login", "POST"},  // Login de cliente
		{"/tenant/resolve", "GET"}, // Resolver tenant
		// {"/user", "POST"} - Removido: agora é tratado no middleware para capturar headers
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
