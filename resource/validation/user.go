package validation

import (
	"lep/repositories/models"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// CreateUserValidation valida dados para criação de usuário
func CreateUserValidation(user *models.User) error {
	return validation.ValidateStruct(user,
		validation.Field(&user.OrganizationId, validation.Required, is.UUID),
		validation.Field(&user.ProjectId, validation.Required, is.UUID),
		validation.Field(&user.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Password, validation.Required, validation.Length(6, 255)),
		validation.Field(&user.Role, validation.Required, validation.In("admin", "manager", "waiter", "kitchen", "cashier")),
		validation.Field(&user.Permissions, validation.Each(validation.Length(1, 50))),
	)
}

// UpdateUserValidation valida dados para atualização de usuário
func UpdateUserValidation(user *models.User) error {
	return validation.ValidateStruct(user,
		validation.Field(&user.Id, validation.Required, is.UUID),
		validation.Field(&user.OrganizationId, validation.Required, is.UUID),
		validation.Field(&user.ProjectId, validation.Required, is.UUID),
		validation.Field(&user.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&user.Email, validation.Required, is.Email),
		// Password é opcional na atualização
		validation.Field(&user.Password, validation.When(user.Password != "", validation.Length(6, 255))),
		validation.Field(&user.Role, validation.Required, validation.In("admin", "manager", "waiter", "kitchen", "cashier")),
		validation.Field(&user.Permissions, validation.Each(validation.Length(1, 50))),
	)
}

// LoginValidation valida dados de login
func LoginValidation(email, password string) error {
	return validation.Errors{
		"email":    validation.Validate(email, validation.Required, is.Email),
		"password": validation.Validate(password, validation.Required, validation.Length(1, 255)),
	}.Filter()
}