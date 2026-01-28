package repositories

import (
	"lep/repositories/models"
	"time"

	"gorm.io/gorm"
)

type resourceUser struct {
	db *gorm.DB
}

type IUserRepository interface {
	GetUserById(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUsersByGroup(groupId string) ([]models.User, error) // Deprecated: usar GetUsersByRole
	GetUsersByRole(role string) ([]models.User, error)
	ListUsersByOrganizationAndProject(orgId, projectId string) ([]models.User, error)
	ListUsersWithRoles(orgId, projectId string) ([]models.UserWithRole, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	UpdateLastAccess(userId string) error
	SoftDeleteUser(id string) error
	DeleteUser(id string) error
	GetUserWithRelations(id string) (*models.UserWithRelations, error)
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &resourceUser{db: db}
}

func (r *resourceUser) GetUserById(id string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *resourceUser) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *resourceUser) GetUsersByGroup(groupId string) ([]models.User, error) {
	var users []models.User
	// Corrigido: usar 'role' que existe no modelo User ao invés de 'group_member' inexistente
	err := r.db.Where("role = ? AND deleted_at IS NULL", groupId).Find(&users).Error
	return users, err
}

func (r *resourceUser) GetUsersByRole(role string) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("role = ? AND deleted_at IS NULL", role).Find(&users).Error
	return users, err
}

func (r *resourceUser) ListUsersByOrganizationAndProject(orgId, projectId string) ([]models.User, error) {
	var users []models.User

	// Construir query base usando user_roles em vez de user_organizations
	query := r.db.
		Distinct("users.*").
		Table("users").
		Joins("INNER JOIN user_roles ON users.id = user_roles.user_id").
		Where("user_roles.organization_id = ? AND user_roles.active = true AND user_roles.deleted_at IS NULL", orgId).
		Where("users.deleted_at IS NULL")

	// Se projectId foi fornecido, filtrar também por projeto
	if projectId != "" {
		query = query.Where("user_roles.project_id = ?", projectId)
	}

	err := query.Find(&users).Error
	return users, err
}

func (r *resourceUser) ListUsersWithRoles(orgId, projectId string) ([]models.UserWithRole, error) {
	// Primeiro buscar os usuários
	users, err := r.ListUsersByOrganizationAndProject(orgId, projectId)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return []models.UserWithRole{}, nil
	}

	// Coletar IDs dos usuários
	userIds := make([]string, len(users))
	for i, u := range users {
		userIds[i] = u.Id.String()
	}

	// Buscar user_roles com role preloaded em uma única query
	var userRoles []models.UserRole
	r.db.Where("user_id IN (?) AND deleted_at IS NULL AND active = true", userIds).
		Preload("Role").
		Find(&userRoles)

	// Mapear userId -> role info (pegar o de maior hierarquia)
	type roleInfo struct {
		Name           string
		DisplayName    string
		HierarchyLevel int
	}
	roleMap := make(map[string]roleInfo)
	for _, ur := range userRoles {
		if ur.Role == nil {
			continue
		}
		existing, ok := roleMap[ur.UserId.String()]
		if !ok || ur.Role.HierarchyLevel > existing.HierarchyLevel {
			roleMap[ur.UserId.String()] = roleInfo{
				Name:           ur.Role.Name,
				DisplayName:    ur.Role.DisplayName,
				HierarchyLevel: ur.Role.HierarchyLevel,
			}
		}
	}

	// Montar resultado
	result := make([]models.UserWithRole, len(users))
	for i, u := range users {
		result[i] = models.UserWithRole{
			User: u,
		}
		if info, ok := roleMap[u.Id.String()]; ok {
			result[i].RoleName = info.Name
			result[i].RoleDisplayName = info.DisplayName
			result[i].RoleHierarchyLevel = info.HierarchyLevel
		}
	}

	return result, nil
}

func (r *resourceUser) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *resourceUser) UpdateUser(user *models.User) error {
	// Se o password estiver vazio, ignora o campo para não sobrescrever o valor existente
	if user.Password == "" {
		return r.db.Omit("Password").Save(user).Error
	}
	return r.db.Save(user).Error
}

func (r *resourceUser) UpdateLastAccess(userId string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userId).Update("last_access_at", time.Now()).Error
}

func (r *resourceUser) SoftDeleteUser(id string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

func (r *resourceUser) DeleteUser(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.User{}).Error
}

// GetUserWithRelations busca um usuário com suas organizações e projetos via user_roles
func (r *resourceUser) GetUserWithRelations(id string) (*models.UserWithRelations, error) {
	// Buscar usuário
	user, err := r.GetUserById(id)
	if err != nil {
		return nil, err
	}

	// Buscar user_roles do usuário
	var userRoles []models.UserRole
	err = r.db.Where("user_id = ? AND deleted_at IS NULL AND active = true", id).
		Preload("Role").
		Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	// Agrupar organizações e projetos únicos
	orgMap := make(map[string]models.UserOrgInfo)
	projMap := make(map[string]models.UserProjInfo)

	for _, ur := range userRoles {
		roleName := ""
		if ur.Role != nil {
			roleName = ur.Role.DisplayName
		}

		// Adicionar organização se existir
		if ur.OrganizationId != nil {
			orgId := ur.OrganizationId.String()
			if _, exists := orgMap[orgId]; !exists {
				orgMap[orgId] = models.UserOrgInfo{
					OrganizationId: *ur.OrganizationId,
					Role:           roleName,
					Active:         ur.Active,
				}
			}
		}

		// Adicionar projeto se existir
		if ur.ProjectId != nil {
			projId := ur.ProjectId.String()
			if _, exists := projMap[projId]; !exists {
				projMap[projId] = models.UserProjInfo{
					ProjectId: *ur.ProjectId,
					Role:      roleName,
					Active:    ur.Active,
				}
			}
		}
	}

	// Converter maps para slices
	orgs := make([]models.UserOrgInfo, 0, len(orgMap))
	for _, org := range orgMap {
		orgs = append(orgs, org)
	}

	projs := make([]models.UserProjInfo, 0, len(projMap))
	for _, proj := range projMap {
		projs = append(projs, proj)
	}

	// Montar DTO
	userWithRelations := &models.UserWithRelations{
		Id:            user.Id,
		Name:          user.Name,
		Email:         user.Email,
		Permissions:   user.Permissions,
		Active:        user.Active,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		DeletedAt:     user.DeletedAt,
		Organizations: orgs,
		Projects:      projs,
	}

	return userWithRelations, nil
}
