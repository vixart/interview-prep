package main

import (
	"fmt"
	"sync"
)

// MutexScoreboardManager хранит:
// - RWMutex для синхронизации
// - общую map, к которой обращаются все горутины
type MutexScoreboardManager struct {
	l          sync.RWMutex   // защита доступа к map
	scoreboard map[string]int // общее состояние (shared memory)
}

// Конструктор — создаёт пустую map.
func NewMutexScoreboardManager() *MutexScoreboardManager {
	return &MutexScoreboardManager{
		scoreboard: map[string]int{},
	}
}

// Update — изменяет состояние.
// Используется обычный Lock (эксклюзивный),
// потому что мы модифицируем map.
func (msm *MutexScoreboardManager) Update(name string, val int) {
	msm.l.Lock()         // блокируем всех (и читателей тоже)
	defer msm.l.Unlock() // гарантированное освобождение
	msm.scoreboard[name] = val
}

// Read — только чтение.
// Используется RLock — параллельные чтения разрешены.
func (msm *MutexScoreboardManager) Read(name string) (int, bool) {
	msm.l.RLock()         // несколько читателей могут работать одновременно
	defer msm.l.RUnlock() // освобождаем read-lock
	val, ok := msm.scoreboard[name]
	return val, ok
}

func main() {

	msm := NewMutexScoreboardManager()

	teams := []string{"Lions", "Tigers", "Bears"}

	var wg sync.WaitGroup
	wg.Add(len(teams))

	// Запускаем по одной горутине на команду.
	for _, v := range teams {
		go func(team string) {
			defer wg.Done()

			for i := 0; i < 10; i++ {

				// Читаем текущее значение
				curScore, found := msm.Read(team)

				if !found {
					curScore = 10
				} else {
					curScore += len(team)
				}

				// Обновляем значение
				msm.Update(team, curScore)
			}

		}(v)
	}

	wg.Wait()

	// Вывод финальных результатов
	for _, v := range teams {
		score, found := msm.Read(v)
		fmt.Println(v, score, found)
	}
}
