package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	sortByColumnStr      string
	numericSort          bool
	reverseSort          bool
	unique               bool
	monthSort            bool
	ignoreTrailingSpaces bool
	checkSorted          bool
	humanReadableSort    bool
)

func init() {
	// Определение флагов для командной строки
	flag.StringVar(&sortByColumnStr, "k", "1", "Sort by column (default is the first column)")
	flag.BoolVar(&numericSort, "n", false, "Sort by numeric value")
	flag.BoolVar(&reverseSort, "r", false, "Sort in reverse order")
	flag.BoolVar(&unique, "u", false, "Suppress duplicate lines")
	flag.BoolVar(&monthSort, "M", false, "Sort by month names")
	flag.BoolVar(&ignoreTrailingSpaces, "b", false, "Ignore trailing spaces")
	flag.BoolVar(&checkSorted, "c", false, "Check if data is sorted")
	flag.BoolVar(&humanReadableSort, "h", false, "Sort by numeric value considering suffixes")
}

func main() {
	// Проверяем, что первый аргумент — это команда "sort"
	if len(os.Args) < 2 || os.Args[1] != "sort" {
		printUsage()
		os.Exit(1)
	}

	// Создание файла input.txt, если он не существует
	ensureInputFileExists()

	// Парсим флаги, начиная со второго аргумента
	flag.CommandLine.Parse(os.Args[2:])

	// Проверка значения флага -k
	sortByColumn, err := strconv.Atoi(sortByColumnStr)
	if err != nil {
		fmt.Println("Column number was not set. Using default column 1.")
		sortByColumn = 1
	}

	// Получаем имя файла из аргументов командной строки
	var inputFile string
	if len(flag.Args()) > 0 {
		inputFile = flag.Args()[0]
	} else {
		inputFile = "input.txt"
	}

	// Чтение строк из файла
	lines, err := readLines(inputFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	// Проверка отсортированности (флаг -c)
	if checkSorted {
		if isSorted(lines, sortByColumn, numericSort, ignoreTrailingSpaces, humanReadableSort, monthSort) {
			fmt.Println("The data is sorted.")
		} else {
			fmt.Println("The data is not sorted.")
		}
		// Завершаем выполнение программы, если включен режим проверки отсортированности
		return
	}

	// Сортировка данных (даже при флаге -c мы производим запись отсортированных данных)
	if unique {
		lines = removeDuplicates(lines)
	}

	if monthSort {
		sort.Slice(lines, func(i, j int) bool {
			return monthIndex(lines[i]) < monthIndex(lines[j])
		})
	} else {
		sort.Slice(lines, func(i, j int) bool {
			return compareLines(lines[i], lines[j], sortByColumn, numericSort, ignoreTrailingSpaces, humanReadableSort) < 0
		})
	}

	if reverseSort {
		sort.Sort(sort.Reverse(sort.StringSlice(lines)))
	}

	// Запись отсортированных строк в новый файл (всегда перезаписываем данные)
	outputFile := "sorted_" + inputFile
	err = writeLines(outputFile, lines)
	if err != nil {
		fmt.Println("Error writing file:", err)
		os.Exit(1)
	}
	fmt.Println("Sorted data written to", outputFile)
}

// Функция, выводящая информацию о том, как использовать программу
func printUsage() {
	fmt.Println(`Usage: go run <filename>.go sort [options] <inputfile>
Options:
  -k <column>       Sort by column (default is the first column)
  -n                Sort by numeric value
  -r                Sort in reverse order
  -u                Suppress duplicate lines
  -M                Sort by month names
  -b                Ignore trailing spaces
  -c                Check if data is sorted
  -h                Sort by numeric value considering suffixes

Example:
  go run develop/dev03/task.go sort -k 2 -n input.txt
  go run develop/dev03/task.go sort -M input.txt`)
}

// Создание файла input.txt, если он не существует
func ensureInputFileExists() {
	const sampleData = `January 12
April 45K
August 7
December 90
February 1.5M
`
	filename := "input.txt"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		err := os.WriteFile(filename, []byte(sampleData), 0644)
		if err != nil {
			fmt.Println("Error creating file:", err)
		} else {
			fmt.Printf("File %s created with sample data.\n", filename)
		}
	}
}

// Чтение строк из файла
func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// Запись строк в файл
func writeLines(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		writer.WriteString(line + "\n")
	}
	return writer.Flush()
}

// Сравнение строк с учетом параметров
func compareLines(line1, line2 string, column int, numeric bool, ignoreSpaces bool, humanReadable bool) int {
	if ignoreSpaces {
		line1 = strings.TrimSpace(line1)
		line2 = strings.TrimSpace(line2)
	}

	fields1 := strings.Fields(line1)
	fields2 := strings.Fields(line2)

	if column-1 >= len(fields1) || column-1 >= len(fields2) {
		return 0
	}

	if humanReadable {
		num1 := parseHumanReadable(fields1[column-1])
		num2 := parseHumanReadable(fields2[column-1])
		return num1 - num2
	}

	if numeric {
		num1, _ := strconv.Atoi(fields1[column-1])
		num2, _ := strconv.Atoi(fields2[column-1])
		return num1 - num2
	}

	return strings.Compare(fields1[column-1], fields2[column-1])
}

// Обработка числовых значений с суффиксами
func parseHumanReadable(value string) int {
	re := regexp.MustCompile(`^([0-9]+)([KkMmGg]?)$`)
	matches := re.FindStringSubmatch(value)
	if len(matches) == 3 {
		num, _ := strconv.Atoi(matches[1])
		suffix := matches[2]
		switch suffix {
		case "K", "k":
			return num * 1000
		case "M", "m":
			return num * 1000000
		case "G", "g":
			return num * 1000000000
		default:
			return num
		}
	}
	return 0
}

// Удаление дубликатов
func removeDuplicates(lines []string) []string {
	keys := make(map[string]bool)
	var result []string
	for _, line := range lines {
		if _, value := keys[line]; !value {
			keys[line] = true
			result = append(result, line)
		}
	}
	return result
}

// Проверка отсортированных данных
func isSorted(lines []string, column int, numeric bool, ignoreSpaces bool, humanReadable bool, monthSort bool) bool {
	// Проверка на пустой ввод
	if len(lines) == 0 {
		return true
	}

	// Переменные для отслеживания, отсортированы ли данные по месяцам и обычной сортировке
	isMonthSorted := true
	isColumnSorted := true

	for i := 1; i < len(lines); i++ {
		// Проверка сортировки по месяцам
		if monthSort {
			if monthIndex(lines[i-1]) == 0 || monthIndex(lines[i]) == 0 {
				// Пропускаем строки с нераспознанными месяцами
				continue
			}
			if monthIndex(lines[i-1]) > monthIndex(lines[i]) {
				isMonthSorted = false
			}
		}

		// Проверка сортировки по указанному столбцу
		compResult := compareLines(lines[i-1], lines[i], column, numeric, ignoreSpaces, humanReadable)
		if compResult > 0 {
			isColumnSorted = false
		}
	}

	// Если данные отсортированы хотя бы по одному критерию, возвращаем true
	return isMonthSorted || isColumnSorted
}

// Индекс месяца
// Индекс месяца
func monthIndex(line string) int {
	months := map[string]int{
		"Jan": 1, "January": 1,
		"Feb": 2, "February": 2,
		"Mar": 3, "March": 3,
		"Apr": 4, "April": 4,
		"May": 5,
		"Jun": 6, "June": 6,
		"Jul": 7, "July": 7,
		"Aug": 8, "August": 8,
		"Sep": 9, "September": 9,
		"Oct": 10, "October": 10,
		"Nov": 11, "November": 11,
		"Dec": 12, "December": 12,
	}

	words := strings.Fields(line)
	if len(words) > 0 {
		if month, exists := months[words[0]]; exists {
			return month
		}
	}
	return 0
}
