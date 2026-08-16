# Решение

## Разбор

1. Реализовать базовый вариант: `map[string]string` + `sync.RWMutex`.
2. Объяснить, где он упирается при высокой нагрузке: единственный мьютекс —
   точка конкуренции, все горутины сериализуются на нём.
3. Реализовать **шардированный** кеш: `N` независимых сегментов,
   каждый со своей мапой и своим мьютексом; шард выбирается по хешу ключа
   (`fnv`/`maphash`), `N` — степень двойки, чтобы взять `hash & (N-1)`.

## Каркас решения

```go
const shardCount = 256

type shard struct {
    mu   sync.RWMutex
    data map[string]string
}

type ShardedCache struct {
    shards [shardCount]*shard
}

func (c *ShardedCache) getShard(k string) *shard {
    h := fnv.New32a()
    _, _ = h.Write([]byte(k))

    return c.shards[h.Sum32()&(shardCount-1)]
}

func (c *ShardedCache) Set(k, v string) {
    s := c.getShard(k)

    s.mu.Lock()
    defer s.mu.Unlock()

    s.data[k] = v
}

func (c *ShardedCache) Get(k string) (string, bool) {
    s := c.getShard(k)

    s.mu.RLock()
    defer s.mu.RUnlock()

    v, ok := s.data[k]
    return v, ok
}
```

## Разбор и подсказки

- Выбор числа шардов (обычно 32–512), false sharing и padding кеш-линий.
- Сравнение с `sync.Map` и с одной мапой под `RWMutex`; когда что выигрывает.
- Расширения: TTL, вытеснение (LRU), метрики, ограничение памяти.
- Как это тестировать и бенчмаркать (`go test -bench`, `-race`, `b.RunParallel`).
