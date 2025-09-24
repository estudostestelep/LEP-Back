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

type ResourceCustomer struct {
	handler *handler.Handlers
}

type IServerCustomer interface {
	ServiceGetCustomer(c *gin.Context)
	ServiceCreateCustomer(c *gin.Context)
	ServiceUpdateCustomer(c *gin.Context)
	ServiceDeleteCustomer(c *gin.Context)
	ServiceListCustomers(c *gin.Context)
}

func (r *ResourceCustomer) ServiceGetCustomer(c *gin.Context) {
	idStr := c.Param("id")

	// Validar formato UUID
	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid customer ID format", err)
		return
	}

	resp, err := r.handler.HandlerCustomer.GetCustomer(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting customer", err)
		return
	}

	if resp == nil {
		utils.SendNotFoundError(c, "Customer")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceCustomer) ServiceCreateCustomer(c *gin.Context) {
	var newCustomer models.Customer
	err := c.BindJSON(&newCustomer)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	newCustomer.OrganizationId, err = uuid.Parse(organizationId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing organization ID", err)
		return
	}
	newCustomer.ProjectId, err = uuid.Parse(projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing project ID", err)
		return
	}

	// Gerar ID se não fornecido
	if newCustomer.Id == uuid.Nil {
		newCustomer.Id = uuid.New()
	}

	// Validações estruturadas
	if err := validation.CreateCustomerValidation(&newCustomer); err != nil {
		utils.SendValidationError(c, "Validation failed", err)
		return
	}

	err = r.handler.HandlerCustomer.CreateCustomer(&newCustomer)
	if err != nil {
		utils.SendInternalServerError(c, "Error creating customer", err)
		return
	}

	utils.SendCreatedSuccess(c, "Customer created successfully", newCustomer)
}

func (r *ResourceCustomer) ServiceUpdateCustomer(c *gin.Context) {
	idStr := c.Param("id")

	// Validar formato UUID
	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid customer ID format", err)
		return
	}

	var updatedCustomer models.Customer
	err = c.BindJSON(&updatedCustomer)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	updatedCustomer.OrganizationId, err = uuid.Parse(organizationId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing organization ID", err)
		return
	}
	updatedCustomer.ProjectId, err = uuid.Parse(projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing project ID", err)
		return
	}
	updatedCustomer.Id, err = uuid.Parse(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing customer ID", err)
		return
	}

	// Validações estruturadas
	if err := validation.UpdateCustomerValidation(&updatedCustomer); err != nil {
		utils.SendValidationError(c, "Validation failed", err)
		return
	}

	err = r.handler.HandlerCustomer.UpdateCustomer(&updatedCustomer)
	if err != nil {
		utils.SendInternalServerError(c, "Error updating customer", err)
		return
	}

	utils.SendOKSuccess(c, "Customer updated successfully", updatedCustomer)
}

func (r *ResourceCustomer) ServiceDeleteCustomer(c *gin.Context) {
	idStr := c.Param("id")

	// Validar formato UUID
	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid customer ID format", err)
		return
	}

	err = r.handler.HandlerCustomer.DeleteCustomer(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error deleting customer", err)
		return
	}

	utils.SendOKSuccess(c, "Customer deleted successfully", nil)
}

func (r *ResourceCustomer) ServiceListCustomers(c *gin.Context) {
	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	resp, err := r.handler.HandlerCustomer.ListCustomers(organizationId, projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing customers", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func NewSourceServerCustomer(handler *handler.Handlers) IServerCustomer {
	return &ResourceCustomer{handler: handler}
}