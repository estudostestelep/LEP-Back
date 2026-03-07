package server

import (
	"lep/handler"
	"lep/repositories/models"
	"lep/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ResourceStaffDashboard struct {
	handler *handler.Handlers
}

type IServerStaffDashboard interface {
	ServiceGetMeta(c *gin.Context)
	ServiceGetRows(c *gin.Context)
	ServiceGetGraphs(c *gin.Context)
	ServiceImportCSV(c *gin.Context)
	ServiceListImportBatches(c *gin.Context)
	ServiceGetStaffReportMeta(c *gin.Context)
	ServiceGetStaffReportRows(c *gin.Context)
}

func NewStaffDashboardServer(h handler.IHandlerStaffDashboard) IServerStaffDashboard {
	return &ResourceStaffDashboard{handler: &handler.Handlers{HandlerStaffDashboard: h}}
}

func (r *ResourceStaffDashboard) ServiceGetMeta(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	resp, err := r.handler.HandlerStaffDashboard.GetDashboardMeta(orgId, projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting metadata", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceStaffDashboard) ServiceGetRows(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	filters := models.DashboardFilters{
		StartDate:    getStringPtr(c.Query("start_date")),
		EndDate:      getStringPtr(c.Query("end_date")),
		Groups:       c.QueryArray("groups"),
		Descriptions: c.QueryArray("descriptions"),
		Employees:    c.QueryArray("employees"),
		Weekdays:     c.QueryArray("weekdays"),
	}

	resp, err := r.handler.HandlerStaffDashboard.GetDashboardRows(orgId, projectId, filters)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting rows", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceStaffDashboard) ServiceGetGraphs(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	resp, err := r.handler.HandlerStaffDashboard.GetDashboardGraphs(orgId, projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting graphs", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceStaffDashboard) ServiceImportCSV(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")
	userId := c.GetString("user_id")

	var req models.ImportCSVRequest
	if err := c.BindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	resp, err := r.handler.HandlerStaffDashboard.ImportCSV(&req, orgId, projectId, userId)
	if err != nil {
		utils.SendInternalServerError(c, "Error importing CSV", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceStaffDashboard) ServiceListImportBatches(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	resp, err := r.handler.HandlerStaffDashboard.ListImportBatches(orgId, projectId, limit)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing batches", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceStaffDashboard) ServiceGetStaffReportMeta(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	resp, err := r.handler.HandlerStaffDashboard.GetStaffReportMeta(orgId, projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting metadata", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceStaffDashboard) ServiceGetStaffReportRows(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	yearStr := c.Query("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		utils.SendBadRequestError(c, "year is required", nil)
		return
	}

	filters := models.StaffReportFilters{
		Year:      year,
		StartDate: getStringPtr(c.Query("start_date")),
		EndDate:   getStringPtr(c.Query("end_date")),
		Weekdays:  c.QueryArray("weekdays"),
		Employees: c.QueryArray("employees"),
		Sectors:   c.QueryArray("sectors"),
	}

	resp, err := r.handler.HandlerStaffDashboard.GetStaffReportRows(orgId, projectId, filters)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting rows", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func getStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
