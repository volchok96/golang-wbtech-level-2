package main

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
Методы API: POST /create_event POST /update_event POST /delete_event GET /events_for_day GET /events_for_week GET /events_for_month
Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Event представляет событие в календаре, содержащее ID, заголовок, ID пользователя, время начала и окончания.
type Event struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	UserID    string    `json:"user_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// JSONResponse представляет стандартную структуру ответа с результатом или сообщением об ошибке.
type JSONResponse struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

// Xранилище в памяти с RW мьютексами для потокобезопасности
var eventsStore = struct {
	sync.RWMutex
	events map[string]Event
}{events: make(map[string]Event)}

// ================== Бизнес-логика ==================

// Создание события
func createEvent(event Event) (Event, error) {
    // Если ID не задан, то генерируем его
    if event.ID == "" {
        event.ID = fmt.Sprintf("%d", time.Now().UnixNano())
    }

    eventsStore.Lock()
    eventsStore.events[event.ID] = event
    eventsStore.Unlock()

    return event, nil
}


// Обновление события
func updateEvent(id string, updated Event) (Event, error) {
	eventsStore.Lock()
	defer eventsStore.Unlock()

	event, exists := eventsStore.events[id]
	if !exists {
		return Event{}, fmt.Errorf("event not found")
	}

	event.Title = updated.Title
	event.UserID = updated.UserID
	event.StartTime = updated.StartTime
	event.EndTime = updated.EndTime

	eventsStore.events[id] = event
	return event, nil
}

// Удаление события
func deleteEvent(id string) error {
	eventsStore.Lock()
	defer eventsStore.Unlock()

	if _, exists := eventsStore.events[id]; !exists {
		return fmt.Errorf("event not found")
	}
	delete(eventsStore.events, id)
	return nil
}

// Получение событий за период (день, неделя, месяц)
func getEventsForPeriod(userID string, start, end time.Time) []Event {
	eventsStore.RLock()
	defer eventsStore.RUnlock()

	var result []Event
	for _, event := range eventsStore.events {
		if event.UserID == userID && event.StartTime.After(start) && event.StartTime.Before(end) {
			result = append(result, event)
		}
	}
	return result
}

// ================== HTTP Handlers ==================

// Middleware для логирования запросов
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
		next.ServeHTTP(w, r)
	})
}

// Создание события
func createEventHandler(w http.ResponseWriter, r *http.Request) {
	event, err := parseAndValidateEvent(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusBadRequest)
		return
	}

	createdEvent, err := createEvent(event)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusServiceUnavailable)
		return
	}

	response, _ := toJSON(createdEvent)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

// Обновление события
func updateEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseAndValidateID(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusBadRequest)
		return
	}

	updatedEvent, err := parseAndValidateEvent(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusBadRequest)
		return
	}

	event, err := updateEvent(id, updatedEvent)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusNotFound)
		return
	}

	response, _ := toJSON(event)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Удаление события
func deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseAndValidateID(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusBadRequest)
		return
	}

	err = deleteEvent(id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusNotFound)
		return
	}

	response, _ := toJSON(JSONResponse{Result: "event deleted"})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Получение событий за день
func eventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	userID, date, err := parseAndValidateUserIDAndDate(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusBadRequest)
		return
	}

	startOfDay := date
	endOfDay := date.Add(24 * time.Hour)

	events := getEventsForPeriod(userID, startOfDay, endOfDay)
	response, _ := toJSON(events)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Получение событий за неделю
func eventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	userID, date, err := parseAndValidateUserIDAndDate(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusBadRequest)
		return
	}

	startOfWeek := date
	endOfWeek := date.Add(7 * 24 * time.Hour)

	events := getEventsForPeriod(userID, startOfWeek, endOfWeek)
	response, _ := toJSON(events)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Получение событий за месяц
func eventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	userID, date, err := parseAndValidateUserIDAndDate(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusBadRequest)
		return
	}

	startOfMonth := date
	endOfMonth := date.Add(30 * 24 * time.Hour)

	events := getEventsForPeriod(userID, startOfMonth, endOfMonth)
	response, _ := toJSON(events)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// ================== Вспомогательные функции ==================

// Преобразование структуры в JSON
func toJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Парсинг ID события из запроса
func parseAndValidateID(r *http.Request) (string, error) {
	id := r.FormValue("id")
	if id == "" {
		return "", fmt.Errorf("missing id")
	}
	return id, nil
}

// Парсинг и валидация данных события из запроса
func parseAndValidateEvent(r *http.Request) (Event, error) {
	var event Event
	if err := r.ParseForm(); err != nil {
		return event, fmt.Errorf("invalid request body: %v", err)
	}

	event.Title = r.FormValue("title")
	event.UserID = r.FormValue("user_id")
	startTimeStr := r.FormValue("start_time")
	endTimeStr := r.FormValue("end_time")

	if event.Title == "" || event.UserID == "" || startTimeStr == "" || endTimeStr == "" {
		return event, fmt.Errorf("missing required fields")
	}

	var err error
	event.StartTime, err = time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		return event, fmt.Errorf("invalid start_time: %v", err)
	}

	event.EndTime, err = time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		return event, fmt.Errorf("invalid end_time: %v", err)
	}

	if event.EndTime.Before(event.StartTime) {
		return event, fmt.Errorf("end_time cannot be before start_time")
	}

	return event, nil
}

// Парсинг user_id и даты из запроса
func parseAndValidateUserIDAndDate(r *http.Request) (string, time.Time, error) {
	userID := r.URL.Query().Get("user_id")
	dateStr := r.URL.Query().Get("date")

	if userID == "" || dateStr == "" {
		return "", time.Time{}, fmt.Errorf("missing user_id or date")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid date format")
	}

	return userID, date, nil
}

// ================== Запуск сервера с graceful shutdown ==================
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/create_event", createEventHandler)
	mux.HandleFunc("/update_event", updateEventHandler)
	mux.HandleFunc("/delete_event", deleteEventHandler)
	mux.HandleFunc("/events_for_day", eventsForDayHandler)
	mux.HandleFunc("/events_for_week", eventsForWeekHandler)
	mux.HandleFunc("/events_for_month", eventsForMonthHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	loggedMux := loggingMiddleware(mux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: loggedMux,
	}

	// Канал для получения сигнала о завершении работы
	idleConnsClosed := make(chan struct{})

	// Обработчик системных сигналов для shutdown
	go func() {
		// Перехватываем системные сигналы
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

		<-sigint // Ожидаем сигнал

		// Начинаем shutdown сервера с тайм-аутом
		log.Println("Shutting down server...")

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed) // Сообщаем, что соединения закрыты
	}()

	log.Printf("Starting server on port %s", port)

	// Запуск сервера
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Could not start server: %v", err)
	}

	// Ожидание завершения всех соединений
	<-idleConnsClosed
	log.Println("Server gracefully stopped")
}
