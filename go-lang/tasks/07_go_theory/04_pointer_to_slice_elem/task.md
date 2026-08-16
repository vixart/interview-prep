# 4.1.4. Указатели на элементы слайса и append

Тип задачи: вопрос с собеседования «что выведет код и почему».

## Условие

Дан код. Нужно сказать, что выведут оба `Println`, и объяснить почему значения расходятся.

```go
package main

import "fmt"

type car struct {
    color   string
    mileage int
}

func main() {
    cars := []car{
        {
            color:   "red",
            mileage: 5000,
        },
        {
            color:   "green",
            mileage: 10000,
        },
        {
            color:   "blue",
            mileage: 7000,
        },
    }

    carPtr := &cars[0]
    carPtr.mileage += 100

    cars = append(cars, car{color: "yellow", mileage: 15000})
    carPtr.mileage += 50

    fmt.Println(cars[0].mileage, cars[0].color) // Выводим пробег первого элемента слайса
    fmt.Println(carPtr.mileage, carPtr.color)   // Выводим пробег через указатель
}
```

**Задание:** определи, что выведет программа (и скомпилируется ли она), и объясни почему.
