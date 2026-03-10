package handler

import (
	"errors"
	"fmt"
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
	GetCustomerByEmail(orgId, projectId uuid.UUID, email string) (*models.Customer, error)
	CreateCustomer(customer *models.Customer) error
	UpdateCustomer(updatedCustomer *models.Customer) error
	DeleteCustomer(id string) error
	ListCustomers(orgId, projectId string) ([]models.Customer, error)
}

func (r *resourceCustomer) GetCustomerByEmail(orgId, projectId uuid.UUID, email string) (*models.Customer, error) {
	return r.repo.Customers.GetCustomerByEmail(orgId, projectId, email)
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
	// Verificar se já existe cliente com o mesmo email no projeto
	if customer.Email != "" {
		exists, err := r.repo.Customers.CheckCustomerEmailExists(customer.OrganizationId, customer.ProjectId, customer.Email, nil)
		if err != nil {
			return fmt.Errorf("erro ao verificar duplicata: %w", err)
		}
		if exists {
			return errors.New("already_exists: customer with this email already exists in this project")
		}
	}

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
	// Verificar se já existe outro cliente com o mesmo email no projeto
	if updatedCustomer.Email != "" {
		exists, err := r.repo.Customers.CheckCustomerEmailExists(updatedCustomer.OrganizationId, updatedCustomer.ProjectId, updatedCustomer.Email, &updatedCustomer.Id)
		if err != nil {
			return fmt.Errorf("erro ao verificar duplicata: %w", err)
		}
		if exists {
			return errors.New("already_exists: customer with this email already exists in this project")
		}
	}

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