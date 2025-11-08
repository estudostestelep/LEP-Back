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

// ResetToDefaults reseta tema para valores padrão
func (r *ThemeCustomizationRepository) ResetToDefaults(projectId uuid.UUID) (*models.ThemeCustomization, error) {
	// Busca tema existente
	theme, err := r.GetThemeByProject(projectId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Se não existe, cria com padrões
	if err == gorm.ErrRecordNotFound {
		// Light Mode defaults
		primaryLight := "#1E293B"
		secondaryLight := "#8B5CF6"
		backgroundLight := "#FFFFFF"
		cardBackgroundLight := "#FFFFFF"
		textLight := "#0F172A"
		textSecondaryLight := "#64748B"
		accentLight := "#EC4899"

		// Dark Mode defaults
		primaryDark := "#F8FAFC"
		secondaryDark := "#A78BFA"
		backgroundDark := "#0F172A"
		cardBackgroundDark := "#1E293B"
		textDark := "#F8FAFC"
		textSecondaryDark := "#94A3B8"
		accentDark := "#F472B6"

		// Light Mode semantic colors
		destructiveLight := "#EF4444"
		successLight := "#10B981"
		warningLight := "#F59E0B"
		borderLight := "#E5E7EB"
		priceLight := "#10B981"

		// Dark Mode semantic colors
		destructiveDark := "#DC2626"
		successDark := "#34D399"
		warningDark := "#FBBF24"
		borderDark := "#475569"
		priceDark := "#34D399"

		// Light Mode system
		focusRingLight := "#3B82F6"
		inputBgLight := "#F3F4F6"

		// Dark Mode system
		focusRingDark := "#93C5FD"
		inputBgDark := "#1F2937"

		// Numeric configs
		disabledOpacity := 0.50
		shadowIntensity := 1.00

		newTheme := &models.ThemeCustomization{
			ID:             uuid.New(),
			ProjectID:      projectId,
			OrganizationID: uuid.Nil, // Será preenchido pelo handler

			// Light Mode cores principais
			PrimaryColorLight:        &primaryLight,
			SecondaryColorLight:      &secondaryLight,
			BackgroundColorLight:     &backgroundLight,
			CardBackgroundColorLight: &cardBackgroundLight,
			TextColorLight:           &textLight,
			TextSecondaryColorLight:  &textSecondaryLight,
			AccentColorLight:         &accentLight,

			// Dark Mode cores principais
			PrimaryColorDark:        &primaryDark,
			SecondaryColorDark:      &secondaryDark,
			BackgroundColorDark:     &backgroundDark,
			CardBackgroundColorDark: &cardBackgroundDark,
			TextColorDark:           &textDark,
			TextSecondaryColorDark:  &textSecondaryDark,
			AccentColorDark:         &accentDark,

			// Light Mode semantic colors
			DestructiveColorLight: &destructiveLight,
			SuccessColorLight:     &successLight,
			WarningColorLight:     &warningLight,
			BorderColorLight:      &borderLight,
			PriceColorLight:       &priceLight,

			// Dark Mode semantic colors
			DestructiveColorDark: &destructiveDark,
			SuccessColorDark:     &successDark,
			WarningColorDark:     &warningDark,
			BorderColorDark:      &borderDark,
			PriceColorDark:       &priceDark,

			// Light Mode system
			FocusRingColorLight:      &focusRingLight,
			InputBackgroundColorLight: &inputBgLight,

			// Dark Mode system
			FocusRingColorDark:       &focusRingDark,
			InputBackgroundColorDark: &inputBgDark,

			// Numeric configurations
			DisabledOpacity: &disabledOpacity,
			ShadowIntensity: &shadowIntensity,

			IsActive:  false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = r.CreateTheme(newTheme)
		if err != nil {
			return nil, err
		}
		return newTheme, nil
	}

	// Se existe, reseta para padrões
	// Light Mode cores principais
	primaryLight := "#1E293B"
	secondaryLight := "#8B5CF6"
	backgroundLight := "#FFFFFF"
	cardBackgroundLight := "#FFFFFF"
	textLight := "#0F172A"
	textSecondaryLight := "#64748B"
	accentLight := "#EC4899"

	theme.PrimaryColorLight = &primaryLight
	theme.SecondaryColorLight = &secondaryLight
	theme.BackgroundColorLight = &backgroundLight
	theme.CardBackgroundColorLight = &cardBackgroundLight
	theme.TextColorLight = &textLight
	theme.TextSecondaryColorLight = &textSecondaryLight
	theme.AccentColorLight = &accentLight

	// Dark Mode cores principais
	primaryDark := "#F8FAFC"
	secondaryDark := "#A78BFA"
	backgroundDark := "#0F172A"
	cardBackgroundDark := "#1E293B"
	textDark := "#F8FAFC"
	textSecondaryDark := "#94A3B8"
	accentDark := "#F472B6"

	theme.PrimaryColorDark = &primaryDark
	theme.SecondaryColorDark = &secondaryDark
	theme.BackgroundColorDark = &backgroundDark
	theme.CardBackgroundColorDark = &cardBackgroundDark
	theme.TextColorDark = &textDark
	theme.TextSecondaryColorDark = &textSecondaryDark
	theme.AccentColorDark = &accentDark

	// Resetar Light Mode semantic colors
	destructiveLight := "#EF4444"
	successLight := "#10B981"
	warningLight := "#F59E0B"
	borderLight := "#E5E7EB"
	priceLight := "#10B981"

	theme.DestructiveColorLight = &destructiveLight
	theme.SuccessColorLight = &successLight
	theme.WarningColorLight = &warningLight
	theme.BorderColorLight = &borderLight
	theme.PriceColorLight = &priceLight

	// Resetar Dark Mode semantic colors
	destructiveDark := "#DC2626"
	successDark := "#34D399"
	warningDark := "#FBBF24"
	borderDark := "#475569"
	priceDark := "#34D399"

	theme.DestructiveColorDark = &destructiveDark
	theme.SuccessColorDark = &successDark
	theme.WarningColorDark = &warningDark
	theme.BorderColorDark = &borderDark
	theme.PriceColorDark = &priceDark

	// Resetar Light Mode system
	focusRingLight := "#3B82F6"
	inputBgLight := "#F3F4F6"

	theme.FocusRingColorLight = &focusRingLight
	theme.InputBackgroundColorLight = &inputBgLight

	// Resetar Dark Mode system
	focusRingDark := "#93C5FD"
	inputBgDark := "#1F2937"

	theme.FocusRingColorDark = &focusRingDark
	theme.InputBackgroundColorDark = &inputBgDark

	// Resetar configurações numéricas
	disabledOpacity := 0.50
	shadowIntensity := 1.00
	theme.DisabledOpacity = &disabledOpacity
	theme.ShadowIntensity = &shadowIntensity

	theme.IsActive = false
	theme.UpdatedAt = time.Now()

	err = r.UpdateTheme(theme)
	if err != nil {
		return nil, err
	}
	return theme, nil
}
