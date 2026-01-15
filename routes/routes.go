package routes

import (
	"lep/handler"
	"lep/middleware"
	"lep/resource"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Public routes (no authentication required)
	r.POST("/login", resource.ServersControllers.SourceAuth.ServiceLogin)
	r.POST("/user", resource.ServersControllers.SourceUsers.ServiceCreateUser)
	r.POST("/create-organization", resource.ServersControllers.SourceOrganization.ServiceCreateOrganizationBootstrap)
	r.POST("/organization", resource.ServersControllers.SourceOrganization.CreateOrganization)                                 // For seeding
	r.POST("/project", resource.ServersControllers.SourceProject.CreateProject)                                                // For seeding with org header
	r.POST("/user-organization/user/:userId", resource.ServersControllers.SourceUserOrganization.ServiceAddUserToOrganization) // For seeding
	r.POST("/user-project/user/:userId", resource.ServersControllers.SourceUserProject.ServiceAddUserToProject)                // For seeding

	// Admin routes (temporary - for password reset)
	r.POST("/admin/reset-passwords", resource.ServersControllers.SourceAdmin.ServiceResetPasswords)

	// =========================================================================
	// DEV MIGRATION ENDPOINT
	// Este endpoint é exclusivo para desenvolvedores (pablo@lep.com)
	// Permite executar migrações de banco de dados através de um botão no frontend
	// Para adicionar novas migrações, edite server/admin.go -> ServiceRunDevMigration
	// =========================================================================
	r.POST("/admin/run-migration", resource.ServersControllers.SourceAdmin.ServiceRunDevMigration)

	// Public routes for menu and reservations (no authentication)
	setupPublicRoutes(r)

	// Upload routes (require authentication)
	setupUploadRoutes(r)

	// Create protected route group with authentication middlewares
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.HeaderValidationMiddleware())

	// Protected auth routes
	protected.POST("/logout", resource.ServersControllers.SourceAuth.ServiceLogout)
	protected.POST("/checkToken", resource.ServersControllers.SourceAuth.ServiceValidateToken)

	// Protected routes (require authentication and organization/project headers)
	// IMPORTANT: Specific routes MUST be registered before generic routes with path parameters
	// Register all specific settings routes first to avoid conflicts with /project/:id
	setupSettingsRoutes(protected)
	setupDisplaySettingsRoutes(protected)
	setupThemeCustomizationRoutes(protected)
	setupSidebarConfigRoutes(protected)

	// Then register generic project routes
	setupOrganizationRoutes(protected)
	setupUserRoutes(protected)
	setupUserOrganizationRoutes(protected)
	setupUserProjectRoutes(protected)
	setupProductRoutes(protected)
	setupTableRoutes(protected)
	setupWaitlistRoutes(protected)
	setupReservationRoutes(protected)
	setupCustomerRoutes(protected)
	setupOrderRoutes(protected)
	setupProjectRoutes(protected)
	setupEnvironmentRoutes(protected)
	setupReportsRoutes(protected)
	setupTagRoutes(protected)
	setupMenuRoutes(protected)
	setupCategoryRoutes(protected)
	setupSubcategoryRoutes(protected)
	setupImageManagementRoutes(protected)
	setupOnboardingRoutes(protected)
	setupRoleRoutes(protected)
	setupModuleRoutes(protected)
	setupPermissionRoutes(protected)
	setupPackageRoutes(protected)
	setupPlanChangeRequestRoutes(protected)

	// Notification routes (mixed public/protected)
	setupNotificationRoutes(r)
}

// setupOnboardingRoutes configura rotas para status de onboarding
func setupOnboardingRoutes(r gin.IRouter) {
	onboardingRoutes := r.Group("/onboarding")
	{
		onboardingRoutes.GET("/status", resource.ServersControllers.SourceOnboarding.GetOnboardingStatus)
	}
}

func setupUserRoutes(r gin.IRouter) {
	userRoutes := r.Group("/user")
	{
		userRoutes.GET("/:id", resource.ServersControllers.SourceUsers.ServiceGetUser)
		userRoutes.GET("/group/:id", resource.ServersControllers.SourceUsers.ServiceGetUserByGroup)
		userRoutes.GET("", resource.ServersControllers.SourceUsers.ServiceListUsers) // Endpoint de listagem
		userRoutes.PUT("/:id", resource.ServersControllers.SourceUsers.ServiceUpdateUser)
		userRoutes.DELETE("/:id", resource.ServersControllers.SourceUsers.ServiceDeleteUser)

		// User Access Management (Master Admin only)
		userRoutes.GET("/:id/organizations-projects",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceUserAccess.ServiceGetUserOrganizationsAndProjects)
		userRoutes.POST("/:id/organizations-projects",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceUserAccess.ServiceUpdateUserOrganizationsAndProjects)
	}
}

func setupUserOrganizationRoutes(r gin.IRouter) {
	// Rotas para gerenciar relacionamento usuário-organização (todas protegidas)
	// Nota: POST /user-organization/user/:userId está registrado como rota pública para seeding
	userOrgRoutes := r.Group("/user-organization")
	{
		// POST removed - registered as public route for seeding
		userOrgRoutes.DELETE("/user/:userId/org/:orgId", resource.ServersControllers.SourceUserOrganization.ServiceRemoveUserFromOrganization)
		userOrgRoutes.PUT("/:id", resource.ServersControllers.SourceUserOrganization.ServiceUpdateUserOrganization)
		userOrgRoutes.GET("/user/:userId", resource.ServersControllers.SourceUserOrganization.ServiceGetUserOrganizations)
		userOrgRoutes.GET("/org/:orgId", resource.ServersControllers.SourceUserOrganization.ServiceGetOrganizationUsers)
	}
}

func setupUserProjectRoutes(r gin.IRouter) {
	// Rotas para gerenciar relacionamento usuário-projeto (todas protegidas)
	// Nota: POST /user-project/user/:userId está registrado como rota pública para seeding
	userProjRoutes := r.Group("/user-project")
	{
		// POST removed - registered as public route for seeding
		userProjRoutes.DELETE("/user/:userId/proj/:projectId", resource.ServersControllers.SourceUserProject.ServiceRemoveUserFromProject)
		userProjRoutes.PUT("/:id", resource.ServersControllers.SourceUserProject.ServiceUpdateUserProject)
		userProjRoutes.GET("/user/:userId", resource.ServersControllers.SourceUserProject.ServiceGetUserProjects)
		userProjRoutes.GET("/user/:userId/org/:orgId", resource.ServersControllers.SourceUserProject.ServiceGetUserProjectsByOrganization)
		userProjRoutes.GET("/proj/:projectId", resource.ServersControllers.SourceUserProject.ServiceGetProjectUsers)
	}
}

func setupProductRoutes(r gin.IRouter) {
	productRoutes := r.Group("/product")
	{
		// Rotas de leitura - requer permissão client_products_view
		productRoutes.GET("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1),
			resource.ServersControllers.SourceProducts.ServiceGetProduct)
		productRoutes.GET("/purchase/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1),
			resource.ServersControllers.SourceProducts.ServiceGetProductByPurchase)
		productRoutes.GET("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1),
			resource.ServersControllers.SourceProducts.ServiceListProducts)
		productRoutes.GET("/by-tag",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1),
			resource.ServersControllers.SourceProducts.ServiceGetProductsByTag)

		// Rotas de escrita - requerem permissões específicas + limite de pacote
		productRoutes.POST("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_create", 1),
			middleware.PackageLimitMiddleware(resource.Handlers.HandlerLimits, handler.LimitProducts),
			resource.ServersControllers.SourceProducts.ServiceCreateProduct)
		productRoutes.PUT("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_edit", 1),
			resource.ServersControllers.SourceProducts.ServiceUpdateProduct)
		productRoutes.PUT("/:id/image",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_edit", 1),
			resource.ServersControllers.SourceProducts.ServiceUpdateProductImage)
		productRoutes.DELETE("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_delete", 1),
			resource.ServersControllers.SourceProducts.ServiceDeleteProduct)

		// Tag management - requer permissão de edição
		productRoutes.GET("/:id/tags",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1),
			resource.ServersControllers.SourceProducts.ServiceGetProductTags)
		productRoutes.POST("/:id/tags",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_edit", 1),
			resource.ServersControllers.SourceProducts.ServiceAddTagToProduct)
		productRoutes.DELETE("/:id/tags/:tagId",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_edit", 1),
			resource.ServersControllers.SourceProducts.ServiceRemoveTagFromProduct)

		// Order and status management - requer permissão de edição
		productRoutes.PUT("/:id/order",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_edit", 1),
			resource.ServersControllers.SourceProducts.ServiceUpdateProductOrder)
		productRoutes.PUT("/:id/status",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_edit", 1),
			resource.ServersControllers.SourceProducts.ServiceUpdateProductStatus)

		// Filtering by menu structure - requer permissão de visualização
		productRoutes.GET("/type/:type",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1),
			resource.ServersControllers.SourceProducts.ServiceGetProductsByType)
		productRoutes.GET("/category/:categoryId",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1),
			resource.ServersControllers.SourceProducts.ServiceGetProductsByCategory)
		productRoutes.GET("/subcategory/:subcategoryId",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1),
			resource.ServersControllers.SourceProducts.ServiceGetProductsBySubcategory)
	}
}

func setupTableRoutes(r gin.IRouter) {
	tableRoutes := r.Group("/table")
	{
		// Rotas de leitura - requer permissão client_tables_view
		tableRoutes.GET("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tables_view", 1),
			resource.ServersControllers.SourceTables.ServiceGetTable)
		tableRoutes.GET("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tables_view", 1),
			resource.ServersControllers.SourceTables.ServiceListTables)

		// Rotas de escrita - requerem permissões específicas + limite de pacote
		tableRoutes.POST("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tables_create", 1),
			middleware.PackageLimitMiddleware(resource.Handlers.HandlerLimits, handler.LimitTables),
			resource.ServersControllers.SourceTables.ServiceCreateTable)
		tableRoutes.PUT("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tables_edit", 1),
			resource.ServersControllers.SourceTables.ServiceUpdateTable)
		tableRoutes.DELETE("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tables_delete", 1),
			resource.ServersControllers.SourceTables.ServiceDeleteTable)
	}
}

func setupWaitlistRoutes(r gin.IRouter) {
	waitlistRoutes := r.Group("/waitlist")
	{
		// Verificar se módulo de fila de espera está disponível
		waitlistRoutes.Use(middleware.ModuleRequiredMiddleware(resource.Handlers.HandlerLimits, "client_waitlist"))

		// Rotas de leitura
		waitlistRoutes.GET("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_waitlist_view", 1),
			resource.ServersControllers.SourceWaitlist.ServiceGetWaitlist)
		waitlistRoutes.GET("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_waitlist_view", 1),
			resource.ServersControllers.SourceWaitlist.ServiceListWaitlists)

		// Rotas de escrita
		waitlistRoutes.POST("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_waitlist_create", 1),
			resource.ServersControllers.SourceWaitlist.ServiceCreateWaitlist)
		waitlistRoutes.PUT("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_waitlist_edit", 1),
			resource.ServersControllers.SourceWaitlist.ServiceUpdateWaitlist)
		waitlistRoutes.DELETE("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_waitlist_delete", 1),
			resource.ServersControllers.SourceWaitlist.ServiceDeleteWaitlist)
	}
}

func setupReservationRoutes(r gin.IRouter) {
	reservationRoutes := r.Group("/reservation")
	{
		// Verificar se módulo de reservas está disponível
		reservationRoutes.Use(middleware.ModuleRequiredMiddleware(resource.Handlers.HandlerLimits, "client_reservations"))

		// Rotas de leitura
		reservationRoutes.GET("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_reservations_view", 1),
			resource.ServersControllers.SourceReservation.ServiceGetReservation)
		reservationRoutes.GET("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_reservations_view", 1),
			resource.ServersControllers.SourceReservation.ServiceListReservations)

		// Rotas de escrita + limite de reservas por dia
		reservationRoutes.POST("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_reservations_create", 1),
			middleware.PackageLimitMiddleware(resource.Handlers.HandlerLimits, handler.LimitReservationsDay),
			resource.ServersControllers.SourceReservation.ServiceCreateReservation)
		reservationRoutes.PUT("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_reservations_edit", 1),
			resource.ServersControllers.SourceReservation.ServiceUpdateReservation)
		reservationRoutes.DELETE("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_reservations_delete", 1),
			resource.ServersControllers.SourceReservation.ServiceDeleteReservation)
	}
}

func setupCustomerRoutes(r gin.IRouter) {
	customerRoutes := r.Group("/customer")
	{
		// Rotas de leitura
		customerRoutes.GET("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_customers_view", 1),
			resource.ServersControllers.SourceCustomer.ServiceGetCustomer)
		customerRoutes.GET("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_customers_view", 1),
			resource.ServersControllers.SourceCustomer.ServiceListCustomers)

		// Rotas de escrita
		customerRoutes.POST("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_customers_create", 1),
			resource.ServersControllers.SourceCustomer.ServiceCreateCustomer)
		customerRoutes.PUT("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_customers_edit", 1),
			resource.ServersControllers.SourceCustomer.ServiceUpdateCustomer)
		customerRoutes.DELETE("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_customers_delete", 1),
			resource.ServersControllers.SourceCustomer.ServiceDeleteCustomer)
	}
}

func setupOrderRoutes(r gin.IRouter) {
	orderRoutes := r.Group("/order")
	{
		// Rotas de leitura
		orderRoutes.GET("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_orders_view", 1),
			resource.ServersControllers.SourceOrders.GetOrderById)
		orderRoutes.GET("/:id/progress",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_orders_view", 1),
			resource.ServersControllers.SourceOrders.GetOrderProgress)
		orderRoutes.GET("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_orders_view", 1),
			resource.ServersControllers.SourceOrders.ListOrders)

		// Rotas de escrita
		orderRoutes.POST("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_orders_create", 1),
			resource.ServersControllers.SourceOrders.CreateOrder)
		orderRoutes.PUT("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_orders_edit", 1),
			resource.ServersControllers.SourceOrders.UpdateOrder)
		orderRoutes.PUT("/:id/status",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_orders_edit", 1),
			resource.ServersControllers.SourceOrders.UpdateOrderStatus)
		orderRoutes.DELETE("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_orders_delete", 1),
			resource.ServersControllers.SourceOrders.SoftDeleteOrder)
	}

	// Kitchen specific routes
	kitchenRoutes := r.Group("/kitchen")
	{
		kitchenRoutes.GET("/queue", resource.ServersControllers.SourceOrders.GetKitchenQueue)
	}
}

// setupProjectRoutes configura rotas para projetos (todas protegidas)
// Nota: POST /project está registrado como rota pública para seeding
func setupProjectRoutes(r gin.IRouter) {
	projectRoutes := r.Group("/project")
	{
		projectRoutes.GET("/:id", resource.ServersControllers.SourceProject.GetProjectById)
		projectRoutes.GET("", resource.ServersControllers.SourceProject.GetProjectsByOrganization)
		projectRoutes.GET("/organization/:orgId", resource.ServersControllers.SourceProject.GetProjectsByOrganizationId)
		projectRoutes.GET("/active", resource.ServersControllers.SourceProject.GetActiveProjects)
		// POST removed - registered as public route for seeding
		projectRoutes.PUT("/:id", resource.ServersControllers.SourceProject.UpdateProject)
		projectRoutes.DELETE("/:id", resource.ServersControllers.SourceProject.SoftDeleteProject)
		projectRoutes.DELETE("/:id/permanent", resource.ServersControllers.SourceProject.HardDeleteProject)
	}
}

// setupSettingsRoutes configura rotas para configurações
func setupSettingsRoutes(r gin.IRouter) {
	settingsRoutes := r.Group("/settings")
	{
		settingsRoutes.GET("", resource.ServersControllers.SourceSettings.GetSettingsByProject)
		settingsRoutes.PUT("", resource.ServersControllers.SourceSettings.UpdateSettings)
	}
}

// setupDisplaySettingsRoutes configura rotas para configurações de exibição de produtos
func setupDisplaySettingsRoutes(r gin.IRouter) {
	displaySettingsRoutes := r.Group("/project/settings/display")
	{
		displaySettingsRoutes.GET("", resource.ServersControllers.SourceDisplaySettings.GetDisplaySettings)
		displaySettingsRoutes.PUT("", resource.ServersControllers.SourceDisplaySettings.UpdateDisplaySettings)
		displaySettingsRoutes.POST("/reset", resource.ServersControllers.SourceDisplaySettings.ResetDisplaySettings)
	}
}

// setupThemeCustomizationRoutes configura rotas para customização de tema
func setupThemeCustomizationRoutes(r gin.IRouter) {
	themeRoutes := r.Group("/project/settings/theme")
	{
		themeRoutes.GET("", resource.ServersControllers.SourceThemeCustomization.GetTheme)
		themeRoutes.POST("", resource.ServersControllers.SourceThemeCustomization.CreateOrUpdateTheme)
		themeRoutes.PUT("", resource.ServersControllers.SourceThemeCustomization.CreateOrUpdateTheme)
		themeRoutes.POST("/reset", resource.ServersControllers.SourceThemeCustomization.ResetTheme)
		themeRoutes.DELETE("", resource.ServersControllers.SourceThemeCustomization.DeleteTheme)
	}
}

// setupSidebarConfigRoutes configura rotas para configuração da sidebar
func setupSidebarConfigRoutes(r gin.IRouter) {
	sidebarConfigRoutes := r.Group("/sidebar-config")
	{
		// GET - Qualquer usuário autenticado pode ler
		sidebarConfigRoutes.GET("", resource.ServersControllers.SourceSidebarConfig.GetConfig)
		// PUT - Apenas Master Admin (verificado no handler)
		sidebarConfigRoutes.PUT("", resource.ServersControllers.SourceSidebarConfig.UpdateConfig)
		// POST /reset - Apenas Master Admin (verificado no handler)
		sidebarConfigRoutes.POST("/reset", resource.ServersControllers.SourceSidebarConfig.ResetConfig)
	}
}

// setupEnvironmentRoutes configura rotas para ambientes
func setupEnvironmentRoutes(r gin.IRouter) {
	environmentRoutes := r.Group("/environment")
	{
		environmentRoutes.GET("/:id", resource.ServersControllers.SourceEnvironment.GetEnvironmentById)
		environmentRoutes.GET("", resource.ServersControllers.SourceEnvironment.GetEnvironmentsByProject)
		environmentRoutes.GET("/active", resource.ServersControllers.SourceEnvironment.GetActiveEnvironments)
		environmentRoutes.POST("", resource.ServersControllers.SourceEnvironment.CreateEnvironment)
		environmentRoutes.PUT("/:id", resource.ServersControllers.SourceEnvironment.UpdateEnvironment)
		environmentRoutes.DELETE("/:id", resource.ServersControllers.SourceEnvironment.SoftDeleteEnvironment)
	}
}

// setupNotificationRoutes configura rotas para notificações (SPRINT 2)
func setupNotificationRoutes(r *gin.Engine) {
	// Webhooks públicos (não requerem autenticação)
	webhookRoutes := r.Group("/webhook")
	{
		// Twilio status callback
		webhookRoutes.POST("/twilio/status", resource.ServersControllers.SourceNotification.TwilioWebhookStatus)
		// Twilio inbound messages (com org/project na URL)
		webhookRoutes.POST("/twilio/inbound/:orgId/:projectId", resource.ServersControllers.SourceNotification.TwilioWebhookInbound)
	}

	// APIs de notificação (protegidas)
	notificationRoutes := r.Group("/notification")
	notificationRoutes.Use(middleware.AuthMiddleware())
	notificationRoutes.Use(middleware.HeaderValidationMiddleware())
	{
		// Enviar notificação manual
		notificationRoutes.POST("/send", resource.ServersControllers.SourceNotification.SendNotification)
		// Processar evento de notificação
		notificationRoutes.POST("/event", resource.ServersControllers.SourceNotification.ProcessEvent)
		// Logs de notificação
		notificationRoutes.GET("/logs/:orgId/:projectId", resource.ServersControllers.SourceNotification.GetNotificationLogs)
		// Templates
		notificationRoutes.GET("/templates/:orgId/:projectId", resource.ServersControllers.SourceNotification.GetNotificationTemplates)
		notificationRoutes.POST("/template", resource.ServersControllers.SourceNotification.CreateNotificationTemplate)
		notificationRoutes.PUT("/template", resource.ServersControllers.SourceNotification.UpdateNotificationTemplate)
		// Configurações
		notificationRoutes.POST("/config", resource.ServersControllers.SourceNotification.CreateOrUpdateNotificationConfig)

		// Fila de revisão de respostas
		notificationRoutes.GET("/review-queue/:orgId/:projectId", resource.ServersControllers.SourceNotification.GetReviewQueue)
		notificationRoutes.POST("/review-queue/:id/approve", resource.ServersControllers.SourceNotification.ApproveReviewItem)
		notificationRoutes.POST("/review-queue/:id/reject", resource.ServersControllers.SourceNotification.RejectReviewItem)
		notificationRoutes.POST("/review-queue/:id/custom", resource.ServersControllers.SourceNotification.ExecuteCustomAction)
	}
}

// setupReportsRoutes configura rotas para relatórios
func setupReportsRoutes(r gin.IRouter) {
	reportsRoutes := r.Group("/reports")
	{
		reportsRoutes.GET("/occupancy", resource.ServersControllers.SourceReports.GetOccupancyReport)
		reportsRoutes.GET("/reservations", resource.ServersControllers.SourceReports.GetReservationReport)
		reportsRoutes.GET("/waitlist", resource.ServersControllers.SourceReports.GetWaitlistReport)
		reportsRoutes.GET("/leads", resource.ServersControllers.SourceReports.GetLeadReport)
		reportsRoutes.GET("/export/:type", resource.ServersControllers.SourceReports.ExportReportToCSV)
	}
}

// setupOrganizationRoutes configura rotas para organizações (todas protegidas)
// Nota: POST /organization está registrado como rota pública para seeding
func setupOrganizationRoutes(r gin.IRouter) {
	organizationRoutes := r.Group("/organization")
	{
		organizationRoutes.GET("/:id", resource.ServersControllers.SourceOrganization.GetOrganizationById)
		organizationRoutes.GET("", resource.ServersControllers.SourceOrganization.ListOrganizations)
		organizationRoutes.GET("/active", resource.ServersControllers.SourceOrganization.ListActiveOrganizations)
		organizationRoutes.GET("/email", resource.ServersControllers.SourceOrganization.GetOrganizationByEmail)
		// POST removed - registered as public route for seeding
		organizationRoutes.PUT("/:id", resource.ServersControllers.SourceOrganization.UpdateOrganization)
		organizationRoutes.DELETE("/:id", resource.ServersControllers.SourceOrganization.SoftDeleteOrganization)
		organizationRoutes.DELETE("/:id/permanent", resource.ServersControllers.SourceOrganization.HardDeleteOrganization)
	}
}

// setupPublicRoutes configura rotas públicas (sem autenticação)
func setupPublicRoutes(r *gin.Engine) {
	publicRoutes := r.Group("/public")
	{
		// Cardápio público
		publicRoutes.GET("/menu/:orgId/:projId", resource.ServersControllers.SourcePublic.ServiceGetPublicMenu)

		// Categorias públicas
		publicRoutes.GET("/categories/:orgId/:projId", resource.ServersControllers.SourcePublic.ServiceGetPublicCategories)

		// Menus públicos
		publicRoutes.GET("/menus/:orgId/:projId", resource.ServersControllers.SourcePublic.ServiceGetPublicMenus)

		// Informações do projeto
		publicRoutes.GET("/project/:orgId/:projId", resource.ServersControllers.SourcePublic.ServiceGetProjectInfo)

		// Horários disponíveis
		publicRoutes.GET("/times/:orgId/:projId", resource.ServersControllers.SourcePublic.ServiceGetAvailableTimes)

		// Criar reserva pública
		publicRoutes.POST("/reservation/:orgId/:projId", resource.ServersControllers.SourcePublic.ServiceCreatePublicReservation)
	}
}

// setupUploadRoutes configura rotas de upload (com autenticação)
func setupUploadRoutes(r *gin.Engine) {
	// Rotas públicas para servir imagens estáticas
	// Nova estrutura: /uploads/orgId/projId/category/filename
	r.GET("/uploads/:orgId/:projId/:category/:filename", resource.ServersControllers.SourceUpload.ServiceGetUploadedFile)
	// Estrutura de compatibilidade: /static/category/filename (evita conflito de rotas)
	r.GET("/static/:category/:filename", resource.ServersControllers.SourceUpload.ServiceGetUploadedFile)

	// Rotas protegidas para upload (requerem autenticação)
	uploadRoutes := r.Group("/upload")
	uploadRoutes.Use(middleware.HeaderValidationMiddleware())
	{
		// Rota genérica para upload de qualquer categoria
		uploadRoutes.POST("/:category/image", resource.ServersControllers.SourceUpload.ServiceUploadImage)

		// Rota de retrocompatibilidade para produtos
		uploadRoutes.POST("/product/image", resource.ServersControllers.SourceUpload.ServiceUploadProductImage)
	}
}

// setupTagRoutes configura rotas para tags
func setupTagRoutes(r gin.IRouter) {
	tagRoutes := r.Group("/tag")
	{
		// Rotas de leitura - requer permissão client_tags_view (nível 1)
		tagRoutes.GET("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tags_view", 1),
			resource.ServersControllers.SourceTag.ServiceGetTag)
		tagRoutes.GET("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tags_view", 1),
			resource.ServersControllers.SourceTag.ServiceListTags)
		tagRoutes.GET("/active",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tags_view", 1),
			resource.ServersControllers.SourceTag.ServiceListActiveTags)
		tagRoutes.GET("/entity/:entityType",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tags_view", 1),
			resource.ServersControllers.SourceTag.ServiceGetTagsByEntityType)

		// Rotas de escrita - requerem permissões específicas (nível 1 = habilitado)
		tagRoutes.POST("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tags_create", 1),
			resource.ServersControllers.SourceTag.ServiceCreateTag)
		tagRoutes.PUT("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tags_edit", 1),
			resource.ServersControllers.SourceTag.ServiceUpdateTag)
		tagRoutes.DELETE("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tags_delete", 1),
			resource.ServersControllers.SourceTag.ServiceDeleteTag)
	}
}

// setupMenuRoutes configura rotas para cardápios/menus
func setupMenuRoutes(r gin.IRouter) {
	menuRoutes := r.Group("/menu")
	{
		// ✨ IMPORTANT: Specific routes FIRST to avoid conflicts with /:id pattern
		// Rotas de seleção inteligente (GET - sem proteção adicional)
		menuRoutes.GET("/active-now", resource.ServersControllers.SourceMenu.ServiceGetActiveMenu)
		menuRoutes.GET("/active", resource.ServersControllers.SourceMenu.ServiceListActiveMenus)
		menuRoutes.GET("/options", resource.ServersControllers.SourceMenu.ServiceGetMenuOptions)

		// Rotas de seleção (PUT/DELETE) - Master Admin Only
		menuRoutes.PUT("/:id/manual-override",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceMenu.ServiceSetMenuAsManualOverride)
		menuRoutes.DELETE("/manual-override",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceMenu.ServiceRemoveManualOverride)

		// Rotas padrão (mais genéricas - devem vir por último)
		menuRoutes.GET("/:id", resource.ServersControllers.SourceMenu.ServiceGetMenu)
		menuRoutes.GET("", resource.ServersControllers.SourceMenu.ServiceListMenus)

		// Rotas de modificação - Master Admin Only
		menuRoutes.POST("",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceMenu.ServiceCreateMenu)
		menuRoutes.PUT("/:id",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceMenu.ServiceUpdateMenu)
		menuRoutes.PUT("/:id/order",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceMenu.ServiceUpdateMenuOrder)
		menuRoutes.PUT("/:id/status",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceMenu.ServiceUpdateMenuStatus)
		menuRoutes.DELETE("/:id",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceMenu.ServiceDeleteMenu)
	}
}

// setupCategoryRoutes configura rotas para categorias
func setupCategoryRoutes(r gin.IRouter) {
	categoryRoutes := r.Group("/category")
	{
		// GET routes (sem proteção adicional)
		categoryRoutes.GET("/:id", resource.ServersControllers.SourceCategory.ServiceGetCategory)
		categoryRoutes.GET("", resource.ServersControllers.SourceCategory.ServiceListCategories)
		categoryRoutes.GET("/active", resource.ServersControllers.SourceCategory.ServiceListActiveCategories)
		categoryRoutes.GET("/menu/:menuId", resource.ServersControllers.SourceCategory.ServiceGetCategoriesByMenu)

		// POST, PUT, DELETE - Master Admin Only
		categoryRoutes.POST("",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceCategory.ServiceCreateCategory)
		categoryRoutes.PUT("/:id",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceCategory.ServiceUpdateCategory)
		categoryRoutes.PUT("/:id/order",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceCategory.ServiceUpdateCategoryOrder)
		categoryRoutes.PUT("/:id/status",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceCategory.ServiceUpdateCategoryStatus)
		categoryRoutes.DELETE("/:id",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceCategory.ServiceDeleteCategory)
	}
}

// setupSubcategoryRoutes configura rotas para subcategorias
func setupSubcategoryRoutes(r gin.IRouter) {
	subcategoryRoutes := r.Group("/subcategory")
	{
		subcategoryRoutes.GET("/:id", resource.ServersControllers.SourceSubcategory.ServiceGetSubcategory)
		subcategoryRoutes.GET("", resource.ServersControllers.SourceSubcategory.ServiceListSubcategories)
		subcategoryRoutes.GET("/active", resource.ServersControllers.SourceSubcategory.ServiceListActiveSubcategories)
		subcategoryRoutes.GET("/category/:categoryId", resource.ServersControllers.SourceSubcategory.ServiceGetSubcategoriesByCategory)
		subcategoryRoutes.POST("", resource.ServersControllers.SourceSubcategory.ServiceCreateSubcategory)
		subcategoryRoutes.PUT("/:id", resource.ServersControllers.SourceSubcategory.ServiceUpdateSubcategory)
		subcategoryRoutes.PUT("/:id/order", resource.ServersControllers.SourceSubcategory.ServiceUpdateSubcategoryOrder)
		subcategoryRoutes.PUT("/:id/status", resource.ServersControllers.SourceSubcategory.ServiceUpdateSubcategoryStatus)
		subcategoryRoutes.DELETE("/:id", resource.ServersControllers.SourceSubcategory.ServiceDeleteSubcategory)

		// Category relationship management
		subcategoryRoutes.POST("/:id/category/:categoryId", resource.ServersControllers.SourceSubcategory.ServiceAddCategoryToSubcategory)
		subcategoryRoutes.DELETE("/:id/category/:categoryId", resource.ServersControllers.SourceSubcategory.ServiceRemoveCategoryFromSubcategory)
		subcategoryRoutes.GET("/:id/categories", resource.ServersControllers.SourceSubcategory.ServiceGetSubcategoryCategories)
	}
}

// setupImageManagementRoutes configura rotas para gerenciamento de imagens (admin)
func setupImageManagementRoutes(r gin.IRouter) {
	adminRoutes := r.Group("/admin/images")
	{
		// Limpar arquivos órfãos (soft deletados, sem referências)
		adminRoutes.POST("/cleanup", resource.ServersControllers.SourceImageManagement.ServiceCleanupOrphanedFiles)

		// Obter estatísticas de imagens
		adminRoutes.GET("/stats", resource.ServersControllers.SourceImageManagement.ServiceGetImageStats)
	}
}

// setupRoleRoutes configura rotas para gerenciamento de cargos
func setupRoleRoutes(r gin.IRouter) {
	roleRoutes := r.Group("/role")
	{
		// Rotas específicas primeiro (antes de /:id para evitar conflitos)
		roleRoutes.GET("/system", resource.ServersControllers.SourceRole.ListSystemRoles)
		roleRoutes.GET("/check", resource.ServersControllers.SourceRole.CheckPermission)
		roleRoutes.GET("/my-permissions", resource.ServersControllers.SourceRole.GetMyPermissions)

		// Atribuição de cargos a usuários
		roleRoutes.POST("/assign", resource.ServersControllers.SourceRole.AssignRoleToUser)
		roleRoutes.POST("/remove", resource.ServersControllers.SourceRole.RemoveRoleFromUser)
		roleRoutes.GET("/user/:userId", resource.ServersControllers.SourceRole.GetUserRoles)
		roleRoutes.GET("/user/:userId/details", resource.ServersControllers.SourceRole.GetUserRolesWithDetails)

		// Níveis de permissão por cargo
		roleRoutes.POST("/permission-level", resource.ServersControllers.SourceRole.SetPermissionLevel)

		// CRUD de cargos (rotas com :id por último)
		roleRoutes.GET("", resource.ServersControllers.SourceRole.ListRoles)
		roleRoutes.POST("", resource.ServersControllers.SourceRole.CreateRole)
		roleRoutes.GET("/:id", resource.ServersControllers.SourceRole.GetRole)
		roleRoutes.PUT("/:id", resource.ServersControllers.SourceRole.UpdateRole)
		roleRoutes.DELETE("/:id", resource.ServersControllers.SourceRole.DeleteRole)
		roleRoutes.GET("/:id/permissions", resource.ServersControllers.SourceRole.GetRolePermissions)
	}
}

// setupModuleRoutes configura rotas para gerenciamento de módulos
func setupModuleRoutes(r gin.IRouter) {
	moduleRoutes := r.Group("/module")
	{
		// Rotas de leitura (usuários autenticados)
		moduleRoutes.GET("", resource.ServersControllers.SourceRole.ListModules)
		moduleRoutes.GET("/with-permissions", resource.ServersControllers.SourceRole.ListModulesWithPermissions)
		moduleRoutes.GET("/available", resource.ServersControllers.SourceRole.GetOrganizationModules)

		// CRUD de módulos (Master Admin only)
		moduleRoutes.POST("",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceRole.CreateModule)
		moduleRoutes.GET("/:id", resource.ServersControllers.SourceRole.GetModule)
		moduleRoutes.PUT("/:id",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceRole.UpdateModule)
		moduleRoutes.DELETE("/:id",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceRole.DeleteModule)
	}
}

// setupPermissionRoutes configura rotas para gerenciamento de permissões
func setupPermissionRoutes(r gin.IRouter) {
	permissionRoutes := r.Group("/permission")
	{
		// Rotas de leitura (usuários autenticados)
		permissionRoutes.GET("", resource.ServersControllers.SourceRole.ListPermissions)

		// CRUD de permissões (Master Admin only)
		permissionRoutes.POST("",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceRole.CreatePermission)
		permissionRoutes.GET("/:id", resource.ServersControllers.SourceRole.GetPermission)
		permissionRoutes.PUT("/:id",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceRole.UpdatePermission)
		permissionRoutes.DELETE("/:id",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceRole.DeletePermission)
	}
}

// setupPackageRoutes configura rotas para gerenciamento de pacotes/planos
func setupPackageRoutes(r gin.IRouter) {
	packageRoutes := r.Group("/package")
	{
		// Rotas específicas primeiro (antes de /:id para evitar conflitos)
		packageRoutes.GET("/subscription", resource.ServersControllers.SourceRole.GetOrganizationSubscription)
		packageRoutes.GET("/usage-limits", resource.ServersControllers.SourceRole.GetUsageAndLimits)
		packageRoutes.POST("/subscribe", resource.ServersControllers.SourceRole.SubscribeOrganization)

		// Lista de todas as assinaturas (Master Admin only)
		packageRoutes.GET("/subscriptions",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceRole.ListAllSubscriptions)

		// Gerenciar assinatura de organização específica (Master Admin only)
		packageRoutes.POST("/subscription/:orgId",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceRole.CreateOrganizationSubscription)
		packageRoutes.PUT("/subscription/:orgId",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceRole.UpdateOrganizationSubscription)
		packageRoutes.DELETE("/subscription/:orgId",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceRole.CancelOrganizationSubscription)

		// Rotas de leitura (usuários autenticados)
		packageRoutes.GET("", resource.ServersControllers.SourceRole.ListPackages)

		// CRUD de pacotes (Master Admin only)
		packageRoutes.POST("",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceRole.CreatePackage)
		packageRoutes.GET("/:id", resource.ServersControllers.SourceRole.GetPackageWithModules)
		packageRoutes.PUT("/:id",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceRole.UpdatePackage)
		packageRoutes.DELETE("/:id",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceRole.DeletePackage)

		// Gerenciar módulos em pacotes (Master Admin only)
		packageRoutes.POST("/:id/modules/:moduleId",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceRole.AddModuleToPackage)
		packageRoutes.DELETE("/:id/modules/:moduleId",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceRole.RemoveModuleFromPackage)

		// Gerenciar limites de pacotes (Master Admin only)
		packageRoutes.GET("/:id/limits", resource.ServersControllers.SourceRole.GetPackageLimits)
		packageRoutes.POST("/:id/limits",
			middleware.MasterAdminOnlyMiddleware(),
			resource.ServersControllers.SourceRole.SetPackageLimit)
	}
}

// setupPlanChangeRequestRoutes configura rotas para solicitações de mudança de plano
func setupPlanChangeRequestRoutes(r gin.IRouter) {
	// Client routes - usuários podem solicitar mudanças de plano
	clientRoutes := r.Group("/plan-change-request")
	{
		clientRoutes.POST("", resource.ServersControllers.SourcePlanChangeRequest.CreateRequest)
		clientRoutes.GET("/my-requests", resource.ServersControllers.SourcePlanChangeRequest.GetMyRequests)
		clientRoutes.GET("/:id", resource.ServersControllers.SourcePlanChangeRequest.GetRequestById)
		clientRoutes.POST("/:id/cancel", resource.ServersControllers.SourcePlanChangeRequest.CancelRequest)
	}

	// Admin routes - apenas admins podem gerenciar solicitações
	adminRoutes := r.Group("/admin/plan-change-request")
	adminRoutes.Use(middleware.MasterAdminOnlyMiddleware())
	{
		adminRoutes.GET("", resource.ServersControllers.SourcePlanChangeRequest.GetAllRequests)
		adminRoutes.GET("/pending", resource.ServersControllers.SourcePlanChangeRequest.GetPendingRequests)
		adminRoutes.GET("/organization/:orgId", resource.ServersControllers.SourcePlanChangeRequest.GetRequestsByOrganization)
		adminRoutes.POST("/:id/approve", resource.ServersControllers.SourcePlanChangeRequest.ApproveRequest)
		adminRoutes.POST("/:id/reject", resource.ServersControllers.SourcePlanChangeRequest.RejectRequest)
	}
}
