// GetOrCompute: кеш дорогих вычислений, безопасный для конкурентного
// доступа, без повторного вычисления одного ключа (анти-thundering herd).
//
// Приём: в кеше храним не готовое значение, а "обещание" — канал, который
// закрывается, когда значение вычислено. Первая горутина по ключу создаёт
// обещание и считает; остальные ждут закрытия канала. Так дорогая функция
// выполняется ровно один раз на ключ (то, что делает singleflight).
package main

import (
	"fmt"
	"sync"
	"time"
)

type promise struct {
	done  chan struct{}
	value any
}

type Cache struct {
	mu   sync.Mutex
	data map[string]*promise
}

func NewCache() *Cache {
	return &Cache{data: make(map[string]*promise)}
}

func (c *Cache) GetOrCompute(key string, compute func() any) any {
	c.mu.Lock()
	if p, ok := c.data[key]; ok {
		c.mu.Unlock()
		<-p.done // значение уже считается или готово — ждём
		return p.value
	}

	p := &promise{done: make(chan struct{})}
	c.data[key] = p
	c.mu.Unlock()

	p.value = compute() // тяжёлая работа выполняется без блокировки кеша
	close(p.done)

	return p.value
}

func main() {
	cache := NewCache()
	calls := 0

	expensive := func() any {
		calls++ // считаем реальные вычисления
		time.Sleep(100 * time.Millisecond)
		return "result"
	}

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(cache.GetOrCompute("db-query", expensive))
		}()
	}
	wg.Wait()

	fmt.Println("вычислений:", calls) // 1 — остальные 4 дождались готового
}
