// Client-side балансировщик: Round Robin + retry на следующем экземпляре.
//
//   - Balancer сам реализует Backend — его можно подставить куда угодно;
//   - счётчик Round Robin — atomic.Uint64, без мьютекса;
//   - при ошибке пробуем следующие экземпляры (не больше len(backends)
//     попыток), уважая отмену контекста;
//   - пустой список бэкендов — явная ошибка.
package main

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
)

type Request any
type Response any

type Backend interface {
	Invoke(ctx context.Context, req Request) (Response, error)
}

type Balancer struct {
	backends []Backend
	next     atomic.Uint64
}

var _ Backend = (*Balancer)(nil)

func NewBalancer(backends []Backend) *Balancer {
	return &Balancer{backends: backends}
}

func (b *Balancer) Invoke(ctx context.Context, req Request) (Response, error) {
	n := len(b.backends)
	if n == 0 {
		return nil, errors.New("balancer: no backends")
	}

	var lastErr error
	for attempt := 0; attempt < n; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		// Add(1)-1: каждый вызов получает свой уникальный номер,
		// поэтому нагрузка распределяется равномерно даже из горутин.
		idx := (b.next.Add(1) - 1) % uint64(n)

		resp, err := b.backends[idx].Invoke(ctx, req)
		if err == nil {
			return resp, nil
		}
		lastErr = err // экземпляр недоступен — пробуем следующий
	}

	return nil, fmt.Errorf("balancer: all backends failed: %w", lastErr)
}

// --- демо ---

type BackendImpl struct {
	addr string
	fail bool
}

func (b *BackendImpl) Invoke(_ context.Context, req Request) (Response, error) {
	if b.fail {
		return nil, fmt.Errorf("%s: unavailable", b.addr)
	}
	return fmt.Sprintf("%v -> handled by %s", req, b.addr), nil
}

func main() {
	lb := NewBalancer([]Backend{
		&BackendImpl{addr: "10.0.0.1"},
		&BackendImpl{addr: "10.0.0.2", fail: true}, // «упавший» экземпляр
		&BackendImpl{addr: "10.0.0.3"},
	})

	ctx := context.Background()
	for i := 1; i <= 6; i++ {
		resp, err := lb.Invoke(ctx, fmt.Sprintf("req-%d", i))
		fmt.Println(resp, err)
	}
	// Запросы распределяются по живым экземплярам; попадание на .2
	// прозрачно ретраится на следующем.
}
