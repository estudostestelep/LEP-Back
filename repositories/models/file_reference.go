package models

import (
	"time"

	"github.com/google/uuid"
)

// FileReference armazena metadados de arquivos com deduplicação via hash
// Uma imagem pode ser usada por múltiplas entidades (Product, Category, etc)
type FileReference struct {
	Id             uuid.UUID  `gorm:"primaryKey" json:"id"`
	OrganizationId uuid.UUID  `json:"organization_id" gorm:"not null;index"`
	ProjectId      uuid.UUID  `json:"project_id" gorm:"not null;index"`

	// Hash SHA-256 do arquivo - único por org/proj
	FileHash       string     `json:"file_hash" gorm:"not null;uniqueIndex:idx_file_hash_org_proj"`

	// Caminho no storage (local ou GCS)
	FilePath       string     `json:"file_path" gorm:"not null"`

	// Metadados do arquivo
	FileSize       int64      `json:"file_size" gorm:"not null"`
	Category       string     `json:"category" gorm:"not null"` // "products", "categories", "menus", "users", "banners"
	MimeType       string     `json:"mime_type" gorm:"not null"`

	// Contador desnormalizado de referências (performance)
	// Incrementa quando EntityFileReference é criado
	// Decrementa quando EntityFileReference é deletado
	ReferenceCount int        `json:"reference_count" gorm:"default:1;not null"`

	// Rastreamento
	CreatedAt      time.Time  `json:"created_at"`
	LastAccessedAt *time.Time `json:"last_accessed_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"` // Soft delete para cleanup posterior

	// Índice composto para evitar duplicatas
	_ struct{} `gorm:"uniqueIndex:idx_file_hash_org_proj,where:deleted_at IS NULL"`
}

// TableName especifica o nome da tabela
func (FileReference) TableName() string {
	return "file_references"
}

// EntityFileReference rastreia qual entidade usa qual imagem
// Relacionamento polimórfico entre FileReference e múltiplas entidades
type EntityFileReference struct {
	Id         uuid.UUID  `gorm:"primaryKey" json:"id"`
	FileId     uuid.UUID  `json:"file_id" gorm:"not null;index;type:uuid;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`

	// Tipo de entidade que usa a imagem
	EntityType string     `json:"entity_type" gorm:"not null;index"` // "product", "category", "menu", "subcategory", "tag"

	// ID da entidade
	EntityId   uuid.UUID  `json:"entity_id" gorm:"not null;index:idx_entity_ref"`

	// Campo da entidade que armazena a imagem
	EntityField string     `json:"entity_field" gorm:"not null"` // "image_url", "photo", "image"

	// Rastreamento
	CreatedAt  time.Time  `json:"created_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`

	// Relacionamento com FileReference
	FileReference FileReference `gorm:"foreignKey:FileId;references:Id" json:"-"`

	// Índices para queries
	// Garantir que uma entidade não tenha múltiplas imagens no mesmo campo
	_ struct{} `gorm:"uniqueIndex:idx_entity_field_unique,where:deleted_at IS NULL"`
	_ struct{} `gorm:"index:idx_entity_ref"`
}

// TableName especifica o nome da tabela
func (EntityFileReference) TableName() string {
	return "entity_file_references"
}
