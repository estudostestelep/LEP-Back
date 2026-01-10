package server

import (
	"lep/handler"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// OnboardingServer handles HTTP requests for onboarding status
type OnboardingServer struct {
	handler handler.IOnboardingHandler
}

// IOnboardingServer interface for onboarding server operations
type IOnboardingServer interface {
	GetOnboardingStatus(c *gin.Context)
}

// NewOnboardingServer creates a new OnboardingServer
func NewOnboardingServer(handler handler.IOnboardingHandler) IOnboardingServer {
	return &OnboardingServer{handler: handler}
}

// GetOnboardingStatus returns the current onboarding status for the project
func (s *OnboardingServer) GetOnboardingStatus(c *gin.Context) {
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

	status, err := s.handler.GetOnboardingStatus(organizationId, projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error fetching onboarding status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": status,
	})
}
