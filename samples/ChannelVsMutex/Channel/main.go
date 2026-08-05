package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// scoreboardManager — ЭТО ВЛАДЕЛЕЦ СОСТОЯНИЯ.
//
// ВАЖНО:
// - только эта горутина имеет прямой доступ к map
// - никакая другая горутина не видит scoreboard напрямую
// - вместо mutex используется последовательное выполнение функций
func scoreboardManager(ctx context.Context, in <-chan func(map[string]int)) {

	// Внутреннее состояние.
	// Эта map никогда не покидает эту горутину.
	scoreboard := map[string]int{}

	for {
		select {

		// Если контекст отменён — корректно завершаемся.
		case <-ctx.Done():
			fmt.Println("scoreboard manager stopped")
			return

		// Получаем функцию-команду из канала.
		// Эта функция будет выполнена в контексте
		// текущей (единственной) горутины-владельца.
		case f := <-in:
			f(scoreboard) // безопасно: нет конкурентного доступа
		}
	}
}

// ChannelScoreboardManager — это канал функций.
// Мы не передаём данные.
// Мы передаём ОПЕРАЦИИ над данными.
type ChannelScoreboardManager chan func(map[string]int)

// Конструктор запускает manager-горутины.
// После этого весь доступ к состоянию идёт через канал.
func NewChannelScoreboardManager(ctx context.Context) ChannelScoreboardManager {

	ch := make(ChannelScoreboardManager)

	// Запускаем отдельную горутину-actor.
	go scoreboardManager(ctx, ch)

	return ch
}

// Update — "fire-and-forget" команда.
// Мы не ждём ответа.
func (csm ChannelScoreboardManager) Update(name string, val int) {

	// Отправляем функцию в канал.
	// Эта функция выполнится внутри manager-горутины.
	csm <- func(m map[string]int) {
		m[name] = val
	}

	// ВАЖНО:
	// здесь нет mutex
	// нет гонки
	// потому что m изменяется только в одной горутине
}

// Read — синхронный запрос.
// Здесь уже нужен ответ обратно.
func (csm ChannelScoreboardManager) Read(name string) (int, bool) {

	// Локальный тип для передачи результата
	type Result struct {
		out int
		ok  bool
	}

	// Канал для ответа.
	// Без буфера → вызов будет блокироваться,
	// пока manager не отправит результат.
	resultCh := make(chan Result)

	// Отправляем функцию-чтение.
	csm <- func(m map[string]int) {

		// Чтение выполняется в manager-горутины.
		out, ok := m[name]

		// Возвращаем результат через канал ответа.
		resultCh <- Result{out, ok}
	}

	// Блокируемся и ждём результат.
	result := <-resultCh

	return result.out, result.ok
}

func main() {

	// Контекст позволяет корректно завершить manager.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	manager := NewChannelScoreboardManager(ctx)

	var wg sync.WaitGroup

	// Демонстрация конкурентных обновлений.
	names := []string{"Alice", "Bob", "Charlie"}

	for i, name := range names {
		wg.Add(1)

		// Запускаем несколько горутин одновременно.
		go func(n string, score int) {
			defer wg.Done()

			// Эти вызовы конкурентные,
			// но внутри manager они выполнятся последовательно.
			manager.Update(n, score)

		}(name, (i+1)*10)
	}

	wg.Wait()

	// Теперь читаем результаты.
	// Чтение тоже безопасно.
	for _, name := range names {

		score, ok := manager.Read(name)

		fmt.Printf("%s -> %d (exists: %v)\n", name, score, ok)
	}

	// Небольшая задержка перед завершением программы,
	// чтобы увидеть сообщение остановки.
	time.Sleep(200 * time.Millisecond)
}
