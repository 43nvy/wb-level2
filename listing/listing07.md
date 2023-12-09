Что выведет программа? Объяснить вывод программы.

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)

	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}

		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v := <-a:
				c <- v
			case v := <-b:
				c <- v
			}
		}
	}()
	return c
}

func main() {

	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4 ,6, 8)
	c := merge(a, b )
	for v := range c {
		fmt.Println(v)
	}
}
```

Ответ:
```
Вывод: 1-8, в неизвестном порядке, а затем бесконечный вывод 0.
Здесь все просто, две горутины, запущенные из двух вызовов asChan, выполняют свою работу - посылают значения в канал и закрывают его.
В то время, как asChan x2 посылают значения, горутина, вызыванная из merge, передает эти значения в другой канал, который в свою очередь,
передает эти значения в main горутину, которая печатает эти значения. Но, после того как asChan x2 горутины закрыли каналы - merge горутина продолжает читать из них,
когда горутина читает закрытый канал, то получает значение "по-умолчанию" типа этого канала. Таким образом горутина merge отправляет дефолтный тип int = 0,
в main горутину, а она, в свою очередь, их печатает в бесконечном цикле.

```