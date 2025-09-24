package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IReservationRepository interface {
	CreateReservation(reservation *models.Reservation) error
	GetReservationById(id uuid.UUID) (*models.Reservation, error)
	ListReservations(OrganizationId, projectId uuid.UUID) ([]models.Reservation, error)
	UpdateReservation(reservation *models.Reservation) error
	SoftDeleteReservation(id uuid.UUID) error
	IsReservationTableAvailable(tableId uuid.UUID, dt time.Time, durationMinutes int) (bool, error)
	GetReservationsByProject(orgId, projectId uuid.UUID) ([]models.Reservation, error)
	GetReservationsByTableAndDateRange(tableId uuid.UUID, startDate, endDate time.Time) ([]models.Reservation, error)
	DeleteReservation(id uuid.UUID) error
}

type ReservationRepository struct {
	db *gorm.DB
}

func NewConnReservation(db *gorm.DB) *ReservationRepository {
	return &ReservationRepository{db}
}

func (r *ReservationRepository) CreateReservation(reservation *models.Reservation) error {
	return r.db.Create(reservation).Error
}

func (r *ReservationRepository) GetReservationById(id uuid.UUID) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.db.First(&reservation, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *ReservationRepository) ListReservations(OrganizationId, projectId uuid.UUID) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Where("organization_id = ? AND project_id = ? AND deleted_at IS NULL", OrganizationId, projectId).Find(&reservations).Error
	return reservations, err
}

func (r *ReservationRepository) UpdateReservation(reservation *models.Reservation) error {
	return r.db.Save(reservation).Error
}

func (r *ReservationRepository) SoftDeleteReservation(id uuid.UUID) error {
	return r.db.Model(&models.Reservation{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

// Verifica disponibilidade de mesa no intervalo (+/- 1h)
func (r *ReservationRepository) IsReservationTableAvailable(tableId uuid.UUID, dt time.Time, durationMinutes int) (bool, error) {
	start := dt.Add(-time.Duration(durationMinutes) * time.Minute)
	end := dt.Add(time.Duration(durationMinutes) * time.Minute)
	var count int64
	err := r.db.Model(&models.Reservation{}).
		Where("table_id = ? AND status = ? AND datetime BETWEEN ? AND ? AND deleted_at IS NULL", tableId, "confirmed", start, end).
		Count(&count).Error
	return count == 0, err
}

func (r *ReservationRepository) GetReservationsByProject(orgId, projectId uuid.UUID) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Where("organization_id = ? AND project_id = ? AND deleted_at IS NULL", orgId, projectId).
		Order("datetime ASC").Find(&reservations).Error
	return reservations, err
}

func (r *ReservationRepository) GetReservationsByTableAndDateRange(tableId uuid.UUID, startDate, endDate time.Time) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Where("table_id = ? AND datetime BETWEEN ? AND ? AND deleted_at IS NULL", tableId, startDate, endDate).
		Find(&reservations).Error
	return reservations, err
}

func (r *ReservationRepository) DeleteReservation(id uuid.UUID) error {
	return r.db.Delete(&models.Reservation{}, id).Error
}
