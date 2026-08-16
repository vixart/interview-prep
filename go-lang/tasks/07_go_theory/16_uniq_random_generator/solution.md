# Решение

## Разбор

1. Реализовать функцию: генерировать случайные числа через `math/rand`
   (`rand.Int()` / `rand.Intn(...)`) и складывать их в `map[int]struct{}`
   до тех пор, пока в множестве не наберётся `n` уникальных значений,
   затем собрать результат в слайс.
2. Не забыть заранее выделить память: `make(map[int]struct{}, n)` и
   `make([]int, 0, n)`.
3. Обсудить, почему `map[int]struct{}` лучше, чем `map[int]bool` (нулевой размер значения),
   и что порядок обхода map в Go случайный.

## Каркас решения

```go
func uniqN(n int) []int {
    m := make(map[int]struct{}, n)

    for len(m) < n {
        m[rand.Int()] = struct{}{}
    }

    result := make([]int, 0, n)
    for k := range m {
        result = append(result, k)
    }

    return result
}
```

Рабочий код: [`solution.go`](solution.go) (запускается: `go run solution.go`).
