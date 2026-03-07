package models

import (
	"time"

	"github.com/google/uuid"
)

// StaffStockItem representa um item de estoque operacional
type StaffStockItem struct {
	Id             uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id" gorm:"type:uuid;not null;index"`
	ProjectId      uuid.UUID  `json:"project_id" gorm:"type:uuid;not null;index"`
	Name           string     `json:"name" gorm:"not null"`
	Category       string     `json:"category,omitempty"`     // categoria do item
	Storage        string     `json:"storage,omitempty"`      // local de armazenamento
	StockMin       *int       `json:"stock_min,omitempty"`    // estoque mínimo
	StockMax       *int       `json:"stock_max,omitempty"`    // estoque máximo
	Sector         string     `json:"sector,omitempty"`       // setor (cozinha, bar, etc)
	WhereToBuy     string     `json:"where_to_buy,omitempty"` // onde comprar
	Notes          string     `json:"notes,omitempty"`        // observações
	Active         bool       `json:"active" gorm:"default:true"`
	Order          int        `json:"order" gorm:"default:0"` // ordem de exibição
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName define o nome da tabela no banco de dados
func (StaffStockItem) TableName() string {
	return "staff_stock_items"
}

// StaffStockRecord representa um registro/snapshot de inventário
type StaffStockRecord struct {
	Id             uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id" gorm:"type:uuid;not null;index"`
	ProjectId      uuid.UUID  `json:"project_id" gorm:"type:uuid;not null;index"`
	RecordDate     time.Time  `json:"record_date" gorm:"not null;index"`
	Sector         string     `json:"sector" gorm:"not null"` // setor inventariado
	CreatedById    *uuid.UUID `json:"created_by_id,omitempty" gorm:"type:uuid"` // quem fez o registro
	Notes          string     `json:"notes,omitempty"`        // itens extras ou observações
	PdfUrl         *string    `json:"pdf_url,omitempty"`      // URL do PDF gerado
	EmailSentAt    *time.Time `json:"email_sent_at,omitempty"` // quando o email foi enviado
	CreatedAt      time.Time  `json:"created_at"`

	// Relacionamentos
	CreatedBy *Client                `json:"created_by,omitempty" gorm:"foreignKey:CreatedById;references:Id"`
	Items     []StaffStockRecordItem `json:"items,omitempty" gorm:"foreignKey:StockRecordId"`
}

// TableName define o nome da tabela no banco de dados
func (StaffStockRecord) TableName() string {
	return "staff_stock_records"
}

// StaffStockRecordItem representa um item do registro de estoque
type StaffStockRecordItem struct {
	Id            uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	StockRecordId uuid.UUID `json:"stock_record_id" gorm:"type:uuid;not null;index"`
	StockItemId   uuid.UUID `json:"stock_item_id" gorm:"type:uuid;not null;index"`
	CurrentStock  int       `json:"current_stock" gorm:"not null"` // quantidade atual
	ToBuy         int       `json:"to_buy" gorm:"default:0"`       // calculado: max - current se current < min
	CreatedAt     time.Time `json:"created_at"`

	// Para eager loading
	StockItem *StaffStockItem `json:"stock_item,omitempty" gorm:"foreignKey:StockItemId;references:Id"`
}

// TableName define o nome da tabela no banco de dados
func (StaffStockRecordItem) TableName() string {
	return "staff_stock_record_items"
}

// ==================== DTOs ====================

// StaffStockItemWithCurrentStock inclui o estoque atual (do último registro)
type StaffStockItemWithCurrentStock struct {
	StaffStockItem
	CurrentStock  *int       `json:"current_stock,omitempty"`
	LastRecordAt  *time.Time `json:"last_record_at,omitempty"`
}

// CreateStockRecordRequest é o payload para criar um registro de estoque
type CreateStockRecordRequest struct {
	Sector string                    `json:"sector" binding:"required"`
	Items  []StockRecordItemInput    `json:"items" binding:"required"`
	Notes  string                    `json:"notes,omitempty"` // itens extras
}

// StockRecordItemInput é o input de um item do registro
type StockRecordItemInput struct {
	StockItemId  uuid.UUID `json:"stock_item_id" binding:"required"`
	CurrentStock int       `json:"current_stock" binding:"required"`
}

// StockRecordWithDetails é o DTO completo do registro
type StockRecordWithDetails struct {
	StaffStockRecord
	CreatedByName string                     `json:"created_by_name,omitempty"`
	Items         []StockRecordItemWithName  `json:"items"`
	TotalToBuy    int                        `json:"total_to_buy"` // soma de todos os to_buy
}

// StockRecordItemWithName inclui o nome do item
type StockRecordItemWithName struct {
	StaffStockRecordItem
	ItemName    string `json:"item_name"`
	ItemMin     *int   `json:"item_min,omitempty"`
	ItemMax     *int   `json:"item_max,omitempty"`
	WhereToBuy  string `json:"where_to_buy,omitempty"`
}

// ShoppingListItem representa um item da lista de compras
type ShoppingListItem struct {
	ItemId       uuid.UUID `json:"item_id"`
	ItemName     string    `json:"item_name"`
	Category     string    `json:"category"`
	CurrentStock int       `json:"current_stock"`
	StockMin     int       `json:"stock_min"`
	StockMax     int       `json:"stock_max"`
	ToBuy        int       `json:"to_buy"`
	WhereToBuy   string    `json:"where_to_buy"`
	Notes        string    `json:"notes"`
}

// ShoppingList é a lista de compras completa
type ShoppingList struct {
	RecordId   uuid.UUID          `json:"record_id"`
	RecordDate time.Time          `json:"record_date"`
	Sector     string             `json:"sector"`
	Items      []ShoppingListItem `json:"items"`
	ExtraNotes string             `json:"extra_notes,omitempty"`
}

// SectorSummary resume os setores disponíveis
type SectorSummary struct {
	Sector     string `json:"sector"`
	ItemCount  int    `json:"item_count"`
}
