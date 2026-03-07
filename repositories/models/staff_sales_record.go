package models

import (
	"time"

	"github.com/google/uuid"
)

// StaffSalesRecord representa um registro de venda importado do sistema externo
// Estrutura baseada no CSV "Relatorio de Vendas em Grid"
type StaffSalesRecord struct {
	Id                   uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrganizationId       uuid.UUID  `json:"organization_id" gorm:"type:uuid;not null;index"`
	ProjectId            uuid.UUID  `json:"project_id" gorm:"type:uuid;not null;index"`

	// Campos do CSV
	Quantity             int        `json:"quantity" gorm:"column:qtd"`                           // qtd
	ServiceFeeCalculated float64    `json:"service_fee_calculated" gorm:"column:vl_servico_calc"` // vl_servico_calculado
	Description          string     `json:"description" gorm:"not null"`                          // descricao
	ServiceFeeInformed   float64    `json:"service_fee_informed" gorm:"column:vl_servico_info"`   // vl_servico_informado
	TotalValue           float64    `json:"total_value" gorm:"column:vl_total"`                   // vl_total
	LaunchDateTime       time.Time  `json:"launch_date_time" gorm:"index"`                        // dt_hr_lancamento
	OriginalPrice        float64    `json:"original_price" gorm:"column:vl_preco_original"`       // vl_preco_original
	AccountingDate       time.Time  `json:"accounting_date" gorm:"not null;index"`                // dt_contabil
	EmployeeName         string     `json:"employee_name"`                                         // nome_funcionario
	Price                float64    `json:"price" gorm:"column:vl_preco"`                         // vl_preco
	SaleMode             string     `json:"sale_mode"`                                             // nome_modo_venda (mesa, etc)
	Discount             float64    `json:"discount" gorm:"column:vl_desconto"`                   // vl_desconto
	GroupName            string     `json:"group_name" gorm:"not null;index"`                     // grupo (categoria)

	// Campos de controle
	ImportedAt           time.Time  `json:"imported_at"`
	ImportBatchId        *uuid.UUID `json:"import_batch_id,omitempty" gorm:"type:uuid;index"` // lote de importação
	CreatedAt            time.Time  `json:"created_at"`
}

// TableName define o nome da tabela no banco de dados
func (StaffSalesRecord) TableName() string {
	return "staff_sales_records"
}

// StaffSalesImportBatch representa um lote de importação
type StaffSalesImportBatch struct {
	Id             uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id" gorm:"type:uuid;not null;index"`
	ProjectId      uuid.UUID  `json:"project_id" gorm:"type:uuid;not null;index"`
	FileName       string     `json:"file_name"`
	RecordsCount   int        `json:"records_count"`
	Status         string     `json:"status" gorm:"default:'pending'"` // pending, processing, completed, failed
	ErrorMessage   *string    `json:"error_message,omitempty"`
	CreatedById    *uuid.UUID `json:"created_by_id,omitempty" gorm:"type:uuid"`
	CreatedAt      time.Time  `json:"created_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
}

// TableName define o nome da tabela no banco de dados
func (StaffSalesImportBatch) TableName() string {
	return "staff_sales_import_batches"
}

// ==================== DTOs ====================

// DashboardMeta retorna metadados para filtros do dashboard
type DashboardMeta struct {
	Groups       []string   `json:"groups"`       // categorias únicas
	Descriptions []string   `json:"descriptions"` // itens únicos
	Employees    []string   `json:"employees"`    // funcionários únicos
	SaleModes    []string   `json:"sale_modes"`   // modos de venda únicos
	MinDate      *time.Time `json:"min_date"`     // data mais antiga
	MaxDate      *time.Time `json:"max_date"`     // data mais recente
}

// DashboardFilters são os filtros para consulta do dashboard
type DashboardFilters struct {
	StartDate    *string  `json:"start_date,omitempty"`    // "2026-01-01"
	EndDate      *string  `json:"end_date,omitempty"`      // "2026-01-31"
	Groups       []string `json:"groups,omitempty"`        // categorias
	Descriptions []string `json:"descriptions,omitempty"`  // itens
	Employees    []string `json:"employees,omitempty"`     // funcionários
	Weekdays     []string `json:"weekdays,omitempty"`      // dias da semana
}

// DashboardRow é uma linha agregada do dashboard
type DashboardRow struct {
	AccountingDate       string  `json:"accounting_date"`
	Weekday              string  `json:"weekday"`
	GroupName            string  `json:"group_name"`
	Description          string  `json:"description"`
	EmployeeName         string  `json:"employee_name"`
	Quantity             int     `json:"quantity"`
	ServiceFeeInformed   float64 `json:"service_fee_informed"`
	ServiceFeeCalculated float64 `json:"service_fee_calculated"`
	TotalValue           float64 `json:"total_value"`
}

// DashboardSummary é o resumo do dashboard
type DashboardSummary struct {
	Rows               []DashboardRow `json:"rows"`
	TotalQuantity      int            `json:"total_quantity"`
	TotalServiceFee    float64        `json:"total_service_fee"`
	TotalValue         float64        `json:"total_value"`
	UniqueGroups       int            `json:"unique_groups"`
	UniqueDescriptions int            `json:"unique_descriptions"`
}

// EmployeeGraphPoint representa faturamento por funcionário
type EmployeeGraphPoint struct {
	Employee string  `json:"employee"`
	Total    float64 `json:"total"`
}

// GroupGraphPoint representa faturamento por grupo
type GroupGraphPoint struct {
	Group string  `json:"group"`
	Total float64 `json:"total"`
}

// DateGraphPoint representa faturamento por data
type DateGraphPoint struct {
	Date  string  `json:"date"`
	Total float64 `json:"total"`
}

// WeekdayGraphPoint representa faturamento por dia da semana
type WeekdayGraphPoint struct {
	Weekday string  `json:"weekday"`
	Total   float64 `json:"total"`
}

// DashboardGraphs contém dados para os 4 gráficos resumidos
type DashboardGraphs struct {
	ByEmployee []EmployeeGraphPoint `json:"by_employee"` // top 10 funcionários
	ByGroup    []GroupGraphPoint    `json:"by_group"`    // top 10 grupos
	ByDate     []DateGraphPoint     `json:"by_date"`     // últimos 30 dias
	ByWeekday  []WeekdayGraphPoint  `json:"by_weekday"`  // por dia da semana
}

// StaffReportFilters são os filtros para relatório de colaboradores
type StaffReportFilters struct {
	Year      int      `json:"year" binding:"required"`
	StartDate *string  `json:"start_date,omitempty"`
	EndDate   *string  `json:"end_date,omitempty"`
	Weekdays  []string `json:"weekdays,omitempty"`
	Employees []string `json:"employees,omitempty"`
	Sectors   []string `json:"sectors,omitempty"`
}

// StaffReportMeta retorna metadados para filtros do relatório
type StaffReportMeta struct {
	Years     []int    `json:"years"`
	Employees []string `json:"employees"`
	Sectors   []string `json:"sectors"`
}

// StaffReportRow é uma linha do relatório de colaboradores
type StaffReportRow struct {
	Date       string  `json:"date"`
	Weekday    string  `json:"weekday"`
	Shift      string  `json:"shift"`
	ClientName string  `json:"client_name"`
	Sector     string  `json:"sector"`
	DailyRate  float64 `json:"daily_rate"`
	Consumption float64 `json:"consumption"`
	Commission float64 `json:"commission"`
	Total      float64 `json:"total"`
}

// StaffReportSummary é o relatório completo de colaboradores
type StaffReportSummary struct {
	Rows            []StaffReportRow `json:"rows"`
	TotalDailyRate  float64          `json:"total_daily_rate"`
	TotalConsumption float64         `json:"total_consumption"`
	TotalCommission float64          `json:"total_commission"`
	GrandTotal      float64          `json:"grand_total"`
}

// ImportCSVRequest é o payload para importação de CSV
type ImportCSVRequest struct {
	FileContent string `json:"file_content" binding:"required"` // conteúdo base64 do CSV
	FileName    string `json:"file_name" binding:"required"`
}

// ImportCSVResponse é a resposta da importação
type ImportCSVResponse struct {
	BatchId      uuid.UUID `json:"batch_id"`
	RecordsCount int       `json:"records_count"`
	Status       string    `json:"status"`
	Errors       []string  `json:"errors,omitempty"`
}
