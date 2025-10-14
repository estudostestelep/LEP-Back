package main

import (
	"fmt"
	"lep/repositories/models"
	"lep/resource"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Println("\n🔧 Criando Master Admins diretamente no banco...")

	// Conectar ao banco
	db, err := resource.OpenConnDBPostgres()
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco: %v", err)
	}

	// Buscar a organização existente
	var org models.Organization
	if err := db.First(&org).Error; err != nil {
		log.Fatalf("❌ Nenhuma organização encontrada. Execute o seeder primeiro.")
	}
	orgID := org.Id
	fmt.Printf("📋 Organização encontrada: %s (ID: %s)\n", org.Name, orgID)

	// Buscar o projeto existente (qualquer um)
	var project models.Project
	if err := db.First(&project).Error; err != nil {
		log.Fatalf("❌ Nenhum projeto encontrado. Execute o seeder primeiro.")
	}
	projectID := project.Id
	fmt.Printf("📁 Projeto encontrado: %s (ID: %s, Org: %s)\n\n", project.Name, projectID, project.OrganizationId)

	masterAdmins := []struct {
		ID    uuid.UUID
		Name  string
		Email string
		Role  string
	}{
		{
			ID:    uuid.MustParse("123e4567-e89b-12d3-a456-426614174010"),
			Name:  "Pablo Master Admin",
			Email: "pablo@lep.com",
			Role:  "owner",
		},
		{
			ID:    uuid.MustParse("123e4567-e89b-12d3-a456-426614174011"),
			Name:  "Luan Master Admin",
			Email: "luan@lep.com",
			Role:  "admin",
		},
		{
			ID:    uuid.MustParse("123e4567-e89b-12d3-a456-426614174012"),
			Name:  "Eduardo Master Admin",
			Email: "eduardo@lep.com",
			Role:  "admin",
		},
	}

	// Hash da senha "senha123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Falha ao hashear senha: %v", err)
	}

	now := time.Now()

	for _, admin := range masterAdmins {
		fmt.Printf("  Criando %s...\n", admin.Name)

		// Criar usuário
		user := models.User{
			Id:          admin.ID,
			Name:        admin.Name,
			Email:       admin.Email,
			Password:    string(hashedPassword),
			Permissions: pq.StringArray{"master_admin"},
			Active:      true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		// Upsert no banco (primeiro tenta pelo email)
		var existingUser models.User
		result := db.Where("email = ?", admin.Email).First(&existingUser)
		if result.Error != nil {
			// Não existe, criar
			if err := db.Create(&user).Error; err != nil {
				log.Printf("    ❌ Erro ao criar usuário: %v", err)
				continue
			}
			fmt.Printf("    ✅ Usuário criado: %s\n", admin.Email)
		} else {
			// Já existe, atualizar
			if err := db.Model(&models.User{}).Where("email = ?", admin.Email).Updates(map[string]interface{}{
				"name":        user.Name,
				"password":    user.Password,
				"permissions": user.Permissions,
				"active":      user.Active,
				"updated_at":  now,
			}).Error; err != nil {
				log.Printf("    ❌ Erro ao atualizar usuário: %v", err)
				continue
			}
			fmt.Printf("    ✅ Usuário atualizado: %s\n", admin.Email)
			// Usar o ID do usuário existente
			admin.ID = existingUser.Id
		}

		// Criar relacionamento user-organization
		userOrg := models.UserOrganization{
			Id:             uuid.New(),
			UserId:         admin.ID,
			OrganizationId: orgID,
			Role:           admin.Role,
			Active:         true,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		result = db.Where("user_id = ? AND organization_id = ?", admin.ID, orgID).First(&models.UserOrganization{})
		if result.Error != nil {
			// Não existe, criar
			if err := db.Create(&userOrg).Error; err != nil {
				log.Printf("    ❌ Erro ao criar user-organization: %v", err)
				continue
			}
			fmt.Printf("    ✅ User-Organization criado (role: %s)\n", admin.Role)
		} else {
			// Já existe, atualizar
			if err := db.Model(&models.UserOrganization{}).
				Where("user_id = ? AND organization_id = ?", admin.ID, orgID).
				Updates(map[string]interface{}{
					"role":       admin.Role,
					"active":     true,
					"updated_at": now,
				}).Error; err != nil {
				log.Printf("    ❌ Erro ao atualizar user-organization: %v", err)
				continue
			}
			fmt.Printf("    ✅ User-Organization atualizado (role: %s)\n", admin.Role)
		}

		// Criar relacionamento user-project
		userProj := models.UserProject{
			Id:        uuid.New(),
			UserId:    admin.ID,
			ProjectId: projectID,
			Role:      "admin",
			Active:    true,
			CreatedAt: now,
			UpdatedAt: now,
		}

		result = db.Where("user_id = ? AND project_id = ?", admin.ID, projectID).First(&models.UserProject{})
		if result.Error != nil {
			// Não existe, criar
			if err := db.Create(&userProj).Error; err != nil {
				log.Printf("    ❌ Erro ao criar user-project: %v", err)
				continue
			}
			fmt.Printf("    ✅ User-Project criado (role: admin)\n")
		} else {
			// Já existe, atualizar
			if err := db.Model(&models.UserProject{}).
				Where("user_id = ? AND project_id = ?", admin.ID, projectID).
				Updates(map[string]interface{}{
					"role":       "admin",
					"active":     true,
					"updated_at": now,
				}).Error; err != nil {
				log.Printf("    ❌ Erro ao atualizar user-project: %v", err)
				continue
			}
			fmt.Printf("    ✅ User-Project atualizado (role: admin)\n")
		}
	}

	fmt.Println("\n✅ Master Admins criados/atualizados com sucesso!")
	fmt.Println("\n📋 Credenciais de Login:")
	fmt.Println("  • pablo@lep.com / senha123")
	fmt.Println("  • luan@lep.com / senha123")
	fmt.Println("  • eduardo@lep.com / senha123")
}
