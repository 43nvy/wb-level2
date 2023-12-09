// Паттерн "Команда" - поведенческий паттерн, позволяет инкапсулировать запрос на выполнение определенного действия в виде отдельного объекта.
// Этот объект запроса на действие и называется командой.
// При этом объекты, инициирующие запросы на выполнение действия, отделяются от объектов, которые выполняют это действие.
// Также, может поддерживать отмены операций и ставить операции в очередь.
// Плюсы:
// - изоляция отправителя и получателя - путем изоляции кода отправителя от кода получателя
// - поддержка отмены и повтора операций - позволяет легко добавлять дополнительный функционал отмены и повтора операций
// - гибкость и расширяемость - команды могут быть легко кобминированы и расширены
//
// Минусы:
// - большое количество интерфейсов и структур, что приводит к усложнению кода
// - код может быть избыточным

package pattern

import "fmt"

// Объект получателя, который выполняет операции
type TV struct {
	isOn bool
}

// Методы получателя - включение и выключение
func (tv *TV) turnOn() {
	if !tv.isOn {
		fmt.Println("TV on")
		tv.isOn = true
	} else {
		fmt.Println("TV already on")
	}
}

func (tv *TV) turnOff() {
	if tv.isOn {
		fmt.Println("TV off")
		tv.isOn = false
	} else {
		fmt.Println("TV already off")
	}
}

// Интерфейс команды
type Command interface {
	execute()
}

// Реализации интерфейсов конкретной команды
type TurnOnTVCommand struct {
	tv *TV
}

func (c *TurnOnTVCommand) execute() {
	c.tv.turnOn()
}

type TurnOffTVCommand struct {
	tv *TV
}

func (c *TurnOffTVCommand) execute() {
	c.tv.turnOff()
}

// Обьект-инициатор команд
type RemoteControl struct {
	command Command
}

func (r *RemoteControl) pressButton() {
	r.command.execute()
}

func main() {
	// Создаем обьект получателя команды
	tv := &TV{}

	// Создаем команды
	turnOnCommand := &TurnOnTVCommand{tv: tv}
	turnOffCommand := &TurnOffTVCommand{tv: tv}

	// Создаем инициатора с командой
	remoteControl := &RemoteControl{}

	// Нажимаем кнопку для включения
	remoteControl.command = turnOnCommand
	remoteControl.pressButton()

	// Нажимаем кнопку для выключения
	remoteControl.command = turnOffCommand
	remoteControl.pressButton()

}
