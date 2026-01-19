package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// --- Notification Models (SPRINT 2) ---

// NotificationConfig - Configurações de notificação por evento
type NotificationConfig struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	EventType      string     `json:"event_type"` // "reservation_create", "reservation_update", "table_available", etc.
	Enabled        bool       `json:"enabled" gorm:"default:true"`
	Channels       pq.StringArray `json:"channels" gorm:"type:text[]"` // ["sms", "email", "whatsapp"]
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
	Variables      pq.StringArray `json:"variables" gorm:"type:text[]"` // Lista de variáveis disponíveis
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

	// Campos para rastreamento de reserva e cliente
	ReservationId *uuid.UUID `json:"reservation_id,omitempty"` // Link para reserva identificada
	CustomerId    *uuid.UUID `json:"customer_id,omitempty"`    // Link para cliente identificado

	// Campos para classificação de resposta
	ResponseType     string  `json:"response_type,omitempty"`     // "confirmed", "cancelled", "unknown"
	ConfidenceScore  float64 `json:"confidence_score,omitempty"`  // Score de confiança (0.0 a 1.0)
	ProcessingMethod string  `json:"processing_method,omitempty"` // "pattern_match", "ai_classification"
	ActionTaken      string  `json:"action_taken,omitempty"`      // "reservation_confirmed", "reservation_cancelled", "queued_for_review", "no_action"
}

// NotificationSchedule - Agendamento de notificações futuras
type NotificationSchedule struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	EventType      string     `json:"event_type"`  // "confirmation_request", "reminder", "auto_cancel_warning"
	EntityType     string     `json:"entity_type"` // "reservation"
	EntityId       uuid.UUID  `json:"entity_id"`   // ID da reserva
	ScheduledFor   time.Time  `json:"scheduled_for"`
	Status         string     `json:"status" gorm:"default:'pending'"` // "pending", "sent", "cancelled", "skipped", "failed"
	ProcessedAt    *time.Time `json:"processed_at,omitempty"`
	Metadata       string     `json:"metadata" gorm:"type:json"` // JSON com dados do cliente/reserva
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ResponseReviewQueue - Fila de revisão de respostas de clientes
type ResponseReviewQueue struct {
	Id              uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId  uuid.UUID  `json:"organization_id"`
	ProjectId       uuid.UUID  `json:"project_id"`
	InboundId       uuid.UUID  `json:"inbound_id"`     // Link para NotificationInbound
	ReservationId   uuid.UUID  `json:"reservation_id"` // Link para Reservation
	CustomerId      uuid.UUID  `json:"customer_id"`    // Link para Customer
	CustomerName    string     `json:"customer_name"`  // Para exibição rápida
	CustomerPhone   string     `json:"customer_phone"`
	MessageBody     string     `json:"message_body"`                          // Mensagem original do cliente
	SuggestedAction string     `json:"suggested_action,omitempty"`            // "confirm", "cancel", "none"
	ConfidenceScore float64    `json:"confidence_score,omitempty"`            // Score de confiança da sugestão
	Status          string     `json:"status" gorm:"default:'pending_review'"` // "pending_review", "approved", "rejected", "expired"
	ReviewedBy      *uuid.UUID `json:"reviewed_by,omitempty"`                 // ID do usuário que revisou
	ReviewedAt      *time.Time `json:"reviewed_at,omitempty"`
	ActionTaken     string     `json:"action_taken,omitempty"` // Ação final tomada
	Notes           string     `json:"notes,omitempty"`        // Notas do revisor
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
