package models

import (
	"time"

	"github.com/google/uuid"
)

// --- SISTEMA DE CARDÁPIO ---

// Menu (Cardápio principal)
type Menu struct {
	Id                uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId    uuid.UUID  `json:"organization_id" gorm:"not null"`
	ProjectId         uuid.UUID  `json:"project_id" gorm:"not null"`
	Name              string     `json:"name" gorm:"not null"`
	Styling           *string    `json:"styling,omitempty" gorm:"type:json"` // JSON com configs visuais
	Order             int        `json:"order" gorm:"default:0"`
	Active            bool       `json:"active" gorm:"default:true"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`

	// ✨ Novos campos para seleção automática de cardápio
	TimeRangeStart    *time.Time `json:"time_range_start,omitempty" gorm:"column:time_range_start"`           // Horário de início (ex: 11:00)
	TimeRangeEnd      *time.Time `json:"time_range_end,omitempty" gorm:"column:time_range_end"`             // Horário de fim (ex: 15:00)
	Priority          int        `json:"priority" gorm:"default:999"`                                         // Prioridade (0 = maior)
	IsManualOverride  bool       `json:"is_manual_override" gorm:"default:false"`                             // Se está como override manual
	ApplicableDays    *string    `json:"applicable_days,omitempty" gorm:"type:json"`                          // Array JSON: [0,1,2,3,4] = dias da semana
	ApplicableDates   *string    `json:"applicable_dates,omitempty" gorm:"type:json"`                         // Array JSON: ["2025-12-25", "2025-01-01"]
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
