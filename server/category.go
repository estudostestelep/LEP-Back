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

type ResourceCategory struct {
	handler *handler.Handlers
}

type IServerCategory interface {
	ServiceGetCategory(c *gin.Context)
	ServiceListCategories(c *gin.Context)
	ServiceGetCategoriesByMenu(c *gin.Context)
	ServiceListActiveCategories(c *gin.Context)
	ServiceCreateCategory(c *gin.Context)
	ServiceUpdateCategory(c *gin.Context)
	ServiceUpdateCategoryOrder(c *gin.Context)
	ServiceUpdateCategoryStatus(c *gin.Context)
	ServiceDeleteCategory(c *gin.Context)
}

func (r *ResourceCategory) ServiceGetCategory(c *gin.Context) {
	idStr := c.Param("id")

	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid category ID format", err)
		return
	}

	resp, err := r.handler.HandlerCategory.GetCategory(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting category", err)
		return
	}

	if resp == nil {
		utils.SendNotFoundError(c, "Category")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceCategory) ServiceListCategories(c *gin.Context) {
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	resp, err := r.handler.HandlerCategory.ListCategories(organizationId, projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing categories", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceCategory) ServiceGetCategoriesByMenu(c *gin.Context) {
	menuId := c.Param("menuId")

	_, err := uuid.Parse(menuId)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid menu ID format", err)
		return
	}

	resp, err := r.handler.HandlerCategory.GetCategoriesByMenu(menuId)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting categories by menu", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceCategory) ServiceListActiveCategories(c *gin.Context) {
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	resp, err := r.handler.HandlerCategory.ListActiveCategories(organizationId, projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing active categories", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceCategory) ServiceCreateCategory(c *gin.Context) {
	var newCategory models.Category
	err := c.BindJSON(&newCategory)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	newCategory.OrganizationId, err = uuid.Parse(organizationId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing organization ID", err)
		return
	}
	newCategory.ProjectId, err = uuid.Parse(projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing project ID", err)
		return
	}

	if newCategory.Id == uuid.Nil {
		newCategory.Id = uuid.New()
	}

	if err := validation.CreateCategoryValidation(&newCategory); err != nil {
		utils.SendValidationError(c, "Validation failed", err)
		return
	}

	err = r.handler.HandlerCategory.CreateCategory(&newCategory)
	if err != nil {
		utils.SendInternalServerError(c, "Error creating category", err)
		return
	}

	utils.SendCreatedSuccess(c, "Category created successfully", newCategory)
}

func (r *ResourceCategory) ServiceUpdateCategory(c *gin.Context) {
	idStr := c.Param("id")

	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid category ID format", err)
		return
	}

	var updatedCategory models.Category
	err = c.BindJSON(&updatedCategory)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	updatedCategory.OrganizationId, err = uuid.Parse(organizationId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing organization ID", err)
		return
	}
	updatedCategory.ProjectId, err = uuid.Parse(projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing project ID", err)
		return
	}
	updatedCategory.Id, err = uuid.Parse(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing category ID", err)
		return
	}

	if err := validation.UpdateCategoryValidation(&updatedCategory); err != nil {
		utils.SendValidationError(c, "Validation failed", err)
		return
	}

	err = r.handler.HandlerCategory.UpdateCategory(&updatedCategory)
	if err != nil {
		utils.SendInternalServerError(c, "Error updating category", err)
		return
	}

	utils.SendOKSuccess(c, "Category updated successfully", updatedCategory)
}

func (r *ResourceCategory) ServiceUpdateCategoryOrder(c *gin.Context) {
	idStr := c.Param("id")

	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid category ID format", err)
		return
	}

	var requestBody struct {
		Order int `json:"order" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	err = r.handler.HandlerCategory.UpdateCategoryOrder(idStr, requestBody.Order)
	if err != nil {
		utils.SendInternalServerError(c, "Error updating category order", err)
		return
	}

	utils.SendOKSuccess(c, "Category order updated successfully", nil)
}

func (r *ResourceCategory) ServiceUpdateCategoryStatus(c *gin.Context) {
	idStr := c.Param("id")

	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid category ID format", err)
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

	err = r.handler.HandlerCategory.UpdateCategoryStatus(idStr, *requestBody.Active)
	if err != nil {
		utils.SendInternalServerError(c, "Error updating category status", err)
		return
	}

	utils.SendOKSuccess(c, "Category status updated successfully", nil)
}

func (r *ResourceCategory) ServiceDeleteCategory(c *gin.Context) {
	idStr := c.Param("id")

	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid category ID format", err)
		return
	}

	err = r.handler.HandlerCategory.DeleteCategory(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error deleting category", err)
		return
	}

	utils.SendOKSuccess(c, "Category deleted successfully", nil)
}

func NewSourceServerCategory(handler *handler.Handlers) IServerCategory {
	return &ResourceCategory{handler: handler}
}
