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
}

func setupUserRoutes(r *gin.Engine) {
	userRoutes := r.Group("/user")
	{
		userRoutes.GET("/:id", resource.ServersControllers.SourceUsers.ServiceGetUser)
		userRoutes.GET("/group/:id", resource.ServersControllers.SourceUsers.ServiceGetUserByGroup)
		userRoutes.PUT("/:id", resource.ServersControllers.SourceUsers.ServiceUpdateUser)
		userRoutes.DELETE("/:id", resource.ServersControllers.SourceUsers.ServiceDeleteUser)
	}
}

func setupProductRoutes(r *gin.Engine) {
	productRoutes := r.Group("/product")
	{
		productRoutes.GET("/:id", resource.ServersControllers.SourceProducts.ServiceGetProduct)
		productRoutes.GET("/purchase/:id", resource.ServersControllers.SourceProducts.ServiceGetProductByPurchase)
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
		orderRoutes.POST("", resource.ServersControllers.SourceOrders.CreateOrder)
		orderRoutes.PUT("/:id", resource.ServersControllers.SourceOrders.UpdateOrder)
		orderRoutes.DELETE("/:id", resource.ServersControllers.SourceOrders.SoftDeleteOrder)
	}
	// List orders endpoint (plural)
	r.GET("/orders", resource.ServersControllers.SourceOrders.ListOrders)
}
