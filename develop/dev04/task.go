package main

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"fmt"
	"sort"
	"strings"
)

// Функция для сортировки букв в строке
func sortString(s string) string {
	r := []rune(s)
	sort.Slice(r, func(i, j int) bool {
		return r[i] < r[j]
	})
	return string(r)
}

// Функция для поиска множеств анаграмм по словарю
func findAnagrams(words []string) map[string][]string {
	// Создание карты для группировки анаграмм
	anagrams := make(map[string][]string)
	wordMap := make(map[string]string) // для отслеживания первого встретившегося слова

	// Приведение всех слов к нижнему регистру
	for _, word := range words {
		lowerWord := strings.ToLower(word)
		sortedWord := sortString(lowerWord)

		if _, found := anagrams[sortedWord]; !found {
			wordMap[sortedWord] = lowerWord // первое встретившееся слово
		}
		anagrams[sortedWord] = append(anagrams[sortedWord], lowerWord)
	}

	// Формирование результирующей карты множеств анаграмм
	result := make(map[string][]string)
	for key, group := range anagrams {
		if len(group) > 1 {
			sort.Strings(group) // сортировка по возрастанию
			result[wordMap[key]] = group
		}
	}

	return result
}

func main() {
	words := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "волчок"}
	anagrams := findAnagrams(words)

	for key, group := range anagrams {
		fmt.Printf("%s: %v\n", key, group)
	}
}