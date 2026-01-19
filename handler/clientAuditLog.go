package handler

import (
	"encoding/json"
	"fmt"
	"lep/repositories"
	"lep/repositories/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// ClientAuditLogHandler - Handler para logs de auditoria de cliente
type ClientAuditLogHandler struct {
	repo repositories.IClientAuditLogRepository
}

// IClientAuditLogHandler - Interface do handler
type IClientAuditLogHandler interface {
	// Métodos públicos (API)
	ListLogs(orgId, projectId uuid.UUID, filters models.ClientAuditLogFilters) (*models.ClientAuditLogPaginatedResponse, error)
	GetLogById(id uuid.UUID) (*models.ClientAuditLog, error)

	// Configuração
	GetConfig(orgId uuid.UUID) (*models.ClientAuditConfig, error)
	SaveConfig(config *models.ClientAuditConfig) error
	GetAvailableModules() []struct {
		Code        string `json:"code"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	// Verificação de ativação
	ShouldLog(orgId uuid.UUID, moduleCode string) bool

	// Métodos internos (interceptação de operações) - genérico
	LogAction(orgId, projectId uuid.UUID, userId *uuid.UUID, userEmail, action, entityType string, entityId uuid.UUID, moduleCode string, oldValues, newValues interface{}, changedFields []string, description, ipAddress string) error
}

// NewClientAuditLogHandler - Construtor do handler
func NewClientAuditLogHandler(repo repositories.IClientAuditLogRepository) IClientAuditLogHandler {
	return &ClientAuditLogHandler{repo: repo}
}

// ListLogs - Lista logs de uma organização/projeto com filtros
func (h *ClientAuditLogHandler) ListLogs(orgId, projectId uuid.UUID, filters models.ClientAuditLogFilters) (*models.ClientAuditLogPaginatedResponse, error) {
	return h.repo.ListByProject(orgId, projectId, filters)
}

// GetLogById - Busca um log específico pelo ID
func (h *ClientAuditLogHandler) GetLogById(id uuid.UUID) (*models.ClientAuditLog, error) {
	return h.repo.GetById(id)
}

// GetConfig - Obtém configuração de auditoria de uma organização
func (h *ClientAuditLogHandler) GetConfig(orgId uuid.UUID) (*models.ClientAuditConfig, error) {
	config, err := h.repo.GetConfig(orgId)
	if err != nil {
		return nil, err
	}

	// Se não existe, retornar configuração padrão (desabilitada)
	if config == nil {
		return &models.ClientAuditConfig{
			OrganizationId: orgId,
			Enabled:        false,
			MaxLogsStored:  10000,
			RetentionDays:  90,
			EnabledModules: pq.StringArray{},
		}, nil
	}

	return config, nil
}

// SaveConfig - Salva (cria ou atualiza) configuração de auditoria
func (h *ClientAuditLogHandler) SaveConfig(config *models.ClientAuditConfig) error {
	existing, err := h.repo.GetConfig(config.OrganizationId)
	if err != nil {
		return err
	}

	if existing == nil {
		// Criar nova configuração
		return h.repo.CreateConfig(config)
	}

	// Atualizar configuração existente
	config.Id = existing.Id
	return h.repo.UpdateConfig(config)
}

// GetAvailableModules - Retorna lista de módulos disponíveis para auditoria
func (h *ClientAuditLogHandler) GetAvailableModules() []struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
} {
	return h.repo.GetAvailableModules()
}

// ShouldLog - Verifica se deve logar ações para um módulo específico
func (h *ClientAuditLogHandler) ShouldLog(orgId uuid.UUID, moduleCode string) bool {
	config, err := h.repo.GetConfig(orgId)
	if err != nil || config == nil || !config.Enabled {
		return false
	}

	// Verificar se o módulo está na lista de ativos
	for _, m := range config.EnabledModules {
		if m == moduleCode {
			return true
		}
	}
	return false
}

// LogAction - Registra uma ação no log de auditoria de cliente
func (h *ClientAuditLogHandler) LogAction(
	orgId, projectId uuid.UUID,
	userId *uuid.UUID,
	userEmail, action, entityType string,
	entityId uuid.UUID,
	moduleCode string,
	oldValues, newValues interface{},
	changedFields []string,
	description, ipAddress string,
) error {
	// Verificar se deve logar
	if !h.ShouldLog(orgId, moduleCode) {
		return nil // Módulo não ativo, não logar
	}

	// Converter valores para JSON
	var oldValuesJSON, newValuesJSON []byte
	if oldValues != nil {
		oldValuesJSON, _ = json.Marshal(oldValues)
	}
	if newValues != nil {
		newValuesJSON, _ = json.Marshal(newValues)
	}

	log := &models.ClientAuditLog{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projectId,
		UserId:         userId,
		UserEmail:      userEmail,
		Action:         action,
		EntityType:     entityType,
		EntityId:       entityId,
		ModuleCode:     moduleCode,
		OldValues:      oldValuesJSON,
		NewValues:      newValuesJSON,
		ChangedFields:  changedFields,
		Description:    description,
		IpAddress:      ipAddress,
	}

	if err := h.repo.Create(log); err != nil {
		fmt.Printf("❌ Erro ao criar log de auditoria de cliente: %v\n", err)
		return err
	}

	fmt.Printf("✅ Log de auditoria de cliente registrado: %s em %s (%s)\n", action, entityType, entityId)
	return nil
}
