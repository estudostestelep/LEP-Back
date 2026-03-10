package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DisplaySettingsRepository struct {
	db *gorm.DB
}

type IDisplaySettingsRepository interface {
	GetSettingsByProject(projectId uuid.UUID) (*models.ProjectDisplaySettings, error)
	CreateSettings(settings *models.ProjectDisplaySettings) error
	UpdateSettings(settings *models.ProjectDisplaySettings) error
	ResetToDefaults(projectId uuid.UUID) (*models.ProjectDisplaySettings, error)
}

func NewDisplaySettingsRepository(db *gorm.DB) IDisplaySettingsRepository {
	return &DisplaySettingsRepository{db: db}
}

// GetSettingsByProject busca configurações de exibição por projeto
func (r *DisplaySettingsRepository) GetSettingsByProject(projectId uuid.UUID) (*models.ProjectDisplaySettings, error) {
	var settings models.ProjectDisplaySettings
	err := r.db.Where("project_id = ?", projectId).First(&settings).Error
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

// CreateSettings cria novas configurações de exibição
func (r *DisplaySettingsRepository) CreateSettings(settings *models.ProjectDisplaySettings) error {
	return r.db.Create(settings).Error
}

// UpdateSettings atualiza configurações de exibição existentes
func (r *DisplaySettingsRepository) UpdateSettings(settings *models.ProjectDisplaySettings) error {
	settings.UpdatedAt = time.Now()
	return r.db.Save(settings).Error
}

// ResetToDefaults reseta configurações para valores padrão
func (r *DisplaySettingsRepository) ResetToDefaults(projectId uuid.UUID) (*models.ProjectDisplaySettings, error) {
	// Busca configuração existente
	settings, err := r.GetSettingsByProject(projectId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Se não existe, cria com padrões
	if err == gorm.ErrRecordNotFound {
		newSettings := &models.ProjectDisplaySettings{
			ID:              uuid.New(),
			ProjectID:       projectId,
			OrganizationID:  uuid.Nil, // Será preenchido pelo caller se necessário
			ShowPrepTime:    false,
			ShowRating:      false,
			ShowDescription: true,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		err = r.CreateSettings(newSettings)
		if err != nil {
			return nil, err
		}
		return newSettings, nil
	}

	// Se existe, reseta para padrões
	settings.ShowPrepTime = false
	settings.ShowRating = false
	settings.ShowDescription = true
	settings.UpdatedAt = time.Now()

	err = r.UpdateSettings(settings)
	if err != nil {
		return nil, err
	}
	return settings, nil
}
