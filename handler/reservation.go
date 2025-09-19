package handler

import (
	"lep/repositories"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
)

type resourceReservation struct {
	repo *repositories.DBconn
}

type IHandlerReservation interface {
	GetReservation(id string) (*models.Reservation, error)
	CreateReservation(reservation *models.Reservation) error
	UpdateReservation(updatedReservation *models.Reservation) error
	DeleteReservation(id string) error
	ListReservations(orgId, projectId string) ([]models.Reservation, error)
}

func (r *resourceReservation) GetReservation(id string) (*models.Reservation, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	resp, err := r.repo.Reservations.GetReservationById(uuid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceReservation) CreateReservation(reservation *models.Reservation) error {
	reservation.Id = uuid.New()
	reservation.Status = "confirmed"
	reservation.CreatedAt = time.Now()
	reservation.UpdatedAt = time.Now()
	err := r.repo.Reservations.CreateReservation(reservation)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceReservation) UpdateReservation(updatedReservation *models.Reservation) error {
	updatedReservation.UpdatedAt = time.Now()
	err := r.repo.Reservations.UpdateReservation(updatedReservation)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceReservation) DeleteReservation(id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	err = r.repo.Reservations.SoftDeleteReservation(uuid)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceReservation) ListReservations(orgId, projectId string) ([]models.Reservation, error) {
	orgUuid, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projectUuid, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	resp, err := r.repo.Reservations.ListReservations(orgUuid, projectUuid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewSourceHandlerReservation(repo *repositories.DBconn) IHandlerReservation {
	return &resourceReservation{repo: repo}
}