Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
Вывод: error
Все дело в том, что тип переменной err = *customError. Так как мы возвращаем этот тип в функции test, со значением nil.
Но, переменная err изначально обьявляется как тип error, который является интерфейсом и реализует аналогичный типу customError метод Error.
А значение err это значение, которые хранится в интерфейсе error.
Тут явно видно работу "Duck Typing", при взаимодействии с интерфейсом, грубо говоря - "Если что то выглядит как утка и ведет себя как утка, то, вероятно, это утка!"
Так как, наша структура customError выглядит как error и ведет себя как error, то, вероятно это error.
```
