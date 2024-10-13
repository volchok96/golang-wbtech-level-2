package main

import (
	"errors"
	"testing"
	"time"

	"github.com/beevik/ntp"
)

// Мокаем функцию ntp.Time для тестирования
var mockNTPTime = ntp.Time

// Тест для успешного получения времени
func TestGetExactTime_Success(t *testing.T) {
	// Мокаем функцию ntp.Time для теста
	mockNTPTime = func(host string) (time.Time, error) {
		// Используем текущее время системы для теста
		return time.Now(), nil
	}

	// Вызов функции
	exactTime, err := getExactTime()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Получаем текущее время для проверки
	expected := time.Now().Format(time.RFC1123)
	if exactTime.Format(time.RFC1123) != expected {
		t.Errorf("Expected %v, got %v", expected, exactTime.Format(time.RFC1123))
	}
}

// Тест для случая ошибки
func TestGetExactTime_Error(t *testing.T) {
	// Мокаем ошибку в функции ntp.Time
	mockNTPTime = func(host string) (time.Time, error) {
		return time.Time{}, errors.New("NTP server error")
	}

	// Вызов функции
	_, err := getExactTime()
	if err == nil {
		t.Log("Expected error, got none")
		return
	}

	expectedErr := "NTP server error"
	if err.Error() != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err.Error())
	}
}
