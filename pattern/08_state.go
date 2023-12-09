// Паттерн "состояние" - поведенческий паттерн, который позволяет объекту изменять свое поведение в зависимости от внутреннего состояния.
// Создает различные классы для каждого состояния объекта и делегирует запросы к объекту текущего состояния.
//
// Плюсы:
// - изолирование состояний - каждое состояние инкапсулируется в отдельном классе
// - расширение - паттерн позволяет легко добавлять новые состояния и управлять переходами между ними
//
// Минусы:
// - усложнение кода - может быть большое колличество структур и интерфейсов, а так же проявлять внимание к деталм переходов, что усложняет поддержку

package pattern

import "fmt"

// Интерфейс состояний
type State interface {
	InsertCoin(vendingMachine *VendingMachine)
	EjectCoin(vendingMachine *VendingMachine)
	SelectDrink(vendingMachine *VendingMachine)
	DispenseDrink(vendingMachine *VendingMachine)
}

// Конкретное состояние NoCoinState
type NoCoinState struct{}

// Имплементация методов
func (s *NoCoinState) InsertCoin(vendingMachine *VendingMachine) {
	fmt.Println("Coin inserted")
	vendingMachine.SetState(vendingMachine.HasCoinState)
}

func (s *NoCoinState) EjectCoin(vendingMachine *VendingMachine) {
	fmt.Println("No coin to eject")
}

func (s *NoCoinState) SelectDrink(vendingMachine *VendingMachine) {
	fmt.Println("Insert coin to select a drink")
}

func (s *NoCoinState) DispenseDrink(vendingMachine *VendingMachine) {
	fmt.Println("Insert coin to get a drink")
}

// Конкретное состояне HasCoinState
type HasCoinState struct{}

func (s *HasCoinState) InsertCoin(vendingMachine *VendingMachine) {
	fmt.Println("Coin already inserted")
}

func (s *HasCoinState) EjectCoin(vendingMachine *VendingMachine) {
	fmt.Println("Coin ejected")
	vendingMachine.SetState(vendingMachine.NoCoinState)
}

func (s *HasCoinState) SelectDrink(vendingMachine *VendingMachine) {
	fmt.Println("Drink selected")
	vendingMachine.SetState(vendingMachine.SoldState)
}

func (s *HasCoinState) DispenseDrink(vendingMachine *VendingMachine) {
	fmt.Println("Select a drink first")
}

// Конкретное состояние SoldState
type SoldState struct{}

func (s *SoldState) InsertCoin(vendingMachine *VendingMachine) {
	fmt.Println("Please wait, dispensing a drink")
}

func (s *SoldState) EjectCoin(vendingMachine *VendingMachine) {
	fmt.Println("Sorry, can't eject coin after selecting a drink")
}

func (s *SoldState) SelectDrink(vendingMachine *VendingMachine) {
	fmt.Println("Please wait, dispensing a drink")
}

func (s *SoldState) DispenseDrink(vendingMachine *VendingMachine) {
	fmt.Println("Drink dispensed")
	vendingMachine.ReleaseDrink()
	if vendingMachine.GetCount() > 0 {
		vendingMachine.SetState(vendingMachine.NoCoinState)
	} else {
		vendingMachine.SetState(vendingMachine.SoldOutState)
	}
}

// Конкретное состояние - SoldOutState
type SoldOutState struct{}

func (s *SoldOutState) InsertCoin(vendingMachine *VendingMachine) {
	fmt.Println("No drinks available")
}

func (s *SoldOutState) EjectCoin(vendingMachine *VendingMachine) {
	fmt.Println("Coins returned")
	vendingMachine.SetState(vendingMachine.NoCoinState)
}

func (s *SoldOutState) SelectDrink(vendingMachine *VendingMachine) {
	fmt.Println("No drinks available")
}

func (s *SoldOutState) DispenseDrink(vendingMachine *VendingMachine) {
	fmt.Println("No drinks available")
}

// Обьект конекста, использущий состояния
// Дополнительно хранит в себе количество напитков
type VendingMachine struct {
	NoCoinState  State
	HasCoinState State
	SoldState    State
	SoldOutState State
	currentState State
	drinkCount   int
}

// Функция для создания нового обьекта
func NewVendingMachine(drinkCount int) *VendingMachine {
	vendingMachine := &VendingMachine{
		NoCoinState:  &NoCoinState{},
		HasCoinState: &HasCoinState{},
		SoldState:    &SoldState{},
		SoldOutState: &SoldOutState{},
		drinkCount:   drinkCount,
	}
	// Если при создании монет нет - сразу устанавливается состояние SoldOut
	if drinkCount > 0 {
		vendingMachine.currentState = vendingMachine.NoCoinState
	} else {
		vendingMachine.currentState = vendingMachine.SoldOutState
	}
	return vendingMachine
}

// Описание методов машины
func (vm *VendingMachine) SetState(state State) {
	vm.currentState = state
}

func (vm *VendingMachine) InsertCoin() {
	vm.currentState.InsertCoin(vm)
}

func (vm *VendingMachine) EjectCoin() {
	vm.currentState.EjectCoin(vm)
}

func (vm *VendingMachine) SelectDrink() {
	vm.currentState.SelectDrink(vm)
}

func (vm *VendingMachine) DispenseDrink() {
	vm.currentState.DispenseDrink(vm)
}

func (vm *VendingMachine) GetCount() int {
	return vm.drinkCount
}

// Некоторый дополнительный метод-счетчик
func (vm *VendingMachine) ReleaseDrink() {
	fmt.Println("A drink is released")
	vm.drinkCount--
}

func main() {
	vendingMachine := NewVendingMachine(5)

	vendingMachine.InsertCoin()
	vendingMachine.SelectDrink()
	vendingMachine.DispenseDrink()

	vendingMachine.InsertCoin()
	vendingMachine.EjectCoin()

	vendingMachine.SelectDrink()
	vendingMachine.DispenseDrink()

	vendingMachine.SelectDrink()
	vendingMachine.DispenseDrink()

	vendingMachine.InsertCoin()
	vendingMachine.SelectDrink()
	vendingMachine.DispenseDrink()

	vendingMachine.SelectDrink()
	vendingMachine.DispenseDrink()
}
