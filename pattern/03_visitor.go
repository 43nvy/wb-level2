// Паттерн "Посетитель" - это поведенческий паттерн, который позволяет определить новую операцию без изменения
// структуры обьектов, над которыми эта операция выполняется.
//
// Плюсы:
// - разделение ответственности - паттерн позволяеть разделять операции над обьектами от их структуры
// - не нужно переписывать код - добавление нового функционала, не затрагивая существующие обьекты
//
// Минусы:
// - сложность добавления новых структр обьектов - приходится изменять все классы посетителей
// - нарушение инкапсуляции - посетитель может нарушить инкапсуляцию, так как получает доступ к внутренней структуре
package pattern

import "fmt"

// Интерфейс элемента
type Element interface {
	Accept(visitor Visitor)
}

// Конкретные элементы
type ConcreteTextBox struct{}

func (t *ConcreteTextBox) Accept(visitor Visitor) {
	visitor.VisitTextBox(t)
}

type ConcreteButton struct{}

func (b *ConcreteButton) Accept(visitor Visitor) {
	visitor.VisitButton(b)
}

// Интерфейс посетителя, описывающий методы элементов
type Visitor interface {
	VisitTextBox(textBox *ConcreteTextBox)
	VisitButton(button *ConcreteButton)
}

// Конкретные посетители и их разные методы
type WindowsVisitor struct{}

func (w *WindowsVisitor) VisitTextBox(textBox *ConcreteTextBox) {
	fmt.Println("WindowsVisitor processes TextBox")
}

func (w *WindowsVisitor) VisitButton(button *ConcreteButton) {
	fmt.Println("WindowsVisitor processes Button")
}

type LinuxVisitor struct{}

func (l *LinuxVisitor) VisitTextBox(textBox *ConcreteTextBox) {
	fmt.Println("LinuxVisitor processes TextBox")
}

func (l *LinuxVisitor) VisitButton(button *ConcreteButton) {
	fmt.Println("LinuxVisitor processes Button")
}

func main() {
	// Создаем срез элементов, для удобного использования
	elements := []Element{&ConcreteTextBox{}, &ConcreteButton{}}

	// Примеры разных посетителей
	windowsVisitor := &WindowsVisitor{}
	linuxVisitor := &LinuxVisitor{}

	// Обработка элементов для различных ОС
	for _, element := range elements {
		element.Accept(windowsVisitor)
		element.Accept(linuxVisitor)
	}
}
