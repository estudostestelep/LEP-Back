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
	IsModuleInPlan(moduleId, planId string) (bool, error)
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

// IsModuleInPlan verifica se um módulo está incluído em um plano
func (r *resourceModule) IsModuleInPlan(moduleId, planId string) (bool, error) {
	var count int64
	err := r.db.Model(&models.PlanModule{}).
		Where("module_id = ? AND plan_id = ?", moduleId, planId).
		Count(&count).Error
	return count > 0, err
}

// GetModulesForOrganization retorna todos os módulos disponíveis para uma organização
// baseado no plano contratado
func (r *resourceModule) GetModulesForOrganization(orgId string) ([]models.Module, error) {
	var modules []models.Module

	// Buscar o plano ativo da organização
	var orgPlan models.OrganizationPlan
	err := r.db.Where("organization_id = ? AND active = true AND deleted_at IS NULL", orgId).
		Order("created_at DESC").
		First(&orgPlan).Error
	if err != nil {
		// Se não tem plano, retornar módulos gratuitos
		err = r.db.Where("is_free = true AND deleted_at IS NULL AND active = true").
			Order("display_order ASC").
			Find(&modules).Error
		return modules, err
	}

	// Buscar módulos do plano
	err = r.db.Table("modules").
		Distinct("modules.*").
		Joins("INNER JOIN plan_modules ON modules.id = plan_modules.module_id").
		Where("plan_modules.plan_id = ?", orgPlan.PlanId).
		Where("modules.deleted_at IS NULL AND modules.active = true").
		Order("modules.display_order ASC").
		Find(&modules).Error

	return modules, err
}
