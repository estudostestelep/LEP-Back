package validation

import (
	"lep/repositories/models"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// CreateUserOrganizationValidation valida dados para criação de relacionamento usuário-organização
func CreateUserOrganizationValidation(userOrg *models.UserOrganization) error {
	return validation.ValidateStruct(userOrg,
		validation.Field(&userOrg.UserId, validation.Required, is.UUID),
		validation.Field(&userOrg.OrganizationId, validation.Required, is.UUID),
		validation.Field(&userOrg.Role, validation.Required, validation.In("owner", "admin", "member")),
	)
}

// UpdateUserOrganizationValidation valida dados para atualização de relacionamento usuário-organização
func UpdateUserOrganizationValidation(userOrg *models.UserOrganization) error {
	return validation.ValidateStruct(userOrg,
		validation.Field(&userOrg.Id, validation.Required, is.UUID),
		validation.Field(&userOrg.Role, validation.Required, validation.In("owner", "admin", "member")),
	)
}
