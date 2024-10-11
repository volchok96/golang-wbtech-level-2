package pattern

/*
Паттерн «Команда» (Command) — это поведенческий паттерн проектирования, 
который превращает запросы в объекты, позволяя передавать их как параметры, хранить в очереди запросов, логировать или отменять операции.

Применимость:
- Когда требуется параметризовать объекты выполняемыми операциями, например, для реализации отложенных операций.
- Для реализации системы отмены/повторения действий.
- Для логирования операций, чтобы впоследствии можно было их выполнить повторно.
- Для реализации очередей запросов или задач.

Плюсы:
- Позволяет легко расширять функциональность приложения за счет добавления новых команд без изменения существующего кода.
- Упрощает реализацию отмены и повтора операций.
- Инкапсулирует действия, позволяя передавать их как объекты.

Минусы:
- Увеличивает количество классов в системе, что может усложнить структуру программы.

Применение на практике:
- GUI-приложения: кнопка может быть связана с командой, которая выполняется при ее нажатии.
- Управление транзакциями: паттерн может быть использован для логирования и последующего отката действий.
*/

import "fmt"

// Command - интерфейс команды
type Command interface {
	Execute()
}

// Light - получатель команды (Receiver)
type Light struct {
	isOn bool
}

// TurnOn - включает свет
func (l *Light) TurnOn() {
	l.isOn = true
	fmt.Println("The light is on")
}

// TurnOff - выключает свет
func (l *Light) TurnOff() {
	l.isOn = false
	fmt.Println("The light is off")
}

// TurnOnCommand - конкретная команда для включения света
type TurnOnCommand struct {
	light *Light
}

// NewTurnOnCommand - конструктор для создания команды включения света
func NewTurnOnCommand(light *Light) *TurnOnCommand {
	return &TurnOnCommand{light: light}
}

// Execute - реализация команды включения света
func (c *TurnOnCommand) Execute() {
	c.light.TurnOn()
}

// TurnOffCommand - конкретная команда для выключения света
type TurnOffCommand struct {
	light *Light
}

// NewTurnOffCommand - конструктор для создания команды выключения света
func NewTurnOffCommand(light *Light) *TurnOffCommand {
	return &TurnOffCommand{light: light}
}

// Execute - реализация команды выключения света
func (c *TurnOffCommand) Execute() {
	c.light.TurnOff()
}

// RemoteControl - инициатор (Invoker), который вызывает команды
type RemoteControl struct {
	command Command
}

// SetCommand - задает команду для выполнения
func (r *RemoteControl) SetCommand(command Command) {
	r.command = command
}

// PressButton - выполняет команду
func (r *RemoteControl) PressButton() {
	r.command.Execute()
}

func main04() {
	// Создаем объект получателя команды - лампочку
	light := &Light{}

	// Создаем команды для включения и выключения света
	turnOn := NewTurnOnCommand(light)
	turnOff := NewTurnOffCommand(light)

	// Инициатор - пульт управления
	remote := &RemoteControl{}

	// Включаем свет с помощью команды
	remote.SetCommand(turnOn)
	remote.PressButton() // Вывод: The light is on

	// Выключаем свет с помощью команды
	remote.SetCommand(turnOff)
	remote.PressButton() // Вывод: The light is off
}
