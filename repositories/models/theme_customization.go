package models

import (
	"time"

	"github.com/google/uuid"
)

// ThemeCustomization armazena as cores personalizadas do tema para cada projeto
// Suporta customização separada para Light Mode e Dark Mode
type ThemeCustomization struct {
	ID             uuid.UUID `gorm:"primaryKey" json:"id"`
	ProjectID      uuid.UUID `json:"project_id" gorm:"not null;uniqueIndex"`
	OrganizationID uuid.UUID `json:"organization_id" gorm:"not null"`

	// ==================== CORES PRINCIPAIS - LIGHT MODE (11 campos) ====================
	PrimaryColorLight        *string `json:"primary_color_light" gorm:"default:null"`          // Light: #1E293B
	SecondaryColorLight      *string `json:"secondary_color_light" gorm:"default:null"`        // Light: #8B5CF6
	BackgroundColorLight     *string `json:"background_color_light" gorm:"default:null"`       // Light: #FFFFFF
	CardBackgroundColorLight *string `json:"card_background_color_light" gorm:"default:null"` // Light: #FFFFFF
	TextColorLight           *string `json:"text_color_light" gorm:"default:null"`             // Light: #0F172A
	TextSecondaryColorLight  *string `json:"text_secondary_color_light" gorm:"default:null"`  // Light: #64748B
	AccentColorLight         *string `json:"accent_color_light" gorm:"default:null"`           // Light: #EC4899

	// ==================== CORES PRINCIPAIS - DARK MODE (11 campos) ====================
	PrimaryColorDark        *string `json:"primary_color_dark" gorm:"default:null"`          // Dark: #F8FAFC
	SecondaryColorDark      *string `json:"secondary_color_dark" gorm:"default:null"`        // Dark: #A78BFA
	BackgroundColorDark     *string `json:"background_color_dark" gorm:"default:null"`       // Dark: #0F172A
	CardBackgroundColorDark *string `json:"card_background_color_dark" gorm:"default:null"` // Dark: #1E293B
	TextColorDark           *string `json:"text_color_dark" gorm:"default:null"`             // Dark: #F8FAFC
	TextSecondaryColorDark  *string `json:"text_secondary_color_dark" gorm:"default:null"`  // Dark: #94A3B8
	AccentColorDark         *string `json:"accent_color_dark" gorm:"default:null"`           // Dark: #F472B6

	// ==================== CORES SEMÂNTICAS - LIGHT MODE (5 campos) ====================
	DestructiveColorLight *string `json:"destructive_color_light" gorm:"default:null"` // Light: #EF4444
	SuccessColorLight     *string `json:"success_color_light" gorm:"default:null"`     // Light: #10B981
	WarningColorLight     *string `json:"warning_color_light" gorm:"default:null"`     // Light: #F59E0B
	BorderColorLight      *string `json:"border_color_light" gorm:"default:null"`      // Light: #E5E7EB
	PriceColorLight       *string `json:"price_color_light" gorm:"default:null"`       // Light: #10B981

	// ==================== CORES SEMÂNTICAS - DARK MODE (5 campos) ====================
	DestructiveColorDark *string `json:"destructive_color_dark" gorm:"default:null"` // Dark: #DC2626
	SuccessColorDark     *string `json:"success_color_dark" gorm:"default:null"`     // Dark: #34D399
	WarningColorDark     *string `json:"warning_color_dark" gorm:"default:null"`     // Dark: #FBBF24
	BorderColorDark      *string `json:"border_color_dark" gorm:"default:null"`      // Dark: #475569
	PriceColorDark       *string `json:"price_color_dark" gorm:"default:null"`       // Dark: #34D399

	// ==================== SISTEMA - LIGHT MODE (2 campos) ====================
	FocusRingColorLight      *string `json:"focus_ring_color_light" gorm:"default:null"`       // Light: #3B82F6
	InputBackgroundColorLight *string `json:"input_background_color_light" gorm:"default:null"` // Light: #F3F4F6

	// ==================== SISTEMA - DARK MODE (2 campos) ====================
	FocusRingColorDark       *string `json:"focus_ring_color_dark" gorm:"default:null"`       // Dark: #93C5FD
	InputBackgroundColorDark *string `json:"input_background_color_dark" gorm:"default:null"` // Dark: #1F2937

	// ==================== CONFIGURAÇÕES NUMÉRICAS (2 campos) ====================
	DisabledOpacity *float64 `json:"disabled_opacity" gorm:"default:null"` // Opacidade para estados desabilitados (0.0-1.0, padrão 0.5)
	ShadowIntensity *float64 `json:"shadow_intensity" gorm:"default:null"` // Intensidade de shadows (0.0-2.0, padrão 1.0)

	// Flag para indicar se está ativo ou usando padrão
	IsActive  bool      `json:"is_active" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ThemeCustomizationRequest é o DTO para requisições POST/PUT (sem ID, ProjectID, OrganizationID)
// Suporta customização separada para Light Mode e Dark Mode
type ThemeCustomizationRequest struct {
	// ==================== CORES PRINCIPAIS - LIGHT MODE (11 campos) ====================
	PrimaryColorLight        *string `json:"primary_color_light"`
	SecondaryColorLight      *string `json:"secondary_color_light"`
	BackgroundColorLight     *string `json:"background_color_light"`
	CardBackgroundColorLight *string `json:"card_background_color_light"`
	TextColorLight           *string `json:"text_color_light"`
	TextSecondaryColorLight  *string `json:"text_secondary_color_light"`
	AccentColorLight         *string `json:"accent_color_light"`

	// ==================== CORES PRINCIPAIS - DARK MODE (11 campos) ====================
	PrimaryColorDark        *string `json:"primary_color_dark"`
	SecondaryColorDark      *string `json:"secondary_color_dark"`
	BackgroundColorDark     *string `json:"background_color_dark"`
	CardBackgroundColorDark *string `json:"card_background_color_dark"`
	TextColorDark           *string `json:"text_color_dark"`
	TextSecondaryColorDark  *string `json:"text_secondary_color_dark"`
	AccentColorDark         *string `json:"accent_color_dark"`

	// ==================== CORES SEMÂNTICAS - LIGHT MODE (5 campos) ====================
	DestructiveColorLight *string `json:"destructive_color_light"`
	SuccessColorLight     *string `json:"success_color_light"`
	WarningColorLight     *string `json:"warning_color_light"`
	BorderColorLight      *string `json:"border_color_light"`
	PriceColorLight       *string `json:"price_color_light"`

	// ==================== CORES SEMÂNTICAS - DARK MODE (5 campos) ====================
	DestructiveColorDark *string `json:"destructive_color_dark"`
	SuccessColorDark     *string `json:"success_color_dark"`
	WarningColorDark     *string `json:"warning_color_dark"`
	BorderColorDark      *string `json:"border_color_dark"`
	PriceColorDark       *string `json:"price_color_dark"`

	// ==================== SISTEMA - LIGHT MODE (2 campos) ====================
	FocusRingColorLight      *string `json:"focus_ring_color_light"`
	InputBackgroundColorLight *string `json:"input_background_color_light"`

	// ==================== SISTEMA - DARK MODE (2 campos) ====================
	FocusRingColorDark       *string `json:"focus_ring_color_dark"`
	InputBackgroundColorDark *string `json:"input_background_color_dark"`

	// ==================== CONFIGURAÇÕES NUMÉRICAS (2 campos) ====================
	DisabledOpacity *float64 `json:"disabled_opacity"`
	ShadowIntensity *float64 `json:"shadow_intensity"`

	IsActive bool `json:"is_active"`
}

// TableName especifica o nome da tabela no banco de dados
func (ThemeCustomization) TableName() string {
	return "theme_customizations"
}
