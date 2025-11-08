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

	// ✨ Novos métodos para seleção automática de cardápio
	GetMenuOptions(organizationId, projectId uuid.UUID) ([]models.Menu, error)
	GetActiveMenuByTimeRange(organizationId, projectId uuid.UUID, currentTime time.Time) (*models.Menu, error)
	GetMenuWithHighestPriority(organizationId, projectId uuid.UUID) (*models.Menu, error)
	UpdateManualOverride(organizationId, projectId, menuId uuid.UUID) error

	// 🔍 Validação de menu
	CheckMenuNameExists(organizationId, projectId uuid.UUID, name string, excludeId *uuid.UUID) (bool, error)
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

// ✨ GetMenuOptions retorna lista de cardápios disponíveis para seleção
func (r *MenuRepository) GetMenuOptions(organizationId, projectId uuid.UUID) ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.
		Where("organization_id = ? AND project_id = ? AND active = ? AND deleted_at IS NULL", organizationId, projectId, true).
		Order(`"priority" ASC, "order" ASC, "created_at" DESC`).
		Find(&menus).Error
	return menus, err
}

// ✨ GetActiveMenuByTimeRange retorna o cardápio ativo baseado no range de horário
// Estratégia de seleção:
// 1. Se houver manual_override, retorna esse
// 2. Se houver cardápio com horário ativo, retorna o de maior prioridade
// 3. Caso contrário, retorna o de maior prioridade geral
func (r *MenuRepository) GetActiveMenuByTimeRange(organizationId, projectId uuid.UUID, currentTime time.Time) (*models.Menu, error) {
	// 1. Verificar se há manual override ativo
	var menuOverride models.Menu
	err := r.db.
		Where(
			"organization_id = ? AND project_id = ? AND is_manual_override = ? AND active = ? AND deleted_at IS NULL",
			organizationId, projectId, true, true,
		).
		First(&menuOverride).Error

	if err == nil {
		return &menuOverride, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 2. Buscar cardápios com horário ativo agora
	currentTimeStr := currentTime.Format("15:04") // HH:MM
	var menusInTimeRange []models.Menu

	// Query: verificar se hora atual está dentro do range de horário
	query := r.db.
		Where("organization_id = ? AND project_id = ? AND active = ? AND deleted_at IS NULL", organizationId, projectId, true).
		Where("time_range_start IS NOT NULL AND time_range_end IS NOT NULL").
		Where("TO_CHAR(time_range_start, 'HH24:MI') <= ? AND TO_CHAR(time_range_end, 'HH24:MI') >= ?", currentTimeStr, currentTimeStr).
		Order(`"priority" ASC, "order" ASC`)

	if err := query.Find(&menusInTimeRange).Error; err == nil && len(menusInTimeRange) > 0 {
		return &menusInTimeRange[0], nil
	}

	// 3. Fallback: retorna cardápio com maior prioridade
	var menuDefault models.Menu
	err = r.db.
		Where("organization_id = ? AND project_id = ? AND active = ? AND deleted_at IS NULL", organizationId, projectId, true).
		Order(`"priority" ASC, "order" ASC, "created_at" DESC`).
		First(&menuDefault).Error

	if err != nil {
		return nil, err
	}

	return &menuDefault, nil
}

// ✨ GetMenuWithHighestPriority retorna o cardápio com maior prioridade
func (r *MenuRepository) GetMenuWithHighestPriority(organizationId, projectId uuid.UUID) (*models.Menu, error) {
	var menu models.Menu
	err := r.db.
		Where("organization_id = ? AND project_id = ? AND active = ? AND deleted_at IS NULL", organizationId, projectId, true).
		Order(`"priority" ASC, "order" ASC, "created_at" DESC`).
		First(&menu).Error
	return &menu, err
}

// ✨ UpdateManualOverride define um cardápio como manual override (desativa os outros)
func (r *MenuRepository) UpdateManualOverride(organizationId, projectId, menuId uuid.UUID) error {
	// Desativar manual override de todos os outros cardápios
	if err := r.db.
		Model(&models.Menu{}).
		Where("organization_id = ? AND project_id = ? AND id != ? AND deleted_at IS NULL", organizationId, projectId, menuId).
		Update("is_manual_override", false).Error; err != nil {
		return err
	}

	// Ativar manual override para o cardápio selecionado
	return r.db.
		Model(&models.Menu{}).
		Where("id = ? AND deleted_at IS NULL", menuId).
		Update("is_manual_override", true).Error
}

// 🔍 CheckMenuNameExists verifica se já existe um menu com o mesmo nome no projeto
// excludeId é opcional: se fornecido, exclui esse ID da busca (útil para UPDATE)
func (r *MenuRepository) CheckMenuNameExists(organizationId, projectId uuid.UUID, name string, excludeId *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Where("organization_id = ? AND project_id = ? AND LOWER(name) = LOWER(?) AND deleted_at IS NULL",
		organizationId, projectId, name)

	// Se excludeId foi fornecido, excluir esse ID da busca (para UPDATE)
	if excludeId != nil {
		query = query.Where("id != ?", excludeId)
	}

	err := query.Model(&models.Menu{}).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
