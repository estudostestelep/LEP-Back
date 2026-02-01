package repositories

import (
	"lep/repositories/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type resourceRole struct {
	db *gorm.DB
}

// IRoleRepository interface para operações com Roles
type IRoleRepository interface {
	// CRUD de Roles
	Create(role *models.Role) error
	GetById(id string) (*models.Role, error)
	GetByIdWithPermissions(id string) (*models.Role, error)
	GetByName(name string) (*models.Role, error)
	Update(role *models.Role) error
	Delete(id string) error
	List() ([]models.Role, error)
	ListByScope(scope string) ([]models.Role, error)
	ListByOrganization(orgId string) ([]models.Role, error)
	ListSystemRoles() ([]models.Role, error)

	// Role-Permission (nova pivot simplificada)
	AddPermissionToRole(roleId, permissionId string) error
	RemovePermissionFromRole(roleId, permissionId string) error
	GetRolePermissions(roleId string) ([]models.Permission, error)
	GetRolePermissionCodes(roleId string) ([]string, error)
	SetRolePermissions(roleId string, permissionIds []string) error

	// Admin-Role associations
	AssignRoleToAdmin(adminRole *models.AdminRole) error
	RemoveRoleFromAdmin(adminId, roleId string) error
	GetAdminRoles(adminId string) ([]models.AdminRole, error)
	GetAdminRolesWithPermissions(adminId string) ([]models.RoleWithPermissions, error)
	GetAdminMaxHierarchyLevel(adminId string) (int, error)

	// Client-Role associations
	AssignRoleToClient(clientRole *models.ClientRole) error
	RemoveRoleFromClient(clientId, roleId, orgId string) error
	GetClientRoles(clientId, orgId string) ([]models.ClientRole, error)
	GetClientRolesWithPermissions(clientId, orgId string) ([]models.RoleWithPermissions, error)
	GetClientMaxHierarchyLevel(clientId, orgId string) (int, error)
	DeleteClientRolesByClientId(clientId string) error

	// Permission checks
	UserHasPermission(userId, userType, permission string) (bool, error)
	IsMasterAdmin(userId, userType string) (bool, error)
	GetUserHierarchyLevel(userId, userType, orgId string) (int, error)
	CanManageUser(actorId, targetId, userType, orgId string) (bool, error)

	// Organization counts
	CountClientsByOrganization(orgId string) (int, error)
	ListClientsByOrganization(orgId string) ([]models.ClientRole, error)
}

func NewRoleRepository(db *gorm.DB) IRoleRepository {
	return &resourceRole{db: db}
}

// ==================== CRUD de Roles ====================

// Create cria um novo cargo
func (r *resourceRole) Create(role *models.Role) error {
	if role.Id == uuid.Nil {
		role.Id = uuid.New()
	}
	return r.db.Create(role).Error
}

// GetById busca cargo por ID
func (r *resourceRole) GetById(id string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByIdWithPermissions busca cargo por ID com permissões carregadas
func (r *resourceRole) GetByIdWithPermissions(id string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).
		Preload("Permissions").
		First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByName busca cargo por nome
func (r *resourceRole) GetByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ? AND deleted_at IS NULL", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// Update atualiza cargo
func (r *resourceRole) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

// Delete faz soft delete do cargo
func (r *resourceRole) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Role{}).Error
}

// List lista todos os cargos ativos
func (r *resourceRole) List() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Where("deleted_at IS NULL AND active = true").
		Order("hierarchy_level DESC, name ASC").
		Find(&roles).Error
	return roles, err
}

// ListByScope lista cargos por escopo (admin ou client)
func (r *resourceRole) ListByScope(scope string) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Where("scope = ? AND deleted_at IS NULL AND active = true", scope).
		Order("hierarchy_level DESC, name ASC").
		Find(&roles).Error
	return roles, err
}

// ListByOrganization lista cargos de uma organização + cargos globais
func (r *resourceRole) ListByOrganization(orgId string) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Where("(organization_id = ? OR organization_id IS NULL) AND deleted_at IS NULL AND active = true", orgId).
		Order("hierarchy_level DESC, name ASC").
		Find(&roles).Error
	return roles, err
}

// ListSystemRoles lista apenas cargos do sistema
func (r *resourceRole) ListSystemRoles() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Where("is_system = true AND deleted_at IS NULL AND active = true").
		Order("hierarchy_level DESC, name ASC").
		Find(&roles).Error
	return roles, err
}

// ==================== Role-Permission (nova pivot) ====================

// AddPermissionToRole adiciona uma permissão a um cargo
func (r *resourceRole) AddPermissionToRole(roleId, permissionId string) error {
	roleUUID, err := uuid.Parse(roleId)
	if err != nil {
		return err
	}
	permUUID, err := uuid.Parse(permissionId)
	if err != nil {
		return err
	}

	// Verificar se já existe
	var existing models.RolePermission
	err = r.db.Where("role_id = ? AND permission_id = ?", roleId, permissionId).
		First(&existing).Error
	if err == nil {
		return nil // Já existe
	}

	rp := models.RolePermission{
		RoleId:       roleUUID,
		PermissionId: permUUID,
	}
	return r.db.Create(&rp).Error
}

// RemovePermissionFromRole remove uma permissão de um cargo
func (r *resourceRole) RemovePermissionFromRole(roleId, permissionId string) error {
	return r.db.Where("role_id = ? AND permission_id = ?", roleId, permissionId).
		Delete(&models.RolePermission{}).Error
}

// GetRolePermissions retorna todas as permissões de um cargo
func (r *resourceRole) GetRolePermissions(roleId string) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Table("permissions").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleId).
		Where("permissions.deleted_at IS NULL AND permissions.active = true").
		Order("permissions.module, permissions.action").
		Find(&permissions).Error
	return permissions, err
}

// GetRolePermissionCodes retorna os códigos das permissões de um cargo
func (r *resourceRole) GetRolePermissionCodes(roleId string) ([]string, error) {
	var codes []string
	err := r.db.Table("permissions").
		Select("permissions.code").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleId).
		Where("permissions.deleted_at IS NULL AND permissions.active = true").
		Pluck("code", &codes).Error
	return codes, err
}

// SetRolePermissions define todas as permissões de um cargo (substitui existentes)
func (r *resourceRole) SetRolePermissions(roleId string, permissionIds []string) error {
	// Remover permissões existentes
	if err := r.db.Where("role_id = ?", roleId).Delete(&models.RolePermission{}).Error; err != nil {
		return err
	}

	// Adicionar novas permissões
	roleUUID, err := uuid.Parse(roleId)
	if err != nil {
		return err
	}

	for _, permId := range permissionIds {
		permUUID, err := uuid.Parse(permId)
		if err != nil {
			continue
		}
		rp := models.RolePermission{
			RoleId:       roleUUID,
			PermissionId: permUUID,
		}
		if err := r.db.Create(&rp).Error; err != nil {
			return err
		}
	}

	return nil
}

// ==================== Admin-Role ====================

// AssignRoleToAdmin associa um cargo a um admin
func (r *resourceRole) AssignRoleToAdmin(adminRole *models.AdminRole) error {
	if adminRole.Id == uuid.Nil {
		adminRole.Id = uuid.New()
	}

	// Verificar se já existe associação ativa
	var existing models.AdminRole
	query := r.db.Where("admin_id = ? AND role_id = ? AND deleted_at IS NULL",
		adminRole.AdminId, adminRole.RoleId)

	if adminRole.OrganizationId != nil {
		query = query.Where("organization_id = ?", adminRole.OrganizationId)
	} else {
		query = query.Where("organization_id IS NULL")
	}

	err := query.First(&existing).Error
	if err == nil {
		if !existing.Active {
			existing.Active = true
			return r.db.Save(&existing).Error
		}
		return nil
	}

	return r.db.Create(adminRole).Error
}

// RemoveRoleFromAdmin remove associação de cargo do admin
func (r *resourceRole) RemoveRoleFromAdmin(adminId, roleId string) error {
	return r.db.Where("admin_id = ? AND role_id = ?", adminId, roleId).
		Delete(&models.AdminRole{}).Error
}

// GetAdminRoles busca todos os cargos de um admin
func (r *resourceRole) GetAdminRoles(adminId string) ([]models.AdminRole, error) {
	var adminRoles []models.AdminRole
	err := r.db.Where("admin_id = ? AND deleted_at IS NULL AND active = true", adminId).
		Preload("Role").
		Find(&adminRoles).Error
	return adminRoles, err
}

// GetAdminRolesWithPermissions busca cargos com permissões
func (r *resourceRole) GetAdminRolesWithPermissions(adminId string) ([]models.RoleWithPermissions, error) {
	adminRoles, err := r.GetAdminRoles(adminId)
	if err != nil {
		return nil, err
	}

	var result []models.RoleWithPermissions
	for _, ar := range adminRoles {
		if ar.Role == nil {
			continue
		}

		codes, err := r.GetRolePermissionCodes(ar.RoleId.String())
		if err != nil {
			codes = []string{}
		}

		result = append(result, models.RoleWithPermissions{
			Role:            *ar.Role,
			PermissionCodes: codes,
		})
	}

	return result, nil
}

// GetAdminMaxHierarchyLevel retorna o maior nível de hierarquia do admin
func (r *resourceRole) GetAdminMaxHierarchyLevel(adminId string) (int, error) {
	var maxLevel int
	err := r.db.Model(&models.AdminRole{}).
		Select("COALESCE(MAX(roles.hierarchy_level), 0)").
		Joins("JOIN roles ON roles.id = admin_roles.role_id").
		Where("admin_roles.admin_id = ? AND admin_roles.deleted_at IS NULL AND admin_roles.active = true", adminId).
		Where("roles.deleted_at IS NULL AND roles.active = true").
		Scan(&maxLevel).Error

	return maxLevel, err
}

// ==================== Client-Role ====================

// AssignRoleToClient associa um cargo a um client
func (r *resourceRole) AssignRoleToClient(clientRole *models.ClientRole) error {
	if clientRole.Id == uuid.Nil {
		clientRole.Id = uuid.New()
	}

	// Verificar se já existe associação ativa
	var existing models.ClientRole
	query := r.db.Where("client_id = ? AND role_id = ? AND organization_id = ? AND deleted_at IS NULL",
		clientRole.ClientId, clientRole.RoleId, clientRole.OrganizationId)

	if clientRole.ProjectId != nil {
		query = query.Where("project_id = ?", clientRole.ProjectId)
	} else {
		query = query.Where("project_id IS NULL")
	}

	err := query.First(&existing).Error
	if err == nil {
		if !existing.Active {
			existing.Active = true
			return r.db.Save(&existing).Error
		}
		return nil
	}

	return r.db.Create(clientRole).Error
}

// RemoveRoleFromClient remove associação de cargo do client
func (r *resourceRole) RemoveRoleFromClient(clientId, roleId, orgId string) error {
	return r.db.Where("client_id = ? AND role_id = ? AND organization_id = ?", clientId, roleId, orgId).
		Delete(&models.ClientRole{}).Error
}

// GetClientRoles busca todos os cargos de um client
func (r *resourceRole) GetClientRoles(clientId, orgId string) ([]models.ClientRole, error) {
	var clientRoles []models.ClientRole
	query := r.db.Where("client_id = ? AND deleted_at IS NULL AND active = true", clientId)

	if orgId != "" {
		query = query.Where("organization_id = ?", orgId)
	}

	err := query.Preload("Role").Find(&clientRoles).Error
	return clientRoles, err
}

// GetClientRolesWithPermissions busca cargos com permissões
func (r *resourceRole) GetClientRolesWithPermissions(clientId, orgId string) ([]models.RoleWithPermissions, error) {
	clientRoles, err := r.GetClientRoles(clientId, orgId)
	if err != nil {
		return nil, err
	}

	var result []models.RoleWithPermissions
	for _, cr := range clientRoles {
		if cr.Role == nil {
			continue
		}

		codes, err := r.GetRolePermissionCodes(cr.RoleId.String())
		if err != nil {
			codes = []string{}
		}

		result = append(result, models.RoleWithPermissions{
			Role:            *cr.Role,
			PermissionCodes: codes,
		})
	}

	return result, nil
}

// GetClientMaxHierarchyLevel retorna o maior nível de hierarquia do client
func (r *resourceRole) GetClientMaxHierarchyLevel(clientId, orgId string) (int, error) {
	var maxLevel int
	query := r.db.Model(&models.ClientRole{}).
		Select("COALESCE(MAX(roles.hierarchy_level), 0)").
		Joins("JOIN roles ON roles.id = client_roles.role_id").
		Where("client_roles.client_id = ? AND client_roles.deleted_at IS NULL AND client_roles.active = true", clientId).
		Where("roles.deleted_at IS NULL AND roles.active = true")

	if orgId != "" {
		query = query.Where("client_roles.organization_id = ?", orgId)
	}

	err := query.Scan(&maxLevel).Error
	return maxLevel, err
}

// DeleteClientRolesByClientId remove todos os cargos de um client
func (r *resourceRole) DeleteClientRolesByClientId(clientId string) error {
	return r.db.Where("client_id = ?", clientId).Delete(&models.ClientRole{}).Error
}

// ==================== Permission Checks ====================

// UserHasPermission verifica se usuário tem uma permissão específica
func (r *resourceRole) UserHasPermission(userId, userType, permission string) (bool, error) {
	// Master admin tem todas as permissões
	isMaster, err := r.IsMasterAdmin(userId, userType)
	if err == nil && isMaster {
		return true, nil
	}

	var count int64

	if userType == "admin" {
		// Buscar permissões via admin_roles
		err = r.db.Table("permissions").
			Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
			Joins("INNER JOIN admin_roles ON role_permissions.role_id = admin_roles.role_id").
			Where("admin_roles.admin_id = ?", userId).
			Where("admin_roles.deleted_at IS NULL AND admin_roles.active = true").
			Where("permissions.code = ?", permission).
			Where("permissions.deleted_at IS NULL AND permissions.active = true").
			Count(&count).Error
	} else {
		// Buscar permissões via client_roles
		err = r.db.Table("permissions").
			Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
			Joins("INNER JOIN client_roles ON role_permissions.role_id = client_roles.role_id").
			Where("client_roles.client_id = ?", userId).
			Where("client_roles.deleted_at IS NULL AND client_roles.active = true").
			Where("permissions.code = ?", permission).
			Where("permissions.deleted_at IS NULL AND permissions.active = true").
			Count(&count).Error
	}

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// IsMasterAdmin verifica se usuário é master admin (hierarchy >= 10)
func (r *resourceRole) IsMasterAdmin(userId, userType string) (bool, error) {
	var maxLevel int

	if userType == "admin" {
		maxLevel, _ = r.GetAdminMaxHierarchyLevel(userId)
	} else {
		maxLevel, _ = r.GetClientMaxHierarchyLevel(userId, "")
	}

	return maxLevel >= 10, nil
}

// GetUserHierarchyLevel retorna o nível de hierarquia do usuário
func (r *resourceRole) GetUserHierarchyLevel(userId, userType, orgId string) (int, error) {
	if userType == "admin" {
		return r.GetAdminMaxHierarchyLevel(userId)
	}
	return r.GetClientMaxHierarchyLevel(userId, orgId)
}

// CanManageUser verifica se actor pode gerenciar target
func (r *resourceRole) CanManageUser(actorId, targetId, userType, orgId string) (bool, error) {
	actorLevel, err := r.GetUserHierarchyLevel(actorId, userType, orgId)
	if err != nil {
		return false, err
	}

	targetLevel, err := r.GetUserHierarchyLevel(targetId, userType, orgId)
	if err != nil {
		return false, err
	}

	return actorLevel >= targetLevel, nil
}

// ==================== Organization ====================

// CountClientsByOrganization conta clientes de uma organização
func (r *resourceRole) CountClientsByOrganization(orgId string) (int, error) {
	var count int64
	err := r.db.Model(&models.ClientRole{}).
		Select("COUNT(DISTINCT client_id)").
		Where("organization_id = ? AND deleted_at IS NULL AND active = true", orgId).
		Scan(&count).Error
	return int(count), err
}

// ListClientsByOrganization lista client_roles de uma organização
func (r *resourceRole) ListClientsByOrganization(orgId string) ([]models.ClientRole, error) {
	var clientRoles []models.ClientRole
	err := r.db.Where("organization_id = ? AND deleted_at IS NULL AND active = true", orgId).
		Preload("Client").
		Preload("Role").
		Find(&clientRoles).Error
	return clientRoles, err
}
