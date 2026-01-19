package server

import (
	"fmt"
	"lep/handler"
	"lep/repositories/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoleServer struct {
	handler           *handler.RoleHandler
	limitHandler      *handler.LimitHandler
	adminAuditHandler handler.IAdminAuditLogHandler
}

func NewRoleServer(h *handler.RoleHandler) *RoleServer {
	return &RoleServer{handler: h}
}

// SetAdminAuditHandler configura o handler de auditoria (injetado separadamente)
func (s *RoleServer) SetAdminAuditHandler(h handler.IAdminAuditLogHandler) {
	s.adminAuditHandler = h
}

// SetLimitHandler configura o handler de limites (injetado separadamente)
func (s *RoleServer) SetLimitHandler(h *handler.LimitHandler) {
	s.limitHandler = h
}

// ==================== Role CRUD ====================

// CreateRole godoc
// @Summary Cria um novo cargo
// @Tags Roles
// @Accept json
// @Produce json
// @Param role body models.Role true "Dados do cargo"
// @Success 201 {object} models.Role
// @Router /role [post]
func (s *RoleServer) CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos", "error": err.Error()})
		return
	}

	userId := c.GetString("user_id")
	userEmail := c.GetString("user_email")
	orgId := c.GetString("organization_id")

	// Se não especificado, associar à organização atual
	var orgUUID *uuid.UUID
	if role.OrganizationId == nil {
		parsed, err := uuid.Parse(orgId)
		if err == nil {
			orgUUID = &parsed
			role.OrganizationId = orgUUID
		}
	} else {
		orgUUID = role.OrganizationId
	}

	if err := s.handler.CreateRole(&role, userId, orgId); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	// 📝 LOG DE AUDITORIA: Registrar criação de cargo
	if s.adminAuditHandler != nil {
		actorUUID, _ := uuid.Parse(userId)
		ipAddress := c.ClientIP()
		userAgent := c.Request.UserAgent()

		go func() {
			s.adminAuditHandler.LogGenericAction(handler.AuditLogParams{
				ActorId:       actorUUID,
				ActorEmail:    userEmail,
				TargetId:      role.Id,
				TargetEmail:   "",
				Action:        models.AdminAuditActionCreate,
				EntityType:    models.AdminAuditEntityRole,
				OrgId:         orgUUID,
				ProjectId:     role.ProjectId,
				IsAdminZone:   true,
				OldValues:     nil,
				NewValues:     map[string]interface{}{"name": role.Name, "display_name": role.DisplayName, "scope": role.Scope},
				ChangedFields: []string{"*"},
				IpAddress:     ipAddress,
				UserAgent:     userAgent,
			})
		}()
	}

	c.JSON(http.StatusCreated, gin.H{"data": role})
}

// GetRole godoc
// @Summary Busca um cargo por ID
// @Tags Roles
// @Produce json
// @Param id path string true "ID do cargo"
// @Success 200 {object} models.Role
// @Router /role/{id} [get]
func (s *RoleServer) GetRole(c *gin.Context) {
	id := c.Param("id")

	role, err := s.handler.GetRole(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Cargo não encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": role})
}

// UpdateRole godoc
// @Summary Atualiza um cargo
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path string true "ID do cargo"
// @Param role body models.Role true "Dados atualizados"
// @Success 200 {object} models.Role
// @Router /role/{id} [put]
func (s *RoleServer) UpdateRole(c *gin.Context) {
	id := c.Param("id")

	// Capturar estado anterior para auditoria
	oldRole, _ := s.handler.GetRole(id)

	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos", "error": err.Error()})
		return
	}

	roleUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID inválido"})
		return
	}
	role.Id = roleUUID

	userId := c.GetString("user_id")
	userEmail := c.GetString("user_email")
	orgId := c.GetString("organization_id")

	if err := s.handler.UpdateRole(&role, userId, orgId); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	// 📝 LOG DE AUDITORIA: Registrar atualização de cargo
	if s.adminAuditHandler != nil && oldRole != nil {
		actorUUID, _ := uuid.Parse(userId)
		var orgUUID *uuid.UUID
		if orgId != "" {
			parsed, _ := uuid.Parse(orgId)
			orgUUID = &parsed
		}
		ipAddress := c.ClientIP()
		userAgent := c.Request.UserAgent()

		go func() {
			s.adminAuditHandler.LogGenericAction(handler.AuditLogParams{
				ActorId:       actorUUID,
				ActorEmail:    userEmail,
				TargetId:      roleUUID,
				TargetEmail:   "",
				Action:        models.AdminAuditActionUpdate,
				EntityType:    models.AdminAuditEntityRole,
				OrgId:         orgUUID,
				ProjectId:     role.ProjectId,
				IsAdminZone:   true,
				OldValues:     map[string]interface{}{"name": oldRole.Name, "display_name": oldRole.DisplayName},
				NewValues:     map[string]interface{}{"name": role.Name, "display_name": role.DisplayName},
				ChangedFields: []string{"name", "display_name", "permissions"},
				IpAddress:     ipAddress,
				UserAgent:     userAgent,
			})
		}()
	}

	c.JSON(http.StatusOK, gin.H{"data": role})
}

// DeleteRole godoc
// @Summary Remove um cargo
// @Tags Roles
// @Produce json
// @Param id path string true "ID do cargo"
// @Success 200 {object} map[string]string
// @Router /role/{id} [delete]
func (s *RoleServer) DeleteRole(c *gin.Context) {
	id := c.Param("id")
	userId := c.GetString("user_id")
	userEmail := c.GetString("user_email")
	orgId := c.GetString("organization_id")

	// Capturar dados do cargo ANTES de deletar
	oldRole, _ := s.handler.GetRole(id)

	if err := s.handler.DeleteRole(id, userId, orgId); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	// 📝 LOG DE AUDITORIA: Registrar exclusão de cargo
	if s.adminAuditHandler != nil && oldRole != nil {
		actorUUID, _ := uuid.Parse(userId)
		roleUUID, _ := uuid.Parse(id)
		var orgUUID *uuid.UUID
		if orgId != "" {
			parsed, _ := uuid.Parse(orgId)
			orgUUID = &parsed
		}
		ipAddress := c.ClientIP()
		userAgent := c.Request.UserAgent()

		go func() {
			s.adminAuditHandler.LogGenericAction(handler.AuditLogParams{
				ActorId:       actorUUID,
				ActorEmail:    userEmail,
				TargetId:      roleUUID,
				TargetEmail:   "",
				Action:        models.AdminAuditActionDelete,
				EntityType:    models.AdminAuditEntityRole,
				OrgId:         orgUUID,
				ProjectId:     oldRole.ProjectId,
				IsAdminZone:   true,
				OldValues:     map[string]interface{}{"name": oldRole.Name, "display_name": oldRole.DisplayName},
				NewValues:     nil,
				ChangedFields: []string{"*"},
				IpAddress:     ipAddress,
				UserAgent:     userAgent,
			})
		}()
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cargo removido com sucesso"})
}

// ListRoles godoc
// @Summary Lista cargos
// @Tags Roles
// @Produce json
// @Param scope query string false "Filtrar por escopo (admin/client)"
// @Success 200 {array} models.Role
// @Router /role [get]
func (s *RoleServer) ListRoles(c *gin.Context) {
	scope := c.Query("scope")
	orgId := c.GetString("organization_id")

	roles, err := s.handler.ListRoles(scope, orgId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao listar cargos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": roles})
}

// ListSystemRoles godoc
// @Summary Lista cargos do sistema
// @Tags Roles
// @Produce json
// @Success 200 {array} models.Role
// @Router /role/system [get]
func (s *RoleServer) ListSystemRoles(c *gin.Context) {
	roles, err := s.handler.ListSystemRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao listar cargos do sistema"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": roles})
}

// ==================== User-Role Assignment ====================

// AssignRoleToUser godoc
// @Summary Atribui um cargo a um usuário
// @Tags Roles
// @Accept json
// @Produce json
// @Param assignment body AssignRoleRequest true "Dados da atribuição"
// @Success 200 {object} map[string]string
// @Router /role/assign [post]
func (s *RoleServer) AssignRoleToUser(c *gin.Context) {
	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos", "error": err.Error()})
		return
	}

	actorUserId := c.GetString("user_id")
	actorEmail := c.GetString("user_email")
	orgId := c.GetString("organization_id")

	userUUID, err := uuid.Parse(req.UserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID do usuário inválido"})
		return
	}

	roleUUID, err := uuid.Parse(req.RoleId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID do cargo inválido"})
		return
	}

	userRole := &models.UserRole{
		UserId: userUUID,
		RoleId: roleUUID,
		Active: true,
	}

	// Se orgId foi fornecido, adicionar ao contexto
	// Se vazio, OrganizationId fica nil (cargo admin global)
	var orgUUID *uuid.UUID
	if orgId != "" {
		parsed, err := uuid.Parse(orgId)
		if err == nil {
			orgUUID = &parsed
			userRole.OrganizationId = orgUUID
		}
	}

	// Se tiver projectId, adicionar
	var projUUID *uuid.UUID
	if req.ProjectId != "" {
		parsed, err := uuid.Parse(req.ProjectId)
		if err == nil {
			projUUID = &parsed
			userRole.ProjectId = projUUID
		}
	}

	if err := s.handler.AssignRoleToUser(userRole, actorUserId); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	// 📝 LOG DE AUDITORIA: Registrar atribuição de cargo
	if s.adminAuditHandler != nil {
		// Buscar informações do cargo
		role, _ := s.handler.GetRole(req.RoleId)

		roleName := ""
		if role != nil {
			roleName = role.DisplayName
		}

		actorUUID, _ := uuid.Parse(actorUserId)
		ipAddress := c.ClientIP()
		userAgent := c.Request.UserAgent()

		go func() {
			if err := s.adminAuditHandler.LogRoleAssignment(
				actorUUID, actorEmail,
				userUUID, "", // targetEmail será vazio, pois não temos acesso fácil
				roleUUID, roleName,
				orgUUID, projUUID,
				ipAddress, userAgent,
			); err != nil {
				fmt.Printf("⚠️ Erro ao registrar log de auditoria (ASSIGN_ROLE): %v\n", err)
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cargo atribuído com sucesso"})
}

// RemoveRoleFromUser godoc
// @Summary Remove um cargo de um usuário
// @Tags Roles
// @Accept json
// @Produce json
// @Param assignment body RemoveRoleRequest true "Dados da remoção"
// @Success 200 {object} map[string]string
// @Router /role/remove [post]
func (s *RoleServer) RemoveRoleFromUser(c *gin.Context) {
	var req RemoveRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos", "error": err.Error()})
		return
	}

	actorUserId := c.GetString("user_id")
	actorEmail := c.GetString("user_email")
	orgId := c.GetString("organization_id")

	// Buscar informações do cargo ANTES de remover (para o log)
	var roleName string
	role, _ := s.handler.GetRole(req.RoleId)
	if role != nil {
		roleName = role.DisplayName
	}

	if err := s.handler.RemoveRoleFromUser(req.UserId, req.RoleId, orgId, actorUserId); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	// 📝 LOG DE AUDITORIA: Registrar remoção de cargo
	if s.adminAuditHandler != nil {
		userUUID, _ := uuid.Parse(req.UserId)
		roleUUID, _ := uuid.Parse(req.RoleId)
		actorUUID, _ := uuid.Parse(actorUserId)

		var orgUUID *uuid.UUID
		if orgId != "" {
			parsed, _ := uuid.Parse(orgId)
			orgUUID = &parsed
		}

		ipAddress := c.ClientIP()
		userAgent := c.Request.UserAgent()

		go func() {
			if err := s.adminAuditHandler.LogRoleRemoval(
				actorUUID, actorEmail,
				userUUID, "",
				roleUUID, roleName,
				orgUUID, nil,
				ipAddress, userAgent,
			); err != nil {
				fmt.Printf("⚠️ Erro ao registrar log de auditoria (REMOVE_ROLE): %v\n", err)
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cargo removido com sucesso"})
}

// GetUserRoles godoc
// @Summary Lista cargos de um usuário
// @Tags Roles
// @Produce json
// @Param userId path string true "ID do usuário"
// @Success 200 {array} models.UserRole
// @Router /role/user/{userId} [get]
func (s *RoleServer) GetUserRoles(c *gin.Context) {
	userId := c.Param("userId")
	orgId := c.GetString("organization_id")

	roles, err := s.handler.GetUserRoles(userId, orgId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao buscar cargos do usuário"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": roles})
}

// GetUserRolesWithDetails godoc
// @Summary Lista cargos de um usuário com detalhes de permissões
// @Tags Roles
// @Produce json
// @Param userId path string true "ID do usuário"
// @Success 200 {array} models.RoleWithPermissionLevels
// @Router /role/user/{userId}/details [get]
func (s *RoleServer) GetUserRolesWithDetails(c *gin.Context) {
	userId := c.Param("userId")
	orgId := c.GetString("organization_id")

	roles, err := s.handler.GetUserRolesWithDetails(userId, orgId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao buscar detalhes dos cargos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": roles})
}

// ==================== Permission Levels ====================

// SetPermissionLevel godoc
// @Summary Define o nível de uma permissão para um cargo
// @Tags Roles
// @Accept json
// @Produce json
// @Param data body SetPermissionLevelRequest true "Dados do nível"
// @Success 200 {object} map[string]string
// @Router /role/permission-level [post]
func (s *RoleServer) SetPermissionLevel(c *gin.Context) {
	var req SetPermissionLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos", "error": err.Error()})
		return
	}

	userId := c.GetString("user_id")
	orgId := c.GetString("organization_id")

	if err := s.handler.SetRolePermissionLevel(req.RoleId, req.PermissionId, req.Level, userId, orgId); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nível de permissão definido com sucesso"})
}

// GetRolePermissions godoc
// @Summary Lista níveis de permissão de um cargo
// @Tags Roles
// @Produce json
// @Param id path string true "ID do cargo"
// @Success 200 {array} models.RolePermissionLevel
// @Router /role/{id}/permissions [get]
func (s *RoleServer) GetRolePermissions(c *gin.Context) {
	roleId := c.Param("id")

	levels, err := s.handler.GetRolePermissionLevels(roleId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao buscar permissões"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": levels})
}

// ==================== Permission Checking ====================

// CheckPermission godoc
// @Summary Verifica se o usuário tem uma permissão
// @Tags Roles
// @Produce json
// @Param permission query string true "Código da permissão"
// @Param level query int false "Nível mínimo (default: 1)"
// @Success 200 {object} PermissionCheckResponse
// @Router /role/check [get]
func (s *RoleServer) CheckPermission(c *gin.Context) {
	userId := c.GetString("user_id")
	orgId := c.GetString("organization_id")
	permission := c.Query("permission")
	level := 1 // Padrão: visualização

	if l := c.Query("level"); l != "" {
		if l == "2" {
			level = 2
		}
	}

	hasPermission, err := s.handler.HasPermission(userId, orgId, permission, level)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"has_permission": false,
			"error":          err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"has_permission": hasPermission,
		"permission":     permission,
		"level":          level,
	})
}

// GetMyPermissions godoc
// @Summary Retorna todas as permissões efetivas do usuário atual
// @Tags Roles
// @Produce json
// @Success 200 {object} MyPermissionsResponse
// @Router /role/my-permissions [get]
func (s *RoleServer) GetMyPermissions(c *gin.Context) {
	userId := c.GetString("user_id")
	orgId := c.GetString("organization_id")

	// Buscar cargos com detalhes
	roles, err := s.handler.GetUserRolesWithDetails(userId, orgId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao buscar permissões"})
		return
	}

	// Buscar nível de hierarquia
	hierarchyLevel, _ := s.handler.GetUserMaxHierarchyLevel(userId, orgId)

	// Agregar permissões efetivas
	effectivePermissions := make(map[string]int)
	for _, role := range roles {
		for _, pl := range role.PermissionLevels {
			if pl.Permission != nil {
				codeName := pl.Permission.CodeName
				if currentLevel, exists := effectivePermissions[codeName]; !exists || pl.Level > currentLevel {
					effectivePermissions[codeName] = pl.Level
				}
			}
		}
	}

	// Buscar módulos disponíveis
	modules, _ := s.handler.GetOrganizationModules(orgId)

	c.JSON(http.StatusOK, gin.H{
		"roles":                 roles,
		"effective_permissions": effectivePermissions,
		"hierarchy_level":       hierarchyLevel,
		"available_modules":     modules,
	})
}

// ==================== Module & Permission Listing ====================

// ListModules godoc
// @Summary Lista módulos
// @Tags Modules
// @Produce json
// @Param scope query string false "Filtrar por escopo (admin/client)"
// @Success 200 {array} models.Module
// @Router /module [get]
func (s *RoleServer) ListModules(c *gin.Context) {
	scope := c.Query("scope")

	modules, err := s.handler.ListModules(scope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao listar módulos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": modules})
}

// ListModulesWithPermissions godoc
// @Summary Lista módulos com suas permissões
// @Tags Modules
// @Produce json
// @Success 200 {array} models.Module
// @Router /module/with-permissions [get]
func (s *RoleServer) ListModulesWithPermissions(c *gin.Context) {
	modules, err := s.handler.ListModulesWithPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao listar módulos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": modules})
}

// ListPermissions godoc
// @Summary Lista permissões
// @Tags Permissions
// @Produce json
// @Param moduleId query string false "Filtrar por módulo"
// @Success 200 {array} models.Permission
// @Router /permission [get]
func (s *RoleServer) ListPermissions(c *gin.Context) {
	moduleId := c.Query("moduleId")

	permissions, err := s.handler.ListPermissions(moduleId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao listar permissões"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": permissions})
}

// GetOrganizationModules godoc
// @Summary Lista módulos disponíveis para a organização
// @Tags Modules
// @Produce json
// @Success 200 {array} models.Module
// @Router /module/available [get]
func (s *RoleServer) GetOrganizationModules(c *gin.Context) {
	orgId := c.GetString("organization_id")

	modules, err := s.handler.GetOrganizationModules(orgId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao buscar módulos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": modules})
}

// ==================== Package Management ====================

// ListPackages godoc
// @Summary Lista pacotes
// @Tags Packages
// @Produce json
// @Param public query bool false "Apenas pacotes públicos"
// @Success 200 {array} models.Package
// @Router /package [get]
func (s *RoleServer) ListPackages(c *gin.Context) {
	publicOnly := c.Query("public") == "true"

	packages, err := s.handler.ListPackages(publicOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao listar pacotes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": packages})
}

// GetPackageWithModules godoc
// @Summary Busca um pacote com seus módulos
// @Tags Packages
// @Produce json
// @Param id path string true "ID do pacote"
// @Success 200 {object} models.Package
// @Router /package/{id} [get]
func (s *RoleServer) GetPackageWithModules(c *gin.Context) {
	id := c.Param("id")

	pkg, err := s.handler.GetPackageWithModules(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Pacote não encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": pkg})
}

// GetOrganizationSubscription godoc
// @Summary Retorna a assinatura da organização
// @Tags Packages
// @Produce json
// @Success 200 {object} models.OrganizationPackage
// @Router /package/subscription [get]
func (s *RoleServer) GetOrganizationSubscription(c *gin.Context) {
	orgId := c.GetString("organization_id")

	subscription, err := s.handler.GetOrganizationSubscription(orgId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Nenhuma assinatura encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": subscription})
}

// SubscribeOrganization godoc
// @Summary Inscreve a organização em um pacote
// @Tags Packages
// @Accept json
// @Produce json
// @Param data body SubscribeRequest true "Dados da assinatura"
// @Success 200 {object} map[string]string
// @Router /package/subscribe [post]
func (s *RoleServer) SubscribeOrganization(c *gin.Context) {
	var req SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos", "error": err.Error()})
		return
	}

	orgId := c.GetString("organization_id")

	if err := s.handler.SubscribeOrganization(orgId, req.PackageId, req.BillingCycle, req.CustomPrice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assinatura realizada com sucesso"})
}

// ==================== Request/Response Types ====================

type AssignRoleRequest struct {
	UserId    string `json:"user_id" binding:"required"`
	RoleId    string `json:"role_id" binding:"required"`
	ProjectId string `json:"project_id,omitempty"`
}

type RemoveRoleRequest struct {
	UserId string `json:"user_id" binding:"required"`
	RoleId string `json:"role_id" binding:"required"`
}

type SetPermissionLevelRequest struct {
	RoleId       string `json:"role_id" binding:"required"`
	PermissionId string `json:"permission_id" binding:"required"`
	Level        int    `json:"level" binding:"required,min=0,max=2"`
}

type PermissionCheckResponse struct {
	HasPermission bool   `json:"has_permission"`
	Permission    string `json:"permission"`
	Level         int    `json:"level"`
	Error         string `json:"error,omitempty"`
}

type MyPermissionsResponse struct {
	Roles                []models.RoleWithPermissionLevels `json:"roles"`
	EffectivePermissions map[string]int                    `json:"effective_permissions"`
	HierarchyLevel       int                               `json:"hierarchy_level"`
	AvailableModules     []models.Module                   `json:"available_modules"`
}

type SubscribeRequest struct {
	PackageId    string   `json:"package_id" binding:"required"`
	BillingCycle string   `json:"billing_cycle" binding:"required,oneof=monthly yearly"`
	CustomPrice  *float64 `json:"custom_price,omitempty"`
}

type UpdateSubscriptionRequest struct {
	PackageId    string   `json:"package_id,omitempty"`
	BillingCycle string   `json:"billing_cycle,omitempty"`
	CustomPrice  *float64 `json:"custom_price,omitempty"`
	Active       *bool    `json:"active,omitempty"`
}

type SetPackageLimitRequest struct {
	LimitType  string `json:"limit_type" binding:"required"`
	LimitValue int    `json:"limit_value" binding:"required"`
}

// ==================== Package CRUD (Master Admin) ====================

// CreatePackage godoc
// @Summary Cria um novo pacote (Master Admin)
// @Tags Packages
// @Accept json
// @Produce json
// @Param package body models.Package true "Dados do pacote"
// @Success 201 {object} models.Package
// @Router /package [post]
func (s *RoleServer) CreatePackage(c *gin.Context) {
	var pkg models.Package
	if err := c.ShouldBindJSON(&pkg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos", "error": err.Error()})
		return
	}

	userId := c.GetString("user_id")
	userEmail := c.GetString("user_email")

	if pkg.Id == uuid.Nil {
		pkg.Id = uuid.New()
	}

	if err := s.handler.CreatePackage(&pkg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao criar pacote", "error": err.Error()})
		return
	}

	// 📝 LOG DE AUDITORIA
	if s.adminAuditHandler != nil {
		actorUUID, _ := uuid.Parse(userId)
		go func() {
			s.adminAuditHandler.LogGenericAction(handler.AuditLogParams{
				ActorId:       actorUUID,
				ActorEmail:    userEmail,
				TargetId:      pkg.Id,
				Action:        models.AdminAuditActionCreate,
				EntityType:    models.AdminAuditEntityPackage,
				IsAdminZone:   true,
				NewValues:     map[string]interface{}{"code_name": pkg.CodeName, "display_name": pkg.DisplayName},
				ChangedFields: []string{"*"},
				IpAddress:     c.ClientIP(),
				UserAgent:     c.Request.UserAgent(),
			})
		}()
	}

	c.JSON(http.StatusCreated, gin.H{"data": pkg})
}

// UpdatePackage godoc
// @Summary Atualiza um pacote (Master Admin)
// @Tags Packages
// @Accept json
// @Produce json
// @Param id path string true "ID do pacote"
// @Param package body models.Package true "Dados atualizados"
// @Success 200 {object} models.Package
// @Router /package/{id} [put]
func (s *RoleServer) UpdatePackage(c *gin.Context) {
	id := c.Param("id")
	userId := c.GetString("user_id")
	userEmail := c.GetString("user_email")

	var pkg models.Package
	if err := c.ShouldBindJSON(&pkg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos", "error": err.Error()})
		return
	}

	pkgUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID inválido"})
		return
	}
	pkg.Id = pkgUUID

	if err := s.handler.UpdatePackage(&pkg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao atualizar pacote", "error": err.Error()})
		return
	}

	// 📝 LOG DE AUDITORIA
	if s.adminAuditHandler != nil {
		actorUUID, _ := uuid.Parse(userId)
		go func() {
			s.adminAuditHandler.LogGenericAction(handler.AuditLogParams{
				ActorId:       actorUUID,
				ActorEmail:    userEmail,
				TargetId:      pkgUUID,
				Action:        models.AdminAuditActionUpdate,
				EntityType:    models.AdminAuditEntityPackage,
				IsAdminZone:   true,
				NewValues:     map[string]interface{}{"code_name": pkg.CodeName, "display_name": pkg.DisplayName},
				ChangedFields: []string{"code_name", "display_name", "prices"},
				IpAddress:     c.ClientIP(),
				UserAgent:     c.Request.UserAgent(),
			})
		}()
	}

	c.JSON(http.StatusOK, gin.H{"data": pkg})
}

// DeletePackage godoc
// @Summary Remove um pacote (Master Admin)
// @Tags Packages
// @Produce json
// @Param id path string true "ID do pacote"
// @Success 200 {object} map[string]string
// @Router /package/{id} [delete]
func (s *RoleServer) DeletePackage(c *gin.Context) {
	id := c.Param("id")
	userId := c.GetString("user_id")
	userEmail := c.GetString("user_email")

	// Capturar dados antes de deletar
	oldPkg, _ := s.handler.GetPackageWithModules(id)

	if err := s.handler.DeletePackage(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao remover pacote", "error": err.Error()})
		return
	}

	// 📝 LOG DE AUDITORIA
	if s.adminAuditHandler != nil {
		actorUUID, _ := uuid.Parse(userId)
		pkgUUID, _ := uuid.Parse(id)
		var oldValues map[string]interface{}
		if oldPkg != nil {
			oldValues = map[string]interface{}{"code_name": oldPkg.CodeName, "display_name": oldPkg.DisplayName}
		}
		go func() {
			s.adminAuditHandler.LogGenericAction(handler.AuditLogParams{
				ActorId:       actorUUID,
				ActorEmail:    userEmail,
				TargetId:      pkgUUID,
				Action:        models.AdminAuditActionDelete,
				EntityType:    models.AdminAuditEntityPackage,
				IsAdminZone:   true,
				OldValues:     oldValues,
				ChangedFields: []string{"*"},
				IpAddress:     c.ClientIP(),
				UserAgent:     c.Request.UserAgent(),
			})
		}()
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pacote removido com sucesso"})
}

// AddModuleToPackage godoc
// @Summary Adiciona um módulo a um pacote (Master Admin)
// @Tags Packages
// @Produce json
// @Param id path string true "ID do pacote"
// @Param moduleId path string true "ID do módulo"
// @Success 200 {object} map[string]string
// @Router /package/{id}/modules/{moduleId} [post]
func (s *RoleServer) AddModuleToPackage(c *gin.Context) {
	packageId := c.Param("id")
	moduleId := c.Param("moduleId")

	if err := s.handler.AddModuleToPackage(packageId, moduleId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao adicionar módulo", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Módulo adicionado ao pacote com sucesso"})
}

// RemoveModuleFromPackage godoc
// @Summary Remove um módulo de um pacote (Master Admin)
// @Tags Packages
// @Produce json
// @Param id path string true "ID do pacote"
// @Param moduleId path string true "ID do módulo"
// @Success 200 {object} map[string]string
// @Router /package/{id}/modules/{moduleId} [delete]
func (s *RoleServer) RemoveModuleFromPackage(c *gin.Context) {
	packageId := c.Param("id")
	moduleId := c.Param("moduleId")

	if err := s.handler.RemoveModuleFromPackage(packageId, moduleId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao remover módulo", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Módulo removido do pacote com sucesso"})
}

// SetPackageLimit godoc
// @Summary Define um limite para o pacote (Master Admin)
// @Tags Packages
// @Accept json
// @Produce json
// @Param id path string true "ID do pacote"
// @Param data body SetPackageLimitRequest true "Dados do limite"
// @Success 200 {object} map[string]string
// @Router /package/{id}/limits [post]
func (s *RoleServer) SetPackageLimit(c *gin.Context) {
	packageId := c.Param("id")

	var req SetPackageLimitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos", "error": err.Error()})
		return
	}

	if err := s.handler.SetPackageLimit(packageId, req.LimitType, req.LimitValue); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao definir limite", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Limite definido com sucesso"})
}

// GetPackageLimits godoc
// @Summary Lista limites de um pacote
// @Tags Packages
// @Produce json
// @Param id path string true "ID do pacote"
// @Success 200 {array} models.PackageLimit
// @Router /package/{id}/limits [get]
func (s *RoleServer) GetPackageLimits(c *gin.Context) {
	packageId := c.Param("id")

	limits, err := s.handler.GetPackageLimits(packageId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao buscar limites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": limits})
}

// CreateOrganizationSubscription godoc
// @Summary Cria uma assinatura para uma organização (Master Admin)
// @Tags Packages
// @Accept json
// @Produce json
// @Param orgId path string true "ID da organização"
// @Param data body SubscribeRequest true "Dados da assinatura"
// @Success 200 {object} map[string]string
// @Router /package/subscription/{orgId} [post]
func (s *RoleServer) CreateOrganizationSubscription(c *gin.Context) {
	orgId := c.Param("orgId")

	var req SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos", "error": err.Error()})
		return
	}

	if err := s.handler.SubscribeOrganization(orgId, req.PackageId, req.BillingCycle, req.CustomPrice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assinatura criada com sucesso"})
}

// UpdateOrganizationSubscription godoc
// @Summary Atualiza a assinatura de uma organização (Master Admin)
// @Tags Packages
// @Accept json
// @Produce json
// @Param orgId path string true "ID da organização"
// @Param data body UpdateSubscriptionRequest true "Dados atualizados"
// @Success 200 {object} map[string]string
// @Router /package/subscription/{orgId} [put]
func (s *RoleServer) UpdateOrganizationSubscription(c *gin.Context) {
	orgId := c.Param("orgId")

	var req UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos", "error": err.Error()})
		return
	}

	if err := s.handler.UpdateOrganizationSubscription(orgId, req.PackageId, req.BillingCycle, req.CustomPrice, req.Active); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao atualizar assinatura", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assinatura atualizada com sucesso"})
}

// CancelOrganizationSubscription godoc
// @Summary Cancela a assinatura de uma organização (Master Admin)
// @Tags Packages
// @Produce json
// @Param orgId path string true "ID da organização"
// @Success 200 {object} map[string]string
// @Router /package/subscription/{orgId} [delete]
func (s *RoleServer) CancelOrganizationSubscription(c *gin.Context) {
	orgId := c.Param("orgId")

	if err := s.handler.CancelOrganizationSubscription(orgId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao cancelar assinatura", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assinatura cancelada com sucesso"})
}

// ListAllSubscriptions godoc
// @Summary Lista todas as assinaturas (Master Admin)
// @Tags Packages
// @Produce json
// @Success 200 {array} models.OrganizationPackage
// @Router /package/subscriptions [get]
func (s *RoleServer) ListAllSubscriptions(c *gin.Context) {
	subscriptions, err := s.handler.ListAllSubscriptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao listar assinaturas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": subscriptions})
}

// ==================== Module CRUD (Master Admin) ====================

// CreateModule godoc
// @Summary Cria um novo módulo (Master Admin)
// @Tags Modules
// @Accept json
// @Produce json
// @Param module body models.Module true "Dados do módulo"
// @Success 201 {object} models.Module
// @Router /module [post]
func (s *RoleServer) CreateModule(c *gin.Context) {
	var module models.Module
	if err := c.ShouldBindJSON(&module); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos", "error": err.Error()})
		return
	}

	if module.Id == uuid.Nil {
		module.Id = uuid.New()
	}

	if err := s.handler.CreateModule(&module); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao criar módulo", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": module})
}

// UpdateModule godoc
// @Summary Atualiza um módulo (Master Admin)
// @Tags Modules
// @Accept json
// @Produce json
// @Param id path string true "ID do módulo"
// @Param module body models.Module true "Dados atualizados"
// @Success 200 {object} models.Module
// @Router /module/{id} [put]
func (s *RoleServer) UpdateModule(c *gin.Context) {
	id := c.Param("id")

	var module models.Module
	if err := c.ShouldBindJSON(&module); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos", "error": err.Error()})
		return
	}

	moduleUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID inválido"})
		return
	}
	module.Id = moduleUUID

	if err := s.handler.UpdateModule(&module); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao atualizar módulo", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": module})
}

// DeleteModule godoc
// @Summary Remove um módulo (Master Admin)
// @Tags Modules
// @Produce json
// @Param id path string true "ID do módulo"
// @Success 200 {object} map[string]string
// @Router /module/{id} [delete]
func (s *RoleServer) DeleteModule(c *gin.Context) {
	id := c.Param("id")

	if err := s.handler.DeleteModule(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao remover módulo", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Módulo removido com sucesso"})
}

// GetModule godoc
// @Summary Busca um módulo por ID
// @Tags Modules
// @Produce json
// @Param id path string true "ID do módulo"
// @Success 200 {object} models.Module
// @Router /module/{id} [get]
func (s *RoleServer) GetModule(c *gin.Context) {
	id := c.Param("id")

	module, err := s.handler.GetModule(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Módulo não encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": module})
}

// ==================== Permission CRUD (Master Admin) ====================

// CreatePermission godoc
// @Summary Cria uma nova permissão (Master Admin)
// @Tags Permissions
// @Accept json
// @Produce json
// @Param permission body models.Permission true "Dados da permissão"
// @Success 201 {object} models.Permission
// @Router /permission [post]
func (s *RoleServer) CreatePermission(c *gin.Context) {
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos", "error": err.Error()})
		return
	}

	if permission.Id == uuid.Nil {
		permission.Id = uuid.New()
	}

	if err := s.handler.CreatePermission(&permission); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao criar permissão", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": permission})
}

// UpdatePermission godoc
// @Summary Atualiza uma permissão (Master Admin)
// @Tags Permissions
// @Accept json
// @Produce json
// @Param id path string true "ID da permissão"
// @Param permission body models.Permission true "Dados atualizados"
// @Success 200 {object} models.Permission
// @Router /permission/{id} [put]
func (s *RoleServer) UpdatePermission(c *gin.Context) {
	id := c.Param("id")

	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos", "error": err.Error()})
		return
	}

	permUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID inválido"})
		return
	}
	permission.Id = permUUID

	if err := s.handler.UpdatePermission(&permission); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao atualizar permissão", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": permission})
}

// DeletePermission godoc
// @Summary Remove uma permissão (Master Admin)
// @Tags Permissions
// @Produce json
// @Param id path string true "ID da permissão"
// @Success 200 {object} map[string]string
// @Router /permission/{id} [delete]
func (s *RoleServer) DeletePermission(c *gin.Context) {
	id := c.Param("id")

	if err := s.handler.DeletePermission(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao remover permissão", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permissão removida com sucesso"})
}

// GetPermission godoc
// @Summary Busca uma permissão por ID
// @Tags Permissions
// @Produce json
// @Param id path string true "ID da permissão"
// @Success 200 {object} models.Permission
// @Router /permission/{id} [get]
func (s *RoleServer) GetPermission(c *gin.Context) {
	id := c.Param("id")

	permission, err := s.handler.GetPermission(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Permissão não encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": permission})
}

// ==================== Usage & Limits ====================

// GetUsageAndLimits godoc
// @Summary Retorna o uso atual e limites do plano da organização
// @Description Retorna informações sobre uso de recursos (mesas, usuários, produtos, reservas) e os limites do plano atual
// @Tags Packages
// @Produce json
// @Success 200 {object} handler.UsageLimitsResponse
// @Failure 403 {object} map[string]string "Organização sem plano ativo"
// @Router /package/usage-limits [get]
func (s *RoleServer) GetUsageAndLimits(c *gin.Context) {
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	if s.limitHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Handler de limites não configurado",
			"message": "Erro interno do servidor",
		})
		return
	}

	response, err := s.limitHandler.GetUsageAndLimits(orgId, projectId)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Erro ao buscar dados do plano",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
