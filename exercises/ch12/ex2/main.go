package main

import "fmt"

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			ch1 <- i
		}
		close(ch1)
	}()

	go func() {
		for i := 0; i < 10; i++ {
			ch2 <- i
		}
		close(ch2)
	}()

	for ch1 != nil || ch2 != nil {
		select {
		case v, ok := <-ch1:
			if !ok {
				ch1 = nil
				continue
			}
			fmt.Println("ch1:", v)
		case v, ok := <-ch2:
			if !ok {
				ch2 = nil
				continue
			}
			fmt.Println("ch2:", v)
		}
	}
}
