package server

import (
	"lep/handler"
	"lep/repositories/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EnvironmentServer struct {
	handler handler.IEnvironmentHandler
}

type IEnvironmentServer interface {
	GetEnvironmentById(c *gin.Context)
	GetEnvironmentsByProject(c *gin.Context)
	CreateEnvironment(c *gin.Context)
	UpdateEnvironment(c *gin.Context)
	SoftDeleteEnvironment(c *gin.Context)
	GetActiveEnvironments(c *gin.Context)
}

func NewEnvironmentServer(handler handler.IEnvironmentHandler) IEnvironmentServer {
	return &EnvironmentServer{handler: handler}
}

// GetEnvironmentById busca ambiente por ID
func (s *EnvironmentServer) GetEnvironmentById(c *gin.Context) {
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
	environment, err := s.handler.GetEnvironmentById(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	// Verificar se ambiente pertence à organização/projeto
	if environment.OrganizationId.String() != organizationId ||
	   environment.ProjectId.String() != projectId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, environment)
}

// GetEnvironmentsByProject busca ambientes por projeto
func (s *EnvironmentServer) GetEnvironmentsByProject(c *gin.Context) {
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

	environments, err := s.handler.GetEnvironmentsByProject(organizationId, projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching environments"})
		return
	}

	c.JSON(http.StatusOK, environments)
}

// CreateEnvironment cria novo ambiente
func (s *EnvironmentServer) CreateEnvironment(c *gin.Context) {
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

	var environment models.Environment
	if err := c.ShouldBindJSON(&environment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Definir organização e projeto dos headers
	orgUUID, err := uuid.Parse(organizationId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	environment.OrganizationId = orgUUID
	environment.ProjectId = projUUID

	// Validações básicas
	if strings.TrimSpace(environment.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Environment name is required"})
		return
	}
	if environment.Capacity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Capacity must be greater than 0"})
		return
	}

	err = s.handler.CreateEnvironment(&environment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating environment"})
		return
	}

	c.JSON(http.StatusCreated, environment)
}

// UpdateEnvironment atualiza ambiente
func (s *EnvironmentServer) UpdateEnvironment(c *gin.Context) {
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
	envId, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}

	// Verificar se ambiente existe e pertence à organização/projeto
	existingEnv, err := s.handler.GetEnvironmentById(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	if existingEnv.OrganizationId.String() != organizationId ||
	   existingEnv.ProjectId.String() != projectId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var updateData models.Environment
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validações básicas
	if strings.TrimSpace(updateData.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Environment name is required"})
		return
	}
	if updateData.Capacity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Capacity must be greater than 0"})
		return
	}

	// Manter dados imutáveis
	updateData.Id = envId
	updateData.OrganizationId = existingEnv.OrganizationId
	updateData.ProjectId = existingEnv.ProjectId
	updateData.CreatedAt = existingEnv.CreatedAt

	err = s.handler.UpdateEnvironment(&updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating environment"})
		return
	}

	c.JSON(http.StatusOK, updateData)
}

// SoftDeleteEnvironment remove ambiente logicamente
func (s *EnvironmentServer) SoftDeleteEnvironment(c *gin.Context) {
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

	// Verificar se ambiente existe e pertence à organização/projeto
	environment, err := s.handler.GetEnvironmentById(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	if environment.OrganizationId.String() != organizationId ||
	   environment.ProjectId.String() != projectId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	err = s.handler.SoftDeleteEnvironment(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting environment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Environment deleted successfully"})
}

// GetActiveEnvironments busca ambientes ativos
func (s *EnvironmentServer) GetActiveEnvironments(c *gin.Context) {
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

	environments, err := s.handler.GetActiveEnvironments(organizationId, projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching active environments"})
		return
	}

	c.JSON(http.StatusOK, environments)
}