// Отрефакторенный вариант.
//
// Что изменилось относительно исходника:
//   - убран лишний мьютекс (канал потокобезопасен);
//   - убран select с единственным case (эквивалентен обычному чтению);
//   - канал закрывается после wg.Wait() в отдельной горутине,
//     чтение через for range -> программа завершается, нет deadlock;
//   - fmt.Print вместо fmt.Printf(result) — данные не должны быть
//     формат-строкой (go vet это ловит);
//   - нормальные имена: wg, results; %d вместо strconv.Itoa + %s.
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	results := make(chan string, 5)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results <- fmt.Sprintf("Current gorutine number: %d\n", i)
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Print(result)
	}
}
