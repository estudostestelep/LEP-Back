package handler

import (
	"lep/repositories"
	"lep/repositories/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CustomerHandler struct {
	repo *repositories.CustomerRepository
}

func NewCustomerHandler(repo *repositories.CustomerRepository) *CustomerHandler {
	return &CustomerHandler{repo}
}

func (h *CustomerHandler) Create(c *gin.Context) {
	var req models.Customer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}
	req.Id = uuid.New()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	if err := h.repo.CreateCustomer(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating customer"})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *CustomerHandler) GetById(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	customer, err := h.repo.GetCustomer(id)
	if err != nil || customer == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) List(c *gin.Context) {
	OrganizationId, _ := uuid.Parse(c.Query("org_id"))
	projectId, _ := uuid.Parse(c.Query("project_id"))
	customers, err := h.repo.GetCustomerList(OrganizationId, projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error listing customers"})
		return
	}
	c.JSON(http.StatusOK, customers)
}

func (h *CustomerHandler) Update(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	customer, err := h.repo.GetCustomer(id)
	if err != nil || customer == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	if err := c.ShouldBindJSON(customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}
	customer.UpdatedAt = time.Now()
	if err := h.repo.UpdateCustomer(customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating customer"})
		return
	}
	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) SoftDelete(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	err := h.repo.SoftDelete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting customer"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted"})
}
