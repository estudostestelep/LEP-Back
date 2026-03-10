package server

import (
	"fmt"
	"lep/handler"
	"lep/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ResourceImageManagement gerencia endpoints de imagens
type ResourceImageManagement struct {
	handlerImageManagement handler.IHandlerImageManagement
}

// IServerImageManagement define endpoints de gerenciamento de imagens
type IServerImageManagement interface {
	ServiceCleanupOrphanedFiles(c *gin.Context)
	ServiceGetImageStats(c *gin.Context)
}

// NewServerImageManagement cria nova instância do servidor
func NewServerImageManagement(h handler.IHandlerImageManagement) IServerImageManagement {
	return &ResourceImageManagement{
		handlerImageManagement: h,
	}
}

// ServiceCleanupOrphanedFiles limpa arquivos órfãos (rota: POST /admin/images/cleanup)
func (r *ResourceImageManagement) ServiceCleanupOrphanedFiles(c *gin.Context) {
	// Parâmetro opcional: dias (default = 0, deleta imediatamente)
	daysStr := c.DefaultQuery("days", "0")
	var days int
	_, err := fmt.Sscanf(daysStr, "%d", &days)
	if err != nil {
		days = 0
	}

	// Executar cleanup
	response, err := r.handlerImageManagement.CleanupOrphanedFiles(days)
	if err != nil {
		utils.SendInternalServerError(c, "Error cleaning up orphaned files", err)
		return
	}

	// Retornar resposta
	utils.SendOKSuccess(c, "Cleanup completed", gin.H{
		"success":       response.Success,
		"files_deleted": response.FilesDeleted,
		"disk_freed":    response.DiskFreed,
		"disk_freed_mb": float64(response.DiskFreed) / (1024 * 1024),
		"error_count":   response.ErrorCount,
		"message":       response.Message,
	})
}

// ServiceGetImageStats retorna estatísticas de imagens (rota: GET /admin/images/stats)
func (r *ResourceImageManagement) ServiceGetImageStats(c *gin.Context) {
	// Headers validados pelo middleware
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	if organizationId == "" || projectId == "" {
		utils.SendBadRequestError(c, "Organization ID and Project ID are required", nil)
		return
	}

	// Parsear UUIDs
	orgUUID, err := uuid.Parse(organizationId)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid organization ID", err)
		return
	}

	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid project ID", err)
		return
	}

	// Obter estatísticas
	stats, err := r.handlerImageManagement.GetImageStats(orgUUID, projUUID)
	if err != nil {
		utils.SendInternalServerError(c, "Error fetching image statistics", err)
		return
	}

	// Retornar resposta
	utils.SendOKSuccess(c, "Image statistics retrieved", gin.H{
		"total_files":           stats.TotalFiles,
		"unique_files":          stats.UniqueFiles,
		"total_references":      stats.TotalReferences,
		"duplicated_references": stats.DuplicatedReferences,
		"total_disk_usage":      stats.TotalDiskUsage,
		"total_disk_usage_mb":   float64(stats.TotalDiskUsage) / (1024 * 1024),
		"estimated_savings":     stats.EstimatedSavings,
		"estimated_savings_mb":  float64(stats.EstimatedSavings) / (1024 * 1024),
	})
}
