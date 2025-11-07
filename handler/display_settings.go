package handler

import (
	"lep/repositories"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
)

type DisplaySettingsHandler struct {
	displaySettingsRepo repositories.IDisplaySettingsRepository
}

type IDisplaySettingsHandler interface {
	GetSettingsByProject(projectId string) (*models.ProjectDisplaySettings, error)
	UpdateSettings(settings *models.ProjectDisplaySettings) error
	ResetToDefaults(projectId string) (*models.ProjectDisplaySettings, error)
}

func NewDisplaySettingsHandler(displaySettingsRepo repositories.IDisplaySettingsRepository) IDisplaySettingsHandler {
	return &DisplaySettingsHandler{displaySettingsRepo: displaySettingsRepo}
}

// GetSettingsByProject busca configurações de exibição de produtos
func (h *DisplaySettingsHandler) GetSettingsByProject(projectId string) (*models.ProjectDisplaySettings, error) {
	projectUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}

	settings, err := h.displaySettingsRepo.GetSettingsByProject(projectUUID)
	if err != nil {
		// Se não encontra, retorna com padrões
		defaultSettings := &models.ProjectDisplaySettings{
			ID:              uuid.New(),
			ProjectID:       projectUUID,
			ShowPrepTime:    true,
			ShowRating:      true,
			ShowDescription: true,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		return defaultSettings, nil
	}
	return settings, nil
}

// UpdateSettings atualiza configurações de exibição
func (h *DisplaySettingsHandler) UpdateSettings(settings *models.ProjectDisplaySettings) error {
	settings.UpdatedAt = time.Now()
	return h.displaySettingsRepo.UpdateSettings(settings)
}

// ResetToDefaults reseta configurações para valores padrão
func (h *DisplaySettingsHandler) ResetToDefaults(projectId string) (*models.ProjectDisplaySettings, error) {
	projectUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}

	return h.displaySettingsRepo.ResetToDefaults(projectUUID)
}
