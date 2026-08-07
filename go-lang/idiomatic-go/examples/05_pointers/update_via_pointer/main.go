// Разница между «переприсвоить указатель» и «изменить значение по указателю»:
// failedUpdate меняет локальную копию указателя (снаружи ничего не произошло),
// update пишет через *px — вот это и видно вызывающему.
package main

import "fmt"

func failedUpdate(px *int) {
	x2 := 20
	px = &x2
	// переприсваиваем КОПИЮ указателя — снаружи ничего не изменилось
}

func update(px *int) {
	*px = 20
	// а вот запись ПО указателю меняет исходную переменную
}

func main() {
	x := 10
	failedUpdate(&x)
	fmt.Println(x) // prints 10
	update(&x)
	fmt.Println(x) // prints 20
}
