package middleware

import (
	"lep/constants"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ==================== Interfaces ====================

// IAuthorizationHandler interface para verificações de autorização
type IAuthorizationHandler interface {
	// IsMasterAdmin verifica se usuário é master admin (hierarchy >= 10)
	IsMasterAdmin(userId, userType string) (bool, error)
	// GetUserHierarchyLevel retorna o maior nível de hierarquia do usuário
	GetUserHierarchyLevel(userId, userType string) (int, error)
	// UserHasPermission verifica se usuário tem uma permissão específica via roles
	UserHasPermission(userId, userType, permission string) (bool, error)
	// OrganizationHasModule verifica se organização tem acesso ao módulo via plan
	OrganizationHasModule(orgId, moduleCode string) (bool, error)
	// CanManageUser verifica se actor pode gerenciar target baseado em hierarquia
	CanManageUser(actorId, targetId, userType string) (bool, error)
}

// ==================== Middleware de Permissão ====================

// RequirePermission verifica se usuário tem permissão no formato module:action
// Fluxo:
// 1. Bypass se master_admin (hierarchy >= 10)
// 2. Verificar se organização tem o módulo (via Plan)
// 3. Verificar se usuário tem a permissão (via Role)
func RequirePermission(authHandler IAuthorizationHandler, perm constants.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString("user_id")
		orgId := c.GetString("organization_id")
		userType := c.GetString("user_type")

		if userId == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Usuário não autenticado",
				"code":  "UNAUTHENTICATED",
			})
			c.Abort()
			return
		}

		// 1. Bypass para master_admin
		isMaster, err := authHandler.IsMasterAdmin(userId, userType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao verificar permissões",
				"code":  "AUTH_CHECK_ERROR",
			})
			c.Abort()
			return
		}
		if isMaster {
			c.Next()
			return
		}

		module, _ := constants.ParsePermission(perm)

		// 2. Verificar se organização tem o módulo (via Plan)
		if orgId != "" && module != "" {
			hasModule, err := authHandler.OrganizationHasModule(orgId, module)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Erro ao verificar módulo",
					"code":  "MODULE_CHECK_ERROR",
				})
				c.Abort()
				return
			}
			if !hasModule {
				c.JSON(http.StatusForbidden, gin.H{
					"error":   "Módulo não disponível no seu plano",
					"code":    "MODULE_NOT_AVAILABLE",
					"module":  module,
					"upgrade": true,
				})
				c.Abort()
				return
			}
		}

		// 3. Verificar se usuário tem a permissão (via Role)
		hasPerm, err := authHandler.UserHasPermission(userId, userType, string(perm))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao verificar permissão",
				"code":  "PERMISSION_CHECK_ERROR",
			})
			c.Abort()
			return
		}
		if !hasPerm {
			c.JSON(http.StatusForbidden, gin.H{
				"error":      "Permissão negada",
				"code":       "PERMISSION_DENIED",
				"permission": string(perm),
				"module":     module,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission verifica se usuário tem pelo menos uma das permissões
func RequireAnyPermission(authHandler IAuthorizationHandler, perms ...constants.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString("user_id")
		userType := c.GetString("user_type")

		if userId == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Usuário não autenticado",
				"code":  "UNAUTHENTICATED",
			})
			c.Abort()
			return
		}

		// Bypass para master_admin
		isMaster, _ := authHandler.IsMasterAdmin(userId, userType)
		if isMaster {
			c.Next()
			return
		}

		// Verificar se tem pelo menos uma permissão
		for _, perm := range perms {
			hasPerm, _ := authHandler.UserHasPermission(userId, userType, string(perm))
			if hasPerm {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "Nenhuma das permissões requeridas encontrada",
			"code":  "PERMISSION_DENIED",
		})
		c.Abort()
	}
}

// RequireAllPermissions verifica se usuário tem todas as permissões
func RequireAllPermissions(authHandler IAuthorizationHandler, perms ...constants.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString("user_id")
		userType := c.GetString("user_type")

		if userId == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Usuário não autenticado",
				"code":  "UNAUTHENTICATED",
			})
			c.Abort()
			return
		}

		// Bypass para master_admin
		isMaster, _ := authHandler.IsMasterAdmin(userId, userType)
		if isMaster {
			c.Next()
			return
		}

		// Verificar se tem todas as permissões
		for _, perm := range perms {
			hasPerm, _ := authHandler.UserHasPermission(userId, userType, string(perm))
			if !hasPerm {
				c.JSON(http.StatusForbidden, gin.H{
					"error":      "Permissão requerida não encontrada",
					"code":       "PERMISSION_DENIED",
					"permission": string(perm),
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// ==================== Middleware de Módulo ====================

// RequireModule verifica se organização tem acesso ao módulo
func RequireModule(authHandler IAuthorizationHandler, moduleCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString("user_id")
		orgId := c.GetString("organization_id")
		userType := c.GetString("user_type")

		if userId == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Usuário não autenticado",
				"code":  "UNAUTHENTICATED",
			})
			c.Abort()
			return
		}

		// Bypass para master_admin
		isMaster, _ := authHandler.IsMasterAdmin(userId, userType)
		if isMaster {
			c.Next()
			return
		}

		if orgId == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Organização não especificada",
				"code":  "ORG_REQUIRED",
			})
			c.Abort()
			return
		}

		hasModule, err := authHandler.OrganizationHasModule(orgId, moduleCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao verificar módulo",
				"code":  "MODULE_CHECK_ERROR",
			})
			c.Abort()
			return
		}
		if !hasModule {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Módulo não disponível no seu plano",
				"code":    "MODULE_NOT_AVAILABLE",
				"module":  moduleCode,
				"upgrade": true,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ==================== Middleware de Hierarquia ====================

// RequireHierarchy verifica se usuário pode gerenciar outro usuário baseado em hierarquia
func RequireHierarchy(authHandler IAuthorizationHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorId := c.GetString("user_id")
		userType := c.GetString("user_type")

		// Tentar obter ID do usuário alvo de várias fontes
		targetId := c.Param("userId")
		if targetId == "" {
			targetId = c.Param("id")
		}
		if targetId == "" {
			targetId = c.Query("userId")
		}

		// Se não há usuário alvo ou é o próprio usuário, permitir
		if targetId == "" || targetId == actorId {
			c.Next()
			return
		}

		// Master admin pode gerenciar qualquer usuário
		isMaster, _ := authHandler.IsMasterAdmin(actorId, userType)
		if isMaster {
			c.Next()
			return
		}

		canManage, err := authHandler.CanManageUser(actorId, targetId, userType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao verificar hierarquia",
				"code":  "HIERARCHY_CHECK_ERROR",
			})
			c.Abort()
			return
		}
		if !canManage {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Hierarquia insuficiente para gerenciar este usuário",
				"code":  "HIERARCHY_DENIED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireMinHierarchy verifica se usuário tem nível mínimo de hierarquia
func RequireMinHierarchy(authHandler IAuthorizationHandler, minLevel int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString("user_id")
		userType := c.GetString("user_type")

		if userId == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Usuário não autenticado",
				"code":  "UNAUTHENTICATED",
			})
			c.Abort()
			return
		}

		level, err := authHandler.GetUserHierarchyLevel(userId, userType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao verificar hierarquia",
				"code":  "HIERARCHY_CHECK_ERROR",
			})
			c.Abort()
			return
		}

		if level < minLevel {
			c.JSON(http.StatusForbidden, gin.H{
				"error":         "Nível de hierarquia insuficiente",
				"code":          "HIERARCHY_DENIED",
				"required":      minLevel,
				"current_level": level,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ==================== Middleware de Admin ====================

// RequireMasterAdmin permite apenas usuários com hierarchy >= 10
func RequireMasterAdmin(authHandler IAuthorizationHandler) gin.HandlerFunc {
	return RequireMinHierarchy(authHandler, constants.HierarchyMasterAdmin)
}

// RequireAdminScope permite usuários com roles de escopo "admin"
func RequireAdminScope(authHandler IAuthorizationHandler) gin.HandlerFunc {
	return RequireMinHierarchy(authHandler, constants.HierarchyAdmin)
}

// ==================== Helpers para Contexto ====================

// GetUserIdFromContext extrai user_id do contexto
func GetUserIdFromContext(c *gin.Context) string {
	return c.GetString("user_id")
}

// GetOrgIdFromContext extrai organization_id do contexto
func GetOrgIdFromContext(c *gin.Context) string {
	return c.GetString("organization_id")
}

// GetUserTypeFromContext extrai user_type do contexto
func GetUserTypeFromContext(c *gin.Context) string {
	return c.GetString("user_type")
}

// IsMasterAdminFromContext verifica se usuário no contexto é master admin
func IsMasterAdminFromContext(c *gin.Context, authHandler IAuthorizationHandler) bool {
	userId := c.GetString("user_id")
	userType := c.GetString("user_type")
	if userId == "" {
		return false
	}
	isMaster, _ := authHandler.IsMasterAdmin(userId, userType)
	return isMaster
}
