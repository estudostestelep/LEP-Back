package repositories

import (
	"context"
	"fmt"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// IEntityFileReferenceRepository define operações para gerenciar referências de entidades a arquivos
type IEntityFileReferenceRepository interface {
	// Registrar referência entre entidade e arquivo
	Create(ctx context.Context, entityRef *models.EntityFileReference) error

	// Buscar referência de uma entidade específica
	GetByEntity(ctx context.Context, entityType string, entityId uuid.UUID, entityField string) (*models.EntityFileReference, error)

	// Deletar referência (soft delete)
	SoftDelete(ctx context.Context, entityType string, entityId uuid.UUID, entityField string) error

	// Hard delete (após soft delete)
	HardDelete(ctx context.Context, entityType string, entityId uuid.UUID, entityField string) error

	// Contar quantas entidades usam um arquivo
	CountByFileID(ctx context.Context, fileId uuid.UUID) (int64, error)

	// Listar todas as referências de um arquivo
	ListByFileID(ctx context.Context, fileId uuid.UUID) ([]models.EntityFileReference, error)

	// Listar referências de uma entidade
	ListByEntity(ctx context.Context, entityType string, entityId uuid.UUID) ([]models.EntityFileReference, error)

	// Deletar todas as referências órfãs (deleted_at NOT NULL)
	CleanupDeletedReferences(ctx context.Context) error
}

// entityFileReferenceRepository implementa IEntityFileReferenceRepository
type entityFileReferenceRepository struct {
	db *gorm.DB
}

// NewEntityFileReferenceRepository cria nova instância do repository
func NewEntityFileReferenceRepository(db *gorm.DB) IEntityFileReferenceRepository {
	return &entityFileReferenceRepository{db: db}
}

// Create registra referência entre entidade e arquivo
func (r *entityFileReferenceRepository) Create(ctx context.Context, entityRef *models.EntityFileReference) error {
	if entityRef.Id == uuid.Nil {
		entityRef.Id = uuid.New()
	}

	entityRef.CreatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(entityRef).Error; err != nil {
		return fmt.Errorf("erro ao criar referência de entidade: %w", err)
	}

	return nil
}

// GetByEntity busca referência de uma entidade específica
func (r *entityFileReferenceRepository) GetByEntity(ctx context.Context, entityType string, entityId uuid.UUID, entityField string) (*models.EntityFileReference, error) {
	var entityRef models.EntityFileReference

	err := r.db.WithContext(ctx).
		Where("entity_type = ? AND entity_id = ? AND entity_field = ? AND deleted_at IS NULL",
			entityType, entityId, entityField).
		First(&entityRef).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar referência de entidade: %w", err)
	}

	return &entityRef, nil
}

// SoftDelete marca referência como deletada
func (r *entityFileReferenceRepository) SoftDelete(ctx context.Context, entityType string, entityId uuid.UUID, entityField string) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Model(&models.EntityFileReference{}).
		Where("entity_type = ? AND entity_id = ? AND entity_field = ?", entityType, entityId, entityField).
		Update("deleted_at", now).Error; err != nil {
		return fmt.Errorf("erro ao soft delete referência: %w", err)
	}

	return nil
}

// HardDelete deleta referência do banco
func (r *entityFileReferenceRepository) HardDelete(ctx context.Context, entityType string, entityId uuid.UUID, entityField string) error {
	if err := r.db.WithContext(ctx).
		Unscoped().
		Where("entity_type = ? AND entity_id = ? AND entity_field = ?", entityType, entityId, entityField).
		Delete(&models.EntityFileReference{}).Error; err != nil {
		return fmt.Errorf("erro ao hard delete referência: %w", err)
	}

	return nil
}

// CountByFileID conta quantas entidades usam um arquivo
func (r *entityFileReferenceRepository) CountByFileID(ctx context.Context, fileId uuid.UUID) (int64, error) {
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&models.EntityFileReference{}).
		Where("file_id = ? AND deleted_at IS NULL", fileId).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("erro ao contar referências: %w", err)
	}

	return count, nil
}

// ListByFileID lista todas as referências de um arquivo
func (r *entityFileReferenceRepository) ListByFileID(ctx context.Context, fileId uuid.UUID) ([]models.EntityFileReference, error) {
	var refs []models.EntityFileReference

	if err := r.db.WithContext(ctx).
		Where("file_id = ? AND deleted_at IS NULL", fileId).
		Find(&refs).Error; err != nil {
		return nil, fmt.Errorf("erro ao listar referências por arquivo: %w", err)
	}

	return refs, nil
}

// ListByEntity lista referências de uma entidade
func (r *entityFileReferenceRepository) ListByEntity(ctx context.Context, entityType string, entityId uuid.UUID) ([]models.EntityFileReference, error) {
	var refs []models.EntityFileReference

	if err := r.db.WithContext(ctx).
		Where("entity_type = ? AND entity_id = ? AND deleted_at IS NULL", entityType, entityId).
		Find(&refs).Error; err != nil {
		return nil, fmt.Errorf("erro ao listar referências por entidade: %w", err)
	}

	return refs, nil
}

// CleanupDeletedReferences deleta referências órfãs (deleted_at NOT NULL)
func (r *entityFileReferenceRepository) CleanupDeletedReferences(ctx context.Context) error {
	if err := r.db.WithContext(ctx).
		Unscoped().
		Where("deleted_at IS NOT NULL").
		Delete(&models.EntityFileReference{}).Error; err != nil {
		return fmt.Errorf("erro ao cleanup referências deletadas: %w", err)
	}

	return nil
}
