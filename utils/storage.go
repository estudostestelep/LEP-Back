package utils

import (
	"context"
	"fmt"
	"io"
	"lep/config"
	"lep/constants"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

type StorageService interface {
	UploadFile(file multipart.File, header *multipart.FileHeader, orgId, projId, category string) (*UploadResult, error)
	GetPublicURL(filename, orgId, projId, category string) string
	DeleteFile(filename, orgId, projId, category string) error
}

type UploadResult struct {
	Filename  string
	PublicURL string
	Size      int64
}

type LocalStorageService struct {
	baseDir string
	baseURL string
}

type GCSStorageService struct {
	client     *storage.Client
	bucketName string
	baseURL    string
}

// NewStorageService cria uma instância do serviço de storage baseado no ambiente
func NewStorageService() StorageService {
	if config.IsGCSStorage() {
		// Ambientes stage/prod - usar Cloud Storage

		// Usar BUCKET_NAME se disponível, senão fallback para STORAGE_BUCKET_NAME
		bucketName := config.BUCKET_NAME
		if bucketName == "" {
			bucketName = config.STORAGE_BUCKET_NAME
		}

		if bucketName == "" {
			// Fallback para storage local se bucket não estiver configurado
			return NewLocalStorageService()
		}

		client, err := storage.NewClient(context.Background())
		if err != nil {
			// Fallback para storage local se não conseguir conectar ao GCS
			return NewLocalStorageService()
		}

		return &GCSStorageService{
			client:     client,
			bucketName: bucketName,
			baseURL:    config.BASE_URL,
		}
	}

	// Ambiente dev - usar storage local
	return NewLocalStorageService()
}

// NewLocalStorageService cria uma instância do serviço de storage local
func NewLocalStorageService() *LocalStorageService {
	return &LocalStorageService{
		baseDir: "./uploads",
		baseURL: config.BASE_URL,
	}
}

// LocalStorageService implementation
func (s *LocalStorageService) UploadFile(file multipart.File, header *multipart.FileHeader, orgId, projId, category string) (*UploadResult, error) {
	// Criar estrutura de diretórios orgId/projId/category
	organizationDir := filepath.Join(s.baseDir, orgId, projId, category)
	if err := os.MkdirAll(organizationDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("error creating directory structure: %w", err)
	}

	// Gerar nome único
	fileExt := getFileExtension(header.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), fileExt)
	filePath := filepath.Join(organizationDir, filename)

	// Criar arquivo no servidor
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("error creating file: %w", err)
	}
	defer dst.Close()

	// Copiar dados
	size, err := io.Copy(dst, file)
	if err != nil {
		return nil, fmt.Errorf("error saving file: %w", err)
	}

	// URL pública com nova estrutura
	publicURL := fmt.Sprintf("%s/uploads/%s/%s/%s/%s", s.baseURL, orgId, projId, category, filename)

	// Nome do arquivo retornado inclui path completo para compatibilidade
	fullFilename := fmt.Sprintf("%s/%s/%s/%s", orgId, projId, category, filename)

	return &UploadResult{
		Filename:  fullFilename,
		PublicURL: publicURL,
		Size:      size,
	}, nil
}

func (s *LocalStorageService) GetPublicURL(filename, orgId, projId, category string) string {
	// Se filename já inclui o path completo, usa direto
	if strings.Contains(filename, "/") {
		return fmt.Sprintf("%s/uploads/%s", s.baseURL, filename)
	}
	// Senão, constrói o path completo
	return fmt.Sprintf("%s/uploads/%s/%s/%s/%s", s.baseURL, orgId, projId, category, filename)
}

func (s *LocalStorageService) DeleteFile(filename, orgId, projId, category string) error {
	// Se filename já inclui o path completo, usa direto
	var filePath string
	if strings.Contains(filename, "/") {
		filePath = filepath.Join(s.baseDir, filename)
	} else {
		filePath = filepath.Join(s.baseDir, orgId, projId, category, filename)
	}
	return os.Remove(filePath)
}

// GCSStorageService implementation
func (s *GCSStorageService) UploadFile(file multipart.File, header *multipart.FileHeader, orgId, projId, category string) (*UploadResult, error) {
	// Usar timeout configurado ou fallback
	timeout := time.Duration(config.BUCKET_TIMEOUT) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second // fallback: 30 segundos
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Gerar nome único com estrutura orgId/projId/category
	fileExt := getFileExtension(header.Filename)
	objectName := fmt.Sprintf("%s/%s/%s/%s%s", orgId, projId, category, uuid.New().String(), fileExt)

	// Obter objeto do bucket
	obj := s.client.Bucket(s.bucketName).Object(objectName)

	// Criar writer
	writer := obj.NewWriter(ctx)

	// Definir metadados
	writer.ContentType = header.Header.Get("Content-Type")

	// Usar configuração de cache control ou fallback
	cacheControl := config.BUCKET_CACHE_CONTROL
	if cacheControl == "" {
		cacheControl = "public, max-age=604800" // fallback: 1 semana
	}
	writer.CacheControl = cacheControl

	// Copiar dados
	size, err := io.Copy(writer, file)
	if err != nil {
		writer.Close() // Fechar writer em caso de erro
		return nil, fmt.Errorf("error uploading to GCS: %w", err)
	}

	// IMPORTANTE: Fechar writer primeiro para finalizar o upload
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("error finalizing upload: %w", err)
	}

	// AGORA definir ACL público para leitura (após o objeto existir)
	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		// Log do erro, mas não falha o upload
		fmt.Printf("Warning: Could not set public ACL for %s: %v\n", objectName, err)
	}

	publicURL := fmt.Sprintf("%s/%s", s.baseURL, objectName)

	return &UploadResult{
		Filename:  objectName,
		PublicURL: publicURL,
		Size:      size,
	}, nil
}

func (s *GCSStorageService) GetPublicURL(filename, orgId, projId, category string) string {
	// Se filename já inclui o path completo, usa direto
	if strings.Contains(filename, "/") {
		return fmt.Sprintf("%s/%s", s.baseURL, filename)
	}

	// Senão, constrói o path completo
	return fmt.Sprintf("%s/%s/%s/%s/%s", s.baseURL, orgId, projId, category, filename)
}

func (s *GCSStorageService) DeleteFile(filename, orgId, projId, category string) error {
	ctx := context.Background()

	// Se filename já inclui o path completo, usa direto
	objectName := filename
	if !strings.Contains(filename, "/") {
		objectName = fmt.Sprintf("%s/%s/%s/%s", orgId, projId, category, filename)
	}

	obj := s.client.Bucket(s.bucketName).Object(objectName)
	return obj.Delete(ctx)
}

// Função auxiliar para obter extensão do arquivo
func getFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return ".jpg" // default
	}
	return strings.ToLower(ext)
}

// Validar se o tipo MIME é de uma imagem válida
func IsValidImageType(contentType string) bool {
	return constants.IsValidImageMimeType(contentType)
}

// Validar tamanho do arquivo
func IsValidFileSize(size int64, maxSize int64) bool {
	return size <= maxSize
}
