# Решение

## Разбор

1. Найти место паники: `cache.Load("width").(float64)` — ключа нет,
   из мапы возвращается нулевое значение `interface{}` (то есть `nil`),
   и «жёсткая» type assertion паникует:
   `interface conversion: interface {} is nil, not float64`.
2. Объяснить, почему проверка `height == nil` работает, а `width` — падает,
   и в чём разница между отсутствующим ключом и ключом с `nil` значением
   (нужна форма `v, ok := c.data[key]`).
3. Исправить: всегда использовать безопасную форму assertion

```go
width, ok := cache.Load("width").(float64)
if !ok {
    fmt.Println("Width not found or wrong type")
    return
}
```

   либо сделать `Load` возвращающим `(interface{}, bool)`, либо использовать дженерики.
4. Обсудить ловушку «nil != nil» в интерфейсах: интерфейс равен `nil` только когда
   и тип, и значение равны `nil`.
