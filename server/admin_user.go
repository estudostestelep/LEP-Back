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
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RoleId   string `json:"role_id" binding:"required"` // ID do cargo (role)
	Active   *bool  `json:"active"`
}

// UpdateAdminRequest DTO para atualizar admin
type UpdateAdminRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"omitempty,min=6"` // Opcional no update
	RoleId   string `json:"role_id"`                            // Opcional - se fornecido, muda o cargo
	Active   *bool  `json:"active"`
}

// AdminResponse DTO para retornar admin com informações de role
type AdminResponse struct {
	Id           uuid.UUID      `json:"id"`
	Name         string         `json:"name"`
	Email        string         `json:"email"`
	Permissions  []string       `json:"permissions"`
	Active       bool           `json:"active"`
	LastAccessAt *string        `json:"last_access_at,omitempty"`
	CreatedAt    string         `json:"created_at"`
	UpdatedAt    string         `json:"updated_at"`
	Role         *RoleResponse  `json:"role,omitempty"`
}

// RoleResponse DTO para informações básicas do role
type RoleResponse struct {
	Id             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	DisplayName    string    `json:"display_name"`
	HierarchyLevel int       `json:"hierarchy_level"`
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

	// Construir resposta com informações de role
	adminResp := AdminResponse{
		Id:          admin.Id,
		Name:        admin.Name,
		Email:       admin.Email,
		Permissions: []string{},
		Active:      admin.Active,
		CreatedAt:   admin.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   admin.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if admin.LastAccessAt != nil {
		lastAccess := admin.LastAccessAt.Format("2006-01-02T15:04:05Z")
		adminResp.LastAccessAt = &lastAccess
	}

	// Buscar role do admin e suas permissões
	adminRoles, err := r.handler.GetAdminRoles(admin.Id.String())
	if err == nil && len(adminRoles) > 0 {
		for _, ar := range adminRoles {
			if ar.Active && ar.Role != nil {
				adminResp.Role = &RoleResponse{
					Id:             ar.Role.Id,
					Name:           ar.Role.Name,
					DisplayName:    ar.Role.DisplayName,
					HierarchyLevel: ar.Role.HierarchyLevel,
				}
				// Buscar permissões do role
				permissions, _ := r.handler.GetPermissionsFromRole(ar.RoleId.String())
				adminResp.Permissions = permissions
				break
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": adminResp})
}

func (r *ResourceAdminUsers) ServiceListAdmins(c *gin.Context) {
	admins, err := r.handler.ListAdmins()
	if err != nil {
		utils.SendInternalServerError(c, "Error listing admins", err)
		return
	}

	// Construir resposta com informações de role
	var response []AdminResponse
	for _, admin := range admins {
		adminResp := AdminResponse{
			Id:          admin.Id,
			Name:        admin.Name,
			Email:       admin.Email,
			Permissions: []string{},
			Active:      admin.Active,
			CreatedAt:   admin.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   admin.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}

		if admin.LastAccessAt != nil {
			lastAccess := admin.LastAccessAt.Format("2006-01-02T15:04:05Z")
			adminResp.LastAccessAt = &lastAccess
		}

		// Buscar role do admin e suas permissões
		adminRoles, err := r.handler.GetAdminRoles(admin.Id.String())
		if err == nil && len(adminRoles) > 0 {
			// Pegar o primeiro role ativo (normalmente só há um)
			for _, ar := range adminRoles {
				if ar.Active && ar.Role != nil {
					adminResp.Role = &RoleResponse{
						Id:             ar.Role.Id,
						Name:           ar.Role.Name,
						DisplayName:    ar.Role.DisplayName,
						HierarchyLevel: ar.Role.HierarchyLevel,
					}
					// Buscar permissões do role
					permissions, _ := r.handler.GetPermissionsFromRole(ar.RoleId.String())
					adminResp.Permissions = permissions
					break
				}
			}
		}

		response = append(response, adminResp)
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (r *ResourceAdminUsers) ServiceCreateAdmin(c *gin.Context) {
	var request CreateAdminRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		if strings.Contains(err.Error(), "RoleId") {
			utils.SendBadRequestError(c, "role_id é obrigatório. Informe o ID do cargo.", err)
			return
		}
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// 1. Buscar role
	role, err := r.handler.GetRole(request.RoleId)
	if err != nil {
		utils.SendBadRequestError(c, "Cargo não encontrado", err)
		return
	}

	// 2. Validar que é role de escopo "admin"
	if role.Scope != "admin" {
		utils.SendBadRequestError(c, "Cargo deve ser do escopo admin", nil)
		return
	}

	// 3. Validar HierarchyLevel (admin logado deve ter nível >= ao role)
	actorId, _ := c.Get("user_id")
	actorLevel, err := r.handler.GetUserMaxHierarchyLevel(actorId.(string))
	if err != nil {
		utils.SendInternalServerError(c, "Erro ao verificar nível de hierarquia", err)
		return
	}
	if role.HierarchyLevel > actorLevel {
		utils.SendForbiddenError(c, "Você não pode atribuir um cargo com nível maior que o seu")
		return
	}

	// 4. Buscar permissões do role para popular campo legado
	permissions, err := r.handler.GetPermissionsFromRole(request.RoleId)
	if err != nil || len(permissions) == 0 {
		utils.SendBadRequestError(c, "O cargo selecionado não possui permissões configuradas. Configure as permissões do cargo antes de usá-lo.", nil)
		return
	}

	// 5. Criar admin (permissões vêm via role, não diretamente no admin)
	admin := &models.Admin{
		Id:       uuid.New(),
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
		Active:   true,
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

	// 6. Criar AdminRole para o admin
	roleUUID, _ := uuid.Parse(request.RoleId)
	adminRole := &models.AdminRole{
		AdminId:        admin.Id,
		RoleId:         roleUUID,
		OrganizationId: nil,
		Active:         true,
	}
	if err := r.handler.AssignRoleToAdmin(adminRole); err != nil {
		// Se falhar ao criar AdminRole, deletar o admin criado
		r.handler.DeleteAdmin(admin.Id.String())
		utils.SendInternalServerError(c, "Erro ao atribuir cargo ao admin", err)
		return
	}

	// Construir resposta com informações de role
	adminResp := AdminResponse{
		Id:          admin.Id,
		Name:        admin.Name,
		Email:       admin.Email,
		Permissions: permissions, // Permissões do role
		Active:      admin.Active,
		CreatedAt:   admin.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   admin.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		Role: &RoleResponse{
			Id:             role.Id,
			Name:           role.Name,
			DisplayName:    role.DisplayName,
			HierarchyLevel: role.HierarchyLevel,
		},
	}

	utils.SendCreatedSuccess(c, "Admin created successfully", adminResp)
}

func (r *ResourceAdminUsers) ServiceUpdateAdmin(c *gin.Context) {
	id, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "admin")
	if !ok {
		return
	}

	var request UpdateAdminRequest
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

	// Validar hierarquia - só pode editar admin com nível <= ao seu
	actorId, _ := c.Get("user_id")
	actorLevel, err := r.handler.GetUserMaxHierarchyLevel(actorId.(string))
	if err != nil {
		utils.SendInternalServerError(c, "Erro ao verificar nível de hierarquia", err)
		return
	}

	targetLevel, _ := r.handler.GetUserMaxHierarchyLevel(existing.Id.String())
	if targetLevel > actorLevel {
		utils.SendForbiddenError(c, "Você não pode editar um admin com nível maior que o seu")
		return
	}

	// Se forneceu novo role_id, validar e atualizar
	if request.RoleId != "" {
		role, err := r.handler.GetRole(request.RoleId)
		if err != nil {
			utils.SendBadRequestError(c, "Cargo não encontrado", err)
			return
		}

		if role.Scope != "admin" {
			utils.SendBadRequestError(c, "Cargo deve ser do escopo admin", nil)
			return
		}

		if role.HierarchyLevel > actorLevel {
			utils.SendForbiddenError(c, "Você não pode atribuir um cargo com nível maior que o seu")
			return
		}

		// Atualizar AdminRole (permissões vêm via role)
		roleUUID, _ := uuid.Parse(request.RoleId)
		adminRole := &models.AdminRole{
			AdminId:        existing.Id,
			RoleId:         roleUUID,
			OrganizationId: nil,
			Active:         true,
		}
		if err := r.handler.AssignRoleToAdmin(adminRole); err != nil {
			utils.SendInternalServerError(c, "Erro ao atualizar cargo do admin", err)
			return
		}
	}

	// Atualizar campos básicos
	existing.Name = request.Name
	existing.Email = request.Email
	if request.Password != "" {
		existing.Password = request.Password
	}
	if request.Active != nil {
		existing.Active = *request.Active
	}

	if err := r.handler.UpdateAdmin(existing); err != nil {
		utils.SendInternalServerError(c, "Error updating admin", err)
		return
	}

	// Construir resposta com informações de role
	adminResp := AdminResponse{
		Id:          existing.Id,
		Name:        existing.Name,
		Email:       existing.Email,
		Permissions: []string{},
		Active:      existing.Active,
		CreatedAt:   existing.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   existing.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if existing.LastAccessAt != nil {
		lastAccess := existing.LastAccessAt.Format("2006-01-02T15:04:05Z")
		adminResp.LastAccessAt = &lastAccess
	}

	// Buscar role do admin e suas permissões
	adminRoles, err := r.handler.GetAdminRoles(existing.Id.String())
	if err == nil && len(adminRoles) > 0 {
		for _, ar := range adminRoles {
			if ar.Active && ar.Role != nil {
				adminResp.Role = &RoleResponse{
					Id:             ar.Role.Id,
					Name:           ar.Role.Name,
					DisplayName:    ar.Role.DisplayName,
					HierarchyLevel: ar.Role.HierarchyLevel,
				}
				// Buscar permissões do role
				permissions, _ := r.handler.GetPermissionsFromRole(ar.RoleId.String())
				adminResp.Permissions = permissions
				break
			}
		}
	}

	utils.SendOKSuccess(c, "Admin updated successfully", adminResp)
}

func (r *ResourceAdminUsers) ServiceDeleteAdmin(c *gin.Context) {
	id, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "admin")
	if !ok {
		return
	}

	// Verificar se admin existe
	existing, err := r.handler.GetAdminById(id.String())
	if err != nil || existing == nil {
		utils.SendNotFoundError(c, "Admin")
		return
	}

	// Validar hierarquia - só pode deletar admin com nível <= ao seu
	actorId, _ := c.Get("user_id")
	actorLevel, err := r.handler.GetUserMaxHierarchyLevel(actorId.(string))
	if err != nil {
		utils.SendInternalServerError(c, "Erro ao verificar nível de hierarquia", err)
		return
	}

	targetLevel, _ := r.handler.GetUserMaxHierarchyLevel(existing.Id.String())
	if targetLevel > actorLevel {
		utils.SendForbiddenError(c, "Você não pode deletar um admin com nível maior que o seu")
		return
	}

	if err := r.handler.DeleteAdmin(id.String()); err != nil {
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
