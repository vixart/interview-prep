// Состояние гонки и его исправление: getCounter пишет в общую переменную
// из пяти горутин без синхронизации, getCounterSafe делает то же под мьютексом.
// Детектор гонок: go test -tags racedemo -race ./13_testing/race_detector
package race

import "sync"

// getCounter содержит состояние гонки: пять горутин пишут в counter без синхронизации.
// Запусти `go test -tags racedemo -race ./13_testing/race_detector` — детектор покажет гонку.
func getCounter() int {
	var counter int
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				counter++
				// ГОНКА: пять горутин читают и пишут counter одновременно
			}
		}()
	}
	wg.Wait()
	return counter
}

// getCounterSafe — исправленная версия: доступ к counter защищен мьютексом.
func getCounterSafe() int {
	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				mu.Lock()
				// исправление: критическая секция под мьютексом
				counter++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return counter
}
