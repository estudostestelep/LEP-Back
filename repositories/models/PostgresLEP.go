package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// --- Organization (organização mãe) ---
type Organization struct {
	Id          uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string     `json:"name" gorm:"not null"`
	Email       string     `gorm:"unique" json:"email"`
	Phone       string     `json:"phone,omitempty"`
	Address     string     `json:"address,omitempty"`
	Website     string     `json:"website,omitempty"`
	Description string     `json:"description,omitempty"`
	Active      bool       `gorm:"default:true" json:"active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type BannedLists struct {
	BannedListId uuid.UUID `gorm:"primaryKey;autoIncrement" json:"banned_list_id"`
	Token        string    `gorm:"type:varchar(300);unique" json:"token"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type LoggedLists struct {
	LoggedListId uuid.UUID `gorm:"primaryKey;autoIncrement" json:"logged_list_id"`
	Token        string    `gorm:"type:varchar(300);unique" json:"token"`
	UserEmail    string    `gorm:"type:varchar(300)" json:"user_email"`
	UserId       uuid.UUID `json:"user_id"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// --- User (funcionário/admin) ---
type User struct {
	Id             uuid.UUID      `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID      `json:"organization_id"`
	ProjectId      uuid.UUID      `json:"project_id"`
	Name           string         `json:"name"`
	Email          string         `gorm:"unique" json:"email"`
	Password       string         `json:"password"` // armazenar hash!
	Role           string         `json:"role"`     // ex: "waiter", "admin"
	Permissions    pq.StringArray `gorm:"type:text[]" json:"permissions"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      *time.Time     `json:"deleted_at,omitempty"`
}

// --- Customer (cliente do restaurante) ---""
type Customer struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	Phone          string     `json:"phone"`
	BirthDate      string     `json:"birth_date,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// --- Table (mesa) ---
type Table struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	EnvironmentId  *uuid.UUID `json:"environment_id,omitempty"` // vinculação com ambiente
	Number         int        `json:"number"`
	Capacity       int        `json:"capacity"`
	Location       string     `json:"location,omitempty"`            // descrição adicional da localização
	Status         string     `json:"status" gorm:"default:'livre'"` // "livre", "ocupada", "reservada"
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// --- Product (item do cardápio) ---
type Product struct {
	Id              uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId  uuid.UUID  `json:"organization_id"`
	ProjectId       uuid.UUID  `json:"project_id"`
	Name            string     `json:"name"`
	Category        string     `json:"category"`
	Description     string     `json:"description"`
	Price           float64    `json:"price"`
	Available       bool       `json:"available"`
	Stock           *int       `json:"stock,omitempty"`   // opcional
	PrepTimeMinutes int        `json:"prep_time_minutes"` // tempo de preparo em minutos
	ImageUrl        *string    `gorm:"column:image_url" json:"image_url,omitempty"` // URL da imagem do produto
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

// --- Reservation (reserva de mesa) ---
type Reservation struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	CustomerId     uuid.UUID  `json:"customer_id"`
	TableId        uuid.UUID  `json:"table_id"`
	Datetime       string     `json:"datetime"`
	PartySize      int        `json:"party_size"`
	Note           string     `json:"note,omitempty"`
	Status         string     `json:"status"` // ex: "confirmed", "cancelled"
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// --- Waitlist (fila de espera) ---
type Waitlist struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	CustomerId     uuid.UUID  `json:"customer_id"`
	People         int        `json:"people"`
	Status         string     `json:"status"` // ex: "waiting", "notified", "seated"
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// --- OrderItem (item do pedido) ---
type OrderItem struct {
	ProductId uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"` // valor unitário
	Notes     string    `json:"notes,omitempty"` // observações do item
}

// OrderItems é um tipo customizado para array de OrderItem que funciona com JSONB
type OrderItems []OrderItem

// Value implementa driver.Valuer para serializar para o banco
func (oi OrderItems) Value() (driver.Value, error) {
	if len(oi) == 0 {
		return "[]", nil
	}
	return json.Marshal(oi)
}

// Scan implementa sql.Scanner para deserializar do banco
func (oi *OrderItems) Scan(value interface{}) error {
	if value == nil {
		*oi = OrderItems{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("cannot scan value into OrderItems: unsupported type")
	}

	// Se for uma string vazia ou array vazio, inicializar como array vazio
	if len(bytes) == 0 || string(bytes) == "" || string(bytes) == "null" {
		*oi = OrderItems{}
		return nil
	}

	// Tentar deserializar como array primeiro
	var items []OrderItem
	if err := json.Unmarshal(bytes, &items); err == nil {
		*oi = OrderItems(items)
		return nil
	}

	// Se falhou como array, pode ser um objeto único (dados antigos)
	// Tentar deserializar como objeto único e envolver em array
	var singleItem OrderItem
	if err := json.Unmarshal(bytes, &singleItem); err == nil {
		*oi = OrderItems{singleItem}
		return nil
	}

	// Se ambos falharam, retornar erro
	return errors.New("cannot unmarshal value into OrderItems: invalid JSON format")
}

// --- Order (pedido) ---
type Order struct {
	Id                    uuid.UUID   `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId        uuid.UUID   `json:"organization_id"`
	ProjectId             uuid.UUID   `json:"project_id"`
	TableId               *uuid.UUID  `json:"table_id,omitempty"`
	TableNumber           *int        `json:"table_number,omitempty"` // Para pedidos públicos
	CustomerId            *uuid.UUID  `json:"customer_id,omitempty"`
	Items                 OrderItems  `gorm:"type:jsonb" json:"items"`
	TotalAmount           float64     `json:"total_amount"`
	Note                  string      `json:"note,omitempty"`
	Source                string      `json:"source"`                            // "internal" ou "public"
	Status                string      `json:"status"`                            // "pending", "preparing", "ready", "delivered", "cancelled"
	EstimatedPrepTime     int         `json:"estimated_prep_time_minutes"`       // tempo estimado total em minutos
	EstimatedDeliveryTime *time.Time  `json:"estimated_delivery_time,omitempty"` // hora estimada de entrega
	StartedAt             *time.Time  `json:"started_at,omitempty"`              // quando começou a preparar
	ReadyAt               *time.Time  `json:"ready_at,omitempty"`                // quando ficou pronto
	DeliveredAt           *time.Time  `json:"delivered_at,omitempty"`            // quando foi entregue
	CreatedAt             time.Time   `json:"created_at"`
	UpdatedAt             time.Time   `json:"updated_at"`
	DeletedAt             *time.Time  `json:"deleted_at,omitempty"`
}

// --- Project (configurações centralizadas) ---
type Project struct {
	Id             uuid.UUID `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description,omitempty"`

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

	// Configurações gerais
	TimeZone  string     `json:"timezone" gorm:"default:'America/Sao_Paulo'"`
	Active    bool       `json:"active" gorm:"default:true"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

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

// --- Environment (ambientes do restaurante) ---
type Environment struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	Name           string     `json:"name"` // ex: "Salão Principal", "Varanda"
	Description    string     `json:"description,omitempty"`
	Capacity       int        `json:"capacity"` // capacidade máxima do ambiente
	Active         bool       `json:"active" gorm:"default:true"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// --- Log/Audit (auditoria de ações) ---
type AuditLog struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	UserId         *uuid.UUID `json:"user_id,omitempty"`
	Action         string     `json:"action"` // ex: "create_reservation", "cancel_order"
	Entity         string     `json:"entity"` // ex: "reservation", "order"
	EntityId       uuid.UUID  `json:"entity_id"`
	Description    string     `json:"description,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

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

// --- SPRINT 4: Validações Avançadas ---

// BlockedPeriod - Períodos bloqueados para reservas
type BlockedPeriod struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	Name           string     `json:"name"` // ex: "Manutenção", "Evento Privado"
	Description    string     `json:"description,omitempty"`
	StartDateTime  time.Time  `json:"start_datetime"`
	EndDateTime    time.Time  `json:"end_datetime"`
	RecurringType  string     `json:"recurring_type,omitempty"` // "none", "weekly", "monthly"
	Active         bool       `json:"active" gorm:"default:true"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// --- SPRINT 5: Features Avançadas ---

// Lead - Sistema básico de CRM
type Lead struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	Name           string     `json:"name"`
	Email          string     `json:"email,omitempty"`
	Phone          string     `json:"phone,omitempty"`
	Source         string     `json:"source"` // "waitlist", "reservation", "walk_in"
	Status         string     `json:"status"` // "new", "contacted", "converted", "lost"
	Notes          string     `json:"notes,omitempty"`
	LastContact    *time.Time `json:"last_contact,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// ReportMetric - Métricas básicas para relatórios
type ReportMetric struct {
	Id             uuid.UUID `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID `json:"organization_id"`
	ProjectId      uuid.UUID `json:"project_id"`
	MetricType     string    `json:"metric_type"`                         // "reservation_count", "occupancy_rate", "revenue"
	MetricDate     time.Time `json:"metric_date"`                         // data da métrica (dia)
	Value          float64   `json:"value"`                               // valor da métrica
	Metadata       string    `json:"metadata,omitempty" gorm:"type:json"` // dados extras em JSON
	CreatedAt      time.Time `json:"created_at"`
}
