# Решение

## Проблемы

1. **Мьютекс объявлен внутри функции** (`var mu sync.Mutex`) — у каждого вызова
   свой мьютекс, никакой синхронизации нет. Мьютекс должен лежать рядом с данными
   (в структуре кеша) или быть пакетным.
2. **`mu.Unlock()` на незалоченном мьютексе** в ветке `ok == true` →
   `fatal error: sync: unlock of unlocked mutex` (паника).
3. **Гонка данных на глобальной мапе** `cache` — при конкурентных вызовах
   `fatal error: concurrent map read and map write`.
4. **Не атомарная связка «проверили → посчитали → записали»**: одно и то же дорогое
   вычисление выполнится многими горутинами (thundering herd) —
   лечится `singleflight` или хранением «обещания» результата.
5. **`time.Duration(secondsToSleep)`** — грубая ошибка: `float64` секунд приводится
   к `Duration`, то есть трактуется как **наносекунды**; нужно
   `time.Duration(secondsToSleep * float64(time.Second))`.
6. **Глобальное изменяемое состояние** — кеш стоит завернуть в тип с методами,
   передавать зависимостью, а не держать пакетной переменной.
7. Кеш растёт неограниченно — нет TTL/вытеснения.
8. Нет `context` для отмены долгого вычисления.

## Правильный вариант

```go
type Cache struct {
    mu   sync.RWMutex
    data map[int]int
}

func (c *Cache) Get(n int) int {
    c.mu.RLock()
    v, ok := c.data[n]
    c.mu.RUnlock()

    if ok {
        return v
    }

    value := LongCalculation(n)

    c.mu.Lock()
    defer c.mu.Unlock()

    if v, ok := c.data[n]; ok { // повторная проверка
        return v
    }
    c.data[n] = value

    return value
}
```

Полный исправленный вариант: [`fixed.go`](fixed.go) (запускается: `go run fixed.go`).
