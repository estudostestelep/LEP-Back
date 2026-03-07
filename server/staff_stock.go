package server

import (
	"lep/handler"
	"lep/repositories/models"
	"lep/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ResourceStaffStock struct {
	handler *handler.Handlers
}

type IServerStaffStock interface {
	// Items
	ServiceListItems(c *gin.Context)
	ServiceListItemsBySector(c *gin.Context)
	ServiceListSectors(c *gin.Context)
	ServiceCreateItem(c *gin.Context)
	ServiceUpdateItem(c *gin.Context)
	ServiceDeleteItem(c *gin.Context)
	// Records
	ServiceListRecords(c *gin.Context)
	ServiceGetRecordById(c *gin.Context)
	ServiceCreateRecord(c *gin.Context)
	ServiceGenerateShoppingList(c *gin.Context)
}

func NewStaffStockServer(h handler.IHandlerStaffStock) IServerStaffStock {
	return &ResourceStaffStock{handler: &handler.Handlers{HandlerStaffStock: h}}
}

// ==================== Stock Items ====================

func (r *ResourceStaffStock) ServiceListItems(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")
	sector := c.Query("sector")

	var resp []models.StaffStockItem
	var err error

	if sector != "" {
		resp, err = r.handler.HandlerStaffStock.ListItemsBySector(orgId, projectId, sector)
	} else {
		resp, err = r.handler.HandlerStaffStock.ListItems(orgId, projectId)
	}

	if err != nil {
		utils.SendInternalServerError(c, "Error listing items", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceStaffStock) ServiceListItemsBySector(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")
	sector := c.Param("sector")

	resp, err := r.handler.HandlerStaffStock.ListItemsBySector(orgId, projectId, sector)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing items", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceStaffStock) ServiceListSectors(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	resp, err := r.handler.HandlerStaffStock.ListSectors(orgId, projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing sectors", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceStaffStock) ServiceCreateItem(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	var req models.StaffStockItem
	if err := c.BindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	orgUUID, _ := uuid.Parse(orgId)
	projUUID, _ := uuid.Parse(projectId)
	req.OrganizationId = orgUUID
	req.ProjectId = projUUID

	err := r.handler.HandlerStaffStock.CreateItem(&req)
	if err != nil {
		utils.SendInternalServerError(c, "Error creating item", err)
		return
	}

	utils.SendCreatedSuccess(c, "Item created successfully", req)
}

func (r *ResourceStaffStock) ServiceUpdateItem(c *gin.Context) {
	id := c.Param("id")

	existing, err := r.handler.HandlerStaffStock.GetItemById(id)
	if err != nil {
		utils.SendNotFoundError(c, "Item")
		return
	}

	var req models.StaffStockItem
	if err := c.BindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	existing.Name = req.Name
	existing.Category = req.Category
	existing.Storage = req.Storage
	existing.StockMin = req.StockMin
	existing.StockMax = req.StockMax
	existing.Sector = req.Sector
	existing.WhereToBuy = req.WhereToBuy
	existing.Notes = req.Notes
	existing.Active = req.Active
	existing.Order = req.Order

	err = r.handler.HandlerStaffStock.UpdateItem(existing)
	if err != nil {
		utils.SendInternalServerError(c, "Error updating item", err)
		return
	}

	utils.SendOKSuccess(c, "Item updated successfully", existing)
}

func (r *ResourceStaffStock) ServiceDeleteItem(c *gin.Context) {
	id := c.Param("id")

	err := r.handler.HandlerStaffStock.DeleteItem(id)
	if err != nil {
		utils.SendInternalServerError(c, "Error deleting item", err)
		return
	}

	utils.SendOKSuccess(c, "Item deleted successfully", nil)
}

// ==================== Stock Records ====================

func (r *ResourceStaffStock) ServiceListRecords(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")
	sector := c.Query("sector")
	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)

	var resp []models.StaffStockRecord
	var err error

	if sector != "" {
		resp, err = r.handler.HandlerStaffStock.ListRecordsBySector(orgId, projectId, sector, limit)
	} else {
		resp, err = r.handler.HandlerStaffStock.ListRecords(orgId, projectId, limit)
	}

	if err != nil {
		utils.SendInternalServerError(c, "Error listing records", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceStaffStock) ServiceGetRecordById(c *gin.Context) {
	id := c.Param("id")

	resp, err := r.handler.HandlerStaffStock.GetRecordByIdWithItems(id)
	if err != nil {
		utils.SendNotFoundError(c, "Record")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceStaffStock) ServiceCreateRecord(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")
	userId := c.GetString("user_id")

	var req models.CreateStockRecordRequest
	if err := c.BindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	resp, err := r.handler.HandlerStaffStock.CreateRecordWithItems(&req, orgId, projectId, userId)
	if err != nil {
		utils.SendInternalServerError(c, "Error creating record", err)
		return
	}

	utils.SendCreatedSuccess(c, "Record created successfully", resp)
}

func (r *ResourceStaffStock) ServiceGenerateShoppingList(c *gin.Context) {
	id := c.Param("id")

	resp, err := r.handler.HandlerStaffStock.GenerateShoppingList(id)
	if err != nil {
		utils.SendInternalServerError(c, "Error generating shopping list", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
