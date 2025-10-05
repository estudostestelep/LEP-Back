package handler

import (
	"lep/repositories"
	"lep/repositories/models"

	"github.com/google/uuid"
)

type resourceProducts struct {
	repo *repositories.DBconn
}

type IHandlerProducts interface {
	GetProduct(id string) (*models.Product, error)
	GetProductByPurchase(id string) ([]models.Product, error)
	ListProducts(orgId, projectId string) ([]models.Product, error)
	CreateProduct(product *models.Product) error
	UpdateProduct(updatedProduct *models.Product) error
	UpdateProductOrder(id string, order int) error
	UpdateProductStatus(id string, active bool) error
	DeleteProduct(id string) error
	DeleteProductsByPurchase(purchaseId string) error
	// Tag management
	AddTagToProduct(productId, tagId string) error
	RemoveTagFromProduct(productId, tagId string) error
	GetProductTags(productId string) ([]models.Tag, error)
	GetProductsByTag(tagId string) ([]models.Product, error)
	// Filtros de cardápio
	GetProductsByType(orgId, projectId, productType string) ([]models.Product, error)
	GetProductsByCategory(categoryId string) ([]models.Product, error)
	GetProductsBySubcategory(subcategoryId string) ([]models.Product, error)
}

func (r *resourceProducts) GetProduct(id string) (*models.Product, error) {
	// Validar UUID
	productId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	resp, err := r.repo.Products.GetProductById(productId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceProducts) CreateProduct(product *models.Product) error {
	product.Id = uuid.New()
	err := r.repo.Products.CreateProduct(product)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceProducts) UpdateProduct(updatedProduct *models.Product) error {
	err := r.repo.Products.UpdateProduct(updatedProduct)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceProducts) DeleteProduct(id string) error {
	// Validar UUID
	productId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	err = r.repo.Products.SoftDeleteProduct(productId)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceProducts) GetProductByPurchase(id string) ([]models.Product, error) {
	resp, err := r.repo.Products.GetProductByPurchase(id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceProducts) ListProducts(orgId, projectId string) ([]models.Product, error) {
	// Converter strings para UUID
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}

	projectUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}

	resp, err := r.repo.Products.ListProducts(orgUUID, projectUUID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceProducts) DeleteProductsByPurchase(purchaseId string) error {
	err := r.repo.Products.DeleteProductsByPurchase(purchaseId)
	if err != nil {
		return err
	}
	return nil
}

// AddTagToProduct adiciona uma tag a um produto
func (r *resourceProducts) AddTagToProduct(productId, tagId string) error {
	prodUUID, err := uuid.Parse(productId)
	if err != nil {
		return err
	}
	tagUUID, err := uuid.Parse(tagId)
	if err != nil {
		return err
	}
	return r.repo.Products.AddTagToProduct(prodUUID, tagUUID)
}

// RemoveTagFromProduct remove uma tag de um produto
func (r *resourceProducts) RemoveTagFromProduct(productId, tagId string) error {
	prodUUID, err := uuid.Parse(productId)
	if err != nil {
		return err
	}
	tagUUID, err := uuid.Parse(tagId)
	if err != nil {
		return err
	}
	return r.repo.Products.RemoveTagFromProduct(prodUUID, tagUUID)
}

// GetProductTags retorna todas as tags de um produto
func (r *resourceProducts) GetProductTags(productId string) ([]models.Tag, error) {
	prodUUID, err := uuid.Parse(productId)
	if err != nil {
		return nil, err
	}
	return r.repo.Products.GetProductTags(prodUUID)
}

// GetProductsByTag retorna todos os produtos que possuem uma tag específica
func (r *resourceProducts) GetProductsByTag(tagId string) ([]models.Product, error) {
	tagUUID, err := uuid.Parse(tagId)
	if err != nil {
		return nil, err
	}
	return r.repo.Products.GetProductsByTag(tagUUID)
}

// UpdateProductOrder atualiza a ordem de um produto
func (r *resourceProducts) UpdateProductOrder(id string, order int) error {
	prodUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.repo.Products.UpdateProductOrder(prodUUID, order)
}

// UpdateProductStatus atualiza o status de um produto
func (r *resourceProducts) UpdateProductStatus(id string, active bool) error {
	prodUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.repo.Products.UpdateProductStatus(prodUUID, active)
}

// GetProductsByType filtra produtos por tipo
func (r *resourceProducts) GetProductsByType(orgId, projectId, productType string) ([]models.Product, error) {
	orgUUID, _ := uuid.Parse(orgId)
	projectUUID, _ := uuid.Parse(projectId)
	return r.repo.Products.GetProductsByType(orgUUID, projectUUID, productType)
}

// GetProductsByCategory retorna produtos de uma categoria
func (r *resourceProducts) GetProductsByCategory(categoryId string) ([]models.Product, error) {
	catUUID, err := uuid.Parse(categoryId)
	if err != nil {
		return nil, err
	}
	return r.repo.Products.GetProductsByCategory(catUUID)
}

// GetProductsBySubcategory retorna produtos de uma subcategoria
func (r *resourceProducts) GetProductsBySubcategory(subcategoryId string) ([]models.Product, error) {
	subUUID, err := uuid.Parse(subcategoryId)
	if err != nil {
		return nil, err
	}
	return r.repo.Products.GetProductsBySubcategory(subUUID)
}

func NewSourceHandlerProducts(repo *repositories.DBconn) IHandlerProducts {
	return &resourceProducts{repo: repo}
}
