package main

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Создаем флаг для таймаута соединения
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")

	// Пропускаем подкоманду (go-telnet) и разбираем флаги
	if len(os.Args) < 2 || os.Args[1] != "go-telnet" {
		fmt.Println("Usage: go-telnet --timeout=10s host port")
		return
	}

	// Парсим флаги, начиная с 2-го аргумента
	flag.CommandLine.Parse(os.Args[2:])

	// Проверяем, что указаны хост и порт
	if flag.NArg() < 2 {
		fmt.Println("Usage: go-telnet --timeout=10s host port")
		return
	}

	// Извлекаем хост и порт из аргументов
	host := flag.Arg(0)
	port := flag.Arg(1)

	// Формируем адрес для подключения
	address := net.JoinHostPort(host, port)

	// Запускаем клиента Telnet
	runTelnetClient(address, *timeout)
}

func runTelnetClient(address string, timeout time.Duration) {
	// Устанавливаем TCP соединение с указанным адресом
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close() // Закрываем соединение при выходе из функции
	fmt.Printf("Connected to %s\n", address)

	// Создаем каналы для сигнализации о завершении работы
	doneReading := make(chan struct{})
	doneWriting := make(chan struct{})

	// Канал для перехвата системных сигналов (например, Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Горутина для чтения данных из соединения и вывода их в стандартный вывод
	go func() {
		defer close(doneReading) // Закрываем канал после завершения чтения
		if _, err := io.Copy(os.Stdout, conn); err != nil {
			fmt.Println("Error reading from connection:", err)
		}
	}()

	// Горутина для записи данных из стандартного ввода в соединение
	go func() {
		defer close(doneWriting) // Закрываем канал после завершения записи
		if _, err := io.Copy(conn, os.Stdin); err != nil {
			fmt.Println("Error writing to connection:", err)
		}
	}()

	// Ожидаем завершения одной из горутин или получения сигнала
	select {
	case <-doneReading:
		fmt.Println("Server closed the connection.")
	case <-doneWriting:
		fmt.Println("Connection closed by client.")
	case sig := <-sigChan:
		fmt.Printf("Received signal: %s, closing connection.\n", sig)
	}

	// Закрываем соединение и завершаем программу
	conn.Close()
	fmt.Println("Telnet client terminated.")
}
