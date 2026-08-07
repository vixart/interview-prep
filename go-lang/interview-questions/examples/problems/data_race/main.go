// Data race: две горутины обращаются к одной переменной, хотя бы одна пишет,
// и между обращениями нет отношения «happens before».
//
// Запусти дважды и сравни:
//
//	go run ./problems/data_race          # просто неверный результат, тишина
//	go run -race ./problems/data_race    # WARNING: DATA RACE с двумя стеками
//
// Важно: гонка — это НЕ «неверное число». Гонка — это неопределенное поведение:
// компилятор и процессор вправе переупорядочить операции, а counter++ вообще
// не атомарен (чтение, инкремент, запись — три шага).
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// racy: классика — общий счетчик без синхронизации.
func racy() int {
	counter := 0
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100_000; j++ {
				counter++ // ГОНКА: read-modify-write без синхронизации
			}
		}()
	}
	wg.Wait()
	return counter // почти никогда не 500000
}

// fixedMutex: критическая секция под мьютексом.
func fixedMutex() int {
	counter := 0
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100_000; j++ {
				mu.Lock()
				counter++
				mu.Unlock() // здесь без defer: он в цикле стоил бы дорого
			}
		}()
	}
	wg.Wait()
	return counter
}

// fixedAtomic: для одного числа мьютекс не нужен — хватит atomic.
func fixedAtomic() int64 {
	var counter atomic.Int64
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100_000; j++ {
				counter.Add(1) // одна атомарная инструкция процессора
			}
		}()
	}
	wg.Wait()
	return counter.Load()
}

// fixedSharding: гонки нет вообще — у каждой горутины свой счетчик.
// Часто это самый быстрый вариант: нет ни блокировок, ни борьбы за кеш-линию.
func fixedSharding() int {
	partial := make([]int, 5) // каждый пишет в свой индекс
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 100_000; j++ {
				partial[i]++
			}
		}(i)
	}
	wg.Wait()
	total := 0
	for _, v := range partial {
		total += v
	}
	return total
}

func main() {
	fmt.Println("racy         :", racy(), "(ожидалось 500000)")
	fmt.Println("fixedMutex   :", fixedMutex())
	fmt.Println("fixedAtomic  :", fixedAtomic())
	fmt.Println("fixedSharding:", fixedSharding())

	// Что стоит сказать про детектор гонок:
	// - включается флагом -race у run/test/build;
	// - находит только те гонки, которые РЕАЛЬНО произошли на этом прогоне,
	//   поэтому его гоняют в CI на тестах, а не «один раз посмотрели»;
	// - замедляет программу в 2-20 раз и увеличивает потребление памяти,
	//   в проде не включают.
}
