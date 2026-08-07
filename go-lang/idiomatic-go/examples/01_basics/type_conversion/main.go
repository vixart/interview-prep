// В Go нет неявных числовых преобразований: int + float64 и int + byte
// не компилируются без явного float64(x) / int(b). Убери любое преобразование — сборка упадет.
package main

import "fmt"

// remove either type conversion and the code will not compile!
func main() {
	var x int = 10
	var y float64 = 30.2

	// converts the int to a float64, so it can be added to a float64
	var sum1 float64 = float64(x) + y
	// без float64(x) — ошибка компиляции: int и float64 не смешиваются

	// converts the float64 to an int, so it can be added to an int
	var sum2 int = x + int(y)
	fmt.Println(sum1, sum2)

	var b byte = 100

	// need to convert ints to bytes or vice versa as well
	var sum3 int = x + int(b)
	var sum4 byte = byte(x) + b
	// то же правило для int и byte: неявных преобразований в Go нет
	fmt.Println(sum3, sum4)
}
