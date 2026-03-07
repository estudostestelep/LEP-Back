package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type resourceStaffAttendance struct {
	db *gorm.DB
}

type IStaffAttendanceRepository interface {
	// Attendance CRUD
	GetById(id uuid.UUID) (*models.StaffAttendance, error)
	GetByIdWithDetails(id uuid.UUID) (*models.StaffAttendance, error)
	ListByDateRange(orgId, projectId uuid.UUID, startDate, endDate time.Time) ([]models.StaffAttendance, error)
	ListByClient(orgId, projectId, clientId uuid.UUID) ([]models.StaffAttendance, error)
	Create(attendance *models.StaffAttendance) error
	Update(attendance *models.StaffAttendance) error
	Delete(id uuid.UUID) error
	SoftDelete(id uuid.UUID) error

	// Transport management
	CreateTransport(transport *models.StaffAttendanceTransport) error
	DeleteTransportsByAttendance(attendanceId uuid.UUID) error

	// Consumption management
	CreateConsumption(consumption *models.StaffAttendanceConsumption) error
	DeleteConsumptionsByAttendance(attendanceId uuid.UUID) error

	// Consumption Products CRUD
	GetConsumptionProductById(id uuid.UUID) (*models.StaffConsumptionProduct, error)
	ListConsumptionProducts(orgId, projectId uuid.UUID) ([]models.StaffConsumptionProduct, error)
	CreateConsumptionProduct(product *models.StaffConsumptionProduct) error
	UpdateConsumptionProduct(product *models.StaffConsumptionProduct) error
	SoftDeleteConsumptionProduct(id uuid.UUID) error
}

func NewStaffAttendanceRepository(db *gorm.DB) IStaffAttendanceRepository {
	return &resourceStaffAttendance{db: db}
}

// ==================== Attendance ====================

func (r *resourceStaffAttendance) GetById(id uuid.UUID) (*models.StaffAttendance, error) {
	var attendance models.StaffAttendance
	err := r.db.First(&attendance, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (r *resourceStaffAttendance) GetByIdWithDetails(id uuid.UUID) (*models.StaffAttendance, error) {
	var attendance models.StaffAttendance
	err := r.db.
		Preload("Client").
		Preload("Transports").
		Preload("Transports.RideWith").
		Preload("Consumptions").
		Preload("Consumptions.Product").
		First(&attendance, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (r *resourceStaffAttendance) ListByDateRange(orgId, projectId uuid.UUID, startDate, endDate time.Time) ([]models.StaffAttendance, error) {
	var attendances []models.StaffAttendance
	err := r.db.
		Preload("Client").
		Preload("Transports").
		Preload("Consumptions").
		Preload("Consumptions.Product").
		Where("organization_id = ? AND project_id = ? AND attendance_date >= ? AND attendance_date <= ? AND deleted_at IS NULL",
			orgId, projectId, startDate, endDate).
		Order("attendance_date DESC, shift ASC").
		Find(&attendances).Error
	return attendances, err
}

func (r *resourceStaffAttendance) ListByClient(orgId, projectId, clientId uuid.UUID) ([]models.StaffAttendance, error) {
	var attendances []models.StaffAttendance
	err := r.db.
		Preload("Transports").
		Preload("Consumptions").
		Preload("Consumptions.Product").
		Where("organization_id = ? AND project_id = ? AND client_id = ? AND deleted_at IS NULL", orgId, projectId, clientId).
		Order("attendance_date DESC").
		Find(&attendances).Error
	return attendances, err
}

func (r *resourceStaffAttendance) Create(attendance *models.StaffAttendance) error {
	if attendance.Id == uuid.Nil {
		attendance.Id = uuid.New()
	}
	return r.db.Create(attendance).Error
}

func (r *resourceStaffAttendance) Update(attendance *models.StaffAttendance) error {
	return r.db.Save(attendance).Error
}

func (r *resourceStaffAttendance) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.StaffAttendance{}, "id = ?", id).Error
}

func (r *resourceStaffAttendance) SoftDelete(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.StaffAttendance{}).Where("id = ?", id).Update("deleted_at", now).Error
}

// ==================== Transport ====================

func (r *resourceStaffAttendance) CreateTransport(transport *models.StaffAttendanceTransport) error {
	if transport.Id == uuid.Nil {
		transport.Id = uuid.New()
	}
	return r.db.Create(transport).Error
}

func (r *resourceStaffAttendance) DeleteTransportsByAttendance(attendanceId uuid.UUID) error {
	return r.db.Delete(&models.StaffAttendanceTransport{}, "attendance_id = ?", attendanceId).Error
}

// ==================== Consumption ====================

func (r *resourceStaffAttendance) CreateConsumption(consumption *models.StaffAttendanceConsumption) error {
	if consumption.Id == uuid.Nil {
		consumption.Id = uuid.New()
	}
	return r.db.Create(consumption).Error
}

func (r *resourceStaffAttendance) DeleteConsumptionsByAttendance(attendanceId uuid.UUID) error {
	return r.db.Delete(&models.StaffAttendanceConsumption{}, "attendance_id = ?", attendanceId).Error
}

// ==================== Consumption Products ====================

func (r *resourceStaffAttendance) GetConsumptionProductById(id uuid.UUID) (*models.StaffConsumptionProduct, error) {
	var product models.StaffConsumptionProduct
	err := r.db.First(&product, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *resourceStaffAttendance) ListConsumptionProducts(orgId, projectId uuid.UUID) ([]models.StaffConsumptionProduct, error) {
	var products []models.StaffConsumptionProduct
	err := r.db.
		Where("organization_id = ? AND project_id = ? AND deleted_at IS NULL AND active = true", orgId, projectId).
		Order("name ASC").
		Find(&products).Error
	return products, err
}

func (r *resourceStaffAttendance) CreateConsumptionProduct(product *models.StaffConsumptionProduct) error {
	if product.Id == uuid.Nil {
		product.Id = uuid.New()
	}
	return r.db.Create(product).Error
}

func (r *resourceStaffAttendance) UpdateConsumptionProduct(product *models.StaffConsumptionProduct) error {
	return r.db.Save(product).Error
}

func (r *resourceStaffAttendance) SoftDeleteConsumptionProduct(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.StaffConsumptionProduct{}).Where("id = ?", id).Update("deleted_at", now).Error
}
