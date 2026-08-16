# Решение

## Проблемы

1. **`time.Sleep` вместо синхронизации** — программа завершится, не дождавшись
   ответов. Нужен `sync.WaitGroup` — это первое исправление.
2. **Отмена не реализована**: каждый запрос получает `context.Background()`,
   поэтому ошибка одного запроса не останавливает остальные.
   Нужен общий `ctx, cancel := context.WithCancel(...)` и `cancel()` при первой ошибке
   (или `errgroup.WithContext`, который делает это из коробки).
3. **Захват переменной цикла** `url` в замыкании (актуально до Go 1.22) —
   передавать параметром в горутину.
4. **Нет `defer cancel()`** → утечка контекста.
5. В `fetch`: не закрывается `resp.Body`, не проверяется `resp.StatusCode`,
   не задан `Timeout` у клиента.
6. Ошибки только печатаются, никуда не агрегируются.

## Эталонный вариант

```go
g, ctx := errgroup.WithContext(context.Background())

for _, url := range urls {
    url := url
    g.Go(func() error {
        return fetch(ctx, url)
    })
}

if err := g.Wait(); err != nil {
    fmt.Println("failed:", err)
}
```
