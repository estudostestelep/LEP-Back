package server

import (
	"lep/handler"
	"lep/repositories/models"
	"lep/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ResourceStaffAttendance struct {
	handler *handler.Handlers
}

type IServerStaffAttendance interface {
	// Attendance
	ServiceGetById(c *gin.Context)
	ServiceListByDateRange(c *gin.Context)
	ServiceCreate(c *gin.Context)
	ServiceDelete(c *gin.Context)
	// Consumption Products
	ServiceListConsumptionProducts(c *gin.Context)
	ServiceCreateConsumptionProduct(c *gin.Context)
	ServiceUpdateConsumptionProduct(c *gin.Context)
	ServiceDeleteConsumptionProduct(c *gin.Context)
}

func NewStaffAttendanceServer(h handler.IHandlerStaffAttendance) IServerStaffAttendance {
	return &ResourceStaffAttendance{handler: &handler.Handlers{HandlerStaffAttendance: h}}
}

// ==================== Attendance ====================

func (r *ResourceStaffAttendance) ServiceGetById(c *gin.Context) {
	id := c.Param("id")

	resp, err := r.handler.HandlerStaffAttendance.GetByIdWithDetails(id)
	if err != nil {
		utils.SendNotFoundError(c, "Attendance")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceStaffAttendance) ServiceListByDateRange(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		utils.SendBadRequestError(c, "start_date and end_date are required", nil)
		return
	}

	resp, err := r.handler.HandlerStaffAttendance.ListByDateRange(orgId, projectId, startDate, endDate)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing attendances", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceStaffAttendance) ServiceCreate(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	var req models.CreateStaffAttendanceRequest
	if err := c.BindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	resp, err := r.handler.HandlerStaffAttendance.CreateWithDetails(&req, orgId, projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error creating attendance", err)
		return
	}

	utils.SendCreatedSuccess(c, "Attendance created successfully", resp)
}

func (r *ResourceStaffAttendance) ServiceDelete(c *gin.Context) {
	id := c.Param("id")

	err := r.handler.HandlerStaffAttendance.Delete(id)
	if err != nil {
		utils.SendInternalServerError(c, "Error deleting attendance", err)
		return
	}

	utils.SendOKSuccess(c, "Attendance deleted successfully", nil)
}

// ==================== Consumption Products ====================

func (r *ResourceStaffAttendance) ServiceListConsumptionProducts(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	resp, err := r.handler.HandlerStaffAttendance.ListConsumptionProducts(orgId, projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing products", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceStaffAttendance) ServiceCreateConsumptionProduct(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	var req models.StaffConsumptionProduct
	if err := c.BindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	orgUUID, _ := uuid.Parse(orgId)
	projUUID, _ := uuid.Parse(projectId)
	req.OrganizationId = orgUUID
	req.ProjectId = projUUID

	err := r.handler.HandlerStaffAttendance.CreateConsumptionProduct(&req)
	if err != nil {
		utils.SendInternalServerError(c, "Error creating product", err)
		return
	}

	utils.SendCreatedSuccess(c, "Product created successfully", req)
}

func (r *ResourceStaffAttendance) ServiceUpdateConsumptionProduct(c *gin.Context) {
	id := c.Param("id")

	existing, err := r.handler.HandlerStaffAttendance.GetConsumptionProductById(id)
	if err != nil {
		utils.SendNotFoundError(c, "Product")
		return
	}

	var req models.StaffConsumptionProduct
	if err := c.BindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	existing.Name = req.Name
	existing.Category = req.Category
	existing.UnitCost = req.UnitCost
	existing.Active = req.Active

	err = r.handler.HandlerStaffAttendance.UpdateConsumptionProduct(existing)
	if err != nil {
		utils.SendInternalServerError(c, "Error updating product", err)
		return
	}

	utils.SendOKSuccess(c, "Product updated successfully", existing)
}

func (r *ResourceStaffAttendance) ServiceDeleteConsumptionProduct(c *gin.Context) {
	id := c.Param("id")

	err := r.handler.HandlerStaffAttendance.DeleteConsumptionProduct(id)
	if err != nil {
		utils.SendInternalServerError(c, "Error deleting product", err)
		return
	}

	utils.SendOKSuccess(c, "Product deleted successfully", nil)
}
