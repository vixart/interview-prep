// Защита общей map мьютексом: RWMutex дает много одновременных читателей (RLock)
// и эксклюзивную запись (Lock), разблокировка всегда через defer.
// Тот же сценарий на каналах — в channel_vs_mutex/channel.
package main

import (
	"fmt"
	"sync"
)

type MutexScoreboardManager struct {
	l          sync.RWMutex
	scoreboard map[string]int
}

func NewMutexScoreboardManager() *MutexScoreboardManager {
	return &MutexScoreboardManager{
		// возвращаем УКАЗАТЕЛЬ: мьютекс нельзя копировать
		scoreboard: map[string]int{},
	}
}

func (msm *MutexScoreboardManager) Update(name string, val int) {
	msm.l.Lock()
	// эксклюзивная блокировка на запись — не пустит ни писателей, ни читателей
	defer msm.l.Unlock()
	// разблокировка через defer — обязательная привычка
	msm.scoreboard[name] = val
}

func (msm *MutexScoreboardManager) Read(name string) (int, bool) {
	msm.l.RLock()
	// читателей может быть много одновременно
	defer msm.l.RUnlock()
	val, ok := msm.scoreboard[name]
	return val, ok
}

func main() {
	msm := NewMutexScoreboardManager()
	teams := []string{"Lions", "Tigers", "Bears"}
	var wg sync.WaitGroup
	wg.Add(len(teams))
	for _, v := range teams {
		go func(team string) {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				curScore, found := msm.Read(team)
				if !found {
					curScore = 10
				} else {
					curScore += len(team)
				}
				msm.Update(team, curScore)
			}
		}(v)
	}
	wg.Wait()
	for _, v := range teams {
		score, found := msm.Read(v)
		fmt.Println(v, score, found)
	}
}
