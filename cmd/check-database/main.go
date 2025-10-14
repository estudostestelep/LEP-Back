package main

import (
	"fmt"
	"lep/repositories/models"
	"lep/resource"
	"log"
)

func main() {
	fmt.Println("\n🔍 Verificando estrutura do banco de dados...")

	// Conectar ao banco
	db, err := resource.OpenConnDBPostgres()
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco: %v", err)
	}

	// Verificar Organizations
	var orgs []models.Organization
	db.Find(&orgs)
	fmt.Printf("\n📋 Organizations: %d encontradas\n", len(orgs))
	for _, org := range orgs {
		fmt.Printf("  • %s (ID: %s, Email: %s)\n", org.Name, org.Id, org.Email)
	}

	// Verificar Projects
	var projects []models.Project
	db.Find(&projects)
	fmt.Printf("\n📁 Projects: %d encontrados\n", len(projects))
	for _, proj := range projects {
		fmt.Printf("  • %s (ID: %s, Org: %s)\n", proj.Name, proj.Id, proj.OrganizationId)
	}

	// Verificar Users
	var users []models.User
	db.Find(&users)
	fmt.Printf("\n👥 Users: %d encontrados\n", len(users))
	for _, user := range users {
		fmt.Printf("  • %s (%s) - Permissions: %v\n", user.Name, user.Email, user.Permissions)
	}

	// Verificar UserOrganizations
	var userOrgs []models.UserOrganization
	db.Find(&userOrgs)
	fmt.Printf("\n🔗 User-Organizations: %d encontrados\n", len(userOrgs))
	for _, uo := range userOrgs {
		fmt.Printf("  • User: %s, Org: %s, Role: %s\n", uo.UserId, uo.OrganizationId, uo.Role)
	}

	// Verificar UserProjects
	var userProjs []models.UserProject
	db.Find(&userProjs)
	fmt.Printf("\n🔗 User-Projects: %d encontrados\n", len(userProjs))
	for _, up := range userProjs {
		fmt.Printf("  • User: %s, Project: %s, Role: %s\n", up.UserId, up.ProjectId, up.Role)
	}

	// Verificar Master Admins especificamente
	fmt.Println("\n🔴 Master Admins:")
	var masterAdmins []models.User
	db.Where("'master_admin' = ANY(permissions)").Find(&masterAdmins)
	for _, admin := range masterAdmins {
		fmt.Printf("  • %s (%s)\n", admin.Name, admin.Email)

		// Ver suas organizações
		var orgs []models.UserOrganization
		db.Where("user_id = ?", admin.Id).Find(&orgs)
		for _, org := range orgs {
			fmt.Printf("    - Org: %s (Role: %s)\n", org.OrganizationId, org.Role)
		}

		// Ver seus projetos
		var projs []models.UserProject
		db.Where("user_id = ?", admin.Id).Find(&projs)
		for _, proj := range projs {
			fmt.Printf("    - Project: %s (Role: %s)\n", proj.ProjectId, proj.Role)
		}
	}

	fmt.Println()
}
