# 4.4.2. Исправление бага: безопасная работа с мапой (общая мапа в канале)

Тип задачи: «что выведет программа и почему» + найти и исправить проблемы.

## Условие

> Что выведет программа и почему?
> Также проанализируй код и исправь все потенциальные проблемы, если таковые есть.

```go
package main

import (
    "fmt"
    "time"
)

func UpdateProductStock() <-chan map[string]int {
    stockUpdates := make(chan map[string]int)

    go func() {
        currentStock := map[string]int{
            "Apples":  50,
            "Bananas": 30,
            "Oranges": 20,
            "Grapes":  15,
        }

        for i := 0; i < 5; i++ {
            for product, quantity := range currentStock {
                currentStock[product] = int(float64(quantity) * 0.95)
            }

            stockUpdates <- currentStock

            time.Sleep(150 * time.Millisecond)
        }
    }()

    return stockUpdates
}

func main() {
    stockStream := UpdateProductStock()

    var stockHistory []map[string]int

    for i := 0; i < 5; i++ {
        stock := <-stockStream
        stockHistory = append(stockHistory, stock)
    }

    for i, stock := range stockHistory {
        fmt.Printf("Iteration %d: %v\n", i+1, stock)
    }
}
```

**Задание:** скажи, что выведет программа, найди проблемы и исправь их.
