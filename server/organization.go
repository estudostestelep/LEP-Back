package server

import (
	"lep/handler"
	"lep/repositories/models"
	"lep/resource/validation"
	"lep/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ResourceOrganization struct {
	handler *handler.Handlers
}

type IServerOrganization interface {
	GetOrganizationById(c *gin.Context)
	GetOrganizationByEmail(c *gin.Context)
	ListOrganizations(c *gin.Context)
	ListActiveOrganizations(c *gin.Context)
	CreateOrganization(c *gin.Context)
	UpdateOrganization(c *gin.Context)
	SoftDeleteOrganization(c *gin.Context)
	HardDeleteOrganization(c *gin.Context)
	ServiceCreateOrganizationBootstrap(c *gin.Context)
}

func (r *ResourceOrganization) GetOrganizationById(c *gin.Context) {
	idStr := c.Param("id")

	// Validar formato UUID
	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid organization ID format", err)
		return
	}

	organization, err := r.handler.HandlerOrganization.GetOrganizationById(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting organization", err)
		return
	}

	if organization == nil {
		utils.SendNotFoundError(c, "Organization")
		return
	}

	c.JSON(http.StatusOK, organization)
}

func (r *ResourceOrganization) GetOrganizationByEmail(c *gin.Context) {
	email := c.Query("email")
	if strings.TrimSpace(email) == "" {
		utils.SendBadRequestError(c, "Email parameter is required", nil)
		return
	}

	organization, err := r.handler.HandlerOrganization.GetOrganizationByEmail(email)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting organization by email", err)
		return
	}

	if organization == nil {
		utils.SendNotFoundError(c, "Organization")
		return
	}

	c.JSON(http.StatusOK, organization)
}

func (r *ResourceOrganization) ListOrganizations(c *gin.Context) {
	organizations, err := r.handler.HandlerOrganization.ListOrganizations()
	if err != nil {
		utils.SendInternalServerError(c, "Error listing organizations", err)
		return
	}

	c.JSON(http.StatusOK, organizations)
}

func (r *ResourceOrganization) ListActiveOrganizations(c *gin.Context) {
	organizations, err := r.handler.HandlerOrganization.ListActiveOrganizations()
	if err != nil {
		utils.SendInternalServerError(c, "Error listing active organizations", err)
		return
	}

	c.JSON(http.StatusOK, organizations)
}

func (r *ResourceOrganization) CreateOrganization(c *gin.Context) {
	var newOrganization models.Organization
	err := c.BindJSON(&newOrganization)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Validações estruturadas
	if err := validation.CreateOrganizationValidation(&newOrganization); err != nil {
		utils.SendValidationError(c, "Validation failed", err)
		return
	}

	err = r.handler.HandlerOrganization.CreateOrganization(&newOrganization)
	if err != nil {
		utils.SendInternalServerError(c, "Error creating organization", err)
		return
	}

	utils.SendCreatedSuccess(c, "Organization created successfully", newOrganization)
}

func (r *ResourceOrganization) UpdateOrganization(c *gin.Context) {
	idStr := c.Param("id")

	// Validar formato UUID
	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid organization ID format", err)
		return
	}

	var updateData models.Organization
	err = c.BindJSON(&updateData)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Verificar se a organização existe
	existingOrg, err := r.handler.HandlerOrganization.GetOrganizationById(idStr)
	if err != nil {
		utils.SendNotFoundError(c, "Organization")
		return
	}

	// Preservar o ID da organização existente
	updateData.Id = existingOrg.Id

	// Validações estruturadas
	if err := validation.UpdateOrganizationValidation(&updateData); err != nil {
		utils.SendValidationError(c, "Validation failed", err)
		return
	}

	err = r.handler.HandlerOrganization.UpdateOrganization(&updateData)
	if err != nil {
		utils.SendInternalServerError(c, "Error updating organization", err)
		return
	}

	utils.SendOKSuccess(c, "Organization updated successfully", updateData)
}

func (r *ResourceOrganization) SoftDeleteOrganization(c *gin.Context) {
	idStr := c.Param("id")

	// Validar formato UUID
	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid organization ID format", err)
		return
	}

	err = r.handler.HandlerOrganization.SoftDeleteOrganization(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error deleting organization", err)
		return
	}

	utils.SendOKSuccess(c, "Organization deleted successfully", nil)
}

func (r *ResourceOrganization) HardDeleteOrganization(c *gin.Context) {
	idStr := c.Param("id")

	// Validar formato UUID
	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid organization ID format", err)
		return
	}

	err = r.handler.HandlerOrganization.HardDeleteOrganization(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error permanently deleting organization", err)
		return
	}

	utils.SendOKSuccess(c, "Organization permanently deleted", nil)
}

type CreateOrganizationBootstrapRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (r *ResourceOrganization) ServiceCreateOrganizationBootstrap(c *gin.Context) {
	var request CreateOrganizationBootstrapRequest
	err := c.BindJSON(&request)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Validar campos obrigatórios
	if strings.TrimSpace(request.Name) == "" {
		utils.SendBadRequestError(c, "Nome da organização é obrigatório", nil)
		return
	}

	if strings.TrimSpace(request.Password) == "" {
		utils.SendBadRequestError(c, "Senha é obrigatória", nil)
		return
	}

	// Chamar handler para bootstrap
	response, err := r.handler.HandlerOrganization.CreateOrganizationBootstrap(request.Name, request.Password)
	if err != nil {
		if strings.Contains(err.Error(), "senha inválida") {
			utils.SendValidationError(c, "Senha inválida", err)
			return
		}
		if strings.Contains(err.Error(), "já existe uma organização") {
			utils.SendValidationError(c, "Já existe uma organização com esse nome", err)
			return
		}
		utils.SendInternalServerError(c, "Erro ao criar organização", err)
		return
	}

	utils.SendCreatedSuccess(c, "Organização criada com sucesso", response)
}

func NewSourceServerOrganization(handler *handler.Handlers) IServerOrganization {
	return &ResourceOrganization{handler: handler}
}