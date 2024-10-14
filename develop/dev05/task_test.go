package main

import (
	"io"
	"os"
	"testing"
)

func TestGrepFile(t *testing.T) {
	// Создаем временный файл для теста
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Записываем тестовые данные в файл
	testData := `This is a test file.
Line with keyword.
Another line without keyword.
Final line with keyword again.`

	if _, err := tempFile.WriteString(testData); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Определяем шаблон для поиска и создаем регулярное выражение
	pattern := "keyword"
	options := grepOptions{after: 1, before: 1, context: 0, count: false, ignoreCase: false, invert: false, fixed: false, lineNum: false}
	patternRegex := compilePattern(pattern, options)

	// Захватываем вывод функции
	output := captureOutput(func() {
		grepFile(tempFile.Name(), patternRegex, options)
	})

	// Ожидаемый вывод
	expected := `This is a test file.
Line with keyword.
Another line without keyword.
Final line with keyword again.
`

	if output != expected {
		t.Errorf("Expected output:\n%s\nBut got:\n%s", expected, output)
	}
}

func TestGrepFileWithCount(t *testing.T) {
	// Создаем временный файл для теста
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Записываем тестовые данные в файл
	testData := `This is a test file.
Line with keyword.
Another line without keyword.
Final line with keyword again.`

	if _, err := tempFile.WriteString(testData); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Определяем шаблон для поиска и создаем регулярное выражение
	pattern := "keyword"
	options := grepOptions{count: true}
	patternRegex := compilePattern(pattern, options)

	// Захватываем вывод функции
	output := captureOutput(func() {
		grepFile(tempFile.Name(), patternRegex, options)
	})

	// Ожидаемый вывод
	expected := "3\n"

	if output != expected {
		t.Errorf("Expected output:\n%s\nBut got:\n%s", expected, output)
	}
}

// Вспомогательная функция для захвата вывода
func captureOutput(f func()) string {
	old := os.Stdout // Сохраняем старый вывод
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = old // Восстанавливаем вывод

	return string(out)
}
