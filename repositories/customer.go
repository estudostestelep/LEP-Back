package repositories

import (
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

type ICustomersRepository interface {
	GetCustomer(id uuid.UUID) (*models.Customer, error)
	GetCustomerById(id uuid.UUID) (*models.Customer, error)
	GetCustomerList(OrganizationId, projectId uuid.UUID) ([]models.Customer, error)
	ListCustomers(OrganizationId, projectId uuid.UUID) ([]models.Customer, error)
	CreateCustomer(customer *models.Customer) error
	UpdateCustomer(updatedCustomer *models.Customer) error
	SoftDelete(id uuid.UUID) error
	SoftDeleteCustomer(id uuid.UUID) error
}

func NewConnCustomer(db *gorm.DB) ICustomersRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) CreateCustomer(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

func (r *CustomerRepository) GetCustomer(id uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.First(&customer, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepository) GetCustomerById(id uuid.UUID) (*models.Customer, error) {
	return r.GetCustomer(id)
}

func (r *CustomerRepository) GetCustomerList(OrganizationId, projectId uuid.UUID) ([]models.Customer, error) {
	var customers []models.Customer
	err := r.db.Where("organization_id = ? AND project_id = ? AND deleted_at IS NULL", OrganizationId, projectId).Find(&customers).Error
	return customers, err
}

func (r *CustomerRepository) ListCustomers(OrganizationId, projectId uuid.UUID) ([]models.Customer, error) {
	return r.GetCustomerList(OrganizationId, projectId)
}

func (r *CustomerRepository) UpdateCustomer(customer *models.Customer) error {
	return r.db.Save(customer).Error
}

func (r *CustomerRepository) SoftDelete(id uuid.UUID) error {
	return r.db.Model(&models.Customer{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

func (r *CustomerRepository) SoftDeleteCustomer(id uuid.UUID) error {
	return r.SoftDelete(id)
}
