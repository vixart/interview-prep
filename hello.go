package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg1 sync.WaitGroup
	wg1.Add(2)
	ch := make(chan int)
	go func() {
		defer wg1.Done()
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()
	go func() {
		defer wg1.Done()
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	go func() {
		wg1.Wait()
		close(ch)
	}()

	var wg2 sync.WaitGroup
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		for v := range ch {
			fmt.Println(v)
		}
	}()
	wg2.Wait()
}
