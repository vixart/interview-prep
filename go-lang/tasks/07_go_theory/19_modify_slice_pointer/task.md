# 4.3.7. Модификация слайса через указатель на элемент

Раздел: `07_modify_slice`

Тип задачи: вопрос с собеседования «что выведет программа и почему».

## Условие

```go
package main

import "fmt"

type person struct {
    age int
}

func main() {
    people := make([]person, 2)

    p1 := &people[1]
    fmt.Printf("%p", p1)

    p1.age++

    people = append(people, person{}, person{}, person{})
    fmt.Println(cap(people))

    p1.age++

    fmt.Println(people[1].age)
}
```

**Задание:** определи, что выведет программа (и скомпилируется ли она), и объясни почему.
