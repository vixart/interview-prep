# Решение

## Разбор

Реализовать worker pool:

```go
func processImages(images []Image, workers int) []Result {
    jobs := make(chan Image)
    results := make(chan Result, len(images))

    var wg sync.WaitGroup
    wg.Add(workers)

    for i := 0; i < workers; i++ {
        go func() {
            defer wg.Done()

            for img := range jobs {
                results <- process(img)
            }
        }()
    }

    for _, img := range images {
        jobs <- img
    }
    close(jobs)

    wg.Wait()
    close(results)

    // ...собрать результаты
}
```

## Разбор и подсказки

- Фиксированное число воркеров, а не «горутина на задачу».
- Канал задач закрывается **после** отправки всех задач; воркеры выходят по `for range`.
- `results` буферизуем (или читаем параллельно), иначе воркеры заблокируются.
- Сбор ошибок и отмена через `context`; альтернатива — `errgroup.SetLimit`.
- Обсуждается разница между worker pool и семафором:
  пул переиспользует горутины, семафор ограничивает одновременность.
- Сколько воркеров брать: CPU-bound → `runtime.NumCPU()`, IO-bound → больше.
