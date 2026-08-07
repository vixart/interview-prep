// Базовые операции с map: запись, чтение, чтение отсутствующего ключа (нулевое значение)
// и ++ по отсутствующему ключу (работает, потому что читается 0).
package main

import "fmt"

func main() {
	totalWins := map[string]int{}
	totalWins["Orcas"] = 1
	totalWins["Lions"] = 2
	fmt.Println(totalWins["Orcas"])
	fmt.Println(totalWins["Kittens"])
	// ключа нет → нулевое значение типа значения (0), а не ошибка
	totalWins["Kittens"]++
	// работает даже для отсутствующего ключа: читается 0, записывается 1
	fmt.Println(totalWins["Kittens"])
	totalWins["Lions"] = 3
	fmt.Println(totalWins["Lions"])
}
