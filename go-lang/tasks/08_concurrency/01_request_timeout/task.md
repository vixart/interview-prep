# 4.6.1. Таймаут для запроса и конкурентная среда

## Условие

> Требуется доработать метод `SimulateRequest` так, чтобы:
> - Код работал в конкурентной среде
> - При долгом ожидании запрос отваливался по таймауту
> - На консоль печаталось время выполнения запроса

```go
package main

import (
    "log"
    "math/rand"
    "time"
)

var counter int64

func SimulateRequest() int64 {
    time.Sleep(time.Duration(rand.Int63n(5)) * time.Second)
    counter++
    return counter
}

func main() {
    val := SimulateRequest()
    log.Printf("Значение счетчика: %d\n", val)
}
```
