// Паттерн "Фабричный метод" - пораждающий паттерн. Предоставляет интерфейс для создания обьекта в суперклассе,
// нопозволяет подклассам изменять тип создаваемых обьектов,
// тоесть делегирует ответственность по созданию экземпляра подклассам.
//
// Плюсы:
// - избегание привязка к конкретным классам - клиентский код не зависит от конкретных классов
// - расширяемость - добавление новых типов фабрики не требует изменений кода
//
// Минусы:
// - усложнение кода
package pattern

import "fmt"

// Интерфейс продукта
type Logger interface {
	Log(message string)
}

// Конкретные типы продуктов и их методы
type FileLogger struct{}

// Метод для записи в файл
func (f *FileLogger) Log(message string) {
	fmt.Printf("Write in file: %s\n", message)
}

type ConsoleLogger struct{}

// Метод для вывода лога в консоль
func (c *ConsoleLogger) Log(message string) {
	fmt.Printf("Console log: %s\n", message)
}

// Интерфейс фабрики
type LoggerFactory interface {
	CreateLogger() Logger
}

// Конкретная фабрика для FileLogger
type FileLoggerFactory struct{}

func (f *FileLoggerFactory) CreateLogger() Logger {
	return &FileLogger{}
}

// Конкретная фабрика для ConsoleLogger
type ConsoleLoggerFactory struct{}

func (c *ConsoleLoggerFactory) CreateLogger() Logger {
	return &ConsoleLogger{}
}

func main() {
	// Создаем фабрики
	fileLoggerFactory := &FileLoggerFactory{}
	consoleLoggerFactory := &ConsoleLoggerFactory{}

	// Используем фабрики для создания логгеров
	fileLogger := fileLoggerFactory.CreateLogger()
	consoleLogger := consoleLoggerFactory.CreateLogger()

	// Записываем логи
	fileLogger.Log("Message")
	consoleLogger.Log("Message")
}
