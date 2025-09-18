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
	GetProduct(id int) (*models.Product, error)
	GetProductById(id uuid.UUID) (*models.Product, error)
	GetProductsByIds(ids []uuid.UUID) ([]models.Product, error)
	GetProductByPurchase(id string) ([]models.Product, error)
	ListProducts(OrganizationId, projectId uuid.UUID) ([]models.Product, error)
	CreateProduct(product *models.Product) error
	UpdateProduct(product *models.Product) error
	DeleteProduct(id int) error
	DeleteProductsByPurchase(purchaseId string) error
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

func (r *resourceProduct) GetProduct(id int) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *resourceProduct) GetProductByPurchase(id string) ([]models.Product, error) {
	var products []models.Product
	// Implementação simplificada - pode precisar ajustar conforme a lógica de negócio
	err := r.db.Find(&products, "deleted_at IS NULL").Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *resourceProduct) DeleteProduct(id int) error {
	return r.db.Delete(&models.Product{}, id).Error
}

func (r *resourceProduct) DeleteProductsByPurchase(purchaseId string) error {
	// Implementação simplificada - pode precisar ajustar conforme a lógica de negócio
	return r.db.Where("purchase_id = ?", purchaseId).Delete(&models.Product{}).Error
}

func (r *resourceProduct) SoftDeleteProduct(id uuid.UUID) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

// GetProductsByIds busca múltiplos produtos por IDs
func (r *resourceProduct) GetProductsByIds(ids []uuid.UUID) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("id IN ? AND deleted_at IS NULL", ids).Find(&products).Error
	return products, err
}
