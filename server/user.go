package server

import (
	"fmt"
	"lep/constants"
	"lep/handler"
	"lep/repositories/models"
	"lep/resource/validation"
	"lep/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ResourceUsers struct {
	handler          handler.IHandlerUser
	adminAuditHandler handler.IAdminAuditLogHandler
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

	// Buscar usuário com suas organizações e projetos
	resp, err := r.handler.GetUserWithRelations(idStr)
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

	fmt.Printf("📦 ServiceCreateUser - Contexto recebido: orgId='%s', projectId='%s', userId='%s', roleId='%s'\n", organizationId, projectId, newUser.Id, roleId)

	err = r.handler.CreateUser(newUser, organizationId, projectId, roleId)
	if err != nil {
		utils.SendInternalServerError(c, "Error creating user", err)
		return
	}

	// 📝 LOG DE AUDITORIA: Registrar criação de usuário se o ator for Master Admin
	if r.isMasterAdmin(c) && r.adminAuditHandler != nil {
		actor := r.getActorFromContext(c)
		if actor != nil {
			orgId, projId, isAdminZone := r.getAuditContext(c)
			ipAddress := c.ClientIP()
			userAgent := c.Request.UserAgent()

			go func() {
				if err := r.adminAuditHandler.LogUserCreate(actor, newUser, orgId, projId, isAdminZone, ipAddress, userAgent); err != nil {
					fmt.Printf("⚠️ Erro ao registrar log de auditoria (CREATE): %v\n", err)
				}
			}()
		}
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

	// 📝 AUDITORIA: Capturar estado anterior do usuário ANTES do update
	var oldUser *models.User
	if r.isMasterAdmin(c) && r.adminAuditHandler != nil {
		oldUser, _ = r.handler.GetUser(idStr)
	}

	var updatedUser models.User
	err = c.BindJSON(&updatedUser)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Verificar se é reset de senha (password não vazio no request)
	isPasswordReset := updatedUser.Password != ""

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

	// 📝 LOG DE AUDITORIA: Registrar atualização de usuário se o ator for Master Admin
	if r.isMasterAdmin(c) && r.adminAuditHandler != nil && oldUser != nil {
		actor := r.getActorFromContext(c)
		if actor != nil {
			orgId, projId, isAdminZone := r.getAuditContext(c)
			ipAddress := c.ClientIP()
			userAgent := c.Request.UserAgent()

			go func() {
				// Se foi reset de senha, logar separadamente
				if isPasswordReset {
					if err := r.adminAuditHandler.LogPasswordReset(actor, &updatedUser, orgId, projId, isAdminZone, ipAddress, userAgent); err != nil {
						fmt.Printf("⚠️ Erro ao registrar log de auditoria (RESET_PASSWORD): %v\n", err)
					}
				}

				// Logar alterações de outros campos
				if err := r.adminAuditHandler.LogUserUpdate(actor, oldUser, &updatedUser, orgId, projId, isAdminZone, ipAddress, userAgent); err != nil {
					fmt.Printf("⚠️ Erro ao registrar log de auditoria (UPDATE): %v\n", err)
				}
			}()
		}
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

	// 📝 AUDITORIA: Capturar dados do usuário ANTES da exclusão
	var targetUser *models.User
	if r.isMasterAdmin(c) && r.adminAuditHandler != nil {
		targetUser, _ = r.handler.GetUser(idStr)
	}

	err = r.handler.DeleteUser(idStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error deleting user", err)
		return
	}

	// 📝 LOG DE AUDITORIA: Registrar exclusão de usuário se o ator for Master Admin
	if r.isMasterAdmin(c) && r.adminAuditHandler != nil && targetUser != nil {
		actor := r.getActorFromContext(c)
		if actor != nil {
			orgId, projId, isAdminZone := r.getAuditContext(c)
			ipAddress := c.ClientIP()
			userAgent := c.Request.UserAgent()

			go func() {
				if err := r.adminAuditHandler.LogUserDelete(actor, targetUser, orgId, projId, isAdminZone, ipAddress, userAgent); err != nil {
					fmt.Printf("⚠️ Erro ao registrar log de auditoria (DELETE): %v\n", err)
				}
			}()
		}
	}

	utils.SendOKSuccess(c, "User deleted successfully", nil)
}

func NewSourceServerUsers(handler *handler.Handlers) IServerUsers {
	return &ResourceUsers{
		handler:          handler.HandlerUser,
		adminAuditHandler: handler.HandlerAdminAuditLog,
	}
}

// isMasterAdmin verifica se o usuário atual é um Master Admin
func (r *ResourceUsers) isMasterAdmin(c *gin.Context) bool {
	permissions, exists := c.Get("user_permissions")
	if !exists {
		return false
	}
	permList, ok := permissions.(pq.StringArray)
	if !ok {
		return false
	}
	return constants.HasPermission(permList, constants.PermissionMasterAdmin)
}

// getActorFromContext obtém o usuário ator (quem está fazendo a ação) do contexto
func (r *ResourceUsers) getActorFromContext(c *gin.Context) *models.User {
	userIdStr := c.GetString("user_id")
	if userIdStr == "" {
		return nil
	}
	actor, err := r.handler.GetUser(userIdStr)
	if err != nil {
		return nil
	}
	return actor
}

// getAuditContext determina o contexto da auditoria (zona admin ou org/projeto)
func (r *ResourceUsers) getAuditContext(c *gin.Context) (orgId, projectId *uuid.UUID, isAdminZone bool) {
	path := c.Request.URL.Path
	isAdminZone = strings.HasPrefix(path, "/admin")

	if !isAdminZone {
		orgIdStr := c.GetString("organization_id")
		projectIdStr := c.GetString("project_id")

		if orgIdStr != "" {
			if parsed, err := uuid.Parse(orgIdStr); err == nil {
				orgId = &parsed
			}
		}
		if projectIdStr != "" {
			if parsed, err := uuid.Parse(projectIdStr); err == nil {
				projectId = &parsed
			}
		}
	}

	return orgId, projectId, isAdminZone
}
