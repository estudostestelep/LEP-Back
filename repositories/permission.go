package repositories

import (
	"lep/repositories/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type resourcePermission struct {
	db *gorm.DB
}

type IPermissionRepository interface {
	// CRUD de Permissions
	Create(permission *models.Permission) error
	GetById(id string) (*models.Permission, error)
	GetByCodeName(codeName string) (*models.Permission, error)
	Update(permission *models.Permission) error
	Delete(id string) error
	List() ([]models.Permission, error)
	ListByModule(moduleId string) ([]models.Permission, error)

	// Bulk operations
	CreateBulk(permissions []models.Permission) error
	GetByCodeNames(codeNames []string) ([]models.Permission, error)
}

func NewPermissionRepository(db *gorm.DB) IPermissionRepository {
	return &resourcePermission{db: db}
}

// Create cria uma nova permissão
func (r *resourcePermission) Create(permission *models.Permission) error {
	if permission.Id == uuid.Nil {
		permission.Id = uuid.New()
	}
	return r.db.Create(permission).Error
}

// GetById busca permissão por ID
func (r *resourcePermission) GetById(id string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).
		Preload("Module").
		First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetByCodeName busca permissão pelo código técnico
func (r *resourcePermission) GetByCodeName(codeName string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("code_name = ? AND deleted_at IS NULL", codeName).
		Preload("Module").
		First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// Update atualiza permissão
func (r *resourcePermission) Update(permission *models.Permission) error {
	return r.db.Save(permission).Error
}

// Delete faz soft delete da permissão
func (r *resourcePermission) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Permission{}).Error
}

// List lista todas as permissões ativas
func (r *resourcePermission) List() ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Where("deleted_at IS NULL AND active = true").
		Preload("Module").
		Order("module_id, code_name ASC").
		Find(&permissions).Error
	return permissions, err
}

// ListByModule lista permissões de um módulo específico
func (r *resourcePermission) ListByModule(moduleId string) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Where("module_id = ? AND deleted_at IS NULL AND active = true", moduleId).
		Preload("Module").
		Order("code_name ASC").
		Find(&permissions).Error
	return permissions, err
}

// CreateBulk cria múltiplas permissões de uma vez
func (r *resourcePermission) CreateBulk(permissions []models.Permission) error {
	for i := range permissions {
		if permissions[i].Id == uuid.Nil {
			permissions[i].Id = uuid.New()
		}
	}
	return r.db.Create(&permissions).Error
}

// GetByCodeNames busca permissões por lista de códigos
func (r *resourcePermission) GetByCodeNames(codeNames []string) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Where("code_name IN ? AND deleted_at IS NULL", codeNames).
		Preload("Module").
		Find(&permissions).Error
	return permissions, err
}
