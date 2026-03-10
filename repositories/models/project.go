package models

import (
	"time"

	"github.com/google/uuid"
)

// --- Project (configurações centralizadas) ---
type Project struct {
	Id             uuid.UUID `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description,omitempty"`
	Slug           string    `json:"slug" gorm:"size:100;index"`
	IsDefault      bool      `json:"is_default" gorm:"default:false"`

	// Configurações Twilio (SMS + WhatsApp)
	TwilioAccountSid       *string `json:"twilio_account_sid,omitempty"`
	TwilioAuthToken        *string `json:"twilio_auth_token,omitempty"`
	TwilioPhone            *string `json:"twilio_phone,omitempty"`
	WhatsappBusinessNumber *string `json:"whatsapp_business_number,omitempty"`

	// Configurações SMTP (Email)
	SmtpHost     *string `json:"smtp_host,omitempty"`
	SmtpPort     *int    `json:"smtp_port,omitempty"`
	SmtpUsername *string `json:"smtp_username,omitempty"`
	SmtpPassword *string `json:"smtp_password,omitempty"`
	SmtpFrom     *string `json:"smtp_from,omitempty"`

	// Configurações de Notificação
	NotificationResponsiblePhone *string `json:"notification_responsible_phone,omitempty"` // Número que recebe cópia das notificações de reserva

	// Configurações gerais
	TimeZone  string     `json:"timezone" gorm:"default:'America/Sao_Paulo'"`
	Active    bool       `json:"active" gorm:"default:true"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
