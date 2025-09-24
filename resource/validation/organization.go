package validation

import (
	"lep/repositories/models"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// CreateOrganizationValidation valida dados para criação de organização
func CreateOrganizationValidation(org *models.Organization) error {
	return validation.ValidateStruct(org,
		validation.Field(&org.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&org.Email, is.Email),
		validation.Field(&org.Phone, validation.Length(0, 20)),
		validation.Field(&org.Address, validation.Length(0, 500)),
		validation.Field(&org.Website, is.URL),
		validation.Field(&org.Description, validation.Length(0, 1000)),
	)
}

// UpdateOrganizationValidation valida dados para atualização de organização
func UpdateOrganizationValidation(org *models.Organization) error {
	return validation.ValidateStruct(org,
		validation.Field(&org.Id, validation.Required, is.UUID),
		validation.Field(&org.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&org.Email, validation.Required, is.Email),
		validation.Field(&org.Phone, validation.Length(0, 20)),
		validation.Field(&org.Address, validation.Length(0, 500)),
		validation.Field(&org.Website, is.URL),
		validation.Field(&org.Description, validation.Length(0, 1000)),
	)
}
