package repositories

import (
	"lep/repositories/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type resourceModule struct {
	db *gorm.DB
}

type IModuleRepository interface {
	// CRUD de Modules
	Create(module *models.Module) error
	GetById(id string) (*models.Module, error)
	GetByCodeName(codeName string) (*models.Module, error)
	Update(module *models.Module) error
	Delete(id string) error
	List() ([]models.Module, error)
	ListByScope(scope string) ([]models.Module, error)

	// Module with permissions
	GetWithPermissions(id string) (*models.Module, error)
	ListWithPermissions() ([]models.Module, error)

	// Check access
	IsModuleInPackage(moduleId, packageId string) (bool, error)
	GetModulesForOrganization(orgId string) ([]models.Module, error)
}

func NewModuleRepository(db *gorm.DB) IModuleRepository {
	return &resourceModule{db: db}
}

// Create cria um novo módulo
func (r *resourceModule) Create(module *models.Module) error {
	if module.Id == uuid.Nil {
		module.Id = uuid.New()
	}
	return r.db.Create(module).Error
}

// GetById busca módulo por ID
func (r *resourceModule) GetById(id string) (*models.Module, error) {
	var module models.Module
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&module).Error
	if err != nil {
		return nil, err
	}
	return &module, nil
}

// GetByCodeName busca módulo pelo código técnico
func (r *resourceModule) GetByCodeName(codeName string) (*models.Module, error) {
	var module models.Module
	err := r.db.Where("code_name = ? AND deleted_at IS NULL", codeName).First(&module).Error
	if err != nil {
		return nil, err
	}
	return &module, nil
}

// Update atualiza módulo
func (r *resourceModule) Update(module *models.Module) error {
	return r.db.Save(module).Error
}

// Delete faz soft delete do módulo
func (r *resourceModule) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Module{}).Error
}

// List lista todos os módulos ativos
func (r *resourceModule) List() ([]models.Module, error) {
	var modules []models.Module
	err := r.db.Where("deleted_at IS NULL AND active = true").
		Order("display_order ASC, code_name ASC").
		Find(&modules).Error
	return modules, err
}

// ListByScope lista módulos por escopo (admin ou client)
func (r *resourceModule) ListByScope(scope string) ([]models.Module, error) {
	var modules []models.Module
	err := r.db.Where("scope = ? AND deleted_at IS NULL AND active = true", scope).
		Order("display_order ASC, code_name ASC").
		Find(&modules).Error
	return modules, err
}

// GetWithPermissions busca módulo com suas permissões
func (r *resourceModule) GetWithPermissions(id string) (*models.Module, error) {
	var module models.Module
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).
		Preload("Permissions", "deleted_at IS NULL AND active = true").
		First(&module).Error
	if err != nil {
		return nil, err
	}
	return &module, nil
}

// ListWithPermissions lista todos os módulos com suas permissões
func (r *resourceModule) ListWithPermissions() ([]models.Module, error) {
	var modules []models.Module
	err := r.db.Where("deleted_at IS NULL AND active = true").
		Preload("Permissions", "deleted_at IS NULL AND active = true").
		Order("display_order ASC, code_name ASC").
		Find(&modules).Error
	return modules, err
}

// IsModuleInPackage verifica se um módulo está incluído em um pacote
func (r *resourceModule) IsModuleInPackage(moduleId, packageId string) (bool, error) {
	var count int64
	err := r.db.Model(&models.PackageModule{}).
		Where("module_id = ? AND package_id = ? AND deleted_at IS NULL", moduleId, packageId).
		Count(&count).Error
	return count > 0, err
}

// GetModulesForOrganization retorna todos os módulos disponíveis para uma organização
// baseado no pacote contratado
func (r *resourceModule) GetModulesForOrganization(orgId string) ([]models.Module, error) {
	var modules []models.Module

	// Buscar o pacote ativo da organização
	var orgPackage models.OrganizationPackage
	err := r.db.Where("organization_id = ? AND active = true AND deleted_at IS NULL", orgId).
		Order("created_at DESC").
		First(&orgPackage).Error
	if err != nil {
		// Se não tem pacote, retornar módulos gratuitos
		err = r.db.Where("is_free = true AND deleted_at IS NULL AND active = true").
			Order("display_order ASC").
			Find(&modules).Error
		return modules, err
	}

	// Buscar módulos do pacote
	err = r.db.Table("modules").
		Distinct("modules.*").
		Joins("INNER JOIN package_modules ON modules.id = package_modules.module_id").
		Where("package_modules.package_id = ? AND package_modules.deleted_at IS NULL", orgPackage.PackageId).
		Where("modules.deleted_at IS NULL AND modules.active = true").
		Order("modules.display_order ASC").
		Find(&modules).Error

	return modules, err
}
