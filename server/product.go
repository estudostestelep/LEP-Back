package server

import (
	"lep/handler"
	"lep/repositories/models"
	"lep/resource/validation"
	"lep/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ResourceProducts struct {
	handler *handler.Handlers
}

type IServerProducts interface {
	ServiceGetProduct(c *gin.Context)
	ServiceGetProductByPurchase(c *gin.Context)
	ServiceListProducts(c *gin.Context)
	ServiceCreateProduct(c *gin.Context)
	ServiceUpdateProduct(c *gin.Context)
	ServiceDeleteProduct(c *gin.Context)
}

func (r *ResourceProducts) ServiceGetProduct(c *gin.Context) {
	idStr := c.Param("id")

	// Validar formato UUID
	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid product ID format", err)
		return
	}

	resp, err := r.handler.HandlerProducts.GetProduct(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting product", err)
		return
	}

	if resp == nil {
		utils.SendNotFoundError(c, "Product")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceProducts) ServiceGetProductByPurchase(c *gin.Context) {
	id := c.Param("id")
	resp, err := r.handler.HandlerProducts.GetProductByPurchase(id)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting products by purchase", err)
		return
	}

	if resp == nil {
		utils.SendNotFoundError(c, "Products")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceProducts) ServiceCreateProduct(c *gin.Context) {
	var newProduct models.Product
	err := c.BindJSON(&newProduct)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}
	newProduct.OrganizationId, err = uuid.Parse(c.GetHeader("X-Lpe-Organization-Id"))
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing organization ID", err)
		return
	}
	newProduct.ProjectId, err = uuid.Parse(c.GetHeader("X-Lpe-Project-Id"))
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing project ID", err)
		return
	}
	// Validações estruturadas
	if err := validation.CreateProductValidation(&newProduct); err != nil {
		utils.SendValidationError(c, "Validation failed", err)
		return
	}

	err = r.handler.HandlerProducts.CreateProduct(&newProduct)
	if err != nil {
		utils.SendInternalServerError(c, "Error creating product", err)
		return
	}

	utils.SendCreatedSuccess(c, "Product created successfully", newProduct)
}

func (r *ResourceProducts) ServiceUpdateProduct(c *gin.Context) {
	idStr := c.Param("id")

	// Validar formato UUID
	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid product ID format", err)
		return
	}

	var updatedProduct models.Product
	err = c.BindJSON(&updatedProduct)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	updatedProduct.OrganizationId, err = uuid.Parse(organizationId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing organization ID", err)
		return
	}
	updatedProduct.ProjectId, err = uuid.Parse(projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing project ID", err)
		return
	}
	updatedProduct.Id, err = uuid.Parse(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing product ID", err)
		return
	}

	// Validações estruturadas
	if err := validation.UpdateProductValidation(&updatedProduct); err != nil {
		utils.SendValidationError(c, "Validation failed", err)
		return
	}

	err = r.handler.HandlerProducts.UpdateProduct(&updatedProduct)
	if err != nil {
		utils.SendInternalServerError(c, "Error updating product", err)
		return
	}

	utils.SendOKSuccess(c, "Product updated successfully", updatedProduct)
}

func (r *ResourceProducts) ServiceListProducts(c *gin.Context) {
	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	products, err := r.handler.HandlerProducts.ListProducts(organizationId, projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing products", err)
		return
	}

	c.JSON(http.StatusOK, products)
}

func (r *ResourceProducts) ServiceDeleteProduct(c *gin.Context) {
	idStr := c.Param("id")

	// Validar formato UUID
	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid product ID format", err)
		return
	}

	err = r.handler.HandlerProducts.DeleteProduct(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error deleting product", err)
		return
	}

	utils.SendOKSuccess(c, "Product deleted successfully", nil)
}

func NewSourceServerProducts(handler *handler.Handlers) IServerProducts {
	return &ResourceProducts{handler: handler}
}
