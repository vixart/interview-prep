# 4.2. Группировка серверов по стабильности

Раздел: `04_hash_tables / 17_group_by_stability`

## Условие

Напишите функцию, которая по списку метрик серверов строит отображение
«стабильность → список ID серверов».

## Входные данные

- `stats` — список структур с полями `server` (ID сервера) и `stability` (стабильность, `float64`).

## Выходные данные

- `map[string][]int`: ключ — стабильность (строка), значение — список ID серверов.

## Сигнатура

```go
type ServerStat struct {
    Server    int
    Stability float64
}

func groupByStability(stats []ServerStat) map[string][]int { return nil }
```

## Пример

```
stats = [
    {server: 1, stability: 99},
    {server: 2, stability: 97},
    {server: 3, stability: 34},
    {server: 4, stability: 97},
    {server: 5, stability: 97.1},
]

Результат: {"34": [3], "97": [2, 4], "99": [1], "97.1": [5]}
```

Замечание: ключ — строковое представление `float64` без лишних нулей
(`strconv.FormatFloat(v, 'f', -1, 64)`), поэтому `97` и `97.1` — разные ключи.
