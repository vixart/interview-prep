// Синтаксис срезания: x[:2], x[1:], x[1:3], x[:] — какие элементы попадают в результат.
package main

import "fmt"

func main() {
	x := []string{"a", "b", "c", "d"}
	y := x[:2]
	z := x[1:]
	d := x[1:3]
	e := x[:]
	// x[:] — весь срез целиком, тоже без копирования
	fmt.Println("x:", x)
	fmt.Println("y:", y)
	fmt.Println("z:", z)
	fmt.Println("d:", d)
	fmt.Println("e:", e)
}
