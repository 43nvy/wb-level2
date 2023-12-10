package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// Получаем файл, флаги, открываем файл и действуем по флагам

type SortCommand struct {
	inputFile   string
	outputFile  string
	keyColumn   int
	numericSort bool
	reverseSort bool
	uniqueSort  bool
}

type Command interface {
	execute([]string) ([]string, error)
}

func (sc *SortCommand) execute(lines []string) ([]string, error) {
	// Созадем копию слайса, с данными из прочитанного файла, с помощью append
	slice := append([]string(nil), lines...)
	// Инициализируем результирующий слайс, который вернем
	resultSlice := lines[:sc.keyColumn]
	// Смотрим заданные параметры и сортируем
	if sc.keyColumn > 0 {
		slice = slice[sc.keyColumn:]
	}

	if sc.uniqueSort {
		slice = uniqueSlice(slice)
	}

	if sc.numericSort {
		numericSort(slice)
	} else {
		sort.Slice(slice, func(i, j int) bool {
			return slice[i] < slice[j]
		})
	}

	if sc.reverseSort {
		reverseSlice(slice)
	}
	// Добавляем отсортированный слайс к результирующему
	resultSlice = append(resultSlice, slice...)
	return resultSlice, nil
}

func reverseSlice(slice []string) {
	// Проходимся по всему слайсу и переприсваиваем значения
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func numericSort(slice []string) {
	// Сортируем слайс с помощью стабильной сортировки из пакета sort
	sort.SliceStable(slice, func(i, j int) bool {
		// Так, как мы проверили обработали большинство ошибок, то
		// при возникновении ошибки здесь - будет уместно завершить работу программы
		numberI, err := searchNumberInString(slice[i])
		if err != nil {
			panic(err)
		}
		numberJ, err := searchNumberInString(slice[j])
		if err != nil {
			panic(err)
		}
		// Сортируем по возрастанию
		return numberI < numberJ
	})

}

func searchNumberInString(str string) (int, error) {
	// Разбиваем строку на слова(токены)
	tokens := strings.Fields(str)
	// Рассматриваем краевой случай, если в строке одно слово
	if len(tokens) < 2 {
		return 0, nil
	}
	// Преобразовываем токен с числом в int
	number, err := strconv.Atoi(tokens[len(tokens)-1])
	if err != nil {
		return 0, err
	}

	return number, nil
}

func uniqueSlice(slice []string) []string {
	uniqueMap := make(map[string]bool)
	// Проходимся по слайсу и заносим строки в мапу
	// Если строка повторяется - то ключ просто перезапишется
	// Таким образом мы получим все уникальные строки
	for _, item := range slice {
		uniqueMap[item] = true
	}
	// Создаем слайс, длина которого = 0, а обьем равен количеству элментов мапы,
	// для предотвращения ненужных аллокаций
	resultSlice := make([]string, 0, len(uniqueMap))
	// Проходимся по мапе и аппендим ключи в результирующий слайс
	for key := range uniqueMap {
		resultSlice = append(resultSlice, key)
	}

	return resultSlice
}

func main() {
	// Инициализируем флаги
	inputFileFlag := flag.String("i", "", "Input file path")
	outputFileFlag := flag.String("o", "", "Output file path")
	keyColumnFlag := flag.Int("k", 0, "Key column for sorting")
	numericFlag := flag.Bool("n", false, "Sort by numeric value")
	reverseFlag := flag.Bool("r", false, "Reverse the order")
	uniqueFlag := flag.Bool("u", false, "Remove duplicate lines")
	// Собираем флаги
	flag.Parse()
	// Читаем файл и записываем строки в слайс
	lines, err := readFile(*inputFileFlag)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
	// Созадем обьект команды сортировки
	sortCommand := &SortCommand{
		inputFile:   *inputFileFlag,
		outputFile:  *outputFileFlag,
		keyColumn:   *keyColumnFlag,
		numericSort: *numericFlag,
		reverseSort: *reverseFlag,
		uniqueSort:  *uniqueFlag,
	}
	// Вызываем функцию сортировки и передаем туда слайс
	sortedLines, err := sortCommand.execute(lines)
	if err != nil {
		fmt.Printf("Error sorting lines: %v\n", err)
		return
	}
	// Записываем отсортированный слайс в файл
	err = writeToFile(*outputFileFlag, sortedLines)
	if err != nil {
		fmt.Printf("Error writing to output file: %v\n", err)
		return
	}

	fmt.Printf("Sorted file: %s\n", *outputFileFlag)
}

func readFile(filename string) ([]string, error) {
	// Открываем файл и проверяем на ошибку
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	// Сразу откладываем закрытие файла
	defer file.Close()
	// Создаем срез строк из файла
	// и сканнер
	var lines []string
	scanner := bufio.NewScanner(file)
	// Запускаем цикл, который будет работать, пока есть непрочитанные строки
	for scanner.Scan() {
		// Записываем эти строки в слайс
		lines = append(lines, scanner.Text())
	}
	// Проверяем на ошибки при чтении строк
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func writeToFile(filename string, lines []string) error {
	// Создаем файл и открываем его, а если файл с таким названием уже есть - перезаписываем
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	// Сразу откладываем закрытие файла
	defer file.Close()
	// Создаем writer, проходимся по слайсу строк и записываем каждую строку в буффер
	// Буфер нам нужен для того, чтобы за 1 операцию записи, записать все строки
	// Использование буфера необходима, при работе с большими обьемами данных
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	// Теперь записываем все данные из буфера в файл
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}
