package pattern

/*
Паттерн "Фасад" упрощает взаимодействие с системой обработки заказов, скрывая сложные детали
работы с различными внутренними сервисами (инвентаризация, платежи, доставка).

Применимость:
- Упрощает работу с разрозненными сервисами (например, инвентаризация, платежи, доставка).
- Снижает связанность между клиентским кодом и сложными подсистемами.

Плюсы:
- Упрощает работу с системой, предоставляя единый интерфейс.
- Улучшает поддерживаемость кода.

Минусы:
- Добавляет дополнительный уровень абстракции.
- Может усложнить интеграцию новых подсистем.

Применение на практике:
- Системы электронной коммерции, где требуется взаимодействовать с несколькими внутренними системами.
- Платежные системы, работающие с различными шлюзами и методами оплаты.
*/

import (
	"fmt"
)

// InventoryService - Сервис для работы с инвентарем
type InventoryService struct{}

// CheckStock проверяет наличие товара на складе
func (i *InventoryService) CheckStock(itemID string) bool {
	fmt.Printf("Checking stock for item: %s\n", itemID)
	// Предположим, что товар всегда есть на складе
	return true
}

// PaymentService - Сервис для обработки платежей
type PaymentService struct{}

// ProcessPayment обрабатывает платеж
func (p *PaymentService) ProcessPayment(amount float64) bool {
	fmt.Printf("Processing payment of: $%.2f\n", amount)
	// Предположим, что платеж всегда успешен
	return true
}

// ShippingService - Сервис для доставки товара
type ShippingService struct{}

// ShipItem организует доставку товара
func (s *ShippingService) ShipItem(itemID, address string) {
	fmt.Printf("Shipping item: %s to address: %s\n", itemID, address)
}

// OrderProcessingFacade - Фасад для системы обработки заказов
type OrderProcessingFacade struct {
	inventory  *InventoryService
	payment    *PaymentService
	shipping   *ShippingService
}

// NewOrderProcessingFacade - Конструктор фасада
func NewOrderProcessingFacade() *OrderProcessingFacade {
	return &OrderProcessingFacade{
		inventory:  &InventoryService{},
		payment:    &PaymentService{},
		shipping:   &ShippingService{},
	}
}

// PlaceOrder - Метод фасада для обработки заказа
func (o *OrderProcessingFacade) PlaceOrder(itemID string, amount float64, address string) error {
	// Проверка наличия товара
	if !o.inventory.CheckStock(itemID) {
		return fmt.Errorf("item %s is out of stock", itemID)
	}

	// Обработка платежа
	if !o.payment.ProcessPayment(amount) {
		return fmt.Errorf("payment failed for amount $%.2f", amount)
	}

	// Организация доставки
	o.shipping.ShipItem(itemID, address)

	fmt.Println("Order placed successfully")
	return nil
}

func main01() {
	// Создаем фасад для обработки заказов
	orderFacade := NewOrderProcessingFacade()

	// Используем фасад для размещения заказа
	err := orderFacade.PlaceOrder("item1", 99.99, "Main St, DC")
	if err != nil {
		fmt.Println("Failed to place order:", err)
	} else {
		fmt.Println("Order placed successfully!")
	}
}
