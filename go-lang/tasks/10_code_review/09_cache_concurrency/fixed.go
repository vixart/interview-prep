// Исправленный кеш дорогих вычислений.
//
// Что изменилось относительно исходника:
//   - мьютекс живёт в структуре кеша рядом с данными (в оригинале
//     var mu объявлялась внутри функции — у каждого вызова свой мьютекс,
//     синхронизации не было вовсе);
//   - убран Unlock незалоченного мьютекса (была паника
//     "sync: unlock of unlocked mutex");
//   - глобальная мапа заменена типом Cache с методами;
//   - повторная проверка под Lock: если пока считали, кто-то уже записал
//     значение — возвращаем его (double-checked locking);
//   - исправлен перевод секунд в Duration: float64 секунд надо умножать
//     на time.Second, иначе Duration трактует число как наносекунды.
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func LongCalculation(n int) int {
	secondsToSleep := rand.Float64() * float64(n)
	time.Sleep(time.Duration(secondsToSleep * float64(time.Second) / 100)) // /100 — ускорено для демо
	return n + 1
}

type Cache struct {
	mu   sync.RWMutex
	data map[int]int
}

func NewCache() *Cache {
	return &Cache{data: make(map[int]int)}
}

func (c *Cache) Get(n int) int {
	c.mu.RLock()
	v, ok := c.data[n]
	c.mu.RUnlock()

	if ok {
		return v
	}

	value := LongCalculation(n)

	c.mu.Lock()
	defer c.mu.Unlock()

	if v, ok := c.data[n]; ok { // кто-то успел раньше — берём его результат
		return v
	}
	c.data[n] = value

	return value
}

func main() {
	cache := NewCache()
	nums := []int{5, 10, 22}

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ { // несколько конкурентных проходов
		for _, n := range nums {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				fmt.Println(n, "->", cache.Get(n))
			}(n)
		}
	}
	wg.Wait()
}
