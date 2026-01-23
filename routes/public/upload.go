package public

import (
	"lep/middleware"
	"lep/resource"

	"github.com/gin-gonic/gin"
)

// SetupUploadRoutes configura rotas de upload
func SetupUploadRoutes(r *gin.Engine) {
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
