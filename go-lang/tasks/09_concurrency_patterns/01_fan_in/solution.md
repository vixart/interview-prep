# Решение

## Эталонное решение

```go
func merge(channels ...chan int64) <-chan int64 {
    out := make(chan int64)

    var wg sync.WaitGroup
    wg.Add(len(channels))

    for _, ch := range channels {
        go func(ch <-chan int64) {
            defer wg.Done()

            for v := range ch {
                out <- v
            }
        }(ch)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

## Разбор и подсказки

- Закрывать `out` можно только один раз и только после того, как все читатели
  входных каналов завершились → отдельная горутина с `wg.Wait()`.
- Порядок значений не гарантирован.
- Продвинутый вариант: добавить `ctx` для отмены и `select` при отправке в `out`,
  чтобы горутины не залипли, если потребитель перестал читать.
- Обсуждается разница fan-in / fan-out и вариант через рефлексию (`reflect.Select`),
  когда каналы разнотипные.
