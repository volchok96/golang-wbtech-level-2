package pattern

/*
Паттерн "Посетитель" (Visitor) — это поведенческий паттерн проектирования, который позволяет добавлять новые операции для классов, не изменяя их код. 
Он отделяет алгоритмы от структуры данных, что делает возможным добавление новых операций без изменения существующих классов.

Применимость:
- Когда необходимо добавить новые операции к объектам, но не менять их структуру.
= Когда нужно выполнять различные операции над объектами разных типов: 
"Посетитель" позволяет избежать дублирования кода при обработке объектов различных типов, особенно когда это зависит от типа объекта.

Плюсы:
- Расширяемость: легко добавлять новые операции над объектами.
- Разделение ответственности: логика операций и структура данных разделены, что улучшает модульность и поддерживаемость кода.

Минусы:
- Сложность добавления новых типов объектов
- Сложность сопровождения: Увеличение числа классов-посетителей может усложнить структуру проекта, особенно если они зависят друг от друга.

Применение на практике:
- Системы выставления счетов: можно иметь различные типы документов (счета, контракты, накладные), 
и "Посетитель" позволяет обрабатывать их без изменения классов самих документов.
- Парсеры или компиляторы: для обхода различных синтаксических конструкций, позволяя добавлять новые обработки (например, вычисления или оптимизации).
*/

import "fmt"

// Element - интерфейс для всех объектов, которые могут быть "посещены" посетителем.
type Element interface {
	Accept(visitor Visitor)
}

// Product - товар с ценой.
type Product struct {
	Price float64
}

// Accept - метод для принятия посетителя в объекте Product.
func (p *Product) Accept(visitor Visitor) {
	visitor.VisitProduct(p)
}

// Service - услуга с почасовой ставкой.
type Service struct {
	HourlyRate float64
}

// Accept - метод для принятия посетителя в объекте Service.
func (s *Service) Accept(visitor Visitor) {
	visitor.VisitService(s)
}

// Visitor - интерфейс для посетителя, который будет обрабатывать разные объекты.
type Visitor interface {
	VisitProduct(p *Product)
	VisitService(s *Service)
}

// TaxCalculator - калькулятор налогов (один из посетителей).
type TaxCalculator struct {
	totalTax float64
}

// VisitProduct - метод для расчета налога для товара.
func (t *TaxCalculator) VisitProduct(p *Product) {
	tax := p.Price * 0.10 // Налог 10% для товаров.
	t.totalTax += tax
	fmt.Printf("Tax for product: $%.2f\n", tax)
}

// VisitService - метод для расчета налога для услуги.
func (t *TaxCalculator) VisitService(s *Service) {
	tax := s.HourlyRate * 0.15 // Налог 15% для услуг.
	t.totalTax += tax
	fmt.Printf("Tax for service: $%.2f\n", tax)
}

// DiscountCalculator - калькулятор скидок (другой посетитель).
type DiscountCalculator struct {
	totalDiscount float64
}

// VisitProduct - метод для расчета скидки для товара.
func (d *DiscountCalculator) VisitProduct(p *Product) {
	discount := p.Price * 0.05 // Скидка 5% для товаров.
	d.totalDiscount += discount
	fmt.Printf("Discount for product: $%.2f\n", discount)
}

// VisitService - метод для расчета скидки для услуги.
func (d *DiscountCalculator) VisitService(s *Service) {
	discount := s.HourlyRate * 0.07 // Скидка 7% для услуг.
	d.totalDiscount += discount
	fmt.Printf("Discount for service: $%.2f\n", discount)
}

func main03() {
	// Создаем объекты — товар и услугу.
	product := &Product{Price: 100}
	service := &Service{HourlyRate: 50}

	// Создаем посетителей — калькулятор налогов и калькулятор скидок.
	taxCalculator := &TaxCalculator{}
	discountCalculator := &DiscountCalculator{}

	// Применяем калькулятор налогов к товару и услуге.
	product.Accept(taxCalculator)
	service.Accept(taxCalculator)
	fmt.Printf("Total tax: $%.2f\n", taxCalculator.totalTax)

	// Применяем калькулятор скидок к товару и услуге.
	product.Accept(discountCalculator)
	service.Accept(discountCalculator)
	fmt.Printf("Total discount: $%.2f\n", discountCalculator.totalDiscount)
}
