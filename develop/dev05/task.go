package main

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

type grepOptions struct {
	after      int
	before     int
	context    int
	count      bool
	ignoreCase bool
	invert     bool
	fixed      bool
	lineNum    bool
}

func parseFlags() (grepOptions, string) {
	options := grepOptions{}

	// Определение флагов командной строки
	flag.IntVar(&options.after, "A", 0, "Show N lines *after* each matching line")
	flag.IntVar(&options.before, "B", 0, "Show N lines *before* each matching line")
	flag.IntVar(&options.context, "C", 0, "Show N lines around (both before and after) each match")
	flag.BoolVar(&options.count, "c", false, "Print only the total count of matching lines")
	flag.BoolVar(&options.ignoreCase, "i", false, "Ignore case when matching")
	flag.BoolVar(&options.invert, "v", false, "Invert the match: show lines that do not match the pattern")
	flag.BoolVar(&options.fixed, "F", false, "Interpret the pattern as a fixed string (not a regular expression)")
	flag.BoolVar(&options.lineNum, "n", false, "Show line numbers alongside matching lines")

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Error: Pattern must be provided.")
		fmt.Fprintln(os.Stderr, "Usage: grep [options] pattern [file]")
		os.Exit(2)
	}

	pattern := flag.Arg(0)
	return options, pattern
}

func compilePattern(pattern string, options grepOptions) *regexp.Regexp {
	// Интерпретация шаблона как точной строки, а не регулярного выражения
	if options.fixed {
		pattern = regexp.QuoteMeta(pattern)
	}

	flags := ""
	if options.ignoreCase {
		flags = "(?i)"
	}

	return regexp.MustCompile(flags + pattern)
}

func matchLine(pattern *regexp.Regexp, line string, options grepOptions) bool {
	return pattern.MatchString(line) != options.invert
}

func printContext(lines []string, start, end int, printed map[int]bool, lineNum bool) {
	for i := start; i <= end; i++ {
		if i >= 0 && i < len(lines) && !printed[i] { // Проверка, что строка еще не была напечатана
			printed[i] = true
			if lineNum {
				fmt.Printf("%d: ", i+1)
			}
			fmt.Println(lines[i])
		}
	}
}

func grepFile(filename string, pattern *regexp.Regexp, options grepOptions) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open file %s: %v\n", filename, err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading file %s: %v\n", filename, err)
		os.Exit(1)
	}

	matchingLines := 0
	printed := make(map[int]bool) // Мапа для отслеживания напечатанных строк

	for i, line := range lines {
		if matchLine(pattern, line, options) {
			matchingLines++
			if options.count {
				continue
			}
			// Определение диапазона контекста
			before := options.before
			after := options.after
			if options.context > 0 {
				before = options.context
				after = options.context
			}

			// Вывод контекста перед совпадением
			printContext(lines, i-before, i-1, printed, options.lineNum)

			// Вывод совпадающей строки
			if !printed[i] {
				printed[i] = true
				if options.lineNum {
					fmt.Printf("%d: ", i+1)
				}
				fmt.Println(line)
			}

			// Вывод контекста после совпадения
			printContext(lines, i+1, i+after, printed, options.lineNum)
		}
	}

	if options.count {
		fmt.Println(matchingLines)
	}
}

func main() {
	options, pattern := parseFlags()

	patternRegex := compilePattern(pattern, options)

	files := flag.Args()[1:]

	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "Error: At least one file must be provided.")
		fmt.Fprintln(os.Stderr, "Usage: grep [options] pattern [file]")
		os.Exit(2)
	}

	for _, filename := range files {
		grepFile(filename, patternRegex, options)
	}
}
