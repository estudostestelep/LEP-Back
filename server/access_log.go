package server

import (
	"lep/handler"
	"lep/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AccessLogController struct {
	handler *handler.Handlers
}

type IAccessLogServer interface {
	ServiceGetUserAccessLogs(c *gin.Context)
}

func NewAccessLogController(h *handler.Handlers) IAccessLogServer {
	return &AccessLogController{handler: h}
}

// ServiceGetUserAccessLogs retorna os logs de acesso de um usuário
// GET /user/:id/access-logs?page=1&per_page=20
func (r *AccessLogController) ServiceGetUserAccessLogs(c *gin.Context) {
	userIdStr := c.Param("id")

	// Validar formato UUID
	_, err := uuid.Parse(userIdStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid user ID format", err)
		return
	}

	// Obter parâmetros de paginação
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	// Validar limites
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	// Buscar logs de acesso
	logs, err := r.handler.HandlerAuth.GetAccessLogs(userIdStr, page, perPage)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting access logs", err)
		return
	}

	c.JSON(http.StatusOK, logs)
}
