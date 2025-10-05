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

type ResourceSubcategory struct {
	handler *handler.Handlers
}

type IServerSubcategory interface {
	ServiceGetSubcategory(c *gin.Context)
	ServiceListSubcategories(c *gin.Context)
	ServiceGetSubcategoriesByCategory(c *gin.Context)
	ServiceListActiveSubcategories(c *gin.Context)
	ServiceCreateSubcategory(c *gin.Context)
	ServiceUpdateSubcategory(c *gin.Context)
	ServiceUpdateSubcategoryOrder(c *gin.Context)
	ServiceUpdateSubcategoryStatus(c *gin.Context)
	ServiceDeleteSubcategory(c *gin.Context)
	ServiceAddCategoryToSubcategory(c *gin.Context)
	ServiceRemoveCategoryFromSubcategory(c *gin.Context)
	ServiceGetSubcategoryCategories(c *gin.Context)
}

func (r *ResourceSubcategory) ServiceGetSubcategory(c *gin.Context) {
	idStr := c.Param("id")

	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid subcategory ID format", err)
		return
	}

	resp, err := r.handler.HandlerSubcategory.GetSubcategory(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting subcategory", err)
		return
	}

	if resp == nil {
		utils.SendNotFoundError(c, "Subcategory")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceSubcategory) ServiceListSubcategories(c *gin.Context) {
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	resp, err := r.handler.HandlerSubcategory.ListSubcategories(organizationId, projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing subcategories", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceSubcategory) ServiceGetSubcategoriesByCategory(c *gin.Context) {
	categoryId := c.Param("categoryId")

	_, err := uuid.Parse(categoryId)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid category ID format", err)
		return
	}

	resp, err := r.handler.HandlerSubcategory.GetSubcategoriesByCategory(categoryId)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting subcategories by category", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceSubcategory) ServiceListActiveSubcategories(c *gin.Context) {
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	resp, err := r.handler.HandlerSubcategory.ListActiveSubcategories(organizationId, projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing active subcategories", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceSubcategory) ServiceCreateSubcategory(c *gin.Context) {
	var newSubcategory models.Subcategory
	err := c.BindJSON(&newSubcategory)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	newSubcategory.OrganizationId, err = uuid.Parse(organizationId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing organization ID", err)
		return
	}
	newSubcategory.ProjectId, err = uuid.Parse(projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing project ID", err)
		return
	}

	if newSubcategory.Id == uuid.Nil {
		newSubcategory.Id = uuid.New()
	}

	if err := validation.CreateSubcategoryValidation(&newSubcategory); err != nil {
		utils.SendValidationError(c, "Validation failed", err)
		return
	}

	err = r.handler.HandlerSubcategory.CreateSubcategory(&newSubcategory)
	if err != nil {
		utils.SendInternalServerError(c, "Error creating subcategory", err)
		return
	}

	utils.SendCreatedSuccess(c, "Subcategory created successfully", newSubcategory)
}

func (r *ResourceSubcategory) ServiceUpdateSubcategory(c *gin.Context) {
	idStr := c.Param("id")

	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid subcategory ID format", err)
		return
	}

	var updatedSubcategory models.Subcategory
	err = c.BindJSON(&updatedSubcategory)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	updatedSubcategory.OrganizationId, err = uuid.Parse(organizationId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing organization ID", err)
		return
	}
	updatedSubcategory.ProjectId, err = uuid.Parse(projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing project ID", err)
		return
	}
	updatedSubcategory.Id, err = uuid.Parse(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing subcategory ID", err)
		return
	}

	if err := validation.UpdateSubcategoryValidation(&updatedSubcategory); err != nil {
		utils.SendValidationError(c, "Validation failed", err)
		return
	}

	err = r.handler.HandlerSubcategory.UpdateSubcategory(&updatedSubcategory)
	if err != nil {
		utils.SendInternalServerError(c, "Error updating subcategory", err)
		return
	}

	utils.SendOKSuccess(c, "Subcategory updated successfully", updatedSubcategory)
}

func (r *ResourceSubcategory) ServiceUpdateSubcategoryOrder(c *gin.Context) {
	idStr := c.Param("id")

	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid subcategory ID format", err)
		return
	}

	var requestBody struct {
		Order int `json:"order" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	err = r.handler.HandlerSubcategory.UpdateSubcategoryOrder(idStr, requestBody.Order)
	if err != nil {
		utils.SendInternalServerError(c, "Error updating subcategory order", err)
		return
	}

	utils.SendOKSuccess(c, "Subcategory order updated successfully", nil)
}

func (r *ResourceSubcategory) ServiceUpdateSubcategoryStatus(c *gin.Context) {
	idStr := c.Param("id")

	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid subcategory ID format", err)
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

	err = r.handler.HandlerSubcategory.UpdateSubcategoryStatus(idStr, *requestBody.Active)
	if err != nil {
		utils.SendInternalServerError(c, "Error updating subcategory status", err)
		return
	}

	utils.SendOKSuccess(c, "Subcategory status updated successfully", nil)
}

func (r *ResourceSubcategory) ServiceDeleteSubcategory(c *gin.Context) {
	idStr := c.Param("id")

	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid subcategory ID format", err)
		return
	}

	err = r.handler.HandlerSubcategory.DeleteSubcategory(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error deleting subcategory", err)
		return
	}

	utils.SendOKSuccess(c, "Subcategory deleted successfully", nil)
}

func (r *ResourceSubcategory) ServiceAddCategoryToSubcategory(c *gin.Context) {
	subcategoryId := c.Param("id")

	_, err := uuid.Parse(subcategoryId)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid subcategory ID format", err)
		return
	}

	var requestBody struct {
		CategoryId string `json:"category_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	_, err = uuid.Parse(requestBody.CategoryId)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid category ID format", err)
		return
	}

	err = r.handler.HandlerSubcategory.AddCategoryToSubcategory(subcategoryId, requestBody.CategoryId)
	if err != nil {
		utils.SendInternalServerError(c, "Error adding category to subcategory", err)
		return
	}

	utils.SendCreatedSuccess(c, "Category added to subcategory successfully", nil)
}

func (r *ResourceSubcategory) ServiceRemoveCategoryFromSubcategory(c *gin.Context) {
	subcategoryId := c.Param("id")
	categoryId := c.Param("categoryId")

	_, err := uuid.Parse(subcategoryId)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid subcategory ID format", err)
		return
	}

	_, err = uuid.Parse(categoryId)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid category ID format", err)
		return
	}

	err = r.handler.HandlerSubcategory.RemoveCategoryFromSubcategory(subcategoryId, categoryId)
	if err != nil {
		utils.SendInternalServerError(c, "Error removing category from subcategory", err)
		return
	}

	utils.SendOKSuccess(c, "Category removed from subcategory successfully", nil)
}

func (r *ResourceSubcategory) ServiceGetSubcategoryCategories(c *gin.Context) {
	subcategoryId := c.Param("id")

	_, err := uuid.Parse(subcategoryId)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid subcategory ID format", err)
		return
	}

	resp, err := r.handler.HandlerSubcategory.GetSubcategoryCategories(subcategoryId)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting subcategory categories", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func NewSourceServerSubcategory(handler *handler.Handlers) IServerSubcategory {
	return &ResourceSubcategory{handler: handler}
}
