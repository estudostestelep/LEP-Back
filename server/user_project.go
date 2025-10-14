package server

import (
	"lep/handler"
	"lep/repositories/models"
	"lep/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ResourceUserProject struct {
	handler handler.IHandlerUserProject
}

type IServerUserProject interface {
	ServiceAddUserToProject(c *gin.Context)
	ServiceRemoveUserFromProject(c *gin.Context)
	ServiceUpdateUserProject(c *gin.Context)
	ServiceGetUserProjects(c *gin.Context)
	ServiceGetUserProjectsByOrganization(c *gin.Context)
	ServiceGetProjectUsers(c *gin.Context)
}

// POST /user/:userId/project
func (r *ResourceUserProject) ServiceAddUserToProject(c *gin.Context) {
	userId := c.Param("userId")

	// Validar formato UUID
	if _, err := uuid.Parse(userId); err != nil {
		utils.SendBadRequestError(c, "Invalid user ID format", err)
		return
	}

	var userProj models.UserProject
	if err := c.BindJSON(&userProj); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Definir userId do path
	userProj.UserId, _ = uuid.Parse(userId)

	// Gerar ID se necessário
	if userProj.Id == uuid.Nil {
		userProj.Id = uuid.New()
	}

	// Role padrão se não fornecido
	if userProj.Role == "" {
		userProj.Role = "member"
	}

	err := r.handler.AddUserToProject(&userProj)
	if err != nil {
		utils.SendInternalServerError(c, "Error adding user to project", err)
		return
	}

	utils.SendCreatedSuccess(c, "User added to project successfully", userProj)
}

// DELETE /user/:userId/project/:projectId
func (r *ResourceUserProject) ServiceRemoveUserFromProject(c *gin.Context) {
	userId := c.Param("userId")
	projectId := c.Param("projectId")

	// Validar formatos UUID
	if _, err := uuid.Parse(userId); err != nil {
		utils.SendBadRequestError(c, "Invalid user ID format", err)
		return
	}
	if _, err := uuid.Parse(projectId); err != nil {
		utils.SendBadRequestError(c, "Invalid project ID format", err)
		return
	}

	err := r.handler.RemoveUserFromProject(userId, projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error removing user from project", err)
		return
	}

	utils.SendOKSuccess(c, "User removed from project successfully", nil)
}

// PUT /user-project/:id
func (r *ResourceUserProject) ServiceUpdateUserProject(c *gin.Context) {
	id := c.Param("id")

	// Validar formato UUID
	if _, err := uuid.Parse(id); err != nil {
		utils.SendBadRequestError(c, "Invalid ID format", err)
		return
	}

	var userProj models.UserProject
	if err := c.BindJSON(&userProj); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	userProj.Id, _ = uuid.Parse(id)

	err := r.handler.UpdateUserProject(&userProj)
	if err != nil {
		utils.SendInternalServerError(c, "Error updating user project", err)
		return
	}

	utils.SendOKSuccess(c, "User project updated successfully", userProj)
}

// GET /user/:userId/projects
func (r *ResourceUserProject) ServiceGetUserProjects(c *gin.Context) {
	userId := c.Param("userId")

	// Validar formato UUID
	if _, err := uuid.Parse(userId); err != nil {
		utils.SendBadRequestError(c, "Invalid user ID format", err)
		return
	}

	userProjs, err := r.handler.GetUserProjects(userId)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting user projects", err)
		return
	}

	c.JSON(http.StatusOK, userProjs)
}

// GET /user/:userId/organization/:orgId/projects
func (r *ResourceUserProject) ServiceGetUserProjectsByOrganization(c *gin.Context) {
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

	userProjs, err := r.handler.GetUserProjectsByOrganization(userId, orgId)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting user projects by organization", err)
		return
	}

	c.JSON(http.StatusOK, userProjs)
}

// GET /project/:projectId/users
func (r *ResourceUserProject) ServiceGetProjectUsers(c *gin.Context) {
	projectId := c.Param("projectId")

	// Validar formato UUID
	if _, err := uuid.Parse(projectId); err != nil {
		utils.SendBadRequestError(c, "Invalid project ID format", err)
		return
	}

	users, err := r.handler.GetProjectUsers(projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting project users", err)
		return
	}

	c.JSON(http.StatusOK, users)
}

func NewSourceServerUserProject(handler *handler.Handlers) IServerUserProject {
	return &ResourceUserProject{handler: handler.HandlerUserProject}
}
