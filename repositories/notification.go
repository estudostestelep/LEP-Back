package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

type INotificationRepository interface {
	// NotificationLog
	CreateNotificationLog(log *models.NotificationLog) error
	GetNotificationLogById(id uuid.UUID) (*models.NotificationLog, error)
	UpdateNotificationLogStatus(externalId, status string, deliveredAt *time.Time) error
	GetNotificationLogsByProject(orgId, projectId uuid.UUID, limit int) ([]models.NotificationLog, error)

	// NotificationInbound
	CreateNotificationInbound(inbound *models.NotificationInbound) error
	GetUnprocessedInbound(orgId, projectId uuid.UUID) ([]models.NotificationInbound, error)
	MarkInboundAsProcessed(id uuid.UUID) error

	// NotificationEvent
	CreateNotificationEvent(event *models.NotificationEvent) error
	GetUnprocessedEvents(orgId, projectId uuid.UUID) ([]models.NotificationEvent, error)
	MarkEventAsProcessed(id uuid.UUID) error

	// NotificationConfig
	GetNotificationConfigByEvent(orgId, projectId uuid.UUID, eventType string) (*models.NotificationConfig, error)
	CreateOrUpdateNotificationConfig(config *models.NotificationConfig) error

	// NotificationTemplate
	GetNotificationTemplateByChannel(orgId, projectId uuid.UUID, channel string) (*models.NotificationTemplate, error)
	CreateNotificationTemplate(template *models.NotificationTemplate) error
	UpdateNotificationTemplate(template *models.NotificationTemplate) error
	GetNotificationTemplatesByProject(orgId, projectId uuid.UUID) ([]models.NotificationTemplate, error)
}

func NewNotificationRepository(db *gorm.DB) INotificationRepository {
	return &NotificationRepository{db: db}
}

// === NotificationLog ===

func (r *NotificationRepository) CreateNotificationLog(log *models.NotificationLog) error {
	return r.db.Create(log).Error
}

func (r *NotificationRepository) GetNotificationLogById(id uuid.UUID) (*models.NotificationLog, error) {
	var log models.NotificationLog
	err := r.db.First(&log, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *NotificationRepository) UpdateNotificationLogStatus(externalId, status string, deliveredAt *time.Time) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	if deliveredAt != nil {
		updates["delivered_at"] = deliveredAt
	}

	return r.db.Model(&models.NotificationLog{}).
		Where("external_id = ?", externalId).
		Updates(updates).Error
}

func (r *NotificationRepository) GetNotificationLogsByProject(orgId, projectId uuid.UUID, limit int) ([]models.NotificationLog, error) {
	var logs []models.NotificationLog
	query := r.db.Where("organization_id = ? AND project_id = ?", orgId, projectId).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&logs).Error
	return logs, err
}

// === NotificationInbound ===

func (r *NotificationRepository) CreateNotificationInbound(inbound *models.NotificationInbound) error {
	return r.db.Create(inbound).Error
}

func (r *NotificationRepository) GetUnprocessedInbound(orgId, projectId uuid.UUID) ([]models.NotificationInbound, error) {
	var inbound []models.NotificationInbound
	err := r.db.Where("organization_id = ? AND project_id = ? AND processed = false", orgId, projectId).
		Order("created_at ASC").Find(&inbound).Error
	return inbound, err
}

func (r *NotificationRepository) MarkInboundAsProcessed(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.NotificationInbound{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"processed":    true,
			"processed_at": now,
		}).Error
}

// === NotificationEvent ===

func (r *NotificationRepository) CreateNotificationEvent(event *models.NotificationEvent) error {
	return r.db.Create(event).Error
}

func (r *NotificationRepository) GetUnprocessedEvents(orgId, projectId uuid.UUID) ([]models.NotificationEvent, error) {
	var events []models.NotificationEvent
	err := r.db.Where("organization_id = ? AND project_id = ? AND processed = false", orgId, projectId).
		Order("created_at ASC").Find(&events).Error
	return events, err
}

func (r *NotificationRepository) MarkEventAsProcessed(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.NotificationEvent{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"processed":    true,
			"processed_at": now,
		}).Error
}

// === NotificationConfig ===

func (r *NotificationRepository) GetNotificationConfigByEvent(orgId, projectId uuid.UUID, eventType string) (*models.NotificationConfig, error) {
	var config models.NotificationConfig
	err := r.db.Where("organization_id = ? AND project_id = ? AND event_type = ?", orgId, projectId, eventType).
		First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *NotificationRepository) CreateOrUpdateNotificationConfig(config *models.NotificationConfig) error {
	// Verificar se já existe
	var existing models.NotificationConfig
	err := r.db.Where("organization_id = ? AND project_id = ? AND event_type = ?",
		config.OrganizationId, config.ProjectId, config.EventType).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// Criar novo
		config.Id = uuid.New()
		config.CreatedAt = time.Now()
		config.UpdatedAt = time.Now()
		return r.db.Create(config).Error
	} else if err != nil {
		return err
	} else {
		// Atualizar existente
		config.Id = existing.Id
		config.CreatedAt = existing.CreatedAt
		config.UpdatedAt = time.Now()
		return r.db.Save(config).Error
	}
}

// === NotificationTemplate ===

func (r *NotificationRepository) GetNotificationTemplateByChannel(orgId, projectId uuid.UUID, channel string) (*models.NotificationTemplate, error) {
	var template models.NotificationTemplate
	err := r.db.Where("organization_id = ? AND project_id = ? AND channel = ? AND active = true",
		orgId, projectId, channel).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *NotificationRepository) CreateNotificationTemplate(template *models.NotificationTemplate) error {
	template.Id = uuid.New()
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()
	return r.db.Create(template).Error
}

func (r *NotificationRepository) UpdateNotificationTemplate(template *models.NotificationTemplate) error {
	template.UpdatedAt = time.Now()
	return r.db.Save(template).Error
}

func (r *NotificationRepository) GetNotificationTemplatesByProject(orgId, projectId uuid.UUID) ([]models.NotificationTemplate, error) {
	var templates []models.NotificationTemplate
	err := r.db.Where("organization_id = ? AND project_id = ?", orgId, projectId).
		Order("channel ASC, name ASC").Find(&templates).Error
	return templates, err
}