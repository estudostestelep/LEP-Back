package handler

import (
	"lep/repositories"
	"lep/repositories/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TableHandler struct {
	repo *repositories.TableRepository
}

func NewTableHandler(repo *repositories.TableRepository) *TableHandler {
	return &TableHandler{repo}
}

func (h *TableHandler) Create(c *gin.Context) {
	var req models.Table
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}
	req.Id = uuid.New()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	if err := h.repo.Create(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating table"})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *TableHandler) GetById(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	table, err := h.repo.GetById(id)
	if err != nil || table == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
		return
	}
	c.JSON(http.StatusOK, table)
}

func (h *TableHandler) List(c *gin.Context) {
	OrganizationId, _ := uuid.Parse(c.Query("org_id"))
	projectId, _ := uuid.Parse(c.Query("project_id"))
	tables, err := h.repo.List(OrganizationId, projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error listing tables"})
		return
	}
	c.JSON(http.StatusOK, tables)
}

func (h *TableHandler) Update(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	table, err := h.repo.GetById(id)
	if err != nil || table == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
		return
	}
	if err := c.ShouldBindJSON(table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}
	table.UpdatedAt = time.Now()
	if err := h.repo.Update(table); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating table"})
		return
	}
	c.JSON(http.StatusOK, table)
}

func (h *TableHandler) SoftDelete(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	err := h.repo.SoftDelete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting table"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Table deleted"})
}
