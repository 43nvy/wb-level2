// Паттерн "Фасад" предоставляеет унифицированный интерфейс для набора интерфейсов в подсистеме.
//
// Плюсы :
// - простое взаимодейсвтие для клиента, так как он не видит детали реализации
// - реализация принципа единственной ответственности
//
// Из минусов, можно выделить :
// - малая гибкость, так как клиент может получить доступ только к тем функциям, которые описаные в фасаде
// - большие слоя абстракций в некоторых случаях

package pattern

import "fmt"

// Некоторая сложная подсистема базы данных
type DBSubsystem struct {
}

// Методы подсистемы базы данных
// В данном примере их реализация не важна, но, будем думать что вместо fmt.Println, что то сложное,
// например хитрые транзакции
func (d *DBSubsystem) Connect() {
	fmt.Println("DB connection")
}

func (d *DBSubsystem) Query(query string) {
	fmt.Printf("Insert query: %s\n", query)
}

// Фасад для упрощения работы с БД
type DatabaseFacade struct {
	db *DBSubsystem
}

// Функция создания нового фасада
func NewDatabaseFacade() *DatabaseFacade {
	return &DatabaseFacade{
		db: &DBSubsystem{},
	}
}

// Описанные методы фасада
func (f *DatabaseFacade) Initialize() {
	f.db.Connect()
}

func (f *DatabaseFacade) QueryData(query string) {
	f.db.Query(query)
}

func main() {
	facade := NewDatabaseFacade()
	facade.Initialize()
	facade.QueryData("SELECT * FROM users")
}
