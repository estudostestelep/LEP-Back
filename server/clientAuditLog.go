package server

import (
	"lep/handler"
	"lep/repositories/models"
	"lep/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// ResourceClientAuditLog - Servidor para logs de auditoria de cliente
type ResourceClientAuditLog struct {
	handler handler.IClientAuditLogHandler
}

// IClientAuditLogServer - Interface do servidor
type IClientAuditLogServer interface {
	// Logs
	ServiceListClientAuditLogs(c *gin.Context)
	ServiceGetClientAuditLog(c *gin.Context)

	// Configuração
	ServiceGetClientAuditConfig(c *gin.Context)
	ServiceUpdateClientAuditConfig(c *gin.Context)
	ServiceGetAvailableModules(c *gin.Context)
}

// NewClientAuditLogServer - Construtor do servidor
func NewClientAuditLogServer(handler handler.IClientAuditLogHandler) IClientAuditLogServer {
	return &ResourceClientAuditLog{handler: handler}
}

// ServiceListClientAuditLogs - Lista logs de auditoria de cliente
// GET /client-audit-logs?page=1&page_size=15&start_date=2024-01-01&end_date=2024-12-31&user_email=user@test.com&action=CREATE&entity_type=reservation&module_code=reservations
func (r *ResourceClientAuditLog) ServiceListClientAuditLogs(c *gin.Context) {
	// Obter org e projeto do contexto (headers já validados pelo middleware)
	orgIdStr := c.GetString("organization_id")
	projectIdStr := c.GetString("project_id")

	orgId, err := uuid.Parse(orgIdStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid organization ID", err)
		return
	}

	projectId, err := uuid.Parse(projectIdStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid project ID", err)
		return
	}

	// Montar filtros a partir dos query params
	filters := models.ClientAuditLogFilters{
		Page:     1,
		PageSize: 15,
	}

	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			filters.Page = p
		}
	}

	if pageSize := c.Query("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 && ps <= 100 {
			filters.PageSize = ps
		}
	}

	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filters.StartDate = &t
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			// Adicionar 23:59:59 para incluir o dia inteiro
			endOfDay := t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			filters.EndDate = &endOfDay
		}
	}

	if userEmail := c.Query("user_email"); userEmail != "" {
		filters.UserEmail = userEmail
	}

	if action := c.Query("action"); action != "" {
		filters.Action = action
	}

	if entityType := c.Query("entity_type"); entityType != "" {
		filters.EntityType = entityType
	}

	if moduleCode := c.Query("module_code"); moduleCode != "" {
		filters.ModuleCode = moduleCode
	}

	// Buscar logs
	response, err := r.handler.ListLogs(orgId, projectId, filters)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing client audit logs", err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// ServiceGetClientAuditLog - Obtém detalhes de um log específico
// GET /client-audit-logs/:id
func (r *ResourceClientAuditLog) ServiceGetClientAuditLog(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid log ID format", err)
		return
	}

	log, err := r.handler.GetLogById(id)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting client audit log", err)
		return
	}

	if log == nil {
		utils.SendNotFoundError(c, "Client audit log")
		return
	}

	// Verificar se o log pertence à organização/projeto do usuário
	orgIdStr := c.GetString("organization_id")
	projectIdStr := c.GetString("project_id")

	orgId, _ := uuid.Parse(orgIdStr)
	projectId, _ := uuid.Parse(projectIdStr)

	if log.OrganizationId != orgId || log.ProjectId != projectId {
		utils.SendForbiddenError(c, "You don't have permission to view this log")
		return
	}

	c.JSON(http.StatusOK, log)
}

// ServiceGetClientAuditConfig - Obtém configuração de auditoria da organização
// GET /client-audit-config
func (r *ResourceClientAuditLog) ServiceGetClientAuditConfig(c *gin.Context) {
	orgIdStr := c.GetString("organization_id")
	orgId, err := uuid.Parse(orgIdStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid organization ID", err)
		return
	}

	config, err := r.handler.GetConfig(orgId)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting audit config", err)
		return
	}

	c.JSON(http.StatusOK, config)
}

// ServiceUpdateClientAuditConfig - Atualiza configuração de auditoria da organização
// PUT /client-audit-config
func (r *ResourceClientAuditLog) ServiceUpdateClientAuditConfig(c *gin.Context) {
	orgIdStr := c.GetString("organization_id")
	orgId, err := uuid.Parse(orgIdStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid organization ID", err)
		return
	}

	// DTO para receber a requisição
	var request struct {
		Enabled        bool     `json:"enabled"`
		MaxLogsStored  int      `json:"max_logs_stored"`
		RetentionDays  int      `json:"retention_days"`
		EnabledModules []string `json:"enabled_modules"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Validar valores
	if request.MaxLogsStored < 100 || request.MaxLogsStored > 100000 {
		utils.SendBadRequestError(c, "max_logs_stored must be between 100 and 100000", nil)
		return
	}

	if request.RetentionDays < 7 || request.RetentionDays > 365 {
		utils.SendBadRequestError(c, "retention_days must be between 7 and 365", nil)
		return
	}

	// Validar módulos
	availableModules := r.handler.GetAvailableModules()
	validModuleCodes := make(map[string]bool)
	for _, m := range availableModules {
		validModuleCodes[m.Code] = true
	}

	for _, moduleCode := range request.EnabledModules {
		if !validModuleCodes[moduleCode] {
			utils.SendBadRequestError(c, "Invalid module code: "+moduleCode, nil)
			return
		}
	}

	// Montar config
	config := &models.ClientAuditConfig{
		OrganizationId: orgId,
		Enabled:        request.Enabled,
		MaxLogsStored:  request.MaxLogsStored,
		RetentionDays:  request.RetentionDays,
		EnabledModules: pq.StringArray(request.EnabledModules),
	}

	// Salvar
	if err := r.handler.SaveConfig(config); err != nil {
		utils.SendInternalServerError(c, "Error saving audit config", err)
		return
	}

	utils.SendOKSuccess(c, "Audit configuration updated successfully", config)
}

// ServiceGetAvailableModules - Retorna módulos disponíveis para auditoria
// GET /client-audit-modules
func (r *ResourceClientAuditLog) ServiceGetAvailableModules(c *gin.Context) {
	modules := r.handler.GetAvailableModules()
	c.JSON(http.StatusOK, gin.H{
		"modules": modules,
	})
}
