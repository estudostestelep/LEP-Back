package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// StaffAvailability representa a disponibilidade semanal de um funcionário
// Um funcionário marca os dias/turnos que pode trabalhar em uma semana específica
type StaffAvailability struct {
	Id             uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrganizationId uuid.UUID      `json:"organization_id" gorm:"type:uuid;not null;index"`
	ProjectId      uuid.UUID      `json:"project_id" gorm:"type:uuid;not null;index"`
	ClientId       uuid.UUID      `json:"client_id" gorm:"type:uuid;not null;index"` // FK para clients
	WeekStart      time.Time      `json:"week_start" gorm:"not null;index"`          // segunda-feira da semana
	AvailableDays  pq.StringArray `json:"available_days" gorm:"type:text[]"`         // ["qua", "qui", "sex", "sab", "dom_almoco", "dom_noite"]
	NoAvailability bool           `json:"no_availability" gorm:"default:false"`      // sem disponibilidade na semana
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`

	// Relacionamento para eager loading
	Client *Client `json:"client,omitempty" gorm:"foreignKey:ClientId;references:Id"`
}

// TableName define o nome da tabela no banco de dados
func (StaffAvailability) TableName() string {
	return "staff_availabilities"
}

// StaffAvailabilityWithClient é um DTO para retornar disponibilidade com nome do funcionário
type StaffAvailabilityWithClient struct {
	StaffAvailability
	ClientName string `json:"client_name"`
}

// WeekAvailabilitySummary resume as disponibilidades de uma semana
type WeekAvailabilitySummary struct {
	WeekStart      string                       `json:"week_start"`
	TotalResponses int                          `json:"total_responses"`
	ByDay          map[string][]AvailableClient `json:"by_day"` // { "qua": [{id, name}, ...], "qui": [...] }
}

// AvailableClient representa um funcionário disponível para um dia
type AvailableClient struct {
	ClientId   uuid.UUID `json:"client_id"`
	ClientName string    `json:"client_name"`
}
