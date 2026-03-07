package models

import (
	"time"

	"github.com/google/uuid"
)

// StaffAttendance representa o registro de presença de um funcionário
type StaffAttendance struct {
	Id             uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id" gorm:"type:uuid;not null;index"`
	ProjectId      uuid.UUID  `json:"project_id" gorm:"type:uuid;not null;index"`
	ClientId       uuid.UUID  `json:"client_id" gorm:"type:uuid;not null;index"`
	AttendanceDate time.Time  `json:"attendance_date" gorm:"not null;index"`
	Shift          string     `json:"shift" gorm:"not null"` // "almoco" | "noite" | "evento"
	CheckedIn      bool       `json:"checked_in" gorm:"default:true"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relacionamentos
	Client      *Client                       `json:"client,omitempty" gorm:"foreignKey:ClientId;references:Id"`
	Transports  []StaffAttendanceTransport    `json:"transports,omitempty" gorm:"foreignKey:AttendanceId"`
	Consumptions []StaffAttendanceConsumption `json:"consumptions,omitempty" gorm:"foreignKey:AttendanceId"`
}

// TableName define o nome da tabela no banco de dados
func (StaffAttendance) TableName() string {
	return "staff_attendances"
}

// Constantes para turno de presença
const (
	AttendanceShiftLunch = "almoco"
	AttendanceShiftNight = "noite"
	AttendanceShiftEvent = "evento"
)

// StaffAttendanceTransport representa o transporte de ida ou volta
type StaffAttendanceTransport struct {
	Id             uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	AttendanceId   uuid.UUID  `json:"attendance_id" gorm:"type:uuid;not null;index"`
	Direction      string     `json:"direction" gorm:"not null"` // "ida" | "volta"
	Mode           string     `json:"mode,omitempty"`            // "carona" | "onibus" | "uber" | ""
	RideWithClient *uuid.UUID `json:"ride_with_client,omitempty" gorm:"type:uuid"` // se carona, com quem
	BusTickets     *int       `json:"bus_tickets,omitempty"`     // se ônibus, quantas passagens (1-3)
	UberCost       *float64   `json:"uber_cost,omitempty"`       // se uber, valor em R$
	CreatedAt      time.Time  `json:"created_at"`

	// Para eager loading do nome de quem deu carona
	RideWith *Client `json:"ride_with,omitempty" gorm:"foreignKey:RideWithClient;references:Id"`
}

// TableName define o nome da tabela no banco de dados
func (StaffAttendanceTransport) TableName() string {
	return "staff_attendance_transports"
}

// Constantes para direção do transporte
const (
	TransportDirectionIn  = "ida"
	TransportDirectionOut = "volta"
)

// Constantes para modo de transporte
const (
	TransportModeRide = "carona"
	TransportModeBus  = "onibus"
	TransportModeUber = "uber"
)

// StaffAttendanceConsumption representa o consumo de produtos pelo funcionário
type StaffAttendanceConsumption struct {
	Id           uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	AttendanceId uuid.UUID `json:"attendance_id" gorm:"type:uuid;not null;index"`
	ProductId    uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"` // FK para staff_consumption_products
	Quantity     int       `json:"quantity" gorm:"default:1"`
	CreatedAt    time.Time `json:"created_at"`

	// Para eager loading
	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductId;references:Id"`
}

// TableName define o nome da tabela no banco de dados
func (StaffAttendanceConsumption) TableName() string {
	return "staff_attendance_consumptions"
}

// StaffConsumptionProduct representa um produto disponível para consumo dos funcionários
// Diferente de Product (cardápio), esses são para controle interno
type StaffConsumptionProduct struct {
	Id             uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id" gorm:"type:uuid;not null;index"`
	ProjectId      uuid.UUID  `json:"project_id" gorm:"type:uuid;not null;index"`
	Name           string     `json:"name" gorm:"not null"`
	Category       string     `json:"category,omitempty"` // bebida, comida, etc
	UnitCost       *float64   `json:"unit_cost,omitempty"` // custo unitário para controle
	Active         bool       `json:"active" gorm:"default:true"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName define o nome da tabela no banco de dados
func (StaffConsumptionProduct) TableName() string {
	return "staff_consumption_products"
}

// ==================== DTOs ====================

// StaffAttendanceWithDetails é um DTO completo para exibição
type StaffAttendanceWithDetails struct {
	StaffAttendance
	ClientName   string                   `json:"client_name"`
	TransportIn  *TransportDetail         `json:"transport_in,omitempty"`
	TransportOut *TransportDetail         `json:"transport_out,omitempty"`
	Items        []ConsumptionDetail      `json:"items,omitempty"`
}

// TransportDetail detalha um transporte
type TransportDetail struct {
	Mode         string  `json:"mode"`
	RideWithName string  `json:"ride_with_name,omitempty"`
	BusTickets   *int    `json:"bus_tickets,omitempty"`
	UberCost     *float64 `json:"uber_cost,omitempty"`
}

// ConsumptionDetail detalha um consumo
type ConsumptionDetail struct {
	ProductId   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
}

// CreateStaffAttendanceRequest é o payload para criar presença
type CreateStaffAttendanceRequest struct {
	ClientId       uuid.UUID                `json:"client_id" binding:"required"`
	AttendanceDate string                   `json:"attendance_date" binding:"required"` // "2026-01-15"
	Shift          string                   `json:"shift" binding:"required"`           // "almoco" | "noite" | "evento"
	TransportIn    *TransportInput          `json:"transport_in,omitempty"`
	TransportOut   *TransportInput          `json:"transport_out,omitempty"`
	Consumptions   []ConsumptionInput       `json:"consumptions,omitempty"`
}

// TransportInput é o input de transporte
type TransportInput struct {
	Mode           string     `json:"mode,omitempty"`
	RideWithClient *uuid.UUID `json:"ride_with_client,omitempty"`
	BusTickets     *int       `json:"bus_tickets,omitempty"`
	UberCost       *float64   `json:"uber_cost,omitempty"`
}

// ConsumptionInput é o input de consumo
type ConsumptionInput struct {
	ProductId uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity"`
}
