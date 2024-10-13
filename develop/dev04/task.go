package main

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