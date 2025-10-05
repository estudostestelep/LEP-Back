package server

import (
	"lep/handler"
	"lep/repositories/models"
	"lep/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ResourceUserOrganization struct {
	handler handler.IHandlerUserOrganization
}

type IServerUserOrganization interface {
	ServiceAddUserToOrganization(c *gin.Context)
	ServiceRemoveUserFromOrganization(c *gin.Context)
	ServiceUpdateUserOrganization(c *gin.Context)
	ServiceGetUserOrganizations(c *gin.Context)
	ServiceGetOrganizationUsers(c *gin.Context)
}

// POST /user/:userId/organization
func (r *ResourceUserOrganization) ServiceAddUserToOrganization(c *gin.Context) {
	userId := c.Param("userId")

	// Validar formato UUID
	if _, err := uuid.Parse(userId); err != nil {
		utils.SendBadRequestError(c, "Invalid user ID format", err)
		return
	}

	var userOrg models.UserOrganization
	if err := c.BindJSON(&userOrg); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Definir userId do path
	userOrg.UserId, _ = uuid.Parse(userId)

	// Gerar ID se necessário
	if userOrg.Id == uuid.Nil {
		userOrg.Id = uuid.New()
	}

	// Role padrão se não fornecido
	if userOrg.Role == "" {
		userOrg.Role = "member"
	}

	err := r.handler.AddUserToOrganization(&userOrg)
	if err != nil {
		utils.SendInternalServerError(c, "Error adding user to organization", err)
		return
	}

	utils.SendCreatedSuccess(c, "User added to organization successfully", userOrg)
}

// DELETE /user/:userId/organization/:orgId
func (r *ResourceUserOrganization) ServiceRemoveUserFromOrganization(c *gin.Context) {
	userId := c.Param("userId")
	orgId := c.Param("orgId")

	// Validar formatos UUID
	if _, err := uuid.Parse(userId); err != nil {
		utils.SendBadRequestError(c, "Invalid user ID format", err)
		return
	}
	if _, err := uuid.Parse(orgId); err != nil {
		utils.SendBadRequestError(c, "Invalid organization ID format", err)
		return
	}

	err := r.handler.RemoveUserFromOrganization(userId, orgId)
	if err != nil {
		utils.SendInternalServerError(c, "Error removing user from organization", err)
		return
	}

	utils.SendOKSuccess(c, "User removed from organization successfully", nil)
}

// PUT /user-organization/:id
func (r *ResourceUserOrganization) ServiceUpdateUserOrganization(c *gin.Context) {
	id := c.Param("id")

	// Validar formato UUID
	if _, err := uuid.Parse(id); err != nil {
		utils.SendBadRequestError(c, "Invalid ID format", err)
		return
	}

	var userOrg models.UserOrganization
	if err := c.BindJSON(&userOrg); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	userOrg.Id, _ = uuid.Parse(id)

	err := r.handler.UpdateUserOrganization(&userOrg)
	if err != nil {
		utils.SendInternalServerError(c, "Error updating user organization", err)
		return
	}

	utils.SendOKSuccess(c, "User organization updated successfully", userOrg)
}

// GET /user/:userId/organizations
func (r *ResourceUserOrganization) ServiceGetUserOrganizations(c *gin.Context) {
	userId := c.Param("userId")

	// Validar formato UUID
	if _, err := uuid.Parse(userId); err != nil {
		utils.SendBadRequestError(c, "Invalid user ID format", err)
		return
	}

	userOrgs, err := r.handler.GetUserOrganizations(userId)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting user organizations", err)
		return
	}

	c.JSON(http.StatusOK, userOrgs)
}

// GET /organization/:orgId/users
func (r *ResourceUserOrganization) ServiceGetOrganizationUsers(c *gin.Context) {
	orgId := c.Param("orgId")

	// Validar formato UUID
	if _, err := uuid.Parse(orgId); err != nil {
		utils.SendBadRequestError(c, "Invalid organization ID format", err)
		return
	}

	users, err := r.handler.GetOrganizationUsers(orgId)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting organization users", err)
		return
	}

	c.JSON(http.StatusOK, users)
}

func NewSourceServerUserOrganization(handler *handler.Handlers) IServerUserOrganization {
	return &ResourceUserOrganization{handler: handler.HandlerUserOrganization}
}
