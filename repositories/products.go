package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type resourceProduct struct {
	db *gorm.DB
}

type IProductRepository interface {
	GetProductById(id uuid.UUID) (*models.Product, error)
	ListProducts(OrganizationId, projectId uuid.UUID) ([]models.Product, error)
	CreateProduct(product *models.Product) error
	UpdateProduct(product *models.Product) error
	SoftDeleteProduct(id uuid.UUID) error
}

func NewConnProduct(db *gorm.DB) IProductRepository {
	return &resourceProduct{db: db}
}

func (r *resourceProduct) GetProductById(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *resourceProduct) ListProducts(OrganizationId, projectId uuid.UUID) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("org_id = ? AND project_id = ? AND deleted_at IS NULL", OrganizationId, projectId).Find(&products).Error
	return products, err
}

func (r *resourceProduct) CreateProduct(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *resourceProduct) UpdateProduct(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *resourceProduct) SoftDeleteProduct(id uuid.UUID) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}
