// Исправленный кеш: RWMutex против гонок + singleflight-подход против
// повторных дорогих вычислений (thundering herd).
//
// Что изменилось относительно исходника:
//   - в Cache добавлен sync.RWMutex: Get под RLock, Set/Delete под Lock —
//     уходит fatal error: concurrent map writes;
//   - GetOrCompute стал атомарным по ключу: в кеше хранится «обещание»
//     (канал done), поэтому дорогое вычисление для ключа выполняется
//     ровно один раз, остальные горутины ждут готовый результат.
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type entry struct {
	done  chan struct{}
	value string
}

type Cache struct {
	mu   sync.RWMutex
	data map[string]*entry
}

func NewCache() *Cache {
	return &Cache{data: make(map[string]*entry)}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	e, ok := c.data[key]
	c.mu.RUnlock()

	if !ok {
		return "", false
	}
	<-e.done // дождаться готовности значения
	return e.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

var computations atomic.Int64

// expensiveComputation симулирует дорогое вычисление.
func expensiveComputation(key string) string {
	computations.Add(1)
	time.Sleep(100 * time.Millisecond)
	return fmt.Sprintf("result for %s", key)
}

// GetOrCompute получает значение из кеша или вычисляет его ровно один раз.
func GetOrCompute(cache *Cache, key string) string {
	cache.mu.Lock()
	if e, ok := cache.data[key]; ok {
		cache.mu.Unlock()
		<-e.done // кто-то уже считает или посчитал
		return e.value
	}

	e := &entry{done: make(chan struct{})}
	cache.data[key] = e
	cache.mu.Unlock()

	e.value = expensiveComputation(key) // без блокировки кеша
	close(e.done)

	return e.value
}

func main() {
	cache := NewCache()
	var wg sync.WaitGroup

	keys := []string{"user:1", "user:2", "user:1", "user:2", "user:1", "user:1"}
	for _, key := range keys {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			fmt.Println(GetOrCompute(cache, key))
		}(key)
	}
	wg.Wait()

	fmt.Println("вычислений:", computations.Load()) // 2 — по одному на ключ
}
