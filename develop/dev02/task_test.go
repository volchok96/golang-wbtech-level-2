package main

import (
	"testing"
)

func TestUnpackString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		err      bool
	}{
		{"a4bc2d5e", "aaaabccddddde", false},
		{"abcd", "abcd", false},
		{"45", "", true},                  // Неверный формат
		{"", "", false},                   // Пустая строка
		{"qwe\\4\\5", "qwe45", false},     // Экранирование цифр
		{"qwe\\\\5", "qwe\\\\\\\\\\", false},  // Экранирование слеша
	}

	for _, test := range tests {
		result, err := unpackString(test.input)
		if test.err {
			if err == nil {
				t.Logf("Expected error for input %s, but got nil", test.input)
			} else {
				t.Logf("Correctly received error for input %s: %v", test.input, err)
			}
		} else {
			if err != nil {
				t.Logf("Unexpected error for input %s: %v", test.input, err)
			} else if result != test.expected {
				t.Logf("For input %s, expected %s, but got %s", test.input, test.expected, result)
			}
		}
	}
}
