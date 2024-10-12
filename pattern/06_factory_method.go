package pattern

import "fmt"

/*
Паттерн "Фабричный метод" предоставляет интерфейс для создания объектов, но позволяет подклассам изменять тип создаваемого объекта.

Применимость:
- Когда заранее неизвестно, объекты каких типов необходимо создавать.
- Для отделения логики создания объектов от их основной функциональности.

Плюсы:
- Уменьшает связанность кода, отделяя процесс создания объектов от их использования.
- Облегчает добавление новых типов объектов.

Минусы:
- Увеличивает количество классов в системе.
- Может усложнить код из-за множества подклассов.

Примеры использования:
- Фреймворки для работы с графическими интерфейсами создают различные элементы интерфейса через фабричные методы.
- Логирование, в зависимости от уровня или назначения, может использовать разные фабрики для создания логгеров.
- Система для отправки уведомлений (Email, SMS, Push).
*/

// Notifier - интерфейс для отправки уведомлений
type Notifier interface {
	SendNotification(message, recipient string) error
}

// EmailNotifier - отправка уведомлений через Email
type EmailNotifier struct{}

// SendNotification - отправка уведомления по Email
func (e *EmailNotifier) SendNotification(message, recipient string) error {
	// Реальная логика отправки Email
	fmt.Printf("Sending email to %s: %s\n", recipient, message)
	return nil
}

// SMSNotifier - отправка уведомлений через SMS
type SMSNotifier struct{}

// SendNotification - отправка уведомления по SMS
func (s *SMSNotifier) SendNotification(message, recipient string) error {
	// Реальная логика отправки SMS
	fmt.Printf("Sending SMS to %s: %s\n", recipient, message)
	return nil
}

// PushNotifier - отправка push-уведомлений
type PushNotifier struct{}

// SendNotification - отправка push-уведомлений
func (p *PushNotifier) SendNotification(message, recipient string) error {
	// Реальная логика отправки push-уведомления
	fmt.Printf("Sending Push Notification to %s: %s\n", recipient, message)
	return nil
}

// NotificationFactory - интерфейс для фабрики, создающей объекты Notifier
type NotificationFactory interface {
	CreateNotifier() Notifier
}

// EmailFactory - фабрика для создания Email-уведомлений
type EmailFactory struct{}

// CreateNotifier - создание Email-уведомления
func (f *EmailFactory) CreateNotifier() Notifier {
	return &EmailNotifier{}
}

// SMSFactory - фабрика для создания SMS-уведомлений
type SMSFactory struct{}

// CreateNotifier - создание SMS-уведомления
func (f *SMSFactory) CreateNotifier() Notifier {
	return &SMSNotifier{}
}

// PushFactory - фабрика для создания Push-уведомлений
type PushFactory struct{}

// CreateNotifier - создание Push-уведомления
func (f *PushFactory) CreateNotifier() Notifier {
	return &PushNotifier{}
}

// Client - функция для использования фабричного метода
func Client(factory NotificationFactory, message, recipient string) {
	notifier := factory.CreateNotifier()
	err := notifier.SendNotification(message, recipient)
	if err != nil {
		fmt.Println("Error sending notification:", err)
	}
}

func main06() {
	// Пример использования Email-уведомлений
	emailFactory := &EmailFactory{}
	Client(emailFactory, "Welcome to our service!", "user1@example.com")

	// Пример использования SMS-уведомлений
	smsFactory := &SMSFactory{}
	Client(smsFactory, "Your login: NewUser", "+1234567890")

	// Пример использования Push-уведомлений
	pushFactory := &PushFactory{}
	Client(pushFactory, "New message", "user_device_token")
}
