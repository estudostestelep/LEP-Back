package validation

import (
	"lep/repositories/models"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// CreateProductValidation valida dados para criação de produto
func CreateProductValidation(product *models.Product) error {
	return validation.ValidateStruct(product,
		validation.Field(&product.OrganizationId, validation.Required, is.UUID),
		validation.Field(&product.ProjectId, validation.Required, is.UUID),
		validation.Field(&product.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&product.Description, validation.Length(0, 500)),
		validation.Field(&product.Price, validation.Required, validation.Min(0.01)),
		validation.Field(&product.PrepTimeMinutes, validation.Required, validation.Min(1)),
	)
}

// UpdateProductValidation valida dados para atualização de produto
func UpdateProductValidation(product *models.Product) error {
	return validation.ValidateStruct(product,
		validation.Field(&product.Id, validation.Required, is.UUID),
		validation.Field(&product.OrganizationId, validation.Required, is.UUID),
		validation.Field(&product.ProjectId, validation.Required, is.UUID),
		validation.Field(&product.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&product.Description, validation.Length(0, 500)),
		validation.Field(&product.Price, validation.Required, validation.Min(0.01)),
		validation.Field(&product.PrepTimeMinutes, validation.Required, validation.Min(1)),
	)
}