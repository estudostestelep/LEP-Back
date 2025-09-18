package handler

import (
	"lep/repositories"
	"lep/repositories/models"
	"lep/utils"
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
	UpdateOrderStatus(orderId string, status string) error
	GetKitchenQueue(orgId, projectId string) ([]models.Order, error)
	CalculateEstimatedTime(order *models.Order) error
}

// Handler para propósito Order
type OrderHandler struct {
	repo        repositories.IOrderRepository
	productRepo repositories.IProductRepository
	kitchenRepo repositories.IKitchenQueueRepository
}

func NewOrderHandler(repo repositories.IOrderRepository, productRepo repositories.IProductRepository, kitchenRepo repositories.IKitchenQueueRepository) IOrderHandler {
	return &OrderHandler{repo, productRepo, kitchenRepo}
}

func (h *OrderHandler) CreateOrder(order *models.Order) error {
	order.Id = uuid.New()
	order.Status = "pending"
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	// Calcular tempo estimado automaticamente
	err := h.CalculateEstimatedTime(order)
	if err != nil {
		return err
	}

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

// UpdateOrderStatus atualiza o status do pedido com timestamps apropriados
func (h *OrderHandler) UpdateOrderStatus(orderId string, status string) error {
	_, err := uuid.Parse(orderId)
	if err != nil {
		return err
	}

	// Busca o pedido atual
	order, err := h.repo.GetOrderById(orderId)
	if err != nil {
		return err
	}

	// Atualiza status e timestamps usando utils
	utils.UpdateOrderStatus(order, status)

	// Salva no banco
	return h.repo.UpdateOrder(order)
}

// GetKitchenQueue retorna a fila da cozinha
func (h *OrderHandler) GetKitchenQueue(orgId, projectId string) ([]models.Order, error) {
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		return nil, err
	}

	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return nil, err
	}

	return h.kitchenRepo.GetKitchenQueue(orgUUID, projUUID)
}

// CalculateEstimatedTime calcula e define o tempo estimado do pedido
func (h *OrderHandler) CalculateEstimatedTime(order *models.Order) error {
	if len(order.Items) == 0 {
		return nil // Pedido sem itens
	}

	// Buscar informações dos produtos
	var productIds []uuid.UUID
	for _, item := range order.Items {
		productIds = append(productIds, item.ProductId)
	}

	products, err := h.productRepo.GetProductsByIds(productIds)
	if err != nil {
		return err
	}

	// Calcular tempo de preparo
	prepTime := utils.CalculateOrderPrepTime(order.Items, products)
	order.EstimatedPrepTime = prepTime

	// Buscar fila da cozinha para calcular tempo de espera
	activeOrders, err := h.kitchenRepo.GetOrdersInPreparation(order.OrganizationId, order.ProjectId)
	if err != nil {
		return err
	}

	queueTime := utils.GetKitchenQueueTime(activeOrders)

	// Calcular tempo estimado de entrega
	deliveryTime := utils.CalculateEstimatedDeliveryTime(prepTime, queueTime)
	order.EstimatedDeliveryTime = &deliveryTime

	return nil
}
