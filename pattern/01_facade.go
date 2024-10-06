package pattern

/*
Паттерн "Фасад" скрывает сложные детали работы с HTTP-клиентом, предлагая разработчикам 
интуитивно понятный способ отправки запросов и обработки ответов.

Основные компоненты:

1. HTTPFacade - основной фасад для конфигурирования и отправки HTTP-запросов.
2. NewHTTPFacade - конструктор для создания экземпляра HTTPFacade с установленным таймаутом.
3. Get - метод для отправки GET-запросов на указанный URL и получения ответа.
4. Post - метод для отправки POST-запросов с данными в формате JSON и получения ответа.
5. StartServer - вспомогательная функция для запуска локального HTTP-сервера для тестирования.

Применимость:
- Упрощает взаимодействие с сложными API, скрывая внутреннюю реализацию.
- Позволяет изменять внутреннюю логику без необходимости менять код клиентской части.

Плюсы:
- Упрощает код, уменьшая его связность.
- Улучшает читабельность и поддерживаемость.

Минусы:
- Может привести к созданию излишне абстрактных интерфейсов.
- Добавляет дополнительный уровень абстракции, что может усложнить отладку.

Примеры использования:
- Модуль для работы с API, где множество различных запросов могут быть упакованы в единый интерфейс.
- Веб-приложения, которые должны взаимодействовать с несколькими внешними сервисами.
*/


import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// HTTPFacade - Фасад для конфигурирования HTTP-запросов
type HTTPFacade struct {
	client *http.Client
}

// NewHTTPFacade - Конструктор для создания нового HTTP-фасада
func NewHTTPFacade() *HTTPFacade {
	return &HTTPFacade{
		client: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

// Get - Метод фасада для отправки GET-запроса
func (h *HTTPFacade) Get(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create GET request: %w", err)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute GET request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

// Post - Метод фасада для отправки POST-запроса с данными в формате JSON
func (h *HTTPFacade) Post(url string, data []byte) (string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("failed to create POST request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute POST request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

// StartServer - Запускает локальный HTTP-сервер
func StartServer(wg *sync.WaitGroup) {
	defer wg.Done()
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from GET request!")
	})
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		fmt.Fprintf(w, "Hello from POST request! You sent: %s", string(body))
	})
	http.ListenAndServe(":8080", nil)
}

func main() {
	// Используем WaitGroup для корректного завершения программы после выполнения всех горутин
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// Запускаем сервер в отдельной горутине
	go StartServer(wg)

	time.Sleep(time.Second) // Даем серверу время на старт

	// Создаем фасад для работы с HTTP-запросами
	httpFacade := NewHTTPFacade()

	// Используем фасад для отправки GET-запроса
	getResponse, err := httpFacade.Get("http://localhost:8080/get")
	if err != nil {
		fmt.Println("GET request failed:", err)
	} else {
		fmt.Println("GET Response:", getResponse)
	}

	// Используем фасад для отправки POST-запроса
	postData := []byte(`{"message": "Hello, server!"}`)
	postResponse, err := httpFacade.Post("http://localhost:8080/post", postData)
	if err != nil {
		fmt.Println("POST request failed:", err)
	} else {
		fmt.Println("POST Response:", postResponse)
	}

	wg.Wait() // Ждем завершения сервера
}
