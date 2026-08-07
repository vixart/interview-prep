// Каноническая схема «пул воркеров»: канал in с заданиями, канал out с результатами,
// N горутин читают in через for-range (цикл завершится, когда in закроют),
// отдельная горутина наполняет in и закрывает его, main собирает ровно len(inVals) результатов.
package main

import "fmt"

func main() {
	x := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	result := processConcurrently(x)
	fmt.Println(result)
}

func process(val int) int {
	// do something with val
	return val * 2
}

const numGoroutines = 5

func processConcurrently(inVals []int) []int {
	// create the channels
	in := make(chan int, numGoroutines)
	// канал заданий
	out := make(chan int, numGoroutines)
	// канал результатов
	// launch numGoroutines
	for i := 0; i < numGoroutines; i++ {
		go func() {
			for val := range in {
				// воркер читает, пока канал не закроют — тогда цикл сам завершится
				result := process(val)
				out <- result
			}
		}()
	}
	// load the data into the channel in another goroutine
	go func() {
		for _, v := range inVals {
			in <- v
		}
		close(in)
		// закрывает тот, кто ПИШЕТ; иначе воркеры зависли бы навсегда
	}()
	// read the data
	outVals := make([]int, 0, len(inVals))
	for i := 0; i < len(inVals); i++ {
		// читаем ровно столько результатов, сколько отправили заданий
		outVals = append(outVals, <-out)
	}
	return outVals
}
