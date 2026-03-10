package handler

import (
	"errors"
	"fmt"
	"lep/repositories"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
)

// Custom error types for menu operations
var (
	ErrMenuNameAlreadyExists = func(menuName string) error {
		return errors.New(fmt.Sprintf("Menu with name '%s' already exists in this project", menuName))
	}
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
	// 🔍 Verificar se já existe um menu com o mesmo nome no projeto
	// excludeId = nil porque é uma criação (não estamos atualizando)
	exists, err := r.repo.Menus.CheckMenuNameExists(menu.OrganizationId, menu.ProjectId, menu.Name, nil)
	if err != nil {
		return err
	}
	if exists {
		// Retornar erro indicando que o nome já existe
		return ErrMenuNameAlreadyExists(menu.Name)
	}

	menu.Id = uuid.New()
	menu.CreatedAt = time.Now()
	menu.UpdatedAt = time.Now()
	return r.repo.Menus.CreateMenu(menu)
}

func (r *resourceMenu) UpdateMenu(updatedMenu *models.Menu) error {
	// 🔍 Verificar se já existe outro menu com o mesmo nome no projeto
	// excludeId = &updatedMenu.Id para excluir o próprio menu da busca
	exists, err := r.repo.Menus.CheckMenuNameExists(updatedMenu.OrganizationId, updatedMenu.ProjectId, updatedMenu.Name, &updatedMenu.Id)
	if err != nil {
		return err
	}
	if exists {
		// Retornar erro indicando que o nome já existe em outro menu
		return ErrMenuNameAlreadyExists(updatedMenu.Name)
	}

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
