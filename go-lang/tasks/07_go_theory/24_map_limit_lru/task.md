# 4.4.3. Ограничение размера мапы (реализация LRU)

Тип задачи: дописать код + «что выведет программа и почему».

## Условие

Вы разрабатываете сервис для подсчёта уникальных слов, но количество слов может быть
очень большим, и вы решили ограничить максимальное число записей в мапе.
Если в мапу добавляется больше слов, чем указано лимитом, она должна автоматически
удалять самые старые записи.

```go
package main

import (
    "fmt"
)

type WordCounter struct {
    counts map[string]int
    limit  int
}

func NewWordCounter(limit int) *WordCounter {
    return &WordCounter{
        counts: make(map[string]int),
        limit:  limit,
    }
}

func (wc *WordCounter) CountWord(word string) {
    wc.counts[word]++

    if len(wc.counts) > wc.limit {
        // Логика удаления (здесь нужно реализовать)
    }
}

func main() {
    wc := NewWordCounter(3)

    words := []string{"apple", "banana", "apple", "orange", "grape", "banana", "kiwi"}
    for _, word := range words {
        wc.CountWord(word)
    }

    fmt.Println("Количество слов:", wc.counts)
}
```

**Задание:** реализуй ограничение размера с вытеснением самых старых записей.
