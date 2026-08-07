// Анонимная функция, присвоенная переменной, и ее вызов в цикле.
package main

import "fmt"

func main() {
	f := func(j int) {
		// анонимная функция — обычное значение, ее можно положить в переменную
		fmt.Println("printing", j, "from inside of an anonymous function")
	}
	for i := 0; i < 5; i++ {
		f(i)
	}
}
