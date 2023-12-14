package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	// Собираем флаги
	urlFlag := flag.String("url", "", "URL сайта для скачивания")
	outputFlag := flag.String("output", "index.html", "Имя файла для сохранения контента")
	flag.Parse()

	// Проверяем на наличие URL
	if *urlFlag == "" {
		fmt.Println("Пожалуйста, укажите URL для скачивания.")
		os.Exit(1)
		return
	}

	// Если нет схемы у URL - добавляем
	if !strings.HasPrefix(*urlFlag, "http://") && !strings.HasPrefix(*urlFlag, "https://") {
		*urlFlag = "http://" + *urlFlag
	}

	// Разбираем URL
	parsedURL, err := url.Parse(*urlFlag)
	if err != nil {
		fmt.Println("Ошибка при разборе URL:", err)
		os.Exit(1)
		return
	}

	// Делаем запрос на сервер и получаем тело ответа
	content, err := downloadContent(parsedURL)
	if err != nil {
		fmt.Println("Ошибка при скачивании контента:", err)
		os.Exit(1)
		return
	}

	// Сохранение контента в файл
	err = saveToFile(*outputFlag, content)
	if err != nil {
		fmt.Println("Ошибка при сохранении в файл:", err)
		os.Exit(1)
		return
	}

	fmt.Println("Скачивание завершено. Контент сохранен в ", *outputFlag)
}

// downloadContent выполняет запрос на сервер и возвращает тело ответа в виде строки
func downloadContent(parsedURL *url.URL) (string, error) {
	response, err := http.Get(parsedURL.String())
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Читаем тело ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// saveToFile сохраняет контент страницы в файл
func saveToFile(filename string, content string) error {
	// Создаем директорию, если она не существует
	dir := path.Dir(filename)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// Создаем файл
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	// Пишем в файл
	_, err = io.WriteString(file, content)
	if err != nil {
		return err
	}

	return nil
}
