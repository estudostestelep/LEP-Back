package server

import (
	"lep/config"
	"lep/handler"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type ResourceAuthAdmin struct {
	handler *handler.Handlers
}

type IServerAuthAdmin interface {
	ServiceAdminLogin(c *gin.Context)
	ServiceAdminLogout(c *gin.Context)
}

// ServiceAdminLogin lida com login de administradores
func (r *ResourceAuthAdmin) ServiceAdminLogin(c *gin.Context) {
	var loginData handler.AdminLoginRequest

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(400, gin.H{"error": "Dados de login inválidos"})
		return
	}

	// Validar credenciais do admin
	admin, err := r.handler.HandlerAdminUser.ValidateAdminCredentials(loginData.Email, loginData.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "Credenciais inválidas"})
		return
	}

	// Criar token JWT com tipo "admin"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":     admin.Email,
		"user_id":   admin.Id.String(),
		"user_type": "admin",
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.JWT_SECRET_PRIVATE_KEY))
	if err != nil {
		c.JSON(500, gin.H{"error": "Erro ao criar token JWT"})
		return
	}

	// Atualizar último acesso
	_ = r.handler.HandlerAdminUser.UpdateLastAccess(admin.Id.String())

	// Registrar log de acesso
	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()
	if forwardedFor := c.GetHeader("X-Forwarded-For"); forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			clientIP = strings.TrimSpace(ips[0])
		}
	}
	_ = r.handler.HandlerAuth.RecordAccessLog(admin.Id.String(), clientIP, userAgent)

	// Buscar permissões do admin via roles
	var permissions []string
	adminRoles, _ := r.handler.HandlerAdminUser.GetAdminRoles(admin.Id.String())
	if len(adminRoles) > 0 {
		for _, ar := range adminRoles {
			if ar.Active {
				perms, _ := r.handler.HandlerAdminUser.GetPermissionsFromRole(ar.RoleId.String())
				permissions = append(permissions, perms...)
				break // Usar o primeiro role ativo
			}
		}
	}

	// Preparar resposta (sem senha)
	adminResponse := struct {
		Id          string   `json:"id"`
		Name        string   `json:"name"`
		Email       string   `json:"email"`
		Permissions []string `json:"permissions"`
		Active      bool     `json:"active"`
	}{
		Id:          admin.Id.String(),
		Name:        admin.Name,
		Email:       admin.Email,
		Permissions: permissions,
		Active:      admin.Active,
	}

	c.JSON(200, handler.AdminLoginResponse{
		Admin:       adminResponse,
		Token:       tokenString,
		UserType:    "admin",
		Permissions: permissions,
	})
}

// ServiceAdminLogout lida com logout de administradores
func (r *ResourceAuthAdmin) ServiceAdminLogout(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	token := strings.TrimPrefix(tokenString, "Bearer ")

	err := r.handler.HandlerAuth.Logout(token)
	if err != nil {
		c.JSON(200, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Logout realizado com sucesso"})
}

func NewSourceServerAuthAdmin(handler *handler.Handlers) IServerAuthAdmin {
	return &ResourceAuthAdmin{handler: handler}
}
