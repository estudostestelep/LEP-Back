package middleware

import (
	"fmt"
	"lep/handler"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// RolePermissionMiddleware verifica se o usuário tem a permissão necessária
// permissionCode: código da permissão no formato module:action (ex: "orders:read", "orders:create")
// minLevel: deprecated - mantido para compatibilidade, não é mais usado
func RolePermissionMiddleware(roleHandler *handler.RoleHandler, permissionCode string, minLevel int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString("user_id")
		userType := c.GetString("user_type")

		if userId == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Usuário não autenticado",
				"message": "É necessário estar autenticado para acessar este recurso",
			})
			c.Abort()
			return
		}

		// Verificar se é master_admin (bypass de permissões)
		if isMasterAdmin(c) {
			c.Next()
			return
		}

		// Se userType vazio, usar "client" como padrão
		if userType == "" {
			userType = "client"
		}

		// Verificação de permissão
		hasPermission, err := roleHandler.UserHasPermission(userId, userType, permissionCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Erro ao verificar permissão",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"error":      "Acesso negado",
				"message":    fmt.Sprintf("Permissão '%s' não encontrada", permissionCode),
				"code":       "PERMISSION_DENIED",
				"permission": permissionCode,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ModuleAccessMiddleware verifica se a organização tem acesso ao módulo
func ModuleAccessMiddleware(roleHandler *handler.RoleHandler, moduleCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgId := c.GetString("organization_id")

		if orgId == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Organização não identificada",
				"message": "É necessário fornecer o ID da organização",
			})
			c.Abort()
			return
		}

		// Master admin tem acesso a todos os módulos
		if isMasterAdmin(c) {
			c.Next()
			return
		}

		hasAccess, err := roleHandler.HasModuleAccess(orgId, moduleCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Erro ao verificar acesso ao módulo",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		if !hasAccess {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Módulo não disponível",
				"message": fmt.Sprintf("O módulo '%s' não está incluído no seu plano atual", moduleCode),
				"code":    "MODULE_NOT_AVAILABLE",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// HierarchyMiddleware verifica se o usuário pode gerenciar outro usuário
// O ID do usuário alvo deve estar disponível via parâmetro "userId" ou "targetUserId"
func HierarchyMiddleware(roleHandler *handler.RoleHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorUserId := c.GetString("user_id")
		orgId := c.GetString("organization_id")

		// Tentar obter o ID do usuário alvo de diferentes fontes
		targetUserId := c.Param("userId")
		if targetUserId == "" {
			targetUserId = c.Param("targetUserId")
		}
		if targetUserId == "" {
			targetUserId = c.Query("userId")
		}

		// Se não há usuário alvo, permitir (pode ser uma operação sobre si mesmo)
		if targetUserId == "" || targetUserId == actorUserId {
			c.Next()
			return
		}

		// Master admin pode gerenciar qualquer usuário
		if isMasterAdmin(c) {
			c.Next()
			return
		}

		canManage, err := roleHandler.CanManageUser(actorUserId, targetUserId, orgId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Erro ao verificar hierarquia",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		if !canManage {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Hierarquia insuficiente",
				"message": "Você não pode gerenciar usuários com nível de hierarquia maior que o seu",
				"code":    "HIERARCHY_DENIED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermissionLevel é um helper que cria um middleware para verificar nível específico
func RequirePermissionLevel(roleHandler *handler.RoleHandler, permissionCode string, level int) gin.HandlerFunc {
	return RolePermissionMiddleware(roleHandler, permissionCode, level)
}

// RequireView verifica se o usuário tem pelo menos permissão de visualização (nível 1)
func RequireView(roleHandler *handler.RoleHandler, permissionCode string) gin.HandlerFunc {
	return RolePermissionMiddleware(roleHandler, permissionCode, 1)
}

// RequireEdit verifica se o usuário tem permissão de edição completa (nível 2)
func RequireEdit(roleHandler *handler.RoleHandler, permissionCode string) gin.HandlerFunc {
	return RolePermissionMiddleware(roleHandler, permissionCode, 2)
}

// isMasterAdmin verifica se o usuário é master_admin
func isMasterAdmin(c *gin.Context) bool {
	permissions, exists := c.Get("user_permissions")
	if !exists {
		return false
	}

	// Tentar converter para []string primeiro
	if permArray, ok := permissions.([]string); ok {
		for _, perm := range permArray {
			if perm == "master_admin" {
				return true
			}
		}
		return false
	}

	// Tentar converter para pq.StringArray (tipo usado pelo GORM com PostgreSQL)
	if permArray, ok := permissions.(pq.StringArray); ok {
		for _, perm := range permArray {
			if perm == "master_admin" {
				return true
			}
		}
		return false
	}

	// Fallback: usar fmt.Sprintf e strings.Contains
	permStr := fmt.Sprintf("%v", permissions)
	return strings.Contains(permStr, "master_admin")
}

// AdminScopeMiddleware verifica se o usuário tem escopo admin
// Isso significa que é um usuário do tipo "admin" ou tem permissão master_admin
func AdminScopeMiddleware(authHandler handler.IHandlerAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType := c.GetString("user_type")

		// Verificar se é usuário admin
		if userType == "admin" {
			c.Next()
			return
		}

		// Verificar se tem master_admin (super admin pode acessar rotas admin)
		if isMasterAdmin(c) {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Acesso negado",
			"message": "Esta área é restrita para administradores do sistema",
			"code":    "ADMIN_SCOPE_REQUIRED",
		})
		c.Abort()
	}
}

// GetUserHierarchyLevel retorna o nível de hierarquia do usuário atual
func GetUserHierarchyLevel(c *gin.Context, roleHandler *handler.RoleHandler) (int, error) {
	userId := c.GetString("user_id")

	if userId == "" {
		return 0, fmt.Errorf("usuário não identificado")
	}

	// Master admin tem nível máximo
	if isMasterAdmin(c) {
		return 10, nil
	}

	userType := c.GetString("user_type")
	if userType == "" {
		userType = "client"
	}
	return roleHandler.GetUserHierarchyLevel(userId, userType)
}
