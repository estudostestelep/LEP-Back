package handler

import (
	"lep/repositories"
	"lep/repositories/models"
)

type resourceProducts struct {
	repo *repositories.DBconn
}

type IHandlerProducts interface {
	GetProduct(id int) (*models.Product, error)
	GetProductByPurchase(id string) ([]models.Product, error)
	CreateProduct(product *models.Product) error
	UpdateProduct(updatedProduct *models.Product) error
	DeleteProduct(id int) error
	DeleteProductsByPurchase(purchaseId string) error
}

func (r *resourceProducts) GetProduct(id int) (*models.Product, error) {
	resp, err := r.repo.Products.GetProduct(id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceProducts) CreateProduct(product *models.Product) error {
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

func (r *resourceProducts) DeleteProduct(id int) error {
	err := r.repo.Products.DeleteProduct(id)
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
