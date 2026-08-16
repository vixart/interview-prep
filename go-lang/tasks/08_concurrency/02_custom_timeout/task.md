# 4.6.2. Кастомный таймаут: контроль медленной зависимости

## Условие

> Как сделать влияние функции `getDiscount()` на всю программу более контролируемым?

```go
package main

import (
    "fmt"
    "time"
)

// Эта функция лезет по сети в старый монолит и может тупить.
func getDiscount() float64 {
    time.Sleep(2 * time.Second)
    return 12.0
}

func main() {
    fmt.Printf("Ваша скидка: %v", getDiscount())
}
```
