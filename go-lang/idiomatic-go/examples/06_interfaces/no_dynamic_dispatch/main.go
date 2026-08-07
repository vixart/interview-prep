// В Go нет динамической диспетчеризации: Inner.Double вызывает Inner.IntPrinter,
// а не переопределенный Outer.IntPrinter. Вывод — "Inner: 20".
// Встраивание не делает методы виртуальными.
package main

import "fmt"

type Inner struct {
	A int
}

func (i Inner) IntPrinter(val int) string {
	return fmt.Sprintf("Inner: %d", val)
}

func (i Inner) Double() string {
	result := i.A * 2
	return i.IntPrinter(result)
	// ЗДЕСЬ вызовется Inner.IntPrinter, а не Outer.IntPrinter — переопределения нет
}

type Outer struct {
	Inner
	S string
}

func (o Outer) IntPrinter(val int) string {
	// этот метод затеняет одноименный только при вызове через Outer напрямую
	return fmt.Sprintf("Outer: %d", val)
}

func main() {
	o := Outer{
		Inner: Inner{
			A: 10,
		},
		S: "Hello",
	}
	fmt.Println(o.Double())
	// печатает "Inner: 20" — Inner ничего не знает о том, что он встроен
}
