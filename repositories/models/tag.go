package models

import (
	"time"

	"github.com/google/uuid"
)

// --- Tag (etiquetas/categorias para entidades) ---
type Tag struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	Name           string     `json:"name" gorm:"not null"`
	Color          string     `json:"color,omitempty"`       // código hexadecimal ex: "#FF5733"
	Description    string     `json:"description,omitempty"`
	EntityType     string     `json:"entity_type,omitempty"` // "product", "customer", "table" - opcional
	Active         bool       `json:"active" gorm:"default:true"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// --- ProductTag (relacionamento produto-tag) ---
type ProductTag struct {
	Id        uuid.UUID `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductId uuid.UUID `json:"product_id" gorm:"not null"`
	TagId     uuid.UUID `json:"tag_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}
