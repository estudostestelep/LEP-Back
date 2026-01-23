package client

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupReportsRoutes configura rotas de relatórios para client
func SetupReportsRoutes(r gin.IRouter) {
	reports := r.Group("/reports")
	{
		reports.GET("/occupancy", resource.ServersControllers.SourceReports.GetOccupancyReport)
		reports.GET("/reservations", resource.ServersControllers.SourceReports.GetReservationReport)
		reports.GET("/waitlist", resource.ServersControllers.SourceReports.GetWaitlistReport)
		reports.GET("/leads", resource.ServersControllers.SourceReports.GetLeadReport)
		reports.GET("/export/:type", resource.ServersControllers.SourceReports.ExportReportToCSV)
	}
}
