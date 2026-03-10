package server

import (
	"lep/handler"
	"lep/resource/validation"
	"lep/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResourceUserAccess struct {
	handler handler.IHandlerUserAccess
}

type IServerUserAccess interface {
	ServiceGetUserAccess(c *gin.Context)
	ServiceUpdateUserAccess(c *gin.Context)
}

// UpdateUserAccessRequest DTO para atualizar acesso do usuário
type UpdateUserAccessRequest struct {
	OrganizationIds []string `json:"organization_ids"`
	ProjectIds      []string `json:"project_ids"`
}

func (r *ResourceUserAccess) ServiceGetUserAccess(c *gin.Context) {
	userId, ok := validation.ParseAndValidateUUID(c, c.Param("userId"), "user")
	if !ok {
		return
	}

	access, err := r.handler.GetUserAccess(userId.String())
	if err != nil {
		if err.Error() == "usuário não encontrado" {
			utils.SendNotFoundError(c, "User")
			return
		}
		utils.SendInternalServerError(c, "Error getting user access", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"organizations": access.Organizations,
		"projects":      access.Projects,
	})
}

func (r *ResourceUserAccess) ServiceUpdateUserAccess(c *gin.Context) {
	userId, ok := validation.ParseAndValidateUUID(c, c.Param("userId"), "user")
	if !ok {
		return
	}

	var request UpdateUserAccessRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Converter para o formato do handler
	handlerRequest := &handler.UpdateUserAccessRequest{
		OrganizationIds: request.OrganizationIds,
		ProjectIds:      request.ProjectIds,
	}

	result, err := r.handler.UpdateUserAccess(userId.String(), handlerRequest)
	if err != nil {
		if err.Error() == "usuário não encontrado" {
			utils.SendNotFoundError(c, "User")
			return
		}
		utils.SendInternalServerError(c, "Error updating user access", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":               result.Message,
		"organizations_added":   result.OrganizationsAdded,
		"organizations_removed": result.OrganizationsRemoved,
		"projects_added":        result.ProjectsAdded,
		"projects_removed":      result.ProjectsRemoved,
	})
}

func NewSourceServerUserAccess(handler *handler.Handlers) IServerUserAccess {
	return &ResourceUserAccess{
		handler: handler.HandlerUserAccess,
	}
}
