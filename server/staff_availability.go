package server

import (
	"lep/handler"
	"lep/repositories/models"
	"lep/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ResourceStaffAvailability struct {
	handler *handler.Handlers
}

type IServerStaffAvailability interface {
	ServiceGetById(c *gin.Context)
	ServiceListByWeek(c *gin.Context)
	ServiceListByClient(c *gin.Context)
	ServiceUpsert(c *gin.Context)
	ServiceDelete(c *gin.Context)
	ServiceGetWeekSummary(c *gin.Context)
}

func NewStaffAvailabilityServer(h handler.IHandlerStaffAvailability) IServerStaffAvailability {
	return &ResourceStaffAvailability{handler: &handler.Handlers{HandlerStaffAvailability: h}}
}

func (r *ResourceStaffAvailability) ServiceGetById(c *gin.Context) {
	id := c.Param("id")

	resp, err := r.handler.HandlerStaffAvailability.GetById(id)
	if err != nil {
		utils.SendNotFoundError(c, "Availability")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceStaffAvailability) ServiceListByWeek(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")
	weekStart := c.Param("weekStart")

	resp, err := r.handler.HandlerStaffAvailability.ListByWeek(orgId, projectId, weekStart)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing availabilities", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceStaffAvailability) ServiceListByClient(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")
	clientId := c.Param("clientId")

	resp, err := r.handler.HandlerStaffAvailability.ListByClient(orgId, projectId, clientId)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing availabilities", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceStaffAvailability) ServiceUpsert(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	var req models.StaffAvailability
	if err := c.BindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Set org and project from context
	orgUUID, _ := uuid.Parse(orgId)
	projUUID, _ := uuid.Parse(projectId)
	req.OrganizationId = orgUUID
	req.ProjectId = projUUID

	err := r.handler.HandlerStaffAvailability.Upsert(&req)
	if err != nil {
		utils.SendInternalServerError(c, "Error saving availability", err)
		return
	}

	utils.SendCreatedSuccess(c, "Availability saved successfully", req)
}

func (r *ResourceStaffAvailability) ServiceDelete(c *gin.Context) {
	id := c.Param("id")

	err := r.handler.HandlerStaffAvailability.Delete(id)
	if err != nil {
		utils.SendInternalServerError(c, "Error deleting availability", err)
		return
	}

	utils.SendOKSuccess(c, "Availability deleted successfully", nil)
}

func (r *ResourceStaffAvailability) ServiceGetWeekSummary(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")
	weekStart := c.Param("weekStart")

	resp, err := r.handler.HandlerStaffAvailability.GetWeekSummary(orgId, projectId, weekStart)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting summary", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
