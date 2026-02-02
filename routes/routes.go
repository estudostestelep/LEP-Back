package routes

import (
	"lep/handler"
	"lep/middleware"
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas as rotas da aplicação
func SetupRoutes(r *gin.Engine) {

	// =============================================================================
	// 1. ROTAS PÚBLICAS (sem autenticação)
	// =============================================================================

	// Auth
	r.POST("/login", resource.ServersControllers.SourceAuth.ServiceLogin)
	r.POST("/admin/login", resource.ServersControllers.SourceAuthAdmin.ServiceAdminLogin)
	r.POST("/client/login", resource.ServersControllers.SourceAuthClient.ServiceClientLogin)
	r.POST("/create-organization", resource.ServersControllers.SourceOrganization.ServiceCreateOrganizationBootstrap)
	r.POST("/admin/reset-passwords", resource.ServersControllers.SourceAdmin.ServiceResetPasswords)
	r.POST("/admin/run-migration", resource.ServersControllers.SourceAdmin.ServiceRunDevMigration)

	// Tenant
	r.GET("/tenant/resolve", resource.ServersControllers.SourceTenant.ServiceResolveTenant)

	// Seeding
	r.POST("/organization", resource.ServersControllers.SourceOrganization.CreateOrganization)
	r.POST("/project", resource.ServersControllers.SourceProject.CreateProject)

	// Uploads públicos (servir arquivos)
	r.GET("/uploads/:orgId/:projId/:category/:filename", resource.ServersControllers.SourceUpload.ServiceGetUploadedFile)
	r.GET("/static/:category/:filename", resource.ServersControllers.SourceUpload.ServiceGetUploadedFile)

	// Upload protegido por header (sem auth completo)
	uploadRoutes := r.Group("/upload")
	uploadRoutes.Use(middleware.HeaderValidationMiddleware())
	uploadRoutes.POST("/:category/image", resource.ServersControllers.SourceUpload.ServiceUploadImage)
	uploadRoutes.POST("/product/image", resource.ServersControllers.SourceUpload.ServiceUploadProductImage)

	// Webhooks
	webhook := r.Group("/webhook")
	webhook.POST("/twilio/status", resource.ServersControllers.SourceNotification.TwilioWebhookStatus)
	webhook.POST("/twilio/inbound/:orgId/:projectId", resource.ServersControllers.SourceNotification.TwilioWebhookInbound)

	// Rotas públicas de menu/reserva
	publicRoutes := r.Group("/public")
	publicRoutes.GET("/menu/:orgId/:projId", resource.ServersControllers.SourcePublic.ServiceGetPublicMenu)
	publicRoutes.GET("/categories/:orgId/:projId", resource.ServersControllers.SourcePublic.ServiceGetPublicCategories)
	publicRoutes.GET("/menus/:orgId/:projId", resource.ServersControllers.SourcePublic.ServiceGetPublicMenus)
	publicRoutes.GET("/project/:orgId/:projId", resource.ServersControllers.SourcePublic.ServiceGetProjectInfo)
	publicRoutes.GET("/times/:orgId/:projId", resource.ServersControllers.SourcePublic.ServiceGetAvailableTimes)
	publicRoutes.POST("/reservation/:orgId/:projId", resource.ServersControllers.SourcePublic.ServiceCreatePublicReservation)
	publicRoutes.GET("/project/resolve", resource.ServersControllers.SourcePublic.ServiceResolveProject)
	publicRoutes.GET("/menu/org/:orgSlug", resource.ServersControllers.SourcePublic.ServiceGetPublicMenuBySlug)
	publicRoutes.GET("/menu/org/:orgSlug/:projectSlug", resource.ServersControllers.SourcePublic.ServiceGetPublicMenuBySlug)
	publicRoutes.GET("/categories/org/:orgSlug", resource.ServersControllers.SourcePublic.ServiceGetPublicCategoriesBySlug)
	publicRoutes.GET("/categories/org/:orgSlug/:projectSlug", resource.ServersControllers.SourcePublic.ServiceGetPublicCategoriesBySlug)
	publicRoutes.GET("/menus/org/:orgSlug", resource.ServersControllers.SourcePublic.ServiceGetPublicMenusBySlug)
	publicRoutes.GET("/menus/org/:orgSlug/:projectSlug", resource.ServersControllers.SourcePublic.ServiceGetPublicMenusBySlug)
	publicRoutes.GET("/project/org/:orgSlug", resource.ServersControllers.SourcePublic.ServiceGetProjectInfoBySlug)
	publicRoutes.GET("/project/org/:orgSlug/:projectSlug", resource.ServersControllers.SourcePublic.ServiceGetProjectInfoBySlug)
	publicRoutes.GET("/times/org/:orgSlug", resource.ServersControllers.SourcePublic.ServiceGetAvailableTimesBySlug)
	publicRoutes.GET("/times/org/:orgSlug/:projectSlug", resource.ServersControllers.SourcePublic.ServiceGetAvailableTimesBySlug)
	publicRoutes.POST("/reservation/org/:orgSlug", resource.ServersControllers.SourcePublic.ServiceCreatePublicReservationBySlug)
	publicRoutes.POST("/reservation/org/:orgSlug/:projectSlug", resource.ServersControllers.SourcePublic.ServiceCreatePublicReservationBySlug)

	// =============================================================================
	// 2. ROTAS PROTEGIDAS (auth + headers obrigatórios)
	// =============================================================================

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.HeaderValidationMiddleware())

	// Auth protegido
	protected.POST("/logout", resource.ServersControllers.SourceAuth.ServiceLogout)
	protected.POST("/checkToken", resource.ServersControllers.SourceAuth.ServiceValidateToken)

	// =============================================================================
	// 3. ROTAS ADMIN (/admin/*) - Requer AdminScopeMiddleware
	// =============================================================================

	admin := protected.Group("/admin")
	admin.Use(middleware.AdminScopeMiddleware(resource.Handlers.HandlerAuth))

	// Admin > Menu
	adminMenu := admin.Group("/menu")
	adminMenu.POST("", resource.ServersControllers.SourceMenu.ServiceCreateMenu)
	adminMenu.PUT("/:id", resource.ServersControllers.SourceMenu.ServiceUpdateMenu)
	adminMenu.PUT("/:id/order", resource.ServersControllers.SourceMenu.ServiceUpdateMenuOrder)
	adminMenu.PUT("/:id/status", resource.ServersControllers.SourceMenu.ServiceUpdateMenuStatus)
	adminMenu.PUT("/:id/manual-override", resource.ServersControllers.SourceMenu.ServiceSetMenuAsManualOverride)
	adminMenu.DELETE("/manual-override", resource.ServersControllers.SourceMenu.ServiceRemoveManualOverride)
	adminMenu.DELETE("/:id", resource.ServersControllers.SourceMenu.ServiceDeleteMenu)

	// Admin > Category
	adminCategory := admin.Group("/category")
	adminCategory.POST("", resource.ServersControllers.SourceCategory.ServiceCreateCategory)
	adminCategory.PUT("/:id", resource.ServersControllers.SourceCategory.ServiceUpdateCategory)
	adminCategory.PUT("/:id/order", resource.ServersControllers.SourceCategory.ServiceUpdateCategoryOrder)
	adminCategory.PUT("/:id/status", resource.ServersControllers.SourceCategory.ServiceUpdateCategoryStatus)
	adminCategory.DELETE("/:id", resource.ServersControllers.SourceCategory.ServiceDeleteCategory)

	// Admin > Subcategory
	adminSubcategory := admin.Group("/subcategory")
	adminSubcategory.POST("", resource.ServersControllers.SourceSubcategory.ServiceCreateSubcategory)
	adminSubcategory.PUT("/:id", resource.ServersControllers.SourceSubcategory.ServiceUpdateSubcategory)
	adminSubcategory.PUT("/:id/order", resource.ServersControllers.SourceSubcategory.ServiceUpdateSubcategoryOrder)
	adminSubcategory.PUT("/:id/status", resource.ServersControllers.SourceSubcategory.ServiceUpdateSubcategoryStatus)
	adminSubcategory.DELETE("/:id", resource.ServersControllers.SourceSubcategory.ServiceDeleteSubcategory)
	adminSubcategory.POST("/:id/category/:categoryId", resource.ServersControllers.SourceSubcategory.ServiceAddCategoryToSubcategory)
	adminSubcategory.DELETE("/:id/category/:categoryId", resource.ServersControllers.SourceSubcategory.ServiceRemoveCategoryFromSubcategory)

	// Admin > Role
	adminRole := admin.Group("/role")
	adminRole.POST("", resource.ServersControllers.SourceRole.CreateRole)
	adminRole.PUT("/:id", resource.ServersControllers.SourceRole.UpdateRole)
	adminRole.DELETE("/:id", resource.ServersControllers.SourceRole.DeleteRole)

	// Admin > Module
	adminModule := admin.Group("/module")
	adminModule.POST("", resource.ServersControllers.SourceRole.CreateModule)
	adminModule.PUT("/:id", resource.ServersControllers.SourceRole.UpdateModule)
	adminModule.DELETE("/:id", resource.ServersControllers.SourceRole.DeleteModule)

	// Admin > Permission
	adminPermission := admin.Group("/permission")
	adminPermission.POST("", resource.ServersControllers.SourceRole.CreatePermission)
	adminPermission.PUT("/:id", resource.ServersControllers.SourceRole.UpdatePermission)
	adminPermission.DELETE("/:id", resource.ServersControllers.SourceRole.DeletePermission)

	// Admin > Package (Plans)
	adminPackage := admin.Group("/package")
	adminPackage.POST("", resource.ServersControllers.SourceRole.CreatePackage)
	adminPackage.PUT("/:id", resource.ServersControllers.SourceRole.UpdatePackage)
	adminPackage.DELETE("/:id", resource.ServersControllers.SourceRole.DeletePackage)
	adminPackage.POST("/:id/modules/:moduleId", resource.ServersControllers.SourceRole.AddModuleToPackage)
	adminPackage.DELETE("/:id/modules/:moduleId", resource.ServersControllers.SourceRole.RemoveModuleFromPackage)
	adminPackage.POST("/:id/limits", resource.ServersControllers.SourceRole.SetPackageLimit)
	adminPackage.GET("/subscriptions", resource.ServersControllers.SourceRole.ListAllSubscriptions)
	adminPackage.POST("/subscription/:orgId", resource.ServersControllers.SourceRole.CreateOrganizationSubscription)
	adminPackage.PUT("/subscription/:orgId", resource.ServersControllers.SourceRole.UpdateOrganizationSubscription)
	adminPackage.DELETE("/subscription/:orgId", resource.ServersControllers.SourceRole.CancelOrganizationSubscription)
	adminPackage.DELETE("/subscription/:orgId/delete", resource.ServersControllers.SourceRole.DeleteOrganizationSubscription)

	// Admin > Plan Change Requests
	adminPlanChange := admin.Group("/plan-change-request")
	adminPlanChange.GET("", resource.ServersControllers.SourcePlanChangeRequest.GetAllRequests)
	adminPlanChange.GET("/pending", resource.ServersControllers.SourcePlanChangeRequest.GetPendingRequests)
	adminPlanChange.GET("/organization/:orgId", resource.ServersControllers.SourcePlanChangeRequest.GetRequestsByOrganization)
	adminPlanChange.POST("/:id/approve", resource.ServersControllers.SourcePlanChangeRequest.ApproveRequest)
	adminPlanChange.POST("/:id/reject", resource.ServersControllers.SourcePlanChangeRequest.RejectRequest)

	// Admin > Audit Logs
	adminAuditLogs := admin.Group("/audit-logs")
	adminAuditLogs.GET("", resource.ServersControllers.SourceAdminAuditLog.ServiceListAdminAuditLogs)
	adminAuditLogs.GET("/:id", resource.ServersControllers.SourceAdminAuditLog.ServiceGetAdminAuditLog)
	adminAuditLogs.DELETE("/cleanup", resource.ServersControllers.SourceAdminAuditLog.ServiceDeleteOldLogs)

	// Admin > Client Audit Config
	adminClientAuditConfig := admin.Group("/client-audit-config")
	adminClientAuditConfig.GET("", resource.ServersControllers.SourceClientAuditLog.ServiceGetClientAuditConfig)
	adminClientAuditConfig.PUT("", resource.ServersControllers.SourceClientAuditLog.ServiceUpdateClientAuditConfig)
	adminClientAuditConfig.GET("/modules", resource.ServersControllers.SourceClientAuditLog.ServiceGetAvailableModules)

	// Admin > Images
	adminImages := admin.Group("/images")
	adminImages.POST("/cleanup", resource.ServersControllers.SourceImageManagement.ServiceCleanupOrphanedFiles)
	adminImages.GET("/stats", resource.ServersControllers.SourceImageManagement.ServiceGetImageStats)

	// Admin > Sidebar Config
	adminSidebarConfig := admin.Group("/sidebar-config")
	adminSidebarConfig.PUT("", resource.ServersControllers.SourceSidebarConfig.UpdateConfig)
	adminSidebarConfig.POST("/reset", resource.ServersControllers.SourceSidebarConfig.ResetConfig)

	// Admin > Admin Users
	adminUser := admin.Group("/admin-user")
	adminUser.GET("", resource.ServersControllers.SourceAdminUsers.ServiceListAdmins)
	adminUser.GET("/:id", resource.ServersControllers.SourceAdminUsers.ServiceGetAdmin)
	adminUser.POST("", resource.ServersControllers.SourceAdminUsers.ServiceCreateAdmin)
	adminUser.PUT("/:id", resource.ServersControllers.SourceAdminUsers.ServiceUpdateAdmin)
	adminUser.DELETE("/:id", resource.ServersControllers.SourceAdminUsers.ServiceDeleteAdmin)

	// Admin > Client Users
	adminClientUser := admin.Group("/client-user")
	adminClientUser.GET("", resource.ServersControllers.SourceClientUsers.ServiceListClients)
	adminClientUser.GET("/:id", resource.ServersControllers.SourceClientUsers.ServiceGetClient)
	adminClientUser.POST("", resource.ServersControllers.SourceClientUsers.ServiceCreateClient)
	adminClientUser.PUT("/:id", resource.ServersControllers.SourceClientUsers.ServiceUpdateClient)
	adminClientUser.DELETE("/:id", resource.ServersControllers.SourceClientUsers.ServiceDeleteClient)

	// =============================================================================
	// 4. ROTAS CLIENT (recursos do restaurante)
	// =============================================================================

	// Menu (leitura)
	menu := protected.Group("/menu")
	menu.GET("/active-now", resource.ServersControllers.SourceMenu.ServiceGetActiveMenu)
	menu.GET("/active", resource.ServersControllers.SourceMenu.ServiceListActiveMenus)
	menu.GET("/options", resource.ServersControllers.SourceMenu.ServiceGetMenuOptions)
	menu.GET("/:id", resource.ServersControllers.SourceMenu.ServiceGetMenu)
	menu.GET("", resource.ServersControllers.SourceMenu.ServiceListMenus)

	// Category (leitura)
	category := protected.Group("/category")
	category.GET("/:id", resource.ServersControllers.SourceCategory.ServiceGetCategory)
	category.GET("", resource.ServersControllers.SourceCategory.ServiceListCategories)
	category.GET("/active", resource.ServersControllers.SourceCategory.ServiceListActiveCategories)
	category.GET("/menu/:menuId", resource.ServersControllers.SourceCategory.ServiceGetCategoriesByMenu)

	// Subcategory (leitura)
	subcategory := protected.Group("/subcategory")
	subcategory.GET("/:id", resource.ServersControllers.SourceSubcategory.ServiceGetSubcategory)
	subcategory.GET("", resource.ServersControllers.SourceSubcategory.ServiceListSubcategories)
	subcategory.GET("/active", resource.ServersControllers.SourceSubcategory.ServiceListActiveSubcategories)
	subcategory.GET("/category/:categoryId", resource.ServersControllers.SourceSubcategory.ServiceGetSubcategoriesByCategory)
	subcategory.GET("/:id/categories", resource.ServersControllers.SourceSubcategory.ServiceGetSubcategoryCategories)

	// Role (leitura + atribuição)
	role := protected.Group("/role")
	role.GET("/system", resource.ServersControllers.SourceRole.ListSystemRoles)
	role.GET("/check", resource.ServersControllers.SourceRole.CheckPermission)
	role.GET("/my-permissions", resource.ServersControllers.SourceRole.GetMyPermissions)
	role.POST("/assign", resource.ServersControllers.SourceRole.AssignRoleToUser)
	role.POST("/remove", resource.ServersControllers.SourceRole.RemoveRoleFromUser)
	role.GET("/user/:userId", resource.ServersControllers.SourceRole.GetUserRoles)
	role.GET("/user/:userId/details", resource.ServersControllers.SourceRole.GetUserRolesWithDetails)
	role.POST("/permission", resource.ServersControllers.SourceRole.AddPermissionToRole)
	role.GET("", resource.ServersControllers.SourceRole.ListRoles)
	role.GET("/:id", resource.ServersControllers.SourceRole.GetRole)
	role.GET("/:id/permissions", resource.ServersControllers.SourceRole.GetRolePermissions)

	// Module (leitura)
	module := protected.Group("/module")
	module.GET("", resource.ServersControllers.SourceRole.ListModules)
	module.GET("/with-permissions", resource.ServersControllers.SourceRole.ListModulesWithPermissions)
	module.GET("/available", resource.ServersControllers.SourceRole.GetOrganizationModules)
	module.GET("/:id", resource.ServersControllers.SourceRole.GetModule)

	// Permission (leitura)
	permission := protected.Group("/permission")
	permission.GET("", resource.ServersControllers.SourceRole.ListPermissions)
	permission.GET("/:id", resource.ServersControllers.SourceRole.GetPermission)

	// Package (leitura + assinatura)
	pkg := protected.Group("/package")
	pkg.GET("/subscription", resource.ServersControllers.SourceRole.GetOrganizationSubscription)
	pkg.GET("/usage-limits", resource.ServersControllers.SourceRole.GetUsageAndLimits)
	pkg.POST("/subscribe", resource.ServersControllers.SourceRole.SubscribeOrganization)
	pkg.GET("", resource.ServersControllers.SourceRole.ListPackages)
	pkg.GET("/:id", resource.ServersControllers.SourceRole.GetPackageWithModules)
	pkg.GET("/:id/limits", resource.ServersControllers.SourceRole.GetPackageLimits)

	// Organization
	organization := protected.Group("/organization")
	organization.GET("/:id", resource.ServersControllers.SourceOrganization.GetOrganizationById)
	organization.GET("", resource.ServersControllers.SourceOrganization.ListOrganizations)
	organization.GET("/active", resource.ServersControllers.SourceOrganization.ListActiveOrganizations)
	organization.GET("/email", resource.ServersControllers.SourceOrganization.GetOrganizationByEmail)
	organization.PUT("/:id", resource.ServersControllers.SourceOrganization.UpdateOrganization)
	organization.DELETE("/:id", resource.ServersControllers.SourceOrganization.SoftDeleteOrganization)
	organization.DELETE("/:id/permanent", resource.ServersControllers.SourceOrganization.HardDeleteOrganization)

	// Project
	project := protected.Group("/project")
	project.GET("/:id", resource.ServersControllers.SourceProject.GetProjectById)
	project.GET("", resource.ServersControllers.SourceProject.GetProjectsByOrganization)
	project.GET("/organization/:orgId", resource.ServersControllers.SourceProject.GetProjectsByOrganizationId)
	project.GET("/active", resource.ServersControllers.SourceProject.GetActiveProjects)
	project.PUT("/:id", resource.ServersControllers.SourceProject.UpdateProject)
	project.PUT("/:id/set-default", resource.ServersControllers.SourceProject.SetDefaultProject)
	project.DELETE("/:id", resource.ServersControllers.SourceProject.SoftDeleteProject)
	project.DELETE("/:id/permanent", resource.ServersControllers.SourceProject.HardDeleteProject)

	// Settings
	settings := protected.Group("/settings")
	settings.GET("", resource.ServersControllers.SourceSettings.GetSettingsByProject)
	settings.PUT("", resource.ServersControllers.SourceSettings.UpdateSettings)

	// Display Settings
	displaySettings := protected.Group("/project/settings/display")
	displaySettings.GET("", resource.ServersControllers.SourceDisplaySettings.GetDisplaySettings)
	displaySettings.PUT("", resource.ServersControllers.SourceDisplaySettings.UpdateDisplaySettings)
	displaySettings.POST("/reset", resource.ServersControllers.SourceDisplaySettings.ResetDisplaySettings)

	// Theme Customization
	theme := protected.Group("/project/settings/theme")
	theme.GET("", resource.ServersControllers.SourceThemeCustomization.GetTheme)
	theme.POST("", resource.ServersControllers.SourceThemeCustomization.CreateOrUpdateTheme)
	theme.PUT("", resource.ServersControllers.SourceThemeCustomization.CreateOrUpdateTheme)
	theme.POST("/reset", resource.ServersControllers.SourceThemeCustomization.ResetTheme)
	theme.DELETE("", resource.ServersControllers.SourceThemeCustomization.DeleteTheme)

	// Sidebar Config (leitura - escrita é admin)
	sidebarConfig := protected.Group("/sidebar-config")
	sidebarConfig.GET("", resource.ServersControllers.SourceSidebarConfig.GetConfig)

	// Reports
	reports := protected.Group("/reports")
	reports.GET("/occupancy", resource.ServersControllers.SourceReports.GetOccupancyReport)
	reports.GET("/reservations", resource.ServersControllers.SourceReports.GetReservationReport)
	reports.GET("/waitlist", resource.ServersControllers.SourceReports.GetWaitlistReport)
	reports.GET("/leads", resource.ServersControllers.SourceReports.GetLeadReport)
	reports.GET("/export/:type", resource.ServersControllers.SourceReports.ExportReportToCSV)

	// Notification
	notification := protected.Group("/notification")
	notification.POST("/send", resource.ServersControllers.SourceNotification.SendNotification)
	notification.POST("/event", resource.ServersControllers.SourceNotification.ProcessEvent)
	notification.GET("/logs/:orgId/:projectId", resource.ServersControllers.SourceNotification.GetNotificationLogs)
	notification.GET("/templates/:orgId/:projectId", resource.ServersControllers.SourceNotification.GetNotificationTemplates)
	notification.POST("/template", resource.ServersControllers.SourceNotification.CreateNotificationTemplate)
	notification.PUT("/template", resource.ServersControllers.SourceNotification.UpdateNotificationTemplate)
	notification.POST("/config", resource.ServersControllers.SourceNotification.CreateOrUpdateNotificationConfig)
	notification.GET("/review-queue/:orgId/:projectId", resource.ServersControllers.SourceNotification.GetReviewQueue)
	notification.POST("/review-queue/:id/approve", resource.ServersControllers.SourceNotification.ApproveReviewItem)
	notification.POST("/review-queue/:id/reject", resource.ServersControllers.SourceNotification.RejectReviewItem)
	notification.POST("/review-queue/:id/custom", resource.ServersControllers.SourceNotification.ExecuteCustomAction)
	notification.GET("/reminders/:orgId/:projectId", resource.ServersControllers.SourceNotification.GetNotificationReminders)
	notification.POST("/reminder", resource.ServersControllers.SourceNotification.CreateNotificationReminder)
	notification.PUT("/reminder", resource.ServersControllers.SourceNotification.UpdateNotificationReminder)
	notification.DELETE("/reminder/:id", resource.ServersControllers.SourceNotification.DeleteNotificationReminder)
	notification.POST("/trigger-scheduled", resource.ServersControllers.SourceNotification.TriggerScheduledNotifications)

	// Client Audit Logs
	clientAuditLogs := protected.Group("/client-audit-logs")
	clientAuditLogs.GET("", resource.ServersControllers.SourceClientAuditLog.ServiceListClientAuditLogs)
	clientAuditLogs.GET("/:id", resource.ServersControllers.SourceClientAuditLog.ServiceGetClientAuditLog)

	// Plan Change Requests (client)
	planChange := protected.Group("/plan-change-request")
	planChange.POST("", resource.ServersControllers.SourcePlanChangeRequest.CreateRequest)
	planChange.GET("/my-requests", resource.ServersControllers.SourcePlanChangeRequest.GetMyRequests)
	planChange.GET("/:id", resource.ServersControllers.SourcePlanChangeRequest.GetRequestById)
	planChange.POST("/:id/cancel", resource.ServersControllers.SourcePlanChangeRequest.CancelRequest)

	// Onboarding
	onboarding := protected.Group("/onboarding")
	onboarding.GET("/status", resource.ServersControllers.SourceOnboarding.GetOnboardingStatus)

	// Environment
	environment := protected.Group("/environment")
	environment.GET("/:id", resource.ServersControllers.SourceEnvironment.GetEnvironmentById)
	environment.GET("", resource.ServersControllers.SourceEnvironment.GetEnvironmentsByProject)
	environment.GET("/active", resource.ServersControllers.SourceEnvironment.GetActiveEnvironments)
	environment.POST("", resource.ServersControllers.SourceEnvironment.CreateEnvironment)
	environment.PUT("/:id", resource.ServersControllers.SourceEnvironment.UpdateEnvironment)
	environment.DELETE("/:id", resource.ServersControllers.SourceEnvironment.SoftDeleteEnvironment)

	// =============================================================================
	// 5. ROTAS COM RolePermissionMiddleware
	// =============================================================================

	// Product
	product := protected.Group("/product")
	product.GET("/:id", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1), resource.ServersControllers.SourceProducts.ServiceGetProduct)
	product.GET("/purchase/:id", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1), resource.ServersControllers.SourceProducts.ServiceGetProductByPurchase)
	product.GET("", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1), resource.ServersControllers.SourceProducts.ServiceListProducts)
	product.GET("/by-tag", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1), resource.ServersControllers.SourceProducts.ServiceGetProductsByTag)
	product.POST("", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_create", 1), middleware.PackageLimitMiddleware(resource.Handlers.HandlerLimits, handler.LimitProducts), resource.ServersControllers.SourceProducts.ServiceCreateProduct)
	product.PUT("/:id", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_edit", 1), resource.ServersControllers.SourceProducts.ServiceUpdateProduct)
	product.PUT("/:id/image", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_edit", 1), resource.ServersControllers.SourceProducts.ServiceUpdateProductImage)
	product.DELETE("/:id", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_delete", 1), resource.ServersControllers.SourceProducts.ServiceDeleteProduct)
	product.GET("/:id/tags", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1), resource.ServersControllers.SourceProducts.ServiceGetProductTags)
	product.POST("/:id/tags", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_edit", 1), resource.ServersControllers.SourceProducts.ServiceAddTagToProduct)
	product.DELETE("/:id/tags/:tagId", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_edit", 1), resource.ServersControllers.SourceProducts.ServiceRemoveTagFromProduct)
	product.PUT("/:id/order", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_edit", 1), resource.ServersControllers.SourceProducts.ServiceUpdateProductOrder)
	product.PUT("/:id/status", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_edit", 1), resource.ServersControllers.SourceProducts.ServiceUpdateProductStatus)
	product.GET("/type/:type", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1), resource.ServersControllers.SourceProducts.ServiceGetProductsByType)
	product.GET("/category/:categoryId", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1), resource.ServersControllers.SourceProducts.ServiceGetProductsByCategory)
	product.GET("/subcategory/:subcategoryId", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_products_view", 1), resource.ServersControllers.SourceProducts.ServiceGetProductsBySubcategory)

	// Table
	table := protected.Group("/table")
	table.GET("/:id", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tables_view", 1), resource.ServersControllers.SourceTables.ServiceGetTable)
	table.GET("", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tables_view", 1), resource.ServersControllers.SourceTables.ServiceListTables)
	table.POST("", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tables_create", 1), middleware.PackageLimitMiddleware(resource.Handlers.HandlerLimits, handler.LimitTables), resource.ServersControllers.SourceTables.ServiceCreateTable)
	table.PUT("/:id", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tables_edit", 1), resource.ServersControllers.SourceTables.ServiceUpdateTable)
	table.DELETE("/:id", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tables_delete", 1), resource.ServersControllers.SourceTables.ServiceDeleteTable)

	// Reservation (requer módulo)
	reservation := protected.Group("/reservation")
	reservation.Use(middleware.ModuleRequiredMiddleware(resource.Handlers.HandlerLimits, "client_reservations"))
	reservation.GET("/:id", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_reservations_view", 1), resource.ServersControllers.SourceReservation.ServiceGetReservation)
	reservation.GET("", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_reservations_view", 1), resource.ServersControllers.SourceReservation.ServiceListReservations)
	reservation.POST("", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_reservations_create", 1), middleware.PackageLimitMiddleware(resource.Handlers.HandlerLimits, handler.LimitReservationsDay), resource.ServersControllers.SourceReservation.ServiceCreateReservation)
	reservation.PUT("/:id", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_reservations_edit", 1), resource.ServersControllers.SourceReservation.ServiceUpdateReservation)
	reservation.DELETE("/:id", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_reservations_delete", 1), resource.ServersControllers.SourceReservation.ServiceDeleteReservation)

	// Customer
	customer := protected.Group("/customer")
	customer.GET("/:id", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_customers_view", 1), resource.ServersControllers.SourceCustomer.ServiceGetCustomer)
	customer.GET("", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_customers_view", 1), resource.ServersControllers.SourceCustomer.ServiceListCustomers)
	customer.POST("", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_customers_create", 1), resource.ServersControllers.SourceCustomer.ServiceCreateCustomer)
	customer.PUT("/:id", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_customers_edit", 1), resource.ServersControllers.SourceCustomer.ServiceUpdateCustomer)
	customer.DELETE("/:id", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_customers_delete", 1), resource.ServersControllers.SourceCustomer.ServiceDeleteCustomer)

	// Order
	order := protected.Group("/order")
	order.GET("/:id", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_orders_view", 1), resource.ServersControllers.SourceOrders.GetOrderById)
	order.GET("/:id/progress", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_orders_view", 1), resource.ServersControllers.SourceOrders.GetOrderProgress)
	order.GET("", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_orders_view", 1), resource.ServersControllers.SourceOrders.ListOrders)
	order.POST("", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_orders_create", 1), resource.ServersControllers.SourceOrders.CreateOrder)
	order.PUT("/:id", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_orders_edit", 1), resource.ServersControllers.SourceOrders.UpdateOrder)
	order.PUT("/:id/status", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_orders_edit", 1), resource.ServersControllers.SourceOrders.UpdateOrderStatus)
	order.DELETE("/:id", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_orders_delete", 1), resource.ServersControllers.SourceOrders.SoftDeleteOrder)

	// Kitchen
	kitchen := protected.Group("/kitchen")
	kitchen.GET("/queue", resource.ServersControllers.SourceOrders.GetKitchenQueue)

	// Waitlist (requer módulo)
	waitlist := protected.Group("/waitlist")
	waitlist.Use(middleware.ModuleRequiredMiddleware(resource.Handlers.HandlerLimits, "client_waitlist"))
	waitlist.GET("/:id", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_waitlist_view", 1), resource.ServersControllers.SourceWaitlist.ServiceGetWaitlist)
	waitlist.GET("", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_waitlist_view", 1), resource.ServersControllers.SourceWaitlist.ServiceListWaitlists)
	waitlist.POST("", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_waitlist_create", 1), resource.ServersControllers.SourceWaitlist.ServiceCreateWaitlist)
	waitlist.PUT("/:id", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_waitlist_edit", 1), resource.ServersControllers.SourceWaitlist.ServiceUpdateWaitlist)
	waitlist.DELETE("/:id", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_waitlist_delete", 1), resource.ServersControllers.SourceWaitlist.ServiceDeleteWaitlist)

	// Tag
	tag := protected.Group("/tag")
	tag.GET("/:id", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tags_view", 1), resource.ServersControllers.SourceTag.ServiceGetTag)
	tag.GET("", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tags_view", 1), resource.ServersControllers.SourceTag.ServiceListTags)
	tag.GET("/active", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tags_view", 1), resource.ServersControllers.SourceTag.ServiceListActiveTags)
	tag.GET("/entity/:entityType", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tags_view", 1), resource.ServersControllers.SourceTag.ServiceGetTagsByEntityType)
	tag.POST("", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tags_create", 1), resource.ServersControllers.SourceTag.ServiceCreateTag)
	tag.PUT("/:id", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tags_edit", 1), resource.ServersControllers.SourceTag.ServiceUpdateTag)
	tag.DELETE("/:id", middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_tags_delete", 1), resource.ServersControllers.SourceTag.ServiceDeleteTag)
}
