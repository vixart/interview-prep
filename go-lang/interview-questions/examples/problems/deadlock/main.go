// Блокировки: четыре способа получить вечное ожидание и как их избежать.
//
// Каждый «плохой» сценарий запускается в отдельной горутине и ждется с
// таймаутом — иначе программа просто зависла бы. В реальном коде такое
// зависание выглядит как «сервис перестал отвечать», без паники и без логов.
package main

import (
	"fmt"
	"sync"
	"time"
)

// tryFor запускает f и сообщает, успела ли она за отведенное время.
func tryFor(name string, d time.Duration, f func()) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		f()
	}()
	select {
	case <-done:
		fmt.Printf("%-28s завершилось\n", name)
	case <-time.After(d):
		fmt.Printf("%-28s ЗАВИСЛО (горутина осталась заблокированной)\n", name)
	}
}

// 1. Мьютексы в Go НЕ реентерабельны: повторный Lock в той же горутине висит.
func selfLock() {
	var mu sync.Mutex
	mu.Lock()
	mu.Lock() // ждет сам себя
	mu.Unlock()
	mu.Unlock()
}

// Та же ошибка, но неочевидная: метод под блокировкой зовет другой метод,
// который берет ту же блокировку.
type store struct {
	mu sync.Mutex
	m  map[string]int
}

func (s *store) get(k string) int { // публичный метод берет замок
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.getLocked(k)
}
func (s *store) getLocked(k string) int { return s.m[k] } // приватный — уже под замком
func (s *store) sum(keys []string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	total := 0
	for _, k := range keys {
		total += s.getLocked(k) // ПРАВИЛЬНО: зовем *Locked-версию, а не s.get
	}
	return total
}

// 2. Разный порядок захвата двух мьютексов — классический deadlock.
func lockOrdering() {
	var a, b sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { // берет a, потом b
		defer wg.Done()
		a.Lock()
		time.Sleep(10 * time.Millisecond)
		b.Lock()
		b.Unlock()
		a.Unlock()
	}()
	go func() { // берет b, потом a  ← вот здесь и ломается
		defer wg.Done()
		b.Lock()
		time.Sleep(10 * time.Millisecond)
		a.Lock()
		a.Unlock()
		b.Unlock()
	}()
	wg.Wait()
}

// 3. Запись в небуферизованный канал без читателя.
func sendWithoutReceiver() {
	ch := make(chan int)
	ch <- 1 // читателя нет — ждем вечно
}

// 4. RWMutex: писатель встает в очередь и блокирует новых читателей,
// поэтому RLock внутри RLock тоже способен зависнуть.
func rwUpgrade() {
	var mu sync.RWMutex
	mu.RLock()
	mu.Lock() // «повысить» read-lock до write-lock нельзя
	mu.Unlock()
	mu.RUnlock()
}

func main() {
	tryFor("1. повторный Lock", 200*time.Millisecond, selfLock)
	tryFor("2. разный порядок замков", 300*time.Millisecond, lockOrdering)
	tryFor("3. запись без читателя", 200*time.Millisecond, sendWithoutReceiver)
	tryFor("4. RLock → Lock", 200*time.Millisecond, rwUpgrade)

	s := &store{m: map[string]int{"a": 1, "b": 2}}
	fmt.Println("правильный вызов под замком:", s.sum([]string{"a", "b"}))

	// Что помогает не попадать:
	// - блокировка на минимальном участке, всегда defer Unlock;
	// - не вызывать чужой код (колбэки, интерфейсы) под блокировкой;
	// - единый порядок захвата, если замков несколько;
	// - пара методов f/fLocked, чтобы не брать замок дважды;
	// - если все горутины заснули — рантайм сам печатает
	//   "fatal error: all goroutines are asleep - deadlock!" со стеками,
	//   но при живом HTTP-сервере этого не случится: смотри /debug/pprof/goroutine.
}
