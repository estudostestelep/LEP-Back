package routes

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Public routes (no authentication required)
	r.POST("/login", resource.ServersControllers.SourceAuth.ServiceLogin)
	r.POST("/user", resource.ServersControllers.SourceUsers.ServiceCreateUser)

	// Protected routes (authentication required)
	r.POST("/logout", resource.ServersControllers.SourceAuth.ServiceLogout)
	r.POST("/checkToken", resource.ServersControllers.SourceAuth.ServiceValidateToken)

	// User routes
	setupUserRoutes(r)

	// Product routes
	setupProductRoutes(r)

	// Table routes
	setupTableRoutes(r)

	// Waitlist routes
	setupWaitlistRoutes(r)

	// Reservation routes
	setupReservationRoutes(r)

	// Customer routes
	setupCustomerRoutes(r)

	// Order routes
	setupOrderRoutes(r)

	// Project routes (SPRINT 1)
	setupProjectRoutes(r)

	// Settings routes (SPRINT 1)
	setupSettingsRoutes(r)

	// Environment routes (SPRINT 1)
	setupEnvironmentRoutes(r)

	// Notification routes (SPRINT 2)
	setupNotificationRoutes(r)
}

func setupUserRoutes(r *gin.Engine) {
	userRoutes := r.Group("/user")
	{
		userRoutes.GET("/:id", resource.ServersControllers.SourceUsers.ServiceGetUser)
		userRoutes.GET("/group/:id", resource.ServersControllers.SourceUsers.ServiceGetUserByGroup)
		userRoutes.GET("", resource.ServersControllers.SourceUsers.ServiceListUsers) // Endpoint de listagem
		userRoutes.PUT("/:id", resource.ServersControllers.SourceUsers.ServiceUpdateUser)
		userRoutes.DELETE("/:id", resource.ServersControllers.SourceUsers.ServiceDeleteUser)
	}
}

func setupProductRoutes(r *gin.Engine) {
	productRoutes := r.Group("/product")
	{
		productRoutes.GET("/:id", resource.ServersControllers.SourceProducts.ServiceGetProduct)
		productRoutes.GET("/purchase/:id", resource.ServersControllers.SourceProducts.ServiceGetProductByPurchase)
		productRoutes.GET("", resource.ServersControllers.SourceProducts.ServiceListProducts) // Endpoint de listagem
		productRoutes.POST("", resource.ServersControllers.SourceProducts.ServiceCreateProduct)
		productRoutes.PUT("/:id", resource.ServersControllers.SourceProducts.ServiceUpdateProduct)
		productRoutes.DELETE("/:id", resource.ServersControllers.SourceProducts.ServiceDeleteProduct)
	}
}

func setupTableRoutes(r *gin.Engine) {
	tableRoutes := r.Group("/table")
	{
		tableRoutes.GET("/:id", resource.ServersControllers.SourceTables.ServiceGetTable)
		tableRoutes.GET("", resource.ServersControllers.SourceTables.ServiceListTables)
		tableRoutes.POST("", resource.ServersControllers.SourceTables.ServiceCreateTable)
		tableRoutes.PUT("/:id", resource.ServersControllers.SourceTables.ServiceUpdateTable)
		tableRoutes.DELETE("/:id", resource.ServersControllers.SourceTables.ServiceDeleteTable)
	}
}

func setupWaitlistRoutes(r *gin.Engine) {
	waitlistRoutes := r.Group("/waitlist")
	{
		waitlistRoutes.GET("/:id", resource.ServersControllers.SourceWaitlist.ServiceGetWaitlist)
		waitlistRoutes.GET("", resource.ServersControllers.SourceWaitlist.ServiceListWaitlists)
		waitlistRoutes.POST("", resource.ServersControllers.SourceWaitlist.ServiceCreateWaitlist)
		waitlistRoutes.PUT("/:id", resource.ServersControllers.SourceWaitlist.ServiceUpdateWaitlist)
		waitlistRoutes.DELETE("/:id", resource.ServersControllers.SourceWaitlist.ServiceDeleteWaitlist)
	}
}

func setupReservationRoutes(r *gin.Engine) {
	reservationRoutes := r.Group("/reservation")
	{
		reservationRoutes.GET("/:id", resource.ServersControllers.SourceReservation.ServiceGetReservation)
		reservationRoutes.GET("", resource.ServersControllers.SourceReservation.ServiceListReservations)
		reservationRoutes.POST("", resource.ServersControllers.SourceReservation.ServiceCreateReservation)
		reservationRoutes.PUT("/:id", resource.ServersControllers.SourceReservation.ServiceUpdateReservation)
		reservationRoutes.DELETE("/:id", resource.ServersControllers.SourceReservation.ServiceDeleteReservation)
	}
}

func setupCustomerRoutes(r *gin.Engine) {
	customerRoutes := r.Group("/customer")
	{
		customerRoutes.GET("/:id", resource.ServersControllers.SourceCustomer.ServiceGetCustomer)
		customerRoutes.GET("", resource.ServersControllers.SourceCustomer.ServiceListCustomers)
		customerRoutes.POST("", resource.ServersControllers.SourceCustomer.ServiceCreateCustomer)
		customerRoutes.PUT("/:id", resource.ServersControllers.SourceCustomer.ServiceUpdateCustomer)
		customerRoutes.DELETE("/:id", resource.ServersControllers.SourceCustomer.ServiceDeleteCustomer)
	}
}

func setupOrderRoutes(r *gin.Engine) {
	orderRoutes := r.Group("/order")
	{
		orderRoutes.GET("/:id", resource.ServersControllers.SourceOrders.GetOrderById)
		orderRoutes.GET("/:id/progress", resource.ServersControllers.SourceOrders.GetOrderProgress)
		orderRoutes.GET("", resource.ServersControllers.SourceOrders.ListOrders) // Moved to maintain consistency
		orderRoutes.POST("", resource.ServersControllers.SourceOrders.CreateOrder)
		orderRoutes.PUT("/:id", resource.ServersControllers.SourceOrders.UpdateOrder)
		orderRoutes.PUT("/:id/status", resource.ServersControllers.SourceOrders.UpdateOrderStatus)
		orderRoutes.DELETE("/:id", resource.ServersControllers.SourceOrders.SoftDeleteOrder)
	}

	// Kitchen specific routes
	kitchenRoutes := r.Group("/kitchen")
	{
		kitchenRoutes.GET("/queue", resource.ServersControllers.SourceOrders.GetKitchenQueue)
	}
}

// setupProjectRoutes configura rotas para projetos
func setupProjectRoutes(r *gin.Engine) {
	projectRoutes := r.Group("/project")
	{
		projectRoutes.GET("/:id", resource.ServersControllers.SourceProject.GetProjectById)
		projectRoutes.GET("", resource.ServersControllers.SourceProject.GetProjectsByOrganization)
		projectRoutes.GET("/active", resource.ServersControllers.SourceProject.GetActiveProjects)
		projectRoutes.POST("", resource.ServersControllers.SourceProject.CreateProject)
		projectRoutes.PUT("/:id", resource.ServersControllers.SourceProject.UpdateProject)
		projectRoutes.DELETE("/:id", resource.ServersControllers.SourceProject.SoftDeleteProject)
	}
}

// setupSettingsRoutes configura rotas para configurações
func setupSettingsRoutes(r *gin.Engine) {
	settingsRoutes := r.Group("/settings")
	{
		settingsRoutes.GET("", resource.ServersControllers.SourceSettings.GetSettingsByProject)
		settingsRoutes.PUT("", resource.ServersControllers.SourceSettings.UpdateSettings)
	}
}

// setupEnvironmentRoutes configura rotas para ambientes
func setupEnvironmentRoutes(r *gin.Engine) {
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
	}
}
