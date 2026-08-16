// Парковка с ограниченным числом мест — семафор на буферизованном канале.
//
// Главный баг оригинала: поле slots вообще не использовалось, все машины
// парковались одновременно. Семафор ёмкостью N гарантирует, что внутри
// секции Park не больше N горутин; defer гарантирует освобождение места.
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type ParkingLot struct {
	slots  chan struct{}
	inside atomic.Int64 // для самопроверки в демо
}

func NewParkingLot(n int) *ParkingLot {
	return &ParkingLot{slots: make(chan struct{}, n)}
}

func (p *ParkingLot) Park(carID int64) {
	p.slots <- struct{}{}        // занять место (блокируется, если всё занято)
	defer func() { <-p.slots }() // освободить место при любом исходе

	cur := p.inside.Add(1)
	fmt.Printf("Тачка %d припарковывается... (занято мест: %d)\n", carID, cur)
	time.Sleep(200 * time.Millisecond) // имитация времени парковки
	p.inside.Add(-1)
	fmt.Printf("Тачка %d припаркована.\n", carID)
}

func main() {
	parking := NewParkingLot(3)

	var wg sync.WaitGroup
	for _, carID := range []int64{1, 2, 3, 4, 5, 6} {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			parking.Park(id)
		}(carID)
	}

	wg.Wait()
	fmt.Println("Все тачки припаркованы.")
}
