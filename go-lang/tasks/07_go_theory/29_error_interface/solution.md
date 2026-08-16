# Решение

## Разбор

1. Реализовать `handle()` так, чтобы она возвращала ошибку **без** импорта
   `errors`/`fmt`.
2. Объяснить, что `error` — это обычный интерфейс:

```go
type error interface {
    Error() string
}
```

   значит, достаточно объявить свой тип с методом `Error() string`:

```go
type myError struct{}

func (e myError) Error() string { return "что-то пошло не так" }

func handle() error {
    return myError{}
}
```

3. Обсудить, почему `println(handle())` печатает адрес/структуру интерфейса,
   а не текст ошибки, и чем `println` отличается от `fmt.Println`.
4. Затронуть классическую ловушку: возврат `nil` указателя конкретного типа
   в интерфейсе `error` даёт `err != nil` (см. `4_5_3` — `4_5_5`).

Рабочий код: [`solution.go`](solution.go) (запускается: `go run solution.go`).
