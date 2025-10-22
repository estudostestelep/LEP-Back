package server

import (
	"lep/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminController struct {
	DB *gorm.DB
}

// ServiceResetPasswords reseta as senhas de usuários específicos (endpoint administrativo temporário)
func (a *AdminController) ServiceResetPasswords(c *gin.Context) {
	db := a.DB

	// Gerar hash correto para "senha123"
	password := "senha123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.SendInternalServerError(c, "Erro ao gerar hash", err)
		return
	}

	// Usuários para resetar
	userEmails := []string{
		"pablo@lep.com",
		"luan@lep.com",
		"eduardo@lep.com",
		"teste@gmail.com",
		"garcom1@gmail.com",
		"gerente1@gmail.com",
	}

	results := make(map[string]string)

	// Atualizar senha de cada usuário
	for _, email := range userEmails {
		result := db.Table("users").
			Where("email = ? AND deleted_at IS NULL", email).
			Update("password", string(hashedPassword))

		if result.Error != nil {
			results[email] = "erro: " + result.Error.Error()
			continue
		}

		if result.RowsAffected == 0 {
			results[email] = "não encontrado"
		} else {
			results[email] = "atualizado"
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset completed",
		"results": results,
		"password": password,
	})
}