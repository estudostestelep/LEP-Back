package models

import (
	"time"

	"github.com/google/uuid"
)

// SidebarItemBehavior define o comportamento de um item quando o usuário não tem acesso
type SidebarItemBehavior string

const (
	BehaviorShow SidebarItemBehavior = "show" // Sempre visível (desabilitado quando sem permissão)
	BehaviorLock SidebarItemBehavior = "lock" // Mostra com cadeado quando indisponível
	BehaviorHide SidebarItemBehavior = "hide" // Esconde completamente quando indisponível
)

// SidebarItemConfig representa a configuração de um item individual da sidebar
type SidebarItemConfig struct {
	ModuleCode string              `json:"module_code"`
	Behavior   SidebarItemBehavior `json:"behavior"`
}

// SidebarConfig armazena a configuração da sidebar por organização
type SidebarConfig struct {
	Id             uuid.UUID `gorm:"primaryKey" json:"id"`
	OrganizationId uuid.UUID `gorm:"not null;uniqueIndex" json:"organization_id"`

	// ItemConfigs armazena as configurações em formato JSON
	ItemConfigs string `gorm:"type:jsonb;default:'[]'" json:"item_configs"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SidebarConfigResponse é a estrutura de resposta da API
type SidebarConfigResponse struct {
	Id             uuid.UUID           `json:"id"`
	OrganizationId uuid.UUID           `json:"organization_id"`
	Items          []SidebarItemConfig `json:"items"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
}

// SidebarConfigUpdateRequest é a estrutura de requisição para atualizar a configuração
type SidebarConfigUpdateRequest struct {
	Items []SidebarItemConfig `json:"items" binding:"required"`
}
