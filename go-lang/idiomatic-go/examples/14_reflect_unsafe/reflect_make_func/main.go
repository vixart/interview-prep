// reflect.MakeFunc создает функцию во время выполнения: обертка замеряет время
// вызова ЛЮБОЙ функции, не зная ее сигнатуры. Результат приводится к нужному типу
// утверждением: MakeTimedFunction(f).(func(int) int).
package main

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

func MakeTimedFunction(f interface{}) interface{} {
	rf := reflect.TypeOf(f)
	if rf.Kind() != reflect.Func {
		// принимаем ЛЮБУЮ функцию — проверяем это в рантайме
		panic("expects a function")
	}
	vf := reflect.ValueOf(f)
	wrapperF := reflect.MakeFunc(rf, func(in []reflect.Value) []reflect.Value {
		// создаем функцию С ТОЙ ЖЕ сигнатурой, не зная ее заранее
		start := time.Now()
		out := vf.Call(in)
		// вызов исходной функции через рефлексию
		end := time.Now()
		fmt.Printf("calling %s took %v\n", runtime.FuncForPC(vf.Pointer()).Name(), end.Sub(start))
		return out
	})
	return wrapperF.Interface()
}

func timeMe() {
	time.Sleep(1 * time.Second)
}

func timeMeToo(a int) int {
	time.Sleep(time.Duration(a) * time.Second)
	result := a * 2
	return result
}

func main() {
	timed := MakeTimedFunction(timeMe).(func())
	// обратно к нормальной типизации — через утверждение типа
	timed()
	timedToo := MakeTimedFunction(timeMeToo).(func(int) int)
	fmt.Println(timedToo(2))
}
