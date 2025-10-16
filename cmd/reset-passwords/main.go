package main

import (
	"fmt"
	"lep/resource"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Println("🔐 LEP Password Reset Tool")
	fmt.Println("===========================\n")

	// Conectar ao banco
	db, err := resource.OpenConnDBPostgres()
	if err != nil {
		log.Fatalf("❌ Erro ao conectar ao banco: %v", err)
	}
	fmt.Println("✅ Conectado ao banco de dados\n")

	// Gerar hash correto para "senha123"
	password := "senha123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("❌ Erro ao gerar hash: %v", err)
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

	fmt.Println("🔄 Resetando senhas dos usuários...")
	fmt.Printf("Nova senha: %s\n", password)
	fmt.Printf("Hash bcrypt: %s\n\n", string(hashedPassword))

	// Atualizar senha de cada usuário
	for _, email := range userEmails {
		result := db.Table("users").
			Where("email = ?", email).
			Update("password", string(hashedPassword))

		if result.Error != nil {
			fmt.Printf("⚠️  Erro ao atualizar %s: %v\n", email, result.Error)
			continue
		}

		if result.RowsAffected == 0 {
			fmt.Printf("⚠️  Usuário %s não encontrado\n", email)
		} else {
			fmt.Printf("✅ Senha atualizada: %s\n", email)
		}
	}

	fmt.Println("\n✅ Processo concluído!")
	fmt.Println("\nCredenciais atualizadas:")
	fmt.Println("  🔴 Master Admins:")
	fmt.Println("    • pablo@lep.com / senha123")
	fmt.Println("    • luan@lep.com / senha123")
	fmt.Println("    • eduardo@lep.com / senha123")
	fmt.Println("")
	fmt.Println("  🟡 Demo Users:")
	fmt.Println("    • teste@gmail.com / senha123")
	fmt.Println("    • garcom1@gmail.com / senha123")
	fmt.Println("    • gerente1@gmail.com / senha123")
}