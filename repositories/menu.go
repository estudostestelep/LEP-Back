package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MenuRepository struct {
	db *gorm.DB
}

type IMenuRepository interface {
	GetMenu(id uuid.UUID) (*models.Menu, error)
	GetMenuById(id uuid.UUID) (*models.Menu, error)
	GetMenuList(organizationId, projectId uuid.UUID) ([]models.Menu, error)
	GetActiveMenuList(organizationId, projectId uuid.UUID) ([]models.Menu, error)
	CreateMenu(menu *models.Menu) error
	UpdateMenu(menu *models.Menu) error
	UpdateMenuOrder(id uuid.UUID, order int) error
	UpdateMenuStatus(id uuid.UUID, active bool) error
	SoftDelete(id uuid.UUID) error
	SoftDeleteMenu(id uuid.UUID) error
}

func NewConnMenu(db *gorm.DB) IMenuRepository {
	return &MenuRepository{db: db}
}

func (r *MenuRepository) CreateMenu(menu *models.Menu) error {
	return r.db.Create(menu).Error
}

func (r *MenuRepository) GetMenu(id uuid.UUID) (*models.Menu, error) {
	var menu models.Menu
	err := r.db.First(&menu, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *MenuRepository) GetMenuById(id uuid.UUID) (*models.Menu, error) {
	return r.GetMenu(id)
}

func (r *MenuRepository) GetMenuList(organizationId, projectId uuid.UUID) ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.Where("organization_id = ? AND project_id = ? AND deleted_at IS NULL", organizationId, projectId).
		Order(`"order" ASC`).
		Find(&menus).Error
	return menus, err
}

func (r *MenuRepository) GetActiveMenuList(organizationId, projectId uuid.UUID) ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.Where("organization_id = ? AND project_id = ? AND active = ? AND deleted_at IS NULL", organizationId, projectId, true).
		Order(`"order" ASC`).
		Find(&menus).Error
	return menus, err
}

func (r *MenuRepository) UpdateMenu(menu *models.Menu) error {
	return r.db.Save(menu).Error
}

func (r *MenuRepository) UpdateMenuOrder(id uuid.UUID, order int) error {
	return r.db.Model(&models.Menu{}).Where("id = ?", id).Update(`"order"`, order).Error
}

func (r *MenuRepository) UpdateMenuStatus(id uuid.UUID, active bool) error {
	return r.db.Model(&models.Menu{}).Where("id = ?", id).Update("active", active).Error
}

func (r *MenuRepository) SoftDelete(id uuid.UUID) error {
	return r.db.Model(&models.Menu{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

func (r *MenuRepository) SoftDeleteMenu(id uuid.UUID) error {
	return r.SoftDelete(id)
}
