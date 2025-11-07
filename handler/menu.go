package handler

import (
	"lep/repositories"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
)

type resourceMenu struct {
	repo *repositories.DBconn
}

type IHandlerMenu interface {
	GetMenu(id string) (*models.Menu, error)
	ListMenus(orgId, projectId string) ([]models.Menu, error)
	ListActiveMenus(orgId, projectId string) ([]models.Menu, error)
	CreateMenu(menu *models.Menu) error
	UpdateMenu(updatedMenu *models.Menu) error
	UpdateMenuOrder(id string, order int) error
	UpdateMenuStatus(id string, active bool) error
	DeleteMenu(id string) error

	// ✨ Novos métodos para seleção inteligente de cardápio
	GetMenuOptions(orgId, projectId string) ([]models.Menu, error)
	GetActiveMenuByTimeRange(orgId, projectId string) (*models.Menu, error)
	SetMenuAsManualOverride(orgId, projectId, menuId string) error
	RemoveManualOverride(orgId, projectId string) error
}

func (r *resourceMenu) GetMenu(id string) (*models.Menu, error) {
	menuId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return r.repo.Menus.GetMenuById(menuId)
}

func (r *resourceMenu) ListMenus(orgId, projectId string) ([]models.Menu, error) {
	orgUuid, _ := uuid.Parse(orgId)
	projectUuid, _ := uuid.Parse(projectId)
	return r.repo.Menus.GetMenuList(orgUuid, projectUuid)
}

func (r *resourceMenu) ListActiveMenus(orgId, projectId string) ([]models.Menu, error) {
	orgUuid, _ := uuid.Parse(orgId)
	projectUuid, _ := uuid.Parse(projectId)
	return r.repo.Menus.GetActiveMenuList(orgUuid, projectUuid)
}

func (r *resourceMenu) CreateMenu(menu *models.Menu) error {
	menu.Id = uuid.New()
	menu.CreatedAt = time.Now()
	menu.UpdatedAt = time.Now()
	return r.repo.Menus.CreateMenu(menu)
}

func (r *resourceMenu) UpdateMenu(updatedMenu *models.Menu) error {
	updatedMenu.UpdatedAt = time.Now()
	return r.repo.Menus.UpdateMenu(updatedMenu)
}

func (r *resourceMenu) UpdateMenuOrder(id string, order int) error {
	menuId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.repo.Menus.UpdateMenuOrder(menuId, order)
}

func (r *resourceMenu) UpdateMenuStatus(id string, active bool) error {
	menuId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.repo.Menus.UpdateMenuStatus(menuId, active)
}

func (r *resourceMenu) DeleteMenu(id string) error {
	menuId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.repo.Menus.SoftDeleteMenu(menuId)
}

// ✨ GetMenuOptions retorna lista de opções de cardápio
func (r *resourceMenu) GetMenuOptions(orgId, projectId string) ([]models.Menu, error) {
	orgUuid, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projectUuid, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	return r.repo.Menus.GetMenuOptions(orgUuid, projectUuid)
}

// ✨ GetActiveMenuByTimeRange retorna o cardápio ativo baseado em horário
func (r *resourceMenu) GetActiveMenuByTimeRange(orgId, projectId string) (*models.Menu, error) {
	orgUuid, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projectUuid, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	return r.repo.Menus.GetActiveMenuByTimeRange(orgUuid, projectUuid, time.Now())
}

// ✨ SetMenuAsManualOverride define um cardápio como override manual
func (r *resourceMenu) SetMenuAsManualOverride(orgId, projectId, menuId string) error {
	orgUuid, err := uuid.Parse(orgId)
	if err != nil {
		return err
	}
	projectUuid, err := uuid.Parse(projectId)
	if err != nil {
		return err
	}
	menuUuid, err := uuid.Parse(menuId)
	if err != nil {
		return err
	}
	return r.repo.Menus.UpdateManualOverride(orgUuid, projectUuid, menuUuid)
}

// ✨ RemoveManualOverride remove o override manual
func (r *resourceMenu) RemoveManualOverride(orgId, projectId string) error {
	menus, err := r.ListMenus(orgId, projectId)
	if err != nil {
		return err
	}

	for _, menu := range menus {
		if menu.IsManualOverride {
			menu.IsManualOverride = false
			return r.UpdateMenu(&menu)
		}
	}

	return nil
}

func NewSourceHandlerMenu(repo *repositories.DBconn) IHandlerMenu {
	return &resourceMenu{repo: repo}
}
