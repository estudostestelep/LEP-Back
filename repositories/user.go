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
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	SoftDeleteUser(id string) error
	DeleteUser(id string) error
	GetUserWithRelations(id string) (*models.UserWithRelations, error)
	GetUserOrganizations(userId string) ([]models.UserOrganization, error)
	GetUserProjects(userId string) ([]models.UserProject, error)
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

	// Construir query base
	query := r.db.
		Distinct("users.*").
		Table("users").
		Joins("INNER JOIN user_organizations ON users.id = user_organizations.user_id").
		Where("user_organizations.organization_id = ? AND user_organizations.active = true AND user_organizations.deleted_at IS NULL", orgId).
		Where("users.deleted_at IS NULL")

	// Se projectId foi fornecido, filtrar também por projeto
	if projectId != "" {
		query = query.
			Joins("INNER JOIN user_projects ON users.id = user_projects.user_id").
			Where("user_projects.project_id = ? AND user_projects.active = true AND user_projects.deleted_at IS NULL", projectId)
	}

	err := query.Find(&users).Error
	return users, err
}

func (r *resourceUser) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *resourceUser) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *resourceUser) SoftDeleteUser(id string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

func (r *resourceUser) DeleteUser(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.User{}).Error
}

func (r *resourceUser) GetUserOrganizations(userId string) ([]models.UserOrganization, error) {
	var userOrgs []models.UserOrganization
	err := r.db.Where("user_id = ? AND deleted_at IS NULL", userId).Find(&userOrgs).Error
	return userOrgs, err
}

func (r *resourceUser) GetUserProjects(userId string) ([]models.UserProject, error) {
	var userProjs []models.UserProject
	err := r.db.Where("user_id = ? AND deleted_at IS NULL", userId).Find(&userProjs).Error
	return userProjs, err
}

func (r *resourceUser) GetUserWithRelations(id string) (*models.UserWithRelations, error) {
	// Buscar usuário
	user, err := r.GetUserById(id)
	if err != nil {
		return nil, err
	}

	// Buscar organizações do usuário
	userOrgs, err := r.GetUserOrganizations(id)
	if err != nil {
		return nil, err
	}

	// Buscar projetos do usuário
	userProjs, err := r.GetUserProjects(id)
	if err != nil {
		return nil, err
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
		Organizations: make([]models.UserOrgInfo, 0),
		Projects:      make([]models.UserProjInfo, 0),
	}

	// Converter organizações
	for _, org := range userOrgs {
		userWithRelations.Organizations = append(userWithRelations.Organizations, models.UserOrgInfo{
			OrganizationId: org.OrganizationId,
			Role:           org.Role,
			Active:         org.Active,
		})
	}

	// Converter projetos
	for _, proj := range userProjs {
		userWithRelations.Projects = append(userWithRelations.Projects, models.UserProjInfo{
			ProjectId: proj.ProjectId,
			Role:      proj.Role,
			Active:    proj.Active,
		})
	}

	return userWithRelations, nil
}
