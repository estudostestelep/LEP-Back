package server

import (
	"lep/handler"

	"github.com/gin-gonic/gin"
)

type ResourceTenant struct {
	handler *handler.Handlers
}

type IServerTenant interface {
	ServiceResolveTenant(c *gin.Context)
}

// ServiceResolveTenant resolve o tenant baseado no slug
func (r *ResourceTenant) ServiceResolveTenant(c *gin.Context) {
	slug := c.Query("slug")

	if slug == "" {
		c.JSON(400, gin.H{"error": "Slug é obrigatório"})
		return
	}

	// Se o slug for "admin", retornar tipo admin
	if slug == "admin" {
		c.JSON(200, handler.TenantResolveResponse{
			Type:     "admin",
			LoginUrl: "/admin/login",
		})
		return
	}

	// Buscar organização pelo slug
	org, err := r.handler.HandlerOrganization.GetOrganizationBySlug(slug)
	if err != nil || org == nil {
		c.JSON(404, gin.H{"error": "Organização não encontrada"})
		return
	}

	c.JSON(200, handler.TenantResolveResponse{
		Type:             "client",
		OrganizationId:   org.Id.String(),
		OrganizationName: org.Name,
		OrganizationSlug: org.Slug,
		LoginUrl:         "/client/login",
	})
}

func NewSourceServerTenant(handler *handler.Handlers) IServerTenant {
	return &ResourceTenant{handler: handler}
}
