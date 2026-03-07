package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type resourceStaffCommission struct {
	db *gorm.DB
}

type IStaffCommissionRepository interface {
	// Commission CRUD
	GetById(id uuid.UUID) (*models.StaffDailyCommission, error)
	GetByDateAndShift(orgId, projectId uuid.UUID, date time.Time, shift string) (*models.StaffDailyCommission, error)
	ListByDateRange(orgId, projectId uuid.UUID, startDate, endDate time.Time) ([]models.StaffDailyCommission, error)
	Create(commission *models.StaffDailyCommission) error
	Update(commission *models.StaffDailyCommission) error
	Upsert(commission *models.StaffDailyCommission) error
	SoftDelete(id uuid.UUID) error

	// Summary and reports
	GetSummary(orgId, projectId uuid.UUID, startDate, endDate time.Time) (*models.CommissionSummary, error)
}

func NewStaffCommissionRepository(db *gorm.DB) IStaffCommissionRepository {
	return &resourceStaffCommission{db: db}
}

func (r *resourceStaffCommission) GetById(id uuid.UUID) (*models.StaffDailyCommission, error) {
	var commission models.StaffDailyCommission
	err := r.db.First(&commission, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &commission, nil
}

func (r *resourceStaffCommission) GetByDateAndShift(orgId, projectId uuid.UUID, date time.Time, shift string) (*models.StaffDailyCommission, error) {
	var commission models.StaffDailyCommission
	err := r.db.Where(
		"organization_id = ? AND project_id = ? AND commission_date = ? AND shift = ? AND deleted_at IS NULL",
		orgId, projectId, date, shift,
	).First(&commission).Error
	if err != nil {
		return nil, err
	}
	return &commission, nil
}

func (r *resourceStaffCommission) ListByDateRange(orgId, projectId uuid.UUID, startDate, endDate time.Time) ([]models.StaffDailyCommission, error) {
	var commissions []models.StaffDailyCommission
	err := r.db.
		Where("organization_id = ? AND project_id = ? AND commission_date >= ? AND commission_date <= ? AND deleted_at IS NULL",
			orgId, projectId, startDate, endDate).
		Order("commission_date DESC, shift ASC").
		Find(&commissions).Error
	return commissions, err
}

func (r *resourceStaffCommission) Create(commission *models.StaffDailyCommission) error {
	if commission.Id == uuid.Nil {
		commission.Id = uuid.New()
	}
	return r.db.Create(commission).Error
}

func (r *resourceStaffCommission) Update(commission *models.StaffDailyCommission) error {
	return r.db.Save(commission).Error
}

func (r *resourceStaffCommission) Upsert(commission *models.StaffDailyCommission) error {
	existing, err := r.GetByDateAndShift(
		commission.OrganizationId,
		commission.ProjectId,
		commission.CommissionDate,
		commission.Shift,
	)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if existing != nil {
		existing.CommissionValue = commission.CommissionValue
		existing.Revenue = commission.Revenue
		existing.UpdatedAt = time.Now()
		return r.Update(existing)
	}

	return r.Create(commission)
}

func (r *resourceStaffCommission) SoftDelete(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.StaffDailyCommission{}).Where("id = ?", id).Update("deleted_at", now).Error
}

func (r *resourceStaffCommission) GetSummary(orgId, projectId uuid.UUID, startDate, endDate time.Time) (*models.CommissionSummary, error) {
	var result struct {
		TotalRevenue    float64
		TotalCommission float64
		DaysCount       int
	}

	err := r.db.Model(&models.StaffDailyCommission{}).
		Select("COALESCE(SUM(revenue), 0) as total_revenue, COALESCE(SUM(commission_value), 0) as total_commission, COUNT(DISTINCT commission_date) as days_count").
		Where("organization_id = ? AND project_id = ? AND commission_date >= ? AND commission_date <= ? AND deleted_at IS NULL",
			orgId, projectId, startDate, endDate).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	avgRevenue := 0.0
	if result.DaysCount > 0 {
		avgRevenue = result.TotalRevenue / float64(result.DaysCount)
	}

	return &models.CommissionSummary{
		StartDate:       startDate.Format("2006-01-02"),
		EndDate:         endDate.Format("2006-01-02"),
		TotalRevenue:    result.TotalRevenue,
		TotalCommission: result.TotalCommission,
		DaysCount:       result.DaysCount,
		AverageRevenue:  avgRevenue,
	}, nil
}
