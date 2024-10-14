package main

import (
	"testing"
	"time"
)

// Вспомогательная функция для создания канала, который закрывается через заданное время
func signalAfter(duration time.Duration) <-chan interface{} {
	signal := make(chan interface{})
	go func() {
		defer close(signal)
		time.Sleep(duration)
	}()
	return signal
}

// Тесты для функции or
func TestOr(t *testing.T) {
	t.Run("No channels", func(t *testing.T) {
		if or() != nil {
			t.Error("Expected nil for no channels")
		}
	})

	t.Run("Single channel", func(t *testing.T) {
		sig := signalAfter(10 * time.Millisecond)
		select {
		case <-or(sig):
			// Ожидаем, что or вернет тот же канал и он закроется
		case <-time.After(20 * time.Millisecond):
			t.Error("Expected to receive from single channel")
		}
	})

	t.Run("Multiple channels, first closes first", func(t *testing.T) {
		start := time.Now()
		<-or(
			signalAfter(10*time.Millisecond),
			signalAfter(5*time.Second),
			signalAfter(1*time.Second),
		)
		if time.Since(start) >= 50*time.Millisecond {
			t.Error("Expected to receive from channel within 10ms")
		}
	})

	t.Run("Multiple channels, middle closes first", func(t *testing.T) {
		start := time.Now()
		<-or(
			signalAfter(100*time.Millisecond),
			signalAfter(5*time.Millisecond),
			signalAfter(200*time.Millisecond),
		)
		if time.Since(start) >= 50*time.Millisecond {
			t.Error("Expected to receive from channel within 5ms")
		}
	})

	t.Run("Multiple channels, last closes first", func(t *testing.T) {
		start := time.Now()
		<-or(
			signalAfter(100*time.Millisecond),
			signalAfter(50*time.Millisecond),
			signalAfter(5*time.Millisecond),
		)
		if time.Since(start) >= 50*time.Millisecond {
			t.Error("Expected to receive from channel within 5ms")
		}
	})
}
