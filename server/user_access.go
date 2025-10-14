package server

import (
	"lep/handler"
	"lep/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserAccessServer struct {
	Handler handler.UserAccessHandler
}

func NewUserAccessServer(h handler.UserAccessHandler) *UserAccessServer {
	return &UserAccessServer{Handler: h}
}

// isMasterAdmin verifica se um usuário tem permissão de Master Admin
func isMasterAdmin(c *gin.Context) bool {
	userPermissions, exists := c.Get("user_permissions")
	if !exists {
		return false
	}

	// ✅ CORREÇÃO: Aceitar tanto []string quanto pq.StringArray
	var permissions []string
	switch v := userPermissions.(type) {
	case []string:
		permissions = v
	case pq.StringArray:
		permissions = []string(v)
	default:
		return false
	}

	for _, p := range permissions {
		if p == "master_admin" || p == "all" {
			return true
		}
	}
	return false
}

// ServiceGetUserOrganizationsAndProjects retorna organizações e projetos de um usuário
// GET /user/:id/organizations-projects
func (s *UserAccessServer) ServiceGetUserOrganizationsAndProjects(c *gin.Context) {
	// Verificar se é Master Admin
	if !isMasterAdmin(c) {
		utils.SendForbiddenError(c, "Acesso negado: apenas Master Admin pode visualizar acessos de usuários")
		return
	}

	// Obter ID do usuário da URL
	userIdStr := c.Param("id")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		utils.SendBadRequestError(c, "ID de usuário inválido", err)
		return
	}

	// Buscar acessos do usuário
	response, err := s.Handler.GetUserOrganizationsAndProjects(userId)
	if err != nil {
		if err.Error() == "usuário não encontrado" {
			utils.SendNotFoundError(c, "Usuário")
			return
		}
		utils.SendInternalServerError(c, "Erro ao buscar acessos do usuário", err)
		return
	}

	utils.SendOKSuccess(c, "Acessos do usuário obtidos com sucesso", response)
}

// ServiceUpdateUserOrganizationsAndProjects atualiza organizações e projetos de um usuário
// POST /user/:id/organizations-projects
func (s *UserAccessServer) ServiceUpdateUserOrganizationsAndProjects(c *gin.Context) {
	// Verificar se é Master Admin
	if !isMasterAdmin(c) {
		utils.SendForbiddenError(c, "Acesso negado: apenas Master Admin pode gerenciar acessos de usuários")
		return
	}

	// Obter ID do usuário da URL
	userIdStr := c.Param("id")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		utils.SendBadRequestError(c, "ID de usuário inválido", err)
		return
	}

	// Fazer bind do request body
	var req handler.UpdateUserAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendBadRequestError(c, "Request body inválido", err)
		return
	}

	// Atualizar acessos
	response, err := s.Handler.UpdateUserOrganizationsAndProjects(userId, &req)
	if err != nil {
		if err.Error() == "usuário não encontrado" {
			utils.SendNotFoundError(c, "Usuário")
			return
		}
		if err.Error() == "uma ou mais organizações não existem" ||
			err.Error() == "um ou mais projetos não existem" ||
			err.Error() == "um ou mais projetos não pertencem às organizações selecionadas" {
			utils.SendBadRequestError(c, err.Error(), err)
			return
		}
		utils.SendInternalServerError(c, "Erro ao atualizar acessos do usuário", err)
		return
	}

	c.JSON(http.StatusOK, response)
}
