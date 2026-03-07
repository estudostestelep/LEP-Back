package models

import (
	"time"

	"github.com/google/uuid"
)

// StaffSchedule representa a escalação de um funcionário para um dia/turno específico
type StaffSchedule struct {
	Id             uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id" gorm:"type:uuid;not null;index"`
	ProjectId      uuid.UUID  `json:"project_id" gorm:"type:uuid;not null;index"`
	ClientId       uuid.UUID  `json:"client_id" gorm:"type:uuid;not null;index"`   // FK para clients
	ScheduleDate   time.Time  `json:"schedule_date" gorm:"not null;index"`         // data específica
	Shift          string     `json:"shift" gorm:"not null"`                       // "almoco" | "noite"
	Status         string     `json:"status" gorm:"default:'scheduled'"`           // scheduled, confirmed, cancelled
	SlotNumber     int        `json:"slot_number" gorm:"default:0"`                // posição na escala (1-15)
	EmailSentAt    *time.Time `json:"email_sent_at,omitempty"`                     // quando o email foi enviado
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relacionamento para eager loading
	Client *Client `json:"client,omitempty" gorm:"foreignKey:ClientId;references:Id"`
}

// TableName define o nome da tabela no banco de dados
func (StaffSchedule) TableName() string {
	return "staff_schedules"
}

// Constantes para status da escala
const (
	ScheduleStatusScheduled = "scheduled"
	ScheduleStatusConfirmed = "confirmed"
	ScheduleStatusCancelled = "cancelled"
)

// Constantes para turnos
const (
	ShiftLunch = "almoco"
	ShiftNight = "noite"
)

// StaffScheduleWithClient é um DTO para retornar escala com nome do funcionário
type StaffScheduleWithClient struct {
	StaffSchedule
	ClientName  string `json:"client_name"`
	ClientEmail string `json:"client_email"`
}

// WeekScheduleSummary resume as escalas de uma semana por dia
type WeekScheduleSummary struct {
	WeekStart string                          `json:"week_start"`
	ByDay     map[string][]ScheduledClient    `json:"by_day"` // { "2026-01-15": [{...}], ... }
}

// ScheduledClient representa um funcionário escalado com detalhes
type ScheduledClient struct {
	ScheduleId uuid.UUID `json:"schedule_id"`
	ClientId   uuid.UUID `json:"client_id"`
	ClientName string    `json:"client_name"`
	Shift      string    `json:"shift"`
	Status     string    `json:"status"`
	SlotNumber int       `json:"slot_number"`
}

// SendScheduleEmailsRequest representa a requisição para enviar emails da escala
type SendScheduleEmailsRequest struct {
	WeekStart string `json:"week_start" binding:"required"` // formato: "2026-01-13"
}

// ScheduleUpdateDiff representa as diferenças em uma atualização de escala
type ScheduleUpdateDiff struct {
	Added   []ScheduledClient `json:"added"`
	Removed []ScheduledClient `json:"removed"`
}
