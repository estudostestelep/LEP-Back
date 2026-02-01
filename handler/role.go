package handler

import (
	"fmt"
	"lep/constants"
	"lep/repositories"
	"lep/repositories/models"

	"github.com/google/uuid"
)

// RoleHandler gerencia roles, permissões e planos
// Implementa IAuthorizationHandler para uso no middleware
type RoleHandler struct {
	roleRepo          repositories.IRoleRepository
	permissionRepo    repositories.IPermissionRepository
	moduleRepo        repositories.IModuleRepository
	planRepo          repositories.IPlanRepository
	adminRepo         repositories.IAdminRepository
	clientRepo        repositories.IClientRepository
	adminAuditHandler IAdminAuditLogHandler
}

func NewRoleHandler(
	roleRepo repositories.IRoleRepository,
	permissionRepo repositories.IPermissionRepository,
	moduleRepo repositories.IModuleRepository,
	planRepo repositories.IPlanRepository,
	adminRepo repositories.IAdminRepository,
	clientRepo repositories.IClientRepository,
) *RoleHandler {
	return &RoleHandler{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		moduleRepo:     moduleRepo,
		planRepo:       planRepo,
		adminRepo:      adminRepo,
		clientRepo:     clientRepo,
	}
}

// SetAdminAuditHandler configura o handler de auditoria (injetado separadamente)
func (h *RoleHandler) SetAdminAuditHandler(handler IAdminAuditLogHandler) {
	h.adminAuditHandler = handler
}

// ==================== Role CRUD ====================

// CreateRole cria um novo cargo
func (h *RoleHandler) CreateRole(role *models.Role, actorUserId, actorUserType, orgId string) error {
	// Validar hierarquia - usuário só pode criar cargos com nível igual ou menor
	actorLevel, err := h.roleRepo.GetUserHierarchyLevel(actorUserId, actorUserType, orgId)
	if err != nil {
		return fmt.Errorf("erro ao verificar hierarquia: %w", err)
	}

	if role.HierarchyLevel > actorLevel {
		return fmt.Errorf("você não pode criar cargos com nível de hierarquia maior que o seu (%d)", actorLevel)
	}

	return h.roleRepo.Create(role)
}

// GetRole busca um cargo por ID
func (h *RoleHandler) GetRole(id string) (*models.Role, error) {
	return h.roleRepo.GetById(id)
}

// UpdateRole atualiza um cargo
func (h *RoleHandler) UpdateRole(role *models.Role, actorUserId, actorUserType, orgId string) error {
	// Buscar o cargo alvo para verificar hierarquia
	targetRole, err := h.roleRepo.GetById(role.Id.String())
	if err != nil {
		return fmt.Errorf("cargo não encontrado: %w", err)
	}

	// Buscar nível de hierarquia do ator
	actorLevel, err := h.roleRepo.GetUserHierarchyLevel(actorUserId, actorUserType, orgId)
	if err != nil {
		return fmt.Errorf("erro ao verificar hierarquia: %w", err)
	}

	// Master admin pode tudo
	if actorLevel >= constants.HierarchyMasterAdmin {
		return h.roleRepo.Update(role)
	}

	// Ator pode gerenciar roles com nível igual ou menor
	if actorLevel < targetRole.HierarchyLevel {
		return fmt.Errorf("você não tem permissão para modificar este cargo")
	}

	return h.roleRepo.Update(role)
}

// DeleteRole remove um cargo
func (h *RoleHandler) DeleteRole(id, actorUserId, actorUserType, orgId string) error {
	// Verificar se é um cargo do sistema
	role, err := h.roleRepo.GetById(id)
	if err != nil {
		return err
	}

	if role.IsSystem {
		return fmt.Errorf("cargos do sistema não podem ser excluídos")
	}

	// Buscar nível de hierarquia do ator
	actorLevel, err := h.roleRepo.GetUserHierarchyLevel(actorUserId, actorUserType, orgId)
	if err != nil {
		return fmt.Errorf("erro ao verificar hierarquia: %w", err)
	}

	// Master admin pode tudo
	if actorLevel >= constants.HierarchyMasterAdmin {
		return h.roleRepo.Delete(id)
	}

	// Ator pode gerenciar roles com nível igual ou menor
	if actorLevel < role.HierarchyLevel {
		return fmt.Errorf("você não tem permissão para excluir este cargo")
	}

	return h.roleRepo.Delete(id)
}

// ListRoles lista cargos com base no escopo
func (h *RoleHandler) ListRoles(scope, orgId string) ([]models.Role, error) {
	if scope != "" {
		return h.roleRepo.ListByScope(scope)
	}
	if orgId != "" {
		return h.roleRepo.ListByOrganization(orgId)
	}
	return h.roleRepo.List()
}

// ListSystemRoles lista apenas cargos do sistema
func (h *RoleHandler) ListSystemRoles() ([]models.Role, error) {
	return h.roleRepo.ListSystemRoles()
}

// ==================== IAuthorizationHandler Implementation ====================

// IsMasterAdmin verifica se usuário é master admin (hierarchy >= 10)
func (h *RoleHandler) IsMasterAdmin(userId, userType string) (bool, error) {
	return h.roleRepo.IsMasterAdmin(userId, userType)
}

// GetUserHierarchyLevel retorna o maior nível de hierarquia do usuário
func (h *RoleHandler) GetUserHierarchyLevel(userId, userType string) (int, error) {
	return h.roleRepo.GetUserHierarchyLevel(userId, userType, "")
}

// UserHasPermission verifica se usuário tem uma permissão específica via roles
func (h *RoleHandler) UserHasPermission(userId, userType, permission string) (bool, error) {
	return h.roleRepo.UserHasPermission(userId, userType, permission)
}

// OrganizationHasModule verifica se organização tem acesso ao módulo via plan
func (h *RoleHandler) OrganizationHasModule(orgId, moduleCode string) (bool, error) {
	return h.planRepo.OrganizationHasModule(orgId, moduleCode)
}

// CanManageUser verifica se actor pode gerenciar target baseado em hierarquia
func (h *RoleHandler) CanManageUser(actorId, targetId, userType string) (bool, error) {
	actorLevel, err := h.roleRepo.GetUserHierarchyLevel(actorId, userType, "")
	if err != nil {
		return false, err
	}

	targetLevel, err := h.roleRepo.GetUserHierarchyLevel(targetId, userType, "")
	if err != nil {
		return false, err
	}

	return actorLevel > targetLevel, nil
}

// ==================== Client-Role Assignment ====================

// AssignRoleToClient atribui um cargo a um cliente
func (h *RoleHandler) AssignRoleToClient(clientRole *models.ClientRole, actorUserId, actorUserType string) error {
	// Validar se o ator pode atribuir este cargo
	role, err := h.roleRepo.GetById(clientRole.RoleId.String())
	if err != nil {
		return fmt.Errorf("erro ao buscar cargo: %w", err)
	}

	orgId := ""
	if clientRole.OrganizationId != uuid.Nil {
		orgId = clientRole.OrganizationId.String()
	}

	// Verificar hierarquia do ator
	actorLevel, err := h.roleRepo.GetUserHierarchyLevel(actorUserId, actorUserType, orgId)
	if err != nil {
		return fmt.Errorf("erro ao verificar hierarquia: %w", err)
	}

	// Ator deve ter hierarquia maior ou igual ao cargo que está atribuindo
	if role.HierarchyLevel > actorLevel && actorLevel < constants.HierarchyMasterAdmin {
		return fmt.Errorf("você não tem permissão para atribuir este cargo (nível %d)", role.HierarchyLevel)
	}

	return h.roleRepo.AssignRoleToClient(clientRole)
}

// RemoveRoleFromClient remove um cargo de um cliente
func (h *RoleHandler) RemoveRoleFromClient(clientId, roleId, orgId, actorUserId, actorUserType string) error {
	// Validar se o ator pode remover este cargo
	role, err := h.roleRepo.GetById(roleId)
	if err != nil {
		return fmt.Errorf("erro ao buscar cargo: %w", err)
	}

	// Verificar hierarquia do ator
	actorLevel, err := h.roleRepo.GetUserHierarchyLevel(actorUserId, actorUserType, orgId)
	if err != nil {
		return fmt.Errorf("erro ao verificar hierarquia: %w", err)
	}

	// Ator deve ter hierarquia maior ou igual ao cargo que está removendo
	if role.HierarchyLevel > actorLevel && actorLevel < constants.HierarchyMasterAdmin {
		return fmt.Errorf("você não tem permissão para remover este cargo")
	}

	return h.roleRepo.RemoveRoleFromClient(clientId, roleId, orgId)
}

// GetClientRoles retorna todos os cargos de um cliente
func (h *RoleHandler) GetClientRoles(clientId, orgId string) ([]models.ClientRole, error) {
	return h.roleRepo.GetClientRoles(clientId, orgId)
}

// ==================== Admin-Role Assignment ====================

// AssignRoleToAdmin atribui um cargo a um admin
func (h *RoleHandler) AssignRoleToAdmin(adminRole *models.AdminRole, actorUserId, actorUserType string) error {
	// Validar se o ator pode atribuir este cargo
	role, err := h.roleRepo.GetById(adminRole.RoleId.String())
	if err != nil {
		return fmt.Errorf("erro ao buscar cargo: %w", err)
	}

	orgId := ""
	if adminRole.OrganizationId != nil {
		orgId = adminRole.OrganizationId.String()
	}

	// Verificar hierarquia do ator
	actorLevel, err := h.roleRepo.GetUserHierarchyLevel(actorUserId, actorUserType, orgId)
	if err != nil {
		return fmt.Errorf("erro ao verificar hierarquia: %w", err)
	}

	// Ator deve ter hierarquia maior ou igual ao cargo que está atribuindo
	if role.HierarchyLevel > actorLevel && actorLevel < constants.HierarchyMasterAdmin {
		return fmt.Errorf("você não tem permissão para atribuir este cargo (nível %d)", role.HierarchyLevel)
	}

	return h.roleRepo.AssignRoleToAdmin(adminRole)
}

// RemoveRoleFromAdmin remove um cargo de um admin
func (h *RoleHandler) RemoveRoleFromAdmin(adminId, roleId, actorUserId, actorUserType string) error {
	// Validar se o ator pode remover este cargo
	role, err := h.roleRepo.GetById(roleId)
	if err != nil {
		return fmt.Errorf("erro ao buscar cargo: %w", err)
	}

	// Verificar hierarquia do ator (admins não tem orgId obrigatório)
	actorLevel, err := h.roleRepo.GetUserHierarchyLevel(actorUserId, actorUserType, "")
	if err != nil {
		return fmt.Errorf("erro ao verificar hierarquia: %w", err)
	}

	// Ator deve ter hierarquia maior ou igual ao cargo que está removendo
	if role.HierarchyLevel > actorLevel && actorLevel < constants.HierarchyMasterAdmin {
		return fmt.Errorf("você não tem permissão para remover este cargo")
	}

	return h.roleRepo.RemoveRoleFromAdmin(adminId, roleId)
}

// GetAdminRoles retorna todos os cargos de um admin
func (h *RoleHandler) GetAdminRoles(adminId string) ([]models.AdminRole, error) {
	return h.roleRepo.GetAdminRoles(adminId)
}

// GetAdminRolesWithPermissions retorna cargos de admin com suas permissões
func (h *RoleHandler) GetAdminRolesWithPermissions(adminId string) ([]models.RoleWithPermissions, error) {
	return h.roleRepo.GetAdminRolesWithPermissions(adminId)
}

// GetClientRolesWithPermissions retorna cargos de client com suas permissões
func (h *RoleHandler) GetClientRolesWithPermissions(clientId, orgId string) ([]models.RoleWithPermissions, error) {
	return h.roleRepo.GetClientRolesWithPermissions(clientId, orgId)
}

// AddPermissionToRole adiciona uma permissão a um cargo
func (h *RoleHandler) AddPermissionToRole(roleId, permissionId string) error {
	return h.roleRepo.AddPermissionToRole(roleId, permissionId)
}

// RemovePermissionFromRole remove uma permissão de um cargo
func (h *RoleHandler) RemovePermissionFromRole(roleId, permissionId string) error {
	return h.roleRepo.RemovePermissionFromRole(roleId, permissionId)
}

// GetRolePermissions retorna todas as permissões de um cargo
func (h *RoleHandler) GetRolePermissions(roleId string) ([]models.Permission, error) {
	return h.roleRepo.GetRolePermissions(roleId)
}

// GetRolePermissionCodes retorna os códigos das permissões de um cargo
func (h *RoleHandler) GetRolePermissionCodes(roleId string) ([]string, error) {
	return h.roleRepo.GetRolePermissionCodes(roleId)
}

// GetRoleByName busca um cargo pelo nome
func (h *RoleHandler) GetRoleByName(name string) (*models.Role, error) {
	return h.roleRepo.GetByName(name)
}

// HasModuleAccess verifica se a organização tem acesso a um módulo
func (h *RoleHandler) HasModuleAccess(orgId, moduleCodeName string) (bool, error) {
	return h.planRepo.OrganizationHasModule(orgId, moduleCodeName)
}

// CanManageUserInOrg verifica se um usuário pode gerenciar outro na mesma organização
func (h *RoleHandler) CanManageUserInOrg(actorUserId, targetUserId, userType, orgId string) (bool, error) {
	return h.roleRepo.CanManageUser(actorUserId, targetUserId, userType, orgId)
}

// ==================== Module & Permission Listing ====================

// ListModules lista todos os módulos
func (h *RoleHandler) ListModules(scope string) ([]models.Module, error) {
	if scope != "" {
		return h.moduleRepo.ListByScope(scope)
	}
	return h.moduleRepo.List()
}

// ListModulesWithPermissions lista módulos com suas permissões
func (h *RoleHandler) ListModulesWithPermissions() ([]models.Module, error) {
	return h.moduleRepo.ListWithPermissions()
}

// ListPermissions lista todas as permissões
func (h *RoleHandler) ListPermissions(moduleId string) ([]models.Permission, error) {
	if moduleId != "" {
		return h.permissionRepo.ListByModule(moduleId)
	}
	return h.permissionRepo.List()
}

// GetOrganizationModules retorna módulos disponíveis para a organização
func (h *RoleHandler) GetOrganizationModules(orgId string) ([]models.Module, error) {
	return h.moduleRepo.GetModulesForOrganization(orgId)
}

// ==================== Plan Management ====================

// ListPlans lista todos os planos
func (h *RoleHandler) ListPlans(publicOnly bool) ([]models.Plan, error) {
	if publicOnly {
		return h.planRepo.ListPublic()
	}
	return h.planRepo.List()
}

// GetPlanWithModules retorna um plano com seus módulos
func (h *RoleHandler) GetPlanWithModules(id string) (*models.Plan, error) {
	return h.planRepo.GetPlanWithModules(id)
}

// SubscribeOrganization inscreve uma organização em um plano
func (h *RoleHandler) SubscribeOrganization(orgId, planId string, billingCycle string, customPrice *float64) error {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return fmt.Errorf("ID da organização inválido: %w", err)
	}

	planUUID, err := uuid.Parse(planId)
	if err != nil {
		return fmt.Errorf("ID do plano inválido: %w", err)
	}

	// Buscar plano para obter preço
	plan, err := h.planRepo.GetById(planId)
	if err != nil {
		return fmt.Errorf("plano não encontrado: %w", err)
	}

	price := plan.PriceMonthly
	if billingCycle == "yearly" {
		price = plan.PriceYearly
	}
	if customPrice != nil {
		price = *customPrice
	}

	orgPlan := &models.OrganizationPlan{
		OrganizationId: orgUUID,
		PlanId:         planUUID,
		BillingCycle:   billingCycle,
		CustomPrice:    customPrice,
		Active:         true,
	}

	// Se o preço customizado não foi definido, usar o preço padrão
	if customPrice == nil {
		orgPlan.CustomPrice = &price
	}

	return h.planRepo.SubscribeOrganization(orgPlan)
}

// GetOrganizationSubscription retorna a assinatura ativa da organização
func (h *RoleHandler) GetOrganizationSubscription(orgId string) (*models.OrganizationPlan, error) {
	return h.planRepo.GetOrganizationPlan(orgId)
}

// ==================== Plan CRUD (Master Admin) ====================

// CreatePlan cria um novo plano
func (h *RoleHandler) CreatePlan(plan *models.Plan) error {
	return h.planRepo.Create(plan)
}

// UpdatePlan atualiza um plano
func (h *RoleHandler) UpdatePlan(plan *models.Plan) error {
	return h.planRepo.Update(plan)
}

// DeletePlan remove um plano
func (h *RoleHandler) DeletePlan(id string) error {
	return h.planRepo.Delete(id)
}

// AddModuleToPlan adiciona um módulo a um plano
func (h *RoleHandler) AddModuleToPlan(planId, moduleId string) error {
	return h.planRepo.AddModuleToPlan(planId, moduleId)
}

// RemoveModuleFromPlan remove um módulo de um plano
func (h *RoleHandler) RemoveModuleFromPlan(planId, moduleId string) error {
	return h.planRepo.RemoveModuleFromPlan(planId, moduleId)
}

// SetPlanLimit define um limite para o plano
func (h *RoleHandler) SetPlanLimit(planId, limitType string, limitValue int) error {
	return h.planRepo.SetPlanLimit(planId, limitType, limitValue)
}

// GetPlanLimits retorna os limites de um plano
func (h *RoleHandler) GetPlanLimits(planId string) ([]models.PlanLimit, error) {
	return h.planRepo.GetPlanLimits(planId)
}

// UpdateOrganizationSubscription atualiza a assinatura de uma organização
func (h *RoleHandler) UpdateOrganizationSubscription(orgId, planId, billingCycle string, customPrice *float64, active *bool) error {
	// Buscar assinatura existente
	existingSubscription, err := h.planRepo.GetOrganizationPlan(orgId)
	if err != nil {
		return fmt.Errorf("assinatura não encontrada: %w", err)
	}

	// Atualizar campos se fornecidos
	if planId != "" {
		planUUID, err := uuid.Parse(planId)
		if err != nil {
			return fmt.Errorf("ID do plano inválido: %w", err)
		}
		existingSubscription.PlanId = planUUID
	}

	if billingCycle != "" {
		existingSubscription.BillingCycle = billingCycle
	}

	if customPrice != nil {
		existingSubscription.CustomPrice = customPrice
	}

	if active != nil {
		existingSubscription.Active = *active
	}

	return h.planRepo.UpdateOrganizationPlan(existingSubscription)
}

// CancelOrganizationSubscription cancela a assinatura de uma organização
func (h *RoleHandler) CancelOrganizationSubscription(orgId string) error {
	return h.planRepo.CancelOrganizationPlan(orgId)
}

// DeleteOrganizationSubscription exclui permanentemente a assinatura de uma organização
func (h *RoleHandler) DeleteOrganizationSubscription(orgId string) error {
	return h.planRepo.DeleteOrganizationPlan(orgId)
}

// ListAllSubscriptions lista todas as assinaturas ativas
func (h *RoleHandler) ListAllSubscriptions() ([]models.OrganizationPlan, error) {
	return h.planRepo.ListAllSubscriptions()
}

// ==================== Module CRUD (Master Admin) ====================

// CreateModule cria um novo módulo
func (h *RoleHandler) CreateModule(module *models.Module) error {
	return h.moduleRepo.Create(module)
}

// UpdateModule atualiza um módulo
func (h *RoleHandler) UpdateModule(module *models.Module) error {
	return h.moduleRepo.Update(module)
}

// DeleteModule remove um módulo
func (h *RoleHandler) DeleteModule(id string) error {
	return h.moduleRepo.Delete(id)
}

// GetModule busca um módulo por ID
func (h *RoleHandler) GetModule(id string) (*models.Module, error) {
	return h.moduleRepo.GetById(id)
}

// ==================== Permission CRUD (Master Admin) ====================

// CreatePermission cria uma nova permissão
func (h *RoleHandler) CreatePermission(permission *models.Permission) error {
	return h.permissionRepo.Create(permission)
}

// UpdatePermission atualiza uma permissão
func (h *RoleHandler) UpdatePermission(permission *models.Permission) error {
	return h.permissionRepo.Update(permission)
}

// DeletePermission remove uma permissão
func (h *RoleHandler) DeletePermission(id string) error {
	return h.permissionRepo.Delete(id)
}

// GetPermission busca uma permissão por ID
func (h *RoleHandler) GetPermission(id string) (*models.Permission, error) {
	return h.permissionRepo.GetById(id)
}

// ==================== Métodos WithContext para Auditoria ====================

// CreateRoleWithContext cria um cargo e registra auditoria
func (h *RoleHandler) CreateRoleWithContext(ctx *RequestContext, role *models.Role, orgId string) error {
	// Executar criação normal
	if err := h.CreateRole(role, ctx.UserId.String(), ctx.UserType, orgId); err != nil {
		return err
	}

	// Registrar auditoria
	if h.adminAuditHandler != nil {
		go func() {
			var orgUUID *uuid.UUID
			if orgId != "" {
				parsed, _ := uuid.Parse(orgId)
				orgUUID = &parsed
			}
			h.adminAuditHandler.LogGenericAction(AuditLogParams{
				ActorId:       ctx.UserId,
				ActorEmail:    ctx.UserEmail,
				TargetId:      role.Id,
				Action:        models.AdminAuditActionCreate,
				EntityType:    models.AdminAuditEntityRole,
				OrgId:         orgUUID,
				ProjectId:     nil,
				IsAdminZone:   true,
				NewValues:     map[string]interface{}{"name": role.Name, "display_name": role.DisplayName, "scope": role.Scope},
				ChangedFields: []string{"*"},
				IpAddress:     ctx.IpAddress,
				UserAgent:     ctx.UserAgent,
			})
		}()
	}

	return nil
}

// UpdateRoleWithContext atualiza um cargo e registra auditoria
func (h *RoleHandler) UpdateRoleWithContext(ctx *RequestContext, roleId string, role *models.Role, orgId string) error {
	// Capturar estado anterior para auditoria
	oldRole, _ := h.GetRole(roleId)

	// Executar atualização normal
	if err := h.UpdateRole(role, ctx.UserId.String(), ctx.UserType, orgId); err != nil {
		return err
	}

	// Registrar auditoria
	if h.adminAuditHandler != nil && oldRole != nil {
		go func() {
			var orgUUID *uuid.UUID
			if orgId != "" {
				parsed, _ := uuid.Parse(orgId)
				orgUUID = &parsed
			}
			h.adminAuditHandler.LogGenericAction(AuditLogParams{
				ActorId:       ctx.UserId,
				ActorEmail:    ctx.UserEmail,
				TargetId:      role.Id,
				Action:        models.AdminAuditActionUpdate,
				EntityType:    models.AdminAuditEntityRole,
				OrgId:         orgUUID,
				ProjectId:     nil,
				IsAdminZone:   true,
				OldValues:     map[string]interface{}{"name": oldRole.Name, "display_name": oldRole.DisplayName},
				NewValues:     map[string]interface{}{"name": role.Name, "display_name": role.DisplayName},
				ChangedFields: []string{"name", "display_name", "permissions"},
				IpAddress:     ctx.IpAddress,
				UserAgent:     ctx.UserAgent,
			})
		}()
	}

	return nil
}

// DeleteRoleWithContext remove um cargo e registra auditoria
func (h *RoleHandler) DeleteRoleWithContext(ctx *RequestContext, roleId, orgId string) error {
	// Capturar estado anterior para auditoria
	oldRole, _ := h.GetRole(roleId)

	// Executar exclusão normal
	if err := h.DeleteRole(roleId, ctx.UserId.String(), ctx.UserType, orgId); err != nil {
		return err
	}

	// Registrar auditoria
	if h.adminAuditHandler != nil && oldRole != nil {
		go func() {
			roleUUID, _ := uuid.Parse(roleId)
			var orgUUID *uuid.UUID
			if orgId != "" {
				parsed, _ := uuid.Parse(orgId)
				orgUUID = &parsed
			}
			h.adminAuditHandler.LogGenericAction(AuditLogParams{
				ActorId:       ctx.UserId,
				ActorEmail:    ctx.UserEmail,
				TargetId:      roleUUID,
				Action:        models.AdminAuditActionDelete,
				EntityType:    models.AdminAuditEntityRole,
				OrgId:         orgUUID,
				ProjectId:     nil,
				IsAdminZone:   true,
				OldValues:     map[string]interface{}{"name": oldRole.Name, "display_name": oldRole.DisplayName},
				ChangedFields: []string{"*"},
				IpAddress:     ctx.IpAddress,
				UserAgent:     ctx.UserAgent,
			})
		}()
	}

	return nil
}

// AssignRoleToClientWithContext atribui um cargo a um cliente e registra auditoria
func (h *RoleHandler) AssignRoleToClientWithContext(ctx *RequestContext, clientRole *models.ClientRole) error {
	// Buscar informações do cargo para o log
	role, _ := h.GetRole(clientRole.RoleId.String())
	roleName := ""
	if role != nil {
		roleName = role.DisplayName
	}

	// Executar atribuição normal
	if err := h.AssignRoleToClient(clientRole, ctx.UserId.String(), "admin"); err != nil {
		return err
	}

	// Registrar auditoria
	if h.adminAuditHandler != nil {
		go func() {
			if err := h.adminAuditHandler.LogRoleAssignment(
				ctx.UserId, ctx.UserEmail,
				clientRole.ClientId, "",
				clientRole.RoleId, roleName,
				&clientRole.OrganizationId, clientRole.ProjectId,
				ctx.IpAddress, ctx.UserAgent,
			); err != nil {
				fmt.Printf("⚠️ Erro ao registrar log de auditoria (ASSIGN_ROLE): %v\n", err)
			}
		}()
	}

	return nil
}

// RemoveRoleFromClientWithContext remove um cargo de um cliente e registra auditoria
func (h *RoleHandler) RemoveRoleFromClientWithContext(ctx *RequestContext, clientId, roleId, orgId string) error {
	// Buscar informações do cargo para o log
	role, _ := h.GetRole(roleId)
	roleName := ""
	if role != nil {
		roleName = role.DisplayName
	}

	// Executar remoção normal
	if err := h.RemoveRoleFromClient(clientId, roleId, orgId, ctx.UserId.String(), "admin"); err != nil {
		return err
	}

	// Registrar auditoria
	if h.adminAuditHandler != nil {
		go func() {
			clientUUID, _ := uuid.Parse(clientId)
			roleUUID, _ := uuid.Parse(roleId)
			var orgUUID *uuid.UUID
			if orgId != "" {
				parsed, _ := uuid.Parse(orgId)
				orgUUID = &parsed
			}
			if err := h.adminAuditHandler.LogRoleRemoval(
				ctx.UserId, ctx.UserEmail,
				clientUUID, "",
				roleUUID, roleName,
				orgUUID, nil,
				ctx.IpAddress, ctx.UserAgent,
			); err != nil {
				fmt.Printf("⚠️ Erro ao registrar log de auditoria (REMOVE_ROLE): %v\n", err)
			}
		}()
	}

	return nil
}

// CreatePlanWithContext cria um plano e registra auditoria
func (h *RoleHandler) CreatePlanWithContext(ctx *RequestContext, plan *models.Plan) error {
	// Executar criação normal
	if err := h.CreatePlan(plan); err != nil {
		return err
	}

	// Registrar auditoria
	if h.adminAuditHandler != nil {
		go func() {
			h.adminAuditHandler.LogGenericAction(AuditLogParams{
				ActorId:       ctx.UserId,
				ActorEmail:    ctx.UserEmail,
				TargetId:      plan.Id,
				Action:        models.AdminAuditActionCreate,
				EntityType:    models.AdminAuditEntityPackage,
				IsAdminZone:   true,
				NewValues:     map[string]interface{}{"code": plan.Code, "name": plan.Name},
				ChangedFields: []string{"*"},
				IpAddress:     ctx.IpAddress,
				UserAgent:     ctx.UserAgent,
			})
		}()
	}

	return nil
}

// UpdatePlanWithContext atualiza um plano e registra auditoria
func (h *RoleHandler) UpdatePlanWithContext(ctx *RequestContext, plan *models.Plan) error {
	// Executar atualização normal
	if err := h.UpdatePlan(plan); err != nil {
		return err
	}

	// Registrar auditoria
	if h.adminAuditHandler != nil {
		go func() {
			h.adminAuditHandler.LogGenericAction(AuditLogParams{
				ActorId:       ctx.UserId,
				ActorEmail:    ctx.UserEmail,
				TargetId:      plan.Id,
				Action:        models.AdminAuditActionUpdate,
				EntityType:    models.AdminAuditEntityPackage,
				IsAdminZone:   true,
				NewValues:     map[string]interface{}{"code": plan.Code, "name": plan.Name},
				ChangedFields: []string{"code", "name", "prices"},
				IpAddress:     ctx.IpAddress,
				UserAgent:     ctx.UserAgent,
			})
		}()
	}

	return nil
}

// DeletePlanWithContext remove um plano e registra auditoria
func (h *RoleHandler) DeletePlanWithContext(ctx *RequestContext, planId string) error {
	// Capturar estado anterior para auditoria
	oldPlan, _ := h.GetPlanWithModules(planId)

	// Executar exclusão normal
	if err := h.DeletePlan(planId); err != nil {
		return err
	}

	// Registrar auditoria
	if h.adminAuditHandler != nil {
		go func() {
			planUUID, _ := uuid.Parse(planId)
			var oldValues map[string]interface{}
			if oldPlan != nil {
				oldValues = map[string]interface{}{"code": oldPlan.Code, "name": oldPlan.Name}
			}
			h.adminAuditHandler.LogGenericAction(AuditLogParams{
				ActorId:       ctx.UserId,
				ActorEmail:    ctx.UserEmail,
				TargetId:      planUUID,
				Action:        models.AdminAuditActionDelete,
				EntityType:    models.AdminAuditEntityPackage,
				IsAdminZone:   true,
				OldValues:     oldValues,
				ChangedFields: []string{"*"},
				IpAddress:     ctx.IpAddress,
				UserAgent:     ctx.UserAgent,
			})
		}()
	}

	return nil
}

