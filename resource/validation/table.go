package validation

import (
	"lep/repositories/models"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// CreateTableValidation valida dados para criação de mesa
func CreateTableValidation(table *models.Table) error {
	return validation.ValidateStruct(table,
		validation.Field(&table.OrganizationId, validation.Required, is.UUID),
		validation.Field(&table.ProjectId, validation.Required, is.UUID),
		validation.Field(&table.Number, validation.Required, validation.Min(1)),
		validation.Field(&table.Capacity, validation.Required, validation.Min(1), validation.Max(20)),
		validation.Field(&table.Location, validation.Length(0, 100)),
		validation.Field(&table.Status, validation.Required, validation.In("livre", "ocupada", "reservada")),
	)
}

// UpdateTableValidation valida dados para atualização de mesa
func UpdateTableValidation(table *models.Table) error {
	return validation.ValidateStruct(table,
		validation.Field(&table.Id, validation.Required, is.UUID),
		validation.Field(&table.OrganizationId, validation.Required, is.UUID),
		validation.Field(&table.ProjectId, validation.Required, is.UUID),
		validation.Field(&table.Number, validation.Required, validation.Min(1)),
		validation.Field(&table.Capacity, validation.Required, validation.Min(1), validation.Max(20)),
		validation.Field(&table.Location, validation.Length(0, 100)),
		validation.Field(&table.Status, validation.Required, validation.In("livre", "ocupada", "reservada")),
	)
}