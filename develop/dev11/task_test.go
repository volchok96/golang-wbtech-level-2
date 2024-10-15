package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func initTestStore() {
	eventsStore.Lock()
	eventsStore.events = make(map[string]Event)
	eventsStore.Unlock()
}

// Тест для создания события
func TestCreateEvent(t *testing.T) {
	initTestStore()

	event := Event{
		Title:     "Test Event",
		UserID:    "1",
		StartTime: time.Now().Add(1 * time.Hour),
		EndTime:   time.Now().Add(2 * time.Hour),
	}

	form := url.Values{}
	form.Set("title", event.Title)
	form.Set("user_id", event.UserID)
	form.Set("start_time", event.StartTime.Format(time.RFC3339))
	form.Set("end_time", event.EndTime.Format(time.RFC3339))

	req, err := http.NewRequest("POST", "/create_event", bytes.NewBufferString(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createEventHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var createdEvent Event
	err = json.Unmarshal(rr.Body.Bytes(), &createdEvent)
	if err != nil {
		t.Fatalf("Could not parse response: %v", err)
	}

	if createdEvent.Title != event.Title || createdEvent.UserID != event.UserID {
		t.Errorf("Unexpected event data: got %v want %v", createdEvent, event)
	}
}

// Тест для обновления события
func TestUpdateEvent(t *testing.T) {
	initTestStore()

	// Создадим событие для обновления
	event := Event{
		ID:        "1",
		Title:     "Test Event",
		UserID:    "1",
		StartTime: time.Now().Add(1 * time.Hour),
		EndTime:   time.Now().Add(2 * time.Hour),
	}
	eventsStore.Lock()
	eventsStore.events[event.ID] = event
	eventsStore.Unlock()

	updatedEvent := Event{
		Title:     "Updated Event",
		UserID:    "1",
		StartTime: time.Now().Add(3 * time.Hour),
		EndTime:   time.Now().Add(4 * time.Hour),
	}
	form := url.Values{}
	form.Set("title", updatedEvent.Title)
	form.Set("user_id", updatedEvent.UserID)
	form.Set("start_time", updatedEvent.StartTime.Format(time.RFC3339))
	form.Set("end_time", updatedEvent.EndTime.Format(time.RFC3339))
	form.Set("id", event.ID)

	req, err := http.NewRequest("POST", "/update_event", bytes.NewBufferString(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(updateEventHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var updatedEventResp Event
	err = json.Unmarshal(rr.Body.Bytes(), &updatedEventResp)
	if err != nil {
		t.Fatalf("Could not parse response: %v", err)
	}

	if updatedEventResp.Title != updatedEvent.Title {
		t.Errorf("Unexpected event data: got %v want %v", updatedEventResp, updatedEvent)
	}
}

// Тест для удаления события
func TestDeleteEvent(t *testing.T) {
	// Инициализируем хранилище
	initTestStore()

	// Создадим событие для последующего удаления
	event := Event{
		ID:        "1",
		Title:     "Test Event",
		UserID:    "1",
		StartTime: time.Now().Add(1 * time.Hour),
		EndTime:   time.Now().Add(2 * time.Hour),
	}

	// Добавляем событие в хранилище
	eventsStore.Lock()
	eventsStore.events[event.ID] = event
	eventsStore.Unlock()

	// Инициализируем форму для передачи ID события
	form := url.Values{}
	form.Set("id", event.ID)

	// Создаем запрос на удаление события
	req, err := http.NewRequest("POST", "/delete_event", bytes.NewBufferString(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteEventHandler)
	handler.ServeHTTP(rr, req)

	// Проверяем код ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверяем тело ответа
	expected := `{"result":"event deleted"}`
	if rr.Body.String() != expected {
		t.Errorf("Unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Проверяем, что событие действительно удалено
	eventsStore.RLock()
	defer eventsStore.RUnlock()
	if _, exists := eventsStore.events[event.ID]; exists {
		t.Errorf("Event was not deleted from the store")
	}
}

// Тест для получения событий за день
func TestGetEventsForDay(t *testing.T) {
	// Инициализируем хранилище
	initTestStore()

	// Установим дату для теста
	testDate := time.Date(2024, 10, 15, 0, 0, 0, 0, time.UTC)

	// Создадим два события, которые должны попасть в этот день
	event1 := Event{
		ID:        "1",
		Title:     "Morning Meeting",
		UserID:    "1",
		StartTime: testDate.Add(9 * time.Hour),  // 09:00
		EndTime:   testDate.Add(10 * time.Hour), // 10:00
	}

	event2 := Event{
		ID:        "2",
		Title:     "Evening Workout",
		UserID:    "1",
		StartTime: testDate.Add(18 * time.Hour), // 18:00
		EndTime:   testDate.Add(19 * time.Hour), // 19:00
	}

	// Добавляем события в хранилище
	eventsStore.Lock()
	eventsStore.events[event1.ID] = event1
	eventsStore.events[event2.ID] = event2
	eventsStore.Unlock()

	// Создаем запрос на получение событий за день
	req, err := http.NewRequest("GET", fmt.Sprintf("/events_for_day?user_id=1&date=%s", testDate.Format("2006-01-02")), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(eventsForDayHandler)
	handler.ServeHTTP(rr, req)

	// Проверяем код ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверяем полученные события
	var events []Event
	err = json.NewDecoder(rr.Body).Decode(&events)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(events))
	}
}
