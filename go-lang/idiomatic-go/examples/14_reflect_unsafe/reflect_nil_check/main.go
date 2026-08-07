// Как проверить, что внутри интерфейса лежит nil: обычное `== nil` возвращает false
// для типизированного nil, а reflect (IsValid + IsNil для Ptr/Slice/Map/Func/Interface) — true.
// Смотри третью строку вывода: `false true`.
package main

import (
	"fmt"
	"reflect"
)

func hasNoValue(i interface{}) bool {
	iv := reflect.ValueOf(i)
	if !iv.IsValid() {
		// невалидное Value = внутри интерфейса вообще ничего нет
		return true
	}
	switch iv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
		return iv.IsNil()
		// IsNil допустим только для этих разновидностей — иначе паника
	default:
		return false
	}
}

func main() {
	var a interface{}
	fmt.Println(a == nil, hasNoValue(a)) // prints true true

	var b *int
	fmt.Println(b == nil, hasNoValue(b)) // prints true true

	var c interface{} = b
	// типизированный nil: c == nil даст false...
	fmt.Println(c == nil, hasNoValue(c)) // prints false true
	// ...а рефлексия честно скажет, что значения нет

	var d int
	fmt.Println(hasNoValue(d)) // prints false

	var e interface{} = d
	fmt.Println(e == nil, hasNoValue(e)) // prints false false
}
