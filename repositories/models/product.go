package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// --- Product (item do cardápio - REFATORADO) ---
type Product struct {
	// Campos base
	Id              uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId  uuid.UUID  `json:"organization_id" gorm:"not null"`
	ProjectId       uuid.UUID  `json:"project_id" gorm:"not null"`
	Name            string     `json:"name" gorm:"not null"`
	Description     string     `json:"description,omitempty"`
	ImageUrl        *string    `gorm:"column:image_url" json:"image_url,omitempty"`

	// NOVOS - Tipo e organização
	Type            string     `json:"type" gorm:"not null"` // "prato" | "bebida" | "vinho"
	Order           int        `json:"order" gorm:"default:0"`
	Active          bool       `json:"active" gorm:"default:true"`
	PDVCode         *string    `json:"pdv_code,omitempty"` // código para integração PDV

	// NOVOS - Relacionamentos (substituem category string)
	CategoryId      *uuid.UUID `json:"category_id,omitempty"`
	SubcategoryId   *uuid.UUID `json:"subcategory_id,omitempty"`

	// Campos de preço (refatorados)
	PriceNormal     float64    `json:"price_normal" gorm:"not null"`
	PricePromo      *float64   `json:"price_promo,omitempty"`

	// NOVOS - Campos para Bebida/Vinho
	Volume          *int       `json:"volume,omitempty"`           // ml
	AlcoholContent  *float64   `json:"alcohol_content,omitempty"`  // % teor alcoólico

	// NOVOS - Campos específicos de Vinho
	Vintage         *string        `json:"vintage,omitempty"`          // safra
	Country         *string        `json:"country,omitempty"`          // país de origem
	Region          *string        `json:"region,omitempty"`           // região
	Winery          *string        `json:"winery,omitempty"`           // vinícola
	WineType        *string        `json:"wine_type,omitempty"`        // tinto, branco, rosé, etc
	Grapes          pq.StringArray `json:"grapes,omitempty" gorm:"type:text[]"` // uvas (multi-select)
	PriceBottle     *float64       `json:"price_bottle,omitempty"`     // preço garrafa
	PriceHalfBottle *float64       `json:"price_half_bottle,omitempty"` // preço meia garrafa
	PriceGlass      *float64       `json:"price_glass,omitempty"`      // preço taça

	// Campos existentes
	Stock           *int      `json:"stock,omitempty"`
	PrepTimeMinutes int       `json:"prep_time_minutes,omitempty"`

	// Timestamps
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	// Relacionamento muitos-para-muitos com Tags (para eager loading)
	Tags []Tag `gorm:"many2many:product_tags;" json:"tags,omitempty"`
}
