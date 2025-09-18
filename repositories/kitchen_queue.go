package repositories

import (
	"lep/repositories/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KitchenQueueRepository struct {
	db *gorm.DB
}

type IKitchenQueueRepository interface {
	GetActiveOrders(orgId, projectId uuid.UUID) ([]models.Order, error)
	GetOrdersByStatus(orgId, projectId uuid.UUID, status string) ([]models.Order, error)
	GetKitchenQueue(orgId, projectId uuid.UUID) ([]models.Order, error)
	UpdateOrderStatus(orderId uuid.UUID, status string) error
	GetOrdersInPreparation(orgId, projectId uuid.UUID) ([]models.Order, error)
}

func NewKitchenQueueRepository(db *gorm.DB) IKitchenQueueRepository {
	return &KitchenQueueRepository{db: db}
}

// GetActiveOrders retorna todos os pedidos ativos (não cancelados, não entregues)
func (r *KitchenQueueRepository) GetActiveOrders(orgId, projectId uuid.UUID) ([]models.Order, error) {
	var orders []models.Order

	err := r.db.Where(
		"organization_id = ? AND project_id = ? AND status IN (?, ?, ?) AND deleted_at IS NULL",
		orgId, projectId, "pending", "preparing", "ready",
	).Order("created_at ASC").Find(&orders).Error

	return orders, err
}

// GetOrdersByStatus retorna pedidos por status específico
func (r *KitchenQueueRepository) GetOrdersByStatus(orgId, projectId uuid.UUID, status string) ([]models.Order, error) {
	var orders []models.Order

	err := r.db.Where(
		"organization_id = ? AND project_id = ? AND status = ? AND deleted_at IS NULL",
		orgId, projectId, status,
	).Order("created_at ASC").Find(&orders).Error

	return orders, err
}

// GetKitchenQueue retorna a fila da cozinha ordenada por prioridade
func (r *KitchenQueueRepository) GetKitchenQueue(orgId, projectId uuid.UUID) ([]models.Order, error) {
	var orders []models.Order

	// Ordena por: pedidos preparando primeiro, depois por tempo de criação
	err := r.db.Where(
		"organization_id = ? AND project_id = ? AND status IN (?, ?) AND deleted_at IS NULL",
		orgId, projectId, "pending", "preparing",
	).Order("CASE WHEN status = 'preparing' THEN 0 ELSE 1 END, created_at ASC").Find(&orders).Error

	return orders, err
}

// UpdateOrderStatus atualiza apenas o status do pedido
func (r *KitchenQueueRepository) UpdateOrderStatus(orderId uuid.UUID, status string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", orderId).Update("status", status).Error
}

// GetOrdersInPreparation retorna pedidos que estão sendo preparados
func (r *KitchenQueueRepository) GetOrdersInPreparation(orgId, projectId uuid.UUID) ([]models.Order, error) {
	var orders []models.Order

	err := r.db.Where(
		"organization_id = ? AND project_id = ? AND status = ? AND deleted_at IS NULL",
		orgId, projectId, "preparing",
	).Find(&orders).Error

	return orders, err
}