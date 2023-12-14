package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

// handleInterrupt обрабатывает сигнал прерывания и закрывает соединение
func handleInterrupt(conn net.Conn) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("\nСоединение закрыто.")
		conn.Close()
		os.Exit(0)
	}()
}

// connect устанавливает соединение с сервером
func connect(host string, port int, timeout time.Duration) (net.Conn, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к %s: %v", address, err)
	}
	fmt.Printf("Подключено к %s\n", address)
	return conn, nil
}

// runTelnetClient запускает telnet клиент
func runTelnetClient(conn net.Conn) {
	// Чтение из STDIN и запись в сокет
	go func() {
		_, err := io.Copy(conn, os.Stdin)
		if err != nil {
			fmt.Printf("Ошибка при копировании из STDIN в сокет: %v\n", err)
		}
	}()

	// Чтение из сокета и запись в STDOUT
	_, err := io.Copy(os.Stdout, conn)
	if err != nil {
		fmt.Printf("Ошибка при копировании из сокета в STDOUT: %v\n", err)
	}
}

func main() {
	var (
		host    string
		port    int
		timeout time.Duration
	)

	// Парсим флаги
	flag.StringVar(&host, "host", "", "Хост для подключения")
	flag.IntVar(&port, "port", 0, "Порт для подключения")
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "Таймаут подключения")
	flag.Parse()

	// Валидируем флаги
	if host == "" || port == 0 {
		fmt.Println("Использование: go-telnet --timeout=<timeout> host port")
		os.Exit(1)
	}

	// Вызываем функцию установки соединения
	conn, err := connect(host, port, timeout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	// Вызываем функцию обработки сигнала заверешения работы
	handleInterrupt(conn)

	// Вызываем функцию запуска телнет клиента
	runTelnetClient(conn)
}
