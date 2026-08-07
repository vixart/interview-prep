// sync.WaitGroup: Add(3) вызывается ДО запуска горутин, каждая делает defer wg.Done(),
// main блокируется на wg.Wait(), пока счетчик не станет нулем.
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3)
	// Add вызывается ДО запуска горутин, иначе возможна гонка со Wait
	go func() {
		defer wg.Done()
		// через defer — чтобы счетчик уменьшился даже при панике
		doThing1()
	}()
	go func() {
		defer wg.Done()
		doThing2()
	}()
	go func() {
		defer wg.Done()
		doThing3()
	}()
	wg.Wait()
	// блокирует main, пока счетчик не станет нулем
}

func doThing1() {
	fmt.Println("Thing 1 done!")
}

func doThing2() {
	fmt.Println("Thing 2 done!")
}

func doThing3() {
	fmt.Println("Thing 3 done!")
}
