package validate

import (
	"lep/repositories/models"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

func CreateOrderRequestValidation(e *models.Order) error {
	return validation.ValidateStruct(e,
		validation.Field(&e.OrganizationId, validation.Required, is.UUID),
		validation.Field(&e.ProjectId, validation.Required, is.UUID),
		validation.Field(&e.CustomerId, validation.Required, is.UUID),
		validation.Field(&e.TableNumber, validation.Required, validation.Length(1, 20)),
		validation.Field(&e.Items, validation.Required, validation.Length(1, 0)),
		validation.Field(&e.Status, validation.Required, validation.In("draft", "pending", "completed", "canceled")),
		validation.Field(&e.TableId, validation.Required, is.UUID),
	)
}
