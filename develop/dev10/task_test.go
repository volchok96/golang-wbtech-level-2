package main

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func TestTelnetConnectionSimplified(t *testing.T) {
	// Создаем тестовый TCP сервер (эхо-сервер)
	listener, err := net.Listen("tcp", "127.0.0.1:0") // Используем порт 0 для автоматического выбора свободного порта
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	defer listener.Close()

	serverAddr := listener.Addr().String()

	// Запускаем сервер в отдельной горутине
	go func() {
		conn, err := listener.Accept() // Принимаем входящее соединение
		if err != nil {
			return
		}
		defer conn.Close()
		io.Copy(conn, conn) // Эхо-сервер, который просто возвращает то, что получает
	}()

	// Запускаем клиента telnet
	conn, err := net.DialTimeout("tcp", serverAddr, 5*time.Second)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Пишем сообщение на сервер
	testMessage := "Hello, Telnet!"
	fmt.Fprintln(conn, testMessage)

	// Читаем ответ от сервера
	buf := make([]byte, len(testMessage)+1)
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("Failed to read from connection: %v", err)
	}

	receivedMessage := string(buf[:n])
	if receivedMessage != testMessage+"\n" {
		t.Errorf("Expected %q, got %q", testMessage+"\n", receivedMessage)
	}

	fmt.Println("Test passed")
}
