# 4.1.7. Выравнивание типов: как Go раскладывает поля структуры в памяти

Раздел: `07_type_alignment`

Тип задачи: вопрос с собеседования «что выведет код и почему».

## Условие

Дан код. Нужно сказать, что выведут `unsafe.Sizeof` для двух структур
с одинаковым набором полей, но разным порядком.

```go
package main

import (
    "fmt"
    "unsafe"
)

type Foo struct {
    aaa bool
    bbb int32
    ccc bool
}

type Bar struct {
    aaa bool
    ccc bool
    bbb int32
}

func main() {
    ff := Foo{}
    bb := Bar{}
    fmt.Println(unsafe.Sizeof(ff))
    fmt.Println(unsafe.Sizeof(bb))
}
```

**Задание:** определи, что выведет программа (и скомпилируется ли она), и объясни почему.
