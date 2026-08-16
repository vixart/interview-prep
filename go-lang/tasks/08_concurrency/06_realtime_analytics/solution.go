// Аналитика в реальном времени: события за последние 5 минут на пользователя.
//
// Вместо хранения всех timestamp'ов (растущая память) — кольцевой буфер
// из 300 секундных бакетов на пользователя + текущая сумма:
//   - HandleEvent: O(1) — инкремент бакета текущей секунды;
//   - GetCount:    O(1) — сумма поддерживается инкрементально, протухшие
//     бакеты обнуляются лениво при обращениях.
//
// Память на пользователя константная (300 int).
package main

import (
	"fmt"
	"sync"
	"time"
)

const windowSec = 5 * 60

type userBuckets struct {
	buckets [windowSec]int
	stamps  [windowSec]int64 // какой секунде принадлежит бакет
	total   int
}

// advance обнуляет бакеты, вышедшие из окна [now-299, now].
func (u *userBuckets) advance(now int64) {
	for i := range u.buckets {
		if u.stamps[i] != 0 && now-u.stamps[i] >= windowSec {
			u.total -= u.buckets[i]
			u.buckets[i] = 0
			u.stamps[i] = 0
		}
	}
}

type Service struct {
	mu    sync.Mutex
	users map[string]*userBuckets
}

func NewService() *Service {
	return &Service{users: make(map[string]*userBuckets)}
}

func (s *Service) HandleEvent(userName string, currentTime time.Time) {
	now := currentTime.Unix()

	s.mu.Lock()
	defer s.mu.Unlock()

	u, ok := s.users[userName]
	if !ok {
		u = &userBuckets{}
		s.users[userName] = u
	}

	idx := now % windowSec
	if u.stamps[idx] != now {
		// бакет от старой секунды — вытесняем его из суммы
		u.total -= u.buckets[idx]
		u.buckets[idx] = 0
		u.stamps[idx] = now
	}
	u.buckets[idx]++
	u.total++
}

func (s *Service) GetCount(userName string, currentTime time.Time) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	u, ok := s.users[userName]
	if !ok {
		return 0
	}

	u.advance(currentTime.Unix())
	return u.total
}

func main() {
	svc := NewService()
	base := time.Now()

	svc.HandleEvent("alice", base.Add(-6*time.Minute)) // протухнет
	svc.HandleEvent("alice", base.Add(-4*time.Minute))
	svc.HandleEvent("alice", base.Add(-1*time.Minute))
	svc.HandleEvent("alice", base)
	svc.HandleEvent("bob", base)

	fmt.Println(svc.GetCount("alice", base))  // 3
	fmt.Println(svc.GetCount("bob", base))    // 1
	fmt.Println(svc.GetCount("nobody", base)) // 0
}
