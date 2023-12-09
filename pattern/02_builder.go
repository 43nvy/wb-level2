// Паттер "Строитель" - это пораждающий паттерн, используется для пошагового создания сложных обьектов,
// используя один и тот же процесс
//
// Плюсы:
// - гибкость: позволяет создавать различные конфигурации обьекта
// - разделение ответственности: раздельные процессы построения обьекта от его представления
// - читаемость: помогает избежать "антипаттерна"
//
// Минусы:
// - избыточность кода: для создания обьекта могут потеброваться дополнительные интерфейсы
package pattern

import "fmt"

// Конечный обьект животного
type Animal struct {
	Kind string
	Roar string
	Legs int
}

// Интерфейс, для пошагового строения животного
type AnimalBuilder interface {
	SetKind(kind string)
	SetRoar(roar string)
	SetLegs(legs int)
	Build() *Animal
}

// Конкретная реализация интерфейса билдера
type ConcreteAnimalBuilder struct {
	animal *Animal
}

// Функция создания нового билдера
func NewAnimalBuilder() *ConcreteAnimalBuilder {
	return &ConcreteAnimalBuilder{animal: &Animal{}}
}

// Методы-сеттеры
func (ca *ConcreteAnimalBuilder) SetRoar(roar string) {
	ca.animal.Roar = roar
}

func (cb *ConcreteAnimalBuilder) SetKind(kind string) {
	cb.animal.Kind = kind
}

func (ca *ConcreteAnimalBuilder) SetLegs(legs int) {
	ca.animal.Legs = legs
}

// Метод построения животного, который возвращает конечный обьект
func (ca *ConcreteAnimalBuilder) Build() *Animal {
	return ca.animal
}

// Структура, управляющая построением животного
type Director struct {
	builder AnimalBuilder
}

// Функция для создаия нового директора
func NewDirector(builder AnimalBuilder) *Director {
	return &Director{builder: builder}
}

func (d *Director) ConstructAnimal() *Animal {
	d.builder.SetKind("Wolf")
	d.builder.SetRoar("Rrrrr")
	d.builder.SetLegs(4)
	return d.builder.Build()
}

func main() {
	// Использование паттерна Строитель
	builder := NewAnimalBuilder()
	director := NewDirector(builder)

	animal := director.ConstructAnimal()

	fmt.Println("Animal:")
	fmt.Printf("Kind: %s\n", animal.Kind)
	fmt.Printf("Roar: %d ГБ\n", animal.Roar)
	fmt.Printf("Legs: %d ГБ\n", animal.Legs)
}
