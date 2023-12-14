package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	// Флаги
	fieldsFlag := flag.String("f", "", "выбрать поля (колонки)")
	delimiterFlag := flag.String("d", " ", "использовать другой разделитель")
	separatedFlag := flag.Bool("s", false, "только строки с разделителем")
	flag.Parse()

	// Проверка наличия хотя бы одного флага
	if flag.NFlag() == 0 {
		fmt.Println("Используйте хотя бы один флаг: -f, -d, -s")
		os.Exit(0)
	}

	needFields, err := fieldsFlagParse(*fieldsFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Создаем матрицу, для удобного хранений строк и колонок
	matrix := make([][]string, 0)

	// Обработка STDIN построчно
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		// Разделение строки на слова
		words := strings.Fields(line)
		if *separatedFlag {
			// Проверяем, содержит ли строка указанный разделитель
			if strings.Contains(line, *delimiterFlag) {
				matrix = append(matrix, words)
			}

		} else {
			// Иначе просто добавляем все слова в матрицу
			matrix = append(matrix, words)
		}
	}

	// Проверка ошибок сканера
	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении STDIN:", err)
		os.Exit(1)
	}

	resultMatrix := make([][]string, 0)

	if len(needFields) == 0 {
		resultMatrix = append(resultMatrix, matrix...)

	} else {
		for _, row := range matrix {
			// Создаем новую строку для resultMatrix
			newRow := make([]string, 0)

			// Проходим по каждому индексу в needFields
			for _, index := range needFields {
				// Проверяем, что индекс находится в пределах длины строки
				if index >= 0 && index < len(row) {
					// Добавляем значение из оригинальной строки в новую строку
					newRow = append(newRow, row[index])
				}
			}
			// Рассматриваем краевой случай
			if len(newRow) == 0 {
				continue

			} else {
				// Добавляем новую строку в resultMatrix
				resultMatrix = append(resultMatrix, newRow)
			}

		}
	}

	if len(resultMatrix) == 0 {
		fmt.Println("Нет результатов, проверьте опции")
		os.Exit(0)
	}

	for _, value := range resultMatrix {
		fmt.Println(value)
	}
}

func fieldsFlagParse(flagData string) ([]int, error) {
	// Обрабатываем кравеой случай
	if flagData == "" {
		return nil, nil
	}
	// Делим строку из флага по разделителю
	stringSlice := strings.Split(flagData, ",")
	// Инициализируем слайс из чисел
	numbers := make([]int, 0, len(stringSlice))
	for _, value := range stringSlice {
		value = strings.TrimSpace(value)
		// Пробуем преобразовать значение в число
		num, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("ошибка при преобразовании числа: '%s'", value)
		}
		// Обрабатываем краевой случай
		if num != 0 {
			num = int(math.Abs(float64(num)))
			// Вычитаем единичку, чтобы корректно отображать строки
			num--
		}

		numbers = append(numbers, num)
	}

	return numbers, nil
}
