package server

import (
	"lep/handler"
	"lep/repositories/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProjectServer struct {
	handler handler.IProjectHandler
}

type IProjectServer interface {
	GetProjectById(c *gin.Context)
	GetProjectsByOrganization(c *gin.Context)
	CreateProject(c *gin.Context)
	UpdateProject(c *gin.Context)
	SoftDeleteProject(c *gin.Context)
	GetActiveProjects(c *gin.Context)
}

func NewProjectServer(handler handler.IProjectHandler) IProjectServer {
	return &ProjectServer{handler: handler}
}

// GetProjectById busca projeto por ID
func (s *ProjectServer) GetProjectById(c *gin.Context) {
	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Organization-Id' cannot be empty",
		})
		return
	}

	idStr := c.Param("id")
	if strings.TrimSpace(idStr) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is required"})
		return
	}

	project, err := s.handler.GetProjectById(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Verificar se o projeto pertence à organização
	if project.OrganizationId.String() != organizationId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// GetProjectsByOrganization busca projetos por organização
func (s *ProjectServer) GetProjectsByOrganization(c *gin.Context) {
	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Organization-Id' cannot be empty",
		})
		return
	}

	projects, err := s.handler.GetProjectsByOrganization(organizationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching projects"})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// CreateProject cria novo projeto
func (s *ProjectServer) CreateProject(c *gin.Context) {
	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Organization-Id' cannot be empty",
		})
		return
	}

	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Definir organização do header
	orgUUID, err := uuid.Parse(organizationId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}
	project.OrganizationId = orgUUID

	// Validações básicas
	if strings.TrimSpace(project.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project name is required"})
		return
	}

	err = s.handler.CreateProject(&project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating project"})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// UpdateProject atualiza projeto
func (s *ProjectServer) UpdateProject(c *gin.Context) {
	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Organization-Id' cannot be empty",
		})
		return
	}

	idStr := c.Param("id")
	projectId, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Verificar se projeto existe e pertence à organização
	existingProject, err := s.handler.GetProjectById(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if existingProject.OrganizationId.String() != organizationId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var updateData models.Project
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Manter dados imutáveis
	updateData.Id = projectId
	updateData.OrganizationId = existingProject.OrganizationId
	updateData.CreatedAt = existingProject.CreatedAt

	err = s.handler.UpdateProject(&updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating project"})
		return
	}

	c.JSON(http.StatusOK, updateData)
}

// SoftDeleteProject remove projeto logicamente
func (s *ProjectServer) SoftDeleteProject(c *gin.Context) {
	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Organization-Id' cannot be empty",
		})
		return
	}

	idStr := c.Param("id")

	// Verificar se projeto existe e pertence à organização
	project, err := s.handler.GetProjectById(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if project.OrganizationId.String() != organizationId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	err = s.handler.SoftDeleteProject(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

// GetActiveProjects busca projetos ativos
func (s *ProjectServer) GetActiveProjects(c *gin.Context) {
	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Organization-Id' cannot be empty",
		})
		return
	}

	projects, err := s.handler.GetActiveProjects(organizationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching active projects"})
		return
	}

	c.JSON(http.StatusOK, projects)
}
