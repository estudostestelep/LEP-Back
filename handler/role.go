package handler

import (
	"fmt"
	"lep/repositories"
	"lep/repositories/models"

	"github.com/google/uuid"
)

type RoleHandler struct {
	roleRepo          repositories.IRoleRepository
	permissionRepo    repositories.IPermissionRepository
	moduleRepo        repositories.IModuleRepository
	packageRepo       repositories.IPackageRepository
	userRepo          repositories.IUserRepository
	adminAuditHandler IAdminAuditLogHandler
}

func NewRoleHandler(
	roleRepo repositories.IRoleRepository,
	permissionRepo repositories.IPermissionRepository,
	moduleRepo repositories.IModuleRepository,
	packageRepo repositories.IPackageRepository,
	userRepo repositories.IUserRepository,
) *RoleHandler {
	return &RoleHandler{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		moduleRepo:     moduleRepo,
		packageRepo:    packageRepo,
		userRepo:       userRepo,
	}
}

// SetAdminAuditHandler configura o handler de auditoria (injetado separadamente)
func (h *RoleHandler) SetAdminAuditHandler(handler IAdminAuditLogHandler) {
	h.adminAuditHandler = handler
}

// ==================== Role CRUD ====================

// CreateRole cria um novo cargo
func (h *RoleHandler) CreateRole(role *models.Role, actorUserId, orgId string) error {
	// Validar hierarquia - usuário só pode criar cargos com nível igual ou menor
	actorLevel, err := h.roleRepo.GetUserMaxHierarchyLevel(actorUserId, orgId)
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
func (h *RoleHandler) UpdateRole(role *models.Role, actorUserId, orgId string) error {
	// Validar se o usuário pode gerenciar este cargo
	canManage, err := h.roleRepo.CanManageRole(actorUserId, role.Id.String(), orgId)
	if err != nil {
		return fmt.Errorf("erro ao verificar permissão: %w", err)
	}

	if !canManage {
		return fmt.Errorf("você não tem permissão para modificar este cargo")
	}

	return h.roleRepo.Update(role)
}

// DeleteRole remove um cargo
func (h *RoleHandler) DeleteRole(id, actorUserId, orgId string) error {
	// Validar se o usuário pode gerenciar este cargo
	canManage, err := h.roleRepo.CanManageRole(actorUserId, id, orgId)
	if err != nil {
		return fmt.Errorf("erro ao verificar permissão: %w", err)
	}

	if !canManage {
		return fmt.Errorf("você não tem permissão para excluir este cargo")
	}

	// Verificar se é um cargo do sistema
	role, err := h.roleRepo.GetById(id)
	if err != nil {
		return err
	}

	if role.IsSystem {
		return fmt.Errorf("cargos do sistema não podem ser excluídos")
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

// ==================== User-Role Assignment ====================

// AssignRoleToUser atribui um cargo a um usuário
func (h *RoleHandler) AssignRoleToUser(userRole *models.UserRole, actorUserId string) error {
	// Obter orgId - pode ser vazio para cargos admin globais
	orgId := ""
	if userRole.OrganizationId != nil {
		orgId = userRole.OrganizationId.String()
	}

	// Validar se o ator pode atribuir este cargo
	canManage, err := h.roleRepo.CanManageRole(actorUserId, userRole.RoleId.String(), orgId)
	if err != nil {
		return fmt.Errorf("erro ao verificar permissão: %w", err)
	}

	if !canManage {
		return fmt.Errorf("você não tem permissão para atribuir este cargo")
	}

	// Verificar se o cargo é super_admin para sincronizar permissão master_admin
	role, err := h.roleRepo.GetById(userRole.RoleId.String())
	if err != nil {
		return fmt.Errorf("erro ao buscar cargo: %w", err)
	}

	// Se o cargo é super_admin, adicionar permissão master_admin ao usuário
	if role.Name == "super_admin" {
		user, err := h.userRepo.GetUserById(userRole.UserId.String())
		if err != nil {
			return fmt.Errorf("erro ao buscar usuário: %w", err)
		}

		// Verificar se já tem a permissão
		hasMasterAdmin := false
		for _, perm := range user.Permissions {
			if perm == "master_admin" {
				hasMasterAdmin = true
				break
			}
		}

		// Adicionar permissão se não tiver
		if !hasMasterAdmin {
			user.Permissions = append(user.Permissions, "master_admin")
			if err := h.userRepo.UpdateUser(user); err != nil {
				return fmt.Errorf("erro ao atualizar permissões do usuário: %w", err)
			}
		}
	}

	return h.roleRepo.AssignRoleToUser(userRole)
}

// RemoveRoleFromUser remove um cargo de um usuário
func (h *RoleHandler) RemoveRoleFromUser(userId, roleId, orgId, actorUserId string) error {
	// Validar se o ator pode gerenciar este cargo
	canManage, err := h.roleRepo.CanManageRole(actorUserId, roleId, orgId)
	if err != nil {
		return fmt.Errorf("erro ao verificar permissão: %w", err)
	}

	if !canManage {
		return fmt.Errorf("você não tem permissão para remover este cargo")
	}

	// Verificar se o cargo é super_admin para remover permissão master_admin
	role, err := h.roleRepo.GetById(roleId)
	if err != nil {
		return fmt.Errorf("erro ao buscar cargo: %w", err)
	}

	// Se o cargo é super_admin, remover permissão master_admin do usuário
	if role.Name == "super_admin" {
		user, err := h.userRepo.GetUserById(userId)
		if err != nil {
			return fmt.Errorf("erro ao buscar usuário: %w", err)
		}

		// Filtrar a permissão master_admin
		newPermissions := make([]string, 0)
		for _, perm := range user.Permissions {
			if perm != "master_admin" {
				newPermissions = append(newPermissions, perm)
			}
		}

		// Atualizar se houve mudança
		if len(newPermissions) != len(user.Permissions) {
			user.Permissions = newPermissions
			if err := h.userRepo.UpdateUser(user); err != nil {
				return fmt.Errorf("erro ao atualizar permissões do usuário: %w", err)
			}
		}
	}

	return h.roleRepo.RemoveRoleFromUser(userId, roleId, orgId)
}

// GetUserRoles retorna todos os cargos de um usuário
func (h *RoleHandler) GetUserRoles(userId, orgId string) ([]models.UserRole, error) {
	return h.roleRepo.GetUserRoles(userId, orgId)
}

// GetUserRolesWithDetails retorna cargos com detalhes de permissões
func (h *RoleHandler) GetUserRolesWithDetails(userId, orgId string) ([]models.RoleWithPermissionLevels, error) {
	return h.roleRepo.GetUserRolesWithDetails(userId, orgId)
}

// ==================== Permission Level Management ====================

// SetRolePermissionLevel define o nível de uma permissão para um cargo
func (h *RoleHandler) SetRolePermissionLevel(roleId, permissionId string, level int, actorUserId, orgId string) error {
	// Validar nível (0, 1, 2)
	if level < 0 || level > 2 {
		return fmt.Errorf("nível de permissão inválido: deve ser 0, 1 ou 2")
	}

	// Validar se o ator pode gerenciar este cargo
	canManage, err := h.roleRepo.CanManageRole(actorUserId, roleId, orgId)
	if err != nil {
		return fmt.Errorf("erro ao verificar permissão: %w", err)
	}

	if !canManage {
		return fmt.Errorf("você não tem permissão para modificar este cargo")
	}

	return h.roleRepo.SetPermissionLevel(roleId, permissionId, level)
}

// GetRolePermissionLevels retorna todos os níveis de permissão de um cargo
func (h *RoleHandler) GetRolePermissionLevels(roleId string) ([]models.RolePermissionLevel, error) {
	return h.roleRepo.GetPermissionLevels(roleId)
}

// ==================== Permission Checking (Herança) ====================

// UserEffectivePermissionLevel retorna o nível efetivo de uma permissão para um usuário
// considerando todos os seus cargos (retorna o maior nível)
func (h *RoleHandler) UserEffectivePermissionLevel(userId, orgId, permissionCodeName string) (int, error) {
	// Buscar a permissão pelo código
	permission, err := h.permissionRepo.GetByCodeName(permissionCodeName)
	if err != nil {
		return 0, fmt.Errorf("permissão não encontrada: %s", permissionCodeName)
	}

	// Buscar todos os cargos do usuário com detalhes
	rolesWithLevels, err := h.roleRepo.GetUserRolesWithDetails(userId, orgId)
	if err != nil {
		return 0, err
	}

	// Encontrar o maior nível para esta permissão
	maxLevel := 0
	for _, roleData := range rolesWithLevels {
		for _, permLevel := range roleData.PermissionLevels {
			if permLevel.PermissionId == permission.Id && permLevel.Level > maxLevel {
				maxLevel = permLevel.Level
			}
		}
	}

	return maxLevel, nil
}

// HasPermission verifica se um usuário tem pelo menos o nível mínimo para uma permissão
func (h *RoleHandler) HasPermission(userId, orgId, permissionCodeName string, minLevel int) (bool, error) {
	effectiveLevel, err := h.UserEffectivePermissionLevel(userId, orgId, permissionCodeName)
	if err != nil {
		return false, err
	}
	return effectiveLevel >= minLevel, nil
}

// HasModuleAccess verifica se a organização tem acesso a um módulo
func (h *RoleHandler) HasModuleAccess(orgId, moduleCodeName string) (bool, error) {
	// Buscar o módulo
	module, err := h.moduleRepo.GetByCodeName(moduleCodeName)
	if err != nil {
		return false, fmt.Errorf("módulo não encontrado: %s", moduleCodeName)
	}

	// Se o módulo é gratuito, todos têm acesso
	if module.IsFree {
		return true, nil
	}

	// Buscar o pacote da organização
	orgPackage, err := h.packageRepo.GetOrganizationPackage(orgId)
	if err != nil {
		// Sem pacote, só módulos gratuitos
		return false, nil
	}

	// Verificar se o módulo está no pacote
	return h.moduleRepo.IsModuleInPackage(module.Id.String(), orgPackage.PackageId.String())
}

// FullPermissionCheck faz verificação completa: módulo + permissão + hierarquia
func (h *RoleHandler) FullPermissionCheck(userId, orgId, permissionCodeName string, minLevel int, targetUserId string) (bool, string, error) {
	// 1. Buscar a permissão para identificar o módulo
	permission, err := h.permissionRepo.GetByCodeName(permissionCodeName)
	if err != nil {
		return false, "Permissão não encontrada", err
	}

	// 2. Verificar acesso ao módulo
	if permission.Module != nil {
		hasModule, err := h.HasModuleAccess(orgId, permission.Module.CodeName)
		if err != nil {
			return false, "Erro ao verificar módulo", err
		}
		if !hasModule {
			return false, fmt.Sprintf("Módulo '%s' não disponível no seu plano", permission.Module.DisplayName), nil
		}
	}

	// 3. Verificar nível de permissão
	hasPermission, err := h.HasPermission(userId, orgId, permissionCodeName, minLevel)
	if err != nil {
		return false, "Erro ao verificar permissão", err
	}
	if !hasPermission {
		return false, fmt.Sprintf("Você não tem permissão para '%s'", permission.DisplayName), nil
	}

	// 4. Se há um usuário alvo, verificar hierarquia
	if targetUserId != "" && targetUserId != userId {
		actorLevel, err := h.roleRepo.GetUserMaxHierarchyLevel(userId, orgId)
		if err != nil {
			return false, "Erro ao verificar hierarquia", err
		}

		targetLevel, err := h.roleRepo.GetUserMaxHierarchyLevel(targetUserId, orgId)
		if err != nil {
			return false, "Erro ao verificar hierarquia do alvo", err
		}

		if targetLevel > actorLevel {
			return false, "Você não pode gerenciar usuários com nível de hierarquia maior", nil
		}
	}

	return true, "", nil
}

// ==================== Helpers ====================

// GetUserMaxHierarchyLevel retorna o maior nível de hierarquia do usuário
func (h *RoleHandler) GetUserMaxHierarchyLevel(userId, orgId string) (int, error) {
	return h.roleRepo.GetUserMaxHierarchyLevel(userId, orgId)
}

// CanManageUser verifica se um usuário pode gerenciar outro
func (h *RoleHandler) CanManageUser(actorUserId, targetUserId, orgId string) (bool, error) {
	actorLevel, err := h.roleRepo.GetUserMaxHierarchyLevel(actorUserId, orgId)
	if err != nil {
		return false, err
	}

	targetLevel, err := h.roleRepo.GetUserMaxHierarchyLevel(targetUserId, orgId)
	if err != nil {
		return false, err
	}

	return actorLevel >= targetLevel, nil
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

// ==================== Package Management ====================

// ListPackages lista todos os pacotes
func (h *RoleHandler) ListPackages(publicOnly bool) ([]models.Package, error) {
	if publicOnly {
		return h.packageRepo.ListPublic()
	}
	return h.packageRepo.List()
}

// GetPackageWithModules retorna um pacote com seus módulos
func (h *RoleHandler) GetPackageWithModules(id string) (*models.Package, error) {
	return h.packageRepo.GetPackageWithModules(id)
}

// SubscribeOrganization inscreve uma organização em um pacote
func (h *RoleHandler) SubscribeOrganization(orgId, packageId string, billingCycle string, customPrice *float64) error {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return fmt.Errorf("ID da organização inválido: %w", err)
	}

	pkgUUID, err := uuid.Parse(packageId)
	if err != nil {
		return fmt.Errorf("ID do pacote inválido: %w", err)
	}

	// Buscar pacote para obter preço
	pkg, err := h.packageRepo.GetById(packageId)
	if err != nil {
		return fmt.Errorf("pacote não encontrado: %w", err)
	}

	price := pkg.PriceMonthly
	if billingCycle == "yearly" {
		price = pkg.PriceYearly
	}
	if customPrice != nil {
		price = *customPrice
	}

	orgPackage := &models.OrganizationPackage{
		OrganizationId: orgUUID,
		PackageId:      pkgUUID,
		BillingCycle:   billingCycle,
		CustomPrice:    customPrice,
		Active:         true,
	}

	// Se o preço customizado não foi definido, usar o preço padrão
	if customPrice == nil {
		orgPackage.CustomPrice = &price
	}

	return h.packageRepo.SubscribeOrganization(orgPackage)
}

// GetOrganizationSubscription retorna a assinatura ativa da organização
func (h *RoleHandler) GetOrganizationSubscription(orgId string) (*models.OrganizationPackage, error) {
	return h.packageRepo.GetOrganizationPackage(orgId)
}

// ==================== Package CRUD (Master Admin) ====================

// CreatePackage cria um novo pacote
func (h *RoleHandler) CreatePackage(pkg *models.Package) error {
	return h.packageRepo.Create(pkg)
}

// UpdatePackage atualiza um pacote
func (h *RoleHandler) UpdatePackage(pkg *models.Package) error {
	return h.packageRepo.Update(pkg)
}

// DeletePackage remove um pacote
func (h *RoleHandler) DeletePackage(id string) error {
	return h.packageRepo.Delete(id)
}

// AddModuleToPackage adiciona um módulo a um pacote
func (h *RoleHandler) AddModuleToPackage(packageId, moduleId string) error {
	return h.packageRepo.AddModuleToPackage(packageId, moduleId)
}

// RemoveModuleFromPackage remove um módulo de um pacote
func (h *RoleHandler) RemoveModuleFromPackage(packageId, moduleId string) error {
	return h.packageRepo.RemoveModuleFromPackage(packageId, moduleId)
}

// SetPackageLimit define um limite para o pacote
func (h *RoleHandler) SetPackageLimit(packageId, limitType string, limitValue int) error {
	return h.packageRepo.SetPackageLimit(packageId, limitType, limitValue)
}

// GetPackageLimits retorna os limites de um pacote
func (h *RoleHandler) GetPackageLimits(packageId string) ([]models.PackageLimit, error) {
	return h.packageRepo.GetPackageLimits(packageId)
}

// UpdateOrganizationSubscription atualiza a assinatura de uma organização
func (h *RoleHandler) UpdateOrganizationSubscription(orgId, packageId, billingCycle string, customPrice *float64, active *bool) error {
	// Buscar assinatura existente
	existingSubscription, err := h.packageRepo.GetOrganizationPackage(orgId)
	if err != nil {
		return fmt.Errorf("assinatura não encontrada: %w", err)
	}

	// Atualizar campos se fornecidos
	if packageId != "" {
		pkgUUID, err := uuid.Parse(packageId)
		if err != nil {
			return fmt.Errorf("ID do pacote inválido: %w", err)
		}
		existingSubscription.PackageId = pkgUUID
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

	return h.packageRepo.UpdateOrganizationPackage(existingSubscription)
}

// CancelOrganizationSubscription cancela a assinatura de uma organização
func (h *RoleHandler) CancelOrganizationSubscription(orgId string) error {
	return h.packageRepo.CancelOrganizationPackage(orgId)
}

// ListAllSubscriptions lista todas as assinaturas ativas
func (h *RoleHandler) ListAllSubscriptions() ([]models.OrganizationPackage, error) {
	return h.packageRepo.ListAllSubscriptions()
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
	if err := h.CreateRole(role, ctx.UserId.String(), orgId); err != nil {
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
				ProjectId:     role.ProjectId,
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
	if err := h.UpdateRole(role, ctx.UserId.String(), orgId); err != nil {
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
				ProjectId:     role.ProjectId,
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
	if err := h.DeleteRole(roleId, ctx.UserId.String(), orgId); err != nil {
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
				ProjectId:     oldRole.ProjectId,
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

// AssignRoleToUserWithContext atribui um cargo e registra auditoria
func (h *RoleHandler) AssignRoleToUserWithContext(ctx *RequestContext, userRole *models.UserRole) error {
	// Buscar informações do cargo para o log
	role, _ := h.GetRole(userRole.RoleId.String())
	roleName := ""
	if role != nil {
		roleName = role.DisplayName
	}

	// Executar atribuição normal
	if err := h.AssignRoleToUser(userRole, ctx.UserId.String()); err != nil {
		return err
	}

	// Registrar auditoria
	if h.adminAuditHandler != nil {
		go func() {
			if err := h.adminAuditHandler.LogRoleAssignment(
				ctx.UserId, ctx.UserEmail,
				userRole.UserId, "",
				userRole.RoleId, roleName,
				userRole.OrganizationId, userRole.ProjectId,
				ctx.IpAddress, ctx.UserAgent,
			); err != nil {
				fmt.Printf("⚠️ Erro ao registrar log de auditoria (ASSIGN_ROLE): %v\n", err)
			}
		}()
	}

	return nil
}

// RemoveRoleFromUserWithContext remove um cargo e registra auditoria
func (h *RoleHandler) RemoveRoleFromUserWithContext(ctx *RequestContext, userId, roleId, orgId string) error {
	// Buscar informações do cargo para o log
	role, _ := h.GetRole(roleId)
	roleName := ""
	if role != nil {
		roleName = role.DisplayName
	}

	// Executar remoção normal
	if err := h.RemoveRoleFromUser(userId, roleId, orgId, ctx.UserId.String()); err != nil {
		return err
	}

	// Registrar auditoria
	if h.adminAuditHandler != nil {
		go func() {
			userUUID, _ := uuid.Parse(userId)
			roleUUID, _ := uuid.Parse(roleId)
			var orgUUID *uuid.UUID
			if orgId != "" {
				parsed, _ := uuid.Parse(orgId)
				orgUUID = &parsed
			}
			if err := h.adminAuditHandler.LogRoleRemoval(
				ctx.UserId, ctx.UserEmail,
				userUUID, "",
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

// CreatePackageWithContext cria um pacote e registra auditoria
func (h *RoleHandler) CreatePackageWithContext(ctx *RequestContext, pkg *models.Package) error {
	// Executar criação normal
	if err := h.CreatePackage(pkg); err != nil {
		return err
	}

	// Registrar auditoria
	if h.adminAuditHandler != nil {
		go func() {
			h.adminAuditHandler.LogGenericAction(AuditLogParams{
				ActorId:       ctx.UserId,
				ActorEmail:    ctx.UserEmail,
				TargetId:      pkg.Id,
				Action:        models.AdminAuditActionCreate,
				EntityType:    models.AdminAuditEntityPackage,
				IsAdminZone:   true,
				NewValues:     map[string]interface{}{"code_name": pkg.CodeName, "display_name": pkg.DisplayName},
				ChangedFields: []string{"*"},
				IpAddress:     ctx.IpAddress,
				UserAgent:     ctx.UserAgent,
			})
		}()
	}

	return nil
}

// UpdatePackageWithContext atualiza um pacote e registra auditoria
func (h *RoleHandler) UpdatePackageWithContext(ctx *RequestContext, pkg *models.Package) error {
	// Executar atualização normal
	if err := h.UpdatePackage(pkg); err != nil {
		return err
	}

	// Registrar auditoria
	if h.adminAuditHandler != nil {
		go func() {
			h.adminAuditHandler.LogGenericAction(AuditLogParams{
				ActorId:       ctx.UserId,
				ActorEmail:    ctx.UserEmail,
				TargetId:      pkg.Id,
				Action:        models.AdminAuditActionUpdate,
				EntityType:    models.AdminAuditEntityPackage,
				IsAdminZone:   true,
				NewValues:     map[string]interface{}{"code_name": pkg.CodeName, "display_name": pkg.DisplayName},
				ChangedFields: []string{"code_name", "display_name", "prices"},
				IpAddress:     ctx.IpAddress,
				UserAgent:     ctx.UserAgent,
			})
		}()
	}

	return nil
}

// DeletePackageWithContext remove um pacote e registra auditoria
func (h *RoleHandler) DeletePackageWithContext(ctx *RequestContext, packageId string) error {
	// Capturar estado anterior para auditoria
	oldPkg, _ := h.GetPackageWithModules(packageId)

	// Executar exclusão normal
	if err := h.DeletePackage(packageId); err != nil {
		return err
	}

	// Registrar auditoria
	if h.adminAuditHandler != nil {
		go func() {
			pkgUUID, _ := uuid.Parse(packageId)
			var oldValues map[string]interface{}
			if oldPkg != nil {
				oldValues = map[string]interface{}{"code_name": oldPkg.CodeName, "display_name": oldPkg.DisplayName}
			}
			h.adminAuditHandler.LogGenericAction(AuditLogParams{
				ActorId:       ctx.UserId,
				ActorEmail:    ctx.UserEmail,
				TargetId:      pkgUUID,
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
