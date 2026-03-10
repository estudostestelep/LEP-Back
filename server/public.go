package server

import (
	"encoding/json"
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
	// Fila de espera pública
	ServiceGetPublicWaitlist(c *gin.Context)
	ServiceGetPublicWaitlistBySlug(c *gin.Context)
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

	// Buscar ou criar cliente
	var customer models.Customer
	if requestData.Customer.Email != "" {
		existingCustomer, emailErr := r.handler.HandlerCustomer.GetCustomerByEmail(orgId, projId, requestData.Customer.Email)
		if emailErr == nil && existingCustomer != nil {
			customer = *existingCustomer
		} else {
			newCustomer := models.Customer{
				OrganizationId: orgId,
				ProjectId:      projId,
				Name:           requestData.Customer.Name,
				Email:          requestData.Customer.Email,
				Phone:          requestData.Customer.Phone,
			}
			if createErr := r.handler.HandlerCustomer.CreateCustomer(&newCustomer); createErr != nil {
				utils.SendInternalServerError(c, "Error creating customer", createErr)
				return
			}
			customer = newCustomer
		}
	} else {
		newCustomer := models.Customer{
			OrganizationId: orgId,
			ProjectId:      projId,
			Name:           requestData.Customer.Name,
			Phone:          requestData.Customer.Phone,
		}
		if createErr := r.handler.HandlerCustomer.CreateCustomer(&newCustomer); createErr != nil {
			utils.SendInternalServerError(c, "Error creating customer", createErr)
			return
		}
		customer = newCustomer
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

	// Carregar configurações do projeto (usadas para diningDuration e AutoConfirm)
	settings, _ := r.handler.HandlerSettings.GetOrCreateSettings(orgIdStr, projIdStr)
	diningDuration := 120
	if settings != nil && settings.DiningDurationMinutes > 0 {
		diningDuration = settings.DiningDurationMinutes
	}

	// Buscar mesa disponível (pegar primeira mesa com capacidade suficiente e sem conflito de reserva)
	tables, err := r.handler.HandlerTables.ListTables(orgIdStr, projIdStr, nil)
	if err != nil {
		utils.SendInternalServerError(c, "Error finding available tables", err)
		return
	}

	// isPendingBySize: excede threshold configurado OU nenhuma mesa comporta o grupo sozinha
	exceedsThreshold := settings != nil && settings.AutoConfirmMaxPartySize > 0 && requestData.Reservation.PartySize > settings.AutoConfirmMaxPartySize
	noTableFits := true
	for _, t := range tables {
		if t.Capacity >= requestData.Reservation.PartySize {
			noTableFits = false
			break
		}
	}
	isPendingBySize := exceedsThreshold || noTableFits

	var selectedTable *models.Table
	for _, table := range tables {
		if table.Capacity < requestData.Reservation.PartySize {
			continue
		}
		available, availErr := r.handler.HandlerReservation.IsTableAvailable(table.Id, datetime, diningDuration)
		if availErr != nil || !available {
			continue
		}
		t := table
		selectedTable = &t
		break
	}

	// Para grupos que exigem intervenção manual, tenta qualquer mesa disponível como placeholder
	if selectedTable == nil && isPendingBySize {
		for _, table := range tables {
			available, availErr := r.handler.HandlerReservation.IsTableAvailable(table.Id, datetime, diningDuration)
			if availErr != nil || !available {
				continue
			}
			t := table
			selectedTable = &t
			break
		}
	}

	if selectedTable == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "no_availability",
			"message": "Estamos sem disponibilidade para o horário e quantidade de pessoas selecionados. Por favor, escolha outro horário ou entre em contato conosco.",
		})
		return
	}

	// Determinar status
	reservationStatus := "confirmed"
	if isPendingBySize {
		reservationStatus = "pending"
	}

	// Criar reserva
	var tableId *uuid.UUID
	if selectedTable != nil {
		tableId = &selectedTable.Id
	}
	newReservation := models.Reservation{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projId,
		CustomerId:     customer.Id,
		TableId:        tableId,
		Datetime:       datetime.Format(time.RFC3339),
		PartySize:      requestData.Reservation.PartySize,
		Status:         reservationStatus,
		Note:           requestData.Reservation.Note,
	}

	err = r.handler.HandlerReservation.CreateReservation(&newReservation)
	if err != nil {
		utils.SendInternalServerError(c, "Error creating reservation", err)
		return
	}

	// Disparar notificação de status pendente
	if reservationStatus == "pending" && r.handler.EventService != nil {
		r.handler.EventService.TriggerReservationStatusChanged(orgId, projId, &newReservation, &customer, selectedTable)
	}

	// Retornar resposta com dados criados
	response := gin.H{
		"customer":    customer,
		"reservation": newReservation,
		"table":       selectedTable,
	}

	utils.SendCreatedSuccess(c, "Reservation created successfully", response)
}

// generateAvailableTimeSlots gera horários disponíveis verificando disponibilidade real no banco.
// Os horários de funcionamento são carregados das configurações do projeto.
func generateAvailableTimeSlots(date time.Time, partySize int, orgId, projId string, h *handler.Handlers) []gin.H {
	// Defaults caso as configurações não sejam encontradas
	lunchStart, lunchEnd := "12:00", "14:30"
	dinnerStart, dinnerEnd := "19:00", "22:00"
	slotInterval := 30
	enableLunch, enableDinner := true, true
	diningDuration := 120

	settings, err := h.HandlerSettings.GetOrCreateSettings(orgId, projId)
	if err == nil && settings != nil {
		if settings.LunchStart != "" {
			lunchStart = settings.LunchStart
		}
		if settings.LunchEnd != "" {
			lunchEnd = settings.LunchEnd
		}
		if settings.DinnerStart != "" {
			dinnerStart = settings.DinnerStart
		}
		if settings.DinnerEnd != "" {
			dinnerEnd = settings.DinnerEnd
		}
		if settings.SlotIntervalMinutes > 0 {
			slotInterval = settings.SlotIntervalMinutes
		}
		enableLunch = settings.EnableLunch
		enableDinner = settings.EnableDinner
		if settings.DiningDurationMinutes > 0 {
			diningDuration = settings.DiningDurationMinutes
		}

		// Agenda semanal: sobrescreve enableLunch/enableDinner para o dia da semana solicitado
		if settings.OperatingScheduleJson != "" {
			type dayConfig struct {
				Enabled      bool `json:"enabled"`
				EnableLunch  bool `json:"enable_lunch"`
				EnableDinner bool `json:"enable_dinner"`
			}
			var schedule map[string]dayConfig
			if jsonErr := json.Unmarshal([]byte(settings.OperatingScheduleJson), &schedule); jsonErr == nil {
				dayKey := strconv.Itoa(int(date.Weekday())) // 0=Domingo ... 6=Sábado
				if dc, ok := schedule[dayKey]; ok {
					if !dc.Enabled {
						return []gin.H{} // restaurante fechado neste dia
					}
					enableLunch = dc.EnableLunch
					enableDinner = dc.EnableDinner
				}
			}
		}
	}

	var timeSlots []string
	if enableLunch {
		timeSlots = append(timeSlots, buildTimeSlots(lunchStart, lunchEnd, slotInterval)...)
	}
	if enableDinner {
		timeSlots = append(timeSlots, buildTimeSlots(dinnerStart, dinnerEnd, slotInterval)...)
	}

	tables, tablesErr := h.HandlerTables.ListTables(orgId, projId, nil)

	// Se nenhuma mesa tem capacidade para o grupo, aceita qualquer mesa disponível (reserva ficará pending)
	noTableFitsParty := true
	if tablesErr == nil {
		for _, t := range tables {
			if t.Capacity >= partySize {
				noTableFitsParty = false
				break
			}
		}
	}

	availableTimes := make([]gin.H, 0)
	for _, slot := range timeSlots {
		slotTime, parseErr := time.Parse("15:04", slot)
		if parseErr != nil {
			continue
		}
		dt := time.Date(date.Year(), date.Month(), date.Day(),
			slotTime.Hour(), slotTime.Minute(), 0, 0, time.UTC)

		hasAvailableTable := false
		if tablesErr == nil {
			for _, table := range tables {
				// Se nenhuma mesa comporta o grupo, verifica disponibilidade sem filtrar capacidade
				if !noTableFitsParty && table.Capacity < partySize {
					continue
				}
				ok, checkErr := h.HandlerReservation.IsTableAvailable(table.Id, dt, diningDuration)
				if checkErr == nil && ok {
					hasAvailableTable = true
					break
				}
			}
		}

		availableTimes = append(availableTimes, gin.H{
			"time":      slot,
			"available": hasAvailableTable,
		})
	}

	return availableTimes
}

// buildTimeSlots gera uma lista de horários (HH:MM) entre start e end com o intervalo dado em minutos.
func buildTimeSlots(start, end string, intervalMinutes int) []string {
	startTime, err := time.Parse("15:04", start)
	if err != nil {
		return nil
	}
	endTime, err := time.Parse("15:04", end)
	if err != nil {
		return nil
	}
	var slots []string
	cur := startTime
	for !cur.After(endTime) {
		slots = append(slots, cur.Format("15:04"))
		cur = cur.Add(time.Duration(intervalMinutes) * time.Minute)
	}
	return slots
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

	// Buscar ou criar cliente
	var customer models.Customer
	if requestData.Customer.Email != "" {
		existingCustomer, emailErr := r.handler.HandlerCustomer.GetCustomerByEmail(orgId, projId, requestData.Customer.Email)
		if emailErr == nil && existingCustomer != nil {
			customer = *existingCustomer
		} else {
			newCustomer := models.Customer{
				OrganizationId: orgId,
				ProjectId:      projId,
				Name:           requestData.Customer.Name,
				Email:          requestData.Customer.Email,
				Phone:          requestData.Customer.Phone,
			}
			if createErr := r.handler.HandlerCustomer.CreateCustomer(&newCustomer); createErr != nil {
				utils.SendInternalServerError(c, "Error creating customer", createErr)
				return
			}
			customer = newCustomer
		}
	} else {
		newCustomer := models.Customer{
			OrganizationId: orgId,
			ProjectId:      projId,
			Name:           requestData.Customer.Name,
			Phone:          requestData.Customer.Phone,
		}
		if createErr := r.handler.HandlerCustomer.CreateCustomer(&newCustomer); createErr != nil {
			utils.SendInternalServerError(c, "Error creating customer", createErr)
			return
		}
		customer = newCustomer
	}

	datetime, err := time.Parse(time.RFC3339, requestData.Reservation.Datetime)
	if err != nil {
		datetime, err = time.Parse("2006-01-02T15:04", requestData.Reservation.Datetime)
		if err != nil {
			utils.SendBadRequestError(c, "Invalid datetime format", err)
			return
		}
	}

	// Carregar configurações do projeto (usadas para diningDuration e AutoConfirm)
	settings, _ := r.handler.HandlerSettings.GetOrCreateSettings(orgId.String(), projId.String())
	diningDuration := 120
	if settings != nil && settings.DiningDurationMinutes > 0 {
		diningDuration = settings.DiningDurationMinutes
	}

	tables, err := r.handler.HandlerTables.ListTables(orgId.String(), projId.String(), nil)
	if err != nil {
		utils.SendInternalServerError(c, "Error finding available tables", err)
		return
	}

	// isPendingBySize: excede threshold configurado OU nenhuma mesa comporta o grupo sozinha
	exceedsThresholdSlug := settings != nil && settings.AutoConfirmMaxPartySize > 0 && requestData.Reservation.PartySize > settings.AutoConfirmMaxPartySize
	noTableFitsSlug := true
	for _, t := range tables {
		if t.Capacity >= requestData.Reservation.PartySize {
			noTableFitsSlug = false
			break
		}
	}
	isPendingBySize := exceedsThresholdSlug || noTableFitsSlug

	var selectedTable *models.Table
	for _, table := range tables {
		if table.Capacity < requestData.Reservation.PartySize {
			continue
		}
		available, availErr := r.handler.HandlerReservation.IsTableAvailable(table.Id, datetime, diningDuration)
		if availErr != nil || !available {
			continue
		}
		t := table
		selectedTable = &t
		break
	}

	// Para grupos que exigem intervenção manual, tenta qualquer mesa disponível como placeholder
	if selectedTable == nil && isPendingBySize {
		for _, table := range tables {
			available, availErr := r.handler.HandlerReservation.IsTableAvailable(table.Id, datetime, diningDuration)
			if availErr != nil || !available {
				continue
			}
			t := table
			selectedTable = &t
			break
		}
	}

	if selectedTable == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "no_availability",
			"message": "Estamos sem disponibilidade para o horário e quantidade de pessoas selecionados. Por favor, escolha outro horário ou entre em contato conosco.",
		})
		return
	}

	// Determinar status
	reservationStatus := "confirmed"
	if isPendingBySize {
		reservationStatus = "pending"
	}

	var tableIdSlug *uuid.UUID
	if selectedTable != nil {
		tableIdSlug = &selectedTable.Id
	}
	newReservation := models.Reservation{
		Id:             uuid.New(),
		OrganizationId: orgId,
		ProjectId:      projId,
		CustomerId:     customer.Id,
		TableId:        tableIdSlug,
		Datetime:       datetime.Format(time.RFC3339),
		PartySize:      requestData.Reservation.PartySize,
		Status:         reservationStatus,
		Note:           requestData.Reservation.Note,
	}

	err = r.handler.HandlerReservation.CreateReservation(&newReservation)
	if err != nil {
		utils.SendInternalServerError(c, "Error creating reservation", err)
		return
	}

	// Disparar notificação de status pendente
	if reservationStatus == "pending" && r.handler.EventService != nil {
		r.handler.EventService.TriggerReservationStatusChanged(orgId, projId, &newReservation, &customer, selectedTable)
	}

	response := gin.H{
		"customer":    customer,
		"reservation": newReservation,
		"table":       selectedTable,
	}

	utils.SendCreatedSuccess(c, "Reservation created successfully", response)
}

// PublicWaitlistEntry representa uma entrada na fila pública (sem dados sensíveis)
type PublicWaitlistEntry struct {
	Position      int    `json:"position"`
	CustomerName  string `json:"customer_name"`
	PartySize     int    `json:"party_size"`
	WaitedMinutes int    `json:"waited_minutes"`
}

// ServiceGetPublicWaitlist retorna a fila de espera sem autenticação (por UUID)
func (r *ResourcePublic) ServiceGetPublicWaitlist(c *gin.Context) {
	orgIdStr := c.Param("orgId")
	projIdStr := c.Param("projId")

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

	waitlist, err := r.handler.HandlerWaitlist.ListWaitlists(orgIdStr, projIdStr)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting waitlist", err)
		return
	}

	var queue []PublicWaitlistEntry
	position := 1
	now := time.Now()
	for _, entry := range waitlist {
		if entry.Status == "waiting" {
			queue = append(queue, PublicWaitlistEntry{
				Position:      position,
				CustomerName:  entry.CustomerName,
				PartySize:     entry.People,
				WaitedMinutes: int(now.Sub(entry.CreatedAt).Minutes()),
			})
			position++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"queue":         queue,
		"total_waiting": len(queue),
	})
}

// ServiceGetPublicWaitlistBySlug retorna a fila de espera sem autenticação (por slug)
func (r *ResourcePublic) ServiceGetPublicWaitlistBySlug(c *gin.Context) {
	orgSlug := c.Param("orgSlug")
	projectSlug := c.Param("projectSlug")

	orgId, projId, err := r.resolveOrgAndProject(orgSlug, projectSlug)
	if err != nil {
		utils.SendNotFoundError(c, "Organization or project not found")
		return
	}

	waitlist, err := r.handler.HandlerWaitlist.ListWaitlists(orgId, projId)
	if err != nil {
		utils.SendInternalServerError(c, "Error getting waitlist", err)
		return
	}

	var queue []PublicWaitlistEntry
	position := 1
	now := time.Now()
	for _, entry := range waitlist {
		if entry.Status == "waiting" {
			queue = append(queue, PublicWaitlistEntry{
				Position:      position,
				CustomerName:  entry.CustomerName,
				PartySize:     entry.People,
				WaitedMinutes: int(now.Sub(entry.CreatedAt).Minutes()),
			})
			position++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"queue":         queue,
		"total_waiting": len(queue),
	})
}

func NewSourceServerPublic(handler *handler.Handlers) IServerPublic {
	return &ResourcePublic{handler: handler}
}
