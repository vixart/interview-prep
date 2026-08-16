# Решение

## Разбор

Реализовать паттерн pipeline: каждый этап — отдельная горутина (или группа горутин),
этапы соединены каналами; каждый этап читает из входного канала и пишет в выходной,
закрывая его по завершении.

## Каркас решения

```go
func generator(ctx context.Context, txs []Transaction) <-chan Transaction {
    out := make(chan Transaction)

    go func() {
        defer close(out)

        for _, tx := range txs {
            select {
            case out <- tx:
            case <-ctx.Done():
                return
            }
        }
    }()

    return out
}

func filter(ctx context.Context, in <-chan Transaction) <-chan Transaction {
    out := make(chan Transaction)

    go func() {
        defer close(out)

        for tx := range in {
            if tx.Amount <= 0 {
                continue
            }
            select {
            case out <- tx:
            case <-ctx.Done():
                return
            }
        }
    }()

    return out
}

func convert(ctx context.Context, in <-chan Transaction) <-chan Transaction { /* ... */ }

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    var result []Transaction
    for tx := range convert(ctx, filter(ctx, generator(ctx, txs))) {
        result = append(result, tx)
    }
}
```

## Разбор и подсказки

- Правило: канал закрывает тот, кто в него пишет.
- Каждый этап должен уметь досрочно выходить по `ctx.Done()`, иначе при раннем
  выходе потребителя горутины останутся навсегда (утечка).
- Узкое место масштабируется fan-out'ом: несколько горутин на тяжёлом этапе + fan-in.
- Буферизация каналов сглаживает скачки, но не решает проблему разной скорости этапов.
