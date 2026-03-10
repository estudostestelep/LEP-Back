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
	ServiceUpdateProductImage(c *gin.Context)
	ServiceUpdateProductOrder(c *gin.Context)
	ServiceUpdateProductStatus(c *gin.Context)
	ServiceDeleteProduct(c *gin.Context)
	// Tag management
	ServiceAddTagToProduct(c *gin.Context)
	ServiceRemoveTagFromProduct(c *gin.Context)
	ServiceGetProductTags(c *gin.Context)
	ServiceGetProductsByTag(c *gin.Context)
	// Filtros de cardápio
	ServiceGetProductsByType(c *gin.Context)
	ServiceGetProductsByCategory(c *gin.Context)
	ServiceGetProductsBySubcategory(c *gin.Context)
}

func (r *ResourceProducts) ServiceGetProduct(c *gin.Context) {
	id, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "product")
	if !ok {
		return
	}

	resp, err := r.handler.HandlerProducts.GetProduct(id.String())
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

	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	newProduct.OrganizationId, err = uuid.Parse(organizationId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing organization ID", err)
		return
	}
	newProduct.ProjectId, err = uuid.Parse(projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing project ID", err)
		return
	}

	// Gerar ID se não fornecido
	if newProduct.Id == uuid.Nil {
		newProduct.Id = uuid.New()
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
	id, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "product")
	if !ok {
		return
	}

	var updatedProduct models.Product
	if err := c.BindJSON(&updatedProduct); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	var err error
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
	updatedProduct.Id = id

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

	// Verificar se deve incluir tags
	includeTags := c.Query("includeTags") == "true"

	// Verificar se há filtros nos query parameters
	categoryIdStr := c.Query("category_id")
	subcategoryIdStr := c.Query("subcategory_id")
	tagIdStr := c.Query("tag_id")
	productType := c.Query("type")
	activeStr := c.Query("active")

	// Se há algum filtro, usar ListProductsWithFilters
	if categoryIdStr != "" || subcategoryIdStr != "" || tagIdStr != "" || productType != "" || activeStr != "" {
		filters := handler.ProductFilters{}

		// Parse category_id
		if categoryIdStr != "" {
			catUUID, err := uuid.Parse(categoryIdStr)
			if err != nil {
				utils.SendBadRequestError(c, "Invalid category_id format", err)
				return
			}
			filters.CategoryId = &catUUID
		}

		// Parse subcategory_id
		if subcategoryIdStr != "" {
			subUUID, err := uuid.Parse(subcategoryIdStr)
			if err != nil {
				utils.SendBadRequestError(c, "Invalid subcategory_id format", err)
				return
			}
			filters.SubcategoryId = &subUUID
		}

		// Parse tag_id
		if tagIdStr != "" {
			tagUUID, err := uuid.Parse(tagIdStr)
			if err != nil {
				utils.SendBadRequestError(c, "Invalid tag_id format", err)
				return
			}
			filters.TagId = &tagUUID
		}

		// Parse type
		if productType != "" {
			filters.Type = &productType
		}

		// Parse active
		if activeStr != "" {
			active := activeStr == "true"
			filters.Active = &active
		}

		products, err := r.handler.HandlerProducts.ListProductsWithFilters(organizationId, projectId, filters)
		if err != nil {
			utils.SendInternalServerError(c, "Error listing products with filters", err)
			return
		}

		c.JSON(http.StatusOK, products)
		return
	}

	// Se includeTags=true, usar listagem com tags
	if includeTags {
		products, err := r.handler.HandlerProducts.ListProductsWithTags(organizationId, projectId)
		if err != nil {
			utils.SendInternalServerError(c, "Error listing products with tags", err)
			return
		}

		c.JSON(http.StatusOK, products)
		return
	}

	// Sem filtros e sem tags, usar listagem normal
	products, err := r.handler.HandlerProducts.ListProducts(organizationId, projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing products", err)
		return
	}

	c.JSON(http.StatusOK, products)
}

func (r *ResourceProducts) ServiceUpdateProductImage(c *gin.Context) {
	id, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "product")
	if !ok {
		return
	}

	// Headers validados pelo middleware
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	// Parse do JSON body para receber nova image_url
	var requestData struct {
		ImageUrl string `json:"image_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Buscar produto existente
	existingProduct, err := r.handler.HandlerProducts.GetProduct(id.String())
	if err != nil {
		utils.SendInternalServerError(c, "Error getting product", err)
		return
	}

	if existingProduct == nil {
		utils.SendNotFoundError(c, "Product")
		return
	}

	// Verificar se o produto pertence à organização/projeto corretos
	if existingProduct.OrganizationId.String() != organizationId ||
		existingProduct.ProjectId.String() != projectId {
		utils.SendBadRequestError(c, "Product does not belong to specified organization/project", nil)
		return
	}

	// Atualizar apenas o campo ImageUrl
	existingProduct.ImageUrl = &requestData.ImageUrl

	// Salvar alteração
	err = r.handler.HandlerProducts.UpdateProduct(existingProduct)
	if err != nil {
		utils.SendInternalServerError(c, "Error updating product image", err)
		return
	}

	utils.SendOKSuccess(c, "Product image updated successfully", existingProduct)
}

func (r *ResourceProducts) ServiceDeleteProduct(c *gin.Context) {
	id, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "product")
	if !ok {
		return
	}

	if err := r.handler.HandlerProducts.DeleteProduct(id.String()); err != nil {
		utils.SendInternalServerError(c, "Error deleting product", err)
		return
	}

	utils.SendOKSuccess(c, "Product deleted successfully", nil)
}

// ServiceAddTagToProduct adiciona uma tag a um produto
func (r *ResourceProducts) ServiceAddTagToProduct(c *gin.Context) {
	productId, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "product")
	if !ok {
		return
	}

	var requestBody struct {
		TagId string `json:"tag_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Validar formato UUID da tag
	tagId, ok := validation.ParseAndValidateUUID(c, requestBody.TagId, "tag")
	if !ok {
		return
	}

	err := r.handler.HandlerProducts.AddTagToProduct(productId.String(), tagId.String())
	if err != nil {
		utils.SendInternalServerError(c, "Error adding tag to product", err)
		return
	}

	utils.SendCreatedSuccess(c, "Tag added to product successfully", nil)
}

// ServiceRemoveTagFromProduct remove uma tag de um produto
func (r *ResourceProducts) ServiceRemoveTagFromProduct(c *gin.Context) {
	productId, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "product")
	if !ok {
		return
	}

	tagId, ok := validation.ParseAndValidateUUID(c, c.Param("tagId"), "tag")
	if !ok {
		return
	}

	err := r.handler.HandlerProducts.RemoveTagFromProduct(productId.String(), tagId.String())
	if err != nil {
		utils.SendInternalServerError(c, "Error removing tag from product", err)
		return
	}

	utils.SendOKSuccess(c, "Tag removed from product successfully", nil)
}

// ServiceGetProductTags retorna todas as tags de um produto
func (r *ResourceProducts) ServiceGetProductTags(c *gin.Context) {
	productId, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "product")
	if !ok {
		return
	}

	tags, err := r.handler.HandlerProducts.GetProductTags(productId.String())
	if err != nil {
		utils.SendInternalServerError(c, "Error getting product tags", err)
		return
	}

	c.JSON(http.StatusOK, tags)
}

// ServiceGetProductsByTag retorna todos os produtos que possuem uma tag específica
func (r *ResourceProducts) ServiceGetProductsByTag(c *gin.Context) {
	tagIdStr := c.Query("tag_id")

	if tagIdStr == "" {
		utils.SendBadRequestError(c, "tag_id query parameter is required", nil)
		return
	}

	tagId, ok := validation.ParseAndValidateUUID(c, tagIdStr, "tag")
	if !ok {
		return
	}

	products, err := r.handler.HandlerProducts.GetProductsByTag(tagId.String())
	if err != nil {
		utils.SendInternalServerError(c, "Error getting products by tag", err)
		return
	}

	c.JSON(http.StatusOK, products)
}

// ServiceUpdateProductOrder atualiza a ordem de exibição de um produto
func (r *ResourceProducts) ServiceUpdateProductOrder(c *gin.Context) {
	productId, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "product")
	if !ok {
		return
	}

	var requestBody struct {
		Order int `json:"order" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	err := r.handler.HandlerProducts.UpdateProductOrder(productId.String(), requestBody.Order)
	if err != nil {
		utils.SendInternalServerError(c, "Error updating product order", err)
		return
	}

	utils.SendOKSuccess(c, "Product order updated successfully", nil)
}

// ServiceUpdateProductStatus atualiza o status (play/pause) de um produto
func (r *ResourceProducts) ServiceUpdateProductStatus(c *gin.Context) {
	productId, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "product")
	if !ok {
		return
	}

	var requestBody struct {
		Active *bool `json:"active" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	if requestBody.Active == nil {
		utils.SendBadRequestError(c, "Active field is required", nil)
		return
	}

	err := r.handler.HandlerProducts.UpdateProductStatus(productId.String(), *requestBody.Active)
	if err != nil {
		utils.SendInternalServerError(c, "Error updating product status", err)
		return
	}

	utils.SendOKSuccess(c, "Product status updated successfully", nil)
}

// ServiceGetProductsByType filtra produtos por tipo (prato, bebida, vinho)
func (r *ResourceProducts) ServiceGetProductsByType(c *gin.Context) {
	productType := c.Param("type")

	if productType == "" {
		utils.SendBadRequestError(c, "type parameter is required", nil)
		return
	}

	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	products, err := r.handler.HandlerProducts.GetProductsByType(organizationId, projectId, productType)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting products by type", err)
		return
	}

	c.JSON(http.StatusOK, products)
}

// ServiceGetProductsByCategory filtra produtos por categoria
func (r *ResourceProducts) ServiceGetProductsByCategory(c *gin.Context) {
	categoryIdStr := c.Param("categoryId")

	if categoryIdStr == "" {
		utils.SendBadRequestError(c, "category_id parameter is required", nil)
		return
	}

	categoryId, ok := validation.ParseAndValidateUUID(c, categoryIdStr, "category")
	if !ok {
		return
	}

	products, err := r.handler.HandlerProducts.GetProductsByCategory(categoryId.String())
	if err != nil {
		utils.SendInternalServerError(c, "Error getting products by category", err)
		return
	}

	c.JSON(http.StatusOK, products)
}

// ServiceGetProductsBySubcategory filtra produtos por subcategoria
func (r *ResourceProducts) ServiceGetProductsBySubcategory(c *gin.Context) {
	subcategoryIdStr := c.Param("subcategoryId")

	if subcategoryIdStr == "" {
		utils.SendBadRequestError(c, "subcategory_id parameter is required", nil)
		return
	}

	subcategoryId, ok := validation.ParseAndValidateUUID(c, subcategoryIdStr, "subcategory")
	if !ok {
		return
	}

	products, err := r.handler.HandlerProducts.GetProductsBySubcategory(subcategoryId.String())
	if err != nil {
		utils.SendInternalServerError(c, "Error getting products by subcategory", err)
		return
	}

	c.JSON(http.StatusOK, products)
}

func NewSourceServerProducts(handler *handler.Handlers) IServerProducts {
	return &ResourceProducts{handler: handler}
}
