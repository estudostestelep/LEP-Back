package handler

import (
	"lep/repositories"
	"lep/repositories/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserWaitlistHandler struct {
	repo *repositories.WaitlistRepository
}

func NewUserWaitlistHandler(repo *repositories.WaitlistRepository) *UserWaitlistHandler {
	return &UserWaitlistHandler{repo}
}

func (h *UserWaitlistHandler) CreateUserWaitlist(c *gin.Context) {
	var req models.Waitlist
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}
	req.Id = uuid.New()
	req.Status = "waiting"
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	if err := h.repo.CreateWaitlist(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating waitlist"})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *UserWaitlistHandler) GetUserWaitlistById(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	waitlist, err := h.repo.GetWaitlistById(id)
	if err != nil || waitlist == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Waitlist not found"})
		return
	}
	c.JSON(http.StatusOK, waitlist)
}

func (h *UserWaitlistHandler) ListUserWaitlists(c *gin.Context) {
	OrganizationId, _ := uuid.Parse(c.Query("org_id"))
	projectId, _ := uuid.Parse(c.Query("project_id"))
	waitlists, err := h.repo.ListWaitlists(OrganizationId, projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error listing waitlists"})
		return
	}
	c.JSON(http.StatusOK, waitlists)
}

func (h *UserWaitlistHandler) UpdateUserWaitlist(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	waitlist, err := h.repo.GetWaitlistById(id)
	if err != nil || waitlist == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Waitlist not found"})
		return
	}
	if err := c.ShouldBindJSON(waitlist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}
	waitlist.UpdatedAt = time.Now()
	if err := h.repo.UpdateWaitlist(waitlist); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating waitlist"})
		return
	}
	c.JSON(http.StatusOK, waitlist)
}

func (h *UserWaitlistHandler) SoftDeleteUserWaitlist(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	err := h.repo.SoftDeleteWaitlist(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting waitlist"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Waitlist deleted"})
}
