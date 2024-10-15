package main

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

// Функция распаковки строки
func unpackString(str string) (string, error) {
	runes := []rune(str)
	var result []rune
	escape := false // Флаг для эскейп последовательностей

	for i := 0; i < len(runes); i++ {
		if runes[i] == '\\' && !escape { // Если это первый символ обратного слеша
			escape = true
			continue
		}

		// Если это цифра и мы не в режиме экранирования
		if unicode.IsDigit(runes[i]) && !escape {
			if i == 0 || unicode.IsDigit(runes[i-1]) { // Проверяем на некорректную строку
				return "", errors.New("некорректная строка")
			}

			// Конвертируем символ в число
			count, _ := strconv.Atoi(string(runes[i]))
			for j := 0; j < count-1; j++ { // Повторяем предыдущий символ count раз
				result = append(result, runes[i-1])
			}
		} else {
			result = append(result, runes[i]) // Добавляем символ в результат
			escape = false                    // Сбрасываем флаг экранирования
		}
	}

	return string(result), nil
}

func main() {
	// Пример использования
	res, err := unpackString("qwe\\4\\5")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", res)
	}
}
