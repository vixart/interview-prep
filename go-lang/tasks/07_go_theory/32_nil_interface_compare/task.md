# 4.5.5. Nil-интерфейсы: почему сравнение с nil не работает

Тип задачи: вопрос с собеседования «что выведет программа и почему».

## Условие

```go
package main

import (
    "fmt"
)

type SomeStruct struct {
    Value int
}

func CheckForNil(i interface{}) {
    if i == nil {
        fmt.Println("Это nil!")
        return
    }

    fmt.Println("Это не nil!")
}

func main() {
    var s *SomeStruct
    CheckForNil(s)
}
```

**Задание:** определи, что выведет программа (и скомпилируется ли она), и объясни почему.
