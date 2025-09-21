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
	DeleteProduct(id string) error
	DeleteProductsByPurchase(purchaseId string) error
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

func NewSourceHandlerProducts(repo *repositories.DBconn) IHandlerProducts {
	return &resourceProducts{repo: repo}
}
