package validation

import (
	"lep/repositories/models"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// CreateCustomerValidation valida dados para criação de cliente
func CreateCustomerValidation(customer *models.Customer) error {
	return validation.ValidateStruct(customer,
		validation.Field(&customer.OrganizationId, validation.Required, is.UUID),
		validation.Field(&customer.ProjectId, validation.Required, is.UUID),
		validation.Field(&customer.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&customer.Email, is.Email),
		validation.Field(&customer.Phone, validation.Required, validation.Length(8, 20)),
	)
}

// UpdateCustomerValidation valida dados para atualização de cliente
func UpdateCustomerValidation(customer *models.Customer) error {
	return validation.ValidateStruct(customer,
		validation.Field(&customer.Id, validation.Required, is.UUID),
		validation.Field(&customer.OrganizationId, validation.Required, is.UUID),
		validation.Field(&customer.ProjectId, validation.Required, is.UUID),
		validation.Field(&customer.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&customer.Email, is.Email),
		validation.Field(&customer.Phone, validation.Required, validation.Length(8, 20)),
	)
}