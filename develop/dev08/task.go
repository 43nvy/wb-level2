package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	// Цикл обработки команд
	for {
		// Отображение приглашения
		fmt.Print("MyShell $ ")

		// Вводи пользователя
		inputReader := bufio.NewReader(os.Stdin)
		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка чтения ввода:", err)
			continue
		}

		// Тримим символ новой строки
		input = strings.TrimSuffix(input, "\n")

		// Запускаем функцию обработки комманд
		processCommand(input)

		// Обработка команды выхода
		if input == "\\quit" {
			fmt.Println("Выход из шелла.")
			break
		}
	}
}

func processCommand(command string) {
	// Партицируем комманду
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "nc":
		if len(parts) < 3 {
			fmt.Println("Использование: netcat <hostname> <port>")
		} else {
			host := parts[1]
			port := parts[2]

			err := startNetcat(host, port)
			if err != nil {
				fmt.Println("Ошибка при использовании netcat:", err)
			}
		}

	case "kill":
		if len(parts) > 1 {
			processID := parts[1]
			err := killProcess(processID)
			if err != nil {
				fmt.Println("Ошибка при завершении процесса:", err)
			}
		} else {
			fmt.Println("Не указан идентификатор процесса для kill.")
		}

	case "ps":
		err := listProcesses()
		if err != nil {
			fmt.Println("Ошибка при выводе списка процессов:", err)
		}

	case "cd":
		if len(parts) > 1 {
			err := os.Chdir(parts[1])
			if err != nil {
				fmt.Println("Ошибка при смене директории:", err)
			}
		} else {
			fmt.Println("Не указан аргумент для cd.")
		}

	case "pwd":
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Ошибка при получении текущей директории:", err)
		} else {
			fmt.Println(currentDir)
		}

	case "echo":
		if len(parts) > 1 {
			fmt.Println(strings.Join(parts[1:], " "))
		} else {
			fmt.Println("Не указан аргумент для echo.")
		}

	default:
		if err := executeCommand(parts); err != nil {
			fmt.Println("Ошибка при выполнении команды:", err)
		}
	}
}

// startNetcat функция TCP connection и обрабатывает I/O bound
func startNetcat(host, port string) error {
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Println("Подключено к", host+":"+port)

	go func() {
		// Чтение данных и вывод на экран
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	// Чтение ввода пользователя и отправка данных
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Fprintln(conn, text)
	}

	return nil
}

// killProcess реализует kill команду
func killProcess(processID string) error {
	pid, err := parsePID(processID)
	if err != nil {
		return err
	}

	// Завершение процесса по PID
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		return err
	}

	return nil
}

// listProcesses реализует ps
func listProcesses() error {
	cmd := exec.Command("ps", "aux")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// parsePID вспомогательная функция - парсит строку в число
func parsePID(processID string) (int, error) {
	// Преобразование строки PID процесса в int
	pid, err := strconv.Atoi(processID)
	if err != nil {
		return 0, fmt.Errorf("некорректный идентификатор процесса: %v", err)
	}
	return pid, nil
}

// executeCommand вспомогательная функция, которая исполняет команду или парсит в пайплайн
func executeCommand(parts []string) error {
	// Проверка наличия символа "|"
	if containsPipe(parts) {
		return executePipeline(parts)
	}

	// Создание нового процесса с использованием fork
	cmd := exec.Command(parts[0], parts[1:]...)

	// Подключение канала вывода процесса к STDOUT текущего процесса
	cmd.Stdout = os.Stdout

	// Запуск нового процесса
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// containsPipe всопомгательная функция, которая проверяет наличие символа "|"
func containsPipe(parts []string) bool {
	for _, part := range parts {
		if part == "|" {
			return true
		}
	}
	return false
}

// executePipeline вспомогательаня функция, которая парсит пайплайн на команды
func executePipeline(parts []string) error {
	// Разделяем команды по символу "|"
	commands := splitByPipe(parts)

	// Запускаем каждую команду в конвейере
	var err error
	var input io.Reader = nil

	for _, cmdParts := range commands {
		// Создаем новый процесс с использованием fork
		cmd := exec.Command(cmdParts[0], cmdParts[1:]...)

		// Подключаем ввод предыдущего процесса к вводу текущего процесса (кроме первой команды)
		if input != nil {
			cmd.Stdin = input
		}

		// Используем bytes.Buffer для хранения вывода процесса
		var output bytes.Buffer
		cmd.Stdout = &output

		// Запускаем процесс
		err = cmd.Run()
		if err != nil {
			break
		}

		// Сохраняем вывод текущей команды для передачи следующей команде
		input = &output
	}

	return err
}

// splitByPipe вспомогательная функция, которая сплитит команды
func splitByPipe(parts []string) [][]string {
	var commands [][]string
	var currentCommand []string

	for _, part := range parts {
		if part == "|" {
			// Завершаем текущую команду и добавляем ее в конвейер
			commands = append(commands, currentCommand)
			currentCommand = nil
		} else {
			// Добавляем часть команды к текущей команде
			currentCommand = append(currentCommand, part)
		}
	}

	// Добавляем последнюю команду в конвейер
	commands = append(commands, currentCommand)

	return commands
}
