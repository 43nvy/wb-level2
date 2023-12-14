package main

import (
	"fmt"
	"sync"
	"time"
)

/*
=== Or channel ===

Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.

Определение функции:
var or func(channels ...<- chan interface{}) <- chan interface{}

Пример использования функции:
sig := func(after time.Duration) <- chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
}()
return c
}

start := time.Now()
<-or (
	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),
)

fmt.Printf(“fone after %v”, time.Since(start))
*/

func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		// Если не переданы каналы - вернуть nil
		return nil
	case 1:
		// Если передан только один канал - вернуть его
		return channels[0]
	}

	var wg sync.WaitGroup
	orDone := make(chan interface{}, 1) // Буферизованный канал

	// Функция, для получения сигнала о закрытии канала
	handleChannel := func(ch <-chan interface{}) {
		defer wg.Done()
		select {
		case <-ch:
			// Если один из каналов закрывается, отправить сигнал в orDone (без блокировки)
			select {
			case orDone <- struct{}{}: // Тут используется пустая структура, чтобы подать сигнал без данных и не производить аллокаций
			default:
			}
		case <-orDone:
		}
	}

	// Добавляем в группу ожидания количество каналов
	wg.Add(len(channels))

	// Запускаем для каждого канала функцию
	for _, ch := range channels {
		go handleChannel(ch)
	}

	// Запускаем дополнительную горутину, которая ждет отчёт от каналов и закрывает канал orDone
	go func() {
		wg.Wait()
		close(orDone)
	}()

	return orDone
}

func main() {
	// Согласно тексту задания описываем функцию
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	// Поправили опечатку :)
	fmt.Printf("Done after %v\n", time.Since(start))
}
