// Функция возвращает функцию: makeMult захватывает base и порождает
// независимые умножители twoBase и threeBase.
package main

import "fmt"

func makeMult(base int) func(int) int {
	return func(factor int) int {
		// возвращаемая функция захватывает base — у каждой свой
		return base * factor
	}
}

func main() {
	twoBase := makeMult(2)
	// два независимых замыкания с разным base
	threeBase := makeMult(3)
	for i := 0; i < 3; i++ {
		fmt.Println(twoBase(i), threeBase(i))
	}
}
