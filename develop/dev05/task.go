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
	"sort"
	"strconv"
	"strings"
)

// Options структура, для передачи параметров
type Options struct {
	after      int
	before     int
	context    int
	count      bool
	ignoreCase bool
	invert     bool
	fixed      bool
	lineNum    bool
}

// sortLinesByNumber функция сортировки строк. Работает только, если установлен флаг -n
func sortLinesByNumber(lines []string) {
	sort.Slice(lines, func(i, j int) bool {
		numI, _ := strconv.Atoi(strings.Split(lines[i], ":")[0])
		numJ, _ := strconv.Atoi(strings.Split(lines[j], ":")[0])
		return numI < numJ
	})
}

// processBeforeContext функция, для вывода строк перед искомой - флаг -B
func processBeforeContext(lines []string, start, end int, resultMap map[string]struct{}, withNum bool) {
	for j := start; j < end; j++ {
		if j >= 0 {
			resultMap[formatLineWithNum(lines[j], j+1, withNum)] = struct{}{}
		}
	}
}

// processBeforeContext функция, для вывода строк после искомой - флаг -A
func processAfterContext(lines []string, start, end int, resultMap map[string]struct{}, withNum bool) {
	for j := start; j < end; j++ {
		resultMap[formatLineWithNum(lines[j], j+1, withNum)] = struct{}{}
	}
}

// formatLineWithNum функция, которая форматирует строки для вывода - флаг -n
func formatLineWithNum(line string, lineNum int, withNum bool) string {
	if withNum {
		return fmt.Sprintf("%d:%s", lineNum, line)
	}
	return line
}

// processContext основная функция нахождения строк
func processContext(lines []string, i int, resultMap map[string]struct{}, opts Options) {
	// Выводим строки перед найденной строкой (флаг -B)
	if opts.context > 0 || opts.before > 0 {
		start := i - opts.context - opts.before
		if start < 0 {
			start = 0
		}
		processBeforeContext(lines, start, i, resultMap, opts.lineNum)
	}

	// Выводим найденную строку
	lineWithNum := formatLineWithNum(lines[i], i+1, opts.lineNum)
	resultMap[lineWithNum] = struct{}{}

	// Выводим строки после найденной строки (флаг -A)
	if opts.context > 0 || opts.after > 0 {
		end := i + opts.context + opts.after + 1
		if end > len(lines) {
			end = len(lines)
		}
		processAfterContext(lines, i+1, end, resultMap, opts.lineNum)
	}
}

// grep функция, которая выполняет нахождение по параметрам
func grep(lines []string, pattern string, opts Options) []string {
	resultMap := make(map[string]struct{})

	for i, line := range lines {
		matched := opts.match(line, pattern)
		if opts.invert {
			matched = !matched
		}

		if matched {
			processContext(lines, i, resultMap, opts)
		}
	}

	var result []string
	for k := range resultMap {
		result = append(result, k)
	}

	// Сортировка по числу, если установлен флаг -n
	if opts.lineNum {
		sortLinesByNumber(result)
	}

	return result
}

// match функция проверяет, соответсвие строки указанному шаблону, с учетом опций
func (o Options) match(line, pattern string) bool {
	// Если установлен флаг -F, используем точное совпадение строк
	if o.fixed {
		return line == pattern
	}
	// Если установлен флаг -i, игнорируем регистр при сравнении
	if o.ignoreCase {
		line = strings.ToLower(line)
		pattern = strings.ToLower(pattern)
	}
	// Или просто используем регулярки для сравнения
	match, _ := regexp.MatchString(pattern, line)
	return match
}

func main() {
	// Инициализируем флаги
	var after, before, context int
	var count, ignoreCase, invert, fixed, lineNum bool

	flag.IntVar(&after, "A", 0, "Print N lines after each match")
	flag.IntVar(&before, "B", 0, "Print N lines before each match")
	flag.IntVar(&context, "C", 0, "Print N lines of output context")
	flag.BoolVar(&count, "c", false, "Print only a count of matching lines")
	flag.BoolVar(&ignoreCase, "i", false, "Case-insensitive matching")
	flag.BoolVar(&invert, "v", false, "Invert the sense of matching")
	flag.BoolVar(&fixed, "F", false, "Fixed, exact matching")
	flag.BoolVar(&lineNum, "n", false, "Print line numbers")
	flag.Parse()

	args := flag.Args()
	// Парсим опции
	if len(args) != 2 {
		fmt.Println("Usage: grep [OPTIONS] PATTERN FILE")
		os.Exit(1)
	}

	pattern := args[0]
	filePath := args[1]
	// Открываем файл
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	// Инициализируем сканнер и читаем файл построчно
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	// Инициализируем структуру для передачи параметров в функцию grep
	options := Options{
		after:      after,
		before:     before,
		context:    context,
		count:      count,
		ignoreCase: ignoreCase,
		invert:     invert,
		fixed:      fixed,
		lineNum:    lineNum,
	}
	// Вызываем функцию grep
	matchingLines := grep(lines, pattern, options)
	// Выводим результат, в соответсвии с опциями
	if count {
		fmt.Println(len(matchingLines))
	} else {
		if len(matchingLines) == 0 {
			fmt.Println("No matches found.")
		} else {
			for _, line := range matchingLines {
				fmt.Println(line)
			}
		}
	}
}
