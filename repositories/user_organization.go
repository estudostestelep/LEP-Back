package repositories

import (
	"lep/repositories/models"

	"gorm.io/gorm"
)

type resourceUserOrganization struct {
	db *gorm.DB
}

type IUserOrganizationRepository interface {
	// CRUD básico
	Create(userOrg *models.UserOrganization) error
	GetById(id string) (*models.UserOrganization, error)
	Update(userOrg *models.UserOrganization) error
	Delete(id string) error

	// Consultas específicas
	GetByUserAndOrganization(userId, orgId string) (*models.UserOrganization, error)
	ListByUser(userId string) ([]models.UserOrganization, error)
	ListByOrganization(orgId string) ([]models.UserOrganization, error)

	// Verificações
	UserBelongsToOrganization(userId, orgId string) (bool, error)
}

func NewUserOrganizationRepository(db *gorm.DB) IUserOrganizationRepository {
	return &resourceUserOrganization{db: db}
}

func (r *resourceUserOrganization) Create(userOrg *models.UserOrganization) error {
	return r.db.Create(userOrg).Error
}

func (r *resourceUserOrganization) GetById(id string) (*models.UserOrganization, error) {
	var userOrg models.UserOrganization
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&userOrg).Error
	if err != nil {
		return nil, err
	}
	return &userOrg, nil
}

func (r *resourceUserOrganization) Update(userOrg *models.UserOrganization) error {
	return r.db.Save(userOrg).Error
}

func (r *resourceUserOrganization) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.UserOrganization{}).Error
}

func (r *resourceUserOrganization) GetByUserAndOrganization(userId, orgId string) (*models.UserOrganization, error) {
	var userOrg models.UserOrganization
	err := r.db.Where("user_id = ? AND organization_id = ? AND deleted_at IS NULL", userId, orgId).First(&userOrg).Error
	if err != nil {
		return nil, err
	}
	return &userOrg, nil
}

func (r *resourceUserOrganization) ListByUser(userId string) ([]models.UserOrganization, error) {
	var userOrgs []models.UserOrganization
	err := r.db.Where("user_id = ? AND deleted_at IS NULL AND active = true", userId).Find(&userOrgs).Error
	return userOrgs, err
}

func (r *resourceUserOrganization) ListByOrganization(orgId string) ([]models.UserOrganization, error) {
	var userOrgs []models.UserOrganization
	err := r.db.Where("organization_id = ? AND deleted_at IS NULL AND active = true", orgId).Find(&userOrgs).Error
	return userOrgs, err
}

func (r *resourceUserOrganization) UserBelongsToOrganization(userId, orgId string) (bool, error) {
	var count int64
	err := r.db.Model(&models.UserOrganization{}).
		Where("user_id = ? AND organization_id = ? AND deleted_at IS NULL AND active = true", userId, orgId).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
