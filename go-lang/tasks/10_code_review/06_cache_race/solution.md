# Решение

## Разбор

1. Найти проблемы конкурентности: `map` читается и пишется из нескольких горутин
   без синхронизации → `fatal error: concurrent map writes` / гонка данных.
2. Заметить проблему «thundering herd»: между `Get` и `Set` в `GetOrCompute` нет
   атомарности, поэтому одно и то же дорогое вычисление выполняется многими горутинами.
3. Исправить: добавить `sync.RWMutex` в `Cache` (`RLock`/`RUnlock` в `Get`,
   `Lock`/`Unlock` в `Set` и `Delete`), а для устранения дублирующих вычислений
   использовать `singleflight` или хранить в кеше «обещание» результата.

## Разбор

```go
type Cache struct {
    mu   sync.RWMutex
    data map[string]string
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    value, ok := c.data[key]
    return value, ok
}

func (c *Cache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.data[key] = value
}
```

Полный исправленный вариант: [`fixed.go`](fixed.go) (запускается: `go run fixed.go`).
