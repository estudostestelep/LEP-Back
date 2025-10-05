package validation

import (
	"lep/repositories/models"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// CreateMenuValidation valida dados para criação de cardápio
func CreateMenuValidation(menu *models.Menu) error {
	return validation.ValidateStruct(menu,
		validation.Field(&menu.OrganizationId, validation.Required, is.UUID),
		validation.Field(&menu.ProjectId, validation.Required, is.UUID),
		validation.Field(&menu.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&menu.Order, validation.Min(0)),
	)
}

// UpdateMenuValidation valida dados para atualização de cardápio
func UpdateMenuValidation(menu *models.Menu) error {
	return validation.ValidateStruct(menu,
		validation.Field(&menu.Id, validation.Required, is.UUID),
		validation.Field(&menu.OrganizationId, validation.Required, is.UUID),
		validation.Field(&menu.ProjectId, validation.Required, is.UUID),
		validation.Field(&menu.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&menu.Order, validation.Min(0)),
	)
}
