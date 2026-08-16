# 4.3.3. Append и изменение слайса — вариант B

Тип задачи: вопрос с собеседования «что выведет код и почему».

## Условие

```go
package main

import "fmt"

func modifyElement(slice []int) {
    slice[1] = 999
}

func addElement(slice []int) {
    slice = append(slice, 100)
    slice[0] = 888
    fmt.Println("Внутри addElement:", slice)
}

func main() {
    original := []int{10, 20, 30}

    fmt.Println("До modifyElement:", original) // [10 20 30]
    modifyElement(original)
    fmt.Println("После modifyElement:", original)

    fmt.Println("До addElement:", original)
    addElement(original)
    fmt.Println("После addElement:", original)
}
```

**Задание:** определи, что выведет программа (и скомпилируется ли она), и объясни почему.
