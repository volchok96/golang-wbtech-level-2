package main

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type cutOptions struct {
	fields    string
	delimiter string
	separated bool
}

func parseFlags() cutOptions {
	options := cutOptions{}

	flag.StringVar(&options.fields, "f", "", "Select fields (columns)")
	flag.StringVar(&options.delimiter, "d", "\t", "Use a different delimiter")
	flag.BoolVar(&options.separated, "s", false, "Only lines with delimiter")

	flag.Parse()

	return options
}

func cutLine(line string, options cutOptions) string {
	if options.separated && !strings.Contains(line, options.delimiter) {
		return ""
	}

	fields := strings.Split(line, options.delimiter)
	selectedFields := parseFieldsOption(options.fields, len(fields))

	var result []string
	for _, field := range selectedFields {
		if field > 0 && field <= len(fields) {
			result = append(result, fields[field-1])
		}
	}

	return strings.Join(result, options.delimiter)
}

func parseFieldsOption(fieldsOption string, numFields int) []int {
	var fields []int
	parts := strings.Split(fieldsOption, ",")
	for _, part := range parts {
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			start, end := parseRange(rangeParts[0], rangeParts[1], numFields)
			for i := start; i <= end; i++ {
				fields = append(fields, i)
			}
		} else {
			field, _ := strconv.Atoi(part)
			fields = append(fields, field)
		}
	}
	return fields
}

func parseRange(startStr, endStr string, numFields int) (int, int) {
	start := 1
	if startStr != "" {
		start, _ = strconv.Atoi(startStr)
	}
	end := numFields
	if endStr != "" {
		end, _ = strconv.Atoi(endStr)
	}
	return start, end
}

func cutLines(input *os.File, options cutOptions) {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		result := cutLine(line, options)
		if result != "" {
			fmt.Println(result)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	options := parseFlags()
	cutLines(os.Stdin, options)
}