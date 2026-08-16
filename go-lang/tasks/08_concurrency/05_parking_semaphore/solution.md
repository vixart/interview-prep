# Решение

## Разбор

1. Заметить главную проблему: поле `slots` вообще не используется —
   на парковку одновременно заезжают все 6 машин, ограничение не соблюдается.
2. Реализовать ограничение через семафор на буферизованном канале
   (или `golang.org/x/sync/semaphore`):

```go
type ParkingLot struct {
    slots chan struct{}
}

func NewParkingLot(n int) *ParkingLot {
    return &ParkingLot{slots: make(chan struct{}, n)}
}

func (p *ParkingLot) Park(carID int64) {
    p.slots <- struct{}{}        // занять место
    defer func() { <-p.slots }() // освободить место

    fmt.Printf("Тачка %d припарковывается...\n", carID)
    time.Sleep(time.Second)
    fmt.Printf("Тачка %d припаркована.\n", carID)
}
```

3. Обсудить сопутствующее: почему нужен `defer` на освобождение слота,
   как добавить `context` для отказа от ожидания, чем плох вариант со счётчиком
   и busy-wait, и почему `sync.Mutex` тут не подходит (нужен доступ N, а не 1).
