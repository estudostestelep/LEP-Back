package models

import (
	"time"

	"github.com/google/uuid"
)

// PlanChangeRequest representa uma solicitação de mudança de plano
type PlanChangeRequest struct {
	Id             uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrganizationId uuid.UUID  `gorm:"type:uuid;not null;index" json:"organization_id"`
	RequestedBy    uuid.UUID  `gorm:"type:uuid;not null" json:"requested_by"` // ID do usuário que solicitou

	// Plano atual e plano desejado
	CurrentPackageId uuid.UUID  `gorm:"type:uuid" json:"current_package_id"`
	RequestedPackageId uuid.UUID  `gorm:"type:uuid;not null" json:"requested_package_id"`

	// Informações do pacote (armazenadas para histórico)
	CurrentPackageName  string `json:"current_package_name,omitempty"`
	RequestedPackageName string `json:"requested_package_name"`

	// Justificativa e informações adicionais
	Reason       string `gorm:"type:text" json:"reason,omitempty"`         // Motivo da mudança
	Notes        string `gorm:"type:text" json:"notes,omitempty"`          // Observações adicionais

	// Status da solicitação
	Status       string    `gorm:"type:varchar(50);default:'pending'" json:"status"` // pending, approved, rejected, cancelled

	// Informações de aprovação/rejeição
	ReviewedBy   *uuid.UUID `gorm:"type:uuid" json:"reviewed_by,omitempty"`   // ID do admin que aprovou/rejeitou
	ReviewedAt   *time.Time `json:"reviewed_at,omitempty"`
	ReviewNotes  string     `gorm:"type:text" json:"review_notes,omitempty"`  // Comentários do admin

	// Preferências de billing (opcional)
	RequestedBillingCycle string   `json:"requested_billing_cycle,omitempty"` // monthly, yearly

	// Timestamps
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName define o nome da tabela no banco de dados
func (PlanChangeRequest) TableName() string {
	return "plan_change_requests"
}

// Status constants
const (
	PlanChangeStatusPending   = "pending"
	PlanChangeStatusApproved  = "approved"
	PlanChangeStatusRejected  = "rejected"
	PlanChangeStatusCancelled = "cancelled"
)

// IsPending verifica se a solicitação está pendente
func (p *PlanChangeRequest) IsPending() bool {
	return p.Status == PlanChangeStatusPending
}

// IsApproved verifica se a solicitação foi aprovada
func (p *PlanChangeRequest) IsApproved() bool {
	return p.Status == PlanChangeStatusApproved
}

// IsRejected verifica se a solicitação foi rejeitada
func (p *PlanChangeRequest) IsRejected() bool {
	return p.Status == PlanChangeStatusRejected
}

// IsCancelled verifica se a solicitação foi cancelada
func (p *PlanChangeRequest) IsCancelled() bool {
	return p.Status == PlanChangeStatusCancelled
}
