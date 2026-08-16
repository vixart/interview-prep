# 4.8.6. Ревью кода: race condition в кеше

Раздел: `code_review / 06_cache_race`

Тип задачи: ревью кода — сделайте ревью кода и исправьте проблемы.

## ТЗ (то, что код должен делать)

Простой in-memory кеш для хранения результатов дорогих вычислений.
Программа кеширует результаты и переиспользует их при повторных запросах.

> Этот код НАМЕРЕННО содержит ошибки для учебных целей! Не запускайте в production!

## Исходный код

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Cache struct {
    data map[string]string
}

func NewCache() *Cache {
    return &Cache{
        data: make(map[string]string),
    }
}

func (c *Cache) Get(key string) (string, bool) {
    value, ok := c.data[key]
    return value, ok
}

func (c *Cache) Set(key, value string) {
    c.data[key] = value
}

func (c *Cache) Delete(key string) {
    delete(c.data, key)
}

// expensiveComputation симулирует дорогое вычисление
func expensiveComputation(key string) string {
    time.Sleep(100 * time.Millisecond)
    return fmt.Sprintf("result for %s", key)
}

// GetOrCompute получает значение из кеша или вычисляет его
func GetOrCompute(cache *Cache, key string) string {
    // Проверяем кеш
    if value, ok := cache.Get(key); ok {
        return value
    }

    // Вычисляем значение
    value := expensiveComputation(key)

    // Сохраняем в кеш
    cache.Set(key, value)

    return value
}

func main() {
    cache := NewCache()
    var wg sync.WaitGroup

    // ... в цикле запускаются горутины, которые конкурентно
    //     вызывают GetOrCompute(cache, key) для набора ключей
    //     и печатают результат; в конце wg.Wait()
}
```

**Задание:** сделай ревью: найди проблемы, объясни их и предложи исправленный вариант.
