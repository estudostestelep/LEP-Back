package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

type ICategoryRepository interface {
	GetCategory(id uuid.UUID) (*models.Category, error)
	GetCategoryById(id uuid.UUID) (*models.Category, error)
	GetCategoryList(organizationId, projectId uuid.UUID) ([]models.Category, error)
	GetCategoriesByMenu(menuId uuid.UUID) ([]models.Category, error)
	GetActiveCategoryList(organizationId, projectId uuid.UUID) ([]models.Category, error)
	CreateCategory(category *models.Category) error
	UpdateCategory(category *models.Category) error
	UpdateCategoryOrder(id uuid.UUID, order int) error
	UpdateCategoryStatus(id uuid.UUID, active bool) error
	SoftDelete(id uuid.UUID) error
	SoftDeleteCategory(id uuid.UUID) error
}

func NewConnCategory(db *gorm.DB) ICategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) CreateCategory(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *CategoryRepository) GetCategory(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) GetCategoryById(id uuid.UUID) (*models.Category, error) {
	return r.GetCategory(id)
}

func (r *CategoryRepository) GetCategoryList(organizationId, projectId uuid.UUID) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Where("organization_id = ? AND project_id = ? AND deleted_at IS NULL", organizationId, projectId).
		Order(`"order" ASC`).
		Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) GetCategoriesByMenu(menuId uuid.UUID) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Where("menu_id = ? AND deleted_at IS NULL", menuId).
		Order(`"order" ASC`).
		Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) GetActiveCategoryList(organizationId, projectId uuid.UUID) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Where("organization_id = ? AND project_id = ? AND active = ? AND deleted_at IS NULL", organizationId, projectId, true).
		Order(`"order" ASC`).
		Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) UpdateCategory(category *models.Category) error {
	// Usar Updates com Select para garantir que campos ponteiros (como Photo) sejam atualizados
	// mesmo quando são nil
	return r.db.Model(&models.Category{}).
		Where("id = ?", category.Id).
		Select("*").  // Selecionar todos os campos para atualização
		Updates(category).Error
}

func (r *CategoryRepository) UpdateCategoryOrder(id uuid.UUID, order int) error {
	return r.db.Model(&models.Category{}).Where("id = ?", id).Update(`"order"`, order).Error
}

func (r *CategoryRepository) UpdateCategoryStatus(id uuid.UUID, active bool) error {
	return r.db.Model(&models.Category{}).Where("id = ?", id).Update("active", active).Error
}

func (r *CategoryRepository) SoftDelete(id uuid.UUID) error {
	return r.db.Model(&models.Category{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

func (r *CategoryRepository) SoftDeleteCategory(id uuid.UUID) error {
	return r.SoftDelete(id)
}
