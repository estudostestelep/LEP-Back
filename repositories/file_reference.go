package repositories

import (
	"context"
	"fmt"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// IFileReferenceRepository define operações para gerenciar referências de arquivos
type IFileReferenceRepository interface {
	// Buscar por hash (deduplicação)
	GetByHash(ctx context.Context, orgId, projId, fileHash string) (*models.FileReference, error)

	// Registrar novo arquivo
	Create(ctx context.Context, fileRef *models.FileReference) error

	// Buscar por ID
	GetByID(ctx context.Context, fileId uuid.UUID) (*models.FileReference, error)

	// Incrementar contador de referências
	IncrementReferenceCount(ctx context.Context, fileId uuid.UUID) error

	// Decrementar contador de referências
	DecrementReferenceCount(ctx context.Context, fileId uuid.UUID) error

	// Soft delete (marcar como deletado)
	SoftDelete(ctx context.Context, fileId uuid.UUID) error

	// Buscar arquivos órfãos (reference_count = 0, deleted_at NOT NULL)
	GetOrphanedFiles(ctx context.Context, olderThanDays int) ([]models.FileReference, error)

	// Hard delete (deletar do banco)
	HardDelete(ctx context.Context, fileId uuid.UUID) error

	// Atualizar last_accessed_at
	UpdateLastAccessed(ctx context.Context, fileId uuid.UUID) error

	// Listar arquivos por categoria
	ListByCategory(ctx context.Context, orgId, projId, category string) ([]models.FileReference, error)
}

// fileReferenceRepository implementa IFileReferenceRepository
type fileReferenceRepository struct {
	db *gorm.DB
}

// NewFileReferenceRepository cria nova instância do repository
func NewFileReferenceRepository(db *gorm.DB) IFileReferenceRepository {
	return &fileReferenceRepository{db: db}
}

// GetByHash busca arquivo por hash (para deduplicação)
func (r *fileReferenceRepository) GetByHash(ctx context.Context, orgId, projId, fileHash string) (*models.FileReference, error) {
	var fileRef models.FileReference

	query := r.db.WithContext(ctx).
		Where("organization_id = ? AND project_id = ? AND file_hash = ? AND deleted_at IS NULL", orgId, projId, fileHash).
		First(&fileRef)

	if query.Error != nil {
		if query.Error == gorm.ErrRecordNotFound {
			return nil, nil // Não é erro, apenas não encontrado
		}
		return nil, fmt.Errorf("erro ao buscar arquivo por hash: %w", query.Error)
	}

	return &fileRef, nil
}

// Create registra novo arquivo
func (r *fileReferenceRepository) Create(ctx context.Context, fileRef *models.FileReference) error {
	if fileRef.Id == uuid.Nil {
		fileRef.Id = uuid.New()
	}

	fileRef.CreatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(fileRef).Error; err != nil {
		return fmt.Errorf("erro ao criar referência de arquivo: %w", err)
	}

	return nil
}

// GetByID busca arquivo por ID
func (r *fileReferenceRepository) GetByID(ctx context.Context, fileId uuid.UUID) (*models.FileReference, error) {
	var fileRef models.FileReference

	if err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", fileId).
		First(&fileRef).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar arquivo por ID: %w", err)
	}

	return &fileRef, nil
}

// IncrementReferenceCount incrementa contador de referências
func (r *fileReferenceRepository) IncrementReferenceCount(ctx context.Context, fileId uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Model(&models.FileReference{}).
		Where("id = ?", fileId).
		Update("reference_count", gorm.Expr("reference_count + ?", 1)).Error; err != nil {
		return fmt.Errorf("erro ao incrementar reference_count: %w", err)
	}

	return nil
}

// DecrementReferenceCount decrementa contador de referências
func (r *fileReferenceRepository) DecrementReferenceCount(ctx context.Context, fileId uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Model(&models.FileReference{}).
		Where("id = ?", fileId).
		Update("reference_count", gorm.Expr("GREATEST(reference_count - 1, 0)")).Error; err != nil {
		return fmt.Errorf("erro ao decrementar reference_count: %w", err)
	}

	return nil
}

// SoftDelete marca arquivo como deletado (soft delete)
func (r *fileReferenceRepository) SoftDelete(ctx context.Context, fileId uuid.UUID) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Model(&models.FileReference{}).
		Where("id = ?", fileId).
		Update("deleted_at", now).Error; err != nil {
		return fmt.Errorf("erro ao soft delete arquivo: %w", err)
	}

	return nil
}

// GetOrphanedFiles busca arquivos órfãos (reference_count = 0, deleted_at NOT NULL)
func (r *fileReferenceRepository) GetOrphanedFiles(ctx context.Context, olderThanDays int) ([]models.FileReference, error) {
	var files []models.FileReference

	cutoffDate := time.Now().AddDate(0, 0, -olderThanDays)

	query := r.db.WithContext(ctx).
		Where("reference_count = 0 AND deleted_at IS NOT NULL AND deleted_at < ?", cutoffDate).
		Order("deleted_at ASC").
		Find(&files)

	if query.Error != nil {
		return nil, fmt.Errorf("erro ao buscar arquivos órfãos: %w", query.Error)
	}

	return files, nil
}

// HardDelete deleta arquivo do banco (após soft delete)
func (r *fileReferenceRepository) HardDelete(ctx context.Context, fileId uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Unscoped(). // Force hard delete
		Where("id = ?", fileId).
		Delete(&models.FileReference{}).Error; err != nil {
		return fmt.Errorf("erro ao hard delete arquivo: %w", err)
	}

	return nil
}

// UpdateLastAccessed atualiza timestamp de último acesso
func (r *fileReferenceRepository) UpdateLastAccessed(ctx context.Context, fileId uuid.UUID) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Model(&models.FileReference{}).
		Where("id = ?", fileId).
		Update("last_accessed_at", now).Error; err != nil {
		return fmt.Errorf("erro ao atualizar last_accessed_at: %w", err)
	}

	return nil
}

// ListByCategory lista arquivos de uma categoria específica
func (r *fileReferenceRepository) ListByCategory(ctx context.Context, orgId, projId, category string) ([]models.FileReference, error) {
	var files []models.FileReference

	query := r.db.WithContext(ctx).
		Where("organization_id = ? AND project_id = ? AND category = ? AND deleted_at IS NULL", orgId, projId, category).
		Order("created_at DESC").
		Find(&files)

	if query.Error != nil {
		return nil, fmt.Errorf("erro ao listar arquivos por categoria: %w", query.Error)
	}

	return files, nil
}
