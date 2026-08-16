# Решение

## Разбор

1. Сделать инкремент счётчика потокобезопасным — `atomic.AddInt64(&counter, 1)`
   (или мьютекс); объяснить, почему `counter++` — это не атомарная операция.
2. Добавить таймаут: обернуть работу в горутину и `select` по `ctx.Done()` /
   `time.After`, использовать `context.WithTimeout`.
3. Замерить и напечатать длительность (`start := time.Now(); ... time.Since(start)`).
4. Обсудить, что горутина после таймаута продолжает работать — как её корректно
   отменить (проброс `ctx` внутрь) и почему буферизованный канал спасает от утечки.

## Эталонный каркас

```go
func SimulateRequest(ctx context.Context) (int64, error) {
    done := make(chan int64, 1)

    go func() {
        time.Sleep(time.Duration(rand.Int63n(5)) * time.Second)
        done <- atomic.AddInt64(&counter, 1)
    }()

    select {
    case v := <-done:
        return v, nil
    case <-ctx.Done():
        return 0, ctx.Err()
    }
}
```
