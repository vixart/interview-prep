// Как правильно синхронизировать несколько горутин — четыре рабочих способа
// и когда какой уместен.
//
//  1. WaitGroup            — «дождаться всех»;
//  2. канал результатов    — «дождаться N значений» и сразу их собрать;
//  3. errgroup своими руками — «первая ошибка отменяет остальных»;
//  4. sync.Once / OnceValue — «сделать ровно один раз».
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// 1. WaitGroup: ждем завершения, результаты не нужны.
func waitAll() {
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1) // Add до go, а не внутри горутины
		go func(id int) {
			defer wg.Done()
			time.Sleep(time.Duration(id) * 10 * time.Millisecond)
		}(i)
	}
	wg.Wait()
	fmt.Println("1. waitAll: все три завершились")
}

// 2. Канал: количество результатов известно заранее, WaitGroup не нужен.
func collectN() {
	const n = 3
	ch := make(chan int, n) // буфер = n, писатели не блокируются
	for i := 1; i <= n; i++ {
		go func(id int) { ch <- id * 10 }(i)
	}
	sum := 0
	for i := 0; i < n; i++ {
		sum += <-ch // ровно n чтений — закрывать канал не обязательно
	}
	fmt.Println("2. collectN: сумма", sum)
}

// 3. «Первая ошибка отменяет остальных» — идея golang.org/x/sync/errgroup
// в двадцати строках. На собеседовании достаточно назвать errgroup,
// но полезно понимать, из чего он собран.
type group struct {
	wg      sync.WaitGroup
	once    sync.Once // ошибку записываем только первую
	err     error
	cancel  context.CancelFunc
	limiter chan struct{} // опционально: ограничение параллелизма
}

func newGroup(ctx context.Context, limit int) (*group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	g := &group{cancel: cancel}
	if limit > 0 {
		g.limiter = make(chan struct{}, limit)
	}
	return g, ctx
}

func (g *group) Go(f func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if g.limiter != nil {
			g.limiter <- struct{}{}        // занять слот
			defer func() { <-g.limiter }() // освободить
		}
		if err := f(); err != nil {
			g.once.Do(func() {
				g.err = err
				g.cancel() // сигнал остальным: можно сворачиваться
			})
		}
	}()
}

func (g *group) Wait() error {
	g.wg.Wait()
	g.cancel() // освобождаем контекст в любом случае
	return g.err
}

func firstErrorCancels() {
	g, ctx := newGroup(context.Background(), 2)
	for i := 1; i <= 4; i++ {
		id := i
		g.Go(func() error {
			select {
			case <-time.After(time.Duration(id) * 20 * time.Millisecond):
				if id == 2 {
					return errors.New("задача 2 упала")
				}
				return nil
			case <-ctx.Done():
				return ctx.Err() // отменили — выходим, а не доделываем
			}
		})
	}
	fmt.Println("3. firstErrorCancels:", g.Wait())
}

// 4. Однократная инициализация из нескольких горутин.
var initOnce = sync.OnceValue(func() string {
	time.Sleep(10 * time.Millisecond)
	return "готово"
})

func onceForAll() {
	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			_ = initOnce() // тело выполнится ровно один раз на всех
		}()
	}
	wg.Wait()
	fmt.Println("4. onceForAll:", initOnce())
}

func main() {
	waitAll()
	collectN()
	firstErrorCancels()
	onceForAll()

	// Чего делать НЕ надо:
	// - time.Sleep вместо синхронизации;
	// - wg.Add внутри горутины (Wait может проскочить раньше);
	// - копировать WaitGroup/Mutex — передавай по указателю;
	// - забывать, что запись в закрытый канал паникует: закрывает писатель.
}
