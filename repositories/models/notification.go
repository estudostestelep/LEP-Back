package models

import (
	"time"

	"github.com/google/uuid"
)

// --- Notification Models (SPRINT 2) ---

// NotificationConfig - Configurações de notificação por evento
type NotificationConfig struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	EventType      string     `json:"event_type"` // "reservation_create", "reservation_update", "table_available", etc.
	Enabled        bool       `json:"enabled" gorm:"default:true"`
	Channels       []string   `json:"channels" gorm:"type:text[]"` // ["sms", "email", "whatsapp"]
	TemplateId     *uuid.UUID `json:"template_id,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// NotificationTemplate - Templates de mensagens
type NotificationTemplate struct {
	Id             uuid.UUID `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID `json:"organization_id"`
	ProjectId      uuid.UUID `json:"project_id"`
	Name           string    `json:"name"`
	Channel        string    `json:"channel"`                      // "sms", "email", "whatsapp"
	Subject        string    `json:"subject,omitempty"`            // Para email
	Body           string    `json:"body"`                         // Conteúdo com variáveis {{nome}}, {{data}}, etc.
	Variables      []string  `json:"variables" gorm:"type:text[]"` // Lista de variáveis disponíveis
	Active         bool      `json:"active" gorm:"default:true"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// NotificationLog - Log de notificações enviadas
type NotificationLog struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	EventType      string     `json:"event_type"`
	Channel        string     `json:"channel"`
	Recipient      string     `json:"recipient"` // telefone, email
	Subject        string     `json:"subject,omitempty"`
	Message        string     `json:"message"`
	Status         string     `json:"status"`                // "sent", "delivered", "failed", "pending"
	ExternalId     string     `json:"external_id,omitempty"` // MessageSid do Twilio, etc.
	ErrorMessage   string     `json:"error_message,omitempty"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// NotificationEvent - Eventos que disparam notificações
type NotificationEvent struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	EventType      string     `json:"event_type"`
	EntityType     string     `json:"entity_type"` // "reservation", "order", "table"
	EntityId       uuid.UUID  `json:"entity_id"`
	Data           string     `json:"data" gorm:"type:json"` // Dados do evento em JSON
	Processed      bool       `json:"processed" gorm:"default:false"`
	ProcessedAt    *time.Time `json:"processed_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

// NotificationInbound - Mensagens recebidas (2-way)
type NotificationInbound struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	Channel        string     `json:"channel"`     // "sms", "whatsapp"
	From           string     `json:"from"`        // número do remetente
	To             string     `json:"to"`          // número do destino (nosso)
	Body           string     `json:"body"`        // conteúdo da mensagem
	ExternalId     string     `json:"external_id"` // MessageSid do Twilio
	Processed      bool       `json:"processed" gorm:"default:false"`
	ProcessedAt    *time.Time `json:"processed_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}
