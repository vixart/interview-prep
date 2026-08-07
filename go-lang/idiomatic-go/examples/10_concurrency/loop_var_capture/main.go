// Захват переменной цикла горутиной.
// До Go 1.22: переменная цикла ОДНА на весь цикл — все горутины видели последнее значение.
// С Go 1.22 (и go >= 1.22 в go.mod): на каждой итерации создается новая переменная.
package main

import (
	"fmt"
	"sort"
	"sync"
)

func main() {
	vals := []int{2, 4, 6, 8, 10}

	// Современный Go: безопасно, каждая горутина получает свою v.
	ch := make(chan int, len(vals))
	var wg sync.WaitGroup
	for _, v := range vals {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch <- v * 2
		}()
	}
	wg.Wait()
	close(ch)

	var got []int
	for v := range ch {
		got = append(got, v)
	}
	sort.Ints(got)
	fmt.Println("captured:", got)

	// Способ, работающий на ЛЮБОЙ версии языка: передать значение параметром.
	// Так же поступают, когда переменная меняется вне горутины.
	ch2 := make(chan int, len(vals))
	wg.Add(len(vals))
	for _, v := range vals {
		go func(val int) {
			defer wg.Done()
			ch2 <- val * 2
		}(v)
	}
	wg.Wait()
	close(ch2)
	fmt.Println("len:", len(ch2))
}
