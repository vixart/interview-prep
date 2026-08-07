// Два свойства defer: порядок LIFO (second печатается раньше first)
// и вычисление аргументов В МОМЕНТ объявления defer (val = 10 и 20, хотя a стала 30).
package main

import "fmt"

func main() {
	deferExample()
}

func deferExample() int {
	a := 10
	defer func(val int) {
		fmt.Println("first:", val)
	}(a)
	// аргумент вычисляется ЗДЕСЬ: val = 10, хотя ниже a изменится
	a = 20
	defer func(val int) {
		fmt.Println("second:", val)
	}(a)
	// val = 20 — снимок значения на момент объявления defer
	a = 30
	fmt.Println("exiting:", a)
	return a
	// при выходе отложенные функции идут в обратном порядке: second, потом first
}
