// Исправленный worker pool.
//
// Что изменилось относительно исходника:
//   - воркеров фиксированное число (параметр), а не «по горутине на задачу»;
//   - jobs закрывается после отправки всех задач — воркеры корректно
//     выходят из for range (в оригинале висели навсегда — утечка);
//   - sync.WaitGroup + горутина-закрыватель results: сбор результатов не
//     завязан на «магическое» знание их количества;
//   - паника в воркере не роняет программу (recover -> Result с ошибкой);
//   - добавлен context для отмены.
package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	ID   int
	Data string
}

type Result struct {
	TaskID int
	Output string
	Err    error
}

func processTask(task Task) (string, error) {
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
	if task.Data == "corrupt" {
		return "", fmt.Errorf("повреждённые данные")
	}
	return "processed " + task.Data, nil
}

func worker(ctx context.Context, jobs <-chan Task, results chan<- Result) {
	for task := range jobs {
		if ctx.Err() != nil {
			return
		}

		func() { // паника одной задачи не должна ронять пул
			defer func() {
				if r := recover(); r != nil {
					results <- Result{TaskID: task.ID, Err: fmt.Errorf("panic: %v", r)}
				}
			}()

			out, err := processTask(task)
			results <- Result{TaskID: task.ID, Output: out, Err: err}
		}()
	}
}

func processAll(ctx context.Context, tasks []Task, workers int) []Result {
	jobs := make(chan Task)
	results := make(chan Result, len(tasks))

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			worker(ctx, jobs, results)
		}()
	}

	go func() {
		defer close(jobs)
		for _, t := range tasks {
			select {
			case jobs <- t:
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var out []Result
	for r := range results {
		out = append(out, r)
	}
	return out
}

func main() {
	tasks := []Task{
		{ID: 1, Data: "image1.jpg"},
		{ID: 2, Data: "image2.jpg"},
		{ID: 3, Data: "corrupt"},
		{ID: 4, Data: "image4.jpg"},
		{ID: 5, Data: "image5.jpg"},
	}

	for _, r := range processAll(context.Background(), tasks, 3) {
		if r.Err != nil {
			fmt.Printf("Task %d failed: %v\n", r.TaskID, r.Err)
			continue
		}
		fmt.Printf("Task %d: %s\n", r.TaskID, r.Output)
	}
}
