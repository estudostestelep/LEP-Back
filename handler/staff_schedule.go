package handler

import (
	"lep/repositories"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
)

type resourceStaffSchedule struct {
	repo *repositories.DBconn
}

type IHandlerStaffSchedule interface {
	GetById(id string) (*models.StaffSchedule, error)
	ListByWeek(orgId, projectId, weekStart string) ([]models.StaffSchedule, error)
	ListByDateRange(orgId, projectId, startDate, endDate string) ([]models.StaffSchedule, error)
	ListByClient(orgId, projectId, clientId string) ([]models.StaffSchedule, error)
	Create(schedule *models.StaffSchedule) error
	CreateBatch(schedules []models.StaffSchedule) error
	Update(schedule *models.StaffSchedule) error
	Delete(id string) error
	GetByDateAndShift(orgId, projectId, date, shift string) ([]models.StaffSchedule, error)
	MarkEmailSent(ids []string) error
	GetWeekSummary(orgId, projectId, weekStart string) (*models.WeekScheduleSummary, error)
}

func NewStaffScheduleHandler(repo *repositories.DBconn) IHandlerStaffSchedule {
	return &resourceStaffSchedule{repo: repo}
}

func (r *resourceStaffSchedule) GetById(id string) (*models.StaffSchedule, error) {
	scheduleId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return r.repo.StaffSchedule.GetById(scheduleId)
}

func (r *resourceStaffSchedule) ListByWeek(orgId, projectId, weekStart string) ([]models.StaffSchedule, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	weekStartDate, err := time.Parse("2006-01-02", weekStart)
	if err != nil {
		return nil, err
	}

	return r.repo.StaffSchedule.ListByWeek(orgUUID, projUUID, weekStartDate)
}

func (r *resourceStaffSchedule) ListByDateRange(orgId, projectId, startDate, endDate string) ([]models.StaffSchedule, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}

	return r.repo.StaffSchedule.ListByDateRange(orgUUID, projUUID, start, end)
}

func (r *resourceStaffSchedule) ListByClient(orgId, projectId, clientId string) ([]models.StaffSchedule, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	clientUUID, err := uuid.Parse(clientId)
	if err != nil {
		return nil, err
	}

	return r.repo.StaffSchedule.ListByClient(orgUUID, projUUID, clientUUID)
}

func (r *resourceStaffSchedule) Create(schedule *models.StaffSchedule) error {
	return r.repo.StaffSchedule.Create(schedule)
}

func (r *resourceStaffSchedule) CreateBatch(schedules []models.StaffSchedule) error {
	return r.repo.StaffSchedule.CreateBatch(schedules)
}

func (r *resourceStaffSchedule) Update(schedule *models.StaffSchedule) error {
	return r.repo.StaffSchedule.Update(schedule)
}

func (r *resourceStaffSchedule) Delete(id string) error {
	scheduleId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.repo.StaffSchedule.SoftDelete(scheduleId)
}

func (r *resourceStaffSchedule) GetByDateAndShift(orgId, projectId, date, shift string) ([]models.StaffSchedule, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	return r.repo.StaffSchedule.GetByDateAndShift(orgUUID, projUUID, dateTime, shift)
}

func (r *resourceStaffSchedule) MarkEmailSent(ids []string) error {
	var uuids []uuid.UUID
	for _, id := range ids {
		u, err := uuid.Parse(id)
		if err != nil {
			return err
		}
		uuids = append(uuids, u)
	}
	return r.repo.StaffSchedule.MarkEmailSent(uuids)
}

func (r *resourceStaffSchedule) GetWeekSummary(orgId, projectId, weekStart string) (*models.WeekScheduleSummary, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	weekStartDate, err := time.Parse("2006-01-02", weekStart)
	if err != nil {
		return nil, err
	}

	return r.repo.StaffSchedule.GetWeekSummary(orgUUID, projUUID, weekStartDate)
}
