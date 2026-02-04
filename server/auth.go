package server

// DEPRECATED: Este arquivo usa autenticação legada. Use auth_admin.go e auth_client.go.
// Mantido para compatibilidade com clientes existentes.

import (
	"fmt"
	"lep/config"
	"lep/handler"
	"lep/repositories/models"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type ResourceAuth struct {
	handler *handler.Handlers
}

type IServerAuth interface {
	ServiceLogin(c *gin.Context)
	ServiceLogout(c *gin.Context)
	ServiceValidateToken(c *gin.Context)
	ServiceValidateTokenIn(c *gin.Context) bool
}

// ServiceLogin tenta autenticar primeiro como Client, depois como Admin
func (r *ResourceAuth) ServiceLogin(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&loginData); err != nil {
		c.JSON(400, gin.H{"error": "Erro ao decodificar dados de login"})
		return
	}

	// Tentar autenticar como Client primeiro
	client, err := r.handler.HandlerClientUser.GetClientByEmail(loginData.Email)
	if err == nil && client != nil {
		if err := bcrypt.CompareHashAndPassword([]byte(client.Password), []byte(loginData.Password)); err == nil {
			r.handleClientLogin(c, client)
			return
		}
	}

	// Se não for Client, tentar como Admin
	admin, err := r.handler.HandlerAdminUser.GetAdminByEmail(loginData.Email)
	if err == nil && admin != nil {
		if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(loginData.Password)); err == nil {
			r.handleAdminLogin(c, admin)
			return
		}
	}

	c.JSON(401, gin.H{"error": "Credenciais inválidas"})
}

func (r *ResourceAuth) handleClientLogin(c *gin.Context, client *models.Client) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":     client.Email,
		"user_id":   client.Id.String(),
		"user_type": "client",
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.JWT_SECRET_PRIVATE_KEY))
	if err != nil {
		c.JSON(500, gin.H{"error": "Erro ao criar token JWT"})
		return
	}

	err = r.handler.HandlerAuth.PostTokenForUser(client.Id, client.Email, tokenString)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

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

	// Buscar organizações e projetos do client
	userOrganizations, _ := r.handler.HandlerAuth.GetClientOrganizationsWithNames(client.Id.String())
	userProjects, _ := r.handler.HandlerAuth.GetClientProjectsWithNames(client.Id.String())

	c.JSON(200, gin.H{
		"user": map[string]interface{}{
			"id":        client.Id,
			"name":      client.Name,
			"email":     client.Email,
			"user_type": "client",
		},
		"token":         tokenString,
		"organizations": userOrganizations,
		"projects":      userProjects,
		"admin_roles":   []handler.UserAdminRoleInfo{},
	})
}

func (r *ResourceAuth) handleAdminLogin(c *gin.Context, admin *models.Admin) {
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

	err = r.handler.HandlerAuth.PostTokenForUser(admin.Id, admin.Email, tokenString)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
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

	// Buscar roles do admin
	adminRoles, _ := r.handler.HandlerAuth.GetAdminRolesInfo(admin.Id.String())

	// Buscar permissões do admin via roles
	var permissions []string
	adminRolesList, _ := r.handler.HandlerAdminUser.GetAdminRoles(admin.Id.String())
	if len(adminRolesList) > 0 {
		for _, ar := range adminRolesList {
			if ar.Active {
				perms, _ := r.handler.HandlerAdminUser.GetPermissionsFromRole(ar.RoleId.String())
				permissions = append(permissions, perms...)
				break // Usar o primeiro role ativo
			}
		}
	}

	c.JSON(200, gin.H{
		"user": map[string]interface{}{
			"id":          admin.Id,
			"name":        admin.Name,
			"email":       admin.Email,
			"user_type":   "admin",
			"permissions": permissions,
		},
		"token":         tokenString,
		"organizations": []interface{}{}, // Admins têm acesso global
		"projects":      []interface{}{}, // Admins têm acesso global
		"admin_roles":   adminRoles,
	})
}

func (r *ResourceAuth) ServiceLogout(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	token := strings.TrimPrefix(tokenString, "Bearer ")

	err := r.handler.HandlerAuth.Logout(token)
	if err != nil {
		c.JSON(200, gin.H{"message": err})
		return
	}
	c.JSON(200, gin.H{"message": "Logout realizado com sucesso"})
}

func (r *ResourceAuth) ServiceValidateToken(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	fmt.Println("ServiceValidateToken", tokenString)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT_SECRET_PRIVATE_KEY), nil
	})

	if err != nil {
		fmt.Println("aqui", err)
		c.JSON(401, gin.H{"error": "Token inválido"})
		return
	}

	user, err := r.handler.HandlerAuth.VerificationToken(tokenString)
	if err != nil {
		fmt.Println("aqui2", err)
		c.JSON(401, gin.H{"error": err})
		return
	}

	c.JSON(200, gin.H{"message": "Token válido", "user": user})
}

func (r *ResourceAuth) ServiceValidateTokenIn(c *gin.Context) bool {
	tokenString := c.GetHeader("Authorization")
	fmt.Println("ServiceValidateTokenIn", tokenString)

	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	fmt.Println("tokenString2", tokenString)

	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT_SECRET_PRIVATE_KEY), nil
	})
	fmt.Println("jwt.Parse", err)

	if err != nil {
		return false
	}
	_, err = r.handler.HandlerAuth.VerificationToken(tokenString)
	if err != nil {
		return false
	}

	return true
}

func NewSourceServerAuth(handler *handler.Handlers) IServerAuth {
	return &ResourceAuth{handler: handler}
}
