// Демонстрация: чужой пакет меняет неэкспортируемое поле b структуры.
// Так делать не надо — пример показывает, что инкапсуляция в Go не является защитой.
package main

import (
	"fmt"
	"interviewprep/examples/14_reflect_unsafe/unsafe_unexported_field/one_package"
	"interviewprep/examples/14_reflect_unsafe/unsafe_unexported_field/other_package"
)

func main() {
	huf := one_package.HasUnexportedField{
		// поле b здесь заполнить нельзя — оно неэкспортируемое
		A: 10,
		C: "hello",
	}
	fmt.Println(huf)
	other_package.SetBUnsafe(&huf)
	// а через reflect + unsafe чужой пакет его меняет: во втором выводе b == true
	fmt.Println(huf)
}
