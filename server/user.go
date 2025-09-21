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

type ResourceUsers struct {
	handler handler.IHandlerUser
}

type IServerUsers interface {
	ServiceGetUser(c *gin.Context)
	ServiceGetUserByGroup(c *gin.Context)
	ServiceListUsers(c *gin.Context)
	ServiceCreateUser(c *gin.Context)
	ServiceUpdateUser(c *gin.Context)
	ServiceDeleteUser(c *gin.Context)
}

func (r *ResourceUsers) ServiceGetUser(c *gin.Context) {
	idStr := c.Param("id")

	// Validar formato UUID
	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid user ID format", err)
		return
	}

	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	resp, err := r.handler.GetUser(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting user", err)
		return
	}

	if resp == nil {
		utils.SendNotFoundError(c, "User")
		return
	}

	// Verificar se o usuário pertence à organização/projeto
	if resp.OrganizationId.String() != organizationId || resp.ProjectId.String() != projectId {
		utils.SendNotFoundError(c, "User")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceUsers) ServiceGetUserByGroup(c *gin.Context) {
	id := c.Param("id")

	resp, err := r.handler.GetUserByGroup(id)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting users by group", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceUsers) ServiceCreateUser(c *gin.Context) {
	var newUser models.User
	err := c.BindJSON(&newUser)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Para criação de usuário, pode ser público (sem headers) ou protegido
	// Se headers existem, usar; senão, validar que são fornecidos no body
	orgHeader := c.GetHeader("X-Lpe-Organization-Id")
	projHeader := c.GetHeader("X-Lpe-Project-Id")

	if orgHeader != "" && projHeader != "" {
		// Rota protegida - usar headers
		newUser.OrganizationId, err = uuid.Parse(orgHeader)
		if err != nil {
			utils.SendBadRequestError(c, "Invalid organization ID format", err)
			return
		}
		newUser.ProjectId, err = uuid.Parse(projHeader)
		if err != nil {
			utils.SendBadRequestError(c, "Invalid project ID format", err)
			return
		}
	}

	// Gerar ID se não fornecido
	if newUser.Id == uuid.Nil {
		newUser.Id = uuid.New()
	}

	// Validações estruturadas
	if err := validation.CreateUserValidation(&newUser); err != nil {
		utils.SendValidationError(c, "Validation failed", err)
		return
	}

	err = r.handler.CreateUser(&newUser)
	if err != nil {
		utils.SendInternalServerError(c, "Error creating user", err)
		return
	}

	utils.SendCreatedSuccess(c, "User created successfully", newUser)
}

func (r *ResourceUsers) ServiceUpdateUser(c *gin.Context) {
	idStr := c.Param("id")

	// Validar formato UUID
	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid user ID format", err)
		return
	}

	var updatedUser models.User
	err = c.BindJSON(&updatedUser)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	updatedUser.OrganizationId, err = uuid.Parse(organizationId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing organization ID", err)
		return
	}
	updatedUser.ProjectId, err = uuid.Parse(projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing project ID", err)
		return
	}
	updatedUser.Id, err = uuid.Parse(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing user ID", err)
		return
	}

	// Validações estruturadas
	if err := validation.UpdateUserValidation(&updatedUser); err != nil {
		utils.SendValidationError(c, "Validation failed", err)
		return
	}

	err = r.handler.UpdateUser(&updatedUser)
	if err != nil {
		utils.SendInternalServerError(c, "Error updating user", err)
		return
	}

	utils.SendOKSuccess(c, "User updated successfully", updatedUser)
}

func (r *ResourceUsers) ServiceListUsers(c *gin.Context) {
	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	resp, err := r.handler.ListUsers(organizationId, projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error listing users"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceUsers) ServiceDeleteUser(c *gin.Context) {
	idStr := c.Param("id")

	// Validar formato UUID
	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid user ID format", err)
		return
	}

	err = r.handler.DeleteUser(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error deleting user", err)
		return
	}

	utils.SendOKSuccess(c, "User deleted successfully", nil)
}

func NewSourceServerUsers(handler *handler.Handlers) IServerUsers {
	return &ResourceUsers{handler: handler.HandlerUser}
}
