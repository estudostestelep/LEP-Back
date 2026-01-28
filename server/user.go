package server

import (
	"fmt"
	"lep/handler"
	"lep/repositories/models"
	"lep/resource/validation"
	"lep/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ResourceUsers struct {
	handler     handler.IHandlerUser
	authHandler handler.IHandlerAuth
}

type IServerUsers interface {
	ServiceGetUser(c *gin.Context)
	ServiceGetUserByGroup(c *gin.Context)
	ServiceListUsers(c *gin.Context)
	ServiceCreateUser(c *gin.Context)
	ServiceUpdateUser(c *gin.Context)
	ServiceDeleteUser(c *gin.Context)
	ServiceGetUserAccess(c *gin.Context)
}

func (r *ResourceUsers) ServiceGetUser(c *gin.Context) {
	id, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "user")
	if !ok {
		return
	}

	// Buscar usuário com suas organizações e projetos
	resp, err := r.handler.GetUserWithRelations(id.String())
	if err != nil {
		utils.SendInternalServerError(c, "Error getting user", err)
		return
	}

	if resp == nil {
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
	// Usar DTO que aceita tanto 'role' quanto 'permissions'
	var request models.CreateUserRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Converter request para User (mapeia role -> permissions se necessário)
	newUser := request.ToUser()

	// Gerar ID se não fornecido
	if newUser.Id == uuid.Nil {
		newUser.Id = uuid.New()
	}

	// Validações estruturadas
	if err := validation.CreateUserValidation(newUser); err != nil {
		utils.SendValidationError(c, "Validation failed", err)
		return
	}

	// Pegar organization_id e project_id dos headers para vincular o usuario
	organizationId := c.Request.Header.Get("X-Lpe-Organization-Id")
	projectId := c.Request.Header.Get("X-Lpe-Project-Id")

	// Pegar roleId do request para atribuir cargo ao usuário
	roleId := request.RoleId

	// Validação: cargo é obrigatório para criar usuário
	if roleId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "É necessário selecionar um cargo para o usuário"})
		return
	}

	fmt.Printf("📦 ServiceCreateUser - Contexto recebido: orgId='%s', projectId='%s', userId='%s', roleId='%s'\n", organizationId, projectId, newUser.Id, roleId)

	// Construir contexto da requisição para auditoria
	ctx := handler.BuildRequestContext(c)

	// Usar método com contexto que já lida com auditoria internamente
	err = r.handler.CreateUserWithContext(ctx, newUser, organizationId, projectId, roleId)
	if err != nil {
		utils.SendInternalServerError(c, "Error creating user", err)
		return
	}

	utils.SendCreatedSuccess(c, "User created successfully", newUser)
}

func (r *ResourceUsers) ServiceUpdateUser(c *gin.Context) {
	id, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "user")
	if !ok {
		return
	}

	var updatedUser models.User
	if err := c.BindJSON(&updatedUser); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	updatedUser.Id = id

	// Validações estruturadas
	if err := validation.UpdateUserValidation(&updatedUser); err != nil {
		utils.SendValidationError(c, "Validation failed", err)
		return
	}

	// Construir contexto da requisição para auditoria
	ctx := handler.BuildRequestContext(c)

	// Handler captura estado anterior internamente para auditoria
	if err := r.handler.UpdateUserWithContext(ctx, id.String(), &updatedUser); err != nil {
		utils.SendInternalServerError(c, "Error updating user", err)
		return
	}

	utils.SendOKSuccess(c, "User updated successfully", updatedUser)
}

func (r *ResourceUsers) ServiceListUsers(c *gin.Context) {
	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	resp, err := r.handler.ListUsersWithRoles(organizationId, projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error listing users"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceUsers) ServiceDeleteUser(c *gin.Context) {
	id, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "user")
	if !ok {
		return
	}

	// Construir contexto da requisição para auditoria
	ctx := handler.BuildRequestContext(c)

	// Handler captura estado anterior internamente e verifica existência
	if err := r.handler.DeleteUserWithContext(ctx, id.String()); err != nil {
		// Verificar se é erro de não encontrado
		if strings.Contains(err.Error(), "não encontrado") {
			utils.SendNotFoundError(c, "User")
			return
		}
		utils.SendInternalServerError(c, "Error deleting user", err)
		return
	}

	utils.SendOKSuccess(c, "User deleted successfully", nil)
}

func (r *ResourceUsers) ServiceGetUserAccess(c *gin.Context) {
	id, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "user")
	if !ok {
		return
	}

	orgs, err := r.authHandler.GetUserOrganizationsWithNames(id.String())
	if err != nil {
		utils.SendInternalServerError(c, "Error getting user organizations", err)
		return
	}

	projects, err := r.authHandler.GetUserProjectsWithNames(id.String())
	if err != nil {
		utils.SendInternalServerError(c, "Error getting user projects", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"organizations": orgs,
			"projects":      projects,
		},
	})
}

func NewSourceServerUsers(handler *handler.Handlers) IServerUsers {
	return &ResourceUsers{
		handler:     handler.HandlerUser,
		authHandler: handler.HandlerAuth,
	}
}
