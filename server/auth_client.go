package server

import (
	"lep/config"
	"lep/handler"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type ResourceAuthClient struct {
	handler *handler.Handlers
}

type IServerAuthClient interface {
	ServiceClientLogin(c *gin.Context)
	ServiceClientLogout(c *gin.Context)
}

// ServiceClientLogin lida com login de clientes
func (r *ResourceAuthClient) ServiceClientLogin(c *gin.Context) {
	var loginData handler.ClientLoginRequest

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(400, gin.H{"error": "Dados de login inválidos"})
		return
	}

	// Validar credenciais do cliente
	client, org, err := r.handler.HandlerClientUser.ValidateClientCredentials(
		loginData.Email,
		loginData.Password,
		loginData.OrgSlug,
	)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	// Criar token JWT com tipo "client"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":     client.Email,
		"user_id":   client.Id.String(),
		"user_type": "client",
		"org_id":    client.OrgId.String(),
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.JWT_SECRET_PRIVATE_KEY))
	if err != nil {
		c.JSON(500, gin.H{"error": "Erro ao criar token JWT"})
		return
	}

	// Atualizar último acesso
	_ = r.handler.HandlerClientUser.UpdateLastAccess(client.Id.String())

	// Registrar log de acesso
	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()
	if forwardedFor := c.GetHeader("X-Forwarded-For"); forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			clientIP = strings.TrimSpace(ips[0])
		}
	}
	_ = r.handler.HandlerAuth.RecordAccessLog(client.Id.String(), clientIP, userAgent)

	// Buscar projetos do cliente
	projects, _ := r.handler.HandlerClientUser.GetClientProjects(client)

	// Buscar permissões do cliente via roles
	var permissions []string
	clientRoles, _ := r.handler.HandlerRole.GetClientRoles(client.Id.String(), client.OrgId.String())
	if len(clientRoles) > 0 {
		for _, cr := range clientRoles {
			if cr.Active {
				perms, _ := r.handler.HandlerAdminUser.GetPermissionsFromRole(cr.RoleId.String())
				permissions = append(permissions, perms...)
				break // Usar o primeiro role ativo
			}
		}
	}

	// Preparar resposta (sem senha)
	clientResponse := struct {
		Id          string   `json:"id"`
		Name        string   `json:"name"`
		Email       string   `json:"email"`
		OrgId       string   `json:"org_id"`
		ProjIds     []string `json:"proj_ids"`
		Permissions []string `json:"permissions"`
		Active      bool     `json:"active"`
	}{
		Id:          client.Id.String(),
		Name:        client.Name,
		Email:       client.Email,
		OrgId:       client.OrgId.String(),
		ProjIds:     client.ProjIds,
		Permissions: permissions,
		Active:      client.Active,
	}

	c.JSON(200, handler.ClientLoginResponse{
		Client:   clientResponse,
		Token:    tokenString,
		UserType: "client",
		Organization: handler.OrganizationInfo{
			Id:   org.Id,
			Name: org.Name,
			Slug: org.Slug,
		},
		Projects:    projects,
		Permissions: permissions,
	})
}

// ServiceClientLogout lida com logout de clientes
func (r *ResourceAuthClient) ServiceClientLogout(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	token := strings.TrimPrefix(tokenString, "Bearer ")

	err := r.handler.HandlerAuth.Logout(token)
	if err != nil {
		c.JSON(200, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Logout realizado com sucesso"})
}

func NewSourceServerAuthClient(handler *handler.Handlers) IServerAuthClient {
	return &ResourceAuthClient{handler: handler}
}
