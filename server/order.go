package server

import (
	"fmt"
	"lep/repositories/models"
	"net/http"
	"strings"
	"time"

	"lep/handler"
	validate "lep/resource/validation"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IOrderServer interface {
	CreateOrder(c *gin.Context)
	GetOrderById(c *gin.Context)
	ListOrders(c *gin.Context)
	UpdateOrder(c *gin.Context)
	SoftDeleteOrder(c *gin.Context)
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

	if err := validate.CreateOrderRequestValidation(createOrderPOST); err != nil {
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
