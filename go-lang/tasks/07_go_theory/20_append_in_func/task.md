# 4.3.8. Обновление слайса через функцию: что изменится

Раздел: `08_append_mistake`

Тип задачи: вопрос с собеседования «что выведет программа и почему».

## Условие

> Задача: Что выведет программа и почему?
> Как правильно исправить `appendLenWrong`?

```go
package main

import (
    "fmt"
)

func appendLenWrong(numbers []*int) {
    size := len(numbers)
    numbers = append(numbers, &size)
}

func main() {
    numbers := make([]*int, 0, 5)
    var number int
    for range 3 {
        number++
        numbers = append(numbers, &number)
    }

    appendLenWrong(numbers)

    for _, number := range numbers {
        fmt.Printf("%d ", *number)
    }
}
```

**Задание:** определи, что выведет программа (и скомпилируется ли она), и объясни почему.
