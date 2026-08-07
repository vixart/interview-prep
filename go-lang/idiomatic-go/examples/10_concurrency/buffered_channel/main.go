// Буферизованный канал там, где заранее известно КОЛИЧЕСТВО результатов:
// 10 горутин, буфер на 10 — никто не блокируется, main спокойно вычитывает все ответы.
package main

import "fmt"

func processChannel(ch chan int) []int {
	const conc = 10
	results := make(chan int, conc)
	// буфер по числу горутин: никто не заблокируется на записи
	for i := 0; i < conc; i++ {
		go func() {
			v := <-ch
			// каждая горутина забирает ровно одно значение
			results <- process(v)
		}()
	}
	var out []int
	for i := 0; i < conc; i++ {
		out = append(out, <-results)
		// собираем conc результатов — их количество известно заранее
	}
	return out
}

func process(i int) int {
	// this should be a more complicated operation to make concurrency worthwhile
	return i * 2
}

func main() {
	vals := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	ch := make(chan int)
	go func() {
		for _, v := range vals {
			ch <- v
		}
	}()
	result := processChannel(ch)
	fmt.Println(result)
}
