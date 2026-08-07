// Обобщенный стек: Stack[T any]. Показаны параметр типа у структуры и у методов,
// идиома `var zero T` для возврата нулевого значения и типобезопасность
// (закомментированный Push строки в Stack[int] не компилируется).
package main

import (
	"fmt"
)

type Stack[T any] struct {
	// T — параметр типа, any — ограничение (годится любой тип)
	vals []T
}

func (s *Stack[T]) Push(val T) {
	// в методах параметр типа повторяется у приемника: Stack[T]
	s.vals = append(s.vals, val)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.vals) == 0 {
		var zero T
		// идиома: так возвращают нулевое значение неизвестного типа
		return zero, false
	}
	top := s.vals[len(s.vals)-1]
	s.vals = s.vals[:len(s.vals)-1]
	return top, true
}

func main() {
	var intStack Stack[int]
	// подставляем конкретный тип; нулевое значение уже готово к работе
	intStack.Push(10)
	intStack.Push(20)
	intStack.Push(30)
	v, ok := intStack.Pop()
	fmt.Println(v, ok)
	v, ok = intStack.Pop()
	fmt.Println(v, ok)
	v, ok = intStack.Pop()
	fmt.Println(v, ok)
	v, ok = intStack.Pop()
	fmt.Println(v, ok)
	// intStack.Push("nope") // ошибка компиляции: Stack[int] принимает только int
}
