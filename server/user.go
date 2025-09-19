package server

import (
	"fmt"
	"lep/handler"
	"lep/repositories/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ResourceUsers struct {
	handler handler.IHandlerUser
}

type IServerUsers interface {
	ServiceGetUser(c *gin.Context)
	ServiceGetUserByGroup(c *gin.Context)
	ServiceListUsers(c *gin.Context)
	ServiceCreateUser(c *gin.Context)
	ServiceUpdateUser(c *gin.Context)
	ServiceDeleteUser(c *gin.Context)
}

func (r *ResourceUsers) ServiceGetUser(c *gin.Context) {
	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("the header param 'X-Lpe-Organization-Id' cannot be empty. Some required params are empty"),
		})
		return
	}

	projectId := c.GetHeader("X-Lpe-Project-Id")
	if strings.TrimSpace(projectId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("the header param 'X-Lpe-Project-Id' cannot be empty. Some required params are empty"),
		})
		return
	}

	id := c.Param("id")
	resp, err := r.handler.GetUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user"})
		return
	}

	if resp == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceUsers) ServiceGetUserByGroup(c *gin.Context) {
	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("the header param 'X-Lpe-Organization-Id' cannot be empty. Some required params are empty"),
		})
		return
	}

	projectId := c.GetHeader("X-Lpe-Project-Id")
	if strings.TrimSpace(projectId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("the header param 'X-Lpe-Project-Id' cannot be empty. Some required params are empty"),
		})
		return
	}

	id := c.Param("id")
	resp, err := r.handler.GetUserByGroup(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting users by group"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceUsers) ServiceCreateUser(c *gin.Context) {
	var newUser models.User
	err := c.BindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = r.handler.CreateUser(&newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

func (r *ResourceUsers) ServiceUpdateUser(c *gin.Context) {
	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("the header param 'X-Lpe-Organization-Id' cannot be empty. Some required params are empty"),
		})
		return
	}

	projectId := c.GetHeader("X-Lpe-Project-Id")
	if strings.TrimSpace(projectId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("the header param 'X-Lpe-Project-Id' cannot be empty. Some required params are empty"),
		})
		return
	}

	var updatedUser models.User
	err := c.BindJSON(&updatedUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = r.handler.UpdateUser(&updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (r *ResourceUsers) ServiceListUsers(c *gin.Context) {
	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	resp, err := r.handler.ListUsers(organizationId, projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error listing users"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceUsers) ServiceDeleteUser(c *gin.Context) {
	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("the header param 'X-Lpe-Organization-Id' cannot be empty. Some required params are empty"),
		})
		return
	}

	projectId := c.GetHeader("X-Lpe-Project-Id")
	if strings.TrimSpace(projectId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("the header param 'X-Lpe-Project-Id' cannot be empty. Some required params are empty"),
		})
		return
	}

	id := c.Param("id")
	err := r.handler.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func NewSourceServerUsers(handler *handler.Handlers) IServerUsers {
	return &ResourceUsers{handler: handler.HandlerUser}
}
