package models

import (
	"time"

	"github.com/google/uuid"
)

type BannedLists struct {
	BannedListId uuid.UUID `gorm:"primaryKey;autoIncrement" json:"banned_list_id"`
	Token        string    `gorm:"type:varchar(300)" json:"token"`
	Date         string    `gorm:"type:varchar(300)" json:"date"`
}

type LoggedLists struct {
	LoggedListId uuid.UUID `gorm:"primaryKey;autoIncrement" json:"logged_list_id"`
	Token        string    `gorm:"type:varchar(300)" json:"token"`
	UserEmail    string    `gorm:"type:varchar(300)" json:"user_email"`
	UserId       uuid.UUID `json:"user_id"`
}

// --- User (funcionário/admin) ---
type User struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	Name           string     `json:"name"`
	Email          string     `gorm:"unique" json:"email"`
	Password       string     `json:"-"`           // armazenar hash!
	Role           string     `json:"role"`        // ex: "waiter", "admin"
	Permissions    []string   `json:"permissions"` // ex: ["view_orders", "create_reservation"]
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// --- Customer (cliente do restaurante) ---""
type Customer struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	Phone          string     `json:"phone"`
	BirthDate      *time.Time `json:"birth_date,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// --- Table (mesa) ---
type Table struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	Number         int        `json:"number"`
	Capacity       int        `json:"capacity"`
	Location       string     `json:"location"`
	IsAvailable    bool       `json:"is_available"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// --- Product (item do cardápio) ---
type Product struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Price          float64    `json:"price"`
	Available      bool       `json:"available"`
	Stock          *int       `json:"stock,omitempty"` // opcional
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// --- Reservation (reserva de mesa) ---
type Reservation struct {
	Id             uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	CustomerId     uuid.UUID  `json:"customer_id"`
	TableId        uuid.UUID  `json:"table_id"`
	Datetime       time.Time  `json:"datetime"`
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
}

// --- Order (pedido) ---
type Order struct {
	Id             uuid.UUID   `json:"id"`
	OrganizationId uuid.UUID   `json:"organization_id"`
	ProjectId      uuid.UUID   `json:"project_id"`
	TableId        *uuid.UUID  `json:"table_id,omitempty"`
	TableNumber    *int        `json:"table_number,omitempty"` // Para pedidos públicos
	CustomerId     *uuid.UUID  `json:"customer_id,omitempty"`
	Items          []OrderItem `json:"items"`
	Total          float64     `json:"total"`
	Note           string      `json:"note,omitempty"`
	Source         string      `json:"source"` // "internal" ou "public"
	Status         string      `json:"status"` // ex: "pending", "completed", "cancelled"
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	DeletedAt      *time.Time  `json:"deleted_at,omitempty"`
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
