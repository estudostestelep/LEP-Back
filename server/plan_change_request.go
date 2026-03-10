package server

import (
	"lep/handler"
	"lep/repositories/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PlanChangeRequestServer struct {
	handler handler.IPlanChangeRequestHandler
}

type IPlanChangeRequestServer interface {
	// Client endpoints
	CreateRequest(c *gin.Context)
	GetMyRequests(c *gin.Context)
	GetRequestById(c *gin.Context)
	CancelRequest(c *gin.Context)

	// Admin endpoints
	GetAllRequests(c *gin.Context)
	GetPendingRequests(c *gin.Context)
	GetRequestsByOrganization(c *gin.Context)
	ApproveRequest(c *gin.Context)
	RejectRequest(c *gin.Context)
}

func NewPlanChangeRequestServer(handler handler.IPlanChangeRequestHandler) IPlanChangeRequestServer {
	return &PlanChangeRequestServer{handler: handler}
}

// ==================== Client Endpoints ====================

// CreateRequest godoc
// @Summary Cria uma nova solicitação de mudança de plano
// @Tags PlanChangeRequest
// @Accept json
// @Produce json
// @Param request body CreatePlanChangeRequestDTO true "Dados da solicitação"
// @Success 201 {object} models.PlanChangeRequest
// @Router /plan-change-request [post]
func (s *PlanChangeRequestServer) CreateRequest(c *gin.Context) {
	userId := c.GetString("user_id")
	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
		return
	}

	orgId := c.GetString("organization_id")
	if orgId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Organization ID is required"})
		return
	}

	var dto CreatePlanChangeRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data", "error": err.Error()})
		return
	}

	// Converter DTO para modelo
	request := dto.ToModel(orgId)

	if err := s.handler.CreateRequest(request, userId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create request", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Plan change request created successfully",
		"data":    request,
	})
}

// GetMyRequests godoc
// @Summary Lista as solicitações do usuário logado
// @Tags PlanChangeRequest
// @Produce json
// @Success 200 {array} models.PlanChangeRequest
// @Router /plan-change-request/my-requests [get]
func (s *PlanChangeRequestServer) GetMyRequests(c *gin.Context) {
	userId := c.GetString("user_id")
	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
		return
	}

	requests, err := s.handler.GetMyRequests(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch requests", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": requests})
}

// GetRequestById godoc
// @Summary Busca uma solicitação por ID
// @Tags PlanChangeRequest
// @Produce json
// @Param id path string true "ID da solicitação"
// @Success 200 {object} models.PlanChangeRequest
// @Router /plan-change-request/{id} [get]
func (s *PlanChangeRequestServer) GetRequestById(c *gin.Context) {
	id := c.Param("id")
	userId := c.GetString("user_id")

	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
		return
	}

	request, err := s.handler.GetRequestById(id, userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Request not found", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": request})
}

// CancelRequest godoc
// @Summary Cancela uma solicitação
// @Tags PlanChangeRequest
// @Produce json
// @Param id path string true "ID da solicitação"
// @Success 200 {object} map[string]string
// @Router /plan-change-request/{id}/cancel [post]
func (s *PlanChangeRequestServer) CancelRequest(c *gin.Context) {
	id := c.Param("id")
	userId := c.GetString("user_id")

	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
		return
	}

	if err := s.handler.CancelRequest(id, userId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to cancel request", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request cancelled successfully"})
}

// ==================== Admin Endpoints ====================

// GetAllRequests godoc
// @Summary Lista todas as solicitações (Admin)
// @Tags PlanChangeRequest
// @Produce json
// @Param status query string false "Filtrar por status (pending, approved, rejected, cancelled)"
// @Success 200 {array} models.PlanChangeRequest
// @Router /admin/plan-change-request [get]
func (s *PlanChangeRequestServer) GetAllRequests(c *gin.Context) {
	status := c.Query("status")

	requests, err := s.handler.GetAllRequests(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch requests", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": requests})
}

// GetPendingRequests godoc
// @Summary Lista todas as solicitações pendentes (Admin)
// @Tags PlanChangeRequest
// @Produce json
// @Success 200 {array} models.PlanChangeRequest
// @Router /admin/plan-change-request/pending [get]
func (s *PlanChangeRequestServer) GetPendingRequests(c *gin.Context) {
	requests, err := s.handler.GetPendingRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch pending requests", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  requests,
		"count": len(requests),
	})
}

// GetRequestsByOrganization godoc
// @Summary Lista solicitações de uma organização específica (Admin)
// @Tags PlanChangeRequest
// @Produce json
// @Param orgId path string true "ID da organização"
// @Param status query string false "Filtrar por status"
// @Success 200 {array} models.PlanChangeRequest
// @Router /admin/plan-change-request/organization/{orgId} [get]
func (s *PlanChangeRequestServer) GetRequestsByOrganization(c *gin.Context) {
	orgId := c.Param("orgId")
	status := c.Query("status")

	requests, err := s.handler.GetRequestsByOrganization(orgId, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch requests", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": requests})
}

// ApproveRequest godoc
// @Summary Aprova uma solicitação de mudança de plano (Admin)
// @Tags PlanChangeRequest
// @Accept json
// @Produce json
// @Param id path string true "ID da solicitação"
// @Param data body ReviewRequestDTO true "Dados da aprovação"
// @Success 200 {object} models.PlanChangeRequest
// @Router /admin/plan-change-request/{id}/approve [post]
func (s *PlanChangeRequestServer) ApproveRequest(c *gin.Context) {
	id := c.Param("id")
	userId := c.GetString("user_id")

	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
		return
	}

	var dto ReviewRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data", "error": err.Error()})
		return
	}

	request, err := s.handler.ApproveRequest(id, userId, dto.ReviewNotes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to approve request", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Request approved successfully",
		"data":    request,
	})
}

// RejectRequest godoc
// @Summary Rejeita uma solicitação de mudança de plano (Admin)
// @Tags PlanChangeRequest
// @Accept json
// @Produce json
// @Param id path string true "ID da solicitação"
// @Param data body ReviewRequestDTO true "Dados da rejeição"
// @Success 200 {object} models.PlanChangeRequest
// @Router /admin/plan-change-request/{id}/reject [post]
func (s *PlanChangeRequestServer) RejectRequest(c *gin.Context) {
	id := c.Param("id")
	userId := c.GetString("user_id")

	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
		return
	}

	var dto ReviewRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data", "error": err.Error()})
		return
	}

	request, err := s.handler.RejectRequest(id, userId, dto.ReviewNotes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to reject request", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Request rejected successfully",
		"data":    request,
	})
}

// ==================== DTOs ====================

type CreatePlanChangeRequestDTO struct {
	RequestedPackageId   string `json:"requested_package_id" binding:"required"`
	CurrentPackageId     string `json:"current_package_id,omitempty"`
	RequestedPackageName string `json:"requested_package_name" binding:"required"`
	CurrentPackageName   string `json:"current_package_name,omitempty"`
	Reason               string `json:"reason,omitempty"`
	Notes                string `json:"notes,omitempty"`
	RequestedBillingCycle string `json:"requested_billing_cycle,omitempty"`
}

func (dto *CreatePlanChangeRequestDTO) ToModel(orgId string) *models.PlanChangeRequest {
	request := &models.PlanChangeRequest{
		RequestedPackageName: dto.RequestedPackageName,
		CurrentPackageName:   dto.CurrentPackageName,
		Reason:               dto.Reason,
		Notes:                dto.Notes,
		RequestedBillingCycle: dto.RequestedBillingCycle,
	}

	// Parse organization ID
	if orgId != "" {
		if orgUUID, err := parseUUID(orgId); err == nil {
			request.OrganizationId = orgUUID
		}
	}

	// Parse requested package ID
	if dto.RequestedPackageId != "" {
		if pkgUUID, err := parseUUID(dto.RequestedPackageId); err == nil {
			request.RequestedPackageId = pkgUUID
		}
	}

	// Parse current package ID (optional)
	if dto.CurrentPackageId != "" {
		if pkgUUID, err := parseUUID(dto.CurrentPackageId); err == nil {
			request.CurrentPackageId = pkgUUID
		}
	}

	return request
}

type ReviewRequestDTO struct {
	ReviewNotes string `json:"review_notes,omitempty"`
}

// Helper function
func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
