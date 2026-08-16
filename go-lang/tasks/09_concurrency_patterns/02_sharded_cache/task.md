# 4.7.2. Sharded Cache — сегментированный in-memory кеш

## Условие

> Напишите реализацию InMemory кэша.

```go
package main

type Cache interface {
    Set(k string, v string)
    Get(k string) (string, bool)
}
```

**Задание:** реализуй потокобезопасный кеш и подумай, как снять конкуренцию за один мьютекс при высокой нагрузке.
