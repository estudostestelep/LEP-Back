package server

import (
	"fmt"
	"lep/config"
	"lep/constants"
	"lep/utils"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ResourceUpload struct {
	storageService  utils.StorageService
	fallbackStorage utils.StorageService
}

type IServerUpload interface {
	ServiceUploadImage(c *gin.Context)        // Método genérico para qualquer categoria
	ServiceUploadProductImage(c *gin.Context) // Mantido para retrocompatibilidade
	ServiceGetUploadedFile(c *gin.Context)
}

// ServiceUploadImage faz upload de imagem genérico para qualquer categoria
func (r *ResourceUpload) ServiceUploadImage(c *gin.Context) {
	// Headers validados pelo middleware - acessar via context
	organizationId := c.GetString("organization_id")
	projectId := c.GetString("project_id")

	if organizationId == "" || projectId == "" {
		utils.SendBadRequestError(c, "Organization ID and Project ID are required", nil)
		return
	}

	// Obter categoria da URL
	category := c.Param("category")
	if !constants.IsValidUploadCategory(category) {
		utils.SendBadRequestError(c, constants.ErrorInvalidCategory, nil)
		return
	}

	// Parse do form multipart
	err := c.Request.ParseMultipartForm(constants.MaxFormSize)
	if err != nil {
		utils.SendBadRequestError(c, constants.ErrorParsingForm, err)
		return
	}

	// Obter o arquivo
	file, handler, err := c.Request.FormFile("image")
	if err != nil {
		utils.SendBadRequestError(c, constants.ErrorNoFile, err)
		return
	}
	defer file.Close()

	// Validações
	contentType := handler.Header.Get("Content-Type")
	if !utils.IsValidImageType(contentType) {
		utils.SendBadRequestError(c, constants.ErrorInvalidFileType, nil)
		return
	}

	if !utils.IsValidFileSize(handler.Size, constants.MaxUploadSize) {
		utils.SendBadRequestError(c, constants.ErrorFileTooLarge, nil)
		return
	}

	// Upload usando serviço híbrido (local ou GCS) com organização por tenant
	result, err := r.uploadWithFallback(file, handler, organizationId, projectId, category)
	if err != nil {
		utils.SendInternalServerError(c, constants.ErrorUploadingFile, err)
		return
	}

	// Retornar resposta de sucesso
	response := gin.H{
		"success":         true,
		"image_url":       result.PublicURL,
		"filename":        result.Filename,
		"size":            result.Size,
		"category":        category,
		"organization_id": organizationId,
		"project_id":      projectId,
	}

	utils.SendCreatedSuccess(c, constants.SuccessUploadMessage, response)
}

// ServiceUploadProductImage faz upload de imagem de produto (retrocompatibilidade)
func (r *ResourceUpload) ServiceUploadProductImage(c *gin.Context) {
	// Simular parâmetro de categoria para reutilizar o método genérico
	c.Params = append(c.Params, gin.Param{Key: "category", Value: constants.UploadCategoryProducts})

	// Chamar método genérico
	r.ServiceUploadImage(c)
}

// ServiceGetUploadedFile serve arquivos uploadados estaticamente
func (r *ResourceUpload) ServiceGetUploadedFile(c *gin.Context) {
	// Suporte para nova estrutura: /uploads/orgId/projId/category/filename
	orgId := c.Param("orgId")
	projId := c.Param("projId")
	category := c.Param("category")
	filename := c.Param("filename")

	// Validar parâmetros obrigatórios
	if orgId == "" || projId == "" || category == "" || filename == "" {
		// Fallback para estrutura antiga para compatibilidade
		category = c.Param("category")
		filename = c.Param("filename")

		if category == "" || filename == "" {
			utils.SendBadRequestError(c, "Invalid file path", nil)
			return
		}

		// Construir caminho do arquivo (estrutura antiga)
		filePath := filepath.Join("uploads", category, filename)

		// Verificar se arquivo existe
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			utils.SendNotFoundError(c, "File")
			return
		}

		r.serveStaticFile(c, filePath, filename)
		return
	}

	// Validar categoria
	if !constants.IsValidUploadCategory(category) {
		utils.SendNotFoundError(c, "Category")
		return
	}

	// Construir caminho do arquivo (nova estrutura)
	filePath := filepath.Join("uploads", orgId, projId, category, filename)

	// Verificar se arquivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.SendNotFoundError(c, "File")
		return
	}

	r.serveStaticFile(c, filePath, filename)
}

// serveStaticFile helper para servir arquivos estáticos
func (r *ResourceUpload) serveStaticFile(c *gin.Context, filePath, filename string) {
	// Determinar Content-Type baseado na extensão
	ext := strings.ToLower(filepath.Ext(filename))
	contentType := "application/octet-stream"
	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".webp":
		contentType = "image/webp"
	case ".gif":
		contentType = "image/gif"
	}

	// Definir headers de cache (1 semana)
	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=604800")
	c.Header("Last-Modified", time.Now().Format(http.TimeFormat))

	// Servir arquivo
	c.File(filePath)
}

// uploadWithFallback tenta upload principal e fallback em caso de erro
func (r *ResourceUpload) uploadWithFallback(file multipart.File, handler *multipart.FileHeader, orgId, projId, category string) (*utils.UploadResult, error) {
	// Tentar storage principal
	result, err := r.storageService.UploadFile(file, handler, orgId, projId, category)
	if err != nil {
		// Se falhou e tem fallback disponível, tentar fallback
		if r.fallbackStorage != nil {
			// Reset file pointer para início
			if _, seekErr := file.Seek(0, 0); seekErr != nil {
				return nil, fmt.Errorf("error resetting file pointer: %w", seekErr)
			}

			result, err = r.fallbackStorage.UploadFile(file, handler, orgId, projId, category)
			if err != nil {
				return nil, fmt.Errorf("both primary and fallback storage failed: %w", err)
			}
		} else {
			return nil, err
		}
	}

	return result, nil
}

// isValidImageType verifica se o tipo MIME é de uma imagem válida
func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/webp",
		"image/gif",
	}

	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}

// getFileExtension extrai a extensão do nome do arquivo
func getFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return ".jpg" // default
	}
	return strings.ToLower(ext)
}

// getBaseURL constrói a URL base do servidor
func getBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	host := c.Request.Host
	if host == "" {
		host = "localhost:8080" // fallback para desenvolvimento
	}

	return fmt.Sprintf("%s://%s", scheme, host)
}

// NewSourceServerUpload cria nova instância do controller de upload
func NewSourceServerUpload() IServerUpload {
	primaryStorage := utils.NewStorageService()
	var fallbackStorage utils.StorageService

	// Se o storage principal for GCS, configurar local como fallback
	if config.IsGCSStorage() {
		fallbackStorage = utils.NewLocalStorageService()
	}

	return &ResourceUpload{
		storageService:  primaryStorage,
		fallbackStorage: fallbackStorage,
	}
}
