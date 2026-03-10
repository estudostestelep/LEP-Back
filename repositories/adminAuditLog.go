package repositories

import (
	"lep/repositories/models"
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AdminAuditLogRepository - Repository para logs de auditoria administrativa
// IMPORTANTE: Este repository é READ-ONLY após criação
// Não possui métodos Update ou Delete para garantir integridade dos logs
type AdminAuditLogRepository struct {
	db *gorm.DB
}

// IAdminAuditLogRepository - Interface do repository
type IAdminAuditLogRepository interface {
	// Create - Apenas para uso interno (não exposto via API)
	Create(log *models.AdminAuditLog) error

	// Read operations (expostas via API)
	ListWithFilters(filters models.AdminAuditLogFilters) (*models.AdminAuditLogPaginatedResponse, error)
	GetById(id uuid.UUID) (*models.AdminAuditLog, error)
	CountAll() (int64, error)

	// Maintenance operations (apenas para Master Admin)
	DeleteOlderThan(days int) (int64, error)
}

// NewAdminAuditLogRepository - Construtor do repository
func NewAdminAuditLogRepository(db *gorm.DB) IAdminAuditLogRepository {
	return &AdminAuditLogRepository{db: db}
}

// Create - Cria um novo log de auditoria (uso interno apenas)
func (r *AdminAuditLogRepository) Create(log *models.AdminAuditLog) error {
	if log.Id == uuid.Nil {
		log.Id = uuid.New()
	}
	return r.db.Create(log).Error
}

// ListWithFilters - Lista logs com filtros e paginação
func (r *AdminAuditLogRepository) ListWithFilters(filters models.AdminAuditLogFilters) (*models.AdminAuditLogPaginatedResponse, error) {
	var logs []models.AdminAuditLog
	var total int64

	// Query base
	query := r.db.Model(&models.AdminAuditLog{})

	// Aplicar filtros
	if filters.StartDate != nil {
		query = query.Where("created_at >= ?", filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("created_at <= ?", filters.EndDate)
	}
	if filters.ActorId != nil {
		query = query.Where("actor_id = ?", filters.ActorId)
	}
	if filters.ActorEmail != "" {
		query = query.Where("actor_email ILIKE ?", "%"+filters.ActorEmail+"%")
	}
	if filters.Action != "" {
		query = query.Where("action = ?", filters.Action)
	}
	if filters.EntityType != "" {
		query = query.Where("entity_type = ?", filters.EntityType)
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

	return &models.AdminAuditLogPaginatedResponse{
		Data:       logs,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetById - Busca um log específico pelo ID
func (r *AdminAuditLogRepository) GetById(id uuid.UUID) (*models.AdminAuditLog, error) {
	var log models.AdminAuditLog
	err := r.db.Where("id = ?", id).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// CountAll - Conta total de logs (para estatísticas)
func (r *AdminAuditLogRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&models.AdminAuditLog{}).Count(&count).Error
	return count, err
}

// DeleteOlderThan - Deleta logs mais antigos que X dias (manutenção)
// Retorna a quantidade de logs deletados
func (r *AdminAuditLogRepository) DeleteOlderThan(days int) (int64, error) {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	result := r.db.Where("created_at < ?", cutoffDate).Delete(&models.AdminAuditLog{})
	return result.RowsAffected, result.Error
}
