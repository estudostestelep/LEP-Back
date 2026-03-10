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
		// Se não encontra, retorna com padrões (ID e OrganizationID serão preenchidos pelo caller)
		defaultSettings := &models.ProjectDisplaySettings{
			ID:              uuid.Nil, // Será gerado no servidor
			ProjectID:       projectUUID,
			OrganizationID:  uuid.Nil, // Será preenchido do header no servidor
			ShowPrepTime:    true,
			ShowRating:      true,
			ShowDescription: true,
			CreatedAt:       time.Time{}, // Será preenchido no servidor
			UpdatedAt:       time.Time{}, // Será preenchido no servidor
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
