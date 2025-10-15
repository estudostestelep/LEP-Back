package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// createOrganizationBootstrap cria organização via bootstrap e retorna org ID, project ID e email do admin
func createOrganizationBootstrap(router *gin.Engine, orgName string) (uuid.UUID, uuid.UUID, string, error) {
	bootstrapRequest := map[string]string{
		"name":     orgName,
		"password": "senha123",
	}
	body, _ := json.Marshal(bootstrapRequest)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/create-organization", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != 201 {
		return uuid.Nil, uuid.Nil, "", fmt.Errorf("failed to create organization bootstrap: status %d - %s", w.Code, w.Body.String())
	}

	// Parse response
	var response struct {
		Data struct {
			Organization struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"organization"`
			Project struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"project"`
			User struct {
				ID    string `json:"id"`
				Email string `json:"email"`
				Name  string `json:"name"`
			} `json:"user"`
		} `json:"data"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		return uuid.Nil, uuid.Nil, "", fmt.Errorf("failed to parse bootstrap response: %v", err)
	}

	orgId, err := uuid.Parse(response.Data.Organization.ID)
	if err != nil {
		return uuid.Nil, uuid.Nil, "", fmt.Errorf("failed to parse organization ID: %v", err)
	}

	projectId, err := uuid.Parse(response.Data.Project.ID)
	if err != nil {
		return uuid.Nil, uuid.Nil, "", fmt.Errorf("failed to parse project ID: %v", err)
	}

	if verbose {
		fmt.Printf("    ✓ Organização: %s\n", response.Data.Organization.Name)
		fmt.Printf("    ✓ Projeto: %s\n", response.Data.Project.Name)
		fmt.Printf("    ✓ Admin: %s (%s)\n", response.Data.User.Name, response.Data.User.Email)
	}

	return orgId, projectId, response.Data.User.Email, nil
}

// loginUser faz login e retorna token JWT
func loginUser(router *gin.Engine, email, password string) (string, error) {
	loginData := map[string]string{
		"email":    email,
		"password": password,
	}
	body, _ := json.Marshal(loginData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		return "", fmt.Errorf("failed to login: status %d - %s", w.Code, w.Body.String())
	}

	var loginResponse map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &loginResponse); err != nil {
		return "", fmt.Errorf("failed to parse login response: %v", err)
	}

	token, ok := loginResponse["token"].(string)
	if !ok {
		return "", fmt.Errorf("no token in login response")
	}

	if verbose {
		fmt.Printf("    ✓ Login realizado com sucesso\n")
	}

	return token, nil
}
