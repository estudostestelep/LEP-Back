package repositories

import (
	"lep/repositories/models"
	"time"

	"gorm.io/gorm"
)

type resourceAdmin struct {
	db *gorm.DB
}

type IAdminRepository interface {
	GetAdminById(id string) (*models.Admin, error)
	GetAdminByEmail(email string) (*models.Admin, error)
	ListAdmins() ([]models.Admin, error)
	ListActiveAdmins() ([]models.Admin, error)
	CreateAdmin(admin *models.Admin) error
	UpdateAdmin(admin *models.Admin) error
	UpdateLastAccess(adminId string) error
	SoftDeleteAdmin(id string) error
	DeleteAdmin(id string) error
	AdminEmailExists(email string) (bool, error)
}

func NewAdminRepository(db *gorm.DB) IAdminRepository {
	return &resourceAdmin{db: db}
}

func (r *resourceAdmin) GetAdminById(id string) (*models.Admin, error) {
	var admin models.Admin
	err := r.db.First(&admin, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *resourceAdmin) GetAdminByEmail(email string) (*models.Admin, error) {
	var admin models.Admin
	err := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *resourceAdmin) ListAdmins() ([]models.Admin, error) {
	var admins []models.Admin
	err := r.db.Where("deleted_at IS NULL").Order("created_at DESC").Find(&admins).Error
	return admins, err
}

func (r *resourceAdmin) ListActiveAdmins() ([]models.Admin, error) {
	var admins []models.Admin
	err := r.db.Where("active = true AND deleted_at IS NULL").Order("created_at DESC").Find(&admins).Error
	return admins, err
}

func (r *resourceAdmin) CreateAdmin(admin *models.Admin) error {
	return r.db.Create(admin).Error
}

func (r *resourceAdmin) UpdateAdmin(admin *models.Admin) error {
	// Se o password estiver vazio, ignora o campo para não sobrescrever
	if admin.Password == "" {
		return r.db.Omit("Password").Save(admin).Error
	}
	return r.db.Save(admin).Error
}

func (r *resourceAdmin) UpdateLastAccess(adminId string) error {
	return r.db.Model(&models.Admin{}).Where("id = ?", adminId).Update("last_access_at", time.Now()).Error
}

func (r *resourceAdmin) SoftDeleteAdmin(id string) error {
	return r.db.Model(&models.Admin{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

func (r *resourceAdmin) DeleteAdmin(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Admin{}).Error
}

func (r *resourceAdmin) AdminEmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Admin{}).Where("email = ? AND deleted_at IS NULL", email).Count(&count).Error
	return count > 0, err
}
