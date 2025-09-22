package handler

import (
	"lep/repositories"
	"lep/repositories/models"

	"github.com/google/uuid"
)

type resourceOrganization struct {
	repo *repositories.DBconn
}

type IHandlerOrganization interface {
	GetOrganizationById(id string) (*models.Organization, error)
	GetOrganizationByEmail(email string) (*models.Organization, error)
	ListOrganizations() ([]models.Organization, error)
	ListActiveOrganizations() ([]models.Organization, error)
	CreateOrganization(organization *models.Organization) error
	UpdateOrganization(organization *models.Organization) error
	SoftDeleteOrganization(id string) error
	HardDeleteOrganization(id string) error
}

func (r *resourceOrganization) GetOrganizationById(id string) (*models.Organization, error) {
	// Validar UUID
	organizationId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	resp, err := r.repo.Organizations.GetOrganizationById(organizationId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceOrganization) GetOrganizationByEmail(email string) (*models.Organization, error) {
	resp, err := r.repo.Organizations.GetOrganizationByEmail(email)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceOrganization) ListOrganizations() ([]models.Organization, error) {
	resp, err := r.repo.Organizations.ListOrganizations()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceOrganization) ListActiveOrganizations() ([]models.Organization, error) {
	resp, err := r.repo.Organizations.ListActiveOrganizations()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceOrganization) CreateOrganization(organization *models.Organization) error {
	organization.Id = uuid.New()
	err := r.repo.Organizations.CreateOrganization(organization)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceOrganization) UpdateOrganization(organization *models.Organization) error {
	err := r.repo.Organizations.UpdateOrganization(organization)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceOrganization) SoftDeleteOrganization(id string) error {
	// Validar UUID
	organizationId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	err = r.repo.Organizations.SoftDeleteOrganization(organizationId)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceOrganization) HardDeleteOrganization(id string) error {
	// Validar UUID
	organizationId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	err = r.repo.Organizations.HardDeleteOrganization(organizationId)
	if err != nil {
		return err
	}
	return nil
}

func NewSourceHandlerOrganization(repo *repositories.DBconn) IHandlerOrganization {
	return &resourceOrganization{repo: repo}
}
