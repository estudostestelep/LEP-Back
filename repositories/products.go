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
	ListProductsWithTags(organizationId, projectId uuid.UUID) ([]models.Product, error)
	ListProductsWithFilters(organizationId, projectId uuid.UUID, filters interface{}) ([]models.Product, error)
	CreateProduct(product *models.Product) error
	UpdateProduct(product *models.Product) error
	UpdateProductOrder(id uuid.UUID, order int) error
	UpdateProductStatus(id uuid.UUID, active bool) error
	DeleteProduct(id int) error
	DeleteProductsByPurchase(purchaseId string) error
	SoftDeleteProduct(id uuid.UUID) error
	// Tag management
	AddTagToProduct(productId, tagId uuid.UUID) error
	RemoveTagFromProduct(productId, tagId uuid.UUID) error
	GetProductTags(productId uuid.UUID) ([]models.Tag, error)
	GetProductsByTag(tagId uuid.UUID) ([]models.Product, error)
	// Filtros avançados (cardápio)
	GetProductsByType(organizationId, projectId uuid.UUID, productType string) ([]models.Product, error)
	GetProductsByCategory(categoryId uuid.UUID) ([]models.Product, error)
	GetProductsBySubcategory(subcategoryId uuid.UUID) ([]models.Product, error)
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
	err := r.db.Where("organization_id = ? AND project_id = ? AND deleted_at IS NULL", OrganizationId, projectId).Find(&products).Error
	return products, err
}

// ListProductsWithTags retorna produtos com tags eager-loaded
func (r *resourceProduct) ListProductsWithTags(organizationId, projectId uuid.UUID) ([]models.Product, error) {
	var products []models.Product
	err := r.db.
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Where("active = ?", true).Where("entity_type = ?", "product")
		}).
		Where("organization_id = ? AND project_id = ? AND deleted_at IS NULL", organizationId, projectId).
		Find(&products).Error
	return products, err
}

func (r *resourceProduct) ListProductsWithFilters(organizationId, projectId uuid.UUID, filters interface{}) ([]models.Product, error) {
	query := r.db.Where("organization_id = ? AND project_id = ? AND deleted_at IS NULL", organizationId, projectId)

	// Type assertion para acessar os filtros
	if f, ok := filters.(map[string]interface{}); ok {
		// Filtro por categoria
		if categoryId, exists := f["CategoryId"]; exists && categoryId != nil {
			if catUUID, ok := categoryId.(*uuid.UUID); ok && catUUID != nil {
				query = query.Where("category_id = ?", catUUID)
			}
		}

		// Filtro por subcategoria
		if subcategoryId, exists := f["SubcategoryId"]; exists && subcategoryId != nil {
			if subUUID, ok := subcategoryId.(*uuid.UUID); ok && subUUID != nil {
				query = query.Where("subcategory_id = ?", subUUID)
			}
		}

		// Filtro por tipo (prato, bebida, vinho)
		if productType, exists := f["Type"]; exists && productType != nil {
			if typeStr, ok := productType.(*string); ok && typeStr != nil {
				query = query.Where("type = ?", *typeStr)
			}
		}

		// Filtro por status ativo
		if active, exists := f["Active"]; exists && active != nil {
			if activePtr, ok := active.(*bool); ok && activePtr != nil {
				query = query.Where("active = ?", *activePtr)
			}
		}

		// Filtro por tag (requer JOIN)
		if tagId, exists := f["TagId"]; exists && tagId != nil {
			if tagUUID, ok := tagId.(*uuid.UUID); ok && tagUUID != nil {
				query = query.Joins("INNER JOIN product_tags ON product_tags.product_id = products.id").
					Where("product_tags.tag_id = ?", tagUUID)
			}
		}
	}

	var products []models.Product
	err := query.Order(`"order" ASC`).Find(&products).Error
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

// AddTagToProduct adiciona uma tag a um produto
func (r *resourceProduct) AddTagToProduct(productId, tagId uuid.UUID) error {
	productTag := models.ProductTag{
		Id:        uuid.New(),
		ProductId: productId,
		TagId:     tagId,
		CreatedAt: time.Now(),
	}
	return r.db.Create(&productTag).Error
}

// RemoveTagFromProduct remove uma tag de um produto
func (r *resourceProduct) RemoveTagFromProduct(productId, tagId uuid.UUID) error {
	return r.db.Where("product_id = ? AND tag_id = ?", productId, tagId).Delete(&models.ProductTag{}).Error
}

// GetProductTags retorna todas as tags de um produto
func (r *resourceProduct) GetProductTags(productId uuid.UUID) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Table("tags").
		Joins("INNER JOIN product_tags ON product_tags.tag_id = tags.id").
		Where("product_tags.product_id = ? AND tags.deleted_at IS NULL", productId).
		Find(&tags).Error
	return tags, err
}

// GetProductsByTag retorna todos os produtos que possuem uma tag específica
func (r *resourceProduct) GetProductsByTag(tagId uuid.UUID) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Table("products").
		Joins("INNER JOIN product_tags ON product_tags.product_id = products.id").
		Where("product_tags.tag_id = ? AND products.deleted_at IS NULL", tagId).
		Order(`products."order" ASC`).
		Find(&products).Error
	return products, err
}

// UpdateProductOrder atualiza a ordem de um produto
func (r *resourceProduct) UpdateProductOrder(id uuid.UUID, order int) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).Update(`"order"`, order).Error
}

// UpdateProductStatus atualiza o status ativo/inativo de um produto
func (r *resourceProduct) UpdateProductStatus(id uuid.UUID, active bool) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).Update("active", active).Error
}

// GetProductsByType retorna produtos filtrados por tipo (prato, bebida, vinho)
func (r *resourceProduct) GetProductsByType(organizationId, projectId uuid.UUID, productType string) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("organization_id = ? AND project_id = ? AND type = ? AND deleted_at IS NULL", organizationId, projectId, productType).
		Order(`"order" ASC`).
		Find(&products).Error
	return products, err
}

// GetProductsByCategory retorna produtos de uma categoria específica
func (r *resourceProduct) GetProductsByCategory(categoryId uuid.UUID) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("category_id = ? AND deleted_at IS NULL", categoryId).
		Order(`"order" ASC`).
		Find(&products).Error
	return products, err
}

// GetProductsBySubcategory retorna produtos de uma subcategoria específica
func (r *resourceProduct) GetProductsBySubcategory(subcategoryId uuid.UUID) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("subcategory_id = ? AND deleted_at IS NULL", subcategoryId).
		Order(`"order" ASC`).
		Find(&products).Error
	return products, err
}
