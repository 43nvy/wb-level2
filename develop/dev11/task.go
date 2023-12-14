package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
Методы API: POST /create_event POST /update_event POST /delete_event GET /events_for_day GET /events_for_week GET /events_for_month
Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

// Event структура представления событий в календаре
type Event struct {
	ID       int       `json:"id"`
	UserID   int       `json:"user_id"`
	Title    string    `json:"title"`
	Date     time.Time `json:"date"`
	Duration int       `json:"duration"`
}

// ServerConfig структура представления конфигурации сервера
type ServerConfig struct {
	Port int `json:"port"`
}

var (
	configMutex  sync.Mutex
	serverConfig ServerConfig
)

var (
	eventsMutex sync.Mutex
	events      []Event
)

// readConfigFromFile вспомогательная функция, для чтения конфига
func readConfigFromFile() error {
	configMutex.Lock()
	defer configMutex.Unlock()

	file, err := os.ReadFile("config.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &serverConfig)
	if err != nil {
		return err
	}

	return nil
}

// serializeEvent вспомогательная функция, для сериализации события
func serializeEvent(event Event) ([]byte, error) {
	return json.Marshal(event)
}

// serializeEvents вспомогательная функция, для сериализации списка событий
func serializeEvents(events []Event) ([]byte, error) {
	return json.Marshal(events)
}

// readEventsFromFile вспомогательная функция, для чтения событий из файла
func readEventsFromFile() error {
	eventsMutex.Lock()
	defer eventsMutex.Unlock()

	file, err := os.ReadFile("events.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &events)
	if err != nil {
		return err
	}

	return nil
}

// writeEventsToFile вспомогательная функция, для записи событий в файл
func writeEventsToFile() error {
	eventsMutex.Lock()
	defer eventsMutex.Unlock()

	file, err := json.MarshalIndent(events, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile("events.json", file, 0644)
	if err != nil {
		return err
	}

	return nil
}

// filterEventsByDate вспомогательная функция, для фильтрации событий по дате
func filterEventsByDate(start, end time.Time) []Event {
	eventsMutex.Lock()
	defer eventsMutex.Unlock()

	var filteredEvents []Event
	for _, event := range events {
		if event.Date.After(start) && event.Date.Before(end) {
			filteredEvents = append(filteredEvents, event)
		}
	}

	return filteredEvents
}

// parseEventParams вспомогательная функция, для парсинга и валидации параметров создания или обновления события
func parseEventParams(r *http.Request) (Event, error) {
	userIDStr := r.FormValue("user_id")
	title := r.FormValue("title")
	dateStr := r.FormValue("date")
	durationStr := r.FormValue("duration")

	// Валидация обязательных параметров
	if userIDStr == "" || title == "" || dateStr == "" || durationStr == "" {
		return Event{}, fmt.Errorf("missing required parameters")
	}

	// Парсинг user_id
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return Event{}, fmt.Errorf("invalid user_id format")
	}

	// Парсинг date в формате "2006-01-02"
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return Event{}, fmt.Errorf("invalid date format")
	}

	// Парсинг duration
	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		return Event{}, fmt.Errorf("invalid duration format")
	}

	// Валидация, длительноси события
	if duration <= 0 {
		return Event{}, fmt.Errorf("duration should be a positive integer")
	}

	return Event{
		UserID:   userID,
		Title:    title,
		Date:     date,
		Duration: duration,
	}, nil
}

// loggingMiddleware вспомогательная функция-мидляр, для логгирования
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		// Передача запроса следующему обработчику в цепочке
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Загрузка конфигурации из файла
	err := readConfigFromFile()
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	// Загрузка событий из файла
	err = readEventsFromFile()
	if err != nil {
		log.Fatal("Error loading events: ", err)
	}

	// Использование middleware
	http.Handle("/", loggingMiddleware(http.NotFoundHandler()))

	// Обработчики для каждого метода
	http.HandleFunc("/create_event", createEventHandler)
	http.HandleFunc("/update_event", updateEventHandler)
	http.HandleFunc("/delete_event", deleteEventHandler)
	http.HandleFunc("/events_for_day", eventsForDayHandler)
	http.HandleFunc("/events_for_week", eventsForWeekHandler)
	http.HandleFunc("/events_for_month", eventsForMonthHandler)

	log.Printf("Server is running on port %d", serverConfig.Port)

	// Запуск HTTP-сервера
	err = http.ListenAndServe(fmt.Sprintf(":%d", serverConfig.Port), nil)
	if err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}

// createEventHandler ручка метода /create_event
func createEventHandler(w http.ResponseWriter, r *http.Request) {
	event, err := parseEventParams(r)
	if err != nil {
		http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	eventsMutex.Lock()
	defer eventsMutex.Unlock()

	// Генерация уникального ID для события
	event.ID = len(events) + 1

	// Добавление события к списку
	events = append(events, event)

	// Сохранение изменений в файл
	err = writeEventsToFile()
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Преобразование события в JSON и возврат успешного результата
	serializedEvent, err := serializeEvent(event)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializedEvent)
}

// updateEventHandler ручка метода /update_event
func updateEventHandler(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.FormValue("event_id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "Bad Request: Invalid event ID", http.StatusBadRequest)
		return
	}

	eventsMutex.Lock()
	defer eventsMutex.Unlock()

	// Поиск события по ID
	var found bool
	var updatedEvent Event
	for i, existingEvent := range events {
		if existingEvent.ID == eventID {
			// Обновление данных события
			events[i].Title = r.FormValue("title")
			dateStr := r.FormValue("date")
			date, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				http.Error(w, "Bad Request: Invalid date format", http.StatusBadRequest)
				return
			}
			events[i].Date = date
			events[i].Duration, _ = strconv.Atoi(r.FormValue("duration"))

			// Сохранение обновленного события
			updatedEvent = events[i]

			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Not Found: Event not found", http.StatusNotFound)
		return
	}

	// Сохранение изменений в файл
	err = writeEventsToFile()
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	serializedEvent, err := serializeEvent(updatedEvent)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializedEvent)
}

// deleteEventHandler ручка метода /delete_event
func deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.FormValue("event_id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "Bad Request: Invalid event ID", http.StatusBadRequest)
		return
	}

	eventsMutex.Lock()
	defer eventsMutex.Unlock()

	// Поиск и удаление события по ID
	var found bool
	var deletedEvent Event
	for i, existingEvent := range events {
		if existingEvent.ID == eventID {
			// Сохранение удаляемого события
			deletedEvent = existingEvent

			// Удаление события из списка
			events = append(events[:i], events[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Not Found: Event not found", http.StatusNotFound)
		return
	}

	// Сохранение изменений в файл
	err = writeEventsToFile()
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	serializedEvent, err := serializeEvent(deletedEvent)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializedEvent)
}

// eventsForDayHandler ручка метода /events_for_day
func eventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	dateStr := r.FormValue("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Bad Request: Invalid date format", http.StatusBadRequest)
		return
	}

	start := date
	end := date.Add(24 * time.Hour)

	filteredEvents := filterEventsByDate(start, end)

	serializedEvents, err := serializeEvents(filteredEvents)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializedEvents)
}

// eventsForWeekHandler ручка метода /events_for_week
func eventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	dateStr := r.FormValue("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Bad Request: Invalid date format", http.StatusBadRequest)
		return
	}

	// Находим ПН, для указанной даты
	weekday := int(date.Weekday())
	start := date.AddDate(0, 0, -weekday)
	end := start.Add(7 * 24 * time.Hour)

	filteredEvents := filterEventsByDate(start, end)

	serializedEvents, err := serializeEvents(filteredEvents)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializedEvents)
}

// eventsForMonthHandler ручка метода /events_for_month
func eventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	dateStr := r.FormValue("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Bad Request: Invalid date format", http.StatusBadRequest)
		return
	}

	// Находим первый день месяца для указанной даты
	start := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	filteredEvents := filterEventsByDate(start, end)

	serializedEvents, err := serializeEvents(filteredEvents)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializedEvents)
}
