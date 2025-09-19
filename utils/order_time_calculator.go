package utils

import (
	"lep/repositories/models"
	"time"
)

// CalculateOrderPrepTime calcula o tempo total de preparo de um pedido
// baseado nos tempos individuais dos produtos e suas quantidades
func CalculateOrderPrepTime(items []models.OrderItem, products []models.Product) int {
	if len(items) == 0 {
		return 0
	}

	// Cria mapa de produtos para busca rápida
	productMap := make(map[string]models.Product)
	for _, product := range products {
		productMap[product.Id.String()] = product
	}

	maxPrepTime := 0
	totalSequentialTime := 0

	// Para cada item do pedido
	for _, item := range items {
		product, exists := productMap[item.ProductId.String()]
		if !exists {
			continue // Produto não encontrado, pula
		}

		// Tempo do item = tempo_preparo * quantidade
		itemTime := product.PrepTimeMinutes * item.Quantity

		// Acumula tempo sequencial (se fosse um por vez)
		totalSequentialTime += itemTime

		// Acha o maior tempo individual (assumindo preparo paralelo)
		if itemTime > maxPrepTime {
			maxPrepTime = itemTime
		}
	}

	// Estratégia híbrida: usa o maior tempo individual + 20% do tempo sequencial
	// Isso considera que alguns itens podem ser preparados em paralelo,
	// mas há alguma sobrecarga sequencial
	estimatedTime := maxPrepTime + int(float64(totalSequentialTime-maxPrepTime)*0.2)

	return estimatedTime
}

// CalculateEstimatedDeliveryTime calcula o horário estimado de entrega
// considerando o tempo de preparo e a fila da cozinha
func CalculateEstimatedDeliveryTime(prepTimeMinutes int, kitchenQueueMinutes int) time.Time {
	now := time.Now()
	totalMinutes := prepTimeMinutes + kitchenQueueMinutes
	return now.Add(time.Duration(totalMinutes) * time.Minute)
}

// GetKitchenQueueTime calcula o tempo atual da fila da cozinha
// baseado nos pedidos que estão sendo preparados
func GetKitchenQueueTime(activeOrders []models.Order) int {
	queueTime := 0

	for _, order := range activeOrders {
		if order.Status == "preparing" {
			// Calcula quanto tempo ainda falta para este pedido
			if order.EstimatedDeliveryTime != nil {
				remaining := int(time.Until(*order.EstimatedDeliveryTime).Minutes())
				if remaining > 0 {
					queueTime += remaining
				}
			}
		}
	}

	// Aplica fator de paralelismo (assumindo que a cozinha pode preparar
	// múltiplos itens ao mesmo tempo, mas com alguma limitação)
	if queueTime > 0 {
		queueTime = int(float64(queueTime) * 0.7) // 70% do tempo total
	}

	return queueTime
}

// UpdateOrderStatus atualiza o status do pedido e timestamps relevantes
func UpdateOrderStatus(order *models.Order, newStatus string) {
	now := time.Now()
	order.Status = newStatus
	order.UpdatedAt = now

	switch newStatus {
	case "preparing":
		order.StartedAt = &now
	case "ready":
		order.ReadyAt = &now
	case "delivered":
		order.DeliveredAt = &now
	}
}

// GetOrderProgress retorna o progresso do pedido em porcentagem
func GetOrderProgress(order models.Order) float64 {
	switch order.Status {
	case "pending":
		return 0.0
	case "preparing":
		if order.EstimatedDeliveryTime != nil {
			elapsed := time.Since(order.CreatedAt).Minutes()
			total := order.EstimatedDeliveryTime.Sub(order.CreatedAt).Minutes()
			if total > 0 {
				progress := elapsed / total * 100
				if progress > 99 {
					return 99.0 // Nunca 100% até estar ready
				}
				return progress
			}
		}
		return 25.0 // Fallback se não tiver tempo estimado
	case "ready":
		return 100.0
	case "delivered":
		return 100.0
	case "cancelled":
		return 0.0
	default:
		return 0.0
	}
}

// GetRemainingTime retorna quantos minutos restam para o pedido ficar pronto
func GetRemainingTime(order models.Order) int {
	if order.Status == "ready" || order.Status == "delivered" {
		return 0
	}

	if order.EstimatedDeliveryTime != nil {
		remaining := int(time.Until(*order.EstimatedDeliveryTime).Minutes())
		if remaining < 0 {
			return 0 // Pedido atrasado
		}
		return remaining
	}

	return order.EstimatedPrepTime // Fallback para tempo estimado inicial
}