package client

import (
	"lep/handler"
	"lep/middleware"
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupReservationRoutes configura rotas de reserva para client
// Requer módulo client_reservations + RolePermissionMiddleware
func SetupReservationRoutes(r gin.IRouter) {
	reservation := r.Group("/reservation")
	{
		// Verificar se módulo de reservas está disponível
		reservation.Use(middleware.ModuleRequiredMiddleware(resource.Handlers.HandlerLimits, "client_reservations"))

		// Rotas de leitura
		reservation.GET("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_reservations_view", 1),
			resource.ServersControllers.SourceReservation.ServiceGetReservation)
		reservation.GET("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_reservations_view", 1),
			resource.ServersControllers.SourceReservation.ServiceListReservations)

		// Rotas de escrita + limite de reservas por dia
		reservation.POST("",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_reservations_create", 1),
			middleware.PackageLimitMiddleware(resource.Handlers.HandlerLimits, handler.LimitReservationsDay),
			resource.ServersControllers.SourceReservation.ServiceCreateReservation)
		reservation.PUT("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_reservations_edit", 1),
			resource.ServersControllers.SourceReservation.ServiceUpdateReservation)
		reservation.DELETE("/:id",
			middleware.RolePermissionMiddleware(resource.Handlers.HandlerRole, "client_reservations_delete", 1),
			resource.ServersControllers.SourceReservation.ServiceDeleteReservation)
	}
}
