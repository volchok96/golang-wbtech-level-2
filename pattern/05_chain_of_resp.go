package pattern

/*
Применимость:
- Когда нужно передать запрос по цепочке потенциальных обработчиков.
- Когда не нужно знать, какой обработчик выполнит запрос.
- Для реализации системы фильтрации, валидации данных или обработки запросов с разными уровнями доступа.

Плюсы:
- Уменьшение зависимости между клиентом и обработчиками.
- Легко добавлять новые обработчики в цепочку.
- Упрощает расширение функциональности, поскольку каждая конкретная реализация обработчика изолирована.

Минусы:
- Не гарантируется, что запрос будет обработан.
- Может быть сложно отладить цепочку обработчиков, если она становится слишком длинной.
- Некоторые запросы могут пройти через много обработчиков, что может повлиять на производительность.

Применение на практике:
- Системы технической поддержки: Запросы передаются от одного уровня поддержки к другому в зависимости от сложности проблемы.
- Межсетевые экраны и системы безопасности: Пакеты данных могут проходить через несколько уровней фильтрации или проверки.
- Логгирование: записи могут обрабатываться разными логгерами в зависимости от важности (уровни логирования — DEBUG, INFO, ERROR).
*/

import (
	"fmt"
	"strings"
)

// LogLevel - тип для определения уровня логирования (DEBUG, INFO, ERROR)
type LogLevel int

const (
	// DEBUG - уровень логирования для отладки
	DEBUG LogLevel = iota // 0
	// INFO - уровень логирования для информационных сообщений
	INFO // 1
	// ERROR - уровень логирования для сообщений об ошибках
	ERROR // 2
)

// Logger - интерфейс обработчика логов
type Logger interface {
	// SetNext - устанавливает следующий логгер в цепочке
	SetNext(logger Logger) Logger
	// Log - выполняет логирование на основании уровня
	Log(level LogLevel, message string)
}

// BaseLogger - базовая структура для логгера, содержащая ссылку на следующий логгер
type BaseLogger struct {
	next Logger
}

// SetNext - установка следующего логгера в цепочке
func (l *BaseLogger) SetNext(logger Logger) Logger {
	l.next = logger
	return logger
}

// Log - передача запроса следующему логгеру, если текущий не может его обработать
func (l *BaseLogger) Log(level LogLevel, message string) {
	if l.next != nil {
		l.next.Log(level, message)
	}
}

// DebugLogger - логгер для уровня DEBUG
type DebugLogger struct {
	BaseLogger
}

// Log - логирование на уровне DEBUG, если уровень соответствует, иначе передача запроса дальше
func (d *DebugLogger) Log(level LogLevel, message string) {
	if level == DEBUG {
		fmt.Printf("[DEBUG]: %s\n", strings.ToUpper(message))
	} else {
		d.BaseLogger.Log(level, message)
	}
}

// InfoLogger - логгер для уровня INFO
type InfoLogger struct {
	BaseLogger
}

// Log - логирование на уровне INFO, если уровень соответствует, иначе передача запроса дальше
func (i *InfoLogger) Log(level LogLevel, message string) {
	if level == INFO {
		fmt.Printf("[INFO]: %s\n", message)
	} else {
		i.BaseLogger.Log(level, message)
	}
}

// ErrorLogger - логгер для уровня ERROR
type ErrorLogger struct {
	BaseLogger
}

// Log - логирование на уровне ERROR, если уровень соответствует, иначе передача запроса дальше
func (e *ErrorLogger) Log(level LogLevel, message string) {
	if level == ERROR {
		fmt.Printf("[ERROR]: %s\n", message)
	} else {
		e.BaseLogger.Log(level, message)
	}
}

func main() {
	// Создаем логгеры для каждого уровня
	debugLogger := &DebugLogger{}
	infoLogger := &InfoLogger{}
	errorLogger := &ErrorLogger{}

	// Строим цепочку логгеров: Debug -> Info -> Error
	debugLogger.SetNext(infoLogger).SetNext(errorLogger)

	// Логируем сообщения разного уровня
	fmt.Println("Logging at different levels:")

	debugLogger.Log(DEBUG, "This is a debug message")  // Логируется на уровне DEBUG
	debugLogger.Log(INFO, "This is an info message")   // Логируется на уровне INFO
	debugLogger.Log(ERROR, "This is an error message") // Логируется на уровне ERROR

	// Сообщение, не соответствующее уровню
	debugLogger.Log(INFO, "User has logged in")
	debugLogger.Log(ERROR, "Server is down!")
}
