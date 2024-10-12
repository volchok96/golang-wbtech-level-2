package pattern

import "fmt"

/*
Паттерн "Стратегия" позволяет определить семейство алгоритмов, инкапсулировать их и делать взаимозаменяемыми. Это позволяет менять алгоритмы независимо от клиентов, которые ими пользуются.

Применимость:
- Когда необходимо выбрать одну из нескольких стратегий выполнения определенной задачи.
- Когда требуется инкапсуляция алгоритмов, которые могут быть изменены в зависимости от контекста.

Плюсы:
- Уменьшает количество условных операторов (if-else) в коде.
- Упрощает расширение и поддержку, поскольку новые стратегии могут быть добавлены без изменения существующего кода.

Минусы:
- Увеличивает количество классов или структур.
- Клиент должен знать, какие стратегии ему доступны.

Пример применения на практике:
- Система скидок в e-commerce: разные стратегии расчета скидок для разных типов клиентов или акций.
- Системы обработки платежей: выбор метода оплаты в зависимости от страны или предпочитаемой платформы.
*/

// DiscountStrategy - интерфейс для различных стратегий расчета скидок
type DiscountStrategy interface {
	CalculateDiscount(totalAmount float64) float64
}

// NoDiscount - стратегия без скидки
type NoDiscount struct{}

// CalculateDiscount - возвращает полную стоимость без скидки
func (d *NoDiscount) CalculateDiscount(totalAmount float64) float64 {
	return totalAmount
}

// PercentageDiscount - стратегия расчета скидки по процентам
type PercentageDiscount struct {
	percent float64
}

// CalculateDiscount - расчет скидки по проценту от суммы
func (d *PercentageDiscount) CalculateDiscount(totalAmount float64) float64 {
	discount := totalAmount * (d.percent / 100)
	return totalAmount - discount
}

// FixedDiscount - стратегия с фиксированной скидкой
type FixedDiscount struct {
	discountAmount float64
}

// CalculateDiscount - расчет скидки фиксированной суммы
func (d *FixedDiscount) CalculateDiscount(totalAmount float64) float64 {
	return totalAmount - d.discountAmount
}

// Context - контекст, который использует стратегию для расчета скидки
type Context struct {
	strategy DiscountStrategy
}

// SetStrategy - установка стратегии расчета скидки
func (c *Context) SetStrategy(strategy DiscountStrategy) {
	c.strategy = strategy
}

// CalculateFinalPrice - расчет финальной стоимости с применением выбранной стратегии скидки
func (c *Context) CalculateFinalPrice(totalAmount float64) float64 {
	return c.strategy.CalculateDiscount(totalAmount)
}

// Пример реальной бизнес-логики для онлайн-магазина
func main07() {
	// Создаем контекст для расчета скидок
	ctx := &Context{}

	// Пример 1: Без скидки
	ctx.SetStrategy(&NoDiscount{})
	totalAmount := 100.0
	finalPrice := ctx.CalculateFinalPrice(totalAmount)
	fmt.Printf("Total amount: $%.2f, Final price (No Discount): $%.2f\n", totalAmount, finalPrice)

	// Пример 2: Скидка 10%
	ctx.SetStrategy(&PercentageDiscount{percent: 10})
	finalPrice = ctx.CalculateFinalPrice(totalAmount)
	fmt.Printf("Total amount: $%.2f, Final price (10%% Discount): $%.2f\n", totalAmount, finalPrice)

	// Пример 3: Фиксированная скидка $20
	ctx.SetStrategy(&FixedDiscount{discountAmount: 20})
	finalPrice = ctx.CalculateFinalPrice(totalAmount)
	fmt.Printf("Total amount: $%.2f, Final price (Fixed Discount $20): $%.2f\n", totalAmount, finalPrice)
}
