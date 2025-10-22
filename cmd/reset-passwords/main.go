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

	// Gerar hash correto para "password"
	password := "password"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("❌ Erro ao gerar hash: %v", err)
	}

	// Usuários criados pelo auto-seed
	userEmails := []string{
		"admin@lep-demo.com",
		"garcom@lep-demo.com",
		"gerente@lep-demo.com",
		"cozinha@lep-demo.com",
		"atendente@lep-demo.com",
		"owner@lep-demo.com",
	}

	fmt.Println("🔄 Resetando senhas dos usuários (auto-seed)...")
	fmt.Printf("Nova senha: %s\n", password)
	fmt.Printf("Hash bcrypt gerado\n\n")

	// Atualizar senha de cada usuário
	updatedCount := 0
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
			updatedCount++
		}
	}

	fmt.Printf("\n✅ Processo concluído! %d senhas atualizadas\n", updatedCount)
	fmt.Println("\nCredenciais atualizadas:")
	fmt.Println("  🔴 Demo Users (senha: password):")
	fmt.Println("    • admin@lep-demo.com / password (Admin)")
	fmt.Println("    • garcom@lep-demo.com / password (Garçom)")
	fmt.Println("    • gerente@lep-demo.com / password (Gerente)")
	fmt.Println("    • cozinha@lep-demo.com / password (Cozinha)")
	fmt.Println("    • atendente@lep-demo.com / password (Atendente)")
	fmt.Println("    • owner@lep-demo.com / password (Owner)")
}
