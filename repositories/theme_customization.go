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
		// Criar values para os campos pointer
		destructive := "#EF4444"
		success := "#10B981"
		warning := "#F59E0B"
		border := "#E5E7EB"
		disabledOpacity := 0.50
		focusRing := "#3B82F6"
		inputBg := "#FFFFFF"
		shadowIntensity := 1.00

		newTheme := &models.ThemeCustomization{
			ID:                  uuid.New(),
			ProjectID:           projectId,
			PrimaryColor:        "#0F172A",
			SecondaryColor:      "#1E293B",
			BackgroundColor:     "#09090b",
			CardBackgroundColor: "#18181b",
			TextColor:           "#fafafa",
			TextSecondaryColor:  "#a1a1aa",
			AccentColor:         "#ec4899",

			// Cores semânticas
			DestructiveColor: &destructive,
			SuccessColor:     &success,
			WarningColor:     &warning,
			BorderColor:      &border,

			// Configurações do sistema
			DisabledOpacity:      &disabledOpacity,
			FocusRingColor:       &focusRing,
			InputBackgroundColor: &inputBg,
			ShadowIntensity:      &shadowIntensity,

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
	theme.PrimaryColor = "#0F172A"
	theme.SecondaryColor = "#1E293B"
	theme.BackgroundColor = "#09090b"
	theme.CardBackgroundColor = "#18181b"
	theme.TextColor = "#fafafa"
	theme.TextSecondaryColor = "#a1a1aa"
	theme.AccentColor = "#ec4899"

	// Resetar cores semânticas
	destructive := "#EF4444"
	success := "#10B981"
	warning := "#F59E0B"
	border := "#E5E7EB"
	theme.DestructiveColor = &destructive
	theme.SuccessColor = &success
	theme.WarningColor = &warning
	theme.BorderColor = &border

	// Resetar configurações do sistema
	disabledOpacity := 0.50
	focusRing := "#3B82F6"
	inputBg := "#FFFFFF"
	shadowIntensity := 1.00
	theme.DisabledOpacity = &disabledOpacity
	theme.FocusRingColor = &focusRing
	theme.InputBackgroundColor = &inputBg
	theme.ShadowIntensity = &shadowIntensity

	theme.IsActive = false
	theme.UpdatedAt = time.Now()

	err = r.UpdateTheme(theme)
	if err != nil {
		return nil, err
	}
	return theme, nil
}
