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
	GetUsersByGroup(groupId string) ([]models.User, error)  // Deprecated: usar GetUsersByRole
	GetUsersByRole(role string) ([]models.User, error)
	ListUsersByOrganizationAndProject(orgId, projectId string) ([]models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	SoftDeleteUser(id string) error
	DeleteUser(id string) error
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
	// Corrigido: usar nomes corretos dos campos no banco
	err := r.db.Where("organization_id = ? AND project_id = ? AND deleted_at IS NULL", orgId, projectId).Find(&users).Error
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
