package handler

import (
	"lep/repositories"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
)

type SettingsHandler struct {
	settingsRepo repositories.ISettingsRepository
}

type ISettingsHandler interface {
	GetSettingsByProject(orgId, projectId string) (*models.Settings, error)
	UpdateSettings(settings *models.Settings) error
	GetOrCreateSettings(orgId, projectId string) (*models.Settings, error)
}

func NewSettingsHandler(settingsRepo repositories.ISettingsRepository) ISettingsHandler {
	return &SettingsHandler{settingsRepo: settingsRepo}
}

// GetSettingsByProject busca configurações por projeto
func (h *SettingsHandler) GetSettingsByProject(orgId, projectId string) (*models.Settings, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}

	projectUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}

	return h.settingsRepo.GetSettingsByProject(orgUUID, projectUUID)
}

// UpdateSettings atualiza configurações
func (h *SettingsHandler) UpdateSettings(settings *models.Settings) error {
	settings.UpdatedAt = time.Now()
	return h.settingsRepo.UpdateSettings(settings)
}

// GetOrCreateSettings busca ou cria configurações padrão
func (h *SettingsHandler) GetOrCreateSettings(orgId, projectId string) (*models.Settings, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}

	projectUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}

	return h.settingsRepo.GetOrCreateSettings(orgUUID, projectUUID)
}