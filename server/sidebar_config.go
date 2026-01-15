package server

import (
	"lep/handler"
	"lep/repositories/models"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type SidebarConfigServer struct {
	handler handler.ISidebarConfigHandler
}

type ISidebarConfigServer interface {
	GetConfig(c *gin.Context)
	UpdateConfig(c *gin.Context)
	ResetConfig(c *gin.Context)
}

func NewSidebarConfigServer(handler handler.ISidebarConfigHandler) ISidebarConfigServer {
	return &SidebarConfigServer{handler: handler}
}

// isMasterAdminCheck verifica se o usuário é master_admin
// Função local para evitar import cycle com middleware
func isMasterAdminCheck(c *gin.Context) bool {
	userPermissions, exists := c.Get("user_permissions")
	if !exists {
		return false
	}

	var permissions []string
	if strArr, ok := userPermissions.([]string); ok {
		permissions = strArr
	} else {
		val := reflect.ValueOf(userPermissions)
		if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
			for i := 0; i < val.Len(); i++ {
				elem := val.Index(i)
				if elem.Kind() == reflect.String {
					permissions = append(permissions, elem.String())
				}
			}
		}
	}

	for _, p := range permissions {
		if p == "master_admin" {
			return true
		}
	}
	return false
}

// GetConfig busca a configuração da sidebar para a organização atual
func (s *SidebarConfigServer) GetConfig(c *gin.Context) {
	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Organization-Id' cannot be empty",
		})
		return
	}

	config, err := s.handler.GetByOrganization(organizationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching sidebar config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    config,
		"message": "Sidebar config retrieved successfully",
	})
}

// UpdateConfig atualiza a configuração da sidebar (apenas Master Admin)
func (s *SidebarConfigServer) UpdateConfig(c *gin.Context) {
	// Verificar permissão de Master Admin
	if !isMasterAdminCheck(c) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied: Master Admin only",
		})
		return
	}

	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Organization-Id' cannot be empty",
		})
		return
	}

	var request models.SidebarConfigUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validar behaviors
	for _, item := range request.Items {
		if item.Behavior != models.BehaviorShow &&
			item.Behavior != models.BehaviorLock &&
			item.Behavior != models.BehaviorHide {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid behavior: must be 'show', 'lock', or 'hide'",
			})
			return
		}
	}

	config, err := s.handler.UpdateConfig(organizationId, request.Items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating sidebar config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    config,
		"message": "Sidebar config updated successfully",
	})
}

// ResetConfig reseta a configuração da sidebar para os valores padrão (apenas Master Admin)
func (s *SidebarConfigServer) ResetConfig(c *gin.Context) {
	// Verificar permissão de Master Admin
	if !isMasterAdminCheck(c) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied: Master Admin only",
		})
		return
	}

	organizationId := c.GetHeader("X-Lpe-Organization-Id")
	if strings.TrimSpace(organizationId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the header param 'X-Lpe-Organization-Id' cannot be empty",
		})
		return
	}

	config, err := s.handler.ResetToDefaults(organizationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error resetting sidebar config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    config,
		"message": "Sidebar config reset to defaults successfully",
	})
}
