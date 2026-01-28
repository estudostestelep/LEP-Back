package repositories

import (
	"lep/repositories/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type resourceRole struct {
	db *gorm.DB
}

type IRoleRepository interface {
	// CRUD de Roles
	Create(role *models.Role) error
	GetById(id string) (*models.Role, error)
	GetByName(name string) (*models.Role, error)
	Update(role *models.Role) error
	Delete(id string) error
	List() ([]models.Role, error)
	ListByScope(scope string) ([]models.Role, error)
	ListByOrganization(orgId string) ([]models.Role, error)
	ListSystemRoles() ([]models.Role, error)

	// User-Role associations
	AssignRoleToUser(userRole *models.UserRole) error
	RemoveRoleFromUser(userId, roleId, orgId string) error
	GetUserRoles(userId, orgId, scope string) ([]models.UserRole, error)
	GetUserRolesWithDetails(userId, orgId, scope string) ([]models.RoleWithPermissionLevels, error)
	GetAllUserRoles(userId string) ([]models.UserRole, error)
	CountUsersByOrganization(orgId string) (int, error)
	ListUsersByOrganization(orgId string) ([]models.UserRole, error)
	DeleteUserRolesByUserId(userId string) error

	// Permission Level management
	SetPermissionLevel(roleId, permissionId string, level int) error
	GetPermissionLevels(roleId string) ([]models.RolePermissionLevel, error)

	// Hierarchy validation
	GetUserMaxHierarchyLevel(userId, orgId string) (int, error)
	CanManageRole(actorUserId, targetRoleId, orgId string) (bool, error)
}

func NewRoleRepository(db *gorm.DB) IRoleRepository {
	return &resourceRole{db: db}
}

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

// ListByOrganization lista cargos de uma organização específica + cargos globais
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

// AssignRoleToUser associa um cargo a um usuário
func (r *resourceRole) AssignRoleToUser(userRole *models.UserRole) error {
	if userRole.Id == uuid.Nil {
		userRole.Id = uuid.New()
	}

	// Verificar se já existe associação ativa
	var existing models.UserRole
	query := r.db.Where("user_id = ? AND role_id = ? AND deleted_at IS NULL",
		userRole.UserId, userRole.RoleId)

	// Tratar OrganizationId null (cargo admin global)
	if userRole.OrganizationId != nil {
		query = query.Where("organization_id = ?", userRole.OrganizationId)
	} else {
		query = query.Where("organization_id IS NULL")
	}

	if userRole.ProjectId != nil {
		query = query.Where("project_id = ?", userRole.ProjectId)
	} else {
		query = query.Where("project_id IS NULL")
	}

	err := query.First(&existing).Error
	if err == nil {
		// Já existe, apenas ativar se estiver inativo
		if !existing.Active {
			existing.Active = true
			return r.db.Save(&existing).Error
		}
		return nil // Já existe e está ativo
	}

	return r.db.Create(userRole).Error
}

// RemoveRoleFromUser remove associação de cargo do usuário
func (r *resourceRole) RemoveRoleFromUser(userId, roleId, orgId string) error {
	query := r.db.Where("user_id = ? AND role_id = ?", userId, roleId)

	// Tratar orgId vazio (cargo admin global)
	if orgId != "" {
		query = query.Where("organization_id = ?", orgId)
	} else {
		query = query.Where("organization_id IS NULL")
	}

	return query.Delete(&models.UserRole{}).Error
}

func (r *resourceRole) DeleteUserRolesByUserId(userId string) error {
	return r.db.Where("user_id = ?", userId).Delete(&models.UserRole{}).Error
}

// GetUserRoles busca todos os cargos de um usuário
// Se orgId for fornecido, filtra por organização específica
// Se orgId for vazio, retorna todos os cargos do usuário (todas as organizações)
// Se scope for fornecido, filtra apenas roles daquele escopo (admin/client)
func (r *resourceRole) GetUserRoles(userId, orgId, scope string) ([]models.UserRole, error) {
	var userRoles []models.UserRole
	query := r.db.Where("user_roles.user_id = ? AND user_roles.deleted_at IS NULL AND user_roles.active = true", userId)

	// Filtrar por organização apenas se fornecido
	if orgId != "" {
		query = query.Where("user_roles.organization_id = ?", orgId)
	}

	// Filtrar por scope se fornecido
	if scope != "" {
		query = query.Joins("JOIN roles ON roles.id = user_roles.role_id").
			Where("roles.scope = ?", scope)
	}

	err := query.Preload("Role").Find(&userRoles).Error
	return userRoles, err
}

// GetUserRolesWithDetails busca cargos com detalhes de permissões
// Se scope for fornecido, filtra apenas roles daquele escopo (admin/client)
func (r *resourceRole) GetUserRolesWithDetails(userId, orgId, scope string) ([]models.RoleWithPermissionLevels, error) {
	userRoles, err := r.GetUserRoles(userId, orgId, scope)
	if err != nil {
		return nil, err
	}

	var result []models.RoleWithPermissionLevels
	for _, ur := range userRoles {
		if ur.Role == nil {
			continue
		}

		levels, err := r.GetPermissionLevels(ur.RoleId.String())
		if err != nil {
			continue
		}

		result = append(result, models.RoleWithPermissionLevels{
			Role:             *ur.Role,
			PermissionLevels: levels,
		})
	}

	return result, nil
}

// SetPermissionLevel define o nível de acesso de uma permissão para um cargo
func (r *resourceRole) SetPermissionLevel(roleId, permissionId string, level int) error {
	roleUUID, err := uuid.Parse(roleId)
	if err != nil {
		return err
	}
	permUUID, err := uuid.Parse(permissionId)
	if err != nil {
		return err
	}

	// Buscar registro existente
	var existing models.RolePermissionLevel
	err = r.db.Where("role_id = ? AND permission_id = ? AND deleted_at IS NULL", roleId, permissionId).
		First(&existing).Error

	if err == nil {
		// Atualizar existente
		existing.Level = level
		return r.db.Save(&existing).Error
	}

	// Criar novo
	newLevel := models.RolePermissionLevel{
		Id:           uuid.New(),
		RoleId:       roleUUID,
		PermissionId: permUUID,
		Level:        level,
	}
	return r.db.Create(&newLevel).Error
}

// GetPermissionLevels busca todos os níveis de permissão de um cargo
func (r *resourceRole) GetPermissionLevels(roleId string) ([]models.RolePermissionLevel, error) {
	var levels []models.RolePermissionLevel
	err := r.db.Where("role_id = ? AND deleted_at IS NULL", roleId).
		Preload("Permission").
		Preload("Permission.Module").
		Find(&levels).Error
	return levels, err
}

// GetUserMaxHierarchyLevel retorna o maior nível de hierarquia do usuário
func (r *resourceRole) GetUserMaxHierarchyLevel(userId, orgId string) (int, error) {
	// Verificar se é um admin (tabela admins) com permissão master_admin
	var admin models.Admin
	err := r.db.Where("id = ? AND deleted_at IS NULL", userId).First(&admin).Error
	if err == nil {
		for _, perm := range admin.Permissions {
			if perm == "master_admin" {
				return 10, nil // Admin Master tem nível máximo
			}
		}
	}

	// Verificar se o usuário é Master Admin pelo sistema legado (tabela users)
	var user models.User
	err = r.db.Where("id = ? AND deleted_at IS NULL", userId).First(&user).Error
	if err == nil {
		for _, perm := range user.Permissions {
			if perm == "master_admin" {
				return 10, nil // Master Admin tem nível máximo
			}
		}
	}

	// Se não é Master Admin pelo sistema legado, buscar da tabela user_roles
	var maxLevel int
	query := r.db.Model(&models.UserRole{}).
		Select("COALESCE(MAX(roles.hierarchy_level), 0)").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND user_roles.deleted_at IS NULL AND user_roles.active = true", userId).
		Where("roles.deleted_at IS NULL AND roles.active = true")

	// Tratar orgId vazio (cargos admin globais - organization_id IS NULL)
	if orgId != "" {
		query = query.Where("user_roles.organization_id = ?", orgId)
	} else {
		query = query.Where("user_roles.organization_id IS NULL")
	}

	err = query.Scan(&maxLevel).Error

	return maxLevel, err
}

// CanManageRole verifica se um usuário pode gerenciar um cargo específico
func (r *resourceRole) CanManageRole(actorUserId, targetRoleId, orgId string) (bool, error) {
	// Buscar nível máximo do ator
	actorLevel, err := r.GetUserMaxHierarchyLevel(actorUserId, orgId)
	if err != nil {
		return false, err
	}

	// Buscar nível do cargo alvo
	targetRole, err := r.GetById(targetRoleId)
	if err != nil {
		return false, err
	}

	// Ator pode gerenciar cargos com nível menor ou igual ao seu
	return actorLevel >= targetRole.HierarchyLevel, nil
}

// CountUsersByOrganization conta quantos usuários únicos estão associados a uma organização
func (r *resourceRole) CountUsersByOrganization(orgId string) (int, error) {
	var count int64
	err := r.db.Model(&models.UserRole{}).
		Select("COUNT(DISTINCT user_id)").
		Where("organization_id = ? AND deleted_at IS NULL AND active = true", orgId).
		Scan(&count).Error
	return int(count), err
}

// ListUsersByOrganization lista todos os user_roles de uma organização
func (r *resourceRole) ListUsersByOrganization(orgId string) ([]models.UserRole, error) {
	var userRoles []models.UserRole
	err := r.db.Where("organization_id = ? AND deleted_at IS NULL AND active = true", orgId).
		Preload("User").
		Preload("Role").
		Find(&userRoles).Error
	return userRoles, err
}

// GetAllUserRoles busca todos os user_roles de um usuário (todas as organizações)
func (r *resourceRole) GetAllUserRoles(userId string) ([]models.UserRole, error) {
	var userRoles []models.UserRole
	err := r.db.Where("user_id = ? AND deleted_at IS NULL AND active = true", userId).
		Preload("User").
		Preload("Role").
		Find(&userRoles).Error
	return userRoles, err
}
