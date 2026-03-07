package handler

import (
	"lep/repositories"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
)

type resourceStaffAvailability struct {
	repo *repositories.DBconn
}

type IHandlerStaffAvailability interface {
	GetById(id string) (*models.StaffAvailability, error)
	GetByClientAndWeek(orgId, projectId, clientId, weekStart string) (*models.StaffAvailability, error)
	ListByWeek(orgId, projectId, weekStart string) ([]models.StaffAvailability, error)
	ListByClient(orgId, projectId, clientId string) ([]models.StaffAvailability, error)
	Upsert(availability *models.StaffAvailability) error
	Delete(id string) error
	GetWeekSummary(orgId, projectId, weekStart string) (*models.WeekAvailabilitySummary, error)
}

func NewStaffAvailabilityHandler(repo *repositories.DBconn) IHandlerStaffAvailability {
	return &resourceStaffAvailability{repo: repo}
}

func (r *resourceStaffAvailability) GetById(id string) (*models.StaffAvailability, error) {
	availabilityId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return r.repo.StaffAvailability.GetById(availabilityId)
}

func (r *resourceStaffAvailability) GetByClientAndWeek(orgId, projectId, clientId, weekStart string) (*models.StaffAvailability, error) {
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
	weekStartDate, err := time.Parse("2006-01-02", weekStart)
	if err != nil {
		return nil, err
	}

	return r.repo.StaffAvailability.GetByClientAndWeek(orgUUID, projUUID, clientUUID, weekStartDate)
}

func (r *resourceStaffAvailability) ListByWeek(orgId, projectId, weekStart string) ([]models.StaffAvailability, error) {
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

	return r.repo.StaffAvailability.ListByWeek(orgUUID, projUUID, weekStartDate)
}

func (r *resourceStaffAvailability) ListByClient(orgId, projectId, clientId string) ([]models.StaffAvailability, error) {
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

	return r.repo.StaffAvailability.ListByClient(orgUUID, projUUID, clientUUID)
}

func (r *resourceStaffAvailability) Upsert(availability *models.StaffAvailability) error {
	return r.repo.StaffAvailability.Upsert(availability)
}

func (r *resourceStaffAvailability) Delete(id string) error {
	availabilityId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.repo.StaffAvailability.Delete(availabilityId)
}

func (r *resourceStaffAvailability) GetWeekSummary(orgId, projectId, weekStart string) (*models.WeekAvailabilitySummary, error) {
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

	return r.repo.StaffAvailability.GetWeekSummary(orgUUID, projUUID, weekStartDate)
}
