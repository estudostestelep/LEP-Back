package public

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupMenuRoutes configura rotas públicas de menu e reserva
func SetupMenuRoutes(r *gin.Engine) {
	publicRoutes := r.Group("/public")
	{
		// === Rotas com UUIDs (compatibilidade) ===
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

		// === Rotas com Slugs (novas - URLs amigáveis) ===
		// Resolver projeto por slugs
		publicRoutes.GET("/project/resolve", resource.ServersControllers.SourcePublic.ServiceResolveProject)

		// Cardápio por slug (projeto default se não especificado)
		publicRoutes.GET("/menu/org/:orgSlug", resource.ServersControllers.SourcePublic.ServiceGetPublicMenuBySlug)
		publicRoutes.GET("/menu/org/:orgSlug/:projectSlug", resource.ServersControllers.SourcePublic.ServiceGetPublicMenuBySlug)

		// Categorias por slug
		publicRoutes.GET("/categories/org/:orgSlug", resource.ServersControllers.SourcePublic.ServiceGetPublicCategoriesBySlug)
		publicRoutes.GET("/categories/org/:orgSlug/:projectSlug", resource.ServersControllers.SourcePublic.ServiceGetPublicCategoriesBySlug)

		// Menus por slug
		publicRoutes.GET("/menus/org/:orgSlug", resource.ServersControllers.SourcePublic.ServiceGetPublicMenusBySlug)
		publicRoutes.GET("/menus/org/:orgSlug/:projectSlug", resource.ServersControllers.SourcePublic.ServiceGetPublicMenusBySlug)

		// Info do projeto por slug
		publicRoutes.GET("/project/org/:orgSlug", resource.ServersControllers.SourcePublic.ServiceGetProjectInfoBySlug)
		publicRoutes.GET("/project/org/:orgSlug/:projectSlug", resource.ServersControllers.SourcePublic.ServiceGetProjectInfoBySlug)

		// Horários por slug
		publicRoutes.GET("/times/org/:orgSlug", resource.ServersControllers.SourcePublic.ServiceGetAvailableTimesBySlug)
		publicRoutes.GET("/times/org/:orgSlug/:projectSlug", resource.ServersControllers.SourcePublic.ServiceGetAvailableTimesBySlug)

		// Reserva por slug
		publicRoutes.POST("/reservation/org/:orgSlug", resource.ServersControllers.SourcePublic.ServiceCreatePublicReservationBySlug)
		publicRoutes.POST("/reservation/org/:orgSlug/:projectSlug", resource.ServersControllers.SourcePublic.ServiceCreatePublicReservationBySlug)
	}
}
