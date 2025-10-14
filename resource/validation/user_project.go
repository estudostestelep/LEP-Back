package validation

import (
	"lep/repositories/models"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// CreateUserProjectValidation valida dados para criação de relacionamento usuário-projeto
func CreateUserProjectValidation(userProj *models.UserProject) error {
	return validation.ValidateStruct(userProj,
		validation.Field(&userProj.UserId, validation.Required, is.UUID),
		validation.Field(&userProj.ProjectId, validation.Required, is.UUID),
		validation.Field(&userProj.Role, validation.Required, validation.In("admin", "manager", "waiter", "member")),
	)
}

// UpdateUserProjectValidation valida dados para atualização de relacionamento usuário-projeto
func UpdateUserProjectValidation(userProj *models.UserProject) error {
	return validation.ValidateStruct(userProj,
		validation.Field(&userProj.Id, validation.Required, is.UUID),
		validation.Field(&userProj.Role, validation.Required, validation.In("admin", "manager", "waiter", "member")),
	)
}
