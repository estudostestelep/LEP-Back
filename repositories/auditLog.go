package repositories

import (
	"lep/repositories/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditLogRepository struct {
	db *gorm.DB
}

type IAuditLogsRepository interface {
	GetAuditLog(OrganizationId, projectId uuid.UUID) ([]models.AuditLog, error)
	CreateAuditLog(AuditLog *models.AuditLog) error
	UpdateAuditLog(id int, updatedAuditLog *models.AuditLog) error
	DeleteAuditLog(id int) error
}

func NewConnAuditLog(db *gorm.DB) IAuditLogsRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) CreateAuditLog(log *models.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *AuditLogRepository) GetAuditLog(OrganizationId, projectId uuid.UUID) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.Where("org_id = ? AND project_id = ?", OrganizationId, projectId).Find(&logs).Error
	return logs, err
}

func (r *AuditLogRepository) UpdateAuditLog(id int, updatedAuditLog *models.AuditLog) error {
	return r.db.Model(&models.AuditLog{}).Where("id = ?", id).Updates(updatedAuditLog).Error
}

func (r *AuditLogRepository) DeleteAuditLog(id int) error {
	return r.db.Delete(&models.AuditLog{}, id).Error
}
