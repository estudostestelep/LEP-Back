package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubcategoryRepository struct {
	db *gorm.DB
}

type ISubcategoryRepository interface {
	GetSubcategory(id uuid.UUID) (*models.Subcategory, error)
	GetSubcategoryById(id uuid.UUID) (*models.Subcategory, error)
	GetSubcategoryList(organizationId, projectId uuid.UUID) ([]models.Subcategory, error)
	GetSubcategoriesByCategory(categoryId uuid.UUID) ([]models.Subcategory, error)
	GetActiveSubcategoryList(organizationId, projectId uuid.UUID) ([]models.Subcategory, error)
	CreateSubcategory(subcategory *models.Subcategory) error
	UpdateSubcategory(subcategory *models.Subcategory) error
	UpdateSubcategoryOrder(id uuid.UUID, order int) error
	UpdateSubcategoryStatus(id uuid.UUID, active bool) error
	SoftDelete(id uuid.UUID) error
	SoftDeleteSubcategory(id uuid.UUID) error
	// Relacionamento N:N com categorias
	AddCategoryToSubcategory(subcategoryId, categoryId uuid.UUID) error
	RemoveCategoryFromSubcategory(subcategoryId, categoryId uuid.UUID) error
	GetSubcategoryCategories(subcategoryId uuid.UUID) ([]models.Category, error)
}

func NewConnSubcategory(db *gorm.DB) ISubcategoryRepository {
	return &SubcategoryRepository{db: db}
}

func (r *SubcategoryRepository) CreateSubcategory(subcategory *models.Subcategory) error {
	return r.db.Create(subcategory).Error
}

func (r *SubcategoryRepository) GetSubcategory(id uuid.UUID) (*models.Subcategory, error) {
	var subcategory models.Subcategory
	err := r.db.First(&subcategory, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &subcategory, nil
}

func (r *SubcategoryRepository) GetSubcategoryById(id uuid.UUID) (*models.Subcategory, error) {
	return r.GetSubcategory(id)
}

func (r *SubcategoryRepository) GetSubcategoryList(organizationId, projectId uuid.UUID) ([]models.Subcategory, error) {
	var subcategories []models.Subcategory
	err := r.db.Where("organization_id = ? AND project_id = ? AND deleted_at IS NULL", organizationId, projectId).
		Order(`"order" ASC`).
		Find(&subcategories).Error
	return subcategories, err
}

func (r *SubcategoryRepository) GetSubcategoriesByCategory(categoryId uuid.UUID) ([]models.Subcategory, error) {
	var subcategories []models.Subcategory
	err := r.db.Table("subcategories").
		Joins("INNER JOIN subcategory_categories ON subcategory_categories.subcategory_id = subcategories.id").
		Where("subcategory_categories.category_id = ? AND subcategories.deleted_at IS NULL", categoryId).
		Order(`subcategories."order" ASC`).
		Find(&subcategories).Error
	return subcategories, err
}

func (r *SubcategoryRepository) GetActiveSubcategoryList(organizationId, projectId uuid.UUID) ([]models.Subcategory, error) {
	var subcategories []models.Subcategory
	err := r.db.Where("organization_id = ? AND project_id = ? AND active = ? AND deleted_at IS NULL", organizationId, projectId, true).
		Order(`"order" ASC`).
		Find(&subcategories).Error
	return subcategories, err
}

func (r *SubcategoryRepository) UpdateSubcategory(subcategory *models.Subcategory) error {
	return r.db.Save(subcategory).Error
}

func (r *SubcategoryRepository) UpdateSubcategoryOrder(id uuid.UUID, order int) error {
	return r.db.Model(&models.Subcategory{}).Where("id = ?", id).Update(`"order"`, order).Error
}

func (r *SubcategoryRepository) UpdateSubcategoryStatus(id uuid.UUID, active bool) error {
	return r.db.Model(&models.Subcategory{}).Where("id = ?", id).Update("active", active).Error
}

func (r *SubcategoryRepository) SoftDelete(id uuid.UUID) error {
	return r.db.Model(&models.Subcategory{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

func (r *SubcategoryRepository) SoftDeleteSubcategory(id uuid.UUID) error {
	return r.SoftDelete(id)
}

// AddCategoryToSubcategory adiciona uma categoria a uma subcategoria
func (r *SubcategoryRepository) AddCategoryToSubcategory(subcategoryId, categoryId uuid.UUID) error {
	relation := models.SubcategoryCategory{
		Id:            uuid.New(),
		SubcategoryId: subcategoryId,
		CategoryId:    categoryId,
		CreatedAt:     time.Now(),
	}
	return r.db.Create(&relation).Error
}

// RemoveCategoryFromSubcategory remove uma categoria de uma subcategoria
func (r *SubcategoryRepository) RemoveCategoryFromSubcategory(subcategoryId, categoryId uuid.UUID) error {
	return r.db.Where("subcategory_id = ? AND category_id = ?", subcategoryId, categoryId).
		Delete(&models.SubcategoryCategory{}).Error
}

// GetSubcategoryCategories retorna todas as categorias de uma subcategoria
func (r *SubcategoryRepository) GetSubcategoryCategories(subcategoryId uuid.UUID) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Table("categories").
		Joins("INNER JOIN subcategory_categories ON subcategory_categories.category_id = categories.id").
		Where("subcategory_categories.subcategory_id = ? AND categories.deleted_at IS NULL", subcategoryId).
		Order("categories.`order` ASC").
		Find(&categories).Error
	return categories, err
}
