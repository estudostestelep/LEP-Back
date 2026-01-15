package handler

import (
	"encoding/json"
	"lep/repositories"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
)

type SidebarConfigHandler struct {
	repo repositories.ISidebarConfigRepository
}

type ISidebarConfigHandler interface {
	GetByOrganization(orgId string) (*models.SidebarConfigResponse, error)
	UpdateConfig(orgId string, items []models.SidebarItemConfig) (*models.SidebarConfigResponse, error)
	ResetToDefaults(orgId string) (*models.SidebarConfigResponse, error)
}

func NewSidebarConfigHandler(repo repositories.ISidebarConfigRepository) ISidebarConfigHandler {
	return &SidebarConfigHandler{repo: repo}
}

// getDefaultSidebarItems retorna as configurações padrão dos itens da sidebar
func getDefaultSidebarItems() []models.SidebarItemConfig {
	return []models.SidebarItemConfig{
		// Módulos gratuitos - sempre visíveis
		{ModuleCode: "client_menu", Behavior: models.BehaviorShow},
		{ModuleCode: "client_orders", Behavior: models.BehaviorShow},
		{ModuleCode: "client_customers", Behavior: models.BehaviorShow},
		{ModuleCode: "client_tables", Behavior: models.BehaviorShow},
		{ModuleCode: "client_products", Behavior: models.BehaviorShow},
		{ModuleCode: "client_users", Behavior: models.BehaviorShow},
		{ModuleCode: "client_settings", Behavior: models.BehaviorShow},
		{ModuleCode: "client_tags", Behavior: models.BehaviorShow},
		// Módulos premium - mostrar com cadeado
		{ModuleCode: "client_reservations", Behavior: models.BehaviorLock},
		{ModuleCode: "client_waitlist", Behavior: models.BehaviorLock},
		{ModuleCode: "client_reports", Behavior: models.BehaviorLock},
		// Módulos premium - esconder
		{ModuleCode: "client_notifications", Behavior: models.BehaviorHide},
	}
}

// getDefaultItemsJSON retorna os items padrão em formato JSON
func getDefaultItemsJSON() string {
	items := getDefaultSidebarItems()
	jsonBytes, _ := json.Marshal(items)
	return string(jsonBytes)
}

// GetByOrganization busca a configuração da sidebar por organização
func (h *SidebarConfigHandler) GetByOrganization(orgId string) (*models.SidebarConfigResponse, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}

	config, err := h.repo.GetOrCreate(orgUUID, getDefaultItemsJSON())
	if err != nil {
		return nil, err
	}

	var items []models.SidebarItemConfig
	if err := json.Unmarshal([]byte(config.ItemConfigs), &items); err != nil {
		return nil, err
	}

	return &models.SidebarConfigResponse{
		Id:             config.Id,
		OrganizationId: config.OrganizationId,
		Items:          items,
		CreatedAt:      config.CreatedAt,
		UpdatedAt:      config.UpdatedAt,
	}, nil
}

// UpdateConfig atualiza a configuração da sidebar
func (h *SidebarConfigHandler) UpdateConfig(orgId string, items []models.SidebarItemConfig) (*models.SidebarConfigResponse, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}

	// Buscar ou criar configuração
	config, err := h.repo.GetOrCreate(orgUUID, getDefaultItemsJSON())
	if err != nil {
		return nil, err
	}

	// Atualizar items
	itemsJSON, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}

	config.ItemConfigs = string(itemsJSON)
	config.UpdatedAt = time.Now()

	if err := h.repo.Update(config); err != nil {
		return nil, err
	}

	return &models.SidebarConfigResponse{
		Id:             config.Id,
		OrganizationId: config.OrganizationId,
		Items:          items,
		CreatedAt:      config.CreatedAt,
		UpdatedAt:      config.UpdatedAt,
	}, nil
}

// ResetToDefaults reseta a configuração para os valores padrão
func (h *SidebarConfigHandler) ResetToDefaults(orgId string) (*models.SidebarConfigResponse, error) {
	return h.UpdateConfig(orgId, getDefaultSidebarItems())
}
