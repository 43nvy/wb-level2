// Паттерн "Цепочка вызовов" - поведенческий паттерн, используется для передачи запроса по цепочке обработчиков.
// Каждый обработчик решает, должен ли он обработать запрос или передать его дальше по цепочке.
//
// Плюсы:
// - гибкость и расширяемость - паттерн позволяет легко настраивать и расширять цепочку оброботчиков
// - отсутсвие зависимости между отправителем и получателем - отправитель не знает конечного получателя
//
// Минусы:
// - не гарантируется обработка запроса
// - сложная отладка - цепочки могут быть очень длинными или некорректно настроены, что усложняет отладку кода
package pattern

import (
	"fmt"
	"net/http"
)

// Обработчик запроса - ручка
type Handler interface {
	HandleRequest(request *http.Request)
}

// Конкретный обработчик - проверка аутентификации
type AuthenticationHandler struct {
	next Handler
}

func (h *AuthenticationHandler) HandleRequest(request *http.Request) {
	fmt.Println("Проверка аутентификации")
	// Абстрактная реализация
	if h.next != nil {
		h.next.HandleRequest(request)
	}
}

// Конкретный обработчик - проверка прав доступа
type AuthorizationHandler struct {
	next Handler
}

func (h *AuthorizationHandler) HandleRequest(request *http.Request) {
	fmt.Println("Проверка прав доступа")
	// Абстрактная реализация
	if h.next != nil {
		h.next.HandleRequest(request)
	}
}

// Конкретный обработчик - логирование
type LoggingHandler struct {
	next Handler
}

func (h *LoggingHandler) HandleRequest(request *http.Request) {
	// Реализация логирования
	fmt.Println("Логирование запроса")
}

func main() {
	// Создаем цепочку обработчиков
	loggingHandler := &LoggingHandler{}
	authorizationHandler := &AuthorizationHandler{next: loggingHandler}
	authenticationHandler := &AuthenticationHandler{next: authorizationHandler}

	// Создаем простой обьект запроса
	// Представим, что здесь происходит обработка HTTP запроса
	request := &http.Request{}
	authenticationHandler.HandleRequest(request)
}
