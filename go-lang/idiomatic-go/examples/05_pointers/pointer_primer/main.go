// Базовая механика указателей: & и *, nil-указатель (разыменование паникует —
// здесь оно перехвачено recover), new(int), и указатель как способ показать
// необязательное поле структуры (MiddleName).
package main

import "fmt"

func main() {
	firstExample()
	secondExample()
	thirdExample()
	fourthExample()
	fifthExample()
}

func firstExample() {
	var x int32 = 10
	var y bool = true
	pointerX := &x
	pointerY := &y
	var pointerZ *string
	// указатель без значения — nil

	fmt.Println(x, y, pointerX, pointerY, pointerZ)
}

func secondExample() {
	x := 10
	pointerToX := &x
	fmt.Println(pointerToX)  // prints a memory address
	fmt.Println(*pointerToX) // prints 10
	// * разыменовывает: печатает 10, а не адрес
	z := 5 + *pointerToX
	fmt.Println(z) // prints 15
}

func thirdExample() {
	// chapter 9 explains panic and recover
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	var x *int
	fmt.Println(x == nil) // prints true
	fmt.Println(*x)       // panics
	// разыменование nil — паника (здесь ее ловит recover выше)
}

func fourthExample() {
	var x = new(int)
	// new выделяет память под нулевое значение и возвращает указатель
	fmt.Println(x == nil) // prints false
	fmt.Println(*x)       // prints 0
}

type person struct {
	FirstName  string
	MiddleName *string
	LastName   string
}

func stringp(s string) *string {
	return &s
}

func fifthExample() {
	p := person{
		FirstName:  "Pat",
		MiddleName: stringp("Perry"), // This works
		// указатель как способ отличить «нет значения» от пустой строки
		LastName: "Peterson",
	}
	fmt.Println(p)
}
