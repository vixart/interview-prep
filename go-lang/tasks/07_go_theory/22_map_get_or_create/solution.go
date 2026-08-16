// ConcurrentMap.GetOrCreate с корректной работой в конкурентной среде.
//
// Ключевой момент — double-checked locking: быстрая проверка под RLock,
// затем ПОВТОРНАЯ проверка под полным Lock. Без повторной проверки две
// горутины могут одновременно пройти первую проверку, и каждая запишет
// своё значение — вызовы вернут разные результаты.
package main

import (
	"fmt"
	"sync"
)

type ConcurrentMap struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{data: make(map[string]string)}
}

func (m *ConcurrentMap) GetOrCreate(key, value string) string {
	// Быстрый путь: значение уже есть, достаточно RLock.
	m.mu.RLock()
	if v, ok := m.data[key]; ok {
		m.mu.RUnlock()
		return v
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	// Повторная проверка: пока мы ждали Lock, другая горутина могла
	// успеть записать значение — тогда возвращаем его, а не своё.
	if v, ok := m.data[key]; ok {
		return v
	}

	m.data[key] = value
	return value
}

func main() {
	cm := NewConcurrentMap()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		val := cm.GetOrCreate("key1", "value1")
		fmt.Println("Goroutine 1 got:", val)
	}()

	go func() {
		defer wg.Done()
		val := cm.GetOrCreate("key1", "value2")
		fmt.Println("Goroutine 2 got:", val)
	}()

	wg.Wait()
	// Обе горутины гарантированно получают ОДНО И ТО ЖЕ значение —
	// то, которое было записано первым (value1 или value2).
}
