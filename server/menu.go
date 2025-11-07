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

type ResourceMenu struct {
	handler *handler.Handlers
}

type IServerMenu interface {
	ServiceGetMenu(c *gin.Context)
	ServiceListMenus(c *gin.Context)
	ServiceListActiveMenus(c *gin.Context)
	ServiceCreateMenu(c *gin.Context)
	ServiceUpdateMenu(c *gin.Context)
	ServiceUpdateMenuOrder(c *gin.Context)
	ServiceUpdateMenuStatus(c *gin.Context)
	ServiceDeleteMenu(c *gin.Context)

	// ✨ Novos endpoints para seleção inteligente de cardápio
	ServiceGetMenuOptions(c *gin.Context)
	ServiceGetActiveMenu(c *gin.Context)
	ServiceSetMenuAsManualOverride(c *gin.Context)
	ServiceRemoveManualOverride(c *gin.Context)
}

func (r *ResourceMenu) ServiceGetMenu(c *gin.Context) {
	idStr := c.Param("id")

	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid menu ID format", err)
		return
	}

	resp, err := r.handler.HandlerMenu.GetMenu(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting menu", err)
		return
	}

	if resp == nil {
		utils.SendNotFoundError(c, "Menu")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceMenu) ServiceListMenus(c *gin.Context) {
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	resp, err := r.handler.HandlerMenu.ListMenus(organizationId, projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing menus", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceMenu) ServiceListActiveMenus(c *gin.Context) {
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	resp, err := r.handler.HandlerMenu.ListActiveMenus(organizationId, projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing active menus", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceMenu) ServiceCreateMenu(c *gin.Context) {
	var newMenu models.Menu
	err := c.BindJSON(&newMenu)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	newMenu.OrganizationId, err = uuid.Parse(organizationId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing organization ID", err)
		return
	}
	newMenu.ProjectId, err = uuid.Parse(projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing project ID", err)
		return
	}

	if newMenu.Id == uuid.Nil {
		newMenu.Id = uuid.New()
	}

	if err := validation.CreateMenuValidation(&newMenu); err != nil {
		utils.SendValidationError(c, "Validation failed", err)
		return
	}

	err = r.handler.HandlerMenu.CreateMenu(&newMenu)
	if err != nil {
		utils.SendInternalServerError(c, "Error creating menu", err)
		return
	}

	utils.SendCreatedSuccess(c, "Menu created successfully", newMenu)
}

func (r *ResourceMenu) ServiceUpdateMenu(c *gin.Context) {
	idStr := c.Param("id")

	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid menu ID format", err)
		return
	}

	var updatedMenu models.Menu
	err = c.BindJSON(&updatedMenu)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	updatedMenu.OrganizationId, err = uuid.Parse(organizationId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing organization ID", err)
		return
	}
	updatedMenu.ProjectId, err = uuid.Parse(projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing project ID", err)
		return
	}
	updatedMenu.Id, err = uuid.Parse(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing menu ID", err)
		return
	}

	if err := validation.UpdateMenuValidation(&updatedMenu); err != nil {
		utils.SendValidationError(c, "Validation failed", err)
		return
	}

	err = r.handler.HandlerMenu.UpdateMenu(&updatedMenu)
	if err != nil {
		utils.SendInternalServerError(c, "Error updating menu", err)
		return
	}

	utils.SendOKSuccess(c, "Menu updated successfully", updatedMenu)
}

func (r *ResourceMenu) ServiceUpdateMenuOrder(c *gin.Context) {
	idStr := c.Param("id")

	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid menu ID format", err)
		return
	}

	var requestBody struct {
		Order int `json:"order" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	err = r.handler.HandlerMenu.UpdateMenuOrder(idStr, requestBody.Order)
	if err != nil {
		utils.SendInternalServerError(c, "Error updating menu order", err)
		return
	}

	utils.SendOKSuccess(c, "Menu order updated successfully", nil)
}

func (r *ResourceMenu) ServiceUpdateMenuStatus(c *gin.Context) {
	idStr := c.Param("id")

	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid menu ID format", err)
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

	err = r.handler.HandlerMenu.UpdateMenuStatus(idStr, *requestBody.Active)
	if err != nil {
		utils.SendInternalServerError(c, "Error updating menu status", err)
		return
	}

	utils.SendOKSuccess(c, "Menu status updated successfully", nil)
}

func (r *ResourceMenu) ServiceDeleteMenu(c *gin.Context) {
	idStr := c.Param("id")

	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid menu ID format", err)
		return
	}

	err = r.handler.HandlerMenu.DeleteMenu(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error deleting menu", err)
		return
	}

	utils.SendOKSuccess(c, "Menu deleted successfully", nil)
}

// ✨ ServiceGetMenuOptions retorna lista de opções de cardápio
func (r *ResourceMenu) ServiceGetMenuOptions(c *gin.Context) {
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	resp, err := r.handler.HandlerMenu.GetMenuOptions(organizationId, projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing menu options", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ✨ ServiceGetActiveMenu retorna o cardápio ativo com lógica de seleção automática
func (r *ResourceMenu) ServiceGetActiveMenu(c *gin.Context) {
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	resp, err := r.handler.HandlerMenu.GetActiveMenuByTimeRange(organizationId, projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting active menu", err)
		return
	}

	if resp == nil {
		utils.SendNotFoundError(c, "Active menu")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ✨ ServiceSetMenuAsManualOverride define um cardápio como override manual
func (r *ResourceMenu) ServiceSetMenuAsManualOverride(c *gin.Context) {
	idStr := c.Param("id")

	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid menu ID format", err)
		return
	}

	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	err = r.handler.HandlerMenu.SetMenuAsManualOverride(organizationId, projectId, idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error setting menu as manual override", err)
		return
	}

	utils.SendOKSuccess(c, "Menu set as manual override successfully", nil)
}

// ✨ ServiceRemoveManualOverride remove o override manual
func (r *ResourceMenu) ServiceRemoveManualOverride(c *gin.Context) {
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	err := r.handler.HandlerMenu.RemoveManualOverride(organizationId, projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error removing manual override", err)
		return
	}

	utils.SendOKSuccess(c, "Manual override removed successfully", nil)
}

func NewSourceServerMenu(handler *handler.Handlers) IServerMenu {
	return &ResourceMenu{handler: handler}
}
