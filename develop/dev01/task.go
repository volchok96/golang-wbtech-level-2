package main

/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/beevik/ntp"
)

// Функция для получения точного времени с обработкой ошибок
func getExactTime() (time.Time, error) {
	exactTime, err := ntp.Time("pool.ntp.org")
	if err != nil {
		return time.Time{}, err
	}
	return exactTime, nil
}

func main() {
	// Получаем текущее время
	currentTime := time.Now()
	fmt.Println("Current time:", currentTime.Format(time.RFC1123))

	// Получаем точное время с использованием функции
	exactTime, err := getExactTime()
	if err != nil {
		// Выводим ошибку в STDERR и завершаем программу с ненулевым кодом
		log.Printf("Error fetching NTP time: %v\n", err)
		os.Exit(1)
	}

	// Выводим точное время
	fmt.Println("Exact time:", exactTime.Format(time.RFC1123))
}
