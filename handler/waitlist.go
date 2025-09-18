package handler

import (
	"lep/repositories"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
)

type resourceWaitlist struct {
	repo *repositories.DBconn
}

type IHandlerWaitlist interface {
	GetWaitlist(id string) (*models.Waitlist, error)
	CreateWaitlist(waitlist *models.Waitlist) error
	UpdateWaitlist(updatedWaitlist *models.Waitlist) error
	DeleteWaitlist(id string) error
	ListWaitlists(orgId, projectId string) ([]models.Waitlist, error)
}

func (r *resourceWaitlist) GetWaitlist(id string) (*models.Waitlist, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	resp, err := r.repo.Waitlists.GetWaitlistById(uuid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceWaitlist) CreateWaitlist(waitlist *models.Waitlist) error {
	waitlist.Id = uuid.New()
	waitlist.Status = "waiting"
	waitlist.CreatedAt = time.Now()
	waitlist.UpdatedAt = time.Now()
	err := r.repo.Waitlists.CreateWaitlist(waitlist)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceWaitlist) UpdateWaitlist(updatedWaitlist *models.Waitlist) error {
	updatedWaitlist.UpdatedAt = time.Now()
	err := r.repo.Waitlists.UpdateWaitlist(updatedWaitlist)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceWaitlist) DeleteWaitlist(id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	err = r.repo.Waitlists.SoftDeleteWaitlist(uuid)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceWaitlist) ListWaitlists(orgId, projectId string) ([]models.Waitlist, error) {
	orgUuid, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projectUuid, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	resp, err := r.repo.Waitlists.ListWaitlists(orgUuid, projectUuid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewSourceHandlerWaitlist(repo *repositories.DBconn) IHandlerWaitlist {
	return &resourceWaitlist{repo: repo}
}