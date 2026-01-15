package middleware

import (
	"fmt"
	"lep/handler"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// RolePermissionMiddleware verifica se o usuário tem a permissão necessária (sistema de níveis)
// permissionCode: código da permissão (ex: "client_orders_view")
// minLevel: nível mínimo requerido (1 = visualizar, 2 = editar/CRUD completo)
func RolePermissionMiddleware(roleHandler *handler.RoleHandler, permissionCode string, minLevel int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString("user_id")
		orgId := c.GetString("organization_id")

		if userId == "" || orgId == "" {
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

		// Verificação completa: módulo + permissão + hierarquia
		hasPermission, reason, err := roleHandler.FullPermissionCheck(userId, orgId, permissionCode, minLevel, "")
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
				"error":   "Acesso negado",
				"message": reason,
				"code":    "PERMISSION_DENIED",
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

// GetUserHierarchyLevel retorna o nível de hierarquia do usuário atual
func GetUserHierarchyLevel(c *gin.Context, roleHandler *handler.RoleHandler) (int, error) {
	userId := c.GetString("user_id")
	orgId := c.GetString("organization_id")

	if userId == "" || orgId == "" {
		return 0, fmt.Errorf("usuário ou organização não identificados")
	}

	// Master admin tem nível máximo
	if isMasterAdmin(c) {
		return 10, nil
	}

	return roleHandler.GetUserMaxHierarchyLevel(userId, orgId)
}
