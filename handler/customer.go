package handler

import (
	"lep/repositories"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
)

type resourceCustomer struct {
	repo *repositories.DBconn
}

type IHandlerCustomer interface {
	GetCustomer(id string) (*models.Customer, error)
	CreateCustomer(customer *models.Customer) error
	UpdateCustomer(updatedCustomer *models.Customer) error
	DeleteCustomer(id string) error
	ListCustomers(orgId, projectId string) ([]models.Customer, error)
}

func (r *resourceCustomer) GetCustomer(id string) (*models.Customer, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	resp, err := r.repo.Customers.GetCustomerById(uuid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resourceCustomer) CreateCustomer(customer *models.Customer) error {
	customer.Id = uuid.New()
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()
	err := r.repo.Customers.CreateCustomer(customer)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceCustomer) UpdateCustomer(updatedCustomer *models.Customer) error {
	updatedCustomer.UpdatedAt = time.Now()
	err := r.repo.Customers.UpdateCustomer(updatedCustomer)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceCustomer) DeleteCustomer(id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	err = r.repo.Customers.SoftDeleteCustomer(uuid)
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceCustomer) ListCustomers(orgId, projectId string) ([]models.Customer, error) {
	orgUuid, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}
	projectUuid, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}
	resp, err := r.repo.Customers.ListCustomers(orgUuid, projectUuid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewSourceHandlerCustomer(repo *repositories.DBconn) IHandlerCustomer {
	return &resourceCustomer{repo: repo}
}