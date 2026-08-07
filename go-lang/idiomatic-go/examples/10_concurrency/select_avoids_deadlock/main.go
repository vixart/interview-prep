// Тот же код, что в deadlock, но main ждет обе операции через select
// и выполняет ту, что готова. Взаимоблокировки нет.
package main

import (
	"fmt"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go func() {
		v := 1
		ch1 <- v
		v2 := <-ch2
		fmt.Println(v, v2)
	}()
	v := 2
	var v2 int
	select {
	// main готов и записать в ch2, и прочитать из ch1 — выполнится то, что готово
	case ch2 <- v:
	case v2 = <-ch1:
		// сработает эта ветвь: горутина уже пишет в ch1
	}
	fmt.Println(v, v2)
}
