# Решение

## Разбор

1. **Deadlock**: бесконечный `for { select { case <-a } }` в `main` никогда
   не завершится, а `wc.Wait()` после него — недостижимый код.
   Когда все 5 значений прочитаны, `main` блокируется навсегда:
   `fatal error: all goroutines are asleep - deadlock!`
2. **Лишний мьютекс** `fff`: канал сам по себе потокобезопасен, блокировка
   вокруг отправки только сериализует горутины и может привести к тому,
   что горутина держит лок, ожидая места в канале.
3. **`select` с одним `case`** эквивалентен обычному чтению — избыточная конструкция.
4. **Канал не закрывается**, поэтому нельзя использовать `for range`.
5. **`fmt.Printf(result)`** — строка пользователя как формат-строка
   (`go vet` ругается); нужно `fmt.Print(result)` или `fmt.Printf("%s", result)`.
6. **Именование**: `wc` вместо `wg`, `fff`, `a` — неинформативно.
7. Передача `wc` параметром и одновременно захват `fff` замыканием — непоследовательно.
8. `strconv.Itoa(i)` внутри `Sprintf("%s")` — достаточно `Sprintf("%d", i)`.

## Правильный вариант

```go
func main() {
    var wg sync.WaitGroup
    results := make(chan string, 5)

    for i := 0; i < 5; i++ {
        wg.Add(1)

        go func(i int) {
            defer wg.Done()
            results <- fmt.Sprintf("Current gorutine number: %d\n", i)
        }(i)
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    for result := range results {
        fmt.Print(result)
    }
}
```

Полный исправленный вариант: [`fixed.go`](fixed.go) (запускается: `go run fixed.go`).
