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

type ResourceTables struct {
	handler *handler.Handlers
}

type IServerTables interface {
	ServiceGetTable(c *gin.Context)
	ServiceCreateTable(c *gin.Context)
	ServiceUpdateTable(c *gin.Context)
	ServiceDeleteTable(c *gin.Context)
	ServiceListTables(c *gin.Context)
}

func (r *ResourceTables) ServiceGetTable(c *gin.Context) {
	idStr := c.Param("id")

	// Validar formato UUID
	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid table ID format", err)
		return
	}

	resp, err := r.handler.HandlerTables.GetTable(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting table", err)
		return
	}

	if resp == nil {
		utils.SendNotFoundError(c, "Table")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceTables) ServiceCreateTable(c *gin.Context) {
	var newTable models.Table
	err := c.BindJSON(&newTable)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	newTable.OrganizationId, err = uuid.Parse(organizationId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing organization ID", err)
		return
	}
	newTable.ProjectId, err = uuid.Parse(projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing project ID", err)
		return
	}

	// Gerar ID se não fornecido
	if newTable.Id == uuid.Nil {
		newTable.Id = uuid.New()
	}

	// Validações estruturadas
	if err := validation.CreateTableValidation(&newTable); err != nil {
		utils.SendValidationError(c, "Validation failed", err)
		return
	}

	err = r.handler.HandlerTables.CreateTable(&newTable)
	if err != nil {
		utils.SendInternalServerError(c, "Error creating table", err)
		return
	}

	utils.SendCreatedSuccess(c, "Table created successfully", newTable)
}

func (r *ResourceTables) ServiceUpdateTable(c *gin.Context) {
	idStr := c.Param("id")

	// Validar formato UUID
	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid table ID format", err)
		return
	}

	var updatedTable models.Table
	err = c.BindJSON(&updatedTable)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	updatedTable.OrganizationId, err = uuid.Parse(organizationId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing organization ID", err)
		return
	}
	updatedTable.ProjectId, err = uuid.Parse(projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing project ID", err)
		return
	}
	updatedTable.Id, err = uuid.Parse(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing table ID", err)
		return
	}

	// Validações estruturadas
	if err := validation.UpdateTableValidation(&updatedTable); err != nil {
		utils.SendValidationError(c, "Validation failed", err)
		return
	}

	err = r.handler.HandlerTables.UpdateTable(&updatedTable)
	if err != nil {
		utils.SendInternalServerError(c, "Error updating table", err)
		return
	}

	utils.SendOKSuccess(c, "Table updated successfully", updatedTable)
}

func (r *ResourceTables) ServiceDeleteTable(c *gin.Context) {
	idStr := c.Param("id")

	// Validar formato UUID
	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid table ID format", err)
		return
	}

	err = r.handler.HandlerTables.DeleteTable(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error deleting table", err)
		return
	}

	utils.SendOKSuccess(c, "Table deleted successfully", nil)
}

func (r *ResourceTables) ServiceListTables(c *gin.Context) {
	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	resp, err := r.handler.HandlerTables.ListTables(organizationId, projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing tables", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func NewSourceServerTables(handler *handler.Handlers) IServerTables {
	return &ResourceTables{handler: handler}
}