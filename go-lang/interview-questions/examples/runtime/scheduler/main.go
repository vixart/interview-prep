// Как работает планировщик Go (модель G-M-P).
//
//	G (goroutine) — сама горутина: стек + контекст выполнения, ~2 КБ на старте;
//	M (machine)   — поток ОС;
//	P (processor) — «право выполнять Go-код», их ровно GOMAXPROCS.
//
// Чтобы выполнить G, потоку M нужен P. У каждого P своя локальная очередь
// горутин (256 штук) плюс есть глобальная очередь. Свободный P ворует половину
// очереди у соседа (work stealing) — так нагрузка размазывается сама.
//
// Переключение (планировщик вытесняющий, начиная с Go 1.14):
//   - блокирующая операция с каналом, мьютексом, сетью, time.Sleep;
//   - системный вызов: M блокируется, P отдается другому M;
//   - вызов функции (кооперативные точки) и асинхронное вытеснение по сигналу
//     примерно раз в 10 мс — поэтому даже пустой цикл не вешает планировщик;
//   - runtime.Gosched() — уступить явно.
//
// Смотреть работу планировщика вживую:
//
//	GODEBUG=schedtrace=1000 go run ./runtime/scheduler
//	GODEBUG=schedtrace=1000,scheddetail=1 go run ./runtime/scheduler
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	fmt.Println("NumCPU     :", runtime.NumCPU())
	fmt.Println("GOMAXPROCS :", runtime.GOMAXPROCS(0), "← столько P, то есть столько горутин РЕАЛЬНО параллельно")
	fmt.Println("NumGoroutine:", runtime.NumGoroutine())

	// 1. Горутины дешевые: миллион штук — это гигабайты? Нет, ~2 КБ на старте.
	const n = 100_000
	var start runtime.MemStats
	runtime.ReadMemStats(&start)

	var wg sync.WaitGroup
	release := make(chan struct{})
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			<-release // держим их живыми, чтобы посчитать память
		}()
	}
	time.Sleep(100 * time.Millisecond)

	var peak runtime.MemStats
	runtime.ReadMemStats(&peak)
	fmt.Printf("\n%d горутин живы, NumGoroutine: %d\n", n, runtime.NumGoroutine())
	fmt.Printf("память на них: %.1f МБ ≈ %.0f байт на горутину\n",
		float64(peak.HeapAlloc-start.HeapAlloc)/1024/1024,
		float64(peak.HeapAlloc-start.HeapAlloc)/n)

	close(release)
	wg.Wait()
	fmt.Println("после завершения NumGoroutine:", runtime.NumGoroutine())

	// 2. Вытеснение: горутина с бесконечным вычислением НЕ блокирует остальных,
	// даже если GOMAXPROCS = 1. До Go 1.14 такой цикл вешал программу.
	old := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(old)

	stop := make(chan struct{})
	done := make(chan string, 2)

	go func() { // «жадная» горутина без единой блокирующей операции
		x := 0
		for {
			select {
			case <-stop:
				done <- "жадная остановлена"
				return
			default:
				x++ // чистое вычисление
			}
		}
	}()

	go func() { // обычная горутина: получает время, хотя P всего один
		time.Sleep(50 * time.Millisecond)
		done <- "обычная горутина отработала при GOMAXPROCS=1"
	}()

	fmt.Println("\n" + <-done)
	close(stop)
	fmt.Println(<-done)

	// Что спрашивают дальше:
	// - «сколько горутин можно создать?» — сотни тысяч; предел не в планировщике,
	//   а в памяти под стеки и в том, что делает сама горутина;
	// - «горутина == поток?» — нет, M:N: тысячи G на десяток M;
	// - «что при syscall?» — M блокируется вместе с G, P уходит другому M,
	//   поэтому блокирующий вызов не останавливает остальные горутины;
	// - «а сетевой вызов?» — не блокирует M: netpoller (epoll/kqueue) паркует G
	//   и будит ее, когда сокет готов;
	// - «зачем GOMAXPROCS менять?» — в контейнере с лимитом CPU
	//   (или брать automaxprocs), иначе рантайм видит все ядра хоста.
}
