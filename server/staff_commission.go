package server

import (
	"lep/handler"
	"lep/repositories/models"
	"lep/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResourceStaffCommission struct {
	handler *handler.Handlers
}

type IServerStaffCommission interface {
	ServiceGetById(c *gin.Context)
	ServiceListByDateRange(c *gin.Context)
	ServiceCreate(c *gin.Context)
	ServiceUpdate(c *gin.Context)
	ServiceDelete(c *gin.Context)
	ServiceGetSummary(c *gin.Context)
}

func NewStaffCommissionServer(h handler.IHandlerStaffCommission) IServerStaffCommission {
	return &ResourceStaffCommission{handler: &handler.Handlers{HandlerStaffCommission: h}}
}

func (r *ResourceStaffCommission) ServiceGetById(c *gin.Context) {
	id := c.Param("id")

	resp, err := r.handler.HandlerStaffCommission.GetById(id)
	if err != nil {
		utils.SendNotFoundError(c, "Commission")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceStaffCommission) ServiceListByDateRange(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		utils.SendBadRequestError(c, "start_date and end_date are required", nil)
		return
	}

	resp, err := r.handler.HandlerStaffCommission.ListByDateRange(orgId, projectId, startDate, endDate)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing commissions", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceStaffCommission) ServiceCreate(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	var req models.CreateCommissionRequest
	if err := c.BindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	resp, err := r.handler.HandlerStaffCommission.Create(&req, orgId, projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error creating commission", err)
		return
	}

	utils.SendCreatedSuccess(c, "Commission created successfully", resp)
}

func (r *ResourceStaffCommission) ServiceUpdate(c *gin.Context) {
	id := c.Param("id")

	existing, err := r.handler.HandlerStaffCommission.GetById(id)
	if err != nil {
		utils.SendNotFoundError(c, "Commission")
		return
	}

	var req models.StaffDailyCommission
	if err := c.BindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	existing.CommissionValue = req.CommissionValue
	existing.Revenue = req.Revenue

	err = r.handler.HandlerStaffCommission.Update(existing)
	if err != nil {
		utils.SendInternalServerError(c, "Error updating commission", err)
		return
	}

	utils.SendOKSuccess(c, "Commission updated successfully", existing)
}

func (r *ResourceStaffCommission) ServiceDelete(c *gin.Context) {
	id := c.Param("id")

	err := r.handler.HandlerStaffCommission.Delete(id)
	if err != nil {
		utils.SendInternalServerError(c, "Error deleting commission", err)
		return
	}

	utils.SendOKSuccess(c, "Commission deleted successfully", nil)
}

func (r *ResourceStaffCommission) ServiceGetSummary(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		utils.SendBadRequestError(c, "start_date and end_date are required", nil)
		return
	}

	resp, err := r.handler.HandlerStaffCommission.GetSummary(orgId, projectId, startDate, endDate)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting summary", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
