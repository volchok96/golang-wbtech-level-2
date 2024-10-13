package main

import (
	"reflect"
	"testing"
)

func TestFindAnagrams(t *testing.T) {
	tests := []struct {
		words    []string
		expected map[string][]string
	}{
		{
			words: []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик"},
			expected: map[string][]string{
				"пятак":  {"пятак", "пятка", "тяпка"},
				"листок": {"листок", "слиток", "столик"},
			},
		},
		{
			words: []string{"кот", "ток"},
			expected: map[string][]string{
				"кот": {"кот", "ток"},
			},
		},
		{
			words:    []string{"волк", "лиса", "собака"},
			expected: map[string][]string{},
		},
	}

	for _, test := range tests {
		result := findAnagrams(test.words)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("For words %v, expected %v, but got %v", test.words, test.expected, result)
		}
	}
}

func TestSortString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"пятак", "акптя"},
		{"слиток", "иклост"},
		{"листок", "иклост"},
	}

	for _, test := range tests {
		result := sortString(test.input)
		if result != test.expected {
			t.Errorf("For input %s, expected %s, but got %s", test.input, test.expected, result)
		}
	}
}