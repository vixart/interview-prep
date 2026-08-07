// Список типов в ограничении и знак ~: `~int` означает «int и любой тип на его основе»,
// поэтому divAndRemainder работает и с uint, и с локальным типом MyInt.
// Без тильды MyInt ограничению не удовлетворял бы.
package main

import (
	"errors"
	"fmt"
)

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
	// тильда = «этот тип И любой пользовательский на его основе»
}

func divAndRemainder[T Integer](num, denom T) (T, T, error) {
	if denom == 0 {
		return 0, 0, errors.New("cannot divide by zero")
	}
	return num / denom, num % denom, nil
	// операторы / и % разрешены именно списком типов в ограничении
}

func main() {
	var a uint = 18_446_744_073_709_551_615
	var b uint = 9_223_372_036_854_775_808
	fmt.Println(divAndRemainder(a, b))

	type MyInt int
	// без ~ в ограничении этот тип не подошел бы
	var myA MyInt = 10
	var myB MyInt = 20
	fmt.Println(divAndRemainder(myA, myB))

}
