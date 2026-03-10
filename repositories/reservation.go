package repositories

import (
	"fmt"
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
	GetPendingConfirmationReservation(orgId, projectId, customerId uuid.UUID) (*models.Reservation, error)
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
	if reservation.Id == uuid.Nil {
		return fmt.Errorf("reservation ID cannot be empty")
	}
	return r.db.Model(reservation).Where("id = ?", reservation.Id).Updates(reservation).Error
}

func (r *ReservationRepository) SoftDeleteReservation(id uuid.UUID) error {
	return r.db.Model(&models.Reservation{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

// IsReservationTableAvailable verifica se uma mesa está disponível para uma nova reserva em dt.
// Considera dois cenários de bloqueio:
// 1. Reservas futuras: existe reserva no intervalo [dt, dt+diningDurationMinutes] → outra pessoa chega logo depois
// 2. Reservas em andamento: existe reserva ativa iniciada hoje antes de dt → grupo ainda pode estar na mesa
func (r *ReservationRepository) IsReservationTableAvailable(tableId uuid.UUID, dt time.Time, diningDurationMinutes int) (bool, error) {
	dtStr := dt.Format(time.RFC3339)
	futureEnd := dt.Add(time.Duration(diningDurationMinutes) * time.Minute).Format(time.RFC3339)
	// Início do dia de dt para verificar reservas em andamento do mesmo dia
	dayStart := time.Date(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, dt.Location()).Format(time.RFC3339)

	var count int64
	err := r.db.Model(&models.Reservation{}).
		Where(`table_id = ? AND deleted_at IS NULL AND status IN ? AND (
			(datetime >= ? AND datetime <= ?) OR
			(datetime >= ? AND datetime < ?)
		)`,
			tableId, []string{"confirmed", "pending"},
			dtStr, futureEnd,
			dayStart, dtStr,
		).Count(&count).Error
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

// GetPendingConfirmationReservation busca a reserva mais próxima do cliente aguardando confirmação
// Retorna reservas futuras com status "confirmed" ou "awaiting_confirmation"
func (r *ReservationRepository) GetPendingConfirmationReservation(orgId, projectId, customerId uuid.UUID) (*models.Reservation, error) {
	var reservation models.Reservation
	now := time.Now()

	// Busca a reserva futura mais próxima para este cliente
	err := r.db.Where(
		"organization_id = ? AND project_id = ? AND customer_id = ? AND datetime > ? AND status IN (?, ?) AND deleted_at IS NULL",
		orgId, projectId, customerId, now, "confirmed", "awaiting_confirmation",
	).Order("datetime ASC").First(&reservation).Error

	if err != nil {
		return nil, err
	}
	return &reservation, nil
}
