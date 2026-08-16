# 4.8.9. Ревью кода: кеш с проблемами конкурентности

Раздел: `code_review / 09_cache_concurrency`

Тип задачи: ревью — «отревьюй кусок кода и предложи доработки, если такие потребуются».

## Исходный код

```go
package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

func LongCalculation(n int) int {
    secondsToSleep := rand.Float64() * float64(n)
    time.Sleep(time.Duration(secondsToSleep))
    return n + 1
}

var cache = map[int]int{}

func CachedLongCalculation(n int) int {
    var mu sync.Mutex

    mu.Lock()
    found, ok := cache[n]
    mu.Unlock()

    if !ok {
        value := LongCalculation(n)
        mu.Lock()
        cache[n] = value
        mu.Unlock()
        return value
    }

    mu.Unlock() // Лишний unlock

    return found
}

func main() {
    nums := []int{5, 10, 22}
    for _, n := range nums {
        // ... параллельные вызовы CachedLongCalculation(n)
    }
}
```

**Задание:** сделай ревью: найди проблемы, объясни их и предложи исправленный вариант.
