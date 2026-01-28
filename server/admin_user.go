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
	"github.com/lib/pq"
)

type ResourceAdminUsers struct {
	handler handler.IHandlerAdminUser
}

type IServerAdminUsers interface {
	ServiceGetAdmin(c *gin.Context)
	ServiceListAdmins(c *gin.Context)
	ServiceCreateAdmin(c *gin.Context)
	ServiceUpdateAdmin(c *gin.Context)
	ServiceDeleteAdmin(c *gin.Context)
}

// CreateAdminRequest DTO para criar admin
type CreateAdminRequest struct {
	Name        string   `json:"name" binding:"required"`
	Email       string   `json:"email" binding:"required,email"`
	Password    string   `json:"password" binding:"required,min=6"`
	Permissions []string `json:"permissions"`
	Active      *bool    `json:"active"`
}

func (r *ResourceAdminUsers) ServiceGetAdmin(c *gin.Context) {
	id, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "admin")
	if !ok {
		return
	}

	admin, err := r.handler.GetAdminById(id.String())
	if err != nil {
		utils.SendInternalServerError(c, "Error getting admin", err)
		return
	}

	if admin == nil {
		utils.SendNotFoundError(c, "Admin")
		return
	}

	// Remover senha da resposta
	admin.Password = ""

	c.JSON(http.StatusOK, gin.H{"data": admin})
}

func (r *ResourceAdminUsers) ServiceListAdmins(c *gin.Context) {
	admins, err := r.handler.ListAdmins()
	if err != nil {
		utils.SendInternalServerError(c, "Error listing admins", err)
		return
	}

	// Remover senhas das respostas
	for i := range admins {
		admins[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{"data": admins})
}

func (r *ResourceAdminUsers) ServiceCreateAdmin(c *gin.Context) {
	var request CreateAdminRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Converter para modelo
	admin := &models.Admin{
		Id:          uuid.New(),
		Name:        request.Name,
		Email:       request.Email,
		Password:    request.Password,
		Permissions: pq.StringArray(request.Permissions),
		Active:      true,
	}

	if request.Active != nil {
		admin.Active = *request.Active
	}

	if err := r.handler.CreateAdmin(admin); err != nil {
		if strings.Contains(err.Error(), "já cadastrado") {
			utils.SendConflictError(c, "Admin with this email already exists", nil)
			return
		}
		utils.SendInternalServerError(c, "Error creating admin", err)
		return
	}

	// Remover senha da resposta
	admin.Password = ""

	utils.SendCreatedSuccess(c, "Admin created successfully", admin)
}

func (r *ResourceAdminUsers) ServiceUpdateAdmin(c *gin.Context) {
	id, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "admin")
	if !ok {
		return
	}

	var request CreateAdminRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Verificar se admin existe
	existing, err := r.handler.GetAdminById(id.String())
	if err != nil || existing == nil {
		utils.SendNotFoundError(c, "Admin")
		return
	}

	// Atualizar campos
	existing.Name = request.Name
	existing.Email = request.Email
	if request.Password != "" {
		existing.Password = request.Password
	}
	if len(request.Permissions) > 0 {
		existing.Permissions = pq.StringArray(request.Permissions)
	}
	if request.Active != nil {
		existing.Active = *request.Active
	}

	if err := r.handler.UpdateAdmin(existing); err != nil {
		utils.SendInternalServerError(c, "Error updating admin", err)
		return
	}

	// Remover senha da resposta
	existing.Password = ""

	utils.SendOKSuccess(c, "Admin updated successfully", existing)
}

func (r *ResourceAdminUsers) ServiceDeleteAdmin(c *gin.Context) {
	id, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "admin")
	if !ok {
		return
	}

	if err := r.handler.DeleteAdmin(id.String()); err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			utils.SendNotFoundError(c, "Admin")
			return
		}
		utils.SendInternalServerError(c, "Error deleting admin", err)
		return
	}

	utils.SendOKSuccess(c, "Admin deleted successfully", nil)
}

func NewSourceServerAdminUsers(handler *handler.Handlers) IServerAdminUsers {
	return &ResourceAdminUsers{
		handler: handler.HandlerAdminUser,
	}
}
