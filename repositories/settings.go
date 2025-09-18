package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SettingsRepository struct {
	db *gorm.DB
}

type ISettingsRepository interface {
	GetSettingsByProject(orgId, projectId uuid.UUID) (*models.Settings, error)
	CreateSettings(settings *models.Settings) error
	UpdateSettings(settings *models.Settings) error
	GetOrCreateSettings(orgId, projectId uuid.UUID) (*models.Settings, error)
}

func NewSettingsRepository(db *gorm.DB) ISettingsRepository {
	return &SettingsRepository{db: db}
}

// GetSettingsByProject busca configurações por projeto
func (r *SettingsRepository) GetSettingsByProject(orgId, projectId uuid.UUID) (*models.Settings, error) {
	var settings models.Settings
	err := r.db.Where("organization_id = ? AND project_id = ?", orgId, projectId).First(&settings).Error
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

// CreateSettings cria novas configurações
func (r *SettingsRepository) CreateSettings(settings *models.Settings) error {
	return r.db.Create(settings).Error
}

// UpdateSettings atualiza configurações existentes
func (r *SettingsRepository) UpdateSettings(settings *models.Settings) error {
	settings.UpdatedAt = time.Now()
	return r.db.Save(settings).Error
}

// GetOrCreateSettings busca ou cria configurações padrão para o projeto
func (r *SettingsRepository) GetOrCreateSettings(orgId, projectId uuid.UUID) (*models.Settings, error) {
	// Tenta buscar configurações existentes
	settings, err := r.GetSettingsByProject(orgId, projectId)
	if err == nil {
		return settings, nil
	}

	// Se não encontrou, cria configurações padrão
	if err == gorm.ErrRecordNotFound {
		defaultSettings := &models.Settings{
			Id:             uuid.New(),
			OrganizationId: orgId,
			ProjectId:      projectId,
			MinAdvanceHours: 2,
			MaxAdvanceDays:  30,
			NotifyReservationCreate: true,
			NotifyReservationUpdate: true,
			NotifyReservationCancel: true,
			NotifyTableAvailable:    true,
			NotifyConfirmation24h:   true,
			DefaultNotificationChannel: "sms",
			EnableSms:      true,
			EnableEmail:    false,
			EnableWhatsapp: false,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		err = r.CreateSettings(defaultSettings)
		if err != nil {
			return nil, err
		}
		return defaultSettings, nil
	}

	return nil, err
}