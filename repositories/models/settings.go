package models

import (
	"time"

	"github.com/google/uuid"
)

// --- Settings (configurações parametrizáveis) ---
type Settings struct {
	Id             uuid.UUID `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID `json:"organization_id"`
	ProjectId      uuid.UUID `json:"project_id"`

	// Configurações de antecedência para reservas
	MinAdvanceHours int `json:"min_advance_hours" gorm:"default:2"` // mínimo 2h antes
	MaxAdvanceDays  int `json:"max_advance_days" gorm:"default:30"` // máximo 30 dias antes

	// Configurações de notificação
	NotifyReservationCreate bool `json:"notify_reservation_create" gorm:"default:true"`
	NotifyReservationUpdate bool `json:"notify_reservation_update" gorm:"default:true"`
	NotifyReservationCancel bool `json:"notify_reservation_cancel" gorm:"default:true"`
	NotifyTableAvailable    bool `json:"notify_table_available" gorm:"default:true"`
	NotifyConfirmation24h   bool `json:"notify_confirmation_24h" gorm:"default:true"`

	// Configurações de agendamento flexível de notificações
	ConfirmationHoursBefore   int `json:"confirmation_hours_before" gorm:"default:24"`   // Horas antes da reserva para enviar confirmação
	ReminderHoursBefore       int `json:"reminder_hours_before" gorm:"default:0"`        // Horas antes para lembrete (0 = desabilitado)
	AutoCancelNoResponseHours int `json:"auto_cancel_no_response_hours" gorm:"default:0"` // Cancelar automaticamente se sem resposta (0 = desabilitado)

	// Limite de pessoas para confirmação automática (0 = todos confirmados diretamente)
	AutoConfirmMaxPartySize int `json:"auto_confirm_max_party_size" gorm:"default:0"`

	// Horários de funcionamento
	LunchStart            string `json:"lunch_start" gorm:"default:'12:00'"`
	LunchEnd              string `json:"lunch_end" gorm:"default:'14:30'"`
	DinnerStart           string `json:"dinner_start" gorm:"default:'19:00'"`
	DinnerEnd             string `json:"dinner_end" gorm:"default:'22:00'"`
	SlotIntervalMinutes   int    `json:"slot_interval_minutes" gorm:"default:30"`
	EnableLunch           bool   `json:"enable_lunch" gorm:"default:true"`
	EnableDinner          bool   `json:"enable_dinner" gorm:"default:true"`
	DiningDurationMinutes int    `json:"dining_duration_minutes" gorm:"default:120"` // tempo médio de permanência na mesa

	// Agenda semanal de funcionamento (JSON)
	// Formato: {"0":{"enabled":false,"enable_lunch":false,"enable_dinner":false},...}
	// Chaves: 0=Domingo, 1=Segunda, ..., 6=Sábado
	// Quando preenchido, sobrescreve EnableLunch/EnableDinner por dia
	OperatingScheduleJson string `json:"operating_schedule_json" gorm:"type:text;default:''"`

	// Modo de processamento de respostas do cliente
	// "automatic" - Sistema processa e atualiza reserva automaticamente
	// "ai_assisted" - IA analisa e sugere ação, aguarda aprovação humana
	// "manual" - Todas as respostas vão para fila de revisão manual
	ResponseProcessingMode string `json:"response_processing_mode" gorm:"default:'automatic'"`

	// Canais de notificação preferenciais
	DefaultNotificationChannel string `json:"default_notification_channel" gorm:"default:'sms'"` // sms, email, whatsapp
	EnableSms                  bool   `json:"enable_sms" gorm:"default:true"`
	EnableEmail                bool   `json:"enable_email" gorm:"default:false"`
	EnableWhatsapp             bool   `json:"enable_whatsapp" gorm:"default:false"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
