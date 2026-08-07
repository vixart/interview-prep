// Константы: типизированные и нетипизированные, на уровне пакета и внутри функции.
// Главное: константа — это имя для литерала, ее нельзя изменить (строки с присвоением
// закомментированы, рядом приведен текст ошибки компилятора).
package main

import "fmt"

const x int64 = 10

// типизированная константа: ее можно присвоить только int64

const (
	idKey   = "id"
	nameKey = "name"
)

const z = 20 * 10

// нетипизированная: считается компилятором, тип возьмется из места использования

// this code will not compile
// ./main.go:23:2: cannot assign to x (constant 10 of type int64)
// ./main.go:24:2: cannot assign to y (untyped string constant "hello")
// on the Go Playground at https://oreil.ly/FdG-W
func main() {
	const y = "hello"
	// константы объявляются и внутри функции — область видимости обычная

	fmt.Println(x)
	fmt.Println(y)

	// Не компилируется: константу нельзя изменить.
	// x = x + 1 // cannot assign to x (constant 10 of type int64)
	// y = "bye" // cannot assign to y (untyped string constant "hello")
}
