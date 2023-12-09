// Паттерн "Стратегия" - поведенческий паттерн, он определяет набор алгоритмов, инкапсулирует каждый из них и делает их взаимозаменяемыми.
// Позволяет изменять алгоритмы независимо от клиентов, которые их используют.
//
// Плюсы:
// - гибкость и поддержание кодовой базы
// - изолирование ответственности - каждый алгоритм инкапсулирован в собственном классе
//
// Минусы:
// - усложнение структуры - создание стратегий может привести к большому количеству структур и их интерфейсов
package pattern

import (
	"fmt"
	"sort"
)

// Интерфейс стратегии для сортировки
type SortStrategy interface {
	Sort(data []int) []int
}

// Конкретная стратегия сортировки пузырьком
type BubbleSort struct{}

// Реализация сортировки пузырьком (моковая)
func (bs *BubbleSort) Sort(data []int) []int {
	fmt.Println("Bubble sort")
	sort.Slice(data, func(i, j int) bool {
		return data[i] < data[j]
	})
	return data
}

// Конкретная стратегия быстрой сортировки
type QuickSort struct{}

// Реализация быстрой сортировки (моковая)
func (qs *QuickSort) Sort(data []int) []int {
	fmt.Println("Quick sort")
	sort.Slice(data, func(i, j int) bool {
		return data[i] < data[j]
	})
	return data
}

// Некоторый контекст, использующий стратегию сортировки
type Context struct {
	strategy SortStrategy
}

// Метод выбора стратегии
func (c *Context) SetStrategy(strategy SortStrategy) {
	c.strategy = strategy
}

func (c *Context) ExecuteStrategy(data []int) []int {
	if c.strategy == nil {
		fmt.Println("Set sort strategy")
		return data
	}
	return c.strategy.Sort(data)
}

func main() {
	// Создаем контекст
	context := &Context{}
	// Инициализируем данные для сортировки
	data := []int{5, 2, 8, 1, 6}

	// Устанавливаем стратегию сортировки пузырьком
	context.SetStrategy(&BubbleSort{})
	result := context.ExecuteStrategy(data)
	fmt.Println("Result:", result)

	// Устанавливаем стратегию быстрой сортировки
	context.SetStrategy(&QuickSort{})
	result = context.ExecuteStrategy(data)
	fmt.Println("Result:", result)
}
