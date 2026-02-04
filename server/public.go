package server

import (
	"lep/handler"
	"lep/repositories/models"
	"lep/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ResourcePublic struct {
	handler *handler.Handlers
}

type IServerPublic interface {
	ServiceGetPublicMenu(c *gin.Context)
	ServiceGetProjectInfo(c *gin.Context)
	ServiceGetAvailableTimes(c *gin.Context)
	ServiceCreatePublicReservation(c *gin.Context)
	ServiceGetPublicCategories(c *gin.Context)
	ServiceGetPublicMenus(c *gin.Context)
	// Novos métodos com slugs
	ServiceResolveProject(c *gin.Context)
	ServiceGetPublicMenuBySlug(c *gin.Context)
	ServiceGetPublicCategoriesBySlug(c *gin.Context)
	ServiceGetPublicMenusBySlug(c *gin.Context)
	ServiceGetProjectInfoBySlug(c *gin.Context)
	ServiceGetAvailableTimesBySlug(c *gin.Context)
	ServiceCreatePublicReservationBySlug(c *gin.Context)
}

// ServiceGetPublicMenu retorna produtos do cardápio sem autenticação
func (r *ResourcePublic) ServiceGetPublicMenu(c *gin.Context) {
	orgIdStr := c.Param("orgId")
	projIdStr := c.Param("projId")

	// Validar UUIDs
	_, err := uuid.Parse(orgIdStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid organization ID format", err)
		return
	}

	_, err = uuid.Parse(projIdStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid project ID format", err)
		return
	}

	// Buscar produtos do cardápio disponíveis
	products, err := r.handler.HandlerProducts.ListProducts(orgIdStr, projIdStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting menu products", err)
		return
	}

	// Filtrar apenas produtos disponíveis
	var availableProducts []models.Product
	for _, product := range products {
		if product.Active {
			availableProducts = append(availableProducts, product)
		}
	}

	c.JSON(http.StatusOK, availableProducts)
}

// ServiceGetProjectInfo retorna informações básicas do projeto
func (r *ResourcePublic) ServiceGetProjectInfo(c *gin.Context) {
	orgIdStr := c.Param("orgId")
	projIdStr := c.Param("projId")

	// Validar UUIDs
	_, err := uuid.Parse(orgIdStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid organization ID format", err)
		return
	}

	_, err = uuid.Parse(projIdStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid project ID format", err)
		return
	}

	// Buscar informações do projeto
	project, err := r.handler.HandlerProject.GetProjectById(projIdStr)
	if err != nil {
		// Se o projeto não existe, retornar informações padrão
		if err.Error() == "record not found" {
			c.JSON(http.StatusOK, gin.H{
				"name":        "Restaurante",
				"description": "",
				"contact_info": gin.H{
					"phone":   "",
					"email":   "",
					"address": "",
				},
			})
			return
		}
		utils.SendInternalServerError(c, "Error getting project info", err)
		return
	}

	if project == nil {
		// Retornar informações padrão se projeto é nil
		c.JSON(http.StatusOK, gin.H{
			"name":        "Restaurante",
			"description": "",
			"contact_info": gin.H{
				"phone":   "",
				"email":   "",
				"address": "",
			},
		})
		return
	}

	// Buscar informações da organização
	organization, err := r.handler.HandlerOrganization.GetOrganizationById(orgIdStr)
	if err != nil {
		// Se a organização não existe, retornar apenas info do projeto
		c.JSON(http.StatusOK, gin.H{
			"name":        project.Name,
			"description": project.Description,
			"contact_info": gin.H{
				"phone":   "",
				"email":   "",
				"address": "",
			},
		})
		return
	}

	// Retornar informações básicas
	projectInfo := gin.H{
		"name":        project.Name,
		"description": project.Description,
		"contact_info": gin.H{
			"phone":   organization.Phone,
			"email":   organization.Email,
			"address": organization.Address,
		},
	}

	c.JSON(http.StatusOK, projectInfo)
}

// ServiceGetAvailableTimes retorna horários disponíveis para reserva
func (r *ResourcePublic) ServiceGetAvailableTimes(c *gin.Context) {
	orgIdStr := c.Param("orgId")
	projIdStr := c.Param("projId")

	// Validar UUIDs
	_, err := uuid.Parse(orgIdStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid organization ID format", err)
		return
	}

	_, err = uuid.Parse(projIdStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid project ID format", err)
		return
	}

	// Obter parâmetros da query
	dateStr := c.Query("date")
	partySizeStr := c.Query("party_size")

	if dateStr == "" || partySizeStr == "" {
		utils.SendBadRequestError(c, "Date and party_size are required", nil)
		return
	}

	partySize, err := strconv.Atoi(partySizeStr)
	if err != nil || partySize < 1 {
		utils.SendBadRequestError(c, "Invalid party_size", err)
		return
	}

	// Parse da data
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid date format. Use YYYY-MM-DD", err)
		return
	}

	// Gerar horários disponíveis (lógica simplificada)
	availableTimes := generateAvailableTimeSlots(date, partySize, orgIdStr, projIdStr, r.handler)

	c.JSON(http.StatusOK, availableTimes)
}

// ServiceCreatePublicReservation cria reserva + cliente sem autenticação
func (r *ResourcePublic) ServiceCreatePublicReservation(c *gin.Context) {
	orgIdStr := c.Param("orgId")
	projIdStr := c.Param("projId")

	// Validar UUIDs
	orgId, err := uuid.Parse(orgIdStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid organization ID format", err)
		return
	}

	projId, err := uuid.Parse(projIdStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid project ID format", err)
		return
	}

	// Estrutura da requisição
	var requestData struct {
		Customer struct {
			Name  string `json:"name" binding:"required"`
			Email string `json:"email"`
			Phone string `json:"phone" binding:"required"`
		} `json:"customer" binding:"required"`
		Reservation struct {
			Datetime  string `json:"datetime" binding:"required"`
			PartySize int    `json:"party_size" binding:"required,min=1"`
			Note      string `json:"note"`
			Source    string `json:"source"`
		} `json:"reservation" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	// Criar cliente
	newCustomer := models.Customer{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projId,
		Name:           requestData.Customer.Name,
		Email:          requestData.Customer.Email,
		Phone:          requestData.Customer.Phone,
	}

	err = r.handler.HandlerCustomer.CreateCustomer(&newCustomer)
	if err != nil {
		utils.SendInternalServerError(c, "Error creating customer", err)
		return
	}

	// Parse datetime
	datetime, err := time.Parse(time.RFC3339, requestData.Reservation.Datetime)
	if err != nil {
		// Tentar formato alternativo
		datetime, err = time.Parse("2006-01-02T15:04", requestData.Reservation.Datetime)
		if err != nil {
			utils.SendBadRequestError(c, "Invalid datetime format", err)
			return
		}
	}

	// Buscar mesa disponível (simplificado - pegar primeira mesa com capacidade)
	tables, err := r.handler.HandlerTables.ListTables(orgIdStr, projIdStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error finding available tables", err)
		return
	}

	var selectedTable *models.Table
	for _, table := range tables {
		if table.Capacity >= requestData.Reservation.PartySize && table.Status == "livre" {
			selectedTable = &table
			break
		}
	}

	if selectedTable == nil {
		utils.SendBadRequestError(c, "No available tables for this party size", nil)
		return
	}

	// Criar reserva
	newReservation := models.Reservation{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projId,
		CustomerId:     newCustomer.Id,
		TableId:        selectedTable.Id,
		Datetime:       datetime.Format(time.RFC3339),
		PartySize:      requestData.Reservation.PartySize,
		Status:         "confirmed",
		Note:           requestData.Reservation.Note,
	}

	err = r.handler.HandlerReservation.CreateReservation(&newReservation)
	if err != nil {
		utils.SendInternalServerError(c, "Error creating reservation", err)
		return
	}

	// Retornar resposta com dados criados
	response := gin.H{
		"customer":    newCustomer,
		"reservation": newReservation,
		"table":       selectedTable,
	}

	utils.SendCreatedSuccess(c, "Reservation created successfully", response)
}

// generateAvailableTimeSlots gera horários disponíveis (lógica simplificada)
func generateAvailableTimeSlots(date time.Time, partySize int, orgId, projId string, handler *handler.Handlers) []gin.H {
	// Horários padrão de funcionamento
	timeSlots := []string{
		"12:00", "12:30", "13:00", "13:30", "14:00", "14:30",
		"19:00", "19:30", "20:00", "20:30", "21:00", "21:30", "22:00",
	}

	var availableTimes []gin.H

	// Para cada horário, verificar disponibilidade (lógica simplificada)
	for _, timeSlot := range timeSlots {
		// Por simplicidade, marcar todos como disponíveis
		// Em produção, verificar conflitos com reservas existentes
		availableTimes = append(availableTimes, gin.H{
			"time":      timeSlot,
			"available": true,
		})
	}

	return availableTimes
}

// ServiceGetPublicCategories retorna categorias ativas sem autenticação
func (r *ResourcePublic) ServiceGetPublicCategories(c *gin.Context) {
	orgIdStr := c.Param("orgId")
	projIdStr := c.Param("projId")

	// Validar UUIDs
	_, err := uuid.Parse(orgIdStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid organization ID format", err)
		return
	}

	_, err = uuid.Parse(projIdStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid project ID format", err)
		return
	}

	// Buscar categorias ativas
	categories, err := r.handler.HandlerCategory.ListCategories(orgIdStr, projIdStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting categories", err)
		return
	}

	// Filtrar apenas categorias ativas
	var activeCategories []models.Category
	for _, category := range categories {
		if category.Active {
			activeCategories = append(activeCategories, category)
		}
	}

	c.JSON(http.StatusOK, activeCategories)
}

// ServiceGetPublicMenus retorna menus ativos sem autenticação
func (r *ResourcePublic) ServiceGetPublicMenus(c *gin.Context) {
	orgIdStr := c.Param("orgId")
	projIdStr := c.Param("projId")

	// Validar UUIDs
	_, err := uuid.Parse(orgIdStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid organization ID format", err)
		return
	}

	_, err = uuid.Parse(projIdStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid project ID format", err)
		return
	}

	// Buscar menus ativos
	menus, err := r.handler.HandlerMenu.ListMenus(orgIdStr, projIdStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting menus", err)
		return
	}

	// Filtrar apenas menus ativos
	var activeMenus []models.Menu
	for _, menu := range menus {
		if menu.Active {
			activeMenus = append(activeMenus, menu)
		}
	}

	c.JSON(http.StatusOK, activeMenus)
}

// ServiceResolveProject resolve organização e projeto por slugs
func (r *ResourcePublic) ServiceResolveProject(c *gin.Context) {
	orgSlug := c.Query("org_slug")
	projectSlug := c.Query("project_slug") // opcional

	if orgSlug == "" {
		utils.SendBadRequestError(c, "org_slug is required", nil)
		return
	}

	// Buscar organização por slug
	org, err := r.handler.HandlerOrganization.GetOrganizationBySlug(orgSlug)
	if err != nil {
		utils.SendNotFoundError(c, "Organization not found")
		return
	}

	// Resolver projeto (por slug ou default/primeiro ativo)
	project, err := r.handler.HandlerProject.ResolveProject(org.Id.String(), projectSlug)
	if err != nil {
		utils.SendNotFoundError(c, "Project not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"organization_id":   org.Id,
		"organization_slug": org.Slug,
		"organization_name": org.Name,
		"project_id":        project.Id,
		"project_slug":      project.Slug,
		"project_name":      project.Name,
		"is_default":        project.IsDefault,
	})
}

// resolveOrgAndProject é um helper interno para resolver org e projeto por slugs
func (r *ResourcePublic) resolveOrgAndProject(orgSlug, projectSlug string) (string, string, error) {
	// Buscar organização por slug
	org, err := r.handler.HandlerOrganization.GetOrganizationBySlug(orgSlug)
	if err != nil {
		return "", "", err
	}

	// Resolver projeto
	project, err := r.handler.HandlerProject.ResolveProject(org.Id.String(), projectSlug)
	if err != nil {
		return "", "", err
	}

	return org.Id.String(), project.Id.String(), nil
}

// ServiceGetPublicMenuBySlug retorna produtos do cardápio usando slugs
func (r *ResourcePublic) ServiceGetPublicMenuBySlug(c *gin.Context) {
	orgSlug := c.Param("orgSlug")
	projectSlug := c.Param("projectSlug") // pode estar vazio

	orgId, projId, err := r.resolveOrgAndProject(orgSlug, projectSlug)
	if err != nil {
		utils.SendNotFoundError(c, "Organization or project not found")
		return
	}

	// Buscar produtos do cardápio disponíveis
	products, err := r.handler.HandlerProducts.ListProducts(orgId, projId)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting menu products", err)
		return
	}

	// Filtrar apenas produtos disponíveis
	var availableProducts []models.Product
	for _, product := range products {
		if product.Active {
			availableProducts = append(availableProducts, product)
		}
	}

	c.JSON(http.StatusOK, availableProducts)
}

// ServiceGetPublicCategoriesBySlug retorna categorias usando slugs
func (r *ResourcePublic) ServiceGetPublicCategoriesBySlug(c *gin.Context) {
	orgSlug := c.Param("orgSlug")
	projectSlug := c.Param("projectSlug")

	orgId, projId, err := r.resolveOrgAndProject(orgSlug, projectSlug)
	if err != nil {
		utils.SendNotFoundError(c, "Organization or project not found")
		return
	}

	categories, err := r.handler.HandlerCategory.ListCategories(orgId, projId)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting categories", err)
		return
	}

	var activeCategories []models.Category
	for _, category := range categories {
		if category.Active {
			activeCategories = append(activeCategories, category)
		}
	}

	c.JSON(http.StatusOK, activeCategories)
}

// ServiceGetPublicMenusBySlug retorna menus usando slugs
func (r *ResourcePublic) ServiceGetPublicMenusBySlug(c *gin.Context) {
	orgSlug := c.Param("orgSlug")
	projectSlug := c.Param("projectSlug")

	orgId, projId, err := r.resolveOrgAndProject(orgSlug, projectSlug)
	if err != nil {
		utils.SendNotFoundError(c, "Organization or project not found")
		return
	}

	menus, err := r.handler.HandlerMenu.ListMenus(orgId, projId)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting menus", err)
		return
	}

	var activeMenus []models.Menu
	for _, menu := range menus {
		if menu.Active {
			activeMenus = append(activeMenus, menu)
		}
	}

	c.JSON(http.StatusOK, activeMenus)
}

// ServiceGetProjectInfoBySlug retorna informações do projeto usando slugs
func (r *ResourcePublic) ServiceGetProjectInfoBySlug(c *gin.Context) {
	orgSlug := c.Param("orgSlug")
	projectSlug := c.Param("projectSlug")

	// Buscar organização por slug
	org, err := r.handler.HandlerOrganization.GetOrganizationBySlug(orgSlug)
	if err != nil {
		utils.SendNotFoundError(c, "Organization not found")
		return
	}

	// Resolver projeto
	project, err := r.handler.HandlerProject.ResolveProject(org.Id.String(), projectSlug)
	if err != nil {
		utils.SendNotFoundError(c, "Project not found")
		return
	}

	projectInfo := gin.H{
		"name":        project.Name,
		"description": project.Description,
		"slug":        project.Slug,
		"is_default":  project.IsDefault,
		"contact_info": gin.H{
			"phone":   org.Phone,
			"email":   org.Email,
			"address": org.Address,
		},
	}

	c.JSON(http.StatusOK, projectInfo)
}

// ServiceGetAvailableTimesBySlug retorna horários disponíveis usando slugs
func (r *ResourcePublic) ServiceGetAvailableTimesBySlug(c *gin.Context) {
	orgSlug := c.Param("orgSlug")
	projectSlug := c.Param("projectSlug")

	orgId, projId, err := r.resolveOrgAndProject(orgSlug, projectSlug)
	if err != nil {
		utils.SendNotFoundError(c, "Organization or project not found")
		return
	}

	dateStr := c.Query("date")
	partySizeStr := c.Query("party_size")

	if dateStr == "" || partySizeStr == "" {
		utils.SendBadRequestError(c, "Date and party_size are required", nil)
		return
	}

	partySize, err := strconv.Atoi(partySizeStr)
	if err != nil || partySize < 1 {
		utils.SendBadRequestError(c, "Invalid party_size", err)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		utils.SendBadRequestError(c, "Invalid date format. Use YYYY-MM-DD", err)
		return
	}

	availableTimes := generateAvailableTimeSlots(date, partySize, orgId, projId, r.handler)
	c.JSON(http.StatusOK, availableTimes)
}

// ServiceCreatePublicReservationBySlug cria reserva usando slugs
func (r *ResourcePublic) ServiceCreatePublicReservationBySlug(c *gin.Context) {
	orgSlug := c.Param("orgSlug")
	projectSlug := c.Param("projectSlug")

	// Buscar organização por slug
	org, err := r.handler.HandlerOrganization.GetOrganizationBySlug(orgSlug)
	if err != nil {
		utils.SendNotFoundError(c, "Organization not found")
		return
	}

	// Resolver projeto
	project, err := r.handler.HandlerProject.ResolveProject(org.Id.String(), projectSlug)
	if err != nil {
		utils.SendNotFoundError(c, "Project not found")
		return
	}

	orgId := org.Id
	projId := project.Id

	var requestData struct {
		Customer struct {
			Name  string `json:"name" binding:"required"`
			Email string `json:"email"`
			Phone string `json:"phone" binding:"required"`
		} `json:"customer" binding:"required"`
		Reservation struct {
			Datetime  string `json:"datetime" binding:"required"`
			PartySize int    `json:"party_size" binding:"required,min=1"`
			Note      string `json:"note"`
			Source    string `json:"source"`
		} `json:"reservation" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		utils.SendBadRequestError(c, "Invalid request body", err)
		return
	}

	newCustomer := models.Customer{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projId,
		Name:           requestData.Customer.Name,
		Email:          requestData.Customer.Email,
		Phone:          requestData.Customer.Phone,
	}

	err = r.handler.HandlerCustomer.CreateCustomer(&newCustomer)
	if err != nil {
		utils.SendInternalServerError(c, "Error creating customer", err)
		return
	}

	datetime, err := time.Parse(time.RFC3339, requestData.Reservation.Datetime)
	if err != nil {
		datetime, err = time.Parse("2006-01-02T15:04", requestData.Reservation.Datetime)
		if err != nil {
			utils.SendBadRequestError(c, "Invalid datetime format", err)
			return
		}
	}

	tables, err := r.handler.HandlerTables.ListTables(orgId.String(), projId.String())
	if err != nil {
		utils.SendInternalServerError(c, "Error finding available tables", err)
		return
	}

	var selectedTable *models.Table
	for _, table := range tables {
		if table.Capacity >= requestData.Reservation.PartySize && table.Status == "livre" {
			selectedTable = &table
			break
		}
	}

	if selectedTable == nil {
		utils.SendBadRequestError(c, "No available tables for this party size", nil)
		return
	}

	newReservation := models.Reservation{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projId,
		CustomerId:     newCustomer.Id,
		TableId:        selectedTable.Id,
		Datetime:       datetime.Format(time.RFC3339),
		PartySize:      requestData.Reservation.PartySize,
		Status:         "confirmed",
		Note:           requestData.Reservation.Note,
	}

	err = r.handler.HandlerReservation.CreateReservation(&newReservation)
	if err != nil {
		utils.SendInternalServerError(c, "Error creating reservation", err)
		return
	}

	response := gin.H{
		"customer":    newCustomer,
		"reservation": newReservation,
		"table":       selectedTable,
	}

	utils.SendCreatedSuccess(c, "Reservation created successfully", response)
}

func NewSourceServerPublic(handler *handler.Handlers) IServerPublic {
	return &ResourcePublic{handler: handler}
}
