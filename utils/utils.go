package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Extrai o userId do contexto (exemplo; ajuste conforme autenticação real)
func GetUserIdFromContext(c *gin.Context) uuid.UUID {
	userIdStr, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil
	}
	userId, err := uuid.Parse(userIdStr.(string))
	if err != nil {
		return uuid.Nil
	}
	return userId
}

// Valida se um UUId é válido
func IsValidUUId(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

// Gera um novo UUId
func NewUUId() uuid.UUID {
	return uuid.New()
}

// Função para calcular total de um pedido
// Recebe slice de OrderItem, retorna soma dos subtotais
func CalculateOrderTotal(items []OrderItem) float64 {
	total := 0.0
	for _, item := range items {
		total += float64(item.Quantity) * item.Price
	}
	return total
}

// OrderItem model simplificado para utils
type OrderItem struct {
	Quantity int
	Price    float64
}
