# 4.8.4. Ревью кода: отмена HTTP-запросов и утечки ресурсов

Раздел: `code_review / 04_http_cancel`

Тип задачи: ревью кода — сделать ревью и добавить отмену запросов при ошибке.

## ТЗ

Программа отправляет HTTP GET запросы к нескольким URL параллельно.
Если хотя бы один запрос завершается с ошибкой, нужно отменить все остальные.

### Требования

- Запросы должны выполняться параллельно.
- При первой ошибке все остальные запросы должны быть отменены.

> Этот код НАМЕРЕННО содержит ошибки для учебных целей! Не запускайте в production!

## Исходный код

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"
)

func main() {
    urls := []string{
        // ...
        "https://github.com",
        "https://stackoverflow.com",
    }

    for _, url := range urls {
        go func() {
            err := fetch(context.Background(), url)
            if err != nil {
                fmt.Printf("Error fetching %s: %v\n", url, err)
                return
            }
            fmt.Printf("Success: %s\n", url)
        }()
    }

    fmt.Println("All requests launched!")
    time.Sleep(400 * time.Millisecond)
    fmt.Println("Done")
}

func fetch(ctx context.Context, url string) error {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    // ...
}
```

**Задание:** сделай ревью: найди проблемы, объясни их и предложи исправленный вариант.
