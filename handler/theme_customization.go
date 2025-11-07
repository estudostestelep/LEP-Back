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
		// Se não encontra, retorna com padrões
		defaultTheme := &models.ThemeCustomization{
			ID:                  uuid.New(),
			ProjectID:           projectUUID,
			PrimaryColor:        "#3b82f6",
			SecondaryColor:      "#8b5cf6",
			BackgroundColor:     "#09090b",
			CardBackgroundColor: "#18181b",
			TextColor:           "#fafafa",
			TextSecondaryColor:  "#a1a1aa",
			AccentColor:         "#ec4899",
			IsActive:            false,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}
		return defaultTheme, nil
	}
	return theme, nil
}

// CreateOrUpdateTheme cria ou atualiza customização de tema
func (h *ThemeCustomizationHandler) CreateOrUpdateTheme(projectId string, organizationId string, theme *models.ThemeCustomization) (*models.ThemeCustomization, error) {
	projectUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}

	orgUUID, err := uuid.Parse(organizationId)
	if err != nil {
		return nil, err
	}

	// Validar cores HEX
	if err := validateHexColor(theme.PrimaryColor); err != nil {
		return nil, fmt.Errorf("invalid primary_color: %w", err)
	}
	if err := validateHexColor(theme.SecondaryColor); err != nil {
		return nil, fmt.Errorf("invalid secondary_color: %w", err)
	}
	if err := validateHexColor(theme.BackgroundColor); err != nil {
		return nil, fmt.Errorf("invalid background_color: %w", err)
	}
	if err := validateHexColor(theme.CardBackgroundColor); err != nil {
		return nil, fmt.Errorf("invalid card_background_color: %w", err)
	}
	if err := validateHexColor(theme.TextColor); err != nil {
		return nil, fmt.Errorf("invalid text_color: %w", err)
	}
	if err := validateHexColor(theme.TextSecondaryColor); err != nil {
		return nil, fmt.Errorf("invalid text_secondary_color: %w", err)
	}
	if err := validateHexColor(theme.AccentColor); err != nil {
		return nil, fmt.Errorf("invalid accent_color: %w", err)
	}

	// Buscar se já existe
	existingTheme, err := h.themeRepo.GetThemeByProject(projectUUID)
	if err == nil && existingTheme != nil {
		// Atualizar existente
		existingTheme.PrimaryColor = theme.PrimaryColor
		existingTheme.SecondaryColor = theme.SecondaryColor
		existingTheme.BackgroundColor = theme.BackgroundColor
		existingTheme.CardBackgroundColor = theme.CardBackgroundColor
		existingTheme.TextColor = theme.TextColor
		existingTheme.TextSecondaryColor = theme.TextSecondaryColor
		existingTheme.AccentColor = theme.AccentColor
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
		ID:                  uuid.New(),
		ProjectID:           projectUUID,
		OrganizationID:      orgUUID,
		PrimaryColor:        theme.PrimaryColor,
		SecondaryColor:      theme.SecondaryColor,
		BackgroundColor:     theme.BackgroundColor,
		CardBackgroundColor: theme.CardBackgroundColor,
		TextColor:           theme.TextColor,
		TextSecondaryColor:  theme.TextSecondaryColor,
		AccentColor:         theme.AccentColor,
		IsActive:            theme.IsActive,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	err = h.themeRepo.CreateTheme(newTheme)
	if err != nil {
		return nil, err
	}
	return newTheme, nil
}

// ResetToDefaults reseta tema para valores padrão
func (h *ThemeCustomizationHandler) ResetToDefaults(projectId string) (*models.ThemeCustomization, error) {
	projectUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}

	return h.themeRepo.ResetToDefaults(projectUUID)
}

// DeleteTheme deleta customização de tema
func (h *ThemeCustomizationHandler) DeleteTheme(projectId string) error {
	projectUUID, err := uuid.Parse(projectId)
	if err != nil {
		return err
	}

	return h.themeRepo.DeleteTheme(projectUUID)
}

// validateHexColor valida se a cor está em formato HEX válido
func validateHexColor(color string) error {
	// Aceita #RRGGBB ou #RRGGBBAA
	hexRegex := regexp.MustCompile("^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{8})$")
	if !hexRegex.MatchString(color) {
		return fmt.Errorf("color must be in hex format (#RRGGBB or #RRGGBBAA)")
	}
	return nil
}
