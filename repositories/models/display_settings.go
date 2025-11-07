package models

import (
	"time"

	"github.com/google/uuid"
)

// ProjectDisplaySettings controla quais campos dos produtos são exibidos
// no sistema (menu público, admin-menu, APIs)
type ProjectDisplaySettings struct {
	ID              uuid.UUID `gorm:"primaryKey" json:"id"`
	ProjectID       uuid.UUID `json:"project_id" gorm:"not null;uniqueIndex"`
	OrganizationID  uuid.UUID `json:"organization_id" gorm:"not null"`

	// Campos de exibição de produtos
	ShowPrepTime    bool `json:"show_prep_time" gorm:"default:true"`      // ⏱️ Tempo de preparo
	ShowRating      bool `json:"show_rating" gorm:"default:true"`         // ⭐ Avaliações
	ShowDescription bool `json:"show_description" gorm:"default:true"`    // 📝 Descrição

	// Campos opcionais para futuras expansões
	ShowPrice        *bool `json:"show_price" gorm:"default:true"`        // 💰 Preço
	ShowAvailability *bool `json:"show_availability" gorm:"default:true"` // ✓ Disponibilidade
	ShowAllergens    *bool `json:"show_allergens" gorm:"default:true"`    // ⚠️ Alergênios

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName especifica o nome da tabela no banco de dados
func (ProjectDisplaySettings) TableName() string {
	return "project_display_settings"
}
