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

	// Canais de notificação preferenciais
	DefaultNotificationChannel string `json:"default_notification_channel" gorm:"default:'sms'"` // sms, email, whatsapp
	EnableSms                  bool   `json:"enable_sms" gorm:"default:true"`
	EnableEmail                bool   `json:"enable_email" gorm:"default:false"`
	EnableWhatsapp             bool   `json:"enable_whatsapp" gorm:"default:false"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
