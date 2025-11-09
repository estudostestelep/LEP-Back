package handler

import (
	"fmt"
	"lep/repositories"
	"lep/repositories/models"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type ThemeCustomizationHandler struct {
	themeRepo repositories.IThemeCustomizationRepository
}

type IThemeCustomizationHandler interface {
	GetThemeByProject(projectId string) (*models.ThemeCustomization, error)
	CreateOrUpdateTheme(projectId string, organizationId string, theme *models.ThemeCustomization) (*models.ThemeCustomization, error)
	ResetToDefaults(projectId string) (*models.ThemeCustomization, error)
	DeleteTheme(projectId string) error
}

func NewThemeCustomizationHandler(themeRepo repositories.IThemeCustomizationRepository) IThemeCustomizationHandler {
	return &ThemeCustomizationHandler{themeRepo: themeRepo}
}

// GetThemeByProject busca customização de tema do projeto
func (h *ThemeCustomizationHandler) GetThemeByProject(projectId string) (*models.ThemeCustomization, error) {
	projectUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}

	theme, err := h.themeRepo.GetThemeByProject(projectUUID)
	if err != nil {
		// Se não encontra, retorna com padrões light/dark profissionais
		defaultTheme := h.buildDefaultTheme(projectUUID, uuid.UUID{})
		return defaultTheme, nil
	}
	return theme, nil
}

// CreateOrUpdateTheme cria ou atualiza customização de tema com suporte a light/dark variants
func (h *ThemeCustomizationHandler) CreateOrUpdateTheme(projectId string, organizationId string, theme *models.ThemeCustomization) (*models.ThemeCustomization, error) {
	projectUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}

	orgUUID, err := uuid.Parse(organizationId)
	if err != nil {
		return nil, err
	}

	// Validar todas as cores light/dark (11 cores principais + 5 semânticas + 2 sistema = 18 campos de cor)
	if err := h.validateAllColors(theme); err != nil {
		return nil, err
	}

	// Validar configurações numéricas
	if err := validateOpacity(theme.DisabledOpacity, "disabled_opacity"); err != nil {
		return nil, err
	}
	if err := validateShadowIntensity(theme.ShadowIntensity, "shadow_intensity"); err != nil {
		return nil, err
	}

	// Buscar se já existe
	existingTheme, err := h.themeRepo.GetThemeByProject(projectUUID)
	if err == nil && existingTheme != nil {
		// Atualizar existente com suporte a atualizações parciais
		existingTheme = h.mergeThemeUpdates(existingTheme, theme)
		existingTheme.IsActive = theme.IsActive
		existingTheme.UpdatedAt = time.Now()

		err = h.themeRepo.UpdateTheme(existingTheme)
		if err != nil {
			return nil, err
		}
		return existingTheme, nil
	}

	// Criar novo
	newTheme := &models.ThemeCustomization{
		ID:             uuid.New(),
		ProjectID:      projectUUID,
		OrganizationID: orgUUID,

		// Light Mode - 11 cores principais
		PrimaryColorLight:        theme.PrimaryColorLight,
		SecondaryColorLight:      theme.SecondaryColorLight,
		BackgroundColorLight:     theme.BackgroundColorLight,
		CardBackgroundColorLight: theme.CardBackgroundColorLight,
		TextColorLight:           theme.TextColorLight,
		TextSecondaryColorLight:  theme.TextSecondaryColorLight,
		AccentColorLight:         theme.AccentColorLight,

		// Dark Mode - 11 cores principais
		PrimaryColorDark:        theme.PrimaryColorDark,
		SecondaryColorDark:      theme.SecondaryColorDark,
		BackgroundColorDark:     theme.BackgroundColorDark,
		CardBackgroundColorDark: theme.CardBackgroundColorDark,
		TextColorDark:           theme.TextColorDark,
		TextSecondaryColorDark:  theme.TextSecondaryColorDark,
		AccentColorDark:         theme.AccentColorDark,

		// Light Mode - 5 cores semânticas
		DestructiveColorLight: theme.DestructiveColorLight,
		SuccessColorLight:     theme.SuccessColorLight,
		WarningColorLight:     theme.WarningColorLight,
		BorderColorLight:      theme.BorderColorLight,
		PriceColorLight:       theme.PriceColorLight,

		// Dark Mode - 5 cores semânticas
		DestructiveColorDark: theme.DestructiveColorDark,
		SuccessColorDark:     theme.SuccessColorDark,
		WarningColorDark:     theme.WarningColorDark,
		BorderColorDark:      theme.BorderColorDark,
		PriceColorDark:       theme.PriceColorDark,

		// Light Mode - 2 cores sistema
		FocusRingColorLight:      theme.FocusRingColorLight,
		InputBackgroundColorLight: theme.InputBackgroundColorLight,

		// Dark Mode - 2 cores sistema
		FocusRingColorDark:       theme.FocusRingColorDark,
		InputBackgroundColorDark: theme.InputBackgroundColorDark,

		// Configurações numéricas
		DisabledOpacity: theme.DisabledOpacity,
		ShadowIntensity: theme.ShadowIntensity,

		IsActive:  theme.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = h.themeRepo.CreateTheme(newTheme)
	if err != nil {
		return nil, err
	}
	return newTheme, nil
}

// ResetToDefaults apaga todas as cores customizadas (seta para nil)
// deixando o frontend lidar com a cor padrão
func (h *ThemeCustomizationHandler) ResetToDefaults(projectId string) (*models.ThemeCustomization, error) {
	projectUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}

	// Delega para o repositório que vai zerar todas as cores
	theme, err := h.themeRepo.ResetToDefaults(projectUUID)
	if err != nil {
		return nil, err
	}

	return theme, nil
}

// DeleteTheme deleta customização de tema
func (h *ThemeCustomizationHandler) DeleteTheme(projectId string) error {
	projectUUID, err := uuid.Parse(projectId)
	if err != nil {
		return err
	}

	return h.themeRepo.DeleteTheme(projectUUID)
}

// ==================== Helper Functions ====================

// buildDefaultTheme constrói tema padrão com cores profissionais light/dark
func (h *ThemeCustomizationHandler) buildDefaultTheme(projectID uuid.UUID, orgID uuid.UUID) *models.ThemeCustomization {
	// LIGHT MODE - Cores profissionais claras
	lightPrimary := "#1E293B"
	lightSecondary := "#8B5CF6"
	lightBackground := "#FFFFFF"
	lightCardBackground := "#FFFFFF"
	lightText := "#0F172A"
	lightTextSecondary := "#64748B"
	lightAccent := "#EC4899"

	lightDestructive := "#EF4444"
	lightSuccess := "#10B981"
	lightWarning := "#F59E0B"
	lightBorder := "#E5E7EB"
	lightPrice := "#10B981"

	lightFocusRing := "#3B82F6"
	lightInputBackground := "#F3F4F6"

	// DARK MODE - Cores profissionais escuras
	darkPrimary := "#F8FAFC"
	darkSecondary := "#A78BFA"
	darkBackground := "#0F172A"
	darkCardBackground := "#1E293B"
	darkText := "#F8FAFC"
	darkTextSecondary := "#94A3B8"
	darkAccent := "#F472B6"

	darkDestructive := "#DC2626"
	darkSuccess := "#34D399"
	darkWarning := "#FBBF24"
	darkBorder := "#475569"
	darkPrice := "#34D399"

	darkFocusRing := "#93C5FD"
	darkInputBackground := "#1F2937"

	// Configurações numéricas
	disabledOpacity := 0.5
	shadowIntensity := 1.0

	return &models.ThemeCustomization{
		ID:             uuid.New(),
		ProjectID:      projectID,
		OrganizationID: orgID,

		// Light Mode - 11 cores principais
		PrimaryColorLight:        &lightPrimary,
		SecondaryColorLight:      &lightSecondary,
		BackgroundColorLight:     &lightBackground,
		CardBackgroundColorLight: &lightCardBackground,
		TextColorLight:           &lightText,
		TextSecondaryColorLight:  &lightTextSecondary,
		AccentColorLight:         &lightAccent,

		// Dark Mode - 11 cores principais
		PrimaryColorDark:        &darkPrimary,
		SecondaryColorDark:      &darkSecondary,
		BackgroundColorDark:     &darkBackground,
		CardBackgroundColorDark: &darkCardBackground,
		TextColorDark:           &darkText,
		TextSecondaryColorDark:  &darkTextSecondary,
		AccentColorDark:         &darkAccent,

		// Light Mode - 5 cores semânticas
		DestructiveColorLight: &lightDestructive,
		SuccessColorLight:     &lightSuccess,
		WarningColorLight:     &lightWarning,
		BorderColorLight:      &lightBorder,
		PriceColorLight:       &lightPrice,

		// Dark Mode - 5 cores semânticas
		DestructiveColorDark: &darkDestructive,
		SuccessColorDark:     &darkSuccess,
		WarningColorDark:     &darkWarning,
		BorderColorDark:      &darkBorder,
		PriceColorDark:       &darkPrice,

		// Light Mode - 2 cores sistema
		FocusRingColorLight:      &lightFocusRing,
		InputBackgroundColorLight: &lightInputBackground,

		// Dark Mode - 2 cores sistema
		FocusRingColorDark:       &darkFocusRing,
		InputBackgroundColorDark: &darkInputBackground,

		// Configurações numéricas
		DisabledOpacity: &disabledOpacity,
		ShadowIntensity: &shadowIntensity,

		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// mergeThemeUpdates mescla atualizações parciais mantendo campos existentes
func (h *ThemeCustomizationHandler) mergeThemeUpdates(existing *models.ThemeCustomization, updates *models.ThemeCustomization) *models.ThemeCustomization {
	// Light Mode - 11 cores principais
	if updates.PrimaryColorLight != nil {
		existing.PrimaryColorLight = updates.PrimaryColorLight
	}
	if updates.SecondaryColorLight != nil {
		existing.SecondaryColorLight = updates.SecondaryColorLight
	}
	if updates.BackgroundColorLight != nil {
		existing.BackgroundColorLight = updates.BackgroundColorLight
	}
	if updates.CardBackgroundColorLight != nil {
		existing.CardBackgroundColorLight = updates.CardBackgroundColorLight
	}
	if updates.TextColorLight != nil {
		existing.TextColorLight = updates.TextColorLight
	}
	if updates.TextSecondaryColorLight != nil {
		existing.TextSecondaryColorLight = updates.TextSecondaryColorLight
	}
	if updates.AccentColorLight != nil {
		existing.AccentColorLight = updates.AccentColorLight
	}

	// Dark Mode - 11 cores principais
	if updates.PrimaryColorDark != nil {
		existing.PrimaryColorDark = updates.PrimaryColorDark
	}
	if updates.SecondaryColorDark != nil {
		existing.SecondaryColorDark = updates.SecondaryColorDark
	}
	if updates.BackgroundColorDark != nil {
		existing.BackgroundColorDark = updates.BackgroundColorDark
	}
	if updates.CardBackgroundColorDark != nil {
		existing.CardBackgroundColorDark = updates.CardBackgroundColorDark
	}
	if updates.TextColorDark != nil {
		existing.TextColorDark = updates.TextColorDark
	}
	if updates.TextSecondaryColorDark != nil {
		existing.TextSecondaryColorDark = updates.TextSecondaryColorDark
	}
	if updates.AccentColorDark != nil {
		existing.AccentColorDark = updates.AccentColorDark
	}

	// Light Mode - 5 cores semânticas
	if updates.DestructiveColorLight != nil {
		existing.DestructiveColorLight = updates.DestructiveColorLight
	}
	if updates.SuccessColorLight != nil {
		existing.SuccessColorLight = updates.SuccessColorLight
	}
	if updates.WarningColorLight != nil {
		existing.WarningColorLight = updates.WarningColorLight
	}
	if updates.BorderColorLight != nil {
		existing.BorderColorLight = updates.BorderColorLight
	}
	if updates.PriceColorLight != nil {
		existing.PriceColorLight = updates.PriceColorLight
	}

	// Dark Mode - 5 cores semânticas
	if updates.DestructiveColorDark != nil {
		existing.DestructiveColorDark = updates.DestructiveColorDark
	}
	if updates.SuccessColorDark != nil {
		existing.SuccessColorDark = updates.SuccessColorDark
	}
	if updates.WarningColorDark != nil {
		existing.WarningColorDark = updates.WarningColorDark
	}
	if updates.BorderColorDark != nil {
		existing.BorderColorDark = updates.BorderColorDark
	}
	if updates.PriceColorDark != nil {
		existing.PriceColorDark = updates.PriceColorDark
	}

	// Light Mode - 2 cores sistema
	if updates.FocusRingColorLight != nil {
		existing.FocusRingColorLight = updates.FocusRingColorLight
	}
	if updates.InputBackgroundColorLight != nil {
		existing.InputBackgroundColorLight = updates.InputBackgroundColorLight
	}

	// Dark Mode - 2 cores sistema
	if updates.FocusRingColorDark != nil {
		existing.FocusRingColorDark = updates.FocusRingColorDark
	}
	if updates.InputBackgroundColorDark != nil {
		existing.InputBackgroundColorDark = updates.InputBackgroundColorDark
	}

	// Configurações numéricas
	if updates.DisabledOpacity != nil {
		existing.DisabledOpacity = updates.DisabledOpacity
	}
	if updates.ShadowIntensity != nil {
		existing.ShadowIntensity = updates.ShadowIntensity
	}

	return existing
}

// validateAllColors valida todas as 22 cores HEX (light + dark)
func (h *ThemeCustomizationHandler) validateAllColors(theme *models.ThemeCustomization) error {
	// Light Mode - 11 cores principais
	if err := validateHexColorPointer(theme.PrimaryColorLight, "primary_color_light"); err != nil {
		return err
	}
	if err := validateHexColorPointer(theme.SecondaryColorLight, "secondary_color_light"); err != nil {
		return err
	}
	if err := validateHexColorPointer(theme.BackgroundColorLight, "background_color_light"); err != nil {
		return err
	}
	if err := validateHexColorPointer(theme.CardBackgroundColorLight, "card_background_color_light"); err != nil {
		return err
	}
	if err := validateHexColorPointer(theme.TextColorLight, "text_color_light"); err != nil {
		return err
	}
	if err := validateHexColorPointer(theme.TextSecondaryColorLight, "text_secondary_color_light"); err != nil {
		return err
	}
	if err := validateHexColorPointer(theme.AccentColorLight, "accent_color_light"); err != nil {
		return err
	}

	// Dark Mode - 11 cores principais
	if err := validateHexColorPointer(theme.PrimaryColorDark, "primary_color_dark"); err != nil {
		return err
	}
	if err := validateHexColorPointer(theme.SecondaryColorDark, "secondary_color_dark"); err != nil {
		return err
	}
	if err := validateHexColorPointer(theme.BackgroundColorDark, "background_color_dark"); err != nil {
		return err
	}
	if err := validateHexColorPointer(theme.CardBackgroundColorDark, "card_background_color_dark"); err != nil {
		return err
	}
	if err := validateHexColorPointer(theme.TextColorDark, "text_color_dark"); err != nil {
		return err
	}
	if err := validateHexColorPointer(theme.TextSecondaryColorDark, "text_secondary_color_dark"); err != nil {
		return err
	}
	if err := validateHexColorPointer(theme.AccentColorDark, "accent_color_dark"); err != nil {
		return err
	}

	// Light Mode - 5 cores semânticas
	if err := validateHexColorPointer(theme.DestructiveColorLight, "destructive_color_light"); err != nil {
		return err
	}
	if err := validateHexColorPointer(theme.SuccessColorLight, "success_color_light"); err != nil {
		return err
	}
	if err := validateHexColorPointer(theme.WarningColorLight, "warning_color_light"); err != nil {
		return err
	}
	if err := validateHexColorPointer(theme.BorderColorLight, "border_color_light"); err != nil {
		return err
	}
	if err := validateHexColorPointer(theme.PriceColorLight, "price_color_light"); err != nil {
		return err
	}

	// Dark Mode - 5 cores semânticas
	if err := validateHexColorPointer(theme.DestructiveColorDark, "destructive_color_dark"); err != nil {
		return err
	}
	if err := validateHexColorPointer(theme.SuccessColorDark, "success_color_dark"); err != nil {
		return err
	}
	if err := validateHexColorPointer(theme.WarningColorDark, "warning_color_dark"); err != nil {
		return err
	}
	if err := validateHexColorPointer(theme.BorderColorDark, "border_color_dark"); err != nil {
		return err
	}
	if err := validateHexColorPointer(theme.PriceColorDark, "price_color_dark"); err != nil {
		return err
	}

	// Light Mode - 2 cores sistema
	if err := validateHexColorPointer(theme.FocusRingColorLight, "focus_ring_color_light"); err != nil {
		return err
	}
	if err := validateHexColorPointer(theme.InputBackgroundColorLight, "input_background_color_light"); err != nil {
		return err
	}

	// Dark Mode - 2 cores sistema
	if err := validateHexColorPointer(theme.FocusRingColorDark, "focus_ring_color_dark"); err != nil {
		return err
	}
	if err := validateHexColorPointer(theme.InputBackgroundColorDark, "input_background_color_dark"); err != nil {
		return err
	}

	return nil
}

// ==================== Validation Functions ====================

// validateHexColor valida se a cor está em formato HEX válido
func validateHexColor(color string) error {
	// Aceita #RRGGBB ou #RRGGBBAA
	hexRegex := regexp.MustCompile("^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{8})$")
	if !hexRegex.MatchString(color) {
		return fmt.Errorf("color must be in hex format (#RRGGBB or #RRGGBBAA)")
	}
	return nil
}

// validateHexColorPointer valida cores em pointers (para campos opcionais)
func validateHexColorPointer(color *string, fieldName string) error {
	if color == nil || *color == "" {
		// Null/empty é permitido para campos opcionais
		return nil
	}
	if err := validateHexColor(*color); err != nil {
		return fmt.Errorf("invalid %s: %w", fieldName, err)
	}
	return nil
}

// validateOpacity valida se o valor está entre 0.0 e 1.0
func validateOpacity(value *float64, fieldName string) error {
	if value == nil {
		// Null é permitido para campos opcionais
		return nil
	}
	if *value < 0.0 || *value > 1.0 {
		return fmt.Errorf("invalid %s: must be between 0.0 and 1.0, got %.2f", fieldName, *value)
	}
	return nil
}

// validateShadowIntensity valida se o valor está entre 0.0 e 2.0
func validateShadowIntensity(value *float64, fieldName string) error {
	if value == nil {
		// Null é permitido para campos opcionais
		return nil
	}
	if *value < 0.0 || *value > 2.0 {
		return fmt.Errorf("invalid %s: must be between 0.0 and 2.0, got %.2f", fieldName, *value)
	}
	return nil
}
