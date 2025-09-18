package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IReportMetricRepository interface {
	CreateMetric(metric *models.ReportMetric) error
	GetMetricById(id uuid.UUID) (*models.ReportMetric, error)
	GetMetricsByProject(orgId, projectId uuid.UUID) ([]models.ReportMetric, error)
	GetMetricsByType(orgId, projectId uuid.UUID, metricType string) ([]models.ReportMetric, error)
	GetMetricsByDateRange(orgId, projectId uuid.UUID, start, end time.Time) ([]models.ReportMetric, error)
	GetMetricsByTypeAndDateRange(orgId, projectId uuid.UUID, metricType string, start, end time.Time) ([]models.ReportMetric, error)
	UpdateMetric(metric *models.ReportMetric) error
	DeleteMetric(id uuid.UUID) error
}

type ReportMetricRepository struct {
	db *gorm.DB
}

func NewReportMetricRepository(db *gorm.DB) IReportMetricRepository {
	return &ReportMetricRepository{db: db}
}

func (r *ReportMetricRepository) CreateMetric(metric *models.ReportMetric) error {
	return r.db.Create(metric).Error
}

func (r *ReportMetricRepository) GetMetricById(id uuid.UUID) (*models.ReportMetric, error) {
	var metric models.ReportMetric
	err := r.db.First(&metric, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &metric, nil
}

func (r *ReportMetricRepository) GetMetricsByProject(orgId, projectId uuid.UUID) ([]models.ReportMetric, error) {
	var metrics []models.ReportMetric
	err := r.db.Where("organization_id = ? AND project_id = ?", orgId, projectId).
		Order("metric_date DESC").Find(&metrics).Error
	return metrics, err
}

func (r *ReportMetricRepository) GetMetricsByType(orgId, projectId uuid.UUID, metricType string) ([]models.ReportMetric, error) {
	var metrics []models.ReportMetric
	err := r.db.Where("organization_id = ? AND project_id = ? AND metric_type = ?", orgId, projectId, metricType).
		Order("metric_date DESC").Find(&metrics).Error
	return metrics, err
}

func (r *ReportMetricRepository) GetMetricsByDateRange(orgId, projectId uuid.UUID, start, end time.Time) ([]models.ReportMetric, error) {
	var metrics []models.ReportMetric
	err := r.db.Where("organization_id = ? AND project_id = ? AND metric_date BETWEEN ? AND ?", orgId, projectId, start, end).
		Order("metric_date DESC").Find(&metrics).Error
	return metrics, err
}

func (r *ReportMetricRepository) GetMetricsByTypeAndDateRange(orgId, projectId uuid.UUID, metricType string, start, end time.Time) ([]models.ReportMetric, error) {
	var metrics []models.ReportMetric
	err := r.db.Where("organization_id = ? AND project_id = ? AND metric_type = ? AND metric_date BETWEEN ? AND ?",
		orgId, projectId, metricType, start, end).
		Order("metric_date DESC").Find(&metrics).Error
	return metrics, err
}

func (r *ReportMetricRepository) UpdateMetric(metric *models.ReportMetric) error {
	return r.db.Save(metric).Error
}

func (r *ReportMetricRepository) DeleteMetric(id uuid.UUID) error {
	return r.db.Delete(&models.ReportMetric{}, id).Error
}