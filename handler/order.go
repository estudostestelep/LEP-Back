package handler

import (
	"lep/repositories"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
)

// Interface para OrderHandler
type IOrderHandler interface {
	CreateOrder(order *models.Order) error
	GetOrderById(id string) (*models.Order, error)
	ListOrders(orgId, projectId string) ([]models.Order, error)
	UpdateOrder(order *models.Order) error
	SoftDeleteOrder(id string) error
}

// Handler para propósito Order
type OrderHandler struct {
	repo repositories.IOrderRepository
}

func NewOrderHandler(repo repositories.IOrderRepository) IOrderHandler {
	return &OrderHandler{repo}
}

func (h *OrderHandler) CreateOrder(order *models.Order) error {
	order.Id = uuid.New()
	order.Status = "pending"
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	//gerar log de criação com defer

	return h.repo.CreateOrder(order)
}

func (h *OrderHandler) GetOrderById(id string) (*models.Order, error) {
	return h.repo.GetOrderById(id)
}

func (h *OrderHandler) ListOrders(orgId, projectId string) ([]models.Order, error) {
	return h.repo.ListOrders(orgId, projectId)
}

func (h *OrderHandler) UpdateOrder(order *models.Order) error {
	order.UpdatedAt = time.Now()
	return h.repo.UpdateOrder(order)
}

func (h *OrderHandler) SoftDeleteOrder(id string) error {
	return h.repo.SoftDeleteOrder(id)
}
