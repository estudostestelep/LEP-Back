package server

import (
	"lep/handler"
	"lep/repositories/models"
	"lep/resource/validation"
	"lep/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ResourceTag struct {
	handler *handler.Handlers
}

type IServerTag interface {
	ServiceGetTag(c *gin.Context)
	ServiceCreateTag(c *gin.Context)
	ServiceUpdateTag(c *gin.Context)
	ServiceDeleteTag(c *gin.Context)
	ServiceListTags(c *gin.Context)
	ServiceListActiveTags(c *gin.Context)
	ServiceGetTagsByEntityType(c *gin.Context)
}

func (r *ResourceTag) ServiceGetTag(c *gin.Context) {
	id, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "tag")
	if !ok {
		return
	}

	resp, err := r.handler.HandlerTag.GetTag(id.String())
	if err != nil {
		utils.SendInternalServerError(c, "Error getting tag", err)
		return
	}

	if resp == nil {
		utils.SendNotFoundError(c, "Tag")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceTag) ServiceCreateTag(c *gin.Context) {
	var newTag models.Tag
	err := c.BindJSON(&newTag)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	newTag.OrganizationId, err = uuid.Parse(organizationId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing organization ID", err)
		return
	}
	newTag.ProjectId, err = uuid.Parse(projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing project ID", err)
		return
	}

	// Gerar ID se não fornecido
	if newTag.Id == uuid.Nil {
		newTag.Id = uuid.New()
	}

	// Validações estruturadas
	if err := validation.CreateTagValidation(&newTag); err != nil {
		utils.SendValidationError(c, "Validation failed", err)
		return
	}

	err = r.handler.HandlerTag.CreateTag(&newTag)
	if err != nil {
		// Verificar se é erro de duplicata
		if strings.Contains(err.Error(), "already_exists") {
			utils.SendConflictError(c, "Tag with this name and type already exists", err)
			return
		}
		utils.SendInternalServerError(c, "Error creating tag", err)
		return
	}

	utils.SendCreatedSuccess(c, "Tag created successfully", newTag)
}

func (r *ResourceTag) ServiceUpdateTag(c *gin.Context) {
	id, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "tag")
	if !ok {
		return
	}

	var updatedTag models.Tag
	if err := c.BindJSON(&updatedTag); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	var err error
	updatedTag.OrganizationId, err = uuid.Parse(organizationId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing organization ID", err)
		return
	}
	updatedTag.ProjectId, err = uuid.Parse(projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error parsing project ID", err)
		return
	}
	updatedTag.Id = id

	// Validações estruturadas
	if err := validation.UpdateTagValidation(&updatedTag); err != nil {
		utils.SendValidationError(c, "Validation failed", err)
		return
	}

	err = r.handler.HandlerTag.UpdateTag(&updatedTag)
	if err != nil {
		// Verificar se é erro de duplicata
		if strings.Contains(err.Error(), "already_exists") {
			utils.SendConflictError(c, "Tag with this name and type already exists", err)
			return
		}
		utils.SendInternalServerError(c, "Error updating tag", err)
		return
	}

	utils.SendOKSuccess(c, "Tag updated successfully", updatedTag)
}

func (r *ResourceTag) ServiceDeleteTag(c *gin.Context) {
	id, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "tag")
	if !ok {
		return
	}

	if err := r.handler.HandlerTag.DeleteTag(id.String()); err != nil {
		utils.SendInternalServerError(c, "Error deleting tag", err)
		return
	}

	utils.SendOKSuccess(c, "Tag deleted successfully", nil)
}

func (r *ResourceTag) ServiceListTags(c *gin.Context) {
	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	resp, err := r.handler.HandlerTag.ListTags(organizationId, projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing tags", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceTag) ServiceListActiveTags(c *gin.Context) {
	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	resp, err := r.handler.HandlerTag.ListActiveTags(organizationId, projectId)
	if err != nil {
		utils.SendInternalServerError(c, "Error listing active tags", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceTag) ServiceGetTagsByEntityType(c *gin.Context) {
	entityType := c.Param("entityType")

	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	resp, err := r.handler.HandlerTag.GetTagsByEntityType(organizationId, projectId, entityType)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting tags by entity type", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func NewSourceServerTag(handler *handler.Handlers) IServerTag {
	return &ResourceTag{handler: handler}
}
