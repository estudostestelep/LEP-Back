package handler

import (
	"lep/repositories"
	"lep/repositories/models"
	"lep/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReservationHandler struct {
	repo      *repositories.ReservationRepository
	auditRepo *repositories.AuditLogRepository
}

func NewReservationHandler(repo *repositories.ReservationRepository, auditRepo *repositories.AuditLogRepository) *ReservationHandler {
	return &ReservationHandler{repo, auditRepo}
}

// Criar reserva
func (h *ReservationHandler) Create(c *gin.Context) {
	var req models.Reservation
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Regra de negócio: validar disponibilidade da mesa
	available, err := h.repo.IsTableAvailable(req.TableId, req.Datetime, 60)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking table availability"})
		return
	}
	if !available {
		c.JSON(http.StatusConflict, gin.H{"error": "Table not available"})
		return
	}

	// Gerar Id e timestamps
	req.Id = uuid.New()
	req.Status = "confirmed"
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	// Salvar reserva
	if err := h.repo.Create(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating reservation"})
		return
	}

	// Geração de log/auditoria
	userId := utils.GetUserIdFromContext(c) // Exemplo de utilitário para extrair userId do contexto/session
	log := &models.AuditLog{
		Id:             uuid.New(),
		OrganizationId: req.OrganizationId,
		ProjectId:      req.ProjectId,
		UserId:         &userId,
		Action:         "create_reservation",
		Entity:         "reservation",
		EntityId:       req.Id,
		Description:    "Reserva criada",
		CreatedAt:      time.Now(),
	}
	_ = h.auditRepo.CreateAuditLog(log) // Não bloqueia em caso de erro de log

	c.JSON(http.StatusCreated, req)
}

// Cancelar reserva
func (h *ReservationHandler) Cancel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation Id"})
		return
	}

	reservation, err := h.repo.GetById(id)
	if err != nil || reservation == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}

	// Regra de negócio: liberar mesa, status cancelado
	reservation.Status = "cancelled"
	reservation.UpdatedAt = time.Now()
	if err := h.repo.Update(reservation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error cancelling reservation"})
		return
	}

	// Log/auditoria
	userId := utils.GetUserIdFromContext(c)
	log := &models.AuditLog{
		Id:             uuid.New(),
		OrganizationId: reservation.OrganizationId,
		ProjectId:      reservation.ProjectId,
		UserId:         &userId,
		Action:         "cancel_reservation",
		Entity:         "reservation",
		EntityId:       reservation.Id,
		Description:    "Reserva cancelada",
		CreatedAt:      time.Now(),
	}
	_ = h.auditRepo.CreateAuditLog(log)

	c.JSON(http.StatusOK, gin.H{"message": "Reservation cancelled"})
}
