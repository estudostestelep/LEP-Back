package validation

import (
	"lep/repositories/models"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// CreateCategoryValidation valida dados para criação de categoria
func CreateCategoryValidation(category *models.Category) error {
	return validation.ValidateStruct(category,
		validation.Field(&category.OrganizationId, validation.Required, is.UUID),
		validation.Field(&category.ProjectId, validation.Required, is.UUID),
		validation.Field(&category.MenuId, validation.Required, is.UUID),
		validation.Field(&category.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&category.Order, validation.Min(0)),
	)
}

// UpdateCategoryValidation valida dados para atualização de categoria
func UpdateCategoryValidation(category *models.Category) error {
	return validation.ValidateStruct(category,
		validation.Field(&category.Id, validation.Required, is.UUID),
		validation.Field(&category.OrganizationId, validation.Required, is.UUID),
		validation.Field(&category.ProjectId, validation.Required, is.UUID),
		validation.Field(&category.MenuId, validation.Required, is.UUID),
		validation.Field(&category.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&category.Order, validation.Min(0)),
	)
}
