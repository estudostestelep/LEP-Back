package server

import (
	"lep/handler"
	"lep/repositories/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ResourceWaitlist struct {
	handler *handler.Handlers
}

type IServerWaitlist interface {
	ServiceGetWaitlist(c *gin.Context)
	ServiceCreateWaitlist(c *gin.Context)
	ServiceUpdateWaitlist(c *gin.Context)
	ServiceDeleteWaitlist(c *gin.Context)
	ServiceListWaitlists(c *gin.Context)
}

func (r *ResourceWaitlist) ServiceGetWaitlist(c *gin.Context) {
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
	resp, err := r.handler.HandlerWaitlist.GetWaitlist(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting waitlist"})
		return
	}

	if resp == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Waitlist not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceWaitlist) ServiceCreateWaitlist(c *gin.Context) {
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

	var newWaitlist models.Waitlist
	err := c.BindJSON(&newWaitlist)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = r.handler.HandlerWaitlist.CreateWaitlist(&newWaitlist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newWaitlist)
}

func (r *ResourceWaitlist) ServiceUpdateWaitlist(c *gin.Context) {
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

	var updatedWaitlist models.Waitlist
	err := c.BindJSON(&updatedWaitlist)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = r.handler.HandlerWaitlist.UpdateWaitlist(&updatedWaitlist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedWaitlist)
}

func (r *ResourceWaitlist) ServiceDeleteWaitlist(c *gin.Context) {
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
	err := r.handler.HandlerWaitlist.DeleteWaitlist(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting waitlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Waitlist deleted successfully"})
}

func (r *ResourceWaitlist) ServiceListWaitlists(c *gin.Context) {
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

	resp, err := r.handler.HandlerWaitlist.ListWaitlists(organizationId, projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error listing waitlists"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func NewSourceServerWaitlist(handler *handler.Handlers) IServerWaitlist {
	return &ResourceWaitlist{handler: handler}
}