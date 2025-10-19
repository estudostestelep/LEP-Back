package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

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
