package pattern

import "fmt"

/*
Паттерн «Состояние» (State) используется для управления поведением объекта в зависимости от его текущего состояния. Этот пример демонстрирует управление состояниями подключения к базе данных.

Применимость:
- Когда требуется четкое управление переходами между состояниями объекта.
- Когда объект может находиться в различных состояниях с разными правилами поведения.

Плюсы:
- Объект меняет поведение динамически в зависимости от текущего состояния.
- Четко изолированы правила для каждого состояния.
- Легко расширяется новыми состояниями.

Минусы:
- Множество состояний могут усложнить реализацию и поддержку.
- Повышается количество классов/структур, что увеличивает сложность системы.

Применение на практике:
- Управление подключением к базе данных или другим удалённым сервисам.
- Управление сессиями пользователя в веб-приложениях.
*/

// DBConnectionState - интерфейс для состояний подключения
type DBConnectionState interface {
	Connect(conn *DBConnection)
	Disconnect(conn *DBConnection)
	Query(conn *DBConnection, query string) error
}

// DBConnection - контекст для управления состояниями подключения
type DBConnection struct {
	State DBConnectionState
}

// SetState - устанавливает текущее состояние подключения
func (conn *DBConnection) SetState(state DBConnectionState) {
	conn.State = state
}

// Connect - вызывает метод подключения в текущем состоянии
func (conn *DBConnection) Connect() {
	conn.State.Connect(conn)
}

// Disconnect - вызывает метод отключения в текущем состоянии
func (conn *DBConnection) Disconnect() {
	conn.State.Disconnect(conn)
}

// Query - выполняет запрос в зависимости от текущего состояния подключения
func (conn *DBConnection) Query(query string) error {
	return conn.State.Query(conn, query)
}

// DisconnectedState - состояние, когда соединение разорвано
type DisconnectedState struct{}

// Connect - устанавливает соединение при переходе из состояния "Отключено"
func (s *DisconnectedState) Connect(conn *DBConnection) {
	fmt.Println("Соединение установлено.")
	conn.SetState(&ConnectedState{})
}

// Disconnect - уведомляет, что соединение уже разорвано
func (s *DisconnectedState) Disconnect(conn *DBConnection) {
	fmt.Println("Соединение разорвано.")
}

// Query - возвращает ошибку при попытке выполнить запрос без активного соединения
func (s *DisconnectedState) Query(conn *DBConnection, query string) error {
	fmt.Println("Ошибка: нет активного соединения.")
	return fmt.Errorf("нет активного соединения")
}

// ConnectedState - состояние, когда соединение активно
type ConnectedState struct{}

// Connect - уведомляет, что соединение уже активно
func (s *ConnectedState) Connect(conn *DBConnection) {
	fmt.Println("Соединение уже активно.")
}

// Disconnect - разрывает активное соединение
func (s *ConnectedState) Disconnect(conn *DBConnection) {
	fmt.Println("Соединение разорвано.")
	conn.SetState(&DisconnectedState{})
}

// Query - выполняет запрос, если соединение активно
func (s *ConnectedState) Query(conn *DBConnection, query string) error {
	fmt.Printf("Выполнение запроса: %s\n", query)
	return nil
}

// ErrorState - состояние, когда подключение находится в ошибке
type ErrorState struct{}

// Connect - уведомляет о невозможности установить соединение из-за ошибки
func (s *ErrorState) Connect(conn *DBConnection) {
	fmt.Println("Невозможно установить соединение: ошибка.")
}

// Disconnect - разрывает соединение после возникновения ошибки
func (s *ErrorState) Disconnect(conn *DBConnection) {
	fmt.Println("Соединение разорвано после ошибки.")
	conn.SetState(&DisconnectedState{})
}

// Query - возвращает ошибку при попытке выполнить запрос в состоянии ошибки
func (s *ErrorState) Query(conn *DBConnection, query string) error {
	fmt.Println("Ошибка в соединении.")
	return fmt.Errorf("ошибка соединения")
}

func main08() {
	// Создаем объект подключения с начальным состоянием "Отключено"
	conn := &DBConnection{State: &DisconnectedState{}}

	// Пробуем выполнить запрос без подключения
	conn.Query("SELECT * FROM users")

	// Подключаемся к базе данных
	conn.Connect()

	// Выполняем запрос в состоянии "Подключено"
	conn.Query("SELECT * FROM users")

	// Отключаемся от базы данных
	conn.Disconnect()

	// Пытаемся снова выполнить запрос после отключения
	conn.Query("SELECT * FROM users")

	// Симулируем ошибку в соединении
	conn.SetState(&ErrorState{})
	conn.Query("SELECT * FROM users")

	// Пробуем отключиться в состоянии ошибки
	conn.Disconnect()
}
