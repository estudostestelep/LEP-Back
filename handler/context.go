package handler

import (
	"lep/constants"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// RequestContext contém informações do contexto da requisição
// Usado para passar dados do servidor para o handler de forma padronizada
type RequestContext struct {
	UserId         uuid.UUID
	UserEmail      string
	OrganizationId *uuid.UUID
	ProjectId      *uuid.UUID
	Permissions    pq.StringArray
	IpAddress      string
	UserAgent      string
	IsAdminZone    bool
}

// IsMasterAdmin verifica se o usuário do contexto é um Master Admin
func (ctx *RequestContext) IsMasterAdmin() bool {
	return constants.HasPermission(ctx.Permissions, constants.PermissionMasterAdmin)
}

// BuildRequestContext constrói um RequestContext a partir do gin.Context
// Esta função extrai todas as informações relevantes do contexto da requisição
func BuildRequestContext(c *gin.Context) *RequestContext {
	ctx := &RequestContext{
		IpAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}

	// Extrair user_id do contexto (definido pelo AuthMiddleware)
	if userIdStr := c.GetString("user_id"); userIdStr != "" {
		if parsed, err := uuid.Parse(userIdStr); err == nil {
			ctx.UserId = parsed
		}
	}

	// Extrair email do usuário
	ctx.UserEmail = c.GetString("user_email")

	// Extrair permissões do usuário
	if permissions, exists := c.Get("user_permissions"); exists {
		if permList, ok := permissions.(pq.StringArray); ok {
			ctx.Permissions = permList
		}
	}

	// Determinar se está na zona admin baseado na URL
	path := c.Request.URL.Path
	ctx.IsAdminZone = strings.HasPrefix(path, "/admin")

	// Se não está na zona admin, extrair org e project dos headers/contexto
	if !ctx.IsAdminZone {
		// Tentar pegar do contexto (definido pelo HeaderValidationMiddleware)
		if orgIdStr := c.GetString("organization_id"); orgIdStr != "" {
			if parsed, err := uuid.Parse(orgIdStr); err == nil {
				ctx.OrganizationId = &parsed
			}
		}
		if projIdStr := c.GetString("project_id"); projIdStr != "" {
			if parsed, err := uuid.Parse(projIdStr); err == nil {
				ctx.ProjectId = &parsed
			}
		}
	}

	return ctx
}

// GetOrganizationIdString retorna o ID da organização como string ou string vazia
func (ctx *RequestContext) GetOrganizationIdString() string {
	if ctx.OrganizationId != nil {
		return ctx.OrganizationId.String()
	}
	return ""
}

// GetProjectIdString retorna o ID do projeto como string ou string vazia
func (ctx *RequestContext) GetProjectIdString() string {
	if ctx.ProjectId != nil {
		return ctx.ProjectId.String()
	}
	return ""
}
