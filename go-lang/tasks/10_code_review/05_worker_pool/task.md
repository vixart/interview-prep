# 4.8.5. Ревью кода: worker pool с ошибками синхронизации

Раздел: `code_review / 05_worker_pool`

Тип задачи: ревью кода — сделайте ревью кода и исправьте проблемы.

## ТЗ

Пул воркеров для обработки задач (например, обработка изображений).
Программа получает задачи из очереди и обрабатывает их параллельно.

> Этот код НАМЕРЕННО содержит ошибки для учебных целей! Не запускайте в production!

## Исходный код

```go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

type Task struct {
    ID   int
    Data string
}

type Result struct {
    TaskID int
    Output string
    Error  error
}

func main() {
    tasks := []Task{
        {ID: 1, Data: "image1.jpg"},
        {ID: 2, Data: "image2.jpg"},
        {ID: 3, Data: "corrupt"},
        {ID: 4, Data: "image4.jpg"},
        {ID: 5, Data: "image5.jpg"},
    }

    // Создаем каналы
    jobChan := make(chan Task)
    resultChan := make(chan Result)

    // Запускаем воркеры для каждой задачи
    for range tasks {
        go worker(jobChan, resultChan)
    }

    // Отправляем задачи
    go func() {
        for _, task := range tasks {
            jobChan <- task
        }
    }()

    // Собираем результаты
    for i := 0; i < len(tasks); i++ {
        result := <-resultChan
        if result.Error != nil {
            fmt.Printf("Task %d failed: %v\n", result.TaskID, result.Error)
            continue
        }
        // ...
    }
}

func worker(jobs <-chan Task, results chan<- Result) {
    // обрабатывает задачи, эмулируя работу через time.Sleep/rand
}
```

**Задание:** сделай ревью: найди проблемы, объясни их и предложи исправленный вариант.
