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
	id, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "table")
	if !ok {
		return
	}

	resp, err := r.handler.HandlerTables.GetTable(id.String())
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
	id, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "table")
	if !ok {
		return
	}

	var updatedTable models.Table
	if err := c.BindJSON(&updatedTable); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	var err error
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
	updatedTable.Id = id

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
	id, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "table")
	if !ok {
		return
	}

	if err := r.handler.HandlerTables.DeleteTable(id.String()); err != nil {
		utils.SendInternalServerError(c, "Error deleting table", err)
		return
	}

	utils.SendOKSuccess(c, "Table deleted successfully", nil)
}

func (r *ResourceTables) ServiceListTables(c *gin.Context) {
	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	var environmentId *string
	if envId := c.Query("environment_id"); envId != "" {
		environmentId = &envId
	}

	resp, err := r.handler.HandlerTables.ListTables(organizationId, projectId, environmentId)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing tables", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func NewSourceServerTables(handler *handler.Handlers) IServerTables {
	return &ResourceTables{handler: handler}
}