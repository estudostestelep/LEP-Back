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
