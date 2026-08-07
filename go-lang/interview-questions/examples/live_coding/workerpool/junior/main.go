// Демо junior-пула: счастливый путь, на котором баги не видны.
//
//	go run ./live_coding/workerpool/junior         # отработает нормально
//	go test  ./live_coding/workerpool/junior       # тест ловит панику send on closed channel
//	go test -race ./live_coding/workerpool/junior  # детектор гонок падает — так и задумано
package main

import (
	"log/slog"
	"time"
)

func main() {
	pool := New(3)

	for i := range 5 {
		// Задача — это func() без параметров: ни контекста, ни возврата ошибки.
		// Отсюда два ограничения junior-версии: задачу нельзя отменить
		// и нельзя узнать, чем она закончилась.
		pool.AddTask(func() {
			slog.Info("задача запущена", "id", i)
			<-time.After(200 * time.Millisecond)
			slog.Info("задача готова", "id", i)
		})
	}

	// Close дожидается всех задач. Здесь все хорошо ровно потому, что
	// AddTask и Close вызываются из ОДНОЙ горутины: гонка из БАГ #3
	// не проявляется. В реальном сервисе Submit зовут хендлеры,
	// а Close — обработчик SIGTERM, то есть параллельно.
	pool.Close()
	slog.Info("все задачи завершены")
}
