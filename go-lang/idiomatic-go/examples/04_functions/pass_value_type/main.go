// Простые типы, строки и структуры копируются при передаче:
// изменения внутри функции наружу не видны.
package main

import "fmt"

type person struct {
	age  int
	name string
}

func modifyFails(i int, s string, p person) {
	i = i * 2
	s = "Goodbye"
	p.name = "Bob"
	// все три параметра — копии, снаружи ничего не изменится
}

func main() {
	p := person{}
	i := 2
	s := "Hello"
	modifyFails(i, s, p)
	fmt.Println(i, s, p)
}
