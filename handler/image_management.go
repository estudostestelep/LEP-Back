package handler

import (
	"context"
	"lep/service"

	"github.com/google/uuid"
)

// resourceImageManagement gerencia operações de imagens (cleanup, stats)
type resourceImageManagement struct {
	svc service.IImageManagementService
}

// IHandlerImageManagement define operações de gerenciamento de imagens
type IHandlerImageManagement interface {
	// Deletar referência de imagem
	DeleteImageReference(entityType string, entityId uuid.UUID, entityField string) (*service.DeleteImageResponse, error)

	// Cleanup de arquivos órfãos
	CleanupOrphanedFiles(olderThanDays int) (*service.CleanupResponse, error)

	// Obter estatísticas de imagens
	GetImageStats(orgId, projId uuid.UUID) (*service.ImageStatsResponse, error)
}

// NewHandlerImageManagement cria nova instância do handler
func NewHandlerImageManagement(svc service.IImageManagementService) IHandlerImageManagement {
	return &resourceImageManagement{
		svc: svc,
	}
}

// DeleteImageReference deleta referência de imagem
func (h *resourceImageManagement) DeleteImageReference(entityType string, entityId uuid.UUID, entityField string) (*service.DeleteImageResponse, error) {
	ctx := context.Background()
	return h.svc.DeleteImageReference(ctx, entityType, entityId, entityField)
}

// CleanupOrphanedFiles executa cleanup de arquivos órfãos
func (h *resourceImageManagement) CleanupOrphanedFiles(olderThanDays int) (*service.CleanupResponse, error) {
	ctx := context.Background()

	// Se olderThanDays não informado, usar padrão 0 (deleta imediatamente)
	if olderThanDays < 0 {
		olderThanDays = 0
	}

	return h.svc.CleanupOrphanedFiles(ctx, olderThanDays)
}

// GetImageStats retorna estatísticas de imagens
func (h *resourceImageManagement) GetImageStats(orgId, projId uuid.UUID) (*service.ImageStatsResponse, error) {
	ctx := context.Background()
	return h.svc.GetImageStats(ctx, orgId, projId)
}
