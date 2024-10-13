package main

import (
	"reflect"
	"sort"
	"testing"
)

// Тест для проверки функции удаления дубликатов
func TestRemoveDuplicates(t *testing.T) {
	lines := []string{
		"January 12",
		"April 45K",
		"August 7",
		"January 12",
		"February 1.5M",
	}

	expected := []string{
		"January 12",
		"April 45K",
		"August 7",
		"February 1.5M",
	}

	result := removeDuplicates(lines)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("removeDuplicates() = %v, expected %v", result, expected)
	}
}

// Тест для функции monthIndex
func TestMonthIndex(t *testing.T) {
	tests := []struct {
		line     string
		expected int
	}{
		{"January 12", 1},
		{"Feb 45", 2},
		{"March 7", 3},
		{"December 90", 12},
		{"NotAMonth 45", 0},
	}

	for _, test := range tests {
		result := monthIndex(test.line)
		if result != test.expected {
			t.Errorf("monthIndex(%s) = %d, expected %d", test.line, result, test.expected)
		}
	}
}

// Тест для функции compareLines
func TestCompareLines(t *testing.T) {
	line1 := "August 7"
	line2 := "December 90"

	// Сравнение без числовой и других сортировок
	if compareLines(line1, line2, 1, false, false, false) >= 0 {
		t.Errorf("compareLines() failed for normal comparison between %s and %s", line1, line2)
	}

	// Сравнение по числовому значению
	line1Numeric := "100"
	line2Numeric := "200"
	if compareLines(line1Numeric, line2Numeric, 1, true, false, false) >= 0 {
		t.Errorf("compareLines() failed for numeric comparison between %s and %s", line1Numeric, line2Numeric)
	}
}

// Тест для проверки функции сортировки по месяцам
func TestMonthSort(t *testing.T) {
	lines := []string{
		"December 90",
		"January 12",
		"August 7",
		"February 1.5M",
	}

	expected := []string{
		"January 12",
		"February 1.5M",
		"August 7",
		"December 90",
	}

	sort.Slice(lines, func(i, j int) bool {
		return monthIndex(lines[i]) < monthIndex(lines[j])
	})

	if !reflect.DeepEqual(lines, expected) {
		t.Errorf("monthSort() = %v, expected %v", lines, expected)
	}
}

// Тест для функции isSorted
func TestIsSorted(t *testing.T) {
	linesSortedByMonth := []string{
		"January 12",
		"February 1.5M",
		"August 7",
		"December 90",
	}

	linesUnsorted := []string{
		"December 90",
		"January 12",
		"August 7",
		"February 1.5M",
	}

	// Проверка сортировки по месяцам
	if !isSorted(linesSortedByMonth, 1, false, false, false, true) {
		t.Errorf("isSorted() failed, should be sorted by month")
	}

	// Проверка несортированных данных
	if isSorted(linesUnsorted, 1, false, false, false, true) {
		t.Errorf("isSorted() failed, should not be sorted by month")
	}
}

// Тест для функции parseHumanReadable
func TestParseHumanReadable(t *testing.T) {
	tests := []struct {
		value    string
		expected int
	}{
		{"1K", 1000},
		{"1M", 1000000},
		{"1G", 1000000000},
		{"123", 123},
		{"2k", 2000},
	}

	for _, test := range tests {
		result := parseHumanReadable(test.value)
		if result != test.expected {
			t.Errorf("parseHumanReadable(%s) = %d, expected %d", test.value, result, test.expected)
		}
	}
}
