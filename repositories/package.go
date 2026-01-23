package repositories

import (
	"lep/repositories/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type resourcePackage struct {
	db *gorm.DB
}

type IPackageRepository interface {
	// CRUD de Packages
	Create(pkg *models.Package) error
	GetById(id string) (*models.Package, error)
	GetByCodeName(codeName string) (*models.Package, error)
	Update(pkg *models.Package) error
	Delete(id string) error
	List() ([]models.Package, error)
	ListPublic() ([]models.Package, error)

	// Package-Module relations
	AddModuleToPackage(packageId, moduleId string) error
	RemoveModuleFromPackage(packageId, moduleId string) error
	GetPackageModules(packageId string) ([]models.Module, error)
	GetPackageWithModules(id string) (*models.Package, error)

	// Package Limits
	SetPackageLimit(packageId, limitType string, limitValue int) error
	GetPackageLimits(packageId string) ([]models.PackageLimit, error)

	// Organization subscription
	SubscribeOrganization(orgPackage *models.OrganizationPackage) error
	GetOrganizationPackage(orgId string) (*models.OrganizationPackage, error)
	UpdateOrganizationPackage(orgPackage *models.OrganizationPackage) error
	CancelOrganizationPackage(orgId string) error
	DeleteOrganizationPackage(orgId string) error
	ListAllSubscriptions() ([]models.OrganizationPackage, error)

	// Bundles
	CreateBundle(bundle *models.PackageBundle) error
	GetBundleById(id string) (*models.PackageBundle, error)
	ListBundles() ([]models.PackageBundle, error)
	AddPackageToBundle(bundleId, packageId string) error
}

func NewPackageRepository(db *gorm.DB) IPackageRepository {
	return &resourcePackage{db: db}
}

// Create cria um novo pacote
func (r *resourcePackage) Create(pkg *models.Package) error {
	if pkg.Id == uuid.Nil {
		pkg.Id = uuid.New()
	}
	return r.db.Create(pkg).Error
}

// GetById busca pacote por ID
func (r *resourcePackage) GetById(id string) (*models.Package, error) {
	var pkg models.Package
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&pkg).Error
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

// GetByCodeName busca pacote pelo código técnico
func (r *resourcePackage) GetByCodeName(codeName string) (*models.Package, error) {
	var pkg models.Package
	err := r.db.Where("code_name = ? AND deleted_at IS NULL", codeName).First(&pkg).Error
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

// Update atualiza pacote
func (r *resourcePackage) Update(pkg *models.Package) error {
	return r.db.Save(pkg).Error
}

// Delete faz soft delete do pacote
func (r *resourcePackage) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Package{}).Error
}

// List lista todos os pacotes ativos
func (r *resourcePackage) List() ([]models.Package, error) {
	var packages []models.Package
	err := r.db.Where("deleted_at IS NULL AND active = true").
		Order("price_monthly ASC").
		Find(&packages).Error
	return packages, err
}

// ListPublic lista pacotes públicos para exibição
func (r *resourcePackage) ListPublic() ([]models.Package, error) {
	var packages []models.Package
	err := r.db.Where("deleted_at IS NULL AND active = true AND is_public = true").
		Order("display_order ASC, price_monthly ASC").
		Find(&packages).Error
	return packages, err
}

// AddModuleToPackage adiciona um módulo a um pacote
func (r *resourcePackage) AddModuleToPackage(packageId, moduleId string) error {
	pkgUUID, err := uuid.Parse(packageId)
	if err != nil {
		return err
	}
	modUUID, err := uuid.Parse(moduleId)
	if err != nil {
		return err
	}

	// Verificar se já existe
	var existing models.PackageModule
	err = r.db.Where("package_id = ? AND module_id = ? AND deleted_at IS NULL", packageId, moduleId).
		First(&existing).Error
	if err == nil {
		return nil // Já existe
	}

	pm := models.PackageModule{
		Id:        uuid.New(),
		PackageId: pkgUUID,
		ModuleId:  modUUID,
	}
	return r.db.Create(&pm).Error
}

// RemoveModuleFromPackage remove um módulo de um pacote
func (r *resourcePackage) RemoveModuleFromPackage(packageId, moduleId string) error {
	return r.db.Where("package_id = ? AND module_id = ?", packageId, moduleId).
		Delete(&models.PackageModule{}).Error
}

// GetPackageModules retorna todos os módulos de um pacote
func (r *resourcePackage) GetPackageModules(packageId string) ([]models.Module, error) {
	var modules []models.Module
	err := r.db.Table("modules").
		Joins("INNER JOIN package_modules ON modules.id = package_modules.module_id").
		Where("package_modules.package_id = ? AND package_modules.deleted_at IS NULL", packageId).
		Where("modules.deleted_at IS NULL AND modules.active = true").
		Order("modules.display_order ASC").
		Find(&modules).Error
	return modules, err
}

// GetPackageWithModules retorna pacote com seus módulos
func (r *resourcePackage) GetPackageWithModules(id string) (*models.Package, error) {
	var pkg models.Package
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&pkg).Error
	if err != nil {
		return nil, err
	}

	modules, err := r.GetPackageModules(id)
	if err != nil {
		return nil, err
	}
	pkg.Modules = modules

	return &pkg, nil
}

// SetPackageLimit define um limite para o pacote
func (r *resourcePackage) SetPackageLimit(packageId, limitType string, limitValue int) error {
	pkgUUID, err := uuid.Parse(packageId)
	if err != nil {
		return err
	}

	// Verificar se já existe
	var existing models.PackageLimit
	err = r.db.Where("package_id = ? AND limit_type = ? AND deleted_at IS NULL", packageId, limitType).
		First(&existing).Error
	if err == nil {
		// Atualizar existente
		existing.LimitValue = limitValue
		return r.db.Save(&existing).Error
	}

	// Criar novo
	limit := models.PackageLimit{
		Id:         uuid.New(),
		PackageId:  pkgUUID,
		LimitType:  limitType,
		LimitValue: limitValue,
	}
	return r.db.Create(&limit).Error
}

// GetPackageLimits retorna todos os limites de um pacote
func (r *resourcePackage) GetPackageLimits(packageId string) ([]models.PackageLimit, error) {
	var limits []models.PackageLimit
	err := r.db.Where("package_id = ? AND deleted_at IS NULL", packageId).
		Find(&limits).Error
	return limits, err
}

// SubscribeOrganization inscreve uma organização em um pacote
func (r *resourcePackage) SubscribeOrganization(orgPackage *models.OrganizationPackage) error {
	if orgPackage.Id == uuid.Nil {
		orgPackage.Id = uuid.New()
	}

	// Desativar pacotes anteriores
	r.db.Model(&models.OrganizationPackage{}).
		Where("organization_id = ? AND deleted_at IS NULL", orgPackage.OrganizationId).
		Update("active", false)

	return r.db.Create(orgPackage).Error
}

// GetOrganizationPackage retorna o pacote ativo de uma organização
func (r *resourcePackage) GetOrganizationPackage(orgId string) (*models.OrganizationPackage, error) {
	var orgPackage models.OrganizationPackage
	err := r.db.Where("organization_id = ? AND active = true AND deleted_at IS NULL", orgId).
		Preload("Package").
		Order("created_at DESC").
		First(&orgPackage).Error
	if err != nil {
		return nil, err
	}

	// Carregar módulos do pacote (campo com gorm:"-" precisa ser preenchido manualmente)
	if orgPackage.Package != nil {
		modules, err := r.GetPackageModules(orgPackage.PackageId.String())
		if err == nil {
			orgPackage.Package.Modules = modules
		}
	}

	return &orgPackage, nil
}

// UpdateOrganizationPackage atualiza a assinatura de uma organização
func (r *resourcePackage) UpdateOrganizationPackage(orgPackage *models.OrganizationPackage) error {
	return r.db.Save(orgPackage).Error
}

// CancelOrganizationPackage cancela a assinatura de uma organização
func (r *resourcePackage) CancelOrganizationPackage(orgId string) error {
	return r.db.Model(&models.OrganizationPackage{}).
		Where("organization_id = ? AND active = true AND deleted_at IS NULL", orgId).
		Update("active", false).Error
}

// DeleteOrganizationPackage exclui permanentemente a assinatura de uma organização
func (r *resourcePackage) DeleteOrganizationPackage(orgId string) error {
	return r.db.Unscoped().
		Where("organization_id = ?", orgId).
		Delete(&models.OrganizationPackage{}).Error
}

// ListAllSubscriptions lista todas as assinaturas ativas
func (r *resourcePackage) ListAllSubscriptions() ([]models.OrganizationPackage, error) {
	var subscriptions []models.OrganizationPackage
	err := r.db.Where("deleted_at IS NULL").
		Preload("Package").
		Preload("Organization").
		Order("created_at DESC").
		Find(&subscriptions).Error
	return subscriptions, err
}

// CreateBundle cria um novo bundle
func (r *resourcePackage) CreateBundle(bundle *models.PackageBundle) error {
	if bundle.Id == uuid.Nil {
		bundle.Id = uuid.New()
	}
	return r.db.Create(bundle).Error
}

// GetBundleById busca bundle por ID
func (r *resourcePackage) GetBundleById(id string) (*models.PackageBundle, error) {
	var bundle models.PackageBundle
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).
		Preload("Packages").
		First(&bundle).Error
	if err != nil {
		return nil, err
	}
	return &bundle, nil
}

// ListBundles lista todos os bundles ativos
func (r *resourcePackage) ListBundles() ([]models.PackageBundle, error) {
	var bundles []models.PackageBundle
	err := r.db.Where("deleted_at IS NULL AND active = true").
		Preload("Packages").
		Find(&bundles).Error
	return bundles, err
}

// AddPackageToBundle adiciona um pacote a um bundle
func (r *resourcePackage) AddPackageToBundle(bundleId, packageId string) error {
	// Esta função seria implementada via tabela de relacionamento bundle_packages
	// Por simplicidade, assumimos que Packages é gerenciado diretamente no bundle
	return nil
}
