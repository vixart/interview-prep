// Ограничение, которое невозможно удовлетворить: базовый тип int И метод String().
// int методов не имеет, а MyInt не проходит из-за отсутствия ~ в списке типов.
// Мораль: список типов и требование методов почти всегда несовместимы.
package main

import (
	"fmt"
	"strconv"
)

type ImpossiblePrintableInt interface {
	// ограничение требует одновременно базовый тип int И метод String()
	int
	String() string
}

type ImpossibleStruct[T ImpossiblePrintableInt] struct {
	val T
}

type MyInt int

func (mi MyInt) String() string {
	// MyInt метод имеет, но в ограничении нет ~int — значит все равно не подходит
	return strconv.Itoa(int(mi))
}

func main() {
	// Ограничение ImpossiblePrintableInt не может быть удовлетворено никогда:
	// требуется одновременно базовый тип int и метод String().
	// s := ImpossibleStruct[int]{10}    // int does not satisfy: missing method String
	// s2 := ImpossibleStruct[MyInt]{10} // MyInt does not satisfy: possibly missing ~ for int
	fmt.Println("см. комментарии в коде: такой интерфейс-ограничение нереализуем")
}
