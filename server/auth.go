package server

// Importe a biblioteca jwt-go
import (
	"fmt"
	"lep/config"
	"lep/handler"
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

func (r *ResourceAuth) ServiceLogin(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&loginData); err != nil {
		c.JSON(400, gin.H{"error": "Erro ao decodificar dados de login"})
		return
	}

	user, err := r.handler.HandlerUser.GetUserByEmail(loginData.Email)
	if err != nil || user == nil {
		c.JSON(401, gin.H{"error": "Credenciais inválidas"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Credenciais inválidas"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Expira em 24 horas
	})

	tokenString, err := token.SignedString([]byte(config.JWT_SECRET_PRIVATE_KEY))
	if err != nil {
		c.JSON(500, gin.H{"error": "Erro ao criar token JWT"})
		return
	}

	err = r.handler.HandlerAuth.PostToken(user, tokenString)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	c.JSON(200, gin.H{"user": user, "token": tokenString})
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
		return []byte(config.JWT_SECRET_PUBLIC_KEY), nil
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
