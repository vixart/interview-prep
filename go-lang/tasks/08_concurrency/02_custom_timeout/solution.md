# Решение

## Разбор

1. Объяснить проблему: синхронный вызов медленной внешней зависимости блокирует
   всю программу на неопределённое время.
2. Обернуть вызов в горутину + канал + `select` с таймаутом,
   пробросить `context.Context` с дедлайном:

```go
func getDiscountCtx(ctx context.Context) (float64, error) {
    ch := make(chan float64, 1) // буфер обязателен, иначе горутина утечёт

    go func() { ch <- getDiscount() }()

    select {
    case v := <-ch:
        return v, nil
    case <-ctx.Done():
        return 0, ctx.Err()
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()

    d, err := getDiscountCtx(ctx)
    if err != nil {
        d = 0 // fallback: скидки нет
    }
    fmt.Printf("Ваша скидка: %v", d)
}
```

3. Обсудить:
   - почему канал должен быть буферизованным (иначе брошенная горутина зависнет навсегда);
   - что таймаут не отменяет саму работу — отмену нужно пробрасывать внутрь
     (например, `http.NewRequestWithContext`);
   - разумный fallback (degrade gracefully), ретраи, circuit breaker.
