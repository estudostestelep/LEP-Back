package repositories

import (
	"lep/repositories/models"
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type resourceAccessLog struct {
	db *gorm.DB
}

type IAccessLogRepository interface {
	Create(log *models.AccessLog) error
	GetByUserId(userId string, page, perPage int) (*models.AccessLogPaginatedResponse, error)
	GetRecentByUserId(userId string, limit int) ([]models.AccessLog, error)
	DeleteOldLogs(olderThan time.Time) error
}

func NewAccessLogRepository(db *gorm.DB) IAccessLogRepository {
	return &resourceAccessLog{db: db}
}

func (r *resourceAccessLog) Create(log *models.AccessLog) error {
	if log.Id == uuid.Nil {
		log.Id = uuid.New()
	}
	log.CreatedAt = time.Now()
	return r.db.Create(log).Error
}

func (r *resourceAccessLog) GetByUserId(userId string, page, perPage int) (*models.AccessLogPaginatedResponse, error) {
	var logs []models.AccessLog
	var total int64

	// Default pagination
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	// Count total
	err := r.db.Model(&models.AccessLog{}).Where("user_id = ?", userId).Count(&total).Error
	if err != nil {
		return nil, err
	}

	// Get paginated data
	err = r.db.Where("user_id = ?", userId).
		Order("login_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&logs).Error
	if err != nil {
		return nil, err
	}

	// Convert to response format
	data := make([]models.AccessLogResponse, len(logs))
	for i, log := range logs {
		data[i] = models.AccessLogResponse{
			Id:        log.Id,
			UserId:    log.UserId,
			IP:        log.IP,
			UserAgent: log.UserAgent,
			Location:  log.Location,
			LoginAt:   log.LoginAt,
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &models.AccessLogPaginatedResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (r *resourceAccessLog) GetRecentByUserId(userId string, limit int) ([]models.AccessLog, error) {
	var logs []models.AccessLog
	if limit < 1 {
		limit = 10
	}
	err := r.db.Where("user_id = ?", userId).
		Order("login_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// DeleteOldLogs remove logs mais antigos que a data especificada
// Útil para limpeza periódica de dados antigos
func (r *resourceAccessLog) DeleteOldLogs(olderThan time.Time) error {
	return r.db.Where("created_at < ?", olderThan).Delete(&models.AccessLog{}).Error
}
