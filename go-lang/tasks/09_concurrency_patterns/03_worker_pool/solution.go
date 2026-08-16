// Worker Pool: параллельная обработка изображений ограниченным числом горутин.
//
// Канонический каркас:
//   - фиксированные N воркеров читают задачи из jobs через for range;
//   - продюсер закрывает jobs после отправки всех задач (сигнал «конец»);
//   - горутина-закрыватель ждёт wg.Wait() и закрывает results;
//   - потребитель читает results через for range.
package main

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Image struct{ Name string }

type Result struct {
	Name string
	Err  error
}

func process(img Image) Result {
	time.Sleep(50 * time.Millisecond) // имитация наложения водяного знака
	if strings.Contains(img.Name, "corrupt") {
		return Result{Name: img.Name, Err: fmt.Errorf("битый файл")}
	}
	return Result{Name: img.Name + ".watermarked"}
}

func processImages(images []Image, workers int) []Result {
	jobs := make(chan Image)
	results := make(chan Result, len(images))

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()

			for img := range jobs { // выходим, когда jobs закрыт и пуст
				results <- process(img)
			}
		}()
	}

	go func() { // продюсер
		defer close(jobs)
		for _, img := range images {
			jobs <- img
		}
	}()

	go func() { // закрыватель results
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
	images := []Image{
		{"img1.jpg"}, {"img2.jpg"}, {"corrupt.jpg"}, {"img4.jpg"},
		{"img5.jpg"}, {"img6.jpg"}, {"img7.jpg"}, {"img8.jpg"},
	}

	start := time.Now()
	results := processImages(images, runtime.NumCPU())

	for _, r := range results {
		if r.Err != nil {
			fmt.Println(r.Name, "ошибка:", r.Err)
			continue
		}
		fmt.Println(r.Name)
	}
	fmt.Println("обработано за", time.Since(start).Round(10*time.Millisecond))
	// 8 задач по 50мс на N ядрах — заметно быстрее последовательных 400мс.
}
