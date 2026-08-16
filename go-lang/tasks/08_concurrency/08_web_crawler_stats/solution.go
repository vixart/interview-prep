// Web crawler: количество загрузок хоста за последние 10 минут.
//
// Та же идея, что и в realtime-аналитике, но бакеты по МИНУТАМ (их всего 10),
// т.к. объёмы большие, а точность до минуты достаточна:
//   - OnPageDownloaded: O(1);
//   - Count: O(10).
//
// Память на хост — 10 счётчиков. Для меньшей конкуренции мьютекс можно
// шардировать по хосту, здесь для простоты один RWMutex.
package main

import (
	"fmt"
	"sync"
	"time"
)

const windowMin = 10

type hostBuckets struct {
	counts [windowMin]int
	stamps [windowMin]int64 // номер минуты, которой принадлежит бакет
}

type Solution struct {
	mu    sync.RWMutex
	hosts map[string]*hostBuckets
	now   func() time.Time // подменяемо в тестах
}

func NewSolution() *Solution {
	return &Solution{hosts: make(map[string]*hostBuckets), now: time.Now}
}

func (s *Solution) OnPageDownloaded(host string, timestamp time.Time) {
	minute := timestamp.Unix() / 60
	idx := minute % windowMin

	s.mu.Lock()
	defer s.mu.Unlock()

	h, ok := s.hosts[host]
	if !ok {
		h = &hostBuckets{}
		s.hosts[host] = h
	}

	if h.stamps[idx] != minute { // бакет от старой минуты — переиспользуем
		h.counts[idx] = 0
		h.stamps[idx] = minute
	}
	h.counts[idx]++
}

func (s *Solution) Count(host string) int {
	nowMin := s.now().Unix() / 60

	s.mu.RLock()
	defer s.mu.RUnlock()

	h, ok := s.hosts[host]
	if !ok {
		return 0
	}

	total := 0
	for i := range h.counts {
		if nowMin-h.stamps[i] < windowMin { // бакет ещё в окне
			total += h.counts[i]
		}
	}

	return total
}

func main() {
	s := NewSolution()
	now := time.Now()

	for i := 0; i < 5; i++ {
		s.OnPageDownloaded("facebook.com", now.Add(-time.Duration(i)*time.Minute))
	}
	s.OnPageDownloaded("facebook.com", now.Add(-15*time.Minute)) // вне окна
	s.OnPageDownloaded("twitter.com", now)

	fmt.Println("facebook.com", s.Count("facebook.com")) // 5
	fmt.Println("twitter.com", s.Count("twitter.com"))   // 1
	fmt.Println("linkedin.com", s.Count("linkedin.com")) // 0
}
