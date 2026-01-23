package public

import (
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupMenuRoutes configura rotas públicas de menu e reserva
func SetupMenuRoutes(r *gin.Engine) {
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
