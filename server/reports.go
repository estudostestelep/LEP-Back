package server

import (
	"lep/handler"
	"lep/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ReportsServer struct {
	handler handler.IReportsHandler
}

type IReportsServer interface {
	GetOccupancyReport(c *gin.Context)
	GetReservationReport(c *gin.Context)
	GetWaitlistReport(c *gin.Context)
	GetLeadReport(c *gin.Context)
	ExportReportToCSV(c *gin.Context)
}

func NewReportsServer(handler handler.IReportsHandler) IReportsServer {
	return &ReportsServer{handler: handler}
}

func (r *ReportsServer) GetOccupancyReport(c *gin.Context) {
	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	// Parâmetros de query para datas
	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid start_date format", err)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid end_date format", err)
		return
	}

	report, err := r.handler.GetOccupancyReport(organizationId, projectId, startDate, endDate)
	if err != nil {
		utils.SendInternalServerError(c, "Error generating occupancy report", err)
		return
	}

	c.JSON(http.StatusOK, report)
}

func (r *ReportsServer) GetReservationReport(c *gin.Context) {
	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	// Parâmetros de query para datas
	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid start_date format", err)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid end_date format", err)
		return
	}

	report, err := r.handler.GetReservationReport(organizationId, projectId, startDate, endDate)
	if err != nil {
		utils.SendInternalServerError(c, "Error generating reservation report", err)
		return
	}

	c.JSON(http.StatusOK, report)
}

func (r *ReportsServer) GetWaitlistReport(c *gin.Context) {
	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	// Parâmetros de query para datas
	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid start_date format", err)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid end_date format", err)
		return
	}

	report, err := r.handler.GetWaitlistReport(organizationId, projectId, startDate, endDate)
	if err != nil {
		utils.SendInternalServerError(c, "Error generating waitlist report", err)
		return
	}

	c.JSON(http.StatusOK, report)
}

func (r *ReportsServer) GetLeadReport(c *gin.Context) {
	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	// Parâmetros de query para datas
	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid start_date format", err)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid end_date format", err)
		return
	}

	report, err := r.handler.GetLeadReport(organizationId, projectId, startDate, endDate)
	if err != nil {
		utils.SendInternalServerError(c, "Error generating lead report", err)
		return
	}

	c.JSON(http.StatusOK, report)
}

func (r *ReportsServer) ExportReportToCSV(c *gin.Context) {
	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	// Parâmetros obrigatórios
	reportType := c.Param("type")
	if reportType == "" {
		utils.SendBadRequestError(c, "Report type is required", nil)
		return
	}

	// Parâmetros de query para datas
	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid start_date format", err)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid end_date format", err)
		return
	}

	csvData, err := r.handler.ExportReportToCSV(reportType, organizationId, projectId, startDate, endDate)
	if err != nil {
		utils.SendInternalServerError(c, "Error exporting report to CSV", err)
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename="+reportType+"_report.csv")
	c.Data(http.StatusOK, "text/csv", csvData)
}