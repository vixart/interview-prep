// Преобразование среза в массив и в указатель на массив.
// [4]int(s) — КОПИЯ (изменение среза не видно в массиве), (*[4]int)(s) — ОБЩАЯ память.
// Последняя функция паникует: массив [5]int длиннее среза из 4 элементов (проверка в рантайме).
package main

import "fmt"

func main() {
	arrayConversions()
	arrayPointerConversions()
	panicArrayConversions()
}

func arrayConversions() {
	xSlice := []int{1, 2, 3, 4}
	xArray := [4]int(xSlice)
	// срез → массив: данные КОПИРУЮТСЯ
	smallArray := [2]int(xSlice)
	xSlice[0] = 10
	// меняем срез после копирования...
	fmt.Println(xSlice)
	fmt.Println(xArray)
	// ...массив остался прежним: [1 2 3 4]
	fmt.Println(smallArray)
}

func arrayPointerConversions() {
	xSlice := []int{1, 2, 3, 4}
	xArrayPointer := (*[4]int)(xSlice)
	// срез → УКАЗАТЕЛЬ на массив: память общая, копирования нет
	xSlice[0] = 10
	xArrayPointer[1] = 20
	// запись через указатель видна в срезе, и наоборот
	fmt.Println(xSlice)
	fmt.Println(xArrayPointer)
}

func panicArrayConversions() {
	xSlice := []int{1, 2, 3, 4}
	panicArray := [5]int(xSlice)
	// ПАНИКА: массив длиннее среза, проверка только в рантайме
	fmt.Println(panicArray)
}
