package server

import (
	"lep/handler"
	"lep/repositories/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoleServer struct {
	handler *handler.RoleHandler
}

func NewRoleServer(h *handler.RoleHandler) *RoleServer {
	return &RoleServer{handler: h}
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
	orgId := c.GetString("organization_id")

	// Se não especificado, associar à organização atual
	if role.OrganizationId == nil {
		orgUUID, err := uuid.Parse(orgId)
		if err == nil {
			role.OrganizationId = &orgUUID
		}
	}

	if err := s.handler.CreateRole(&role, userId, orgId); err != nil {
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

	userId := c.GetString("user_id")
	orgId := c.GetString("organization_id")

	if err := s.handler.UpdateRole(&role, userId, orgId); err != nil {
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
	userId := c.GetString("user_id")
	orgId := c.GetString("organization_id")

	if err := s.handler.DeleteRole(id, userId, orgId); err != nil {
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

	actorUserId := c.GetString("user_id")
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

	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID da organização inválido"})
		return
	}

	userRole := &models.UserRole{
		UserId:         userUUID,
		RoleId:         roleUUID,
		OrganizationId: orgUUID,
		Active:         true,
	}

	// Se tiver projectId, adicionar
	if req.ProjectId != "" {
		projUUID, err := uuid.Parse(req.ProjectId)
		if err == nil {
			userRole.ProjectId = &projUUID
		}
	}

	if err := s.handler.AssignRoleToUser(userRole, actorUserId); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
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
	orgId := c.GetString("organization_id")

	if err := s.handler.RemoveRoleFromUser(req.UserId, req.RoleId, orgId, actorUserId); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
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
// @Param roleId path string true "ID do cargo"
// @Success 200 {array} models.RolePermissionLevel
// @Router /role/{roleId}/permissions [get]
func (s *RoleServer) GetRolePermissions(c *gin.Context) {
	roleId := c.Param("roleId")

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
