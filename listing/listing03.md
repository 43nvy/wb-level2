Что выведет программа? Объяснить вывод программы. Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.

```go
package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)
}
```

Ответ:

```
Вывод:
	<nil>
	False

В функции Foo() мы присвоили значение err = nil
Оно выводится корректно.
Но, когда мы пытаемся сравнить тип error с nil, то мы сравниваем *os.PathError с nil.
Верно, что нулевое значение error это nil, но не интерфейса *os.PathError.
И вот почему:
Интерфейс представляет собой набор методов. Переменная, которая является экземпляром интерфейса может хранить любое значение, реалезующее этот интерфейс.
Так как интерфейсы являются неявными, то тип данных авотматически соответсвует интерфейсу, если он реализует все его методы, как в нашем случае.
---
Также, в Go существует пустой интерфейс. Такие интерфейсы могут хранить значения любого типа.
До появления дженериков - разработчики пользовались пустыми интерфейсами.

```