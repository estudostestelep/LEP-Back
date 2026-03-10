package server

import (
	"lep/handler"
	"lep/repositories/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ResourceReservation struct {
	handler *handler.Handlers
}

type IServerReservation interface {
	ServiceGetReservation(c *gin.Context)
	ServiceCreateReservation(c *gin.Context)
	ServiceUpdateReservation(c *gin.Context)
	ServiceDeleteReservation(c *gin.Context)
	ServiceListReservations(c *gin.Context)
}

func (r *ResourceReservation) ServiceGetReservation(c *gin.Context) {
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

	id := c.Param("id")
	resp, err := r.handler.HandlerReservation.GetReservation(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting reservation"})
		return
	}

	if resp == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceReservation) ServiceCreateReservation(c *gin.Context) {
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

	var newReservation models.Reservation
	err := c.BindJSON(&newReservation)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Setar IDs da organização e projeto a partir dos headers
	newReservation.OrganizationId, err = uuid.Parse(organizationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing organization ID"})
		return
	}
	newReservation.ProjectId, err = uuid.Parse(projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing project ID"})
		return
	}

	// Gerar ID se não fornecido
	if newReservation.Id == uuid.Nil {
		newReservation.Id = uuid.New()
	}

	err = r.handler.HandlerReservation.CreateReservation(&newReservation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newReservation)
}

func (r *ResourceReservation) ServiceUpdateReservation(c *gin.Context) {
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

	// Obter ID da URL
	id := c.Param("id")
	reservationId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation ID"})
		return
	}

	// Buscar reserva atual para comparar status
	currentReservation, _ := r.handler.HandlerReservation.GetReservation(id)

	var updatedReservation models.Reservation
	err = c.BindJSON(&updatedReservation)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Setar IDs obrigatórios
	updatedReservation.Id = reservationId
	updatedReservation.OrganizationId, _ = uuid.Parse(organizationId)
	updatedReservation.ProjectId, _ = uuid.Parse(projectId)

	err = r.handler.HandlerReservation.UpdateReservation(&updatedReservation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Disparar notificação se status mudou para confirmed, not_approved ou pending
	if currentReservation != nil && currentReservation.Status != updatedReservation.Status && r.handler.EventService != nil {
		newStatus := updatedReservation.Status
		if newStatus == "confirmed" || newStatus == "not_approved" || newStatus == "pending" {
			orgUUID, _ := uuid.Parse(organizationId)
			projUUID, _ := uuid.Parse(projectId)

			customer, customerErr := r.handler.HandlerCustomer.GetCustomer(updatedReservation.CustomerId.String())
			var table *models.Table
			if updatedReservation.TableId != nil {
				t, tableErr := r.handler.HandlerTables.GetTable(updatedReservation.TableId.String())
				if tableErr == nil {
					table = t
				}
			}

			if customerErr == nil && customer != nil {
				r.handler.EventService.TriggerReservationStatusChanged(orgUUID, projUUID, &updatedReservation, customer, table)
			}
		}
	}

	c.JSON(http.StatusOK, updatedReservation)
}

func (r *ResourceReservation) ServiceDeleteReservation(c *gin.Context) {
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

	id := c.Param("id")
	err := r.handler.HandlerReservation.DeleteReservation(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting reservation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reservation deleted successfully"})
}

func (r *ResourceReservation) ServiceListReservations(c *gin.Context) {
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

	resp, err := r.handler.HandlerReservation.ListReservations(organizationId, projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error listing reservations"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func NewSourceServerReservation(handler *handler.Handlers) IServerReservation {
	return &ResourceReservation{handler: handler}
}