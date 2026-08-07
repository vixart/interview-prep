// Взаимоблокировка: горутина ждет, пока прочитают из ch1, main ждет, пока прочитают из ch2.
// Программа падает с `fatal error: all goroutines are asleep - deadlock!`.
// Лечение — в примере select_avoids_deadlock.
package main

import (
	"fmt"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go func() {
		inGoroutine := 1
		ch1 <- inGoroutine
		// горутина блокируется здесь: читателя у ch1 пока нет
		fromMain := <-ch2
		fmt.Println("goroutine:", inGoroutine, fromMain)
	}()
	inMain := 2
	ch2 <- inMain
	// main блокируется здесь: горутина до чтения ch2 не дошла → deadlock
	fromGoroutine := <-ch1
	fmt.Println("main:", inMain, fromGoroutine)
}
