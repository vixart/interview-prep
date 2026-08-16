# 4.6.3. Исправление бага с каналами: deadlock и утечки

Тип задачи: найти баг, объяснить и исправить.

## Условие

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch := make(chan bool)
    go func() {
        time.Sleep(3 * time.Second)
        fmt.Println("Отдельная горутина отвисла")
        ch <- false
    }()

    ticker := time.NewTicker(time.Second)
    for {
        select {
        case <-ticker.C:
            fmt.Println("Произошёл тик тикера")
            ch <- true
        case value := <-ch:
            fmt.Printf("Получено значение %t\n", value)
            return
        }
    }
}
```

**Задание:** найди баг, объясни поведение программы и исправь код.
