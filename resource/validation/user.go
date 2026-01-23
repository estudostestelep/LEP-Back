package validation

import (
	"fmt"
	"lep/constants"
	"lep/repositories/models"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
	"github.com/lib/pq"
)

// ValidPermissionRule valida se todas as permissões são válidas
var ValidPermissionRule = validation.By(func(value interface{}) error {
	permissions, ok := value.(pq.StringArray)
	if !ok {
		return nil // Deixar outras validações tratarem tipo inválido
	}

	// Validar cada permissão contra whitelist
	for _, perm := range permissions {
		if perm == "" {
			return fmt.Errorf("permissão não pode ser vazia")
		}
		if !constants.IsValidPermission(perm) {
			return fmt.Errorf("permissão inválida: %s", perm)
		}
	}
	return nil
})

// CreateUserValidation valida dados para criação de usuário
func CreateUserValidation(user *models.User) error {
	return validation.ValidateStruct(user,
		validation.Field(&user.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Password, validation.Required, validation.Length(6, 255)),
		validation.Field(&user.Permissions, ValidPermissionRule),
	)
}

// UpdateUserValidation valida dados para atualização de usuário
func UpdateUserValidation(user *models.User) error {
	return validation.ValidateStruct(user,
		validation.Field(&user.Id, validation.Required, is.UUID),
		validation.Field(&user.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&user.Email, validation.Required, is.Email),
		// Password é opcional na atualização
		validation.Field(&user.Password, validation.When(user.Password != "", validation.Length(6, 255))),
		validation.Field(&user.Permissions, ValidPermissionRule),
	)
}

// LoginValidation valida dados de login
func LoginValidation(email, password string) error {
	return validation.Errors{
		"email":    validation.Validate(email, validation.Required, is.Email),
		"password": validation.Validate(password, validation.Required, validation.Length(1, 255)),
	}.Filter()
}
