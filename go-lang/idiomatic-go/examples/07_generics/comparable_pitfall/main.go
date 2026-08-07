// Дыра в типобезопасности comparable: интерфейс Thinger считается сравнимым на этапе
// компиляции, но если внутри него лежит срез — сравнение ПАНИКУЕТ в рантайме.
// reflect.Value.Comparable() показывает, что реальная сравнимость известна только в рантайме.
package main

import (
	"fmt"
	"reflect"
)

type Thinger interface {
	Thing()
}

type ThingerInt int

func (t ThingerInt) Thing() {
	fmt.Println("ThingInt:", t)
}

type ThingerSlice []int

func (t ThingerSlice) Thing() {
	fmt.Println("ThingSlice:", t)
}

func Comparer[T comparable](t1, t2 T) {
	// интерфейс формально удовлетворяет comparable — компилятор пропустит
	v1 := reflect.ValueOf(t1)
	v2 := reflect.ValueOf(t2)
	fmt.Println("comparable v1:", v1.Comparable())
	// реальная сравнимость известна только в рантайме
	fmt.Println("comparable v2:", v2.Comparable())
	if t1 == t2 {
		fmt.Println("equal!")
	}
}

func main() {
	var a int = 10
	var b int = 10
	Comparer(a, b)

	var a2 ThingerInt = 20
	var b2 ThingerInt = 20
	Comparer(a2, b2)

	var a3 ThingerSlice = []int{1, 2, 3}
	var b3 ThingerSlice = []int{1, 2, 3}
	//Comparer(a3, b3)

	var a4 Thinger = a2
	var b4 Thinger = b2
	Comparer(a4, b4)
	// внутри интерфейсов лежит ThingerInt — сравнение проходит нормально

	a4 = a3
	b4 = b3
	Comparer(a4, b4)
	// ПАНИКА: теперь внутри интерфейсов срезы, а их сравнивать нельзя
	// (до этой строки программа не дойдет — остальное показано для полноты)

	b4 = b2
	Comparer(a4, b4)
}
