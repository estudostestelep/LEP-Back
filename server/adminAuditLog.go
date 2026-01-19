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
)

// AdminAuditLogServer - Controller para logs de auditoria administrativa
// IMPORTANTE: Este server expõe apenas operações de LEITURA (read-only)
type AdminAuditLogServer struct {
	handler handler.IAdminAuditLogHandler
}

// IAdminAuditLogServer - Interface do server
type IAdminAuditLogServer interface {
	// GET /admin/audit-logs - Lista logs com filtros
	ServiceListAdminAuditLogs(c *gin.Context)
	// GET /admin/audit-logs/:id - Detalhes de um log específico
	ServiceGetAdminAuditLog(c *gin.Context)
	// DELETE /admin/audit-logs/cleanup?days=90 - Remove logs mais antigos que X dias
	ServiceDeleteOldLogs(c *gin.Context)
}

// NewAdminAuditLogServer - Construtor do server
func NewAdminAuditLogServer(h handler.IAdminAuditLogHandler) IAdminAuditLogServer {
	return &AdminAuditLogServer{handler: h}
}

// ServiceListAdminAuditLogs - Lista logs de auditoria com filtros e paginação
// GET /admin/audit-logs?start_date=2024-01-01&end_date=2024-12-31&actor_email=admin@example.com&action=UPDATE&page=1&page_size=20
func (s *AdminAuditLogServer) ServiceListAdminAuditLogs(c *gin.Context) {
	// Obter parâmetros de filtro
	filters := models.AdminAuditLogFilters{}

	// Filtro por data inicial
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			// Início do dia
			startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
			filters.StartDate = &startDate
		}
	}

	// Filtro por data final
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			// Final do dia
			endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, time.UTC)
			filters.EndDate = &endDate
		}
	}

	// Filtro por ID do ator
	if actorIdStr := c.Query("actor_id"); actorIdStr != "" {
		actorId, err := uuid.Parse(actorIdStr)
		if err == nil {
			filters.ActorId = &actorId
		}
	}

	// Filtro por email do ator
	if actorEmail := c.Query("actor_email"); actorEmail != "" {
		filters.ActorEmail = actorEmail
	}

	// Filtro por ação
	if action := c.Query("action"); action != "" {
		filters.Action = action
	}

	// Filtro por tipo de entidade
	if entityType := c.Query("entity_type"); entityType != "" {
		filters.EntityType = entityType
	}

	// Paginação
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// Validar limites
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	filters.Page = page
	filters.PageSize = pageSize

	// Buscar logs
	response, err := s.handler.ListLogs(filters)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing admin audit logs", err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// ServiceGetAdminAuditLog - Obtém detalhes de um log específico
// GET /admin/audit-logs/:id
func (s *AdminAuditLogServer) ServiceGetAdminAuditLog(c *gin.Context) {
	idStr := c.Param("id")

	// Validar formato UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid audit log ID format", err)
		return
	}

	// Buscar log
	log, err := s.handler.GetLogById(id)
	if err != nil {
		utils.SendNotFoundError(c, "Admin audit log")
		return
	}

	c.JSON(http.StatusOK, log)
}

// ServiceDeleteOldLogs - Remove logs mais antigos que X dias
// DELETE /admin/audit-logs/cleanup?days=90
func (s *AdminAuditLogServer) ServiceDeleteOldLogs(c *gin.Context) {
	// Obter parâmetro days
	daysStr := c.Query("days")
	if daysStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing required parameter",
			"message": "O parâmetro 'days' é obrigatório",
		})
		return
	}

	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid parameter",
			"message": "O parâmetro 'days' deve ser um número inteiro maior que 0",
		})
		return
	}

	// Deletar logs antigos
	deleted, err := s.handler.DeleteOlderThan(days)
	if err != nil {
		utils.SendInternalServerError(c, "Error deleting old logs", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logs antigos removidos com sucesso",
		"deleted": deleted,
		"days":    days,
	})
}
