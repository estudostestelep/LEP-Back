package server

import (
	"lep/handler"
	"lep/repositories/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoleServer struct {
	handler      *handler.RoleHandler
	limitHandler *handler.LimitHandler
}

func NewRoleServer(h *handler.RoleHandler) *RoleServer {
	return &RoleServer{handler: h}
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

	orgId := c.GetString("organization_id")

	// Se não especificado, associar à organização atual
	if role.OrganizationId == nil {
		parsed, err := uuid.Parse(orgId)
		if err == nil {
			role.OrganizationId = &parsed
		}
	}

	// Construir contexto da requisição para auditoria
	ctx := handler.BuildRequestContext(c)

	// Handler captura auditoria internamente
	if err := s.handler.CreateRoleWithContext(ctx, &role, orgId); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
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

	orgId := c.GetString("organization_id")

	// Construir contexto da requisição para auditoria
	ctx := handler.BuildRequestContext(c)

	// Handler captura estado anterior e auditoria internamente
	if err := s.handler.UpdateRoleWithContext(ctx, id, &role, orgId); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
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
	orgId := c.GetString("organization_id")

	// Construir contexto da requisição para auditoria
	ctx := handler.BuildRequestContext(c)

	// Handler captura estado anterior e auditoria internamente
	if err := s.handler.DeleteRoleWithContext(ctx, id, orgId); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
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

	// Obter dados do ator (quem está fazendo a atribuição)
	actorUserId := c.GetString("user_id")
	actorUserType := c.GetString("user_type")

	// Atribuir cargo baseado no tipo de usuário (admin ou client)
	if req.UserType == "admin" {
		adminRole := &models.AdminRole{
			AdminId: userUUID,
			RoleId:  roleUUID,
			Active:  true,
		}

		// Se tiver orgId, adicionar (para admin pode ser opcional)
		if orgId != "" {
			parsed, err := uuid.Parse(orgId)
			if err == nil {
				adminRole.OrganizationId = &parsed
			}
		}

		// Usar método com contexto para auditoria
		ctx := handler.BuildRequestContext(c)
		if err := s.handler.AssignRoleToAdminWithContext(ctx, adminRole, actorUserType); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
			return
		}
	} else {
		// Para client, orgId é obrigatório
		if orgId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "organization_id é obrigatório para clients"})
			return
		}

		orgUUID, err := uuid.Parse(orgId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "organization_id inválido"})
			return
		}

		clientRole := &models.ClientRole{
			ClientId:       userUUID,
			RoleId:         roleUUID,
			OrganizationId: orgUUID,
			Active:         true,
		}

		// Se tiver projectId, adicionar
		if req.ProjectId != "" {
			parsed, err := uuid.Parse(req.ProjectId)
			if err == nil {
				clientRole.ProjectId = &parsed
			}
		}

		if err := s.handler.AssignRoleToClient(clientRole, actorUserId, actorUserType); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
			return
		}
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

	orgId := c.GetString("organization_id")

	// Remover cargo baseado no tipo de usuário
	if req.UserType == "admin" {
		// Usar método com contexto para auditoria
		ctx := handler.BuildRequestContext(c)
		if err := s.handler.RemoveRoleFromAdminWithContext(ctx, req.UserId, req.RoleId); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
			return
		}
	} else {
		// Usar método com contexto para auditoria
		ctx := handler.BuildRequestContext(c)
		if err := s.handler.RemoveRoleFromClientWithContext(ctx, req.UserId, req.RoleId, orgId); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cargo removido com sucesso"})
}

// GetUserRoles godoc
// @Summary Lista cargos de um usuário (client ou admin)
// @Tags Roles
// @Produce json
// @Param userId path string true "ID do usuário"
// @Param user_type query string true "Tipo de usuário (admin/client)"
// @Success 200 {array} models.ClientRole
// @Router /role/user/{userId} [get]
func (s *RoleServer) GetUserRoles(c *gin.Context) {
	userId := c.Param("userId")
	orgId := c.GetString("organization_id")
	userType := c.Query("user_type")

	if userType == "admin" {
		roles, err := s.handler.GetAdminRoles(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao buscar cargos do admin"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": roles})
	} else {
		roles, err := s.handler.GetClientRoles(userId, orgId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao buscar cargos do client"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": roles})
	}
}

// GetUserRolesWithDetails godoc
// @Summary Lista cargos de um usuário com detalhes de permissões
// @Tags Roles
// @Produce json
// @Param userId path string true "ID do usuário"
// @Param user_type query string true "Tipo de usuário (admin/client)"
// @Success 200 {array} models.RoleWithPermissions
// @Router /role/user/{userId}/details [get]
func (s *RoleServer) GetUserRolesWithDetails(c *gin.Context) {
	userId := c.Param("userId")
	orgId := c.GetString("organization_id")
	userType := c.Query("user_type")

	if userType == "admin" {
		roles, err := s.handler.GetAdminRolesWithPermissions(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao buscar detalhes dos cargos"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": roles})
	} else {
		roles, err := s.handler.GetClientRolesWithPermissions(userId, orgId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao buscar detalhes dos cargos"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": roles})
	}
}

// ==================== Permission Management ====================

// AddPermissionToRole godoc
// @Summary Adiciona uma permissão a um cargo
// @Tags Roles
// @Accept json
// @Produce json
// @Param data body AddPermissionRequest true "Dados da permissão"
// @Success 200 {object} map[string]string
// @Router /role/permission [post]
func (s *RoleServer) AddPermissionToRole(c *gin.Context) {
	var req struct {
		RoleId       string `json:"role_id" binding:"required"`
		PermissionId string `json:"permission_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos", "error": err.Error()})
		return
	}

	if err := s.handler.AddPermissionToRole(req.RoleId, req.PermissionId); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permissão adicionada com sucesso"})
}

// GetRolePermissions godoc
// @Summary Lista permissões de um cargo
// @Tags Roles
// @Produce json
// @Param id path string true "ID do cargo"
// @Success 200 {array} models.Permission
// @Router /role/{id}/permissions [get]
func (s *RoleServer) GetRolePermissions(c *gin.Context) {
	roleId := c.Param("id")

	permissions, err := s.handler.GetRolePermissions(roleId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao buscar permissões"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": permissions})
}

// ==================== Permission Checking ====================

// CheckPermission godoc
// @Summary Verifica se o usuário tem uma permissão
// @Tags Roles
// @Produce json
// @Param permission query string true "Código da permissão (formato: module:action)"
// @Success 200 {object} PermissionCheckResponse
// @Router /role/check [get]
func (s *RoleServer) CheckPermission(c *gin.Context) {
	userId := c.GetString("user_id")
	userType := c.GetString("user_type")
	permission := c.Query("permission")

	hasPermission, err := s.handler.UserHasPermission(userId, userType, permission)
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
	userType := c.GetString("user_type")
	orgId := c.GetString("organization_id")

	// Verificar se é Master Admin
	isMasterAdmin, _ := s.handler.IsMasterAdmin(userId, userType)

	// Buscar cargos com detalhes baseado no tipo de usuário
	var roles []models.RoleWithPermissions
	var err error

	if userType == "admin" {
		roles, err = s.handler.GetAdminRolesWithPermissions(userId)
	} else {
		roles, err = s.handler.GetClientRolesWithPermissions(userId, orgId)
	}

	// Se é Master Admin e não tem cargos, simular cargo "owner"
	if isMasterAdmin && (err != nil || len(roles) == 0) {
		ownerRole, roleErr := s.handler.GetRoleByName("owner")
		if roleErr == nil && ownerRole != nil {
			codes, _ := s.handler.GetRolePermissionCodes(ownerRole.Id.String())
			roles = []models.RoleWithPermissions{{
				Role:            *ownerRole,
				PermissionCodes: codes,
			}}
			err = nil
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao buscar permissões"})
		return
	}

	// Buscar nível de hierarquia
	hierarchyLevel, _ := s.handler.GetUserHierarchyLevel(userId, userType)
	if isMasterAdmin && hierarchyLevel == 0 {
		hierarchyLevel = 10
	}

	// Agregar permissões efetivas (lista de códigos únicos)
	permSet := make(map[string]bool)
	for _, role := range roles {
		for _, code := range role.PermissionCodes {
			permSet[code] = true
		}
	}

	effectivePermissions := make([]string, 0, len(permSet))
	for code := range permSet {
		effectivePermissions = append(effectivePermissions, code)
	}

	// Buscar módulos disponíveis
	modules, _ := s.handler.GetOrganizationModules(orgId)

	c.JSON(http.StatusOK, gin.H{
		"roles":                 roles,
		"effective_permissions": effectivePermissions,
		"hierarchy_level":       hierarchyLevel,
		"available_modules":     modules,
		"is_master_admin":       isMasterAdmin,
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

	packages, err := s.handler.ListPlans(publicOnly)
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

	pkg, err := s.handler.GetPlanWithModules(id)
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

	if err := s.handler.SubscribeOrganization(orgId, req.PlanId, req.BillingCycle, req.CustomPrice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assinatura realizada com sucesso"})
}

// ==================== Request/Response Types ====================

type AssignRoleRequest struct {
	UserId    string `json:"user_id" binding:"required"`
	UserType  string `json:"user_type" binding:"required,oneof=admin client"` // "admin" ou "client"
	RoleId    string `json:"role_id" binding:"required"`
	ProjectId string `json:"project_id,omitempty"`
}

type RemoveRoleRequest struct {
	UserId   string `json:"user_id" binding:"required"`
	UserType string `json:"user_type" binding:"required,oneof=admin client"` // "admin" ou "client"
	RoleId   string `json:"role_id" binding:"required"`
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
	Roles                []models.RoleWithPermissions `json:"roles"`
	EffectivePermissions []string                     `json:"effective_permissions"` // Lista de permissões no formato module:action
	HierarchyLevel       int                          `json:"hierarchy_level"`
	AvailableModules     []models.Module              `json:"available_modules"`
}

type SubscribeRequest struct {
	PlanId       string   `json:"plan_id" binding:"required"`
	BillingCycle string   `json:"billing_cycle" binding:"required,oneof=monthly yearly"`
	CustomPrice  *float64 `json:"custom_price,omitempty"`
}

type UpdateSubscriptionRequest struct {
	PlanId       string   `json:"plan_id,omitempty"`
	BillingCycle string   `json:"billing_cycle,omitempty"`
	CustomPrice  *float64 `json:"custom_price,omitempty"`
	Active       *bool    `json:"active,omitempty"`
}

type SetPlanLimitRequest struct {
	LimitType  string `json:"limit_type" binding:"required"`
	LimitValue int    `json:"limit_value" binding:"required"`
}

// ==================== Package CRUD (Master Admin) ====================

// CreatePackage godoc
// @Summary Cria um novo plano (Master Admin)
// @Tags Packages
// @Accept json
// @Produce json
// @Param package body models.Plan true "Dados do plano"
// @Success 201 {object} models.Plan
// @Router /package [post]
func (s *RoleServer) CreatePackage(c *gin.Context) {
	var plan models.Plan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos", "error": err.Error()})
		return
	}

	if plan.Id == uuid.Nil {
		plan.Id = uuid.New()
	}

	// Construir contexto da requisição para auditoria
	ctx := handler.BuildRequestContext(c)

	// Handler captura auditoria internamente
	if err := s.handler.CreatePlanWithContext(ctx, &plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao criar plano", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": plan})
}

// UpdatePackage godoc
// @Summary Atualiza um plano (Master Admin)
// @Tags Packages
// @Accept json
// @Produce json
// @Param id path string true "ID do plano"
// @Param package body models.Plan true "Dados atualizados"
// @Success 200 {object} models.Plan
// @Router /package/{id} [put]
func (s *RoleServer) UpdatePackage(c *gin.Context) {
	id := c.Param("id")

	var plan models.Plan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos", "error": err.Error()})
		return
	}

	planUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID inválido"})
		return
	}
	plan.Id = planUUID

	// Construir contexto da requisição para auditoria
	ctx := handler.BuildRequestContext(c)

	// Handler captura auditoria internamente
	if err := s.handler.UpdatePlanWithContext(ctx, &plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao atualizar plano", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": plan})
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

	// Construir contexto da requisição para auditoria
	ctx := handler.BuildRequestContext(c)

	// Handler captura estado anterior e auditoria internamente
	if err := s.handler.DeletePlanWithContext(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao remover pacote", "error": err.Error()})
		return
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

	if err := s.handler.AddModuleToPlan(packageId, moduleId); err != nil {
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

	if err := s.handler.RemoveModuleFromPlan(packageId, moduleId); err != nil {
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
// @Param data body SetPlanLimitRequest true "Dados do limite"
// @Success 200 {object} map[string]string
// @Router /package/{id}/limits [post]
func (s *RoleServer) SetPackageLimit(c *gin.Context) {
	packageId := c.Param("id")

	var req SetPlanLimitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos", "error": err.Error()})
		return
	}

	if err := s.handler.SetPlanLimit(packageId, req.LimitType, req.LimitValue); err != nil {
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

	limits, err := s.handler.GetPlanLimits(packageId)
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

	if err := s.handler.SubscribeOrganization(orgId, req.PlanId, req.BillingCycle, req.CustomPrice); err != nil {
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

	if err := s.handler.UpdateOrganizationSubscription(orgId, req.PlanId, req.BillingCycle, req.CustomPrice, req.Active); err != nil {
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

// DeleteOrganizationSubscription godoc
// @Summary Exclui permanentemente a assinatura de uma organização (Master Admin)
// @Tags Packages
// @Produce json
// @Param orgId path string true "ID da organização"
// @Success 200 {object} map[string]string
// @Router /package/subscription/{orgId}/delete [delete]
func (s *RoleServer) DeleteOrganizationSubscription(c *gin.Context) {
	orgId := c.Param("orgId")

	if err := s.handler.DeleteOrganizationSubscription(orgId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao excluir assinatura", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assinatura excluída com sucesso"})
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
