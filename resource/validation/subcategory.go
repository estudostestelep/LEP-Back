package validation

import (
	"lep/repositories/models"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// CreateSubcategoryValidation valida dados para criação de subcategoria
func CreateSubcategoryValidation(subcategory *models.Subcategory) error {
	return validation.ValidateStruct(subcategory,
		validation.Field(&subcategory.OrganizationId, validation.Required, is.UUID),
		validation.Field(&subcategory.ProjectId, validation.Required, is.UUID),
		validation.Field(&subcategory.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&subcategory.Order, validation.Min(0)),
	)
}

// UpdateSubcategoryValidation valida dados para atualização de subcategoria
func UpdateSubcategoryValidation(subcategory *models.Subcategory) error {
	return validation.ValidateStruct(subcategory,
		validation.Field(&subcategory.Id, validation.Required, is.UUID),
		validation.Field(&subcategory.OrganizationId, validation.Required, is.UUID),
		validation.Field(&subcategory.ProjectId, validation.Required, is.UUID),
		validation.Field(&subcategory.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&subcategory.Order, validation.Min(0)),
	)
}
