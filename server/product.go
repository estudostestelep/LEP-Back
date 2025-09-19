package server

import (
	"fmt"
	"lep/handler"
	"lep/repositories/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	id := c.Param("id")
	intNumber, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Erro ao converter string para int:", err)
		return
	}
	resp, err := r.handler.HandlerProducts.GetProduct(intNumber)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao obter o produto")
		return
	}

	if resp == nil {
		c.String(http.StatusNotFound, "Produto não encontrado")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceProducts) ServiceGetProductByPurchase(c *gin.Context) {
	id := c.Param("id")
	resp, err := r.handler.HandlerProducts.GetProductByPurchase(id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao obter o produto")
		return
	}

	if resp == nil {
		c.String(http.StatusNotFound, "Produto não encontrado")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceProducts) ServiceCreateProduct(c *gin.Context) {
	var newProduct models.Product
	err := c.BindJSON(&newProduct)
	if err != nil {
		c.String(http.StatusBadRequest, "Erro ao decodificar dados do produto")
		return
	}

	err = r.handler.HandlerProducts.CreateProduct(&newProduct)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao criar o produto")
		return
	}

	c.String(http.StatusCreated, "Produto criado com sucesso")
}

func (r *ResourceProducts) ServiceUpdateProduct(c *gin.Context) {
	var updatedProduct models.Product
	err := c.BindJSON(&updatedProduct)
	if err != nil {
		c.String(http.StatusBadRequest, "Erro ao decodificar dados do produto")
		return
	}

	err = r.handler.HandlerProducts.UpdateProduct(&updatedProduct)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao atualizar o produto")
		return
	}

	c.String(http.StatusOK, "Produto atualizado com sucesso")
}

func (r *ResourceProducts) ServiceListProducts(c *gin.Context) {
	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	products, err := r.handler.HandlerProducts.ListProducts(organizationId, projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error listing products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (r *ResourceProducts) ServiceDeleteProduct(c *gin.Context) {
	id := c.Param("id")
	intNumber, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Erro ao converter string para int:", err)
		return
	}

	err = r.handler.HandlerProducts.DeleteProduct(intNumber)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao excluir o produto")
		return
	}

	c.String(http.StatusOK, "Produto excluído com sucesso")
}

func NewSourceServerProducts(handler *handler.Handlers) IServerProducts {
	return &ResourceProducts{handler: handler}
}
