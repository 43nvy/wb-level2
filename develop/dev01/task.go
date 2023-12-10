package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/beevik/ntp"
)

/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/
// Структура для хранения локального времени и времени NTP
type currentTime struct {
	timePackage string
	ntpPackage  string
}

func printCurrentTime() (currentTime, error) {
	// Инициализируем структуру и, сразу, указываем локальное время
	nowTime := currentTime{timePackage: time.Now().UTC().Format(time.RFC3339)}
	// Делаем запрос точного времени и, при возникновении ошибки, возвращаем ее
	ntpTime, err := ntp.Time("pool.ntp.org")
	if err != nil {
		return nowTime, fmt.Errorf("failed to get NTP time: %v", err)
	}
	// Если все прошло хорошо - записываем точное время в структуру и возвращаем ее.
	nowTime.ntpPackage = ntpTime.UTC().Format(time.RFC3339)

	return nowTime, nil
}

func main() {
	// Вызываем функцию
	nowInfo, err := printCurrentTime()
	// Проверяем на ошибку, если есть ошибка - заканчиваем работу со статусом 1
	if err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	// Выводим время
	fmt.Printf("Local Time: %s\nNTP Time: %s\n", nowInfo.timePackage, nowInfo.ntpPackage)
}
