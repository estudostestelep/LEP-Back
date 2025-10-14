package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

type ITagRepository interface {
	GetTag(id uuid.UUID) (*models.Tag, error)
	GetTagById(id uuid.UUID) (*models.Tag, error)
	GetTagList(organizationId, projectId uuid.UUID) ([]models.Tag, error)
	GetActiveTagList(organizationId, projectId uuid.UUID) ([]models.Tag, error)
	GetTagsByEntityType(organizationId, projectId uuid.UUID, entityType string) ([]models.Tag, error)
	CreateTag(tag *models.Tag) error
	UpdateTag(tag *models.Tag) error
	SoftDelete(id uuid.UUID) error
	SoftDeleteTag(id uuid.UUID) error
}

func NewConnTag(db *gorm.DB) ITagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) CreateTag(tag *models.Tag) error {
	return r.db.Create(tag).Error
}

func (r *TagRepository) GetTag(id uuid.UUID) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.First(&tag, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepository) GetTagById(id uuid.UUID) (*models.Tag, error) {
	return r.GetTag(id)
}

func (r *TagRepository) GetTagList(organizationId, projectId uuid.UUID) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Where("organization_id = ? AND project_id = ? AND deleted_at IS NULL", organizationId, projectId).Find(&tags).Error
	return tags, err
}

func (r *TagRepository) GetActiveTagList(organizationId, projectId uuid.UUID) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Where("organization_id = ? AND project_id = ? AND active = ? AND deleted_at IS NULL", organizationId, projectId, true).Find(&tags).Error
	return tags, err
}

func (r *TagRepository) GetTagsByEntityType(organizationId, projectId uuid.UUID, entityType string) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Where("organization_id = ? AND project_id = ? AND entity_type = ? AND deleted_at IS NULL", organizationId, projectId, entityType).Find(&tags).Error
	return tags, err
}

func (r *TagRepository) UpdateTag(tag *models.Tag) error {
	return r.db.Save(tag).Error
}

func (r *TagRepository) SoftDelete(id uuid.UUID) error {
	return r.db.Model(&models.Tag{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

func (r *TagRepository) SoftDeleteTag(id uuid.UUID) error {
	return r.SoftDelete(id)
}
