// Утечки горутин: горутина, которая никогда не завершится, держит свой стек
// (минимум 8 КБ) и все, на что ссылается, — и так до конца жизни процесса.
// GC ее не соберет: работающая горутина всегда достижима.
//
// Четыре типовых источника утечки и как их чинить.
// Считаем горутины через runtime.NumGoroutine().
package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"time"
)

// 1. Читатель ушел, писатель остался блокированным на отправке.
func leakBlockedSend() {
	ch := make(chan int) // небуферизованный: запись ждет читателя
	go func() {
		for i := 0; i < 100; i++ {
			ch <- i // после break в main здесь зависнем НАВСЕГДА
		}
		close(ch)
	}()
	for v := range ch {
		if v > 2 {
			break // читатель ушел, а горутина осталась
		}
	}
}

// Лечение: контекст (или канал done) в select — у горутины появляется выход.
func fixedWithContext() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // обязательно, иначе лечения не случится

	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 0; i < 100; i++ {
			select {
			case ch <- i:
			case <-ctx.Done():
				return
			}
		}
	}()
	for v := range ch {
		if v > 2 {
			break
		}
	}
}

// 2. Забытый Ticker: тикер продолжает слать значения, горутина живет вечно.
func leakTicker() {
	t := time.NewTicker(10 * time.Millisecond)
	// нет defer t.Stop() — и нет способа выйти из цикла
	go func() {
		for range t.C {
			_ = 1
		}
	}()
}

// 3. Незакрытое тело ответа: соединение не возвращается в пул,
// а фоновая горутина транспорта продолжает жить.
func leakHTTPBody(url string) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	_ = resp // забыли defer resp.Body.Close() и не дочитали тело
}

// 4. Ожидание на канале, в который никто уже не напишет.
func leakForgottenWaiter() {
	done := make(chan struct{})
	go func() {
		<-done // отправителя нет — горутина заблокирована навсегда
	}()
}

func count(label string, before int) {
	time.Sleep(50 * time.Millisecond) // даем горутинам доработать
	runtime.GC()
	fmt.Printf("%-26s горутин: %d (было %d)\n", label, runtime.NumGoroutine(), before)
}

func main() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}))
	defer srv.Close()

	base := runtime.NumGoroutine()

	leakBlockedSend()
	count("после leakBlockedSend", base)

	fixedWithContext()
	count("после fixedWithContext", base) // прироста нет

	leakTicker()
	count("после leakTicker", base)

	leakHTTPBody(srv.URL)
	count("после leakHTTPBody", base)

	leakForgottenWaiter()
	count("после leakForgottenWaiter", base)

	// Как ловить в проде:
	// - runtime.NumGoroutine() как метрика: монотонный рост = утечка;
	// - net/http/pprof: /debug/pprof/goroutine?debug=1 покажет стеки и где стоят;
	// - в тестах: сравнить NumGoroutine до и после (goleak от Uber делает это же);
	// - при `fatal error: all goroutines are asleep` Go печатает стеки всех горутин.
	//
	// Правило: запуская горутину, сразу отвечай на вопрос «как она завершится».
}
