package validation

import (
	"lep/repositories/models"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// CreateClientValidation valida dados para criação de cliente
// Nota: Permissões são atribuídas via Roles, não diretamente no Client
func CreateClientValidation(client *models.Client) error {
	return validation.ValidateStruct(client,
		validation.Field(&client.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&client.Email, validation.Required, is.Email),
		validation.Field(&client.Password, validation.Required, validation.Length(6, 255)),
		validation.Field(&client.OrgId, validation.Required),
	)
}

// UpdateClientValidation valida dados para atualização de cliente
// Nota: Permissões são atribuídas via Roles, não diretamente no Client
// Nota: OrgId é imutável após criação, não validado aqui
func UpdateClientValidation(client *models.Client) error {
	return validation.ValidateStruct(client,
		validation.Field(&client.Id, validation.Required),
		validation.Field(&client.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&client.Email, validation.Required, is.Email),
		validation.Field(&client.Password, validation.When(client.Password != "", validation.Length(6, 255))),
	)
}
