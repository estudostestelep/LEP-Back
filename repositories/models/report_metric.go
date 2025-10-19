package models

import (
	"time"

	"github.com/google/uuid"
)

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
