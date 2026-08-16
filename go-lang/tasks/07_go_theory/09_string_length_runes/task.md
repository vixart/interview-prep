# 4.2.2. Длина строки: руны, байты и UTF-8

Тип задачи: вопрос с собеседования «что выведет код и почему».

## Условие

Дан код. Нужно сказать, что он выведет, и объяснить разницу между `len()` и
`utf8.RuneCountInString()` для строки с не-ASCII символами.

```go
package main

import (
    "fmt"
    "unicode/utf8"
)

func main() {
    str := "ddЯй漢"
    fmt.Println("Длина через len:", len(str))
    fmt.Println("Длина через RuneCountInString:", utf8.RuneCountInString(str))
}
```

**Задание:** определи, что выведет программа (и скомпилируется ли она), и объясни почему.
