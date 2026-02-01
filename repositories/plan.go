package repositories

import (
	"lep/repositories/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type resourcePlan struct {
	db *gorm.DB
}

// IPlanRepository interface para operações com Plans
type IPlanRepository interface {
	// CRUD de Plans
	Create(plan *models.Plan) error
	GetById(id string) (*models.Plan, error)
	GetByCode(code string) (*models.Plan, error)
	Update(plan *models.Plan) error
	Delete(id string) error
	List() ([]models.Plan, error)
	ListPublic() ([]models.Plan, error)

	// Plan-Module relations
	AddModuleToPlan(planId, moduleId string) error
	RemoveModuleFromPlan(planId, moduleId string) error
	GetPlanModules(planId string) ([]models.Module, error)
	GetPlanWithModules(id string) (*models.Plan, error)

	// Plan Limits
	SetPlanLimit(planId, limitType string, limitValue int) error
	GetPlanLimits(planId string) ([]models.PlanLimit, error)

	// Organization subscription
	SubscribeOrganization(orgPlan *models.OrganizationPlan) error
	GetOrganizationPlan(orgId string) (*models.OrganizationPlan, error)
	UpdateOrganizationPlan(orgPlan *models.OrganizationPlan) error
	CancelOrganizationPlan(orgId string) error
	DeleteOrganizationPlan(orgId string) error
	ListAllSubscriptions() ([]models.OrganizationPlan, error)

	// Module access check
	OrganizationHasModule(orgId, moduleCode string) (bool, error)
	GetOrganizationModules(orgId string) ([]models.Module, error)
}

func NewPlanRepository(db *gorm.DB) IPlanRepository {
	return &resourcePlan{db: db}
}

// ==================== CRUD ====================

// Create cria um novo plano
func (r *resourcePlan) Create(plan *models.Plan) error {
	if plan.Id == uuid.Nil {
		plan.Id = uuid.New()
	}
	return r.db.Create(plan).Error
}

// GetById busca plano por ID
func (r *resourcePlan) GetById(id string) (*models.Plan, error) {
	var plan models.Plan
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// GetByCode busca plano pelo código
func (r *resourcePlan) GetByCode(code string) (*models.Plan, error) {
	var plan models.Plan
	err := r.db.Where("code = ? AND deleted_at IS NULL", code).First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// Update atualiza plano
func (r *resourcePlan) Update(plan *models.Plan) error {
	return r.db.Save(plan).Error
}

// Delete faz soft delete do plano
func (r *resourcePlan) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Plan{}).Error
}

// List lista todos os planos ativos
func (r *resourcePlan) List() ([]models.Plan, error) {
	var plans []models.Plan
	err := r.db.Where("deleted_at IS NULL AND active = true").
		Order("display_order ASC, price_monthly ASC").
		Find(&plans).Error
	return plans, err
}

// ListPublic lista planos públicos para exibição
func (r *resourcePlan) ListPublic() ([]models.Plan, error) {
	var plans []models.Plan
	err := r.db.Where("deleted_at IS NULL AND active = true AND is_public = true").
		Order("display_order ASC, price_monthly ASC").
		Find(&plans).Error
	return plans, err
}

// ==================== Plan-Module Relations ====================

// AddModuleToPlan adiciona um módulo a um plano
func (r *resourcePlan) AddModuleToPlan(planId, moduleId string) error {
	planUUID, err := uuid.Parse(planId)
	if err != nil {
		return err
	}
	modUUID, err := uuid.Parse(moduleId)
	if err != nil {
		return err
	}

	// Verificar se já existe
	var existing models.PlanModule
	err = r.db.Where("plan_id = ? AND module_id = ?", planId, moduleId).
		First(&existing).Error
	if err == nil {
		return nil // Já existe
	}

	pm := models.PlanModule{
		PlanId:   planUUID,
		ModuleId: modUUID,
	}
	return r.db.Create(&pm).Error
}

// RemoveModuleFromPlan remove um módulo de um plano
func (r *resourcePlan) RemoveModuleFromPlan(planId, moduleId string) error {
	return r.db.Where("plan_id = ? AND module_id = ?", planId, moduleId).
		Delete(&models.PlanModule{}).Error
}

// GetPlanModules retorna todos os módulos de um plano
func (r *resourcePlan) GetPlanModules(planId string) ([]models.Module, error) {
	var modules []models.Module
	err := r.db.Table("modules").
		Joins("INNER JOIN plan_modules ON modules.id = plan_modules.module_id").
		Where("plan_modules.plan_id = ?", planId).
		Where("modules.deleted_at IS NULL AND modules.active = true").
		Order("modules.display_order ASC").
		Find(&modules).Error
	return modules, err
}

// GetPlanWithModules retorna plano com seus módulos
func (r *resourcePlan) GetPlanWithModules(id string) (*models.Plan, error) {
	var plan models.Plan
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&plan).Error
	if err != nil {
		return nil, err
	}

	modules, err := r.GetPlanModules(id)
	if err != nil {
		return nil, err
	}
	plan.Modules = modules

	// Carregar limites também
	limits, _ := r.GetPlanLimits(id)
	plan.Limits = limits

	return &plan, nil
}

// ==================== Plan Limits ====================

// SetPlanLimit define um limite para o plano
func (r *resourcePlan) SetPlanLimit(planId, limitType string, limitValue int) error {
	planUUID, err := uuid.Parse(planId)
	if err != nil {
		return err
	}

	// Verificar se já existe
	var existing models.PlanLimit
	err = r.db.Where("plan_id = ? AND limit_type = ? AND deleted_at IS NULL", planId, limitType).
		First(&existing).Error
	if err == nil {
		// Atualizar existente
		existing.LimitValue = limitValue
		return r.db.Save(&existing).Error
	}

	// Criar novo
	limit := models.PlanLimit{
		Id:         uuid.New(),
		PlanId:     planUUID,
		LimitType:  limitType,
		LimitValue: limitValue,
	}
	return r.db.Create(&limit).Error
}

// GetPlanLimits retorna todos os limites de um plano
func (r *resourcePlan) GetPlanLimits(planId string) ([]models.PlanLimit, error) {
	var limits []models.PlanLimit
	err := r.db.Where("plan_id = ? AND deleted_at IS NULL", planId).
		Find(&limits).Error
	return limits, err
}

// ==================== Organization Subscription ====================

// SubscribeOrganization inscreve uma organização em um plano
func (r *resourcePlan) SubscribeOrganization(orgPlan *models.OrganizationPlan) error {
	if orgPlan.Id == uuid.Nil {
		orgPlan.Id = uuid.New()
	}

	// Desativar planos anteriores
	r.db.Model(&models.OrganizationPlan{}).
		Where("organization_id = ? AND deleted_at IS NULL", orgPlan.OrganizationId).
		Update("active", false)

	return r.db.Create(orgPlan).Error
}

// GetOrganizationPlan retorna o plano ativo de uma organização
func (r *resourcePlan) GetOrganizationPlan(orgId string) (*models.OrganizationPlan, error) {
	var orgPlan models.OrganizationPlan
	err := r.db.Where("organization_id = ? AND active = true AND deleted_at IS NULL", orgId).
		Preload("Plan").
		Order("created_at DESC").
		First(&orgPlan).Error
	if err != nil {
		return nil, err
	}

	// Carregar módulos do plano
	if orgPlan.Plan != nil {
		modules, err := r.GetPlanModules(orgPlan.PlanId.String())
		if err == nil {
			orgPlan.Plan.Modules = modules
		}
		limits, err := r.GetPlanLimits(orgPlan.PlanId.String())
		if err == nil {
			orgPlan.Plan.Limits = limits
		}
	}

	return &orgPlan, nil
}

// UpdateOrganizationPlan atualiza a assinatura de uma organização
func (r *resourcePlan) UpdateOrganizationPlan(orgPlan *models.OrganizationPlan) error {
	return r.db.Save(orgPlan).Error
}

// CancelOrganizationPlan cancela a assinatura de uma organização
func (r *resourcePlan) CancelOrganizationPlan(orgId string) error {
	return r.db.Model(&models.OrganizationPlan{}).
		Where("organization_id = ? AND active = true AND deleted_at IS NULL", orgId).
		Update("active", false).Error
}

// DeleteOrganizationPlan exclui permanentemente a assinatura
func (r *resourcePlan) DeleteOrganizationPlan(orgId string) error {
	return r.db.Unscoped().
		Where("organization_id = ?", orgId).
		Delete(&models.OrganizationPlan{}).Error
}

// ListAllSubscriptions lista todas as assinaturas
func (r *resourcePlan) ListAllSubscriptions() ([]models.OrganizationPlan, error) {
	var subscriptions []models.OrganizationPlan
	err := r.db.Where("deleted_at IS NULL").
		Preload("Plan").
		Preload("Organization").
		Order("created_at DESC").
		Find(&subscriptions).Error
	return subscriptions, err
}

// ==================== Module Access ====================

// OrganizationHasModule verifica se organização tem acesso a um módulo
func (r *resourcePlan) OrganizationHasModule(orgId, moduleCode string) (bool, error) {
	// Primeiro, verificar se é um módulo gratuito
	var freeModule models.Module
	err := r.db.Where("code = ? AND is_free = true AND active = true AND deleted_at IS NULL", moduleCode).
		First(&freeModule).Error
	if err == nil {
		return true, nil // Módulo é gratuito
	}

	// Verificar se o plano da organização inclui o módulo
	var count int64
	err = r.db.Table("organization_plans").
		Joins("INNER JOIN plan_modules ON organization_plans.plan_id = plan_modules.plan_id").
		Joins("INNER JOIN modules ON plan_modules.module_id = modules.id").
		Where("organization_plans.organization_id = ?", orgId).
		Where("organization_plans.active = true").
		Where("organization_plans.deleted_at IS NULL").
		Where("modules.code = ?", moduleCode).
		Where("modules.active = true").
		Where("modules.deleted_at IS NULL").
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetOrganizationModules retorna todos os módulos disponíveis para uma organização
func (r *resourcePlan) GetOrganizationModules(orgId string) ([]models.Module, error) {
	var modules []models.Module

	// Buscar módulos do plano ativo
	err := r.db.Table("modules").
		Joins("LEFT JOIN plan_modules ON modules.id = plan_modules.module_id").
		Joins("LEFT JOIN organization_plans ON plan_modules.plan_id = organization_plans.plan_id").
		Where("(organization_plans.organization_id = ? AND organization_plans.active = true AND organization_plans.deleted_at IS NULL) OR modules.is_free = true", orgId).
		Where("modules.active = true AND modules.deleted_at IS NULL").
		Distinct().
		Order("modules.display_order ASC").
		Find(&modules).Error

	return modules, err
}
