package server

import (
	"fmt"
	"lep/repositories/models"
	"lep/utils"
	"net/http"
	"strings"
	"time"

	"lep/handler"
	"lep/resource/validation"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IOrderServer interface {
	CreateOrder(c *gin.Context)
	GetOrderById(c *gin.Context)
	ListOrders(c *gin.Context)
	UpdateOrder(c *gin.Context)
	SoftDeleteOrder(c *gin.Context)
	UpdateOrderStatus(c *gin.Context)
	GetKitchenQueue(c *gin.Context)
	GetOrderProgress(c *gin.Context)
}

type OrderServer struct {
	handler handler.IOrderHandler

}

func NewOrderServer(handler handler.IOrderHandler) IOrderServer {
	return &OrderServer{handler}
}

func (s *OrderServer) CreateOrder(c *gin.Context) {
	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("the header param 'X-Lpe-Organization-Id' cannot be empty. Some required params are empty"),
		})
		return
	}

	projectId := c.GetHeader("X-Lpe-Project-Id")
	if strings.TrimSpace(projectId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("the header param 'X-Lpe-Project-Id' cannot be empty. Some required params are empty"),
		})
		return
	}

	createOrderPOST := new(models.Order)
	err := c.Bind(createOrderPOST)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := validation.CreateOrderValidation(createOrderPOST); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = s.handler.CreateOrder(createOrderPOST)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createOrderPOST)
}


func (s *OrderServer) GetOrderById(c *gin.Context) {
	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("the header param 'X-Lpe-Organization-Id' cannot be empty. Some required params are empty"),
		})
		return
	}

	projectId := c.GetHeader("X-Lpe-Project-Id")
	if strings.TrimSpace(projectId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("the header param 'X-Lpe-Project-Id' cannot be empty. Some required params are empty"),
		})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		return
	}
	order, err := s.handler.GetOrderById(id.String())
	if err != nil || order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}



func (s *OrderServer) ListOrders(c *gin.Context) {
	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("the header param 'X-Lpe-Organization-Id' cannot be empty. Some required params are empty"),
		})
		return
	}

	projectId := c.GetHeader("X-Lpe-Project-Id")
	if strings.TrimSpace(projectId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("the header param 'X-Lpe-Project-Id' cannot be empty. Some required params are empty"),
		})
		return
	}

	orders, err := s.handler.ListOrders(organizationId, projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error listing orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (s *OrderServer) UpdateOrder(c *gin.Context) {
	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("the header param 'X-Lpe-Organization-Id' cannot be empty. Some required params are empty"),
		})
		return
	}

	projectId := c.GetHeader("X-Lpe-Project-Id")
	if strings.TrimSpace(projectId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("the header param 'X-Lpe-Project-Id' cannot be empty. Some required params are empty"),
		})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		return
	}
	order, err := s.handler.GetOrderById(id.String())
	if err != nil || order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	if err := c.ShouldBindJSON(order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	order.UpdatedAt = time.Now()
	if err := s.handler.UpdateOrder(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating order"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (s *OrderServer) SoftDeleteOrder(c *gin.Context) {
	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("the header param 'X-Lpe-Organization-Id' cannot be empty. Some required params are empty"),
		})
		return
	}

	projectId := c.GetHeader("X-Lpe-Project-Id")
	if strings.TrimSpace(projectId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("the header param 'X-Lpe-Project-Id' cannot be empty. Some required params are empty"),
		})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		return
	}
	if err := s.handler.SoftDeleteOrder(id.String()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting order"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order deleted"})
}

// UpdateOrderStatus atualiza o status do pedido
func (s *OrderServer) UpdateOrderStatus(c *gin.Context) {
	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Organization-Id' cannot be empty",
		})
		return
	}

	projectId := c.GetHeader("X-Lpe-Project-Id")
	if strings.TrimSpace(projectId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Project-Id' cannot be empty",
		})
		return
	}

	idStr := c.Param("id")
	_, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		return
	}

	var statusUpdate struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validar status permitidos
	allowedStatuses := []string{"pending", "preparing", "ready", "delivered", "cancelled"}
	validStatus := false
	for _, status := range allowedStatuses {
		if statusUpdate.Status == status {
			validStatus = true
			break
		}
	}

	if !validStatus {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid status. Allowed: pending, preparing, ready, delivered, cancelled",
		})
		return
	}

	if err := s.handler.UpdateOrderStatus(idStr, statusUpdate.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

// GetKitchenQueue retorna a fila da cozinha
func (s *OrderServer) GetKitchenQueue(c *gin.Context) {
	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Organization-Id' cannot be empty",
		})
		return
	}

	projectId := c.GetHeader("X-Lpe-Project-Id")
	if strings.TrimSpace(projectId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Project-Id' cannot be empty",
		})
		return
	}

	orders, err := s.handler.GetKitchenQueue(organizationId, projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching kitchen queue"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrderProgress retorna o progresso e tempo restante do pedido
func (s *OrderServer) GetOrderProgress(c *gin.Context) {
	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Organization-Id' cannot be empty",
		})
		return
	}

	projectId := c.GetHeader("X-Lpe-Project-Id")
	if strings.TrimSpace(projectId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Project-Id' cannot be empty",
		})
		return
	}

	idStr := c.Param("id")
	_, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		return
	}

	order, err := s.handler.GetOrderById(idStr)
	if err != nil || order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	progress := utils.GetOrderProgress(*order)
	remainingTime := utils.GetRemainingTime(*order)

	c.JSON(http.StatusOK, gin.H{
		"order_id":        order.Id,
		"status":          order.Status,
		"progress_percent": progress,
		"remaining_minutes": remainingTime,
		"estimated_delivery": order.EstimatedDeliveryTime,
	})
}
