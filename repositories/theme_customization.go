package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ThemeCustomizationRepository struct {
	db *gorm.DB
}

type IThemeCustomizationRepository interface {
	GetThemeByProject(projectId uuid.UUID) (*models.ThemeCustomization, error)
	CreateTheme(theme *models.ThemeCustomization) error
	UpdateTheme(theme *models.ThemeCustomization) error
	DeleteTheme(projectId uuid.UUID) error
	ResetToDefaults(projectId uuid.UUID) (*models.ThemeCustomization, error)
}

func NewThemeCustomizationRepository(db *gorm.DB) IThemeCustomizationRepository {
	return &ThemeCustomizationRepository{db: db}
}

// GetThemeByProject busca customização de tema por projeto
func (r *ThemeCustomizationRepository) GetThemeByProject(projectId uuid.UUID) (*models.ThemeCustomization, error) {
	var theme models.ThemeCustomization
	err := r.db.Where("project_id = ?", projectId).First(&theme).Error
	if err != nil {
		return nil, err
	}
	return &theme, nil
}

// CreateTheme cria nova customização de tema
func (r *ThemeCustomizationRepository) CreateTheme(theme *models.ThemeCustomization) error {
	return r.db.Create(theme).Error
}

// UpdateTheme atualiza customização de tema existente
func (r *ThemeCustomizationRepository) UpdateTheme(theme *models.ThemeCustomization) error {
	if theme.ID == uuid.Nil {
		return gorm.ErrInvalidData
	}
	theme.UpdatedAt = time.Now()
	return r.db.Model(theme).Where("id = ?", theme.ID).Updates(theme).Error
}

// DeleteTheme deleta customização de tema
func (r *ThemeCustomizationRepository) DeleteTheme(projectId uuid.UUID) error {
	return r.db.Where("project_id = ?", projectId).Delete(&models.ThemeCustomization{}).Error
}

// ResetToDefaults apaga todas as cores customizadas (seta para nil)
// deixando o frontend lidar com a cor padrão
func (r *ThemeCustomizationRepository) ResetToDefaults(projectId uuid.UUID) (*models.ThemeCustomization, error) {
	theme, err := r.GetThemeByProject(projectId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Se não existe, cria com todos os campos nil (padrão do frontend)
	if err == gorm.ErrRecordNotFound {
		newTheme := &models.ThemeCustomization{
			ID:             uuid.New(),
			ProjectID:      projectId,
			OrganizationID: uuid.Nil,
			// Todos os campos de cor setados para nil
			PrimaryColorLight:         nil,
			SecondaryColorLight:       nil,
			BackgroundColorLight:      nil,
			CardBackgroundColorLight:  nil,
			TextColorLight:            nil,
			TextSecondaryColorLight:   nil,
			AccentColorLight:          nil,
			DestructiveColorLight:     nil,
			SuccessColorLight:         nil,
			WarningColorLight:         nil,
			BorderColorLight:          nil,
			PriceColorLight:           nil,
			FocusRingColorLight:       nil,
			InputBackgroundColorLight: nil,
			// Dark mode
			PrimaryColorDark:         nil,
			SecondaryColorDark:       nil,
			BackgroundColorDark:      nil,
			CardBackgroundColorDark:  nil,
			TextColorDark:            nil,
			TextSecondaryColorDark:   nil,
			AccentColorDark:          nil,
			DestructiveColorDark:     nil,
			SuccessColorDark:         nil,
			WarningColorDark:         nil,
			BorderColorDark:          nil,
			PriceColorDark:           nil,
			FocusRingColorDark:       nil,
			InputBackgroundColorDark: nil,
			// Configs também zerados
			DisabledOpacity: nil,
			ShadowIntensity: nil,
			IsActive:        false,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		if err := r.CreateTheme(newTheme); err != nil {
			return nil, err
		}
		return newTheme, nil
	}

	// Se já existe, zera todas as cores customizadas
	theme.PrimaryColorLight = nil
	theme.SecondaryColorLight = nil
	theme.BackgroundColorLight = nil
	theme.CardBackgroundColorLight = nil
	theme.TextColorLight = nil
	theme.TextSecondaryColorLight = nil
	theme.AccentColorLight = nil

	theme.PrimaryColorDark = nil
	theme.SecondaryColorDark = nil
	theme.BackgroundColorDark = nil
	theme.CardBackgroundColorDark = nil
	theme.TextColorDark = nil
	theme.TextSecondaryColorDark = nil
	theme.AccentColorDark = nil

	theme.DestructiveColorLight = nil
	theme.SuccessColorLight = nil
	theme.WarningColorLight = nil
	theme.BorderColorLight = nil
	theme.PriceColorLight = nil

	theme.DestructiveColorDark = nil
	theme.SuccessColorDark = nil
	theme.WarningColorDark = nil
	theme.BorderColorDark = nil
	theme.PriceColorDark = nil

	theme.FocusRingColorLight = nil
	theme.InputBackgroundColorLight = nil
	theme.FocusRingColorDark = nil
	theme.InputBackgroundColorDark = nil
	theme.DisabledOpacity = nil
	theme.ShadowIntensity = nil

	theme.IsActive = false
	theme.UpdatedAt = time.Now()

	if err := r.UpdateTheme(theme); err != nil {
		return nil, err
	}

	return theme, nil
}
