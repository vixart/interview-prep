// Когда вывод типов не работает: T2 встречается только в возвращаемом значении,
// поэтому параметры типа приходится указывать явно — Convert[int, int64](a).
package main

import (
	"fmt"
	"reflect"
)

type Integer interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

func Convert[T1, T2 Integer](in T1) T2 {
	// T2 нигде во ВХОДНЫХ параметрах не встречается
	return T2(in)
}

func main() {
	var a int = 10
	b := Convert[int, int64](a)
	// поэтому типы приходится указывать явно — вывести их не из чего
	fmt.Println(b, reflect.TypeOf(b))
}
