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

func sortString(s string) string {
	// Преобразовываем строку в массив символов
	sChars := strings.Split(s, "")
	// Сортируем
	sort.Strings(sChars)
	// И соединям символы обратно в строку
	return strings.Join(sChars, "")
}

func findAnagrams(words []string) map[string][]string {
	anagramSets := make(map[string][]string)

	for _, word := range words {
		// Приводим слова к нижнему регистру и сортируем символы
		sortedWord := sortString(strings.ToLower(word))
		// Добавляем слова в множество анаграмм
		anagramSets[sortedWord] = append(anagramSets[sortedWord], word)
	}

	// Убираем множества из одного элемента
	for key, value := range anagramSets {
		if len(value) < 2 {
			delete(anagramSets, key)
		} else {
			// Сортируем множество по возрастанию
			sort.Strings(anagramSets[key])
		}
	}

	return anagramSets
}

func main() {
	words := []string{"пятак", "пятка", "пол", "тяпка", "листок", "слиток", "столик", "кот", "ток", "кто", "сам", "стул"}
	fmt.Printf("Исходный срез: %s\n", words)
	anagramSets := findAnagrams(words)

	for _, value := range anagramSets {
		fmt.Printf("Множество анаграмм для %s: %v\n", value[0], value)
	}
}
