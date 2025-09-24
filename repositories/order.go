package repositories

import (
	"lep/repositories/models"
	"time"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

type IOrderRepository interface {
	GetOrderById(id string) (*models.Order, error)
	ListOrders(OrganizationId, projectId string) ([]models.Order, error)
	CreateOrder(order *models.Order) error
	UpdateOrder(order *models.Order) error
	SoftDeleteOrder(id string) error
}

func NewConnOrder(db *gorm.DB) IOrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) GetOrderById(id string) (*models.Order, error) {
	var order models.Order
	err := r.db.First(&order, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) ListOrders(OrganizationId, projectId string) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Where("organization_id = ? AND project_id = ? AND deleted_at IS NULL", OrganizationId, projectId).Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) UpdateOrder(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *OrderRepository) SoftDeleteOrder(id string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}
