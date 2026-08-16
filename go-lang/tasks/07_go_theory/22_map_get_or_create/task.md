# 4.4.1. Конкурентное обновление мапы: GetOrCreate

Тип задачи: написать код (конкурентная структура данных).

## Условие

Напишите функцию `GetOrCreate`, которая или создаёт новый элемент мапы,
если его ещё не было, и возвращает его значение, или просто возвращает значение при наличии.

Важно учесть, что код должен корректно работать в конкурентной среде.

## Каркас

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    cm := NewConcurrentMap()

    wg := sync.WaitGroup{}
    wg.Add(2)

    go func() {
        defer wg.Done()
        val := cm.GetOrCreate("key1", "value1")
        fmt.Println("Goroutine 1 got:", val)
    }()

    go func() {
        defer wg.Done()
        val := cm.GetOrCreate("key1", "value2")
        fmt.Println("Goroutine 2 got:", val)
    }()

    wg.Wait()
}
```

## Что нужно реализовать

- Тип `ConcurrentMap` и конструктор `NewConcurrentMap()`.
- Метод `GetOrCreate(key, value string) string` с корректной синхронизацией.
