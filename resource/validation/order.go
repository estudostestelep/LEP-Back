package validation

import (
	"lep/repositories/models"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

func CreateOrderValidation(e *models.Order) error {
	return validation.ValidateStruct(e,
		validation.Field(&e.OrganizationId, validation.Required, is.UUID),
		validation.Field(&e.ProjectId, validation.Required, is.UUID),
		validation.Field(&e.CustomerId, validation.Required, is.UUID),
		validation.Field(&e.Items, validation.Required),
		validation.Field(&e.TableId, validation.Required, is.UUID),
	)
}

// UpdateOrderValidation valida dados para atualização de pedido
func UpdateOrderValidation(e *models.Order) error {
	return validation.ValidateStruct(e,
		validation.Field(&e.Id, validation.Required, is.UUID),
		validation.Field(&e.OrganizationId, validation.Required, is.UUID),
		validation.Field(&e.ProjectId, validation.Required, is.UUID),
		validation.Field(&e.CustomerId, validation.Required, is.UUID),
		validation.Field(&e.Items, validation.Required),
		validation.Field(&e.Status, validation.Required, validation.In("draft", "pending", "completed", "canceled")),
		validation.Field(&e.TableId, validation.Required, is.UUID),
	)
}
