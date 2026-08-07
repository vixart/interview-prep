// Функция — обычное значение: переменной типа func(string) int можно присвоить
// любую подходящую функцию и вызвать через нее. Нулевое значение такой переменной — nil.
package main

import "fmt"

func f1(a string) int {
	return len(a)
}

func f2(a string) int {
	total := 0
	for _, v := range a {
		total += int(v)
	}
	return total
}

func main() {
	var myFuncVariable func(string) int
	// тип переменной — сигнатура функции; нулевое значение nil
	myFuncVariable = f1
	// в одну переменную ложится любая функция с такой же сигнатурой
	result := myFuncVariable("Hello")
	fmt.Println(result)

	myFuncVariable = f2
	result = myFuncVariable("Hello")
	fmt.Println(result)
}
