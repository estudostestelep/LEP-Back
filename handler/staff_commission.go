package handler

import (
	"lep/repositories"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
)

type resourceStaffCommission struct {
	repo *repositories.DBconn
}

type IHandlerStaffCommission interface {
	GetById(id string) (*models.StaffDailyCommission, error)
	ListByDateRange(orgId, projectId, startDate, endDate string) ([]models.StaffDailyCommission, error)
	Create(req *models.CreateCommissionRequest, orgId, projectId string) (*models.StaffDailyCommission, error)
	Update(commission *models.StaffDailyCommission) error
	Delete(id string) error
	GetSummary(orgId, projectId, startDate, endDate string) (*models.CommissionSummary, error)
}

func NewStaffCommissionHandler(repo *repositories.DBconn) IHandlerStaffCommission {
	return &resourceStaffCommission{repo: repo}
}

func (r *resourceStaffCommission) GetById(id string) (*models.StaffDailyCommission, error) {
	commissionId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return r.repo.StaffCommission.GetById(commissionId)
}

func (r *resourceStaffCommission) ListByDateRange(orgId, projectId, startDate, endDate string) ([]models.StaffDailyCommission, error) {
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

	return r.repo.StaffCommission.ListByDateRange(orgUUID, projUUID, start, end)
}

func (r *resourceStaffCommission) Create(req *models.CreateCommissionRequest, orgId, projectId string) (*models.StaffDailyCommission, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	commissionDate, err := time.Parse("2006-01-02", req.CommissionDate)
	if err != nil {
		return nil, err
	}

	commission := &models.StaffDailyCommission{
		OrganizationId:  orgUUID,
		ProjectId:       projUUID,
		CommissionDate:  commissionDate,
		Shift:           req.Shift,
		CommissionValue: req.CommissionValue,
		Revenue:         req.Revenue,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err = r.repo.StaffCommission.Upsert(commission)
	if err != nil {
		return nil, err
	}

	return commission, nil
}

func (r *resourceStaffCommission) Update(commission *models.StaffDailyCommission) error {
	return r.repo.StaffCommission.Update(commission)
}

func (r *resourceStaffCommission) Delete(id string) error {
	commissionId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.repo.StaffCommission.SoftDelete(commissionId)
}

func (r *resourceStaffCommission) GetSummary(orgId, projectId, startDate, endDate string) (*models.CommissionSummary, error) {
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

	return r.repo.StaffCommission.GetSummary(orgUUID, projUUID, start, end)
}
