package handler

import (
	"lep/repositories"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
)

type resourceStaffAttendance struct {
	repo *repositories.DBconn
}

type IHandlerStaffAttendance interface {
	// Attendance
	GetById(id string) (*models.StaffAttendance, error)
	GetByIdWithDetails(id string) (*models.StaffAttendance, error)
	ListByDateRange(orgId, projectId, startDate, endDate string) ([]models.StaffAttendance, error)
	ListByClient(orgId, projectId, clientId string) ([]models.StaffAttendance, error)
	CreateWithDetails(req *models.CreateStaffAttendanceRequest, orgId, projectId string) (*models.StaffAttendance, error)
	Update(attendance *models.StaffAttendance) error
	Delete(id string) error

	// Consumption Products
	GetConsumptionProductById(id string) (*models.StaffConsumptionProduct, error)
	ListConsumptionProducts(orgId, projectId string) ([]models.StaffConsumptionProduct, error)
	CreateConsumptionProduct(product *models.StaffConsumptionProduct) error
	UpdateConsumptionProduct(product *models.StaffConsumptionProduct) error
	DeleteConsumptionProduct(id string) error
}

func NewStaffAttendanceHandler(repo *repositories.DBconn) IHandlerStaffAttendance {
	return &resourceStaffAttendance{repo: repo}
}

// ==================== Attendance ====================

func (r *resourceStaffAttendance) GetById(id string) (*models.StaffAttendance, error) {
	attendanceId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return r.repo.StaffAttendance.GetById(attendanceId)
}

func (r *resourceStaffAttendance) GetByIdWithDetails(id string) (*models.StaffAttendance, error) {
	attendanceId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return r.repo.StaffAttendance.GetByIdWithDetails(attendanceId)
}

func (r *resourceStaffAttendance) ListByDateRange(orgId, projectId, startDate, endDate string) ([]models.StaffAttendance, error) {
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

	return r.repo.StaffAttendance.ListByDateRange(orgUUID, projUUID, start, end)
}

func (r *resourceStaffAttendance) ListByClient(orgId, projectId, clientId string) ([]models.StaffAttendance, error) {
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

	return r.repo.StaffAttendance.ListByClient(orgUUID, projUUID, clientUUID)
}

func (r *resourceStaffAttendance) CreateWithDetails(req *models.CreateStaffAttendanceRequest, orgId, projectId string) (*models.StaffAttendance, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	attendanceDate, err := time.Parse("2006-01-02", req.AttendanceDate)
	if err != nil {
		return nil, err
	}

	// Create attendance
	attendance := &models.StaffAttendance{
		Id:             uuid.New(),
		OrganizationId: orgUUID,
		ProjectId:      projUUID,
		ClientId:       req.ClientId,
		AttendanceDate: attendanceDate,
		Shift:          req.Shift,
		CheckedIn:      true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := r.repo.StaffAttendance.Create(attendance); err != nil {
		return nil, err
	}

	// Create transports
	if req.TransportIn != nil && req.TransportIn.Mode != "" {
		transport := &models.StaffAttendanceTransport{
			AttendanceId:   attendance.Id,
			Direction:      models.TransportDirectionIn,
			Mode:           req.TransportIn.Mode,
			RideWithClient: req.TransportIn.RideWithClient,
			BusTickets:     req.TransportIn.BusTickets,
			UberCost:       req.TransportIn.UberCost,
			CreatedAt:      time.Now(),
		}
		r.repo.StaffAttendance.CreateTransport(transport)
	}

	if req.TransportOut != nil && req.TransportOut.Mode != "" {
		transport := &models.StaffAttendanceTransport{
			AttendanceId:   attendance.Id,
			Direction:      models.TransportDirectionOut,
			Mode:           req.TransportOut.Mode,
			RideWithClient: req.TransportOut.RideWithClient,
			BusTickets:     req.TransportOut.BusTickets,
			UberCost:       req.TransportOut.UberCost,
			CreatedAt:      time.Now(),
		}
		r.repo.StaffAttendance.CreateTransport(transport)
	}

	// Create consumptions
	for _, c := range req.Consumptions {
		qty := c.Quantity
		if qty == 0 {
			qty = 1
		}
		consumption := &models.StaffAttendanceConsumption{
			AttendanceId: attendance.Id,
			ProductId:    c.ProductId,
			Quantity:     qty,
			CreatedAt:    time.Now(),
		}
		r.repo.StaffAttendance.CreateConsumption(consumption)
	}

	return r.GetByIdWithDetails(attendance.Id.String())
}

func (r *resourceStaffAttendance) Update(attendance *models.StaffAttendance) error {
	return r.repo.StaffAttendance.Update(attendance)
}

func (r *resourceStaffAttendance) Delete(id string) error {
	attendanceId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.repo.StaffAttendance.SoftDelete(attendanceId)
}

// ==================== Consumption Products ====================

func (r *resourceStaffAttendance) GetConsumptionProductById(id string) (*models.StaffConsumptionProduct, error) {
	productId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return r.repo.StaffAttendance.GetConsumptionProductById(productId)
}

func (r *resourceStaffAttendance) ListConsumptionProducts(orgId, projectId string) ([]models.StaffConsumptionProduct, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}

	products, err := r.repo.Products.ListProducts(orgUUID, projUUID)
	if err != nil {
		return nil, err
	}

	result := make([]models.StaffConsumptionProduct, 0, len(products))
	for _, p := range products {
		result = append(result, models.StaffConsumptionProduct{
			Id:     p.Id,
			Name:   p.Name,
			Active: p.Active,
		})
	}
	return result, nil
}

func (r *resourceStaffAttendance) CreateConsumptionProduct(product *models.StaffConsumptionProduct) error {
	return r.repo.StaffAttendance.CreateConsumptionProduct(product)
}

func (r *resourceStaffAttendance) UpdateConsumptionProduct(product *models.StaffConsumptionProduct) error {
	return r.repo.StaffAttendance.UpdateConsumptionProduct(product)
}

func (r *resourceStaffAttendance) DeleteConsumptionProduct(id string) error {
	productId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.repo.StaffAttendance.SoftDeleteConsumptionProduct(productId)
}
