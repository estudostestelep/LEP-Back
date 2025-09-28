package constants

// Upload Categories - Categorias permitidas para upload de arquivos
const (
	// UploadCategoryProducts representa uploads de imagens de produtos do cardápio
	UploadCategoryProducts = "products"

	// UploadCategoryUsers representa uploads de avatars de usuários/funcionários
	UploadCategoryUsers = "users"

	// UploadCategoryRestaurants representa uploads de logos e imagens do restaurante
	UploadCategoryRestaurants = "restaurants"

	// UploadCategoryBanners representa uploads de banners promocionais
	UploadCategoryBanners = "banners"
)

// Upload Limits - Limitações de upload
const (
	// MaxUploadSize representa o tamanho máximo de arquivo para upload (5MB)
	MaxUploadSize = 5 * 1024 * 1024

	// MaxFormSize representa o tamanho máximo do form multipart (10MB)
	MaxFormSize = 10 << 20
)

// Valid Image MIME Types - Tipos MIME válidos para upload de imagens
const (
	// MimeTypeJPEG representa o tipo MIME para imagens JPEG
	MimeTypeJPEG = "image/jpeg"

	// MimeTypeJPG representa o tipo MIME alternativo para imagens JPEG
	MimeTypeJPG = "image/jpg"

	// MimeTypePNG representa o tipo MIME para imagens PNG
	MimeTypePNG = "image/png"

	// MimeTypeWebP representa o tipo MIME para imagens WebP
	MimeTypeWebP = "image/webp"

	// MimeTypeGIF representa o tipo MIME para imagens GIF
	MimeTypeGIF = "image/gif"
)

// Upload Validation Messages - Mensagens de validação para upload
const (
	// ErrorInvalidCategory mensagem para categoria inválida
	ErrorInvalidCategory = "Invalid category. Allowed: products, users, restaurants, banners"

	// ErrorNoFile mensagem quando nenhum arquivo é enviado
	ErrorNoFile = "No image file provided"

	// ErrorInvalidFileType mensagem para tipo de arquivo inválido
	ErrorInvalidFileType = "Invalid file type. Only JPEG, PNG, WebP and GIF are allowed"

	// ErrorFileTooLarge mensagem para arquivo muito grande
	ErrorFileTooLarge = "File too large. Maximum size is 5MB"

	// ErrorParsingForm mensagem para erro no parsing do form
	ErrorParsingForm = "Error parsing multipart form"

	// ErrorUploadingFile mensagem para erro no upload
	ErrorUploadingFile = "Error uploading file"

	// SuccessUploadMessage mensagem de sucesso no upload
	SuccessUploadMessage = "Image uploaded successfully"
)

// GetValidUploadCategories retorna slice com todas as categorias válidas
func GetValidUploadCategories() []string {
	return []string{
		UploadCategoryProducts,
		UploadCategoryUsers,
		UploadCategoryRestaurants,
		UploadCategoryBanners,
	}
}

// GetValidImageMimeTypes retorna slice com todos os tipos MIME válidos
func GetValidImageMimeTypes() []string {
	return []string{
		MimeTypeJPEG,
		MimeTypeJPG,
		MimeTypePNG,
		MimeTypeWebP,
		MimeTypeGIF,
	}
}

// IsValidUploadCategory verifica se uma categoria é válida para upload
func IsValidUploadCategory(category string) bool {
	validCategories := GetValidUploadCategories()

	for _, validCategory := range validCategories {
		if category == validCategory {
			return true
		}
	}
	return false
}

// IsValidImageMimeType verifica se o tipo MIME é válido para imagens
func IsValidImageMimeType(mimeType string) bool {
	validTypes := GetValidImageMimeTypes()

	for _, validType := range validTypes {
		if mimeType == validType {
			return true
		}
	}
	return false
}