package models

import (
	"time"

	"github.com/google/uuid"
)

// ThemeCustomization armazena as cores personalizadas do tema para cada projeto
type ThemeCustomization struct {
	ID             uuid.UUID `gorm:"primaryKey" json:"id"`
	ProjectID      uuid.UUID `json:"project_id" gorm:"not null;uniqueIndex"`
	OrganizationID uuid.UUID `json:"organization_id" gorm:"not null"`

	// Cores do tema - valores HEX (ex: #ffffff, #000000)
	PrimaryColor        string `json:"primary_color" gorm:"default:'#3b82f6'"`           // Azul principal
	SecondaryColor      string `json:"secondary_color" gorm:"default:'#8b5cf6'"`        // Roxo secundário
	BackgroundColor     string `json:"background_color" gorm:"default:'#09090b'"`       // Fundo escuro
	CardBackgroundColor string `json:"card_background_color" gorm:"default:'#18181b'"` // Fundo do card
	TextColor           string `json:"text_color" gorm:"default:'#fafafa'"`             // Texto principal
	TextSecondaryColor  string `json:"text_secondary_color" gorm:"default:'#a1a1aa'"` // Texto secundário
	AccentColor         string `json:"accent_color" gorm:"default:'#ec4899'"`           // Cor de destaque

	// Flag para indicar se está ativo ou usando padrão
	IsActive  bool      `json:"is_active" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ThemeCustomizationRequest é o DTO para requisições POST/PUT (sem ID, ProjectID, OrganizationID)
type ThemeCustomizationRequest struct {
	PrimaryColor        string `json:"primary_color"`
	SecondaryColor      string `json:"secondary_color"`
	BackgroundColor     string `json:"background_color"`
	CardBackgroundColor string `json:"card_background_color"`
	TextColor           string `json:"text_color"`
	TextSecondaryColor  string `json:"text_secondary_color"`
	AccentColor         string `json:"accent_color"`
	IsActive            bool   `json:"is_active"`
}

// TableName especifica o nome da tabela no banco de dados
func (ThemeCustomization) TableName() string {
	return "theme_customizations"
}
