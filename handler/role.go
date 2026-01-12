package handler

import (
	"fmt"
	"lep/repositories"
	"lep/repositories/models"

	"github.com/google/uuid"
)

type RoleHandler struct {
	roleRepo       repositories.IRoleRepository
	permissionRepo repositories.IPermissionRepository
	moduleRepo     repositories.IModuleRepository
	packageRepo    repositories.IPackageRepository
}

func NewRoleHandler(
	roleRepo repositories.IRoleRepository,
	permissionRepo repositories.IPermissionRepository,
	moduleRepo repositories.IModuleRepository,
	packageRepo repositories.IPackageRepository,
) *RoleHandler {
	return &RoleHandler{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		moduleRepo:     moduleRepo,
		packageRepo:    packageRepo,
	}
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
	// Validar se o ator pode atribuir este cargo
	canManage, err := h.roleRepo.CanManageRole(actorUserId, userRole.RoleId.String(), userRole.OrganizationId.String())
	if err != nil {
		return fmt.Errorf("erro ao verificar permissão: %w", err)
	}

	if !canManage {
		return fmt.Errorf("você não tem permissão para atribuir este cargo")
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
