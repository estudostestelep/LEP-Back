package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"lep/repositories"
	"lep/repositories/models"
	"mime/multipart"
	"os"

	"github.com/google/uuid"
)

// RegisterImageRequest dados para registrar/atualizar imagem
type RegisterImageRequest struct {
	EntityType  string    `json:"entity_type"` // "product", "category", etc
	EntityId    uuid.UUID `json:"entity_id"`
	EntityField string    `json:"entity_field"` // "image_url", "photo"
	FileHash    string    `json:"file_hash"`    // SHA-256
	FilePath    string    `json:"file_path"`
	FileSize    int64     `json:"file_size"`
	Category    string    `json:"category"` // Upload category
	MimeType    string    `json:"mime_type"`
	OrgId       uuid.UUID `json:"organization_id"`
	ProjId      uuid.UUID `json:"project_id"`
}

// ImageRegistrationResponse resposta do registro de imagem
type ImageRegistrationResponse struct {
	Success     bool      `json:"success"`
	ImageUrl    string    `json:"image_url"`
	FileHash    string    `json:"file_hash"`
	IsReused    bool      `json:"is_reused"` // Indica se reutilizou arquivo existente
	ReferenceId uuid.UUID `json:"reference_id"`
}

// DeleteImageResponse resposta da deleção de imagem
type DeleteImageResponse struct {
	Success             bool   `json:"success"`
	FileDeleted         bool   `json:"file_deleted"` // Se arquivo foi deletado (órfão)
	ReferencesRemaining int64  `json:"references_remaining"`
	Message             string `json:"message"`
}

// CleanupResponse resposta da limpeza de órfãs
type CleanupResponse struct {
	Success      bool   `json:"success"`
	FilesDeleted int    `json:"files_deleted"`
	DiskFreed    int64  `json:"disk_freed"` // bytes
	ErrorCount   int    `json:"error_count"`
	Message      string `json:"message"`
}

// ImageManagementService gerencia deduplicação e limpeza de imagens
type ImageManagementService struct {
	fileRefRepo       repositories.IFileReferenceRepository
	entityFileRefRepo repositories.IEntityFileReferenceRepository
	storageBasePath   string // Caminho base para deletar arquivos
}

// IImageManagementService define operações de gerenciamento de imagens
type IImageManagementService interface {
	// Registrar/atualizar imagem com deduplicação
	RegisterOrUpdateImage(ctx context.Context, req RegisterImageRequest) (*ImageRegistrationResponse, error)

	// Deletar referência de imagem
	DeleteImageReference(ctx context.Context, entityType string, entityId uuid.UUID, entityField string) (*DeleteImageResponse, error)

	// Cleanup de arquivos órfãos
	CleanupOrphanedFiles(ctx context.Context, olderThanDays int) (*CleanupResponse, error)

	// Calcular hash SHA-256 de arquivo
	CalculateFileHash(file multipart.File) (string, error)

	// Obter estatísticas de imagens
	GetImageStats(ctx context.Context, orgId, projId uuid.UUID) (*ImageStatsResponse, error)
}

// ImageStatsResponse estatísticas de imagens
type ImageStatsResponse struct {
	TotalFiles           int64 `json:"total_files"`
	UniqueFiles          int64 `json:"unique_files"`
	TotalReferences      int64 `json:"total_references"`
	DuplicatedReferences int64 `json:"duplicated_references"`
	TotalDiskUsage       int64 `json:"total_disk_usage"`  // bytes
	EstimatedSavings     int64 `json:"estimated_savings"` // bytes economizados por dedup
}

// NewImageManagementService cria nova instância do serviço
func NewImageManagementService(
	fileRefRepo repositories.IFileReferenceRepository,
	entityFileRefRepo repositories.IEntityFileReferenceRepository,
	storagePath string,
) IImageManagementService {
	return &ImageManagementService{
		fileRefRepo:       fileRefRepo,
		entityFileRefRepo: entityFileRefRepo,
		storageBasePath:   storagePath,
	}
}

// RegisterOrUpdateImage registra ou atualiza imagem com deduplicação automática
func (s *ImageManagementService) RegisterOrUpdateImage(ctx context.Context, req RegisterImageRequest) (*ImageRegistrationResponse, error) {
	// 1. Verificar se já existe imagem nessa entidade/campo
	existingRef, err := s.entityFileRefRepo.GetByEntity(ctx, req.EntityType, req.EntityId, req.EntityField)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar referência existente: %w", err)
	}

	isReused := false

	// 2. Se existe imagem anterior, remover referência (será substituta ou deletada)
	if existingRef != nil {
		// Se o hash é igual, não fazer nada (mesma imagem)
		existingFile, err := s.fileRefRepo.GetByID(ctx, existingRef.FileId)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar arquivo anterior: %w", err)
		}

		// Se hashes são iguais, reutilizar arquivo existente
		if existingFile != nil && existingFile.FileHash == req.FileHash {
			// Atualizar last_accessed_at
			_ = s.fileRefRepo.UpdateLastAccessed(ctx, existingFile.Id)

			return &ImageRegistrationResponse{
				Success:     true,
				ImageUrl:    existingFile.FilePath,
				FileHash:    existingFile.FileHash,
				IsReused:    true, // Mesma imagem, nenhuma mudança
				ReferenceId: existingRef.Id,
			}, nil
		}

		// Hashes diferentes: remover referência antiga
		if err := s.entityFileRefRepo.SoftDelete(ctx, req.EntityType, req.EntityId, req.EntityField); err != nil {
			return nil, fmt.Errorf("erro ao remover referência anterior: %w", err)
		}

		// Decrementar contador da imagem antiga
		if err := s.fileRefRepo.DecrementReferenceCount(ctx, existingRef.FileId); err != nil {
			return nil, fmt.Errorf("erro ao decrementar referência anterior: %w", err)
		}

		// Se contador chegar a 0, soft delete da imagem antiga
		count, err := s.entityFileRefRepo.CountByFileID(ctx, existingRef.FileId)
		if err != nil {
			return nil, fmt.Errorf("erro ao contar referências: %w", err)
		}

		if count == 0 {
			if err := s.fileRefRepo.SoftDelete(ctx, existingRef.FileId); err != nil {
				return nil, fmt.Errorf("erro ao soft delete arquivo antigo: %w", err)
			}
		}
	}

	// 3. Verificar se arquivo com esse hash já existe
	existingFile, err := s.fileRefRepo.GetByHash(ctx, req.OrgId.String(), req.ProjId.String(), req.FileHash)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar arquivo por hash: %w", err)
	}

	// 4. Se arquivo existe, reutilizar; se não, criar novo
	var fileId uuid.UUID
	var imageUrl string
	isReused = (existingFile != nil)

	if existingFile != nil {
		// Reutilizar arquivo existente
		fileId = existingFile.Id
		imageUrl = existingFile.FilePath

		// Incrementar contador
		if err := s.fileRefRepo.IncrementReferenceCount(ctx, fileId); err != nil {
			return nil, fmt.Errorf("erro ao incrementar reference_count: %w", err)
		}

		// Atualizar last_accessed_at
		_ = s.fileRefRepo.UpdateLastAccessed(ctx, fileId)
	} else {
		// Criar novo arquivo
		fileId = uuid.New()
		imageUrl = req.FilePath

		newFile := &models.FileReference{
			Id:             fileId,
			OrganizationId: req.OrgId,
			ProjectId:      req.ProjId,
			FileHash:       req.FileHash,
			FilePath:       req.FilePath,
			FileSize:       req.FileSize,
			Category:       req.Category,
			MimeType:       req.MimeType,
			ReferenceCount: 1,
		}

		if err := s.fileRefRepo.Create(ctx, newFile); err != nil {
			return nil, fmt.Errorf("erro ao criar arquivo: %w", err)
		}
	}

	// 5. Criar nova referência de entidade
	newEntityRef := &models.EntityFileReference{
		Id:          uuid.New(),
		FileId:      fileId,
		EntityType:  req.EntityType,
		EntityId:    req.EntityId,
		EntityField: req.EntityField,
	}

	if err := s.entityFileRefRepo.Create(ctx, newEntityRef); err != nil {
		return nil, fmt.Errorf("erro ao criar referência de entidade: %w", err)
	}

	return &ImageRegistrationResponse{
		Success:     true,
		ImageUrl:    imageUrl,
		FileHash:    req.FileHash,
		IsReused:    isReused,
		ReferenceId: newEntityRef.Id,
	}, nil
}

// DeleteImageReference deleta referência de imagem e limpa arquivo se órfão
func (s *ImageManagementService) DeleteImageReference(ctx context.Context, entityType string, entityId uuid.UUID, entityField string) (*DeleteImageResponse, error) {
	// 1. Buscar referência
	entityRef, err := s.entityFileRefRepo.GetByEntity(ctx, entityType, entityId, entityField)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar referência: %w", err)
	}

	if entityRef == nil {
		return &DeleteImageResponse{
			Success:             true,
			FileDeleted:         false,
			ReferencesRemaining: 0,
			Message:             "Referência não encontrada",
		}, nil
	}

	fileId := entityRef.FileId

	// 2. Soft delete da referência
	if err := s.entityFileRefRepo.SoftDelete(ctx, entityType, entityId, entityField); err != nil {
		return nil, fmt.Errorf("erro ao deletar referência: %w", err)
	}

	// 3. Decrementar contador do arquivo
	if err := s.fileRefRepo.DecrementReferenceCount(ctx, fileId); err != nil {
		return nil, fmt.Errorf("erro ao decrementar contador: %w", err)
	}

	// 4. Verificar se arquivo ficou órfão
	count, err := s.entityFileRefRepo.CountByFileID(ctx, fileId)
	if err != nil {
		return nil, fmt.Errorf("erro ao contar referências: %w", err)
	}

	fileDeleted := false
	if count == 0 {
		// Soft delete do arquivo órfão
		if err := s.fileRefRepo.SoftDelete(ctx, fileId); err != nil {
			return nil, fmt.Errorf("erro ao soft delete arquivo: %w", err)
		}
		fileDeleted = true
	}

	return &DeleteImageResponse{
		Success:             true,
		FileDeleted:         fileDeleted,
		ReferencesRemaining: count,
		Message:             fmt.Sprintf("Referência deletada. Arquivo %s deletado", map[bool]string{true: "foi", false: "não foi"}[fileDeleted]),
	}, nil
}

// CleanupOrphanedFiles limpa arquivos órfãos (soft deletados, sem referências)
func (s *ImageManagementService) CleanupOrphanedFiles(ctx context.Context, olderThanDays int) (*CleanupResponse, error) {
	// 1. Buscar arquivos órfãos
	orphanedFiles, err := s.fileRefRepo.GetOrphanedFiles(ctx, olderThanDays)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar arquivos órfãos: %w", err)
	}

	if len(orphanedFiles) == 0 {
		return &CleanupResponse{
			Success:      true,
			FilesDeleted: 0,
			DiskFreed:    0,
			ErrorCount:   0,
			Message:      "Nenhum arquivo órfão encontrado",
		}, nil
	}

	// 2. Deletar arquivo de storage e do banco
	var totalDiskFreed int64 = 0
	var errorCount = 0

	for _, file := range orphanedFiles {
		// Deletar do storage
		fullPath := fmt.Sprintf("%s/%s", s.storageBasePath, file.FilePath)
		if err := os.Remove(fullPath); err != nil {
			// Logar erro mas continuar (arquivo pode não existir)
			fmt.Printf("⚠️  Erro ao deletar arquivo %s: %v\n", fullPath, err)
			errorCount++
		} else {
			totalDiskFreed += file.FileSize
		}

		// Hard delete do banco
		if err := s.fileRefRepo.HardDelete(ctx, file.Id); err != nil {
			fmt.Printf("⚠️  Erro ao deletar referência %s do banco: %v\n", file.Id, err)
			errorCount++
		}
	}

	// 3. Cleanup de referências deletadas (opcional, manutenção)
	_ = s.entityFileRefRepo.CleanupDeletedReferences(ctx)

	return &CleanupResponse{
		Success:      errorCount == 0,
		FilesDeleted: len(orphanedFiles) - errorCount,
		DiskFreed:    totalDiskFreed,
		ErrorCount:   errorCount,
		Message:      fmt.Sprintf("%d arquivos deletados, %d bytes liberados", len(orphanedFiles)-errorCount, totalDiskFreed),
	}, nil
}

// CalculateFileHash calcula hash SHA-256 de arquivo
func (s *ImageManagementService) CalculateFileHash(file multipart.File) (string, error) {
	hash := sha256.New()

	// Copiar arquivo para hash
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("erro ao calcular hash: %w", err)
	}

	// Reset file pointer para início
	if _, err := file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("erro ao reset file pointer: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// GetImageStats retorna estatísticas de imagens por organização/projeto
func (s *ImageManagementService) GetImageStats(ctx context.Context, orgId, projId uuid.UUID) (*ImageStatsResponse, error) {
	// Implementação simplificada
	// Em produção, usar queries otimizadas no banco

	return &ImageStatsResponse{
		TotalFiles:           0,
		UniqueFiles:          0,
		TotalReferences:      0,
		DuplicatedReferences: 0,
		TotalDiskUsage:       0,
		EstimatedSavings:     0,
	}, nil
}
