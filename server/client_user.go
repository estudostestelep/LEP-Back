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

type ResourceClientUsers struct {
	handler     handler.IHandlerClientUser
	roleHandler *handler.RoleHandler
}

type IServerClientUsers interface {
	ServiceGetClient(c *gin.Context)
	ServiceListClients(c *gin.Context)
	ServiceCreateClient(c *gin.Context)
	ServiceUpdateClient(c *gin.Context)
	ServiceDeleteClient(c *gin.Context)
}

// CreateClientRequest DTO para criar client
type CreateClientRequest struct {
	Name        string   `json:"name" binding:"required"`
	Email       string   `json:"email" binding:"required,email"`
	Password    string   `json:"password" binding:"required,min=6"`
	OrgId       string   `json:"org_id" binding:"required"`
	ProjIds     []string `json:"proj_ids"`
	Permissions []string `json:"permissions"`
	Active      *bool    `json:"active"`
}

// UpdateClientRequest DTO para atualizar client
// OrgId é usado para identificar o contexto da organização, não para alterar
type UpdateClientRequest struct {
	Name        string   `json:"name" binding:"required"`
	Email       string   `json:"email" binding:"required,email"`
	Password    string   `json:"password" binding:"omitempty,min=6"`
	OrgId       string   `json:"org_id"` // Para contexto, não altera o OrgId do cliente
	ProjIds     []string `json:"proj_ids"`
	Permissions []string `json:"permissions"`
	Active      *bool    `json:"active"`
}

func (r *ResourceClientUsers) ServiceGetClient(c *gin.Context) {
	id, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "client")
	if !ok {
		return
	}

	client, err := r.handler.GetClientById(id.String())
	if err != nil {
		utils.SendInternalServerError(c, "Error getting client", err)
		return
	}

	if client == nil {
		utils.SendNotFoundError(c, "Client")
		return
	}

	// Remover senha da resposta
	client.Password = ""

	c.JSON(http.StatusOK, gin.H{"data": client})
}

func (r *ResourceClientUsers) ServiceListClients(c *gin.Context) {
	// Pegar org_id e project_id do header para filtrar clientes
	orgId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	var clients []models.Client
	var err error

	// Filtrar por projeto se disponível, senão por organização
	if projectId != "" && orgId != "" {
		clients, err = r.handler.ListClientsByProject(orgId, projectId)
	} else if orgId != "" {
		clients, err = r.handler.ListClientsByOrgId(orgId)
	} else {
		clients, err = r.handler.ListClients()
	}

	if err != nil {
		utils.SendInternalServerError(c, "Error listing clients", err)
		return
	}

	// DTO com informações de role
	type ClientWithRole struct {
		Id                 uuid.UUID      `json:"id"`
		Name               string         `json:"name"`
		Email              string         `json:"email"`
		Active             bool           `json:"active"`
		OrgId              uuid.UUID      `json:"org_id"`
		ProjIds            pq.StringArray `json:"proj_ids"`
		LastAccessAt       interface{}    `json:"last_access_at"`
		CreatedAt          interface{}    `json:"created_at"`
		UpdatedAt          interface{}    `json:"updated_at"`
		RoleDisplayName    string         `json:"role_display_name"`
		RoleHierarchyLevel int            `json:"role_hierarchy_level"`
	}

	result := make([]ClientWithRole, 0, len(clients))
	for _, client := range clients {
		cwr := ClientWithRole{
			Id:           client.Id,
			Name:         client.Name,
			Email:        client.Email,
			Active:       client.Active,
			OrgId:        client.OrgId,
			ProjIds:      client.ProjIds,
			LastAccessAt: client.LastAccessAt,
			CreatedAt:    client.CreatedAt,
			UpdatedAt:    client.UpdatedAt,
		}

		// Buscar role do cliente
		if r.roleHandler != nil {
			roles, _ := r.roleHandler.GetClientRoles(client.Id.String(), orgId)
			if len(roles) > 0 && roles[0].Role != nil {
				cwr.RoleDisplayName = roles[0].Role.DisplayName
				cwr.RoleHierarchyLevel = roles[0].Role.HierarchyLevel
			}
		}

		result = append(result, cwr)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (r *ResourceClientUsers) ServiceCreateClient(c *gin.Context) {
	var request CreateClientRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Validar org_id
	orgUUID, err := uuid.Parse(request.OrgId)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid org_id format", err)
		return
	}

	// Converter para modelo (permissões vêm via roles, não diretamente no client)
	client := &models.Client{
		Id:       uuid.New(),
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
		OrgId:    orgUUID,
		ProjIds:  pq.StringArray(request.ProjIds),
		Active:   true,
	}

	if request.Active != nil {
		client.Active = *request.Active
	}

	// Validação estruturada do modelo
	if err := validation.CreateClientValidation(client); err != nil {
		utils.SendValidationError(c, "Validation failed", err)
		return
	}

	if err := r.handler.CreateClient(client); err != nil {
		// Verificar erro de validação prévia (mensagem em português)
		if strings.Contains(err.Error(), "já cadastrado") {
			utils.SendConflictError(c, "Email already registered in this organization", nil)
			return
		}
		// Verificar erro de constraint do banco (fallback para race conditions)
		if utils.IsDuplicateKeyError(err) {
			utils.SendConflictError(c, "Email already registered in this organization", nil)
			return
		}
		utils.SendInternalServerError(c, "Error creating client", err)
		return
	}

	// Remover senha da resposta
	client.Password = ""

	utils.SendCreatedSuccess(c, "Client created successfully", client)
}

func (r *ResourceClientUsers) ServiceUpdateClient(c *gin.Context) {
	id, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "client")
	if !ok {
		return
	}

	var request UpdateClientRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Verificar se client existe
	existing, err := r.handler.GetClientById(id.String())
	if err != nil || existing == nil {
		utils.SendNotFoundError(c, "Client")
		return
	}

	// Atualizar campos (OrgId é imutável, não pode ser alterado)
	existing.Name = request.Name
	existing.Email = request.Email
	if request.Password != "" {
		existing.Password = request.Password
	}
	if len(request.ProjIds) > 0 {
		existing.ProjIds = pq.StringArray(request.ProjIds)
	}
	// Nota: Permissões agora são gerenciadas via roles, não diretamente no client
	if request.Active != nil {
		existing.Active = *request.Active
	}

	// Validação estruturada do modelo
	if err := validation.UpdateClientValidation(existing); err != nil {
		utils.SendValidationError(c, "Validation failed", err)
		return
	}

	if err := r.handler.UpdateClient(existing); err != nil {
		utils.SendInternalServerError(c, "Error updating client", err)
		return
	}

	// Remover senha da resposta
	existing.Password = ""

	utils.SendOKSuccess(c, "Client updated successfully", existing)
}

func (r *ResourceClientUsers) ServiceDeleteClient(c *gin.Context) {
	id, ok := validation.ParseAndValidateUUID(c, c.Param("id"), "client")
	if !ok {
		return
	}

	if err := r.handler.DeleteClient(id.String()); err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			utils.SendNotFoundError(c, "Client")
			return
		}
		utils.SendInternalServerError(c, "Error deleting client", err)
		return
	}

	utils.SendOKSuccess(c, "Client deleted successfully", nil)
}

func NewSourceServerClientUsers(handler *handler.Handlers) IServerClientUsers {
	return &ResourceClientUsers{
		handler:     handler.HandlerClientUser,
		roleHandler: handler.HandlerRole,
	}
}
