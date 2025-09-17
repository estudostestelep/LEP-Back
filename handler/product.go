package handler

import (
	"lep/repositories"
	"lep/repositories/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type resourceProducts struct {
	repo *repositories.DBconn
}

type IHandlerProducts interface {
	GetProduct(id int) (*models.Products, error)
	GetProductByPurchase(id string) ([]models.Products, error)
	CreateProduct(product *models.Products) error
	UpdateProduct(updatedProduct *models.Products) error
	DeleteProduct(id int) error
	DeleteProductsByPurchase(purchaseId string) error
}

func (r *resourceProducts) GetProduct(id int) (*models.Products, error) {
	resp, err := r.repo.Products.GetProduct(id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceProducts) GetProductByPurchase(id string) ([]models.Products, error) {
	resp, err := r.repo.Products.GetProductByPurchase(id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceProducts) CreateProduct(product *models.Products) error {
	err := r.repo.Products.CreateProduct(product)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceProducts) UpdateProduct(updatedProduct *models.Products) error {
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

func (r *resourceProducts) DeleteProductsByPurchase(purchaseId string) error {
	err := r.repo.Products.DeleteProductsByPurchase(purchaseId)
	if err != nil {
		return err
	}
	return nil
}

func (h *resourceProducts) Create(c *gin.Context) {
	var req models.Product
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}
	req.Id = uuid.New()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	if err := h.repo.Create(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating product"})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *resourceProducts) GetById(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	product, err := h.repo.GetById(id)
	if err != nil || product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *resourceProducts) List(c *gin.Context) {
	OrganizationId, _ := uuid.Parse(c.Query("org_id"))
	projectId, _ := uuid.Parse(c.Query("project_id"))
	products, err := h.repo.List(OrganizationId, projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error listing products"})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h *resourceProducts) Update(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	product, err := h.repo.GetById(id)
	if err != nil || product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	if err := c.ShouldBindJSON(product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}
	product.UpdatedAt = time.Now()
	if err := h.repo.Update(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating product"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *resourceProducts) SoftDelete(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	err := h.repo.SoftDelete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

func NewSourceHandlerProducts(repo *repositories.DBconn) IHandlerProducts {
	return &resourceProducts{repo: repo}
}
