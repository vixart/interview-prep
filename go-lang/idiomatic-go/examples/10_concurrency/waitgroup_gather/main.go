// Сбор результатов от N воркеров: воркеры пишут в общий канал out,
// ОТДЕЛЬНАЯ горутина ждет wg.Wait() и закрывает out — только так for-range по out
// завершится, и при этом никто не закроет канал раньше времени.
package main

import (
	"fmt"
	"sync"
)

func processAndGather[T, R any](in <-chan T, processor func(T) R, num int) []R {
	out := make(chan R, num)
	// буфер по числу воркеров уменьшает блокировки на записи
	var wg sync.WaitGroup
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			defer wg.Done()
			for v := range in {
				out <- processor(v)
			}
		}()
	}
	go func() {
		wg.Wait()
		// отдельная горутина ждет всех воркеров...
		close(out)
		// ...и только потом закрывает out — иначе цикл ниже не завершится
	}()
	var result []R
	for v := range out {
		// читаем, пока out не закроют
		result = append(result, v)
	}
	return result
}

func main() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 20; i++ {
			ch <- i
		}
		close(ch)
	}()
	results := processAndGather(ch, func(i int) int {
		return i * 2
	}, 3)
	fmt.Println(results)
}
