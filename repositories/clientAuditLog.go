package repositories

import (
	"lep/repositories/models"
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ClientAuditLogRepository - Repository para logs de auditoria de cliente
type ClientAuditLogRepository struct {
	db *gorm.DB
}

// IClientAuditLogRepository - Interface do repository
type IClientAuditLogRepository interface {
	// Logs
	Create(log *models.ClientAuditLog) error
	ListByProject(orgId, projectId uuid.UUID, filters models.ClientAuditLogFilters) (*models.ClientAuditLogPaginatedResponse, error)
	GetById(id uuid.UUID) (*models.ClientAuditLog, error)
	CountByOrganization(orgId uuid.UUID) (int64, error)
	CleanupOldLogs(orgId uuid.UUID, retentionDays int, maxLogs int) (int64, error)

	// Config
	GetConfig(orgId uuid.UUID) (*models.ClientAuditConfig, error)
	CreateConfig(config *models.ClientAuditConfig) error
	UpdateConfig(config *models.ClientAuditConfig) error
	GetAvailableModules() []struct {
		Code        string `json:"code"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
}

// NewClientAuditLogRepository - Construtor do repository
func NewClientAuditLogRepository(db *gorm.DB) IClientAuditLogRepository {
	return &ClientAuditLogRepository{db: db}
}

// Create - Cria um novo log de auditoria de cliente
func (r *ClientAuditLogRepository) Create(log *models.ClientAuditLog) error {
	if log.Id == uuid.Nil {
		log.Id = uuid.New()
	}
	return r.db.Create(log).Error
}

// ListByProject - Lista logs de uma organização/projeto com filtros
func (r *ClientAuditLogRepository) ListByProject(orgId, projectId uuid.UUID, filters models.ClientAuditLogFilters) (*models.ClientAuditLogPaginatedResponse, error) {
	var logs []models.ClientAuditLog
	var total int64

	// Query base - filtrar por org e projeto
	query := r.db.Model(&models.ClientAuditLog{}).
		Where("organization_id = ? AND project_id = ?", orgId, projectId)

	// Aplicar filtros
	if filters.StartDate != nil {
		query = query.Where("created_at >= ?", filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("created_at <= ?", filters.EndDate)
	}
	if filters.UserId != nil {
		query = query.Where("user_id = ?", filters.UserId)
	}
	if filters.UserEmail != "" {
		query = query.Where("user_email ILIKE ?", "%"+filters.UserEmail+"%")
	}
	if filters.Action != "" {
		query = query.Where("action = ?", filters.Action)
	}
	if filters.EntityType != "" {
		query = query.Where("entity_type = ?", filters.EntityType)
	}
	if filters.ModuleCode != "" {
		query = query.Where("module_code = ?", filters.ModuleCode)
	}

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Paginação
	page := filters.Page
	if page < 1 {
		page = 1
	}
	pageSize := filters.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	// Buscar registros ordenados por data (mais recente primeiro)
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, err
	}

	// Calcular total de páginas
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &models.ClientAuditLogPaginatedResponse{
		Data:       logs,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetById - Busca um log específico pelo ID
func (r *ClientAuditLogRepository) GetById(id uuid.UUID) (*models.ClientAuditLog, error) {
	var log models.ClientAuditLog
	err := r.db.Where("id = ?", id).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// CountByOrganization - Conta total de logs de uma organização
func (r *ClientAuditLogRepository) CountByOrganization(orgId uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.ClientAuditLog{}).Where("organization_id = ?", orgId).Count(&count).Error
	return count, err
}

// CleanupOldLogs - Remove logs antigos respeitando retention e limite
func (r *ClientAuditLogRepository) CleanupOldLogs(orgId uuid.UUID, retentionDays int, maxLogs int) (int64, error) {
	var deleted int64

	// 1. Deletar logs mais antigos que retentionDays
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)
	result := r.db.Where("organization_id = ? AND created_at < ?", orgId, cutoffDate).
		Delete(&models.ClientAuditLog{})
	if result.Error != nil {
		return 0, result.Error
	}
	deleted = result.RowsAffected

	// 2. Se ainda excede maxLogs, deletar os mais antigos
	if maxLogs > 0 {
		var currentCount int64
		r.db.Model(&models.ClientAuditLog{}).Where("organization_id = ?", orgId).Count(&currentCount)

		if currentCount > int64(maxLogs) {
			excessCount := currentCount - int64(maxLogs)

			// Encontrar IDs dos logs mais antigos a serem deletados
			var oldestLogs []models.ClientAuditLog
			r.db.Where("organization_id = ?", orgId).
				Order("created_at ASC").
				Limit(int(excessCount)).
				Select("id").
				Find(&oldestLogs)

			if len(oldestLogs) > 0 {
				var ids []uuid.UUID
				for _, log := range oldestLogs {
					ids = append(ids, log.Id)
				}
				result := r.db.Where("id IN ?", ids).Delete(&models.ClientAuditLog{})
				if result.Error == nil {
					deleted += result.RowsAffected
				}
			}
		}
	}

	return deleted, nil
}

// GetConfig - Obtém configuração de auditoria de uma organização
func (r *ClientAuditLogRepository) GetConfig(orgId uuid.UUID) (*models.ClientAuditConfig, error) {
	var config models.ClientAuditConfig
	err := r.db.Where("organization_id = ?", orgId).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Não configurado ainda
		}
		return nil, err
	}
	return &config, nil
}

// CreateConfig - Cria configuração de auditoria para uma organização
func (r *ClientAuditLogRepository) CreateConfig(config *models.ClientAuditConfig) error {
	if config.Id == uuid.Nil {
		config.Id = uuid.New()
	}
	return r.db.Create(config).Error
}

// UpdateConfig - Atualiza configuração de auditoria
func (r *ClientAuditLogRepository) UpdateConfig(config *models.ClientAuditConfig) error {
	return r.db.Save(config).Error
}

// GetAvailableModules - Retorna lista de módulos disponíveis
func (r *ClientAuditLogRepository) GetAvailableModules() []struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
} {
	return models.AvailableClientAuditModules
}
