package validation

import (
	"errors"
	"lep/repositories/models"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// CreateProductValidation valida dados para criação de produto
func CreateProductValidation(product *models.Product) error {
	// Validação base
	if err := validation.ValidateStruct(product,
		validation.Field(&product.OrganizationId, validation.Required, is.UUID),
		validation.Field(&product.ProjectId, validation.Required, is.UUID),
		validation.Field(&product.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&product.Description, validation.Length(0, 500)),
	); err != nil {
		return err
	}

	// Validação de preço condicional por tipo
	return validateProductPrice(product)
}

// UpdateProductValidation valida dados para atualização de produto
func UpdateProductValidation(product *models.Product) error {
	// Validação base
	if err := validation.ValidateStruct(product,
		validation.Field(&product.Id, validation.Required, is.UUID),
		validation.Field(&product.OrganizationId, validation.Required, is.UUID),
		validation.Field(&product.ProjectId, validation.Required, is.UUID),
		validation.Field(&product.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&product.Description, validation.Length(0, 500)),
	); err != nil {
		return err
	}

	// Validação de preço condicional por tipo
	return validateProductPrice(product)
}

// validateProductPrice valida preço baseado no tipo do produto
// Vinhos podem ter price_normal=0 se tiverem price_bottle ou price_glass
func validateProductPrice(product *models.Product) error {
	// Validação de preço promocional
	if product.UsePromo {
		if product.PricePromo == nil || *product.PricePromo <= 0 {
			return errors.New("promotional price must be set and greater than 0 when use_promo is true")
		}
		if product.PriceNormal > 0 && *product.PricePromo >= product.PriceNormal {
			return errors.New("promotional price must be less than normal price")
		}
	}

	// Para vinhos: pelo menos um preço deve estar definido
	if product.Type == "vinho" {
		hasBottlePrice := product.PriceBottle != nil && *product.PriceBottle > 0
		hasGlassPrice := product.PriceGlass != nil && *product.PriceGlass > 0
		hasNormalPrice := product.PriceNormal > 0

		if !hasBottlePrice && !hasGlassPrice && !hasNormalPrice {
			return errors.New("wine products must have at least one price (price_normal, price_bottle, or price_glass)")
		}
		return nil
	}

	// Para outros produtos: price_normal é obrigatório e deve ser > 0
	if product.PriceNormal <= 0 {
		return errors.New("price_normal is required and must be greater than 0")
	}
	return nil
}
