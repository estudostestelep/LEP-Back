package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IBlockedPeriodRepository interface {
	CreateBlockedPeriod(period *models.BlockedPeriod) error
	GetBlockedPeriodById(id uuid.UUID) (*models.BlockedPeriod, error)
	GetBlockedPeriodsByProject(orgId, projectId uuid.UUID) ([]models.BlockedPeriod, error)
	UpdateBlockedPeriod(period *models.BlockedPeriod) error
	SoftDeleteBlockedPeriod(id uuid.UUID) error
	CheckPeriodBlocked(orgId, projectId uuid.UUID, datetime time.Time) (bool, error)
	GetActiveBlockedPeriodsInRange(orgId, projectId uuid.UUID, start, end time.Time) ([]models.BlockedPeriod, error)
}

type BlockedPeriodRepository struct {
	db *gorm.DB
}

func NewBlockedPeriodRepository(db *gorm.DB) IBlockedPeriodRepository {
	return &BlockedPeriodRepository{db: db}
}

func (r *BlockedPeriodRepository) CreateBlockedPeriod(period *models.BlockedPeriod) error {
	return r.db.Create(period).Error
}

func (r *BlockedPeriodRepository) GetBlockedPeriodById(id uuid.UUID) (*models.BlockedPeriod, error) {
	var period models.BlockedPeriod
	err := r.db.First(&period, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &period, nil
}

func (r *BlockedPeriodRepository) GetBlockedPeriodsByProject(orgId, projectId uuid.UUID) ([]models.BlockedPeriod, error) {
	var periods []models.BlockedPeriod
	err := r.db.Where("organization_id = ? AND project_id = ? AND deleted_at IS NULL", orgId, projectId).
		Order("start_datetime ASC").Find(&periods).Error
	return periods, err
}

func (r *BlockedPeriodRepository) UpdateBlockedPeriod(period *models.BlockedPeriod) error {
	period.UpdatedAt = time.Now()
	return r.db.Save(period).Error
}

func (r *BlockedPeriodRepository) SoftDeleteBlockedPeriod(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.BlockedPeriod{}).Where("id = ?", id).Update("deleted_at", now).Error
}

// CheckPeriodBlocked - Verifica se um horário específico está bloqueado
func (r *BlockedPeriodRepository) CheckPeriodBlocked(orgId, projectId uuid.UUID, datetime time.Time) (bool, error) {
	var count int64
	err := r.db.Model(&models.BlockedPeriod{}).
		Where("organization_id = ? AND project_id = ? AND active = true AND deleted_at IS NULL", orgId, projectId).
		Where("start_datetime <= ? AND end_datetime >= ?", datetime, datetime).
		Count(&count).Error

	return count > 0, err
}

// GetActiveBlockedPeriodsInRange - Busca períodos bloqueados em um intervalo
func (r *BlockedPeriodRepository) GetActiveBlockedPeriodsInRange(orgId, projectId uuid.UUID, start, end time.Time) ([]models.BlockedPeriod, error) {
	var periods []models.BlockedPeriod
	err := r.db.Where("organization_id = ? AND project_id = ? AND active = true AND deleted_at IS NULL", orgId, projectId).
		Where("(start_datetime <= ? AND end_datetime >= ?) OR (start_datetime >= ? AND start_datetime <= ?)",
			end, start, start, end).
		Find(&periods).Error
	return periods, err
}