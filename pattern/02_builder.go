package pattern

/*
Паттерн «Строитель» (Builder) используется для пошагового построения сложных объектов. 
Строитель разделяет процесс создания объекта на отдельные шаги (или части), которые можно вызывать поочередно, изменять или комбинировать. 
Это позволяет создавать различные представления одного и того же объекта с разными комбинациями данных или конфигураций.

Применимость:
- Когда объект имеет много опциональных или сложных параметров, и его создание без паттерна может вызвать громоздкие и запутанные конструкторы.
- Когда нужно отделить логику создания объекта от его структуры или представления.
- Когда важно иметь возможность создавать разные вариации объекта с различными конфигурациями.

Плюсы:
- Позволяет разделить логику построения объекта на отдельные шаги, делая код более читабельным.
- Упрощает создание сложных объектов с множеством конфигураций.
- Позволяет использовать один и тот же процесс построения для создания различных представлений объекта.
- Код становится более гибким и поддерживаемым.

Минусы:
- Увеличивает сложность системы, если объект простой и не требует множества шагов для создания.
- Требует создания дополнительного класса для строителя, что может оказаться избыточным в простых случаях.

Применение на практике:
- Создание объектов с множеством опциональных параметров. Например, построение запроса HTTP или объекта конфигурации.
- Генерация сложных документов (например, PDF, HTML), где каждый шаг добавляет разные элементы к документу.
- Построение объектов в играх (например, создание персонажей или уровней, где есть много опций и параметров).
*/

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

// HTTPRequest представляет структуру HTTP-запроса.
type HTTPRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
	Timeout time.Duration
}

// HTTPRequestBuilder определяет интерфейс для построения HTTP-запросов.
type HTTPRequestBuilder interface {
	SetMethod(method string) HTTPRequestBuilder
	SetURL(url string) HTTPRequestBuilder
	AddHeader(key, value string) HTTPRequestBuilder
	SetBody(body []byte) HTTPRequestBuilder
	SetTimeout(timeout time.Duration) HTTPRequestBuilder
	Build() (*HTTPRequest, error)
}

// HTTPRequestBuilderImpl представляет конкретного строителя HTTP-запросов.
type HTTPRequestBuilderImpl struct {
	request HTTPRequest
}

// NewHTTPRequestBuilder возвращает новый экземпляр строителя HTTP-запросов.
func NewHTTPRequestBuilder() HTTPRequestBuilder {
	return &HTTPRequestBuilderImpl{
		request: HTTPRequest{
			Headers: make(map[string]string), // Инициализация пустых заголовков
		},
	}
}

// SetMethod устанавливает HTTP-метод для запроса.
func (b *HTTPRequestBuilderImpl) SetMethod(method string) HTTPRequestBuilder {
	b.request.Method = method
	return b
}

// SetURL устанавливает URL для HTTP-запроса.
func (b *HTTPRequestBuilderImpl) SetURL(url string) HTTPRequestBuilder {
	b.request.URL = url
	return b
}

// AddHeader добавляет заголовок к HTTP-запросу.
func (b *HTTPRequestBuilderImpl) AddHeader(key, value string) HTTPRequestBuilder {
	b.request.Headers[key] = value
	return b
}

// SetBody устанавливает тело HTTP-запроса.
func (b *HTTPRequestBuilderImpl) SetBody(body []byte) HTTPRequestBuilder {
	b.request.Body = body
	return b
}

// SetTimeout устанавливает тайм-аут для HTTP-запроса.
func (b *HTTPRequestBuilderImpl) SetTimeout(timeout time.Duration) HTTPRequestBuilder {
	b.request.Timeout = timeout
	return b
}

// Build завершает создание HTTP-запроса и возвращает его.
func (b *HTTPRequestBuilderImpl) Build() (*HTTPRequest, error) {
	// Валидация обязательных параметров
	if b.request.Method == "" {
		return nil, fmt.Errorf("method cannot be empty")
	}
	if b.request.URL == "" {
		return nil, fmt.Errorf("URL cannot be empty")
	}
	return &b.request, nil
}

// Send выполняет отправку HTTP-запроса.
func (r *HTTPRequest) Send() (*http.Response, error) {
	client := &http.Client{Timeout: r.Timeout}
	req, err := http.NewRequest(r.Method, r.URL, bytes.NewBuffer(r.Body))
	if err != nil {
		return nil, err
	}
	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}
	return client.Do(req)
}

func main02() {
	// Создание запроса с помощью строителя
	builder := NewHTTPRequestBuilder()
	request, err := builder.
		SetMethod("POST").
		SetURL("https://example.com/api").
		AddHeader("Content-Type", "application/json").
		SetBody([]byte(`{"message": "Golang is a programming language"}`)).
		SetTimeout(5 * time.Second).
		Build()

	if err != nil {
		fmt.Println("Error building request:", err)
		return
	}

	// Отправка HTTP-запроса
	response, err := request.Send()
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer response.Body.Close()

	fmt.Println("Response Status:", response.Status)
}
