# Решение

## Разбор

Реализовать семафор на буферизованном канале:

```go
type Semaphore chan struct{}

func NewSemaphore(n int) Semaphore { return make(Semaphore, n) }

func (s Semaphore) Acquire() { s <- struct{}{} }
func (s Semaphore) Release() { <-s }

// с поддержкой контекста
func (s Semaphore) AcquireCtx(ctx context.Context) error {
    select {
    case s <- struct{}{}:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

Использование:

```go
sem := NewSemaphore(limit)

var wg sync.WaitGroup
for _, url := range urls {
    wg.Add(1)

    go func(url string) {
        defer wg.Done()

        sem.Acquire()
        defer sem.Release()

        download(url)
    }(url)
}
wg.Wait()
```

## Разбор и подсказки

- Отличие семафора от worker pool: горутина на задачу + ограничение
  одновременности vs фиксированный пул воркеров.
- Обязательный `defer Release()` — иначе слот утекает при панике/ошибке.
- Динамически изменяемый лимит: `golang.org/x/sync/semaphore` (weighted semaphore)
  либо пересоздание канала / счётчик с `sync.Cond`.
- `errgroup.Group` + `SetLimit(n)` как готовое решение.
