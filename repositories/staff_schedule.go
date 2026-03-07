package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type resourceStaffSchedule struct {
	db *gorm.DB
}

type IStaffScheduleRepository interface {
	GetById(id uuid.UUID) (*models.StaffSchedule, error)
	ListByDateRange(orgId, projectId uuid.UUID, startDate, endDate time.Time) ([]models.StaffSchedule, error)
	ListByWeek(orgId, projectId uuid.UUID, weekStart time.Time) ([]models.StaffSchedule, error)
	ListByClient(orgId, projectId, clientId uuid.UUID) ([]models.StaffSchedule, error)
	Create(schedule *models.StaffSchedule) error
	CreateBatch(schedules []models.StaffSchedule) error
	Update(schedule *models.StaffSchedule) error
	Delete(id uuid.UUID) error
	SoftDelete(id uuid.UUID) error
	GetByDateAndShift(orgId, projectId uuid.UUID, date time.Time, shift string) ([]models.StaffSchedule, error)
	MarkEmailSent(ids []uuid.UUID) error
	GetWeekSummary(orgId, projectId uuid.UUID, weekStart time.Time) (*models.WeekScheduleSummary, error)
}

func NewStaffScheduleRepository(db *gorm.DB) IStaffScheduleRepository {
	return &resourceStaffSchedule{db: db}
}

func (r *resourceStaffSchedule) GetById(id uuid.UUID) (*models.StaffSchedule, error) {
	var schedule models.StaffSchedule
	err := r.db.Preload("Client").First(&schedule, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *resourceStaffSchedule) ListByDateRange(orgId, projectId uuid.UUID, startDate, endDate time.Time) ([]models.StaffSchedule, error) {
	var schedules []models.StaffSchedule
	err := r.db.
		Preload("Client").
		Where("organization_id = ? AND project_id = ? AND schedule_date >= ? AND schedule_date <= ? AND deleted_at IS NULL",
			orgId, projectId, startDate, endDate).
		Order("schedule_date ASC, shift ASC, slot_number ASC").
		Find(&schedules).Error
	return schedules, err
}

func (r *resourceStaffSchedule) ListByWeek(orgId, projectId uuid.UUID, weekStart time.Time) ([]models.StaffSchedule, error) {
	weekEnd := weekStart.AddDate(0, 0, 6)
	return r.ListByDateRange(orgId, projectId, weekStart, weekEnd)
}

func (r *resourceStaffSchedule) ListByClient(orgId, projectId, clientId uuid.UUID) ([]models.StaffSchedule, error) {
	var schedules []models.StaffSchedule
	err := r.db.
		Where("organization_id = ? AND project_id = ? AND client_id = ? AND deleted_at IS NULL", orgId, projectId, clientId).
		Order("schedule_date DESC").
		Find(&schedules).Error
	return schedules, err
}

func (r *resourceStaffSchedule) Create(schedule *models.StaffSchedule) error {
	if schedule.Id == uuid.Nil {
		schedule.Id = uuid.New()
	}
	return r.db.Create(schedule).Error
}

func (r *resourceStaffSchedule) CreateBatch(schedules []models.StaffSchedule) error {
	for i := range schedules {
		if schedules[i].Id == uuid.Nil {
			schedules[i].Id = uuid.New()
		}
	}
	return r.db.Create(&schedules).Error
}

func (r *resourceStaffSchedule) Update(schedule *models.StaffSchedule) error {
	return r.db.Save(schedule).Error
}

func (r *resourceStaffSchedule) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.StaffSchedule{}, "id = ?", id).Error
}

func (r *resourceStaffSchedule) SoftDelete(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.StaffSchedule{}).Where("id = ?", id).Update("deleted_at", now).Error
}

func (r *resourceStaffSchedule) GetByDateAndShift(orgId, projectId uuid.UUID, date time.Time, shift string) ([]models.StaffSchedule, error) {
	var schedules []models.StaffSchedule
	err := r.db.
		Preload("Client").
		Where("organization_id = ? AND project_id = ? AND schedule_date = ? AND shift = ? AND deleted_at IS NULL",
			orgId, projectId, date, shift).
		Order("slot_number ASC").
		Find(&schedules).Error
	return schedules, err
}

func (r *resourceStaffSchedule) MarkEmailSent(ids []uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.StaffSchedule{}).Where("id IN ?", ids).Update("email_sent_at", now).Error
}

func (r *resourceStaffSchedule) GetWeekSummary(orgId, projectId uuid.UUID, weekStart time.Time) (*models.WeekScheduleSummary, error) {
	schedules, err := r.ListByWeek(orgId, projectId, weekStart)
	if err != nil {
		return nil, err
	}

	summary := &models.WeekScheduleSummary{
		WeekStart: weekStart.Format("2006-01-02"),
		ByDay:     make(map[string][]models.ScheduledClient),
	}

	for _, sched := range schedules {
		dateKey := sched.ScheduleDate.Format("2006-01-02")
		clientName := ""
		if sched.Client != nil {
			clientName = sched.Client.Name
		}

		summary.ByDay[dateKey] = append(summary.ByDay[dateKey], models.ScheduledClient{
			ScheduleId: sched.Id,
			ClientId:   sched.ClientId,
			ClientName: clientName,
			Shift:      sched.Shift,
			Status:     sched.Status,
			SlotNumber: sched.SlotNumber,
		})
	}

	return summary, nil
}
