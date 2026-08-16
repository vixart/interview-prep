// Fan-in: merge N каналов в один.
//
// На каждый входной канал — горутина, переливающая значения в out.
// Отдельная горутина ждёт wg.Wait() и закрывает out — закрыть канал можно
// только один раз и только когда все писатели закончили.
package main

import "sync"

func merge(channels ...chan int64) <-chan int64 {
	out := make(chan int64)

	var wg sync.WaitGroup
	wg.Add(len(channels))

	for _, ch := range channels {
		go func(ch <-chan int64) {
			defer wg.Done()

			for v := range ch {
				out <- v
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	channels := make([]chan int64, 10)
	for i := range channels {
		channels[i] = make(chan int64)
	}

	for i := range channels {
		go func(i int) {
			channels[i] <- int64(i)
			close(channels[i])
		}(i)
	}

	sum := int64(0)
	n := 0
	for v := range merge(channels...) {
		println(v) // порядок не гарантирован
		sum += v
		n++
	}
	println("count:", n, "sum:", sum) // count: 10 sum: 45
}
