package main

import (
	"fmt"
	"lep/resource"
	"log"
)

func main() {
	// Conectar ao banco
	db, err := resource.OpenConnDBPostgres()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco: %v", err)
	}

	// Buscar user_organizations do Pablo
	type UserOrg struct {
		ID             string
		UserID         string
		OrganizationID string
		Role           string
		Active         bool
	}

	var userOrgs []UserOrg
	err = db.Table("user_organizations").
		Where("user_id = ?", "123e4567-e89b-12d3-a456-426614174010").
		Scan(&userOrgs).Error

	if err != nil {
		log.Fatalf("Erro ao buscar user_organizations: %v", err)
	}

	fmt.Println("User Organizations para Pablo (123e4567-e89b-12d3-a456-426614174010):")
	fmt.Printf("Total: %d\n\n", len(userOrgs))
	for _, uo := range userOrgs {
		fmt.Printf("ID: %s\n", uo.ID)
		fmt.Printf("User ID: %s\n", uo.UserID)
		fmt.Printf("Organization ID: %s\n", uo.OrganizationID)
		fmt.Printf("Role: %s\n", uo.Role)
		fmt.Printf("Active: %t\n", uo.Active)
		fmt.Println("---")
	}

	// Buscar user_projects do Pablo
	type UserProj struct {
		ID        string
		UserID    string
		ProjectID string
		Role      string
		Active    bool
	}

	var userProjs []UserProj
	err = db.Table("user_projects").
		Where("user_id = ?", "123e4567-e89b-12d3-a456-426614174010").
		Scan(&userProjs).Error

	if err != nil {
		log.Fatalf("Erro ao buscar user_projects: %v", err)
	}

	fmt.Println("\nUser Projects para Pablo (123e4567-e89b-12d3-a456-426614174010):")
	fmt.Printf("Total: %d\n\n", len(userProjs))
	for _, up := range userProjs {
		fmt.Printf("ID: %s\n", up.ID)
		fmt.Printf("User ID: %s\n", up.UserID)
		fmt.Printf("Project ID: %s\n", up.ProjectID)
		fmt.Printf("Role: %s\n", up.Role)
		fmt.Printf("Active: %t\n", up.Active)
		fmt.Println("---")
	}
}
