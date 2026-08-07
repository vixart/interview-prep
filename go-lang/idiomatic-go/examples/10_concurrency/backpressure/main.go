// Противодавление: буферизованный канал играет роль семафора (limit слотов).
// Если свободного слота нет, select уходит в default и сразу возвращает ошибку —
// сервер отвечает 429 вместо того, чтобы копить очередь.
// Запуск: go run . , затем несколько параллельных curl localhost:8080/request
package main

import (
	"errors"
	"net/http"
	"time"
)

type PressureGauge struct {
	ch chan struct{}
}

func New(limit int) *PressureGauge {
	return &PressureGauge{
		ch: make(chan struct{}, limit),
		// буфер = количество одновременно выполняемых задач
	}
}

func (pg *PressureGauge) Process(f func()) error {
	select {
	case pg.ch <- struct{}{}:
		// занимаем слот; если буфер полон, эта ветвь не готова
		f()
		<-pg.ch
		// освобождаем слот
		return nil
	default:
		return errors.New("no more capacity")
		// это и есть противодавление: не ждем очередь, а сразу отказываем
	}
}

func doThingThatShouldBeLimited() string {
	time.Sleep(2 * time.Second)
	return "done"
}

func main() {
	pg := New(10)
	http.HandleFunc("/request", func(w http.ResponseWriter, r *http.Request) {
		err := pg.Process(func() {
			w.Write([]byte(doThingThatShouldBeLimited()))
		})
		if err != nil {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Too many requests"))
		}
	})
	http.ListenAndServe(":8080", nil)
}
