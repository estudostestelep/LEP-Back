package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"lep/utils"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var (
	baseURL     string
	verbose     bool
	environment string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "seed-remote",
		Short: "LEP Remote Database Seeder",
		Long:  `Populate the LEP database via HTTP API calls to remote server (staging/production).`,
		Run:   runRemoteSeed,
	}

	rootCmd.Flags().StringVar(&baseURL, "url", "https://lep-system-516622888070.us-central1.run.app", "Base URL of the API")
	rootCmd.Flags().StringVar(&environment, "environment", "stage", "Environment to seed (stage, prod)")
	rootCmd.Flags().BoolVar(&verbose, "verbose", false, "Enable verbose logging")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runRemoteSeed(cmd *cobra.Command, args []string) {
	fmt.Println("\n🌐 LEP Remote Database Seeder")
	fmt.Println("==============================")
	fmt.Printf("Target URL: %s\n", baseURL)
	fmt.Printf("Environment: %s\n", environment)
	fmt.Printf("Verbose: %t\n\n", verbose)

	// Testar conectividade
	fmt.Println("🔍 Testing API connectivity...")
	if err := testConnectivity(); err != nil {
		log.Fatalf("Failed to connect to API: %v", err)
	}
	fmt.Println("✅ API is reachable\n")

	// Gerar dados de seed
	fmt.Println("📊 Generating seed data...")
	seedData := utils.GenerateCompleteData()

	// Executar seeding via HTTP
	fmt.Println("🌱 Seeding database via HTTP API...")
	err := seedDatabaseViaHTTP(seedData)
	if err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	fmt.Println("\n✅ Remote database seeding completed successfully!")
	fmt.Println("🚀 You can now access the application and login.")
	fmt.Println("\nCredentials:")
	fmt.Println("  🔴 Master Admins:")
	fmt.Println("    • pablo@lep.com / senha123")
	fmt.Println("    • luan@lep.com / senha123")
	fmt.Println("    • eduardo@lep.com / senha123")
	fmt.Println("")
	fmt.Println("  🟡 Demo Users:")
	fmt.Println("    • teste@gmail.com / password")
	fmt.Println("    • garcom1@gmail.com / password")
	fmt.Println("    • gerente1@gmail.com / password")
}

// Testar conectividade com a API
func testConnectivity() error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(baseURL + "/ping")
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if verbose {
		fmt.Printf("  Response: %s\n", string(body))
	}

	return nil
}

// Executar seeding completo via HTTP
func seedDatabaseViaHTTP(data *utils.SeedData) error {
	client := &http.Client{Timeout: 30 * time.Second}

	// 1. Criar organização (sem headers)
	fmt.Println("  📋 Criando organização...")
	orgId, err := createOrganizationHTTP(client, data.Organizations[0])
	if err != nil {
		return fmt.Errorf("failed to create organization: %v", err)
	}

	// 2. Criar projeto (apenas org header)
	fmt.Println("  📁 Criando projeto...")
	projectId, err := createProjectHTTP(client, data.Projects[0], orgId)
	if err != nil {
		return fmt.Errorf("failed to create project: %v", err)
	}

	// 3. Criar usuário admin e obter token
	fmt.Println("  👥 Criando usuário admin...")
	adminToken, err := createAdminUserHTTP(client, data.Users[0], orgId, projectId, data.UserOrganizations[0], data.UserProjects[0])
	if err != nil {
		return fmt.Errorf("failed to create admin user: %v", err)
	}

	// 4. Configurar headers para requisições autenticadas
	headers := map[string]string{
		"Authorization":         "Bearer " + adminToken,
		"X-Lpe-Organization-Id": orgId.String(),
		"X-Lpe-Project-Id":      projectId.String(),
	}

	// 5. Criar demais usuários
	if len(data.Users) > 1 {
		fmt.Println("  👥 Criando demais usuários...")
		for i := 1; i < len(data.Users); i++ {
			user := data.Users[i]
			userOrg := data.UserOrganizations[i]
			userProj := data.UserProjects[i]

			if err := createUserHTTP(client, user, orgId, projectId, userOrg, userProj, headers); err != nil {
				return fmt.Errorf("failed to create user %s: %v", user.Name, err)
			}
		}
	}

	// 6. Criar environments
	if len(data.Environments) > 0 {
		fmt.Println("  🏢 Criando environments...")
		for _, env := range data.Environments {
			env.OrganizationId = orgId
			env.ProjectId = projectId
			if err := createResourceHTTP(client, "/environment", env, headers); err != nil {
				return fmt.Errorf("failed to create environment %s: %v", env.Name, err)
			}
		}
	}

	// 7. Criar customers
	if len(data.Customers) > 0 {
		fmt.Println("  👤 Criando customers...")
		for _, customer := range data.Customers {
			customer.OrganizationId = orgId
			customer.ProjectId = projectId
			if err := createResourceHTTP(client, "/customer", customer, headers); err != nil {
				return fmt.Errorf("failed to create customer %s: %v", customer.Name, err)
			}
		}
	}

	// 8. Criar menus
	if len(data.Menus) > 0 {
		fmt.Println("  📖 Criando menus...")
		for _, menu := range data.Menus {
			menu.OrganizationId = orgId
			menu.ProjectId = projectId
			if err := createResourceHTTP(client, "/menu", menu, headers); err != nil {
				return fmt.Errorf("failed to create menu %s: %v", menu.Name, err)
			}
		}
	}

	// 9. Criar categories
	if len(data.Categories) > 0 {
		fmt.Println("  📂 Criando categories...")
		for _, category := range data.Categories {
			category.OrganizationId = orgId
			category.ProjectId = projectId
			if err := createResourceHTTP(client, "/category", category, headers); err != nil {
				return fmt.Errorf("failed to create category %s: %v", category.Name, err)
			}
		}
	}

	// 10. Criar tags
	if len(data.Tags) > 0 {
		fmt.Println("  🏷️  Criando tags...")
		for _, tag := range data.Tags {
			tag.OrganizationId = orgId
			tag.ProjectId = projectId
			if err := createResourceHTTP(client, "/tag", tag, headers); err != nil {
				return fmt.Errorf("failed to create tag %s: %v", tag.Name, err)
			}
		}
	}

	// 11. Criar products
	if len(data.Products) > 0 {
		fmt.Println("  🍽️  Criando products...")
		for _, product := range data.Products {
			product.OrganizationId = orgId
			product.ProjectId = projectId
			if err := createResourceHTTP(client, "/product", product, headers); err != nil {
				return fmt.Errorf("failed to create product %s: %v", product.Name, err)
			}
		}
	}

	// 12. Criar product tags
	if len(data.ProductTags) > 0 {
		fmt.Println("  🔗 Criando product tags...")
		for _, productTag := range data.ProductTags {
			endpoint := fmt.Sprintf("/product/%s/tags", productTag.ProductId.String())
			requestBody := map[string]string{
				"tag_id": productTag.TagId.String(),
			}
			if err := createResourceHTTP(client, endpoint, requestBody, headers); err != nil {
				// Ignorar erros de duplicação para tags
				if verbose {
					fmt.Printf("    ⚠️  Warning creating product tag: %v\n", err)
				}
			}
		}
	}

	// 13. Criar tables
	if len(data.Tables) > 0 {
		fmt.Println("  🪑 Criando tables...")
		for _, table := range data.Tables {
			table.OrganizationId = orgId
			table.ProjectId = projectId
			if err := createResourceHTTP(client, "/table", table, headers); err != nil {
				return fmt.Errorf("failed to create table %d: %v", table.Number, err)
			}
		}
	}

	// 14. Criar orders
	if len(data.Orders) > 0 {
		fmt.Println("  📝 Criando orders...")
		for _, order := range data.Orders {
			order.OrganizationId = orgId
			order.ProjectId = projectId
			if err := createResourceHTTP(client, "/order", order, headers); err != nil {
				return fmt.Errorf("failed to create order: %v", err)
			}
		}
	}

	// 15. Criar reservations
	if len(data.Reservations) > 0 {
		fmt.Println("  🎫 Criando reservations...")
		for _, reservation := range data.Reservations {
			reservation.OrganizationId = orgId
			reservation.ProjectId = projectId
			if err := createResourceHTTP(client, "/reservation", reservation, headers); err != nil {
				return fmt.Errorf("failed to create reservation: %v", err)
			}
		}
	}

	// 16. Criar waitlists
	if len(data.Waitlists) > 0 {
		fmt.Println("  ⏰ Criando waitlists...")
		for _, waitlist := range data.Waitlists {
			waitlist.OrganizationId = orgId
			waitlist.ProjectId = projectId
			if err := createResourceHTTP(client, "/waitlist", waitlist, headers); err != nil {
				return fmt.Errorf("failed to create waitlist: %v", err)
			}
		}
	}

	fmt.Printf("\n✅ Seeding concluído com sucesso!")
	fmt.Printf("\n📋 Organization ID: %s", orgId)
	fmt.Printf("\n📁 Project ID: %s\n", projectId)

	return nil
}

// Criar organização via HTTP (sem headers)
func createOrganizationHTTP(client *http.Client, org interface{}) (uuid.UUID, error) {
	body, _ := json.Marshal(org)

	req, _ := http.NewRequest("POST", baseURL+"/organization", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return uuid.Nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	// Se a organização já existe (409 ou 500 com duplicate key), usar o ID existente
	if resp.StatusCode == 409 || (resp.StatusCode == 500 && bytes.Contains(bodyBytes, []byte("duplicate key"))) {
		if verbose {
			fmt.Printf("    ⚠️  Organization already exists, using existing ID\n")
		}
		// Retornar o ID da organização original
		type OrgWithId struct {
			Id uuid.UUID `json:"id"`
		}
		orgBytes, _ := json.Marshal(org)
		var orgWithId OrgWithId
		json.Unmarshal(orgBytes, &orgWithId)
		return orgWithId.Id, nil
	}

	if resp.StatusCode != 201 {
		return uuid.Nil, fmt.Errorf("status %d - %s", resp.StatusCode, string(bodyBytes))
	}

	// Extrair ID da resposta
	var response map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		// Se não conseguir fazer unmarshal, usar o ID da organização passada
		if orgMap, ok := org.(map[string]interface{}); ok {
			if idStr, ok := orgMap["id"].(string); ok {
				return uuid.Parse(idStr)
			}
		}
		// Fallback: tentar converter org para struct com campo Id
		type OrgWithId struct {
			Id uuid.UUID `json:"id"`
		}
		orgBytes, _ := json.Marshal(org)
		var orgWithId OrgWithId
		json.Unmarshal(orgBytes, &orgWithId)
		return orgWithId.Id, nil
	}

	// Tentar extrair ID da resposta
	if data, ok := response["data"].(map[string]interface{}); ok {
		if idStr, ok := data["id"].(string); ok {
			return uuid.Parse(idStr)
		}
	}

	// Fallback: usar ID da organização original
	type OrgWithId struct {
		Id uuid.UUID `json:"id"`
	}
	orgBytes, _ := json.Marshal(org)
	var orgWithId OrgWithId
	json.Unmarshal(orgBytes, &orgWithId)

	if verbose {
		fmt.Printf("    ✓ Organization created: %s\n", orgWithId.Id)
	}

	return orgWithId.Id, nil
}

// Criar projeto via HTTP (apenas org header)
func createProjectHTTP(client *http.Client, project interface{}, orgId uuid.UUID) (uuid.UUID, error) {
	body, _ := json.Marshal(project)

	req, _ := http.NewRequest("POST", baseURL+"/project", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Lpe-Organization-Id", orgId.String())

	resp, err := client.Do(req)
	if err != nil {
		return uuid.Nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	// Se o projeto já existe, usar o ID existente
	if resp.StatusCode == 409 || resp.StatusCode == 500 && (bytes.Contains(bodyBytes, []byte("duplicate key")) || bytes.Contains(bodyBytes, []byte("already exists"))) {
		if verbose {
			fmt.Printf("    ⚠️  Project already exists, using existing ID\n")
		}
		type ProjectWithId struct {
			Id uuid.UUID `json:"id"`
		}
		projectBytes, _ := json.Marshal(project)
		var projectWithId ProjectWithId
		json.Unmarshal(projectBytes, &projectWithId)
		return projectWithId.Id, nil
	}

	if resp.StatusCode != 201 {
		return uuid.Nil, fmt.Errorf("status %d - %s", resp.StatusCode, string(bodyBytes))
	}

	// Extrair ID do projeto
	type ProjectWithId struct {
		Id uuid.UUID `json:"id"`
	}
	projectBytes, _ := json.Marshal(project)
	var projectWithId ProjectWithId
	json.Unmarshal(projectBytes, &projectWithId)

	if verbose {
		fmt.Printf("    ✓ Project created: %s\n", projectWithId.Id)
	}

	return projectWithId.Id, nil
}

// Criar usuário admin via HTTP e fazer login
func createAdminUserHTTP(client *http.Client, user interface{}, orgId, projectId uuid.UUID, userOrg, userProj interface{}) (string, error) {
	// 1. Criar usuário
	body, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", baseURL+"/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	// Se o usuário já existe (409 ou 500 com E-mail já cadastrado), continuar com login
	if resp.StatusCode == 409 || (resp.StatusCode == 500 && bytes.Contains(bodyBytes, []byte("E-mail já cadastrado"))) {
		if verbose {
			fmt.Printf("    ⚠️  User already exists, attempting login\n")
		}
	} else if resp.StatusCode != 201 {
		return "", fmt.Errorf("failed to create user: status %d - %s", resp.StatusCode, string(bodyBytes))
	}

	// Extrair user ID
	type UserWithId struct {
		Id    uuid.UUID `json:"id"`
		Email string    `json:"email"`
	}
	userBytes, _ := json.Marshal(user)
	var userWithId UserWithId
	json.Unmarshal(userBytes, &userWithId)

	// 2. Criar relacionamento user-organization
	type UserOrgRequest struct {
		UserId         uuid.UUID `json:"user_id"`
		OrganizationId uuid.UUID `json:"organization_id"`
		Role           string    `json:"role"`
		Permissions    []string  `json:"permissions"`
	}

	userOrgReq := UserOrgRequest{
		UserId:         userWithId.Id,
		OrganizationId: orgId,
		Role:           "owner",
		Permissions:    []string{"all"},
	}
	userOrgBody, _ := json.Marshal(userOrgReq)

	req, _ = http.NewRequest("POST", baseURL+"/user-organization/user/"+userWithId.Id.String(), bytes.NewBuffer(userOrgBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ = io.ReadAll(resp.Body)

	// Aceitar 201, 409, ou 500 com mensagem de duplicação
	if resp.StatusCode != 201 && resp.StatusCode != 409 {
		if resp.StatusCode == 500 && bytes.Contains(bodyBytes, []byte("já vinculado")) {
			if verbose {
				fmt.Printf("    ⚠️  User-organization already exists\n")
			}
		} else {
			return "", fmt.Errorf("failed to create user-organization: status %d - %s", resp.StatusCode, string(bodyBytes))
		}
	}

	// 3. Criar relacionamento user-project
	type UserProjRequest struct {
		UserId      uuid.UUID `json:"user_id"`
		ProjectId   uuid.UUID `json:"project_id"`
		Role        string    `json:"role"`
		Permissions []string  `json:"permissions"`
	}

	userProjReq := UserProjRequest{
		UserId:      userWithId.Id,
		ProjectId:   projectId,
		Role:        "admin",
		Permissions: []string{"all"},
	}
	userProjBody, _ := json.Marshal(userProjReq)

	req, _ = http.NewRequest("POST", baseURL+"/user-project/user/"+userWithId.Id.String(), bytes.NewBuffer(userProjBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ = io.ReadAll(resp.Body)

	// Aceitar 201, 409, ou 500 com mensagem de duplicação
	if resp.StatusCode != 201 && resp.StatusCode != 409 {
		if resp.StatusCode == 500 && bytes.Contains(bodyBytes, []byte("já vinculado")) {
			if verbose {
				fmt.Printf("    ⚠️  User-project already exists\n")
			}
		} else {
			return "", fmt.Errorf("failed to create user-project: status %d - %s", resp.StatusCode, string(bodyBytes))
		}
	}

	// 4. Fazer login para obter token
	loginData := map[string]string{
		"email":    userWithId.Email,
		"password": "senha123",
	}
	loginBody, _ := json.Marshal(loginData)

	req, _ = http.NewRequest("POST", baseURL+"/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	loginBodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to login: status %d - %s", resp.StatusCode, string(loginBodyBytes))
	}

	var loginResponse map[string]interface{}
	if err := json.Unmarshal(loginBodyBytes, &loginResponse); err != nil {
		return "", fmt.Errorf("failed to parse login response: %v", err)
	}

	token, ok := loginResponse["token"].(string)
	if !ok {
		return "", fmt.Errorf("no token in login response")
	}

	if verbose {
		fmt.Printf("    ✓ Admin user created and logged in: %s\n", userWithId.Email)
	}

	return token, nil
}

// Criar usuário comum via HTTP
func createUserHTTP(client *http.Client, user interface{}, orgId, projectId uuid.UUID, userOrg, userProj interface{}, headers map[string]string) error {
	// 1. Criar usuário
	body, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", baseURL+"/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 && resp.StatusCode != 409 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d - %s", resp.StatusCode, string(bodyBytes))
	}

	// Extrair user ID
	type UserWithId struct {
		Id uuid.UUID `json:"id"`
	}
	userBytes, _ := json.Marshal(user)
	var userWithId UserWithId
	json.Unmarshal(userBytes, &userWithId)

	// 2. Criar relacionamento user-organization
	type UserOrgRequest struct {
		OrganizationId uuid.UUID `json:"organization_id"`
		Role           string    `json:"role"`
		Permissions    []string  `json:"permissions"`
	}

	// Extrair role do userOrg
	type UserOrgWithRole struct {
		Role        string   `json:"role"`
		Permissions []string `json:"permissions"`
	}
	userOrgBytes, _ := json.Marshal(userOrg)
	var userOrgWithRole UserOrgWithRole
	json.Unmarshal(userOrgBytes, &userOrgWithRole)

	userOrgReq := UserOrgRequest{
		OrganizationId: orgId,
		Role:           userOrgWithRole.Role,
		Permissions:    userOrgWithRole.Permissions,
	}
	userOrgBody, _ := json.Marshal(userOrgReq)

	req, _ = http.NewRequest("POST", baseURL+"/user-organization/user/"+userWithId.Id.String(), bytes.NewBuffer(userOrgBody))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 && resp.StatusCode != 409 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create user-organization: status %d - %s", resp.StatusCode, string(bodyBytes))
	}

	// 3. Criar relacionamento user-project
	type UserProjRequest struct {
		ProjectId   uuid.UUID `json:"project_id"`
		Role        string    `json:"role"`
		Permissions []string  `json:"permissions"`
	}

	type UserProjWithRole struct {
		Role        string   `json:"role"`
		Permissions []string `json:"permissions"`
	}
	userProjBytes, _ := json.Marshal(userProj)
	var userProjWithRole UserProjWithRole
	json.Unmarshal(userProjBytes, &userProjWithRole)

	userProjReq := UserProjRequest{
		ProjectId:   projectId,
		Role:        userProjWithRole.Role,
		Permissions: userProjWithRole.Permissions,
	}
	userProjBody, _ := json.Marshal(userProjReq)

	req, _ = http.NewRequest("POST", baseURL+"/user-project/user/"+userWithId.Id.String(), bytes.NewBuffer(userProjBody))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 && resp.StatusCode != 409 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create user-project: status %d - %s", resp.StatusCode, string(bodyBytes))
	}

	if verbose {
		type UserWithEmail struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		var userWithEmail UserWithEmail
		json.Unmarshal(userBytes, &userWithEmail)
		fmt.Printf("    ✓ %s (%s)\n", userWithEmail.Name, userWithEmail.Email)
	}

	return nil
}

// Criar recurso genérico via HTTP
func createResourceHTTP(client *http.Client, endpoint string, resource interface{}, headers map[string]string) error {
	body, _ := json.Marshal(resource)

	req, _ := http.NewRequest("POST", baseURL+endpoint, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 && resp.StatusCode != 409 && resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d - %s", resp.StatusCode, string(bodyBytes))
	}

	if verbose {
		type ResourceWithName struct {
			Name string `json:"name"`
		}
		resourceBytes, _ := json.Marshal(resource)
		var resourceWithName ResourceWithName
		json.Unmarshal(resourceBytes, &resourceWithName)
		if resourceWithName.Name != "" {
			fmt.Printf("    ✓ %s\n", resourceWithName.Name)
		} else {
			fmt.Printf("    ✓ Resource created\n")
		}
	}

	return nil
}
