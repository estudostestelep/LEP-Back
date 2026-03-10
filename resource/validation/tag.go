package validation

import (
	"lep/repositories/models"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// CreateTagValidation valida dados para criação de tag
func CreateTagValidation(tag *models.Tag) error {
	return validation.ValidateStruct(tag,
		validation.Field(&tag.OrganizationId, validation.Required, is.UUID),
		validation.Field(&tag.ProjectId, validation.Required, is.UUID),
		validation.Field(&tag.Name, validation.Required, validation.Length(1, 50)),
		// Color é opcional - só validar formato se preenchido
		validation.Field(&tag.Color, validation.When(tag.Color != "", validation.Match(colorRegex).Error("must be a valid hex color (e.g., #FF5733)"))),
		validation.Field(&tag.EntityType, validation.In("product", "customer", "table", "reservation", "order", "")),
	)
}

// UpdateTagValidation valida dados para atualização de tag
func UpdateTagValidation(tag *models.Tag) error {
	return validation.ValidateStruct(tag,
		validation.Field(&tag.Id, validation.Required, is.UUID),
		validation.Field(&tag.OrganizationId, validation.Required, is.UUID),
		validation.Field(&tag.ProjectId, validation.Required, is.UUID),
		validation.Field(&tag.Name, validation.Required, validation.Length(1, 50)),
		// Color é opcional - só validar formato se preenchido
		validation.Field(&tag.Color, validation.When(tag.Color != "", validation.Match(colorRegex).Error("must be a valid hex color (e.g., #FF5733)"))),
		validation.Field(&tag.EntityType, validation.In("product", "customer", "table", "reservation", "order", "")),
	)
}
