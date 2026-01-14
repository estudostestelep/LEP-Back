package middleware

import (
	"fmt"
	"lep/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PackageLimitMiddleware verifica se a organização está dentro do limite
// antes de permitir a criação de novos recursos
func PackageLimitMiddleware(limitHandler *handler.LimitHandler, limitType handler.LimitType) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgId := c.GetString("organization_id")
		projectId := c.GetString("project_id")

		if orgId == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Organização não identificada",
				"message": "É necessário fornecer o ID da organização",
			})
			c.Abort()
			return
		}

		// Verificar se é master admin (bypass de limites)
		if isMasterAdmin(c) {
			fmt.Printf("🔓 Master admin bypass de limite para %s\n", limitType)
			c.Next()
			return
		}

		// Verificar limite
		canCreate, current, limit, err := limitHandler.CheckLimit(orgId, projectId, limitType)
		if err != nil {
			// Se não tem plano ativo, bloquear
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Plano não encontrado",
				"message": err.Error(),
				"code":    "NO_ACTIVE_PLAN",
			})
			c.Abort()
			return
		}

		// Se não pode criar, retornar erro com detalhes
		if !canCreate {
			var message string
			if limit == 0 {
				message = fmt.Sprintf("Funcionalidade de %s não está disponível no seu plano atual. Faça upgrade para acessar.", getLimitDisplayName(limitType))
			} else {
				message = fmt.Sprintf("Limite de %s atingido (%d/%d). Faça upgrade do seu plano para criar mais.", getLimitDisplayName(limitType), current, limit)
			}

			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Limite atingido",
				"message": message,
				"code":    "LIMIT_EXCEEDED",
				"details": gin.H{
					"limit_type": limitType,
					"current":    current,
					"limit":      limit,
				},
			})
			c.Abort()
			return
		}

		// Limite OK, continuar
		c.Next()
	}
}

// ModuleRequiredMiddleware verifica se a organização tem acesso a um módulo específico
func ModuleRequiredMiddleware(limitHandler *handler.LimitHandler, moduleCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgId := c.GetString("organization_id")

		if orgId == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Organização não identificada",
				"message": "É necessário fornecer o ID da organização",
			})
			c.Abort()
			return
		}

		// Verificar se é master admin (bypass)
		if isMasterAdmin(c) {
			c.Next()
			return
		}

		// Verificar se tem acesso ao módulo
		hasModule, err := limitHandler.HasModule(orgId, moduleCode)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Plano não encontrado",
				"message": err.Error(),
				"code":    "NO_ACTIVE_PLAN",
			})
			c.Abort()
			return
		}

		if !hasModule {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Módulo não disponível",
				"message": fmt.Sprintf("O módulo '%s' não está disponível no seu plano atual. Faça upgrade para acessar.", moduleCode),
				"code":    "MODULE_NOT_AVAILABLE",
				"details": gin.H{
					"module_code": moduleCode,
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// getLimitDisplayName retorna o nome amigável do tipo de limite
func getLimitDisplayName(limitType handler.LimitType) string {
	switch limitType {
	case handler.LimitTables:
		return "mesas"
	case handler.LimitUsers:
		return "usuários"
	case handler.LimitProducts:
		return "produtos"
	case handler.LimitReservationsDay:
		return "reservas por dia"
	default:
		return string(limitType)
	}
}
