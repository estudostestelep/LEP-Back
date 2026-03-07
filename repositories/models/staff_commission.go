package models

import (
	"time"

	"github.com/google/uuid"
)

// StaffDailyCommission representa a comissão de um dia/turno
type StaffDailyCommission struct {
	Id              uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrganizationId  uuid.UUID  `json:"organization_id" gorm:"type:uuid;not null;index"`
	ProjectId       uuid.UUID  `json:"project_id" gorm:"type:uuid;not null;index"`
	CommissionDate  time.Time  `json:"commission_date" gorm:"not null;index"`
	Shift           string     `json:"shift" gorm:"not null"` // "almoco" | "noite"
	CommissionValue float64    `json:"commission_value" gorm:"not null"` // valor da comissão
	Revenue         float64    `json:"revenue" gorm:"not null"`          // faturamento do turno
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName define o nome da tabela no banco de dados
func (StaffDailyCommission) TableName() string {
	return "staff_daily_commissions"
}

// ==================== DTOs ====================

// CreateCommissionRequest é o payload para criar comissão
type CreateCommissionRequest struct {
	CommissionDate  string  `json:"commission_date" binding:"required"` // "2026-01-15"
	Shift           string  `json:"shift" binding:"required"`           // "almoco" | "noite"
	CommissionValue float64 `json:"commission_value" binding:"required"`
	Revenue         float64 `json:"revenue" binding:"required"`
}

// PaymentReportRequest é o payload para gerar relatório de pagamentos
type PaymentReportRequest struct {
	StartDate string `json:"start_date" binding:"required"` // "2026-01-01"
	EndDate   string `json:"end_date" binding:"required"`   // "2026-01-31"
}

// PaymentReportRow representa uma linha do relatório de pagamentos
type PaymentReportRow struct {
	Date          string  `json:"date"`
	Weekday       string  `json:"weekday"`
	Shift         string  `json:"shift"`
	ClientId      uuid.UUID `json:"client_id"`
	ClientName    string  `json:"client_name"`
	Sector        string  `json:"sector,omitempty"`
	DailyRate     float64 `json:"daily_rate"`      // diária fixa
	Consumption   float64 `json:"consumption"`     // valor do consumo
	Commission    float64 `json:"commission"`      // comissão individual
	TransportCost float64 `json:"transport_cost"`  // custo de transporte
	Total         float64 `json:"total"`           // total a pagar
}

// PaymentReport é o relatório completo de pagamentos
type PaymentReport struct {
	StartDate       string             `json:"start_date"`
	EndDate         string             `json:"end_date"`
	Rows            []PaymentReportRow `json:"rows"`
	TotalDailyRate  float64            `json:"total_daily_rate"`
	TotalConsumption float64           `json:"total_consumption"`
	TotalCommission float64            `json:"total_commission"`
	TotalTransport  float64            `json:"total_transport"`
	GrandTotal      float64            `json:"grand_total"`
}

// StaffPaymentSummary resume os pagamentos de um funcionário
type StaffPaymentSummary struct {
	ClientId        uuid.UUID `json:"client_id"`
	ClientName      string    `json:"client_name"`
	TotalDays       int       `json:"total_days"`
	TotalDailyRate  float64   `json:"total_daily_rate"`
	TotalConsumption float64  `json:"total_consumption"`
	TotalCommission float64   `json:"total_commission"`
	TotalTransport  float64   `json:"total_transport"`
	GrandTotal      float64   `json:"grand_total"`
}

// CommissionSummary resume as comissões de um período
type CommissionSummary struct {
	StartDate      string  `json:"start_date"`
	EndDate        string  `json:"end_date"`
	TotalRevenue   float64 `json:"total_revenue"`
	TotalCommission float64 `json:"total_commission"`
	DaysCount      int     `json:"days_count"`
	AverageRevenue float64 `json:"average_revenue"`
}
