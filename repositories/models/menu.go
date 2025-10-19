package models

import (
	"time"

	"github.com/google/uuid"
)

// --- SISTEMA DE CARDÁPIO ---

// Menu (Cardápio principal)
type Menu struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id" gorm:"not null"`
	ProjectId      uuid.UUID  `json:"project_id" gorm:"not null"`
	Name           string     `json:"name" gorm:"not null"`
	Styling        *string    `json:"styling,omitempty" gorm:"type:json"` // JSON com configs visuais
	Order          int        `json:"order" gorm:"default:0"`
	Active         bool       `json:"active" gorm:"default:true"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// Category (Categoria do cardápio)
type Category struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id" gorm:"not null"`
	ProjectId      uuid.UUID  `json:"project_id" gorm:"not null"`
	MenuId         uuid.UUID  `json:"menu_id" gorm:"not null"`
	Name           string     `json:"name" gorm:"not null"`
	ImageUrl       *string    `json:"image_url" gorm:"column:photo"`
	Notes          *string    `json:"notes,omitempty"`
	Order          int        `json:"order" gorm:"default:0"`
	Active         bool       `json:"active" gorm:"default:true"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// Subcategory (Subcategoria)
type Subcategory struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id" gorm:"not null"`
	ProjectId      uuid.UUID  `json:"project_id" gorm:"not null"`
	Name           string     `json:"name" gorm:"not null"`
	Photo          *string    `json:"photo,omitempty"`
	Notes          *string    `json:"notes,omitempty"`
	Order          int        `json:"order" gorm:"default:0"`
	Active         bool       `json:"active" gorm:"default:true"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// SubcategoryCategory (Relacionamento N:N entre Subcategoria e Categoria)
type SubcategoryCategory struct {
	Id            uuid.UUID `gorm:"primaryKey;autoIncrement" json:"id"`
	SubcategoryId uuid.UUID `json:"subcategory_id" gorm:"not null"`
	CategoryId    uuid.UUID `json:"category_id" gorm:"not null"`
	CreatedAt     time.Time `json:"created_at"`
}
