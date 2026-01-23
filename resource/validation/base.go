package validation

import (
	"fmt"
	"regexp"

	"lep/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// ValidationRule representa uma regra de validação customizada
type ValidationRule struct {
	Field   string
	Rules   []validation.Rule
	Message string
}

// Regex patterns
var (
	colorRegex = regexp.MustCompile(`^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`)
)

// Common validation rules that can be reused
var (
	RequiredUUID     = []validation.Rule{validation.Required, is.UUID}
	RequiredString   = []validation.Rule{validation.Required, validation.Length(1, 255)}
	RequiredEmail    = []validation.Rule{validation.Required, is.Email}
	RequiredPhone    = []validation.Rule{validation.Required, validation.Length(8, 20)}
	OptionalString   = []validation.Rule{validation.Length(0, 255)}
	OptionalURL      = []validation.Rule{is.URL}
	RequiredPositive = []validation.Rule{validation.Required, validation.Min(0.01)}
)

// ValidateUUIDParam valida se uma string é um UUID válido
func ValidateUUIDParam(value string, fieldName string) error {
	return validation.Validate(value, validation.Required.Error(fieldName+" is required"), is.UUID.Error("Invalid "+fieldName+" format"))
}

// ValidateRequiredString valida se uma string não está vazia
func ValidateRequiredString(value string, fieldName string) error {
	return validation.Validate(value,
		validation.Required.Error(fieldName+" is required"),
		validation.Length(1, 255).Error(fieldName+" must be between 1 and 255 characters"))
}

// ValidateEmail valida formato de email
func ValidateEmail(value string, fieldName string) error {
	return validation.Validate(value,
		validation.Required.Error(fieldName+" is required"),
		is.Email.Error("Invalid "+fieldName+" format"))
}

// ValidatePositiveNumber valida se um número é positivo
func ValidatePositiveNumber(value float64, fieldName string) error {
	return validation.Validate(value,
		validation.Required.Error(fieldName+" is required"),
		validation.Min(0.01).Error(fieldName+" must be greater than 0"))
}

// ParseAndValidateUUID valida e retorna UUID parseado
// Retorna false se inválido (erro já enviado ao client)
func ParseAndValidateUUID(c *gin.Context, idStr string, entityName string) (uuid.UUID, bool) {
	if err := ValidateUUIDParam(idStr, entityName+" ID"); err != nil {
		utils.SendBadRequestError(c, fmt.Sprintf("Invalid %s ID format", entityName), err)
		return uuid.Nil, false
	}
	parsed, _ := uuid.Parse(idStr)
	return parsed, true
}